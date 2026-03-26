import React, { useEffect, useState } from "react";
import { testingTasks } from "@/lib/static";
import TaskCard from "./TaskCard";

function TaskGrid() {

    type Task = {
        id: number;
        title: string;
        dueDate: Date;
        description: string;
        priority: number;
        status: string;
        assignedTo: string;
    };

    const [todoTasks, setTodoTasks] = useState<Task[]>([]);
    const [inProgressTasks, setInProgressTasks] = useState<Task[]>([]);
    const [reviewTasks, setReviewTasks] = useState<Task[]>([]);
    const [completedTasks, setCompleteTasks] = useState<Task[]>([]);

    useEffect(() => {
        setTodoTasks(testingTasks.filter((t) => t.status === "todo"));
        setInProgressTasks(testingTasks.filter((t) => t.status === "inProgress"));
        setReviewTasks(testingTasks.filter((t) => t.status === "review"));
        setCompleteTasks(testingTasks.filter((t) => t.status === "completed"));
    }, []);

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <div className="bg-blue-100 h-96 hover:shadow-md transition-shadow rounded-xl flex flex-col text-black">
                <div className="font-bold text-center pt-1">To Do</div>
                <div className="overflow-y-auto no-scrollbar flex-1">
                    {todoTasks.map((task) => {
                        return <TaskCard key={task.id} task={task} />
                    })}
                </div>
            </div>
            <div className="bg-blue-100 h-96 hover:shadow-md transition-shadow rounded-xl flex flex-col text-black">
                <div className="font-bold text-center pt-1">In Progress</div>
                <div className="overflow-y-auto no-scrollbar flex-1">
                    {inProgressTasks.map((task) => {
                        return <TaskCard key={task.id} task={task} />
                    })}
                </div>
            </div>
            <div className="bg-blue-100 h-96 hover:shadow-md transition-shadow rounded-xl flex flex-col text-black">
                <div className="font-bold text-center pt-1">Staging Review</div>
                <div className="overflow-y-auto no-scrollbar flex-1">
                    {reviewTasks.map((task) => {
                        return <TaskCard key={task.id} task={task} />
                    })}
                </div>
            </div>
            <div className="bg-blue-100 h-96 hover:shadow-md transition-shadow rounded-xl flex flex-col text-black">
                <div className="font-bold text-center pt-1">Completed</div>
                <div className="overflow-y-auto no-scrollbar flex-1">
                    {completedTasks.map((task) => {
                        return <TaskCard key={task.id} task={task} />
                    })}
                </div>
            </div>
        </div>
    );
}

export default TaskGrid;