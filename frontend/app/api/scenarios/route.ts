import { NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';

export async function GET() {
    try {
        const scenariosDir = path.join(process.cwd(), '..', 'scenarios');
        const entries = fs.readdirSync(scenariosDir, { withFileTypes: true });

        const scenarios = entries
            .filter(entry => entry.isDirectory())
            .map(entry => entry.name);

        return NextResponse.json(scenarios);
    } catch (error) {
        return NextResponse.json({ error: 'Failed to list studies' }, { status: 500 });
    }
}
