import { NextRequest, NextResponse } from 'next/server';
import { spawn } from 'child_process';
import path from 'path';
import fs from 'fs';
import { getExperiment } from '@/lib/api';

export async function POST(
    req: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id: idStr } = await params;
        const id = parseInt(idStr);
        if (isNaN(id)) {
            return NextResponse.json({ error: 'Invalid ID' }, { status: 400 });
        }

        const exp = await getExperiment(id);
        if (!exp) {
            return NextResponse.json({ error: 'Experiment not found' }, { status: 404 });
        }

        // Determine experiment directory from the result path stored in DB
        // resultsPath is typically .../experiments/runs/<folder>/results.json
        const resultsPath = exp.results_path;
        const experimentDir = path.dirname(resultsPath);

        const rootDir = path.resolve(process.cwd(), '..');
        const binaryPath = path.join(rootDir, 'tenkai');

        if (!fs.existsSync(binaryPath)) {
            console.error(`Tenkai binary not found at: ${binaryPath}`);
            return NextResponse.json({ error: 'Backend binary not found' }, { status: 500 });
        }

        // Trigger pre-built Tenkai binary in report-only mode
        // Use absolute path for command
        const cmd = binaryPath;
        const args = ['-report-only', experimentDir];

        console.log(`Executing: ${cmd} ${args.join(' ')}`);

        // We run synchronously-ish or detach? 
        // Report generation is fast, so we can probably wait for it to ensure it's done before returning.

        return new Promise<NextResponse>((resolve, reject) => {
            const child = spawn(cmd, args, {
                cwd: rootDir
            });

            let stderrData = '';
            let stdoutData = '';

            child.stdout.on('data', (data) => {
                stdoutData += data.toString();
                console.log(`[Tenkai CLI] ${data}`);
            });

            child.stderr.on('data', (data) => {
                stderrData += data.toString();
                console.error(`[Tenkai CLI Error] ${data}`);
            });

            child.on('close', (code) => {
                if (code === 0) {
                    resolve(NextResponse.json({ message: 'Report regenerated successfully' }));
                } else {
                    console.error(`Regeneration failed (code ${code}). Stderr: ${stderrData}`);
                    resolve(NextResponse.json({ error: `Process exited with code ${code}: ${stderrData}` }, { status: 500 }));
                }
            });

            child.on('error', (err) => {
                console.error('Spawn error:', err);
                resolve(NextResponse.json({ error: 'Failed to spawn process' }, { status: 500 }));
            });
        });

    } catch (error) {
        console.error('Failed to regenerate report:', error);
        return NextResponse.json({ error: 'Failed to regenerate report' }, { status: 500 });
    }
}
