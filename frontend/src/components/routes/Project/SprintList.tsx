import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Plus, Trash2, Pencil, ChevronDown, ChevronRight, Calendar, Clock, AlertTriangle, X } from "lucide-react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

type Task = {
    id: string;
    title: string;
    status: string;
    priority: string;
    due_date?: string;
    is_past_due?: boolean;
    assignee?: { name: string };
    sprint_id?: string;
};

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

function SprintList({ projectId }: { projectId: string }) {
    const [sprints, setSprints] = useState<any[]>([]);
    const [tasks, setTasks] = useState<Task[]>([]);
    const [collapsed, setCollapsed] = useState<Record<string, boolean>>({});
    const [open, setOpen] = useState(false);
    const [editingSprint, setEditingSprint] = useState<any>(null);
    const [name, setName] = useState("");
    const [startDate, setStartDate] = useState("");
    const [endDate, setEndDate] = useState("");
    const navigate = useNavigate();

    const fetchSprints = async () => {
        try {
            const res = await axios.get(`/api/projects/${projectId}/sprints`);
            setSprints(res.data.sprints || []);
        } catch (err) {
            console.error("Failed to fetch sprints:", err);
        }
    };

    const fetchTasks = async () => {
        try {
            const res = await axios.get(`/api/projects/${projectId}/tasks`);
            setTasks(res.data.data || []);
        } catch (err) {
            console.error("Failed to fetch tasks:", err);
        }
    };

    useEffect(() => {
        fetchSprints();
        fetchTasks();
    }, [projectId]);

    const resetForm = () => {
        setName("");
        setStartDate("");
        setEndDate("");
        setEditingSprint(null);
    };

    const handleSave = async () => {
        if (!name || !startDate || !endDate) return;
        try {
            const payload = {
                name,
                start_date: new Date(startDate).toISOString(),
                end_date: new Date(endDate).toISOString(),
            };
            if (editingSprint) {
                await axios.put(`/api/sprints/${editingSprint.id}`, payload);
            } else {
                await axios.post(`/api/projects/${projectId}/sprints`, payload);
            }
            resetForm();
            setOpen(false);
            fetchSprints();
        } catch (err) {
            console.error("Failed to save sprint:", err);
        }
    };

    const handleEdit = (sprint: any) => {
        setEditingSprint(sprint);
        setName(sprint.name);
        setStartDate(sprint.start_date?.slice(0, 10) || "");
        setEndDate(sprint.end_date?.slice(0, 10) || "");
        setOpen(true);
    };

    const handleDelete = async (id: string) => {
        if (!window.confirm("Delete this sprint?")) return;
        try {
            await axios.delete(`/api/sprints/${id}`);
            fetchSprints();
            fetchTasks();
        } catch (err) {
            console.error("Failed to delete sprint:", err);
        }
    };

    const handleToggleActive = async (sprint: any) => {
        try {
            await axios.put(`/api/sprints/${sprint.id}`, { is_active: !sprint.is_active });
            fetchSprints();
        } catch (err) {
            console.error("Failed to toggle sprint:", err);
        }
    };

    const handleAddTask = async (sprintId: string, taskId: string) => {
        try {
            await axios.post(`/api/sprints/${sprintId}/tasks`, { task_id: taskId });
            fetchSprints();
            fetchTasks();
        } catch (err) {
            console.error("Failed to add task to sprint:", err);
        }
    };

    const handleRemoveTask = async (sprintId: string, taskId: string) => {
        try {
            await axios.delete(`/api/sprints/${sprintId}/tasks`, { data: { task_id: taskId } });
            fetchSprints();
            fetchTasks();
        } catch (err) {
            console.error("Failed to remove task:", err);
        }
    };

    const toggle = (id: string) => setCollapsed((prev) => ({ ...prev, [id]: !prev[id] }));

    const formatDate = (d: string) =>
        new Date(d).toLocaleDateString("en-GB", { day: "2-digit", month: "short" });

    const daysLeft = (end: string) => {
        const diff = Math.ceil((new Date(end).getTime() - Date.now()) / (1000 * 60 * 60 * 24));
        if (diff < 0) return `${Math.abs(diff)}d overdue`;
        if (diff === 0) return "Ends today";
        return `${diff}d left`;
    };

    const getSprintTasks = (sprintId: string) => tasks.filter((t) => t.sprint_id === sprintId);
    const backlogTasks = tasks.filter((t) => !t.sprint_id);

    const statusCounts = (items: Task[]) => {
        const c: Record<string, number> = { todo: 0, in_progress: 0, staging_review: 0, done: 0 };
        items.forEach((t) => { c[t.status] = (c[t.status] || 0) + 1; });
        return c;
    };

    const TaskRow = ({ task, sprintId }: { task: Task; sprintId?: string }) => (
        <div
            className="flex items-center gap-3 px-4 py-2.5 border-b last:border-b-0 hover:bg-slate-50 cursor-pointer group"
            onClick={() => navigate(`/task/${task.id}`)}
        >
            <div className={`w-2 h-2 rounded-full shrink-0 ${priorityDot[task.priority] || "bg-slate-300"}`} />
            <span className="text-sm flex-1 truncate">{task.title}</span>
            <span className={`text-[10px] font-semibold px-2 py-0.5 rounded ${statusBadge[task.status] || ""}`}>
                {statusLabel[task.status] || task.status}
            </span>
            <span className="text-xs text-muted-foreground w-24 text-right truncate">
                {task.assignee?.name || "Unassigned"}
            </span>
            {task.due_date && (
                <span className={`text-xs w-16 text-right ${task.is_past_due ? "text-red-600 font-medium" : "text-muted-foreground"}`}>
                    {formatDate(task.due_date)}
                </span>
            )}
            {sprintId && (
                <button
                    className="opacity-0 group-hover:opacity-100 text-muted-foreground hover:text-red-500 transition-opacity"
                    onClick={(e) => { e.stopPropagation(); handleRemoveTask(sprintId, task.id); }}
                    title="Remove from sprint"
                >
                    <X size={14} />
                </button>
            )}
        </div>
    );

    return (
        <div className="space-y-3">
            {/* top bar */}
            <div className="flex justify-between items-center">
                <h2 className="text-lg font-semibold">Sprints</h2>
                <Dialog open={open} onOpenChange={(val) => { setOpen(val); if (!val) resetForm(); }}>
                    <DialogTrigger asChild>
                        <Button variant="outline" size="sm">
                            <Plus size={16} className="mr-1" />Create Sprint
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-md" aria-describedby={undefined}>
                        <DialogHeader>
                            <DialogTitle>{editingSprint ? "Edit Sprint" : "Create Sprint"}</DialogTitle>
                        </DialogHeader>
                        <div className="space-y-3">
                            <div>
                                <label className="text-sm font-medium">Name</label>
                                <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Sprint 1" />
                            </div>
                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="text-sm font-medium">Start</label>
                                    <Input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
                                </div>
                                <div>
                                    <label className="text-sm font-medium">End</label>
                                    <Input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
                                </div>
                            </div>
                        </div>
                        <DialogFooter>
                            <DialogClose asChild>
                                <Button variant="outline">Cancel</Button>
                            </DialogClose>
                            <Button onClick={handleSave}>{editingSprint ? "Update" : "Create"}</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            {/* sprint panels */}
            {sprints.map((sprint) => {
                const sprintTasks = getSprintTasks(sprint.id);
                const counts = statusCounts(sprintTasks);
                const isCollapsed = collapsed[sprint.id];
                const overdue = new Date(sprint.end_date) < new Date() && sprint.is_active;

                return (
                    <div key={sprint.id} className="border rounded-lg overflow-hidden">
                        {/* sprint header */}
                        <div className={`flex flex-wrap items-center gap-3 px-4 py-3 ${sprint.is_active ? "bg-blue-50" : "bg-slate-50"}`}>
                            <button onClick={() => toggle(sprint.id)} className="shrink-0">
                                {isCollapsed ? <ChevronRight size={18} /> : <ChevronDown size={18} />}
                            </button>
                            <div className="flex-1 min-w-0">
                                <div className="flex items-center gap-2 flex-wrap">
                                    <span
                                        className="font-semibold text-sm hover:text-blue-600 hover:underline cursor-pointer"
                                        onClick={(e) => { e.stopPropagation(); navigate(`/sprint/${sprint.id}`); }}
                                    >{sprint.name}</span>
                                    {sprint.is_active && (
                                        <span className="text-[10px] font-bold bg-blue-600 text-white px-1.5 py-0.5 rounded">ACTIVE</span>
                                    )}
                                    {overdue && (
                                        <span className="text-[10px] font-bold bg-red-100 text-red-700 px-1.5 py-0.5 rounded flex items-center gap-0.5">
                                            <AlertTriangle size={10} />OVERDUE
                                        </span>
                                    )}
                                </div>
                                <div className="flex items-center gap-4 text-xs text-muted-foreground mt-0.5">
                                    <span className="flex items-center gap-1"><Calendar size={12} />{formatDate(sprint.start_date)} – {formatDate(sprint.end_date)}</span>
                                    <span className="flex items-center gap-1"><Clock size={12} />{daysLeft(sprint.end_date)}</span>
                                    <span>{sprintTasks.length} issues</span>
                                </div>
                            </div>

                            {/* progress bar */}
                            {sprintTasks.length > 0 && (
                                <div className="flex h-2 w-20 sm:w-32 rounded-full overflow-hidden bg-slate-200 shrink-0">
                                    {counts.done > 0 && <div className="bg-green-500" style={{ width: `${(counts.done / sprintTasks.length) * 100}%` }} />}
                                    {counts.in_progress > 0 && <div className="bg-blue-500" style={{ width: `${(counts.in_progress / sprintTasks.length) * 100}%` }} />}
                                    {counts.staging_review > 0 && <div className="bg-amber-400" style={{ width: `${(counts.staging_review / sprintTasks.length) * 100}%` }} />}
                                </div>
                            )}

                            <div className="flex gap-1 shrink-0">
                                <Button
                                    variant={sprint.is_active ? "default" : "outline"}
                                    size="sm"
                                    className="text-xs h-7"
                                    onClick={() => handleToggleActive(sprint)}
                                >
                                    {sprint.is_active ? "Complete Sprint" : "Start Sprint"}
                                </Button>
                                <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => handleEdit(sprint)}>
                                    <Pencil size={13} />
                                </Button>
                                <Button variant="ghost" size="icon" className="h-7 w-7 text-red-500 hover:text-red-700" onClick={() => handleDelete(sprint.id)}>
                                    <Trash2 size={13} />
                                </Button>
                            </div>
                        </div>

                        {/* tasks inside sprint */}
                        {!isCollapsed && (
                            <div>
                                {sprintTasks.length === 0 ? (
                                    <div className="px-4 py-6 text-center text-sm text-muted-foreground border-t">
                                        No issues in this sprint yet.
                                    </div>
                                ) : (
                                    sprintTasks.map((task) => <TaskRow key={task.id} task={task} sprintId={sprint.id} />)
                                )}
                            </div>
                        )}
                    </div>
                );
            })}

            {/* backlog panel */}
            <div className="border rounded-lg overflow-hidden">
                <div className="flex items-center gap-3 px-4 py-3 bg-slate-50">
                    <button onClick={() => toggle("backlog")} className="shrink-0">
                        {collapsed["backlog"] ? <ChevronRight size={18} /> : <ChevronDown size={18} />}
                    </button>
                    <span className="font-semibold text-sm">Backlog</span>
                    <span className="text-xs text-muted-foreground">({backlogTasks.length} issues)</span>
                </div>
                {!collapsed["backlog"] && (
                    <div>
                        {backlogTasks.length === 0 ? (
                            <div className="px-4 py-6 text-center text-sm text-muted-foreground border-t">Backlog is empty.</div>
                        ) : (
                            backlogTasks.map((task) => (
                                <div key={task.id} className="flex items-center border-b last:border-b-0">
                                    <div className="flex-1"><TaskRow task={task} /></div>
                                    {sprints.length > 0 && (
                                        <div className="pr-3 shrink-0">
                                            <select
                                                className="text-xs border rounded px-2 py-1 bg-white cursor-pointer"
                                                defaultValue=""
                                                onChange={(e) => {
                                                    if (e.target.value) handleAddTask(e.target.value, task.id);
                                                    e.target.value = "";
                                                }}
                                            >
                                                <option value="" disabled>+ Sprint</option>
                                                {sprints.map((s) => (
                                                    <option key={s.id} value={s.id}>{s.name}</option>
                                                ))}
                                            </select>
                                        </div>
                                    )}
                                </div>
                            ))
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}

export default SprintList;
