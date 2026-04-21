import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { useNavigate } from "react-router-dom";

function TaskCard({ task } : any) {

    const { id, title, due_date, description, priority, status, assignee } = task;

    const date = due_date ? new Date(due_date) : null;

    const navigate = useNavigate();

    return (
        <div className="px-2 py-1">
            <Card onClick={() => navigate(`/task/${id}`)} className="hover:shadow-md transition-shadow group relative cursor-pointer">
                <CardContent className="my-0 py-0">
                    <div className="text-sm font-semibold pb-1">{title}</div>
                    <div className="flex flex-col justify-between space-y-0">
                        {date && <p className="text-xs text-muted-foreground font-bold">Due Date: {date.toLocaleDateString("en-GB")}</p>}
                        <p className="text-xs text-black truncate bg-green-300">{description}</p>
                        <p className="text-xs">Assigned To: {assignee?.name ?? "Unassigned"}</p>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}

export default TaskCard;