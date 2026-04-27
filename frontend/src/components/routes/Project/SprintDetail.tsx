import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useNavigate, useParams } from "react-router-dom";
import { Calendar, Clock, AlertTriangle, ArrowLeft, CheckCircle2 } from "lucide-react";
import axios from "axios";

const statusBadge: Record<string, string> = {
    todo: "bg-slate-200 text-slate-700",
    in_progress: "bg-blue-100 text-blue-700",
    staging_review: "bg-amber-100 text-amber-700",
    done: "bg-green-100 text-green-700",
};

const statusLabel: Record<string, string> = {
    todo: "TO DO",
    in_progress: "IN PROGRESS",
    staging_review: "REVIEW",
    done: "DONE",
};

const priorityDot: Record<string, string> = {
    p1: "bg-red-500",
    p2: "bg-orange-400",
    p3: "bg-blue-400",
    p4: "bg-slate-300",
};

function SprintDetail() {
    const { id: sprintId } = useParams<{ id: string }>();
    const [sprint, setSprint] = useState<any>(null);
    const [tasks, setTasks] = useState<any[]>([]);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchSprint = async () => {
            try {
                const res = await axios.get(`/api/sprints/${sprintId}`);
                setSprint(res.data.sprint);
            } catch (err) {
                console.error("Failed to fetch sprint:", err);
            }
        };
        const fetchTasks = async () => {
            try {
                // get sprint to find project_id, then filter tasks
                const sRes = await axios.get(`/api/sprints/${sprintId}`);
                const projectId = sRes.data.sprint?.project_id;
                if (projectId) {
                    const tRes = await axios.get(`/api/projects/${projectId}/tasks`);
                    const all = tRes.data.data || [];
                    setTasks(all.filter((t: any) => t.sprint_id === sprintId));
                }
            } catch (err) {
                console.error("Failed to fetch sprint tasks:", err);
            }
        };
        if (sprintId) {
            fetchSprint();
            fetchTasks();
        }
    }, [sprintId]);

    if (!sprint) {
        return (
            <div className="min-h-screen flex flex-col">
                <Header />
                <div className="px-4 pt-20 pb-5 text-muted-foreground flex-1">Loading...</div>
                <Footer />
            </div>
        );
    }

    const formatDate = (d: string) =>
        new Date(d).toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });

    const daysLeft = () => {
        const diff = Math.ceil((new Date(sprint.end_date).getTime() - Date.now()) / (1000 * 60 * 60 * 24));
        if (diff < 0) return `${Math.abs(diff)} days overdue`;
        if (diff === 0) return "Ends today";
        return `${diff} days remaining`;
    };

    const overdue = new Date(sprint.end_date) < new Date() && sprint.is_active;

    const counts: Record<string, number> = { todo: 0, in_progress: 0, staging_review: 0, done: 0 };
    tasks.forEach((t) => { counts[t.status] = (counts[t.status] || 0) + 1; });
    const total = tasks.length;

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="max-w-4xl mx-auto px-4 pt-20 pb-8 flex-1">
                <Button variant="ghost" size="sm" onClick={() => navigate(-1)} className="mb-4">
                    <ArrowLeft size={14} className="mr-1" />Back
                </Button>

                <div className="flex items-center gap-3 mb-2">
                    <h1 className="text-3xl font-bold">{sprint.name}</h1>
                    {sprint.is_active && (
                        <span className="text-xs font-bold bg-blue-600 text-white px-2 py-0.5 rounded">ACTIVE</span>
                    )}
                    {overdue && (
                        <span className="text-xs font-bold bg-red-100 text-red-700 px-2 py-0.5 rounded flex items-center gap-1">
                            <AlertTriangle size={12} />OVERDUE
                        </span>
                    )}
                </div>

                {/* sprint metadata */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Start Date</span>
                        <div className="flex items-center gap-1.5 text-sm bg-slate-50 px-3 py-2 rounded-md">
                            <Calendar size={14} />{formatDate(sprint.start_date)}
                        </div>
                    </div>
                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">End Date</span>
                        <div className="flex items-center gap-1.5 text-sm bg-slate-50 px-3 py-2 rounded-md">
                            <Calendar size={14} />{formatDate(sprint.end_date)}
                        </div>
                    </div>
                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Time</span>
                        <div className={`flex items-center gap-1.5 text-sm px-3 py-2 rounded-md ${overdue ? "bg-red-50 text-red-700" : "bg-slate-50"}`}>
                            <Clock size={14} />{daysLeft()}
                        </div>
                    </div>
                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Issues</span>
                        <div className="flex items-center gap-1.5 text-sm bg-slate-50 px-3 py-2 rounded-md">
                            <CheckCircle2 size={14} />{counts.done}/{total} done
                        </div>
                    </div>
                </div>

                {/* progress bar */}
                {total > 0 && (
                    <div className="mb-6">
                        <div className="flex h-3 rounded-full overflow-hidden bg-slate-100">
                            {counts.done > 0 && <div className="bg-green-500" style={{ width: `${(counts.done / total) * 100}%` }} />}
                            {counts.in_progress > 0 && <div className="bg-blue-500" style={{ width: `${(counts.in_progress / total) * 100}%` }} />}
                            {counts.staging_review > 0 && <div className="bg-amber-400" style={{ width: `${(counts.staging_review / total) * 100}%` }} />}
                            {counts.todo > 0 && <div className="bg-slate-300" style={{ width: `${(counts.todo / total) * 100}%` }} />}
                        </div>
                        <div className="flex gap-4 mt-2 text-xs text-muted-foreground">
                            <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-green-500" />Done {counts.done}</span>
                            <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-blue-500" />In Progress {counts.in_progress}</span>
                            <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-amber-400" />Review {counts.staging_review}</span>
                            <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-slate-300" />To Do {counts.todo}</span>
                        </div>
                    </div>
                )}

                {/* task list */}
                <h2 className="text-lg font-semibold mb-3">Issues ({total})</h2>
                {total === 0 ? (
                    <p className="text-sm text-muted-foreground text-center py-8">No issues in this sprint.</p>
                ) : (
                    <div className="border rounded-lg overflow-hidden">
                        {tasks.map((task) => {
                            const dueDate = task.due_date ? new Date(task.due_date) : null;
                            return (
                                <div
                                    key={task.id}
                                    className="flex items-center gap-3 px-4 py-3 border-b last:border-b-0 hover:bg-slate-50 cursor-pointer"
                                    onClick={() => navigate(`/task/${task.id}`)}
                                >
                                    <div className={`w-2 h-2 rounded-full shrink-0 ${priorityDot[task.priority] || "bg-slate-300"}`} />
                                    <span className="text-sm flex-1 truncate font-medium">{task.title}</span>
                                    <span className={`text-[10px] font-semibold px-2 py-0.5 rounded ${statusBadge[task.status] || ""}`}>
                                        {statusLabel[task.status] || task.status}
                                    </span>
                                    <span className="text-xs text-muted-foreground w-24 text-right truncate">
                                        {task.assignee?.name || "Unassigned"}
                                    </span>
                                    {dueDate && (
                                        <span className={`text-xs w-16 text-right ${task.is_past_due ? "text-red-600 font-medium" : "text-muted-foreground"}`}>
                                            {dueDate.toLocaleDateString("en-GB", { day: "2-digit", month: "short" })}
                                        </span>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
            <Footer />
        </div>
    );
}

export default SprintDetail;
