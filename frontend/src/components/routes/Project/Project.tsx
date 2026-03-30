import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import TaskGrid from "./TaskGrid";
import AddEditTask from "../AddEditTask/AddEditTask";

function Project() {
    return (
        <div>
            <Header />
            <div className="px-4 pt-20 pb-5">
                <h1 className="text-6xl font-bold gradient-title mb-4">Project_Name</h1>
                <TaskGrid />
            </div>
            <Footer />
        </div>
    );
}

export default Project;