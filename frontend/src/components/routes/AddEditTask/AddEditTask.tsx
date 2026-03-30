import React from "react";
import Header from "../Header/Header";
import TaskForm from "./TaskForm";
import Footer from "../Header/Footer";

function AddEditTask() {
    return (
        <div>
            <Header />
            <TaskForm />
            <Footer />
        </div>
    );
}

export default AddEditTask;