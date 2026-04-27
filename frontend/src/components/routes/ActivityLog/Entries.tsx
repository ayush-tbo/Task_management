import React, { useEffect, useState } from "react";
import axios from "axios";
import { MessageSquare, Trash2, Edit, History, UserPlus2, UserMinus, CheckSquare, Pencil, Zap, Tag } from "lucide-react";

const actionMessages: Record<string, string> = {
    comment_added: "added a comment",
    comment_changed: "edited a comment",
    comment_deleted: "deleted a comment",
    member_added: "added a member",
    member_removed: "removed a member",
    task_created: "created a task",
    task_updated: "updated a task",
    task_deleted: "deleted a task",
    task_assigned: "assigned a task",
    status_changed: "changed task status",
    sprint_created: "created a sprint",
    sprint_updated: "updated a sprint",
    label_created: "created a label",
    label_deleted: "deleted a label",
};

const actionIcons: Record<string, React.ReactNode> = {
    comment_added: <MessageSquare className="w-4 h-4 text-blue-500" />,
    comment_changed: <Edit className="w-4 h-4 text-amber-500" />,
    comment_deleted: <Trash2 className="w-4 h-4 text-red-500" />,
    member_added: <UserPlus2 className="w-4 h-4 text-green-500" />,
    member_removed: <UserMinus className="w-4 h-4 text-orange-500" />,
    task_created: <CheckSquare className="w-4 h-4 text-indigo-500" />,
    task_updated: <Pencil className="w-4 h-4 text-purple-500" />,
    task_deleted: <Trash2 className="w-4 h-4 text-rose-500" />,
    task_assigned: <UserPlus2 className="w-4 h-4 text-teal-500" />,
    status_changed: <CheckSquare className="w-4 h-4 text-sky-500" />,
    sprint_created: <Zap className="w-4 h-4 text-green-600" />,
    sprint_updated: <Zap className="w-4 h-4 text-amber-600" />,
    label_created: <Tag className="w-4 h-4 text-violet-500" />,
    label_deleted: <Tag className="w-4 h-4 text-red-400" />,
};

function Entries({ projectId, taskId }: { projectId?: string; taskId?: string }) {
    const [activities, setActivities] = useState<any[]>([]);

    useEffect(() => {
        const fetchActivity = async () => {
            try {
                if (projectId) {
                    const res = await axios.get(`/api/projects/${projectId}/activity`);
                    setActivities(res.data.activities || []);
                } else if (taskId) {
                    const res = await axios.get(`/api/tasks/${taskId}/activity`);
                    setActivities(res.data.activities || []);
                }
            } catch (err) {
                console.error("Failed to fetch activity:", err);
            }
        };
        if (projectId || taskId) fetchActivity();
    }, [projectId, taskId]);

    if (activities.length === 0) {
        return <p className="text-sm text-muted-foreground text-center py-8">No activity recorded yet.</p>;
    }

    return (
        <div className="max-w-2xl mx-auto">
            <div className="relative border-l-2 border-slate-200 ml-4">
                {activities.map((activity: any) => (
                    <div key={activity.id} className="relative pl-8 pb-6">
                        <div className="absolute left-[-9px] top-1 w-4 h-4 rounded-full bg-white border-2 border-slate-300 flex items-center justify-center">
                            {actionIcons[activity.action] || <History className="w-3 h-3 text-slate-400" />}
                        </div>
                        <div>
                            <div className="flex items-center gap-2 flex-wrap">
                                <span className="font-semibold text-sm">{activity.user?.name || "Unknown"}</span>
                                <span className="text-sm text-slate-600">
                                    {actionMessages[activity.action] || activity.action.replaceAll("_", " ")}
                                </span>
                            </div>
                            {activity.details?.old_content && (
                                <div className="mt-1 text-xs text-slate-500 bg-slate-50 p-2 rounded border-l-2 border-amber-300">
                                    Previous: "{activity.details.old_content}"
                                </div>
                            )}
                            {activity.details?.info && (
                                <p className="text-xs text-slate-400 italic mt-0.5">{activity.details.info}</p>
                            )}
                            <p className="text-[10px] text-slate-400 mt-1">
                                {new Date(activity.created_at).toLocaleString("en-GB", {
                                    day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit",
                                })}
                            </p>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default Entries;