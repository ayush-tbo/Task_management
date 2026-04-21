import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import TaskGrid from "./TaskGrid";
import { Link, useParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Logs, Plus } from "lucide-react";
import axios from "axios";

function Project() {
    const { id: projectId } = useParams<{ id: string }>();
    const [project, setProject] = useState<any>(null);

    useEffect(() => {
        const fetchProject = async () => {
            try {
                const res = await axios.get(`http://localhost:8080/api/projects/${projectId}`);
                setProject(res.data.project);
            } catch (err) {
                console.error("Failed to fetch project:", err);
            }
        };
        if (projectId) fetchProject();
    }, [projectId]);

    return (
        <div>
            <Header />
            <div className="px-4 pt-20 pb-5">
                <div className="flex flex-row justify-between">
                    <h1 className="text-6xl font-bold gradient-title mb-4">{project?.name ?? "Loading..."}</h1>
                    <div className="flex gap-2 mt-2">
                        <Link to={`/addEdit?projectId=${projectId}`}>
                            <Button variant="outline">
                                <Plus size={18} /><span className="hidden md:inline">Add Task</span>
                            </Button>
                        </Link>
                        <Link to={`/activity?projectId=${projectId}`}>
                            <Button variant="outline">
                                <Logs size={18} /><span className="hidden md:inline">Activity Log</span>
                            </Button>
                        </Link>
                    </div>
                </div>
                <TaskGrid projectId={projectId!} />
            </div>
            <Footer />
        </div>
    );
}

export default Project;