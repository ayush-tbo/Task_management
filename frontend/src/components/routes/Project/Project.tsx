import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import TaskGrid from "./TaskGrid";
import AddEditTask from "../AddEditTask/AddEditTask";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Logs } from "lucide-react";

function Project() {
    const projectId = "69e3e7eca5eca23605cf7108"
    return (
        <div>
            <Header />
            <div className="px-4 pt-20 pb-5">
                <div className="flex flex-row justify-between">
                    <h1 className="text-6xl font-bold gradient-title mb-4">Project_Name</h1>
                    <Link to={`/activity?projectId=${projectId}`} className="mt-2">
                        <Button variant="outline">
                            <Logs size={18} /><span className="hidden md:inline">Activity Log</span>
                        </Button>
                    </Link>
                </div>
                <TaskGrid />
            </div>
            <Footer />
        </div>
    );
}

export default Project;