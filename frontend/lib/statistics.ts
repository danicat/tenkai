export interface RunResult {
    alternative: string;
    scenario: string;
    duration: number; // in nanoseconds
    error: string;
    status: string;
    lint_issues?: number;
    tests_passed?: number;
    tests_failed?: number;
    is_success?: boolean;
    validation_report?: string;
    total_tokens?: number;
    agent_metrics?: {
        total_tokens: number;
        tool_calls: number;
    };
}

export interface StatGroup {
    alternative: string;
    count: number;
    successCount: number;
    successRate: number;
    avgDuration: number; // seconds
    avgTokens: number;
    avgLint: number;
    totalLint: number;
    durations: number[]; // seconds
    successes: number[]; // 1 or 0
}

export interface SignificanceResult {
    isSignificant: boolean;
    p: number;
    level: string; // *, **, ***
    better: boolean;
}

// Group runs by alternative
export function calculateStats(runs: RunResult[]): Record<string, StatGroup> {
    const groups: Record<string, StatGroup> = {};

    runs.forEach(run => {
        if (!groups[run.alternative]) {
            groups[run.alternative] = {
                alternative: run.alternative,
                count: 0,
                successCount: 0,
                successRate: 0,
                avgDuration: 0,
                avgTokens: 0,
                avgLint: 0,
                totalLint: 0,
                durations: [],
                successes: []
            };
        }

        const g = groups[run.alternative];
        g.count++;

        // A run is successful if is_success is true (new validation system) 
        // OR if no error AND at least one test passed (legacy eval system)
        const isSuccess = (run.validation_report && run.validation_report !== "")
            ? !!run.is_success
            : (!run.error && (run.tests_passed || 0) > 0);
        if (isSuccess) {
            g.successCount++;
            g.successes.push(1);
            g.durations.push(run.duration / 1e9); // Convert to seconds

            if (run.total_tokens !== undefined) {
                g.avgTokens += run.total_tokens;
            } else if (run.agent_metrics) {
                g.avgTokens += run.agent_metrics.total_tokens;
            }
        } else {
            g.successes.push(0);
        }

        if (run.lint_issues) {
            g.totalLint += run.lint_issues;
        }
    });

    // Final averages
    Object.values(groups).forEach(g => {
        g.successRate = (g.successCount / g.count) * 100;
        if (g.successCount > 0) {
            g.avgDuration = g.durations.reduce((a, b) => a + b, 0) / g.successCount;
            g.avgTokens = g.avgTokens / g.successCount;
        }
        g.avgLint = g.totalLint / g.count; // Lint is averaged over ALL runs typically
    });

    return groups;
}

// Z-Test for Two Proportions (Success Rate)
export function zTestProportions(p1: number, n1: number, p2: number, n2: number): SignificanceResult {
    if (n1 === 0 || n2 === 0) return { isSignificant: false, p: 1, level: '', better: false };

    // Pooled sample proportion
    const p = (p1 * n1 + p2 * n2) / (n1 + n2);

    // Standard error
    const se = Math.sqrt(p * (1 - p) * (1 / n1 + 1 / n2));

    if (se === 0) return { isSignificant: false, p: 1, level: '', better: false };

    const z = (p1 - p2) / se;

    // Two-tailed p-value
    // Approximate normal CDF
    const pVal = 2 * (1 - normalCDF(Math.abs(z)));

    return {
        isSignificant: pVal < 0.05,
        p: pVal,
        level: getSignificanceLevel(pVal),
        better: p1 > p2 // Only if significant
    };
}

// Welch's t-test for Means (Duration)
export function welchTTest(sample1: number[], sample2: number[]): SignificanceResult {
    const n1 = sample1.length;
    const n2 = sample2.length;

    if (n1 < 2 || n2 < 2) return { isSignificant: false, p: 1, level: '', better: false };

    const m1 = mean(sample1);
    const m2 = mean(sample2);
    const v1 = variance(sample1, m1);
    const v2 = variance(sample2, m2);

    if (v1 === 0 && v2 === 0) return { isSignificant: false, p: 1, level: '', better: false };

    const se = Math.sqrt(v1 / n1 + v2 / n2);
    const t = (m1 - m2) / se;

    // Degrees of freedom (Welch-Satterthwaite)
    const df = Math.pow(v1 / n1 + v2 / n2, 2) / (Math.pow(v1 / n1, 2) / (n1 - 1) + Math.pow(v2 / n2, 2) / (n2 - 1));

    // P-value implementation (approximated via Normal for now as t-dist is complex to implement fully in 1 file without libs)
    // For df > 30, t -> normal. For small samples, this is an approximation.
    const pVal = 2 * (1 - normalCDF(Math.abs(t))); // Two-tailed

    return {
        isSignificant: pVal < 0.05,
        p: pVal,
        level: getSignificanceLevel(pVal),
        better: m1 < m2 // Lower duration is better
    };
}

// Helpers
function mean(arr: number[]): number {
    return arr.reduce((a, b) => a + b, 0) / arr.length;
}

function variance(arr: number[], mean: number): number {
    return arr.reduce((a, b) => a + Math.pow(b - mean, 2), 0) / (arr.length - 1);
}

function normalCDF(x: number): number {
    // Error function approximation
    const t = 1 / (1 + 0.2316419 * x);
    const d = 0.3989423 * Math.exp(-x * x / 2);
    const prob = d * t * (0.3193815 + t * (-0.3565638 + t * (1.781478 + t * (-1.821256 + t * 1.330274))));
    return 1 - prob;
}

function getSignificanceLevel(p: number): string {
    if (p < 0.001) return '***';
    if (p < 0.01) return '**';
    if (p < 0.05) return '*';
    return '';
}
