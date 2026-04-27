import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import axios from "axios";

const statusConfig: Record<string, { label: string; color: string }> = {
    todo: { label: "To Do", color: "#94a3b8" },
    in_progress: { label: "In Progress", color: "#3b82f6" },
    staging_review: { label: "Review", color: "#f59e0b" },
    done: { label: "Done", color: "#22c55e" },
};

const priorityConfig: Record<string, { label: string; color: string }> = {
    p1: { label: "P1 Critical", color: "#ef4444" },
    p2: { label: "P2 High", color: "#f97316" },
    p3: { label: "P3 Medium", color: "#3b82f6" },
    p4: { label: "P4 Low", color: "#94a3b8" },
};

type ChartEntry = { status?: string; priority?: string; count: number };

function Bar({ label, count, total, color }: { label: string; count: number; total: number; color: string }) {
    const pct = total > 0 ? (count / total) * 100 : 0;
    return (
        <div className="flex items-center gap-3">
            <span className="text-sm w-28 text-right shrink-0">{label}</span>
            <div className="flex-1 bg-slate-100 rounded-full h-7 overflow-hidden relative">
                <div
                    className="h-full rounded-full transition-all duration-500"
                    style={{ width: `${pct}%`, backgroundColor: color }}
                />
                {count > 0 && (
                    <span className="absolute inset-0 flex items-center px-3 text-xs font-semibold text-slate-700">
                        {count}
                    </span>
                )}
            </div>
            <span className="text-xs text-muted-foreground w-10">{pct.toFixed(0)}%</span>
        </div>
    );
}

function ProjectCharts({ projectId }: { projectId: string }) {
    const [statusData, setStatusData] = useState<ChartEntry[]>([]);
    const [priorityData, setPriorityData] = useState<ChartEntry[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetch = async () => {
            try {
                const [sRes, pRes] = await Promise.all([
                    axios.get(`/api/projects/${projectId}/charts/status`),
                    axios.get(`/api/projects/${projectId}/charts/priority`),
                ]);
                setStatusData(sRes.data.data || []);
                setPriorityData(pRes.data.data || []);
            } catch (err) {
                console.error("Failed to fetch charts:", err);
            } finally {
                setLoading(false);
            }
        };
        fetch();
    }, [projectId]);

    if (loading) return <p className="text-sm text-muted-foreground py-8 text-center">Loading charts...</p>;

    const statusTotal = statusData.reduce((s, d) => s + d.count, 0);
    const priorityTotal = priorityData.reduce((s, d) => s + d.count, 0);

    return (
        <div className="grid md:grid-cols-2 gap-6">
            <Card>
                <CardHeader className="pb-3">
                    <CardTitle className="text-base">Tasks by Status</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    {statusTotal === 0 ? (
                        <p className="text-sm text-muted-foreground text-center py-4">No tasks yet.</p>
                    ) : (
                        Object.entries(statusConfig).map(([key, cfg]) => {
                            const entry = statusData.find((d) => d.status === key);
                            return (
                                <Bar key={key} label={cfg.label} count={entry?.count || 0} total={statusTotal} color={cfg.color} />
                            );
                        })
                    )}
                    {statusTotal > 0 && (
                        <p className="text-xs text-muted-foreground text-right pt-1">{statusTotal} total tasks</p>
                    )}
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="pb-3">
                    <CardTitle className="text-base">Tasks by Priority</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    {priorityTotal === 0 ? (
                        <p className="text-sm text-muted-foreground text-center py-4">No tasks yet.</p>
                    ) : (
                        Object.entries(priorityConfig).map(([key, cfg]) => {
                            const entry = priorityData.find((d) => d.priority === key);
                            return (
                                <Bar key={key} label={cfg.label} count={entry?.count || 0} total={priorityTotal} color={cfg.color} />
                            );
                        })
                    )}
                    {priorityTotal > 0 && (
                        <p className="text-xs text-muted-foreground text-right pt-1">{priorityTotal} total tasks</p>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}

export default ProjectCharts;
