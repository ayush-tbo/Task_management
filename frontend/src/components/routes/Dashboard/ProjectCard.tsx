import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useNavigate } from "react-router-dom";
import { Users, ListTodo } from "lucide-react";

function ProjectCard({ project }: any) {
    const { id, name, description, member_count, task_count } = project;
    const navigate = useNavigate();

    return (
        <Card onClick={() => navigate(`/project/${id}`)} className="hover:shadow-md transition-shadow cursor-pointer">
            <CardHeader className="pb-2">
                <CardTitle className="text-xl font-bold capitalize">{name}</CardTitle>
                {description && (
                    <p className="text-sm text-muted-foreground line-clamp-2">{description}</p>
                )}
            </CardHeader>
            <CardContent className="pt-0">
                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    <span className="flex items-center gap-1"><Users size={14} />{member_count || 0} members</span>
                    <span className="flex items-center gap-1"><ListTodo size={14} />{task_count || 0} tasks</span>
                </div>
            </CardContent>
        </Card>
    );
}

export default ProjectCard;