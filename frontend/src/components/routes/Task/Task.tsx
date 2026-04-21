import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import Comments from "./Comments";
import { Button } from "@/components/ui/button";
import { Logs, PenBox } from "lucide-react";
import { Link, useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const statusLabels: Record<string, string> = {
    todo: "To Do",
    in_progress: "In Progress",
    staging_review: "Staging Review",
    done: "Completed",
};

function Task() {

    const { id: taskId } = useParams<{ id: string }>();
    const [task, setTask] = useState<any>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchTask = async () => {
            try {
                const res = await axios.get(`http://localhost:8080/api/tasks/${taskId}`);
                setTask(res.data.task);
            } catch (err) {
                console.error("Failed to fetch task:", err);
            }
        };
        if (taskId) fetchTask();
    }, [taskId]);

    const handleEdit = (id: any) => {
        navigate(`/addEdit?id=${id}`);
    };

    if (!task) {
        return (
            <div>
                <Header />
                <div className="px-4 pt-20 pb-5">Loading...</div>
                <Footer />
            </div>
        );
    }

    const dueDate = task.due_date ? new Date(task.due_date).toLocaleDateString("en-GB") : "Not set";

    return (
        <div>
            <Header />
            <div className="px-4 pt-20 pb-5">
                <div className="flex flex-row justify-between">
                    <div className="flex flex-row justify-between">
                        <h1 className="text-4xl font-bold gradient-title">{task.title}</h1>
                        <Button className="mt-1 ml-5" variant={"outline"} onClick={() => handleEdit(task.id)}>
                            <PenBox size={18} /><span className="hidden md:inline">Edit Task</span>
                        </Button>
                    </div>
                    <Link to={`/activity?taskId=${taskId}`}>
                        <Button variant="outline">
                            <Logs size={18} /><span className="hidden md:inline">Activity Log</span>
                        </Button>
                    </Link>
                </div>
                <p className="text-xl text-gray-600 mb-8 max-w-3xl">{task.description}</p>
                <div className="flex items-center mb-5">
                    <h2 className="text-3xl font-bold">Assigned To:</h2>
                    {task.assignee && (
                        <>
                            {task.assignee.avatar_url && <img className="ml-4 rounded-full" src={task.assignee.avatar_url} alt={task.assignee.name} width={40} height={40} />}
                            <div className="text-3xl font-semibold ml-4">{task.assignee.name}</div>
                        </>
                    )}
                    {!task.assignee && <div className="text-3xl font-semibold ml-4">Unassigned</div>}
                </div>
                <div className="text-[#162ab0] text-2xl font-bold mb-5 flex flex-row justify-around">
                    <h2>Status: {statusLabels[task.status] ?? task.status}</h2>
                    <h2>Priority: {task.priority?.toUpperCase()}</h2>
                    <h2>Due Date: {dueDate}</h2>
                </div>
            </div>
            <Comments taskId={taskId!} projectId={task.project_id}/>
            <Footer />
        </div>
    );
}

export default Task;