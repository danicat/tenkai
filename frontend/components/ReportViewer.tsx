'use client';

import { useState, useEffect } from 'react';
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { Checkpoint, ExperimentRecord } from '@/lib/api';
import { welchTTest, zTestProportions } from "@/lib/statistics";

interface ReportViewerProps {
    experiment: ExperimentRecord;
    initialContent: string;
    initialMetrics: any;
    initialCheckpoint: Checkpoint | null;
    runResults: any[];
    stats: any;
    config: any;
}

export default function ReportViewer({ experiment, initialContent, initialMetrics, initialCheckpoint, runResults, stats, config }: ReportViewerProps) {
    const [tab, setTab] = useState<'analytics' | 'history'>('analytics');
    const [checkpoint, setCheckpoint] = useState<Checkpoint | null>(initialCheckpoint);
    const [status, setStatus] = useState(experiment.status);
    const [metrics, setMetrics] = useState(initialMetrics);
    const [loadingAction, setLoadingAction] = useState<string | null>(null);
    const [selectedRun, setSelectedRun] = useState<any>(null);
    const [runDetails, setRunDetails] = useState<{ tools: any[], messages: any[], files: any[], tests: any[], lints: any[] } | null>(null);
    const [selectedFile, setSelectedFile] = useState<any>(null);
    const [loadingDetails, setLoadingDetails] = useState(false);
    const [inspectingTool, setInspectingTool] = useState<any>(null);

    // Derived Logic
    const altNames = Object.keys(stats);
    const referenceAlt = config?.alternatives?.[0]?.name || altNames[0];
    const refStats = stats[referenceAlt];

    // Sync state with props when server refreshes
    useEffect(() => { setMetrics(initialMetrics); }, [initialMetrics]);
    useEffect(() => { setStatus(experiment.status); }, [experiment.status]);

    useEffect(() => {
        if (status === 'running' || status === 'paused') {
            const interval = setInterval(async () => {
                try {
                    const res = await fetch(`/api/experiments/checkpoint?id=${experiment.id}`);
                    if (res.ok) {
                        const data = await res.json();
                        setCheckpoint(data);
                        if (data.status !== status) {
                            setStatus(data.status);
                        }

                        // Always refresh metrics when running to show live simplified stats (duration, lint, etc.)
                        const mRes = await fetch(`/api/experiments/metrics?id=${experiment.id}`);
                        if (mRes.ok) setMetrics(await mRes.json());

                        if (data.status === 'completed' || data.status === 'stopped' || data.status === 'failed') {
                            clearInterval(interval);
                        }
                    }
                } catch (e) {
                    console.error("Failed to poll checkpoint", e);
                }
            }, 2000);
            return () => clearInterval(interval);
        }
    }, [status, experiment.results_path, experiment.id]);

    const handleControl = async (action: 'pause' | 'resume' | 'stop') => {
        setLoadingAction(action);
        try {
            const res = await fetch('/api/experiments/control', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: experiment.id, action }),
            });
            if (res.ok) {
                if (action === 'stop') setStatus('stopped');
                if (action === 'pause') setStatus('paused');
                if (action === 'resume') setStatus('running');
            } else {
                alert(`Failed to ${action} experiment`);
            }
        } catch (e) {
            console.error(e);
            alert(`Error sending ${action} command`);
        } finally {
            setLoadingAction(null);
        }
    };

    const fetchRunDetails = async (run: any) => {
        setSelectedRun(run);
        setLoadingDetails(true);
        try {
            const [tRes, mRes, fRes, testRes, lRes] = await Promise.all([
                fetch(`/api/runs/${run.id}/tools`),
                fetch(`/api/runs/${run.id}/messages`),
                fetch(`/api/runs/${run.id}/files`),
                fetch(`/api/runs/${run.id}/tests`),
                fetch(`/api/runs/${run.id}/lint`)
            ]);
            if (tRes.ok && mRes.ok && fRes.ok && testRes.ok && lRes.ok) {
                const tools = await tRes.json();
                const messages = await mRes.json();
                const files = await fRes.json();
                const tests = await testRes.json();
                const lints = await lRes.json();
                setRunDetails({ tools, messages, files, tests, lints });
                if (files.length > 0) setSelectedFile(files[0]);
            }
        } catch (e) {
            console.error("Failed to fetch run details", e);
        } finally {
            setLoadingDetails(false);
        }
    };

    const handleRegenerate = async () => {
        setLoadingAction('regenerate');
        try {
            const res = await fetch(`/api/experiments/${experiment.id}/regenerate`, { method: 'POST' });
            if (res.ok) {
                window.location.reload();
            } else {
                alert('Failed to regenerate report');
            }
        } catch (e) {
            console.error(e);
            alert('Error regenerating report');
        } finally {
            setLoadingAction(null);
        }
    };

    const isActive = status === 'running' || status === 'paused';

    return (
        <div className="space-y-6">
            {/* Status Banner */}
            {isActive && (
                <div className={`glass p-6 rounded-2xl border transition-all ${status === 'running' ? 'border-blue-500/20 bg-blue-500/5 animate-pulse' : 'border-amber-500/20 bg-amber-500/5'}`}>
                    <div className="flex justify-between items-center mb-4">
                        <div className="flex items-center gap-4">
                            <h3 className={`font-bold flex items-center gap-2 ${status === 'running' ? 'text-blue-400' : 'text-amber-400'}`}>
                                <span className={`w-2 h-2 rounded-full ${status === 'running' ? 'bg-blue-500 animate-ping' : 'bg-amber-500'}`}></span>
                                {status === 'running' ? 'Experiment in Progress...' : 'Experiment Paused'}
                            </h3>

                            {/* Controls */}
                            <div className="flex gap-2">
                                {status === 'running' ? (
                                    <button
                                        disabled={loadingAction !== null}
                                        onClick={() => handleControl('pause')}
                                        className="bg-amber-500/10 hover:bg-amber-500/20 text-amber-500 text-xs font-bold px-3 py-1 rounded border border-amber-500/20 transition-all uppercase tracking-widest"
                                    >
                                        {loadingAction === 'pause' ? '...' : 'Pause'}
                                    </button>
                                ) : (
                                    <button
                                        disabled={loadingAction !== null}
                                        onClick={() => handleControl('resume')}
                                        className="bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-500 text-xs font-bold px-3 py-1 rounded border border-emerald-500/20 transition-all uppercase tracking-widest"
                                    >
                                        {loadingAction === 'resume' ? '...' : 'Resume'}
                                    </button>
                                )}
                                <button
                                    disabled={loadingAction !== null}
                                    onClick={() => handleControl('stop')}
                                    className="bg-red-500/10 hover:bg-red-500/20 text-red-500 text-xs font-bold px-3 py-1 rounded border border-red-500/20 transition-all uppercase tracking-widest"
                                >
                                    {loadingAction === 'stop' ? '...' : 'Stop'}
                                </button>
                            </div>
                        </div>
                        <span className="text-base font-mono">{checkpoint?.percentage.toFixed(1) || 0}%</span>
                    </div>
                    <div className="w-full h-2 bg-white/5 rounded-full overflow-hidden">
                        <div
                            className={`h-full transition-all duration-500 ${status === 'running' ? 'bg-blue-500' : 'bg-amber-500'}`}
                            style={{ width: `${checkpoint?.percentage || 0}%` }}
                        ></div>
                    </div>
                    <p className="text-xs text-zinc-500 mt-2 uppercase tracking-widest font-bold">
                        Job {checkpoint?.completed_jobs || 0} of {checkpoint?.total_jobs || 0} completed
                    </p>
                </div>
            )}

            {/* Tabs & Actions */}
            <div className="flex justify-between items-center bg-white/5 p-1 rounded-xl">
                <div className="flex gap-2">
                    <button
                        onClick={() => setTab('analytics')}
                        className={`px-6 py-2 rounded-lg text-base font-bold transition-all ${tab === 'analytics' ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/20' : 'text-zinc-500 hover:text-white'}`}
                    >
                        Analytics
                    </button>
                    <button
                        onClick={() => setTab('history')}
                        className={`px-6 py-2 rounded-lg text-base font-bold transition-all ${tab === 'history' ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/20' : 'text-zinc-500 hover:text-white'}`}
                    >
                        Runs History
                    </button>
                </div>

                <button
                    onClick={handleRegenerate}
                    disabled={loadingAction !== null || isActive}
                    className="mr-2 px-4 py-2 rounded-lg bg-white/5 hover:bg-white/10 text-sm font-bold uppercase tracking-widest text-zinc-400 hover:text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                    <span className={loadingAction === 'regenerate' ? 'animate-spin' : ''}>🔄</span>
                    {loadingAction === 'regenerate' ? 'Regenerating...' : 'Regenerate'}
                </button>
            </div>

            {tab === 'analytics' ? (
                <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <div className="glass p-6 rounded-2xl border border-white/5">
                            <p className="text-sm font-bold text-zinc-500 uppercase tracking-widest mb-1">Success Rate</p>
                            <p className="text-3xl font-bold text-emerald-400">{metrics?.successRate || '0%'}</p>
                            <p className="text-xs text-zinc-600 mt-1">{metrics?.successful || 0} / {metrics?.total || 0} Runs Passed</p>
                        </div>
                        <div className="glass p-6 rounded-2xl border border-white/5">
                            <p className="text-sm font-bold text-zinc-500 uppercase tracking-widest mb-1">Avg Duration</p>
                            <p className="text-3xl font-bold">{metrics?.avgDuration || '0s'}</p>
                            <p className="text-xs text-zinc-600 mt-1">Per Repetition</p>
                        </div>
                        <div className="glass p-6 rounded-2xl border border-white/5">
                            <p className="text-sm font-bold text-zinc-500 uppercase tracking-widest mb-1">Lint Issues</p>
                            <p className="text-3xl font-bold text-amber-400">{metrics?.totalLint || 0}</p>
                            <p className="text-xs text-zinc-600 mt-1">Cumulative across all runs</p>
                        </div>
                    </div>

                    {/* Performance Analysis Table */}
                    {runResults.length > 0 && (
                        <section className="glass rounded-2xl overflow-hidden border border-white/5">
                            <div className="p-6 border-b border-white/5 bg-white/[0.02] flex justify-between items-center">
                                <h3 className="text-lg font-bold">Performance Comparison</h3>
                                <div className="flex gap-4 text-xs font-bold uppercase tracking-wider text-zinc-500">
                                    <span>* p&lt;0.05</span>
                                    <span>** p&lt;0.01</span>
                                    <span>*** p&lt;0.001</span>
                                </div>
                            </div>
                            <div className="overflow-x-auto">
                                <table className="w-full text-left">
                                    <thead>
                                        <tr className="text-zinc-500 text-xs uppercase font-bold tracking-widest border-b border-white/5">
                                            <th className="px-6 py-4">Alternative</th>
                                            <th className="px-6 py-4 text-center">Success Rate</th>
                                            <th className="px-6 py-4 text-right">Avg Duration</th>
                                            <th className="px-6 py-4 text-right">Avg Tokens</th>
                                            <th className="px-6 py-4 text-right">Avg Lint</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-white/5">
                                        {altNames.map((alt, idx) => {
                                            const s = stats[alt];

                                            // Calculate significance vs reference (if not reference itself)
                                            let durSig = { level: '', better: false };
                                            let successSig = { level: '', better: false };

                                            if (alt !== referenceAlt && refStats) {
                                                durSig = welchTTest(s.durations, refStats.durations);
                                                // Count successes for Z-test
                                                const sSuccess = s.successCount;
                                                const rSuccess = refStats.successCount;
                                                successSig = zTestProportions(s.successRate / 100, s.count, refStats.successRate / 100, refStats.count);
                                            }

                                            return (
                                                <tr key={idx} className={`group hover:bg-white/[0.02] transition-colors ${alt === referenceAlt ? 'bg-blue-500/5' : ''}`}>
                                                    <td className="px-6 py-4 font-bold text-base text-gray-200">
                                                        {alt}
                                                        {alt === referenceAlt && <span className="ml-2 text-[8px] uppercase bg-blue-500 text-white px-1.5 py-0.5 rounded tracking-widest">Ref</span>}
                                                    </td>
                                                    <td className="px-6 py-4 text-center">
                                                        <div className="flex items-center justify-center gap-1">
                                                            <span className={`font-mono font-bold ${s.successRate >= 90 ? 'text-emerald-400' : s.successRate >= 70 ? 'text-amber-400' : 'text-red-400'
                                                                }`}>
                                                                {s.successRate.toFixed(1)}%
                                                            </span>
                                                            {successSig.level && (
                                                                <span className="text-sm text-blue-400 font-bold" title="Statistically Significant">{successSig.level}</span>
                                                            )}
                                                        </div>
                                                        <div className="text-xs text-zinc-600 mt-0.5">{s.successCount}/{s.count} runs</div>
                                                    </td>
                                                    <td className="px-6 py-4 text-right font-mono text-zinc-300">
                                                        {s.avgDuration.toFixed(2)}s
                                                        {durSig.level && (
                                                            <sup className={`ml-1 font-bold ${durSig.better ? 'text-emerald-400' : 'text-red-400'}`}>{durSig.level}</sup>
                                                        )}
                                                    </td>
                                                    <td className="px-6 py-4 text-right font-mono text-zinc-400">
                                                        {(s.avgTokens / 1000).toFixed(1)}k
                                                    </td>
                                                    <td className="px-6 py-4 text-right">
                                                        <span className={`font-mono ${s.avgLint > 0 ? 'text-amber-400' : 'text-emerald-400'}`}>
                                                            {s.avgLint.toFixed(1)}
                                                        </span>
                                                    </td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            </div>
                        </section>
                    )}

                    {/* Config Summary */}
                    <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <div className="glass p-6 rounded-2xl border border-white/5">
                            <h3 className="text-sm font-bold text-zinc-500 uppercase tracking-widest mb-4">Alternatives</h3>
                            <div className="space-y-3">
                                {config?.alternatives
                                    ?.filter((alt: any) => altNames.includes(alt.name))
                                    .map((alt: any, idx: number) => (
                                        <div key={idx} className="flex justify-between items-center group">
                                            <span className={`text-base font-bold ${alt.name === referenceAlt ? 'text-blue-400' : 'text-gray-300'}`}>
                                                {alt.name} {alt.name === referenceAlt && <span className="text-[10px] bg-blue-500/20 text-blue-300 px-1.5 rounded ml-2">REF</span>}
                                            </span>
                                            <span className="text-sm text-zinc-500 font-mono">{alt.model}</span>
                                        </div>
                                    ))}
                                {config?.alternatives?.filter((alt: any) => altNames.includes(alt.name)).length === 0 && (
                                    <p className="text-zinc-500 text-sm italic">No alternatives run.</p>
                                )}
                            </div>
                        </div>
                        <div className="glass p-6 rounded-2xl border border-white/5">
                            <h3 className="text-sm font-bold text-zinc-500 uppercase tracking-widest mb-4">Scenarios</h3>
                            <div className="flex flex-wrap gap-2">
                                {config?.scenarios
                                    ?.filter((scen: any) => runResults.some(r => r.scenario === scen.name))
                                    .map((scen: any, idx: number) => (
                                        <span key={idx} className="px-2.5 py-1 rounded-md bg-white/5 text-sm font-medium text-gray-300 border border-white/5">
                                            {scen.name}
                                        </span>
                                    ))}
                                {config?.scenarios?.filter((scen: any) => runResults.some(r => r.scenario === scen.name)).length === 0 && (
                                    <p className="text-zinc-500 text-sm italic">No scenarios run.</p>
                                )}
                            </div>
                        </div>
                        <div className="glass p-6 rounded-2xl border border-white/5">
                            <h3 className="text-sm font-bold text-zinc-500 uppercase tracking-widest mb-4">Parameters</h3>
                            <div className="space-y-2 text-base">
                                <div className="flex justify-between">
                                    <span className="text-zinc-500">Repetitions</span>
                                    <span className="font-mono">{config?.repetitions || experiment.reps}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-zinc-500">Concurrent</span>
                                    <span className="font-mono">{config?.max_concurrent || experiment.concurrent}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-zinc-500">Timeout</span>
                                    <span className="font-mono">{config?.timeout || "5m"}</span>
                                </div>
                            </div>
                        </div>
                    </section>

                    {/* Markdown Report Content */}
                    {initialContent && (
                        <section className="glass p-8 rounded-2xl border border-white/5 prose prose-invert prose-emerald max-w-none prose-pre:bg-black/40 prose-pre:border prose-pre:border-white/5">
                            <ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeRaw]}>
                                {initialContent}
                            </ReactMarkdown>
                        </section>
                    )}
                </div>
            ) : (
                <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 animate-in fade-in duration-500">
                    {/* Runs List */}
                    <div className="lg:col-span-4 space-y-4">
                        <h3 className="text-sm font-bold text-zinc-500 uppercase tracking-widest px-2">Repetitions</h3>
                        <div className="glass rounded-2xl border border-white/5 overflow-hidden divide-y divide-white/5">
                            {runResults.map((run, i) => (
                                <button
                                    key={i}
                                    onClick={() => fetchRunDetails(run)}
                                    className={`w-full p-4 text-left transition-all hover:bg-white/5 ${selectedRun?.id === run.id ? 'bg-blue-600/10 border-l-4 border-l-blue-500' : ''}`}
                                >
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <p className="font-bold text-sm text-gray-200">{run.alternative}</p>
                                            <p className="text-xs text-zinc-500 mt-0.5 uppercase tracking-tighter">{run.scenario}</p>
                                        </div>
                                        <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-widest ${run.status === 'SUCCESS' ? 'bg-emerald-500/10 text-emerald-400' : 'bg-red-500/10 text-red-400'}`}>
                                            {run.status}
                                        </span>
                                    </div>
                                    <div className="flex gap-4 mt-3 text-[10px] items-center text-zinc-500 font-mono">
                                        <span title="Tokens">🪙 {(run.agent_metrics?.total_tokens / 1000).toFixed(1)}k</span>
                                        <span title="Tools">🛠️ {run.agent_metrics?.tool_calls_count}</span>
                                        <span title="Duration">⏱️ {(run.duration / 1e9).toFixed(1)}s</span>
                                    </div>
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Run Details */}
                    <div className="lg:col-span-8 space-y-6">
                        {!selectedRun ? (
                            <div className="glass h-[600px] rounded-2xl border border-white/5 flex flex-col items-center justify-center text-zinc-500 space-y-4">
                                <span className="text-4xl">🔬</span>
                                <p className="font-bold uppercase tracking-widest text-sm">Select a run to view details</p>
                            </div>
                        ) : loadingDetails ? (
                            <div className="glass h-[600px] rounded-2xl border border-white/5 flex flex-col items-center justify-center text-zinc-500 animate-pulse">
                                <p className="font-bold uppercase tracking-widest text-sm">Hypothesizing data...</p>
                            </div>
                        ) : (
                            <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-500">
                                {/* Run Header Stats */}
                                <div className="grid grid-cols-4 gap-4">
                                    <div className="glass p-4 rounded-xl border border-white/5">
                                        <p className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest mb-1">Tokens</p>
                                        <p className="text-xl font-bold">{selectedRun.agent_metrics?.total_tokens?.toLocaleString()}</p>
                                        <p className="text-[9px] text-zinc-600 mt-1 uppercase">In_{selectedRun.agent_metrics?.input_tokens} Out_{selectedRun.agent_metrics?.output_tokens}</p>
                                    </div>
                                    <div className="glass p-4 rounded-xl border border-white/5">
                                        <p className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest mb-1">Tools</p>
                                        <p className="text-xl font-bold">{selectedRun.agent_metrics?.tool_calls_count}</p>
                                    </div>
                                    <div className="glass p-4 rounded-xl border border-white/5">
                                        <p className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest mb-1">Tests</p>
                                        <p className="text-xl font-bold text-emerald-400">{selectedRun.tests_passed}</p>
                                        <p className="text-[9px] text-red-500 mt-1 uppercase">{selectedRun.tests_failed} Failures</p>
                                    </div>
                                    <div className="glass p-4 rounded-xl border border-white/5">
                                        <p className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest mb-1">Lint</p>
                                        <p className="text-xl font-bold text-amber-400">{selectedRun.lint_issues}</p>
                                    </div>
                                </div>

                                {/* Thread View */}
                                <div className="glass rounded-2xl border border-white/5 overflow-hidden">
                                    <div className="p-4 border-b border-white/5 bg-white/[0.02]">
                                        <h4 className="text-sm font-bold uppercase tracking-widest">Conversation Thread</h4>
                                    </div>
                                    <div className="p-6 space-y-6 max-h-[500px] overflow-y-auto scrollbar-thin scrollbar-thumb-zinc-800">
                                        {runDetails?.messages.length === 0 && runDetails?.tools.length === 0 && (
                                            <p className="text-zinc-500 italic text-sm text-center py-10">No thread data available for this run.</p>
                                        )}
                                        {runDetails?.messages.map((msg, i) => (
                                            <div key={i} className={`flex flex-col ${msg.role === 'user' ? 'items-end' : 'items-start'}`}>
                                                <span className="text-[10px] font-bold text-zinc-600 uppercase tracking-tighter mb-1 ml-1">{msg.role}</span>
                                                <div className={`p-4 rounded-2xl text-sm max-w-[85%] ${msg.role === 'user' ? 'bg-blue-600/20 text-blue-100' : 'bg-white/5 text-gray-200 shadow-xl'}`}>
                                                    <ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeRaw]}>
                                                        {msg.content}
                                                    </ReactMarkdown>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                {/* Tool Usage */}
                                {runDetails?.tools && runDetails.tools.length > 0 && (
                                    <div className="glass rounded-2xl border border-white/5 overflow-hidden">
                                        <div className="p-4 border-b border-white/5 bg-white/[0.02]">
                                            <h4 className="text-sm font-bold uppercase tracking-widest text-amber-500">Tool Executions</h4>
                                        </div>
                                        <div className="overflow-x-auto">
                                            <table className="w-full text-left text-xs">
                                                <thead className="text-zinc-500 uppercase font-bold tracking-widest bg-white/[0.01]">
                                                    <tr>
                                                        <th className="px-4 py-3">Tool</th>
                                                        <th className="px-4 py-3">Status</th>
                                                        <th className="px-4 py-3">Duration</th>
                                                        <th className="px-4 py-3 text-right">Details</th>
                                                    </tr>
                                                </thead>
                                                <tbody className="divide-y divide-white/5">
                                                    {runDetails.tools.map((tool, i) => (
                                                        <tr key={i} className="hover:bg-white/[0.01]">
                                                            <td className="px-4 py-3 font-mono font-bold text-blue-400">{tool.name}</td>
                                                            <td className="px-4 py-3">
                                                                <span className={tool.status === 'success' ? 'text-emerald-400' : 'text-red-400'}>{tool.status}</span>
                                                            </td>
                                                            <td className="px-4 py-3 text-zinc-500 font-mono">{tool.duration.toFixed(3)}s</td>
                                                            <td className="px-4 py-3 text-right">
                                                                <button onClick={() => setInspectingTool(tool)} className="text-[10px] uppercase text-blue-500 hover:text-blue-400 font-bold transition-colors">Inspect</button>
                                                            </td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>
                                )}
                                {/* Test Results */}
                                {runDetails?.tests && runDetails.tests.length > 0 && (
                                    <div className="glass rounded-2xl border border-white/5 overflow-hidden">
                                        <div className="p-4 border-b border-white/5 bg-white/[0.02] flex justify-between items-center">
                                            <h4 className="text-sm font-bold uppercase tracking-widest text-emerald-500">Test Outcomes</h4>
                                            <div className="flex gap-4">
                                                <span className="text-[10px] text-emerald-400 font-bold uppercase tracking-widest">{runDetails.tests.filter((t: any) => t.status === 'PASS').length} Passed</span>
                                                <span className="text-[10px] text-red-500 font-bold uppercase tracking-widest">{runDetails.tests.filter((t: any) => t.status === 'FAIL').length} Failed</span>
                                            </div>
                                        </div>
                                        <div className="overflow-x-auto max-h-[400px]">
                                            <table className="w-full text-left text-xs">
                                                <thead className="text-zinc-500 uppercase font-bold tracking-widest bg-white/[0.01]">
                                                    <tr>
                                                        <th className="px-4 py-3">Test Name</th>
                                                        <th className="px-4 py-3">Status</th>
                                                        <th className="px-4 py-3 text-right">Duration</th>
                                                    </tr>
                                                </thead>
                                                <tbody className="divide-y divide-white/5">
                                                    {runDetails.tests.map((test: any, i: number) => (
                                                        <tr key={i} className="hover:bg-white/[0.01] group">
                                                            <td className="px-4 py-3 font-mono">
                                                                <p className="font-bold text-gray-300">{test.name}</p>
                                                                {test.status === 'FAIL' && test.output && (
                                                                    <pre className="mt-2 text-[10px] text-red-400 bg-red-400/5 p-2 rounded overflow-x-auto whitespace-pre-wrap font-mono leading-relaxed max-h-[200px]">{test.output}</pre>
                                                                )}
                                                            </td>
                                                            <td className="px-4 py-3">
                                                                <span className={`font-bold ${test.status === 'PASS' ? 'text-emerald-400' : 'text-red-500'}`}>{test.status}</span>
                                                            </td>
                                                            <td className="px-4 py-3 text-right text-zinc-500">{(test.duration_ns / 1e6).toFixed(1)}ms</td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>
                                )}

                                {/* Lint Findings */}
                                {runDetails?.lints && runDetails.lints.length > 0 && (
                                    <div className="glass rounded-2xl border border-white/5 overflow-hidden">
                                        <div className="p-4 border-b border-white/5 bg-white/[0.02]">
                                            <h4 className="text-sm font-bold uppercase tracking-widest text-amber-500">Lint Findings</h4>
                                        </div>
                                        <div className="overflow-x-auto max-h-[400px]">
                                            <table className="w-full text-left text-xs">
                                                <thead className="text-zinc-500 uppercase font-bold tracking-widest bg-white/[0.01]">
                                                    <tr>
                                                        <th className="px-4 py-3">Location</th>
                                                        <th className="px-4 py-3">Message</th>
                                                        <th className="px-4 py-3 text-right">Rule</th>
                                                    </tr>
                                                </thead>
                                                <tbody className="divide-y divide-white/5">
                                                    {runDetails.lints.map((lint: any, i: number) => (
                                                        <tr key={i} className="hover:bg-white/[0.01]">
                                                            <td className="px-4 py-3 font-mono text-zinc-400 min-w-[200px]">
                                                                {lint.file}:{lint.line}:{lint.col}
                                                            </td>
                                                            <td className="px-4 py-3 text-gray-300">{lint.message}</td>
                                                            <td className="px-4 py-3 text-right font-mono text-amber-400/80">{lint.rule_id}</td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>
                                )}

                                {runDetails?.files && runDetails.files.length > 0 && (
                                    <div className="glass rounded-2xl border border-white/5 overflow-hidden">
                                        <div className="p-4 border-b border-white/5 bg-white/[0.02] flex justify-between items-center">
                                            <h4 className="text-sm font-bold uppercase tracking-widest text-emerald-500">Workspace Snapshot</h4>
                                            <span className="text-[10px] text-zinc-500 font-mono uppercase italic">Captured from DB</span>
                                        </div>
                                        <div className="grid grid-cols-12 h-[400px]">
                                            <div className="col-span-4 border-r border-white/5 overflow-y-auto bg-black/20">
                                                {runDetails.files.map((file, i) => (
                                                    <button
                                                        key={i}
                                                        onClick={() => setSelectedFile(file)}
                                                        className={`w-full p-3 text-left text-xs font-mono border-b border-white/5 transition-all hover:bg-white/5 ${selectedFile?.id === file.id ? 'bg-emerald-500/10 text-emerald-400 border-l-2 border-l-emerald-500' : 'text-zinc-400'}`}
                                                    >
                                                        {file.path}
                                                    </button>
                                                ))}
                                            </div>
                                            <div className="col-span-8 overflow-y-auto bg-black/40 p-4">
                                                {selectedFile ? (
                                                    <pre className="text-[11px] font-mono text-zinc-300 whitespace-pre-wrap">{selectedFile.content}</pre>
                                                ) : (
                                                    <p className="text-zinc-600 italic text-center text-xs mt-20 uppercase font-bold tracking-tighter">Select a file to view content</p>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                )}

                                {/* Logs */}
                                {(selectedRun.stdout || selectedRun.stderr) && (
                                    <div className="glass rounded-2xl border border-white/5 overflow-hidden">
                                        <div className="p-4 border-b border-white/5 bg-white/[0.02]">
                                            <h4 className="text-sm font-bold uppercase tracking-widest text-zinc-500">System Logs (Stdout/Stderr)</h4>
                                        </div>
                                        <div className="p-4 bg-black/40 font-mono text-[10px] space-y-4 max-h-[300px] overflow-y-auto">
                                            {selectedRun.stdout && (
                                                <div className="space-y-1">
                                                    <p className="text-blue-500/50 uppercase font-bold">--- stdout ---</p>
                                                    <pre className="text-zinc-400 whitespace-pre-wrap">{selectedRun.stdout}</pre>
                                                </div>
                                            )}
                                            {selectedRun.stderr && (
                                                <div className="space-y-1">
                                                    <p className="text-red-500/50 uppercase font-bold">--- stderr ---</p>
                                                    <pre className="text-red-400/80 whitespace-pre-wrap">{selectedRun.stderr}</pre>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* Tool Inspection Modal */}
            {inspectingTool && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
                    <div className="glass w-full max-w-4xl max-h-[90vh] rounded-3xl border border-white/10 shadow-2xl flex flex-col overflow-hidden text-left">
                        <div className="p-6 border-b border-white/5 flex justify-between items-center bg-white/[0.02]">
                            <div>
                                <h3 className="text-xl font-bold text-blue-400 font-mono flex items-center gap-2">
                                    <span className="text-zinc-500 text-sm">tool:</span> {inspectingTool.name}
                                </h3>
                                <p className="text-[10px] text-zinc-500 mt-1 uppercase tracking-widest font-bold">
                                    Executed in {inspectingTool.duration.toFixed(3)}s • {new Date(inspectingTool.timestamp).toLocaleTimeString()}
                                </p>
                            </div>
                            <button
                                onClick={() => setInspectingTool(null)}
                                className="w-10 h-10 rounded-full bg-white/5 hover:bg-white/10 flex items-center justify-center text-zinc-400 hover:text-white transition-all text-xl"
                            >
                                ✕
                            </button>
                        </div>

                        <div className="flex-1 overflow-y-auto p-8 space-y-8 scrollbar-thin scrollbar-thumb-zinc-800">
                            <section className="space-y-3">
                                <h4 className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest flex items-center gap-2">
                                    <span className="w-1 h-1 rounded-full bg-blue-500"></span> Arguments
                                </h4>
                                <pre className="p-4 rounded-xl bg-black/40 border border-white/5 text-xs text-blue-300 whitespace-pre-wrap font-mono leading-relaxed">
                                    {inspectingTool.args ? JSON.stringify(JSON.parse(inspectingTool.args), null, 2) : "No arguments"}
                                </pre>
                            </section>

                            <section className="space-y-3">
                                <h4 className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest flex items-center gap-2">
                                    <span className={`w-1 h-1 rounded-full ${inspectingTool.status === 'success' ? 'bg-emerald-500' : 'bg-red-500'}`}></span>
                                    {inspectingTool.status === 'success' ? 'Output' : 'Error Response'}
                                </h4>
                                <pre className={`p-4 rounded-xl border text-xs whitespace-pre-wrap font-mono leading-relaxed ${inspectingTool.status === 'success' ? 'bg-black/40 border-white/5 text-zinc-300' : 'bg-red-500/5 border-red-500/20 text-red-400'
                                    }`}>
                                    {inspectingTool.status === 'success' ? inspectingTool.output : inspectingTool.error || inspectingTool.output}
                                </pre>
                            </section>
                        </div>

                        <div className="p-4 border-t border-white/5 bg-white/[0.01] flex justify-end">
                            <button
                                onClick={() => setInspectingTool(null)}
                                className="px-6 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-sm font-bold transition-all text-white"
                            >
                                Close
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
