import React, { useEffect, useState } from "react";
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod/src/zod.js";
import { commentSchema } from "@/lib/schema";
import axios from "axios";
import { useAuth } from "@/context/AuthContext";
import { Edit, Trash2 } from "lucide-react";

function Comments({ taskId, projectId } : any) {

    const { user } = useAuth();

    const [comments, setComments] = useState([]);

    const { register, handleSubmit, formState:{errors}, reset} = useForm({
        resolver:zodResolver(commentSchema),
    });

    useEffect(() => {
        handleGetComments();
    }, [taskId]);

    const handleGetComments = async () => {
        try {
            const res = await axios.get(`http://localhost:8080/api/tasks/${taskId}/comments`);
            setComments(res.data.comments);
        }
        catch(err){
            console.error("Error getting comments", err);
        }
    };

    const handleAddComment = async (data: any) => {
        const payload = {...data, project_id: projectId}
        try{
            const res = await axios.post(`http://localhost:8080/api/tasks/${taskId}/comments`, payload);
            handleGetComments();
            reset();
        }
        catch(err){
            console.error("Create comment failed", err);
        }
    };

    const handleUpdate = async (commentId: string) => {
        const newContent = prompt("Edit your comment:");
        if (!newContent) return;
        try{
            const res = await axios.put(`http://localhost:8080/api/comments/${commentId}`, { content: newContent, task_id: taskId, project_id: projectId });
            handleGetComments();
        }
        catch(err){
            console.log("Update failed", err);
        }
    };

    const handleDelete = async (commentId: string) => {
        if (!window.confirm("Delete this comment?")) return;
        try{
            const res = await axios.delete(`http://localhost:8080/api/comments/${commentId}`, {data: { task_id: taskId, project_id: projectId }});
            handleGetComments();
        }
        catch(err){
            console.log("Delete failed", err);
        }
    };

    return (
        <div className="mb-10 bg-[#ecf4f1]">
            <div className="container mx-auto px-4 space-y-1">
                <h2 className="text-3xl font-bold text-center pt-4 mb-5">Comments</h2>
                <form onSubmit={handleSubmit(handleAddComment)}>
                    <Textarea id="content" placeholder="Type your comment here..." {...register("content")}/>
                    {errors.content && (
                        <p className="text-sm text-red-500">{errors.content.message as string}</p>
                    )}
                    <Button type="submit">Add Comment</Button>
                </form>
            </div>
            <div className="container mx-auto py-4 px-4 space-y-2">
                {comments.map((comment: any) => (
                    <Card key={comment.id} className="bg-blue-100">
                        <CardContent>
                            <div className="flex justify-between items-start mb-2">
                                <div className="flex flex-col">
                                    <span className="font-semibold text-sm text-slate-900">{comment.user?.name || "Anonymous"}</span>
                                    <span className="text-[10px] text-slate-400">{new Date(comment.updated_at).toLocaleDateString("en-GB")}</span> 
                                </div>

                                {user?.id == comment.user_id && (
                                    <div className="flex items-center gap-1">
                                        <Button variant="outline" size="icon" className="h-8 w-8 text-slate-500 hover:text-blue-600" onClick={() => handleUpdate(comment.id)}>
                                            <Edit className="h-4 w-4" />
                                        </Button>
                                        <Button variant="outline" size="icon" className="h-8 w-8 text-slate-500 hover:text-red-600" onClick={() => handleDelete(comment.id)}>
                                            <Trash2 className="h-4 w-4" />
                                        </Button>
                                    </div>
                                )}
                            </div>
                            <p className="text-sm text-slate-700 leading-relaxed">{comment.content}</p>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}

export default Comments;