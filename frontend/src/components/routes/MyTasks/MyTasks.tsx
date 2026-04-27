import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useNavigate } from "react-router-dom";
import { AlertTriangle, Clock, ChevronLeft, ChevronRight } from "lucide-react";
import axios from "axios";

const statusLabels: Record<string, string> = {
    todo: "To Do",
    in_progress: "In Progress",
    staging_review: "Review",
    done: "Done",
};

const statusColors: Record<string, string> = {
    todo: "bg-slate-200 text-slate-700",
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

function MyTasks() {
    const [tasks, setTasks] = useState<any[]>([]);
    const [pagination, setPagination] = useState<any>(null);
    const [statusFilter, setStatusFilter] = useState("all");
    const [priorityFilter, setPriorityFilter] = useState("all");
    const [page, setPage] = useState(1);
    const navigate = useNavigate();

    const fetchTasks = async () => {
        try {
            const params: Record<string, string> = { page: String(page), page_size: "20" };
            if (statusFilter !== "all") params.status = statusFilter;
            if (priorityFilter !== "all") params.priority = priorityFilter;
            const res = await axios.get("/api/tasks/my", { params });
            setTasks(res.data.data || []);
            setPagination(res.data.pagination || null);
        } catch (err) {
            console.error("Failed to fetch my tasks:", err);
        }
    };

    useEffect(() => {
        fetchTasks();
    }, [page, statusFilter, priorityFilter]);

    const resetFilters = () => {
        setStatusFilter("all");
        setPriorityFilter("all");
        setPage(1);
    };

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="px-4 pt-20 pb-8 flex-1">
                <h1 className="text-3xl font-bold mb-1">My Tasks</h1>
                <p className="text-muted-foreground text-sm mb-5">Tasks assigned to you across all projects.</p>

                {/* filters */}
                <div className="flex items-center gap-3 mb-4 flex-wrap">
                    <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1); }}>
                        <SelectTrigger className="w-40">
                            <SelectValue placeholder="Status" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">All Statuses</SelectItem>
                            {Object.entries(statusLabels).map(([k, v]) => (
                                <SelectItem key={k} value={k}>{v}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>

                    <Select value={priorityFilter} onValueChange={(v) => { setPriorityFilter(v); setPage(1); }}>
                        <SelectTrigger className="w-40">
                            <SelectValue placeholder="Priority" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">All Priorities</SelectItem>
                            <SelectItem value="p1">P1 Critical</SelectItem>
                            <SelectItem value="p2">P2 High</SelectItem>
                            <SelectItem value="p3">P3 Medium</SelectItem>
                            <SelectItem value="p4">P4 Low</SelectItem>
                        </SelectContent>
                    </Select>

                    {(statusFilter !== "all" || priorityFilter !== "all") && (
                        <Button variant="ghost" size="sm" onClick={resetFilters}>Clear filters</Button>
                    )}
                </div>

                {/* task list */}
                {tasks.length === 0 ? (
                    <div className="text-center py-16">
                        <p className="text-muted-foreground">No tasks found.</p>
                    </div>
                ) : (
                    <div className="space-y-2">
                        {tasks.map((task) => {
                            const dueDate = task.due_date ? new Date(task.due_date) : null;
                            return (
                                <Card
                                    key={task.id}
                                    className="hover:shadow-sm transition-shadow cursor-pointer"
                                    onClick={() => navigate(`/task/${task.id}`)}
                                >
                                    <CardContent className="py-3 px-4 flex flex-wrap items-center gap-3 sm:gap-4">
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center gap-2 mb-0.5">
                                                <span className="font-medium text-sm truncate">{task.title}</span>
                                            </div>
                                            {task.project_id && (
                                                <span className="text-xs text-muted-foreground">
                                                    {task.project_name || task.project_id}
                                                </span>
                                            )}
                                        </div>
                                        <span className={`text-[10px] font-semibold px-2 py-0.5 rounded shrink-0 ${statusColors[task.status] || ""}`}>
                                            {statusLabels[task.status] || task.status}
                                        </span>
                                        <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded shrink-0 ${priorityColors[task.priority] || ""}`}>
                                            {task.priority?.toUpperCase()}
                                        </span>
                                        {dueDate && (
                                            <span className={`flex items-center gap-1 text-xs shrink-0 ${task.is_past_due ? "text-red-600 font-medium" : "text-muted-foreground"}`}>
                                                {task.is_past_due ? <AlertTriangle size={12} /> : <Clock size={12} />}
                                                {dueDate.toLocaleDateString("en-GB", { day: "2-digit", month: "short" })}
                                            </span>
                                        )}
                                    </CardContent>
                                </Card>
                            );
                        })}
                    </div>
                )}

                {/* pagination */}
                {pagination && pagination.total_pages > 1 && (
                    <div className="flex items-center justify-center gap-3 mt-6">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={page <= 1}
                            onClick={() => setPage((p) => p - 1)}
                        >
                            <ChevronLeft size={14} />Prev
                        </Button>
                        <span className="text-sm text-muted-foreground">
                            Page {pagination.page} of {pagination.total_pages}
                        </span>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={page >= pagination.total_pages}
                            onClick={() => setPage((p) => p + 1)}
                        >
                            Next<ChevronRight size={14} />
                        </Button>
                    </div>
                )}
            </div>
            <Footer />
        </div>
    );
}

export default MyTasks;
