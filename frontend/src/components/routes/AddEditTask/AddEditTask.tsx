import React from "react";
import Header from "../Header/Header";
import TaskForm from "./TaskForm";
import Footer from "../Header/Footer";

function AddEditTask() {
    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="flex-1">
                <TaskForm />
            </div>
            <Footer />
        </div>
    );
}

export default AddEditTask;