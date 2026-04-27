import React, { useEffect, useState } from "react";
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { commentSchema } from "@/lib/schema";
import axios from "axios";
import { useAuth } from "@/context/AuthContext";
import { Edit, Trash2, Check, X } from "lucide-react";

function Comments({ taskId, projectId } : any) {

    const { user } = useAuth();

    const [comments, setComments] = useState([]);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [editContent, setEditContent] = useState("");

    const { register, handleSubmit, formState:{errors}, reset} = useForm({
        resolver:zodResolver(commentSchema),
    });

    useEffect(() => {
        handleGetComments();
    }, [taskId]);

    const handleGetComments = async () => {
        try {
            const res = await axios.get(`/api/tasks/${taskId}/comments`);
            setComments(res.data.comments);
        }
        catch(err){
            console.error("Error getting comments", err);
        }
    };

    const handleAddComment = async (data: any) => {
        const payload = {...data, project_id: projectId}
        try{
            await axios.post(`/api/tasks/${taskId}/comments`, payload);
            handleGetComments();
            reset();
        }
        catch(err){
            console.error("Create comment failed", err);
        }
    };

    const startEditing = (comment: any) => {
        setEditingId(comment.id);
        setEditContent(comment.content);
    };

    const cancelEditing = () => {
        setEditingId(null);
        setEditContent("");
    };

    const handleSaveEdit = async (commentId: string) => {
        if (!editContent.trim()) return;
        try {
            await axios.put(`/api/comments/${commentId}`, { content: editContent, task_id: taskId, project_id: projectId });
            setEditingId(null);
            setEditContent("");
            handleGetComments();
        } catch (err) {
            console.error("Update failed", err);
        }
    };

    const handleDelete = async (commentId: string) => {
        if (!window.confirm("Delete this comment?")) return;
        try{
            await axios.delete(`/api/comments/${commentId}`, {data: { task_id: taskId, project_id: projectId }});
            handleGetComments();
        }
        catch(err){
            console.error("Delete failed", err);
        }
    };

    return (
        <div className="mb-10 bg-[#ecf4f1]">
            <div className="container mx-auto px-4 space-y-1">
                <h2 className="text-3xl font-bold text-center pt-4 mb-5">Comments</h2>
                <form onSubmit={handleSubmit(handleAddComment)}>
                    <Textarea id="content" placeholder="Add a comment..." {...register("content")}/>
                    {errors.content && (
                        <p className="text-sm text-red-500">{errors.content.message as string}</p>
                    )}
                    <Button type="submit" className="mt-2">Add Comment</Button>
                </form>
            </div>
            <div className="container mx-auto py-4 px-4 space-y-2">
                {comments.map((comment: any) => (
                    <Card key={comment.id} className="bg-blue-100">
                        <CardContent>
                            <div className="flex justify-between items-start mb-2">
                                <div className="flex flex-col">
                                    <span className="font-semibold text-sm text-slate-900">{comment.user?.name || "Anonymous"}</span>
                                    <span className="text-[10px] text-slate-400">
                                        {new Date(comment.updated_at).toLocaleString("en-GB", { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" })}
                                    </span> 
                                </div>

                                {user?.id == comment.user_id && editingId !== comment.id && (
                                    <div className="flex items-center gap-1">
                                        <Button variant="outline" size="icon" className="h-8 w-8 text-slate-500 hover:text-blue-600" onClick={() => startEditing(comment)}>
                                            <Edit className="h-4 w-4" />
                                        </Button>
                                        <Button variant="outline" size="icon" className="h-8 w-8 text-slate-500 hover:text-red-600" onClick={() => handleDelete(comment.id)}>
                                            <Trash2 className="h-4 w-4" />
                                        </Button>
                                    </div>
                                )}
                            </div>

                            {editingId === comment.id ? (
                                <div className="space-y-2">
                                    <Textarea
                                        value={editContent}
                                        onChange={(e) => setEditContent(e.target.value)}
                                        className="bg-white"
                                        autoFocus
                                    />
                                    <div className="flex gap-2">
                                        <Button size="sm" onClick={() => handleSaveEdit(comment.id)}>
                                            <Check size={14} className="mr-1" />Save
                                        </Button>
                                        <Button size="sm" variant="outline" onClick={cancelEditing}>
                                            <X size={14} className="mr-1" />Cancel
                                        </Button>
                                    </div>
                                </div>
                            ) : (
                                <p className="text-sm text-slate-700 leading-relaxed">{comment.content}</p>
                            )}
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}

export default Comments;