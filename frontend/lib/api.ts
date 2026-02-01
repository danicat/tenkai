
import { ConfigBlock, ExperimentRecord, RunResult, ExperimentSummaryRow, Checkpoint } from "@/types/domain";

const API_BASE = '/api';

// Re-export or alias types to match legacy API client usage
export type { ExperimentRecord, Checkpoint };
export type RunResultRecord = RunResult;
export type ExperimentSummaryRecord = ExperimentSummaryRow;

// Types NOT in domain.ts

export interface ToolUsageRecord {
    id?: number;
    run_id?: number;
    name: string; // UI expects 'name', DB might be 'tool_name' but we map it? Or UI uses 'name'.
    status: string;
    duration: number;
    args?: any;
    output?: string;
    error?: string;
    count?: number; // legacy?
    success_count?: number; // legacy?
    fail_count?: number; // legacy?
}

export interface MessageRecord {
    id: number;
    run_id: number;
    role: string;
    content: string;
    timestamp: string;
}

export interface TestResultRecord {
    name: string;
    status: string; // PASS/FAIL
    duration: number;
    output?: string;
}

export interface LintResultRecord {
    file: string;
    line: number;
    message: string;
    code: string;
}

export interface ToolStatRow {
    alternative: string;
    tool_name: string;
    total_calls: number;
    failed_calls: number; // needed for ToolUsageTable
    avg_calls_per_run: number;
    avg_cost: number;
    failure_rate: number;
}

// --- API Functions ---

async function fetchJSON(endpoint: string, options?: RequestInit) {
    const res = await fetch(`${API_BASE}${endpoint}`, options);
    if (!res.ok) {
        throw new Error(`API Error ${res.status}: ${await res.text()}`);
    }
    return res.json();
}

export async function getExperiments(): Promise<ExperimentRecord[]> {
    return fetchJSON('/experiments');
}

export async function getExperiment(id: number): Promise<ExperimentRecord> {
    return fetchJSON(`/experiments/${id}`);
}

export async function getRunResults(experimentId: number | string, page: number = 1, limit: number = 5000): Promise<RunResultRecord[]> {
    return fetchJSON(`/experiments/${experimentId}/runs?page=${page}&limit=${limit}`);
}

export async function reValidateExperiment(experimentId: number): Promise<{ job_id?: string; error?: string }> {
    return fetchJSON(`/experiments/${experimentId}/revalidate`, { method: 'POST' });
}

export async function getExperimentSummaries(experimentId: number, filter: string = 'all'): Promise<ExperimentSummaryRecord[]> {
    return fetchJSON(`/experiments/${experimentId}/stats?filter=${filter}`);
}

export async function getSimplifiedMetrics(experimentId: number | string) {
    return fetchJSON(`/experiments/${experimentId}/stats`);
}

export async function getCheckpoint(experimentId: string): Promise<Checkpoint | null> {
    try {
        return await fetchJSON(`/experiments/${experimentId}/checkpoint`);
    } catch { return null; }
}

export async function getBlocks(): Promise<ConfigBlock[]> {
    return fetchJSON('/blocks');
}

export async function getScenarios() {
    return fetchJSON('/scenarios');
}

export async function getTemplates() {
    return fetchJSON('/templates');
}

export async function toggleLock(experimentId: number, locked: boolean) {
    return fetchJSON(`/experiments/${experimentId}/lock`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ locked })
    });
}

export async function toggleTemplateLock(name: string, locked: boolean) {
    return fetchJSON(`/templates/${name}/lock`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ locked })
    });
}

export async function toggleScenarioLock(id: string, locked: boolean) {
    return fetchJSON(`/scenarios/${id}/lock`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ locked })
    });
}

export async function saveExperimentAnnotations(id: number, notes: string) {
    return fetchJSON(`/experiments/${id}/annotations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ notes })
    });
}

export async function reValidateRun(runId: number) {
    return fetchJSON(`/runs/${runId}/revalidate`, { method: 'POST' });
}
