import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Clock, Plus, AlertTriangle } from "lucide-react";
import axios from "axios";

type TimeEntry = {
    hours: number;
    description?: string;
    user_name: string;
    created_at: string;
};

type TimeData = {
    estimated_hours: number;
    logged_hours: number;
    entries?: TimeEntry[];
    is_overdue: boolean;
    overdue_duration?: string;
};

function TimeTracking({ taskId }: { taskId: string }) {
    const [data, setData] = useState<TimeData | null>(null);
    const [showForm, setShowForm] = useState(false);
    const [hours, setHours] = useState("");
    const [description, setDescription] = useState("");
    const [saving, setSaving] = useState(false);

    const fetchTime = async () => {
        try {
            const res = await axios.get(`/api/tasks/${taskId}/time`);
            setData(res.data.time_tracking);
        } catch (err) {
            console.error("Failed to fetch time tracking:", err);
        }
    };

    useEffect(() => {
        fetchTime();
    }, [taskId]);

    const handleLog = async () => {
        const h = parseFloat(hours);
        if (isNaN(h) || h <= 0) return;
        setSaving(true);
        try {
            const res = await axios.put(`/api/tasks/${taskId}/time`, {
                hours: h,
                description: description || undefined,
            });
            setData(res.data.time_tracking);
            setHours("");
            setDescription("");
            setShowForm(false);
        } catch (err) {
            console.error("Failed to log time:", err);
        } finally {
            setSaving(false);
        }
    };

    const estimated = data?.estimated_hours || 0;
    const logged = data?.logged_hours || 0;
    const pct = estimated > 0 ? Math.min((logged / estimated) * 100, 100) : 0;
    const overBudget = estimated > 0 && logged > estimated;

    return (
        <Card className="mb-6">
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <CardTitle className="text-base flex items-center gap-2">
                        <Clock size={16} />Time Tracking
                    </CardTitle>
                    <Button variant="outline" size="sm" onClick={() => setShowForm(!showForm)}>
                        <Plus size={14} className="mr-1" />Log Time
                    </Button>
                </div>
            </CardHeader>
            <CardContent className="space-y-3">
                {/* progress bar */}
                <div className="space-y-1.5">
                    <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Logged</span>
                        <span className="font-medium">
                            {logged.toFixed(1)}h {estimated > 0 ? `/ ${estimated.toFixed(1)}h estimated` : ""}
                        </span>
                    </div>
                    {estimated > 0 && (
                        <div className="h-2.5 bg-slate-100 rounded-full overflow-hidden">
                            <div
                                className={`h-full rounded-full transition-all duration-500 ${overBudget ? "bg-red-500" : "bg-blue-500"}`}
                                style={{ width: `${pct}%` }}
                            />
                        </div>
                    )}
                    {data?.is_overdue && data.overdue_duration && (
                        <div className="flex items-center gap-1 text-xs text-red-600">
                            <AlertTriangle size={12} />
                            Overdue by {data.overdue_duration}
                        </div>
                    )}
                </div>

                {/* log time form */}
                {showForm && (
                    <div className="border-t pt-3 space-y-2">
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <label className="text-xs font-medium">Hours</label>
                                <Input
                                    type="number"
                                    step="0.25"
                                    min="0.25"
                                    placeholder="e.g. 1.5"
                                    value={hours}
                                    onChange={(e) => setHours(e.target.value)}
                                />
                            </div>
                            <div>
                                <label className="text-xs font-medium">Description (optional)</label>
                                <Input
                                    placeholder="What did you work on?"
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                />
                            </div>
                        </div>
                        <div className="flex gap-2">
                            <Button size="sm" onClick={handleLog} disabled={saving}>
                                {saving ? "Saving..." : "Log Time"}
                            </Button>
                            <Button size="sm" variant="outline" onClick={() => setShowForm(false)}>
                                Cancel
                            </Button>
                        </div>
                    </div>
                )}

                {/* time entries */}
                {data?.entries && data.entries.length > 0 && (
                    <div className="border-t pt-3 space-y-2">
                        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Log History</span>
                        {data.entries.slice().reverse().map((entry, i) => (
                            <div key={i} className="flex items-start gap-3 text-sm py-1.5">
                                <span className="font-semibold text-blue-600 w-12 shrink-0">{entry.hours}h</span>
                                <div className="flex-1 min-w-0">
                                    {entry.description && <p className="text-slate-700 truncate">{entry.description}</p>}
                                    <p className="text-xs text-muted-foreground">
                                        {entry.user_name} &middot; {new Date(entry.created_at).toLocaleDateString("en-GB", { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" })}
                                    </p>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </CardContent>
        </Card>
    );
}

export default TimeTracking;
