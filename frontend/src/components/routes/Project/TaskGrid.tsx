import React, { useEffect, useState } from "react";
import axios from "axios";
import TaskCard from "./TaskCard";

function TaskGrid({ projectId }: { projectId: string }) {

    type Task = {
        id: string;
        title: string;
        due_date?: string;
        description: string;
        priority: string;
        status: string;
        assignee_id?: string;
        assignee?: { name: string };
    };

    const [todoTasks, setTodoTasks] = useState<Task[]>([]);
    const [inProgressTasks, setInProgressTasks] = useState<Task[]>([]);
    const [reviewTasks, setReviewTasks] = useState<Task[]>([]);
    const [completedTasks, setCompleteTasks] = useState<Task[]>([]);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const res = await axios.get(`http://localhost:8080/api/projects/${projectId}/tasks`);
                const tasks: Task[] = res.data.data || [];
                setTodoTasks(tasks.filter((t) => t.status === "todo"));
                setInProgressTasks(tasks.filter((t) => t.status === "in_progress"));
                setReviewTasks(tasks.filter((t) => t.status === "staging_review"));
                setCompleteTasks(tasks.filter((t) => t.status === "done"));
            } catch (err) {
                console.error("Failed to fetch tasks:", err);
            }
        };
        if (projectId) fetchTasks();
    }, [projectId]);

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