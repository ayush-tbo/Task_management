import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import CreateTask from "./CreateTask";
import TaskGrid from "./TaskGrid";

function Project() {
    return (
        <div>
            <Header />
            <div className="px-4 pt-20 pb-5">
                <div className="flex flex-row justify-between">
                    <h1 className="text-6xl font-bold gradient-title mb-4">Project_Name</h1>
                    <CreateTask />
                </div>
                <TaskGrid />
            </div>
            <Footer />
        </div>
    );
}

export default Project;