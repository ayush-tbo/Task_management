import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { useNavigate } from "react-router-dom";
import { AlertTriangle, Clock } from "lucide-react";

const priorityColors: Record<string, string> = {
    p1: "bg-red-100 text-red-700",
    p2: "bg-orange-100 text-orange-700",
    p3: "bg-blue-100 text-blue-700",
    p4: "bg-slate-100 text-slate-600",
};

function TaskCard({ task, labels = [] }: { task: any; labels?: any[] }) {
    const { id, title, due_date, priority, assignee, is_past_due, label_ids } = task;
    const date = due_date ? new Date(due_date) : null;
    const navigate = useNavigate();
    const taskLabels = labels.filter((l: any) => label_ids?.includes(l.id));

    return (
        <div className="px-1.5 py-1">
            <Card
                onClick={() => navigate(`/task/${id}`)}
                className="hover:shadow-md transition-shadow cursor-pointer border-slate-200"
            >
                <CardContent className="p-3 space-y-1.5">
                    {taskLabels.length > 0 && (
                        <div className="flex flex-wrap gap-1">
                            {taskLabels.map((l: any) => (
                                <span
                                    key={l.id}
                                    className="text-[9px] font-medium px-1.5 py-0.5 rounded text-white"
                                    style={{ backgroundColor: l.color || "#6b7280" }}
                                >
                                    {l.name}
                                </span>
                            ))}
                        </div>
                    )}
                    <div className="flex items-start justify-between gap-2">
                        <span className="text-sm font-medium leading-tight">{title}</span>
                        {priority && (
                            <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded shrink-0 ${priorityColors[priority] || ""}`}>
                                {priority.toUpperCase()}
                            </span>
                        )}
                    </div>
                    <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>{assignee?.name ?? "Unassigned"}</span>
                        {date && (
                            <span className={`flex items-center gap-1 ${is_past_due ? "text-red-600 font-medium" : ""}`}>
                                {is_past_due ? <AlertTriangle size={12} /> : <Clock size={12} />}
                                {date.toLocaleDateString("en-GB", { day: "2-digit", month: "short" })}
                            </span>
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}

export default TaskCard;