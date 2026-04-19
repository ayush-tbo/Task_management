import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { LayoutDashboard, PenBox } from "lucide-react";
import Profile from "./Profile";

function Header() {
    return (
        <div className="fixed top-0 w-full bg-white/80 backdrop-blur-md z-50 border-b">
            <nav className="container flex items-center mx-auto px-4 py-4 justify-between">
                <Link to="/">
                    <img className="size-11 w-auto opacity-100 object-contain" src="/Header/Jira.svg" alt="Jira Logo" />
                </Link>
                <div className="flex items-center space-x-4">
                    <Link to="/dashboard" className="text-gray-600 hover:text-blue-800 flex items-center gap-2">
                        <Button variant="outline">
                            <LayoutDashboard size={18} /><span className="hidden md:inline">Dashboard</span>
                        </Button>
                    </Link>
                    <Link to="/addEdit" className="flex items-center gap-2">
                        <Button>
                            <PenBox size={18} /><span className="hidden md:inline">Add Task</span>
                        </Button>
                    </Link>
                    <Profile />
                </div>
            </nav>
        </div>
    );
}

export default Header;