import React, { useEffect, useState } from "react";
import axios from "axios";
import ProjectCard from "./ProjectCard";
import CreateProject from "./CreateProject";

function ProjectGrid() {

    const [open, setOpen] = useState(false);
    const [projects, setProjects] = useState<any[]>([]);

    const fetchProjects = async () => {
        try {
            const res = await axios.get("http://localhost:8080/api/projects");
            setProjects(res.data.data || []);
        } catch (err) {
            console.error("Failed to fetch projects:", err);
        }
    };

    useEffect(() => {
        fetchProjects();
    }, []);

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 pt-8">
            <CreateProject open={open} setOpen={setOpen} onCreated={fetchProjects} />
            {projects.map((project) => {
                return <ProjectCard key={project.id} project={project} />
            })}
        </div>
    );
}

export default ProjectGrid;