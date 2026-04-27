import React, { useEffect, useState } from "react";
import axios from "axios";
import TaskCard from "./TaskCard";

type Task = {
    id: string;
    title: string;
    due_date?: string;
    description: string;
    priority: string;
    status: string;
    is_past_due?: boolean;
    assignee_id?: string;
    assignee?: { name: string };
};

const columns = [
    { key: "todo", label: "To Do", color: "bg-slate-100" },
    { key: "in_progress", label: "In Progress", color: "bg-blue-50" },
    { key: "staging_review", label: "Review", color: "bg-amber-50" },
    { key: "done", label: "Done", color: "bg-green-50" },
];

function TaskGrid({ projectId }: { projectId: string }) {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [labels, setLabels] = useState<any[]>([]);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const res = await axios.get(`/api/projects/${projectId}/tasks`);
                setTasks(res.data.data || []);
            } catch (err) {
                console.error("Failed to fetch tasks:", err);
            }
        };
        const fetchLabels = async () => {
            try {
                const res = await axios.get(`/api/projects/${projectId}/labels`);
                setLabels(res.data.labels || []);
            } catch (err) {
                console.error("Failed to fetch labels:", err);
            }
        };
        if (projectId) {
            fetchTasks();
            fetchLabels();
        }
    }, [projectId]);

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {columns.map((col) => {
                const filtered = tasks.filter((t) => t.status === col.key);
                return (
                    <div key={col.key} className={`${col.color} rounded-xl flex flex-col min-h-[24rem]`}>
                        <div className="flex items-center justify-between px-3 py-2 border-b">
                            <span className="text-sm font-semibold">{col.label}</span>
                            <span className="text-xs text-muted-foreground bg-white px-1.5 py-0.5 rounded-full">{filtered.length}</span>
                        </div>
                        <div className="overflow-y-auto no-scrollbar flex-1 p-1">
                            {filtered.map((task) => (
                                <TaskCard key={task.id} task={task} labels={labels} />
                            ))}
                        </div>
                    </div>
                );
            })}
        </div>
    );
}

export default TaskGrid;