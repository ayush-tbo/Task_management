import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { format } from 'date-fns';
import { CalendarIcon, Loader2 } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { taskSchema } from "@/lib/schema";
import axios from "axios";
import { useSearchParams } from "react-router-dom";
import { Textarea } from "@/components/ui/textarea";
import { useNavigate } from "react-router-dom";

function TaskForm(){
    
    const [loading, setLoading] = useState(false);
    const [editMode, setEditMode] = useState(false);
    const [members, setMembers] = useState<any[]>([]);
    const [labels, setLabels] = useState<any[]>([]);
    const [sprints, setSprints] = useState<any[]>([]);
    const [selectedLabels, setSelectedLabels] = useState<string[]>([]);
    const [selectedSprint, setSelectedSprint] = useState("");

    const [searchParams] = useSearchParams();
    const taskId = searchParams.get("id");
    const projectId = searchParams.get("projectId");

    // track the resolved project id (from URL or from fetched task)
    const [resolvedProjectId, setResolvedProjectId] = useState(projectId || "");

    const navigate = useNavigate();

    const { register, handleSubmit, formState:{errors}, setValue, watch, reset } = useForm({
        resolver:zodResolver(taskSchema),
        defaultValues:{
            title: "",
            due_date: new Date(),
            description: "",
            priority: "p3" as "p1" | "p2" | "p3" | "p4",
            status: "todo" as "todo" | "in_progress" | "staging_review" | "done",
            assignee_id: "",
        },
    });

    useEffect(() => {
        if(taskId){
            setEditMode(true);
            getTask();
        }
    }, [taskId]);

    useEffect(() => {
        if (projectId) setResolvedProjectId(projectId);
    }, [projectId]);

    useEffect(() => {
        if (!resolvedProjectId) return;
        const fetchMembers = async () => {
            try {
                const res = await axios.get(`/api/projects/${resolvedProjectId}/members`);
                setMembers(res.data.members || []);
            } catch (err) {
                console.error("Failed to fetch members:", err);
            }
        };
        const fetchLabels = async () => {
            try {
                const res = await axios.get(`/api/projects/${resolvedProjectId}/labels`);
                setLabels(res.data.labels || []);
            } catch (err) {
                console.error("Failed to fetch labels:", err);
            }
        };
        const fetchSprints = async () => {
            try {
                const res = await axios.get(`/api/projects/${resolvedProjectId}/sprints`);
                setSprints(res.data.sprints || []);
            } catch (err) {
                console.error("Failed to fetch sprints:", err);
            }
        };
        fetchMembers();
        fetchLabels();
        fetchSprints();
    }, [resolvedProjectId]);

    const dueDate = watch("due_date");
    const priority = watch("priority");

    const onSubmit = async (data: any) => {
        try{
            setLoading(true);
            const payload = {
                ...data,
                due_date: data.due_date.toISOString(),
                assignee_id: data.assignee_id || undefined,
                label_ids: selectedLabels.length > 0 ? selectedLabels : undefined,
                sprint_id: selectedSprint || undefined,
            };
            if (editMode && taskId) {
                await axios.put(`/api/tasks/${taskId}`, payload);
            } else if (projectId) {
                await axios.post(`/api/projects/${projectId}/tasks`, payload);
            }
            navigate(-1);
        }
        catch(err){
            console.error("Failed to save task:", err);
        } finally {
            setLoading(false);
        }
    };

    const getTask = async () => {
        try{
            const res = await axios.get(`/api/tasks/${taskId}`);
            const task = res.data.task;
            if (task.project_id) setResolvedProjectId(task.project_id);
            setSelectedLabels(task.label_ids || []);
            setSelectedSprint(task.sprint_id || "");
            reset({
                title: task.title,
                due_date: task.due_date ? new Date(task.due_date) : new Date(),
                description: task.description || "",
                priority: task.priority || "p3",
                status: task.status || "todo",
                assignee_id: task.assignee_id || "",
            });
        }
        catch(err){
            console.error("Failed to get task:", err);
        }
    };

    return (
        <div className="max-w-3xl mx-auto px-5 pt-20">
            <h1 className="text-3xl sm:text-5xl gradient-title">{editMode ? "Edit" : "Add"} task</h1>
            <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
                <div className="space-y-2">
                    <label htmlFor="title" className="text-sm font-medium">Title</label>
                    <Input id="title" placeholder="e.g., Initialize monorepo repository" {...register("title")} />
                    {errors.title && (
                        <p className="text-sm text-red-500">{errors.title.message}</p>
                    )}
                </div>

                <div className="grid gap-6 md:grid-cols-2">
                    <div className="space-y-2">
                        <label htmlFor="status" className="text-sm font-medium">Task Status</label>
                        <Select onValueChange={(value) => setValue("status", value as "todo" | "in_progress" | "staging_review" | "done")} defaultValue={watch("status")}>
                            <SelectTrigger className="w-full" id="status">
                                <SelectValue placeholder="Select Status" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="todo">To Do</SelectItem>
                                <SelectItem value="in_progress">In Progress</SelectItem>
                                <SelectItem value="staging_review">Staging Review</SelectItem>
                                <SelectItem value="done">Completed</SelectItem>
                            </SelectContent>
                        </Select>
                        {errors.status && (
                            <p className="text-sm text-red-500">{errors.status.message}</p>
                        )}
                    </div>
                    <div className="space-y-2">
                        <label htmlFor="assignee_id" className="text-sm font-medium">Assignee</label>
                        <Select
                            value={watch("assignee_id") || "unassigned"}
                            onValueChange={(value) => setValue("assignee_id", value === "unassigned" ? "" : value)}
                        >
                            <SelectTrigger className="w-full" id="assignee_id">
                                <SelectValue placeholder="Select Assignee" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="unassigned">Unassigned</SelectItem>
                                {members.map((m) => (
                                    <SelectItem key={m.user_id} value={m.user_id}>
                                        {m.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        {errors.assignee_id && (
                            <p className="text-sm text-red-500">{errors.assignee_id.message}</p>
                        )}
                    </div>
                </div>

                <div className="space-y-2">
                    <label htmlFor="description" className="text-sm font-medium">Description</label>
                    <Textarea id="description" placeholder="e.g., Set up the base Go backend, React TypeScript frontend, and Docker compose configurations." {...register("description")} />
                    {errors.description && (
                        <p className="text-sm text-red-500">{errors.description.message}</p>
                    )}
                </div>

                <div className="space-y-2">
                    <label className="text-sm font-medium">Labels</label>
                    {labels.length > 0 ? (
                        <div className="flex flex-wrap gap-2">
                            {labels.map((l: any) => {
                                const isSelected = selectedLabels.includes(l.id);
                                return (
                                    <button
                                        key={l.id}
                                        type="button"
                                        onClick={() => setSelectedLabels((prev) =>
                                            isSelected ? prev.filter((id) => id !== l.id) : [...prev, l.id]
                                        )}
                                        className={`text-xs px-2.5 py-1 rounded-full border transition-colors ${
                                            isSelected
                                                ? "text-white border-transparent"
                                                : "bg-white text-slate-700 border-slate-300 hover:border-slate-400"
                                        }`}
                                        style={isSelected ? { backgroundColor: l.color || "#6b7280" } : {}}
                                    >
                                        {l.name}
                                    </button>
                                );
                            })}
                        </div>
                    ) : (
                        <p className="text-xs text-muted-foreground">No labels in this project.</p>
                    )}
                </div>

                <div className="space-y-2">
                    <label htmlFor="sprint_id" className="text-sm font-medium">Sprint</label>
                    <Select
                        value={selectedSprint || "none"}
                        onValueChange={(value) => setSelectedSprint(value === "none" ? "" : value)}
                    >
                        <SelectTrigger className="w-full" id="sprint_id">
                            <SelectValue placeholder="Select Sprint" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="none">No Sprint</SelectItem>
                            {sprints.map((s: any) => (
                                <SelectItem key={s.id} value={s.id}>{s.name}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>

                <div className="space-y-2">
                    <label className="text-sm font-medium">Due Date</label>
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button variant="outline" className="w-full pl-3 text-left font-normal">
                                {dueDate ? format(dueDate, "PPP") : <span>Pick a Date</span>}
                                <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                            <Calendar
                                mode="single"
                                selected={dueDate}
                                onSelect={(date) => date && setValue("due_date", date)}
                                autoFocus
                            />
                        </PopoverContent>
                    </Popover>
                    {errors.due_date && (
                        <p className="text-sm text-red-500">{errors.due_date.message}</p>
                    )}
                </div>

                <div className="space-y-2">
                    <label htmlFor="priority" className="text-sm font-medium">Priority</label>
                    <Select onValueChange={(value) => setValue("priority", value as "p1" | "p2" | "p3" | "p4")} defaultValue={watch("priority")}>
                        <SelectTrigger className="w-full" id="priority">
                            <SelectValue placeholder="Select Priority" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="p1">P1 - Critical</SelectItem>
                            <SelectItem value="p2">P2 - High</SelectItem>
                            <SelectItem value="p3">P3 - Medium</SelectItem>
                            <SelectItem value="p4">P4 - Low</SelectItem>
                        </SelectContent>
                    </Select>
                    {errors.priority && (
                        <p className="text-sm text-red-500">{errors.priority.message}</p>
                    )}
                </div>

                <div className="flex flex-col sm:flex-row gap-4 pb-20">
                    <Button type="button" variant="outline" className="flex-1" onClick={() => navigate(-1)}>Cancel</Button>
                    <Button type="submit" className="flex-1" disabled={loading}>
                        {loading ? (
                            <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving...</>
                        ) : editMode ? "Update Task" : "Create Task"}
                    </Button>
                </div>
            </form>
        </div>
    );
}

export default TaskForm;