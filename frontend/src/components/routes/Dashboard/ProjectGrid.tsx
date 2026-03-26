import React, { useState } from "react";
import { testingProjects } from "@/lib/static";
import ProjectCard from "./ProjectCard";
import CreateProject from "./CreateProject";

function ProjectGrid() {

    const [open, setOpen] = useState(false);

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 pt-8">
            <CreateProject open={open} setOpen={setOpen} />
            {testingProjects.map((project) => {
                return <ProjectCard key={project.id} project={project} />
            })}
        </div>
    );
}

export default ProjectGrid;