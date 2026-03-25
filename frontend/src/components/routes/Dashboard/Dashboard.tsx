import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import ProjectGrid from "./ProjectGrid";

function Dashboard() {
    return (
        <div>
            <Header />
            <div className="px-4 py-20">
                <h1 className="text-6xl font-bold gradient-title mb-4">Dashboard</h1>
                <ProjectGrid />
            </div>
            <Footer />
        </div>
    );
}

export default Dashboard;