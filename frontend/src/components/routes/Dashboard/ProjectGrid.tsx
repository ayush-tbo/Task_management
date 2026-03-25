import React from "react";
import { testingProjects } from "@/lib/static";
import ProjectCard from "./ProjectCard";

function ProjectGrid() {
    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 pt-8">
            {testingProjects.map((project) => {
                return <ProjectCard key={project.id} project={project} />
            })}
        </div>
    );
}

export default ProjectGrid;