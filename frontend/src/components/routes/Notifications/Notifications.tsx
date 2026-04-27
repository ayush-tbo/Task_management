import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { Button } from "@/components/ui/button";
import { AtSign, Bell, CheckCheck, History, Info, UserPlus, MessageSquare, RefreshCw } from "lucide-react";
import axios from "axios";

function Notifications() {

    const [notifications, setNotifications] = useState([]);

    const getIcon = (type: string) => {
        switch (type) {
            case "reminder": return <Info className="h-4 w-4 text-blue-500" />;
            case "alert": return <RefreshCw className="h-4 w-4 text-orange-500" />;
            case "mention": return <MessageSquare className="h-4 w-4 text-green-500" />;
            case "assignment": return <UserPlus className="h-4 w-4 text-purple-500" />;
            default: return <Bell className="w-4 h-4 text-slate-500" />;
        }
    };

    const handleGetNotifications = async () => {
        try {
            const res = await axios.get(`/api/notifications`);
            setNotifications(res.data.notifications);
            console.log(res.data.notifications);
        }
        catch(err){
            console.error("Failed to load notifications", err);
        }
    };

    const handleMarkRead = async (id: string) => {
        try{
            await axios.put(`/api/notifications/${id}/read`);
            await handleGetNotifications();
        }
        catch(err){
            console.error("Failed to mark read notification", err);
        }
    };

    const handleMarkAllRead = async () => {
        try{
            await axios.put(`/api/notifications/read-all`);
            await handleGetNotifications();
        }
        catch(err){
            console.error("Failed to mark read all notification", err);
        }
    };

    useEffect(() => {
        handleGetNotifications();
    }, []);

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="px-4 py-20 flex-1">
                <div className="flex flex-col sm:flex-row justify-between gap-2">
                    <h4 className="text-3xl sm:text-5xl font-bold gradient-title mb-4">Notifications</h4>
                    <Button variant="ghost" className="text-blue-600 hover:text-blue-700 mt-2" onClick={handleMarkAllRead}>
                        <CheckCheck size={18} /><span className="hidden md:inline">Mark all read</span>
                    </Button>
                </div>
                {notifications.length === 0 ? (
                    <div className="flex h-full items-center justify-center p-8 text-center text-sm text-slate-500">No notifications yet.</div>
                ) : (
                    <div className="flex flex-col">
                        {notifications.map((notification : any) => (
                            <div  key={notification.id}  onClick={() => !notification.is_read && handleMarkRead(notification.id)} className={`flex cursor-pointer items-start gap-3 border-b px-4 py-3 transition-colors hover:bg-slate-50 ${!notification.is_read ? `bg-blue-50/50` : ``}`}>
                                <div className="mt-1">{getIcon(notification.type)}</div>
                                <div className="flex-1 space-y-1">
                                    <p className={`text-sm leading-none ${!notification.is_read ? `font-bold` : `font-medium`}`}>{notification.title}</p>
                                    <p className="text-xs text-slate-500 line-clamp-2">{notification.message}</p>
                                    <p className="text-[10px] text-slate-400">{new Date(notification.created_at).toLocaleDateString("en-GB", { hour: "2-digit", minute: "2-digit" })}</p>
                                </div>
                                {!notification.is_read && (
                                    <div className="mt-2 h-2 w-2 rounded-full bg-blue-600" />
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>
            <Footer />
        </div>
    );
}

export default Notifications;