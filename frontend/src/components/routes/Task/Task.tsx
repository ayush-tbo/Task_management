import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { testingTasks } from "@/lib/static";
import Comments from "./Comments";
import { Button } from "@/components/ui/button";
import { Logs, PenBox } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";

function Task() {

    const taskId = "69e29172c6baba2bb6ddb35e"

    const [taskStatus, setTaskStatus] = useState<string>("");

    const navigate = useNavigate();

    useEffect(() => {
        if(testingTasks[0].status == "todo") setTaskStatus("To Do");
        else if(testingTasks[0].status == "inProgress") setTaskStatus("In Progress");
        else if(testingTasks[0].status == "review") setTaskStatus("Staging Review");
        else setTaskStatus("Completed");
    }, [])

    const handleEdit = (id:any) => {
        navigate(`/addEdit?id=${id}`);
    };

    return (
        <div>
            <Header />
            <div className="px-4 pt-20 pb-5">
                <div className="flex flex-row justify-between">
                    <div className="flex flex-row justify-between">
                        <h1 className="text-4xl font-bold gradient-title">{testingTasks[0].title}</h1>
                        <Button className="mt-1 ml-5" variant={"outline"} onClick={() => handleEdit(testingTasks[0].id)}>
                            <PenBox size={18} /><span className="hidden md:inline">Edit Task</span>
                        </Button>
                    </div>
                    <Link to={`/activity?taskId=${taskId}`}>
                        <Button variant="outline">
                            <Logs size={18} /><span className="hidden md:inline">Activity Log</span>
                        </Button>
                    </Link>
                </div>
                <p className="text-xl text-gray-600 mb-8 max-w-3xl">{testingTasks[0].description}</p>
                <div className="flex items-center mb-5">
                    <h2 className="text-3xl font-bold">Assigned To:</h2>
                    <img className="ml-4 rounded-full" src={"https://randomuser.me/api/portraits/men/75.jpg"} alt={"Michael Chen"} width={40} height={40} />
                    <div className="text-3xl font-semibold ml-4">{testingTasks[0].assignedTo}</div>
                </div>
                <div className="text-[#162ab0] text-2xl font-bold mb-5 flex flex-row justify-around">
                    <h2>Status: {taskStatus}</h2>
                    <h2>Priority: {testingTasks[0].priority}</h2>
                    <h2>Due Date: {testingTasks[0].dueDate.toLocaleDateString("en-GB")}</h2>
                </div>
            </div>
            <Comments taskId={"69e29172c6baba2bb6ddb35e"} projectId={"69e3e7eca5eca23605cf7108"}/>
            <Footer />
        </div>
    );
}

export default Task;