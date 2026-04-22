import React, { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import axios from "axios";
import { MessageSquare, Trash2, Edit, History, UserPlus2, CheckSquare, Pencil } from "lucide-react";

function Entries({ projectId, taskId }: any){

    const [activities, setActivities] = useState([]);

    const getActionMessage = (activity: any) => {
        const {action, details} = activity

        const actions: Record<string, string> = {
            comment_added: "added a comment",
            comment_changed: "edited a comment",
            comment_deleted: "deleted a comment",
            member_added: "added a member",
            task_created: "created a task",
            task_updated: "updated a task",
            task_deleted: "deleted a task",
        };

        return actions[action] || action.replace("_", " ");
    }

    const getIcon = (action: string) => {
        switch (action) {
            case "comment_added": return <MessageSquare className="w-4 h-4 text-blue-500" />;
            case "comment_deleted": return <Trash2 className="w-4 h-4 text-red-500" />;
            case "comment_changed": return <Edit className="w-4 h-4 text-amber-500" />;
            case "member_added": return <UserPlus2 className="w-4 h-4 text-green-500" />;
            case "task_created": return <CheckSquare className="w-4 h-4 text-indigo-500" />;
            case "task_updated": return <Pencil className="w-4 h-4 text-purple-500" />;
            case "task_deleted": return <Trash2 className="w-4 h-4 text-rose-500" />;
            default: return <History className="w-4 h-4 text-slate-500" />;
        }
    }

    const handleGetActivity = async () => {
        try{
            if(projectId){
                const res = await axios.get(`http://localhost:8080/api/projects/${projectId}/activity`);
                setActivities(res.data.activities);
            }
            else if(taskId){
                const res = await axios.get(`http://localhost:8080/api/tasks/${taskId}/activity`);
                setActivities(res.data.activities);
            }
        }
        catch(err){
            console.error("Failed to fetch activity", err);
        }
    };

    useEffect(() => {
        if(projectId || taskId){
            handleGetActivity();
        }
    }, [projectId, taskId]);

    return (
        <div className="container mx-auto px-4 space-y-2">
            {activities?.map((activity: any) => (
                <Card key={activity.id}>
                    <CardContent>
                        <div className="mt-2">
                            <div className="flex flex-col space-y-1">
                                <div className="flex items-center gap-2">
                                    {getIcon(activity.action)}
                                    <span className="font-bold text-slate-900">{activity.user?.name}</span>
                                    <span className="text-slate-600">{getActionMessage(activity)}</span>
                                </div>
                                {activity.details?.old_content && (
                                    <div className="text-xs italic text-slate-400 bg-slate-50 p-2 rounded border-l-2 border-amber-300">
                                        Prev: "{activity.details.old_content}"
                                    </div>
                                )}
                                {activity.details?.info && (
                                    <span className="text-xs text-slate-400 italic">({activity.details.info})</span>
                                )}
                                <div className="text-[10px] font-medium text-slate-400 uppercase tracking-wider">
                                    {new Date(activity.created_at).toLocaleString("en-GB", { day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" })}
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}

export default Entries;