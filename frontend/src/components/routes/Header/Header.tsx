import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Bell, LayoutDashboard, ListChecks } from "lucide-react";
import Profile from "./Profile";
import axios from "axios";

function Header() {
    const [unreadCount, setUnreadCount] = useState(0);

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (!token) return;

        const fetchCount = async () => {
            try {
                const res = await axios.get("/api/notifications");
                const notifications = res.data.notifications || [];
                setUnreadCount(notifications.filter((n: any) => !n.is_read).length);
            } catch {}
        };
        fetchCount();
        const interval = setInterval(fetchCount, 30000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="fixed top-0 w-full bg-white/80 backdrop-blur-md z-50 border-b">
            <nav className="flex items-center px-4 py-3 justify-between gap-2">
                <Link to="/" className="shrink-0">
                    <img className="size-11 w-auto opacity-100 object-contain" src="/Header/fq.png" alt="FloQast Logo" />
                </Link>
                <div className="flex items-center gap-1 sm:gap-3 overflow-x-auto">
                    <Link to="/my-tasks" className="text-gray-600 hover:text-blue-800 shrink-0">
                        <Button variant="outline" size="sm" className="sm:size-auto">
                            <ListChecks size={18} /><span className="hidden sm:inline">My Tasks</span>
                        </Button>
                    </Link>
                    <Link to="/dashboard" className="text-gray-600 hover:text-blue-800 shrink-0">
                        <Button variant="outline" size="sm" className="sm:size-auto">
                            <LayoutDashboard size={18} /><span className="hidden sm:inline">Dashboard</span>
                        </Button>
                    </Link>
                    <Link to="/notifications" className="text-gray-600 hover:text-blue-800 relative shrink-0">
                        <Button size="sm" className="sm:size-auto">
                            <Bell size={18} /><span className="hidden sm:inline">Notifications</span>
                        </Button>
                        {unreadCount > 0 && (
                            <span className="absolute -top-1 -right-1 bg-red-500 text-white text-[10px] font-bold w-5 h-5 flex items-center justify-center rounded-full">
                                {unreadCount > 9 ? "9+" : unreadCount}
                            </span>
                        )}
                    </Link>
                    <Profile />
                </div>
            </nav>
        </div>
    );
}

export default Header;