import Database from 'better-sqlite3';
import path from 'path';

const DB_PATH = path.join(process.cwd(), '..', 'experiments', 'tenkai.db');

export interface ExperimentRecord {
    id: number;
    name: string;
    timestamp: string;
    config_path: string;
    results_path: string;
    report_path: string;
    status: 'running' | 'completed' | 'failed' | 'stopped' | 'paused';
    reps: number;
    concurrent: number;
    description?: string;
    duration_total?: number;
    success_rate?: string;
    config_content: string;
    report_content: string;
    progress?: {
        completed: number;
        total: number;
        percentage: number;
    };
}

export interface Checkpoint {
    total_jobs: number;
    completed_jobs: number;
    percentage: number;
    last_update: string;
    status: string;
}

export interface ToolUsageRecord {
    id: number;
    run_id: number;
    name: string;
    args: string;
    status: string;
    output: string;
    error: string;
    duration: number; // seconds
    timestamp: string;
}

export interface MessageRecord {
    id: number;
    run_id: number;
    role: string;
    content: string;
    timestamp: string;
}

export interface TestResultRecord {
    id: number;
    run_id: number;
    name: string;
    status: string;
    duration_ns: number;
    output: string;
}

export interface LintResultRecord {
    id: number;
    run_id: number;
    file: string;
    line: number;
    col: number;
    message: string;
    severity: string;
    rule_id: string;
}

function getDB() {
    const db = new Database(DB_PATH);
    db.pragma('journal_mode = WAL');
    return db;
}

function safeDate(dateStr: string): string {
    try {
        const cleanDate = dateStr.includes(' m=+') ? dateStr.split(' m=+')[0] : dateStr;
        const date = new Date(cleanDate);
        return isNaN(date.getTime()) ? new Date().toISOString() : date.toISOString();
    } catch (e) {
        return new Date().toISOString();
    }
}

export async function getExperiments(): Promise<ExperimentRecord[]> {
    const db = getDB();
    try {
        const rows = db.prepare(`
            SELECT 
                e.*,
                (SELECT COUNT(*) FROM run_results WHERE experiment_id = e.id AND error = '' AND tests_passed > 0) as success_count,
                (SELECT COUNT(*) FROM run_results WHERE experiment_id = e.id) as total_runs
            FROM experiments e 
            ORDER BY timestamp DESC
        `).all();

        const records = rows.map((row: any) => {
            const successRate = row.total_runs > 0
                ? ((row.success_count / row.total_runs) * 100).toFixed(1) + '%'
                : 'N/A';

            const record: ExperimentRecord = {
                ...row,
                timestamp: safeDate(row.timestamp),
                description: row.description,
                duration_total: row.duration,
                success_rate: successRate,
                config_content: row.config_content,
                report_content: row.report_content
            };

            // Use DB columns
            if (row.total_jobs > 0) {
                record.progress = {
                    completed: row.completed_jobs,
                    total: row.total_jobs,
                    percentage: (row.completed_jobs / row.total_jobs) * 100
                };
            }
            return record;
        }) as ExperimentRecord[];

        return records;
    } catch (error) {
        console.error('Failed to get experiments from DB:', error);
        return [];
    } finally {
        db.close();
    }
}

export async function getExperimentById(id: number | string): Promise<ExperimentRecord | null> {
    const db = getDB();
    try {
        const row = db.prepare('SELECT * FROM experiments WHERE id = ?').get(id) as any;
        if (!row) return null;

        const record: ExperimentRecord = {
            ...row,
            timestamp: safeDate(row.timestamp),
            config_content: row.config_content,
            report_content: row.report_content
        };

        if (row.total_jobs > 0) {
            record.progress = {
                completed: row.completed_jobs,
                total: row.total_jobs,
                percentage: (row.completed_jobs / row.total_jobs) * 100
            };
        }
        return record;
    } finally {
        db.close();
    }
}

export async function getGlobalStats() {
    const db = getDB();
    try {
        const stats = db.prepare(`
            SELECT 
                COUNT(*) as totalJobs,
                AVG(CASE WHEN error = '' OR error IS NULL THEN 1 ELSE 0 END) * 100 as successRate
            FROM run_results
        `).get() as any;

        return {
            totalRuns: stats.totalJobs || 0,
            avgSuccessRate: (stats.successRate || 0).toFixed(1) + '%'
        };
    } finally {
        db.close();
    }
}

export async function getReportContent(reportPath: string): Promise<string> {
    return ""; // Deprecated, use experiment.report_content
}

export async function getCheckpoint(idOrPath: string): Promise<Checkpoint | null> {
    const db = getDB();
    try {
        let row: any;
        if (!isNaN(parseInt(idOrPath)) && !idOrPath.includes('/') && !idOrPath.includes('.')) {
            row = db.prepare('SELECT total_jobs, completed_jobs, status FROM experiments WHERE id = ?').get(idOrPath);
        } else {
            // Fallback to name if not an ID (though we prefer IDs now)
            row = db.prepare('SELECT total_jobs, completed_jobs, status FROM experiments WHERE name = ?').get(idOrPath);
        }

        if (!row) return null;

        const total = row.total_jobs || 1; // avoid division by zero
        return {
            total_jobs: row.total_jobs || 0,
            completed_jobs: row.completed_jobs || 0,
            percentage: ((row.completed_jobs || 0) / total) * 100,
            last_update: new Date().toISOString(), // DB doesn't track update time per row yet, but that's fine for polling
            status: row.status
        };
    } catch (error) {
        return null;
    } finally {
        db.close();
    }
}
// Fetch a single experiment by ID
export function getExperiment(id: number): ExperimentRecord | null {
    const dbPath = process.env.DB_PATH || path.join(process.cwd(), '..', 'experiments', 'tenkai.db');
    const db = new Database(dbPath);
    const stmt = db.prepare(`
        SELECT 
            id, 
            timestamp, 
            name, 
            config_path, 
            results_path, 
            report_path, 
            status,
            reps,
            concurrent,
            total_jobs,
            completed_jobs
        FROM experiments 
        WHERE id = ?
    `);

    const row = stmt.get(id) as any;
    db.close();

    if (!row) return null;

    return {
        id: row.id,
        timestamp: row.timestamp,
        name: row.name,
        config_path: row.config_path,
        results_path: row.results_path,
        report_path: row.report_path,
        status: row.status,
        reps: row.reps,
        concurrent: row.concurrent,
        progress: {
            total: row.total_jobs || 0,
            completed: row.completed_jobs || 0,
            percentage: row.total_jobs > 0 ? (row.completed_jobs / row.total_jobs) * 100 : 0
        }
    };
}

export async function getSimplifiedMetrics(id: number | string) {
    const db = getDB();
    try {
        const results = db.prepare('SELECT * FROM run_results WHERE experiment_id = ?').all(String(id)) as any[];

        if (results.length === 0) return null;

        const total = results.length;
        const successful = results.filter(r => r.error === '' && r.tests_passed > 0).length;
        const totalDuration = results.reduce((acc, r) => acc + (r.duration || 0), 0);
        const totalLint = results.reduce((acc, r) => acc + (r.lint_issues || 0), 0);

        return {
            total,
            successful,
            successRate: (successful / total * 100).toFixed(1) + '%',
            avgDuration: (totalDuration / total / 1e9).toFixed(1) + 's',
            totalLint,
            results: results.map(r => ({ ...r, duration: r.duration / 1e9 }))
        };
    } finally {
        db.close();
    }
}
export async function getRunResults(experimentId: number | string): Promise<any[]> {
    const db = getDB();
    const id = Number(experimentId);
    try {
        const rows = db.prepare('SELECT * FROM run_results WHERE experiment_id = ?').all(id) as any[];

        return rows.map(r => ({
            id: r.id,
            alternative: r.alternative,
            scenario: r.scenario,
            duration: r.duration,
            error: r.error,
            status: r.tests_passed > 0 ? 'SUCCESS' : (r.error ? 'ERROR' : 'FAIL'),
            lint_issues: r.lint_issues,
            tests_passed: r.tests_passed,
            tests_failed: r.tests_failed,
            agent_metrics: {
                total_tokens: r.total_tokens || 0,
                input_tokens: r.input_tokens || 0,
                output_tokens: r.output_tokens || 0,
                tool_calls_count: r.tool_calls_count || 0,
                loop_detected: Boolean(r.loop_detected)
            },
            stdout: r.stdout,
            stderr: r.stderr
        }));
    } finally {
        db.close();
    }
}

export async function getToolUsage(runId: number): Promise<ToolUsageRecord[]> {
    const db = getDB();
    try {
        const rows = db.prepare('SELECT * FROM tool_usage WHERE run_id = ? ORDER BY timestamp ASC').all(runId) as any[];
        return rows.map(r => ({
            ...r,
            duration: (r.duration || 0) / 1e9,
            timestamp: safeDate(r.timestamp)
        }));
    } finally {
        db.close();
    }
}

export async function getMessages(runId: number): Promise<MessageRecord[]> {
    const db = getDB();
    try {
        const rows = db.prepare('SELECT * FROM messages WHERE run_id = ? ORDER BY timestamp ASC').all(runId) as any[];
        return rows.map(r => ({
            ...r,
            timestamp: safeDate(r.timestamp)
        }));
    } finally {
        db.close();
    }
}

export async function getRunFiles(runId: number) {
    const db = getDB();
    try {
        return db.prepare('SELECT * FROM run_files WHERE run_id = ?').all(runId) as any[];
    } finally {
        db.close();
    }
}

export async function getTestResults(runId: number): Promise<TestResultRecord[]> {
    const db = getDB();
    try {
        return db.prepare('SELECT * FROM test_results WHERE run_id = ?').all(runId) as TestResultRecord[];
    } finally {
        db.close();
    }
}

export async function getLintResults(runId: number): Promise<LintResultRecord[]> {
    const db = getDB();
    try {
        return db.prepare('SELECT * FROM lint_results WHERE run_id = ?').all(runId) as LintResultRecord[];
    } finally {
        db.close();
    }
}

export async function getNextExperimentNumber(nameBase: string): Promise<number> {
    const db = getDB();
    try {
        // Look for experiments starting with the nameBase
        const rows = db.prepare('SELECT name FROM experiments WHERE name LIKE ?').all(nameBase + '%') as any[];
        let maxNum = 0;

        const regex = new RegExp(`^${nameBase}_(\\d+)$`);
        for (const row of rows) {
            const match = row.name.match(regex);
            if (match) {
                const num = parseInt(match[1]);
                if (num > maxNum) maxNum = num;
            }
        }
        return maxNum + 1;
    } finally {
        db.close();
    }
}
