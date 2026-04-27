import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import ProjectGrid from "./ProjectGrid";
import { useAuth } from "@/context/AuthContext";

function Dashboard() {
    const { user } = useAuth();

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="px-4 py-20 flex-1">
                <div className="mb-6">
                    <h1 className="text-2xl sm:text-4xl font-bold">
                        Welcome back{user?.name ? `, ${user.name}` : ""}
                    </h1>
                    <p className="text-muted-foreground mt-1">Here are your projects.</p>
                </div>
                <ProjectGrid />
            </div>
            <Footer />
        </div>
    );
}

export default Dashboard;