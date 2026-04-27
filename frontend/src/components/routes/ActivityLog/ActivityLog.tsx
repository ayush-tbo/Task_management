import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import Entries from "./Entries";
import { useSearchParams } from "react-router-dom";

function ActivityLog(){

    const [searchParams] = useSearchParams();
    
    const projectId = searchParams.get("projectId");
    const taskId = searchParams.get("taskId");

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="px-4 py-20 flex-1">
                <h1 className="text-3xl sm:text-5xl font-bold gradient-title mb-4">
                    {taskId ? "Task History" : "Project Activity Log"}
                </h1>
                <Entries projectId={projectId || undefined} taskId={taskId || undefined} />
            </div>
            <Footer />
        </div>
    );
}

export default ActivityLog;