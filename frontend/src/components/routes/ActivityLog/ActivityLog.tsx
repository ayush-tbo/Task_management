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
        <div>
            <Header />
            <div className="px-4 py-20">
                <h1 className="text-6xl font-bold gradient-title mb-4">
                    {taskId ? "Task History" : "Project Activity Log"}
                </h1>
                <Entries projectId={projectId || undefined} taskId={taskId || undefined} />
            </div>
            <Footer />
        </div>
    );
}

export default ActivityLog;