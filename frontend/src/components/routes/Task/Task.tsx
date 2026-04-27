import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import Comments from "./Comments";
import TimeTracking from "./TimeTracking";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Logs, PenBox, Trash2, AlertTriangle, Clock } from "lucide-react";
import { Link, useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const statusLabels: Record<string, string> = {
    todo: "To Do",
    in_progress: "In Progress",
    staging_review: "Staging Review",
    done: "Completed",
};

const statusColors: Record<string, string> = {
    todo: "bg-slate-100 text-slate-700",
    in_progress: "bg-blue-100 text-blue-700",
    staging_review: "bg-amber-100 text-amber-700",
    done: "bg-green-100 text-green-700",
};

const priorityColors: Record<string, string> = {
    p1: "bg-red-100 text-red-700",
    p2: "bg-orange-100 text-orange-700",
    p3: "bg-blue-100 text-blue-700",
    p4: "bg-slate-100 text-slate-600",
};

function Task() {
    const { id: taskId } = useParams<{ id: string }>();
    const [task, setTask] = useState<any>(null);
    const [members, setMembers] = useState<any[]>([]);
    const [labels, setLabels] = useState<any[]>([]);
    const navigate = useNavigate();

    const fetchTask = async () => {
        try {
            const res = await axios.get(`/api/tasks/${taskId}`);
            setTask(res.data.task);
        } catch (err) {
            console.error("Failed to fetch task:", err);
        }
    };

    useEffect(() => {
        if (taskId) fetchTask();
    }, [taskId]);

    useEffect(() => {
        if (!task?.project_id) return;
        const fetchMembers = async () => {
            try {
                const res = await axios.get(`/api/projects/${task.project_id}/members`);
                setMembers(res.data.members || []);
            } catch (err) {
                console.error("Failed to fetch members:", err);
            }
        };
        const fetchLabels = async () => {
            try {
                const res = await axios.get(`/api/projects/${task.project_id}/labels`);
                setLabels(res.data.labels || []);
            } catch (err) {
                console.error("Failed to fetch labels:", err);
            }
        };
        fetchMembers();
        fetchLabels();
    }, [task?.project_id]);

    const handleStatusChange = async (status: string) => {
        try {
            await axios.put(`/api/tasks/${taskId}/status`, { status });
            setTask((prev: any) => ({ ...prev, status }));
        } catch (err) {
            console.error("Failed to update status:", err);
        }
    };

    const handleDelete = async () => {
        if (!window.confirm("Delete this task?")) return;
        try {
            await axios.delete(`/api/tasks/${taskId}`);
            navigate(-1);
        } catch (err) {
            console.error("Failed to delete task:", err);
        }
    };

    const handleAssign = async (userId: string) => {
        try {
            const res = await axios.put(`/api/tasks/${taskId}/assign`, {
                assignee_id: userId === "unassigned" ? "" : userId,
            });
            setTask(res.data.task);
        } catch (err) {
            console.error("Failed to assign task:", err);
        }
    };

    if (!task) {
        return (
            <div className="min-h-screen flex flex-col">
                <Header />
                <div className="px-4 pt-20 pb-5 text-muted-foreground flex-1">Loading...</div>
                <Footer />
            </div>
        );
    }

    const dueDate = task.due_date ? new Date(task.due_date) : null;

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="flex-1">
            <div className="max-w-4xl mx-auto px-4 pt-20 pb-8">
                {/* top bar */}
                <div className="flex flex-col sm:flex-row items-start justify-between mb-6 gap-3">
                    <div className="flex-1 min-w-0">
                        <h1 className="text-2xl sm:text-3xl font-bold tracking-tight mb-1">{task.title}</h1>
                        {task.description && (
                            <p className="text-muted-foreground leading-relaxed">{task.description}</p>
                        )}
                    </div>
                    <div className="flex flex-wrap gap-2 shrink-0">
                        <Button variant="outline" size="sm" onClick={() => navigate(`/addEdit?id=${task.id}`)}>
                            <PenBox size={14} className="mr-1" />Edit
                        </Button>
                        <Button variant="outline" size="sm" className="text-red-600 hover:text-red-700" onClick={handleDelete}>
                            <Trash2 size={14} className="mr-1" />Delete
                        </Button>
                        <Link to={`/activity?taskId=${taskId}`}>
                            <Button size="sm">
                                <Logs size={14} className="mr-1" />History
                            </Button>
                        </Link>
                    </div>
                </div>

                {/* metadata grid */}
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 mb-8">
                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Status</span>
                        <Select value={task.status} onValueChange={handleStatusChange}>
                            <SelectTrigger className="h-9">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {Object.entries(statusLabels).map(([key, label]) => (
                                    <SelectItem key={key} value={key}>{label}</SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Priority</span>
                        <div className={`text-sm font-semibold px-3 py-2 rounded-md ${priorityColors[task.priority] || ""}`}>
                            {task.priority?.toUpperCase()}
                        </div>
                    </div>

                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Due Date</span>
                        <div className={`flex items-center gap-1.5 text-sm px-3 py-2 rounded-md ${task.is_past_due ? "bg-red-50 text-red-700" : "bg-slate-50"}`}>
                            {task.is_past_due ? <AlertTriangle size={14} /> : <Clock size={14} />}
                            {dueDate ? dueDate.toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" }) : "Not set"}
                        </div>
                    </div>

                    <div className="space-y-1">
                        <span className="text-xs text-muted-foreground uppercase tracking-wider">Assignee</span>
                        <Select
                            value={task.assignee_id || "unassigned"}
                            onValueChange={handleAssign}
                        >
                            <SelectTrigger className="h-9">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="unassigned">Unassigned</SelectItem>
                                {members.map((m: any) => (
                                    <SelectItem key={m.user_id} value={m.user_id}>
                                        {m.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                </div>

                {/* labels */}
                {task.label_ids?.length > 0 && labels.length > 0 && (
                    <div className="flex flex-wrap gap-2 mb-8">
                        {labels
                            .filter((l: any) => task.label_ids.includes(l.id))
                            .map((l: any) => (
                                <span
                                    key={l.id}
                                    className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                                    style={{ backgroundColor: l.color || "#6b7280" }}
                                >
                                    {l.name}
                                </span>
                            ))}
                    </div>
                )}
            </div>

            <div className="max-w-4xl mx-auto px-4">
                <TimeTracking taskId={taskId!} />
            </div>

            <Comments taskId={taskId!} projectId={task.project_id} />
            </div>
            <Footer />
        </div>
    );
}

export default Task;