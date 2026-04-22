import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useNavigate } from "react-router-dom";

function ProjectCard({ project }: any) {

    const { id, name, member_count } = project;

    const navigate = useNavigate();

    return (
        <div>
            <Card onClick={() => navigate(`/project/${id}`)} className="hover:shadow-md transition-shadow group relative cursor-pointer">
                <CardHeader className="mb-0 py-0">
                    <CardTitle className=" text-2xl font-bold capitalize">{name}</CardTitle>
                </CardHeader>
                <CardContent className="my-0 py-0">
                    <div className="flex flex-col justify-between space-y-1">
                        <p className="text-sm font-medium">Members: {member_count}</p>
                    </div>
                    <p className="text-sm text-muted-foreground mt-3">Click the card to Edit or View the Board in detail</p>
                </CardContent>
            </Card>
        </div>
    );
}

export default ProjectCard;