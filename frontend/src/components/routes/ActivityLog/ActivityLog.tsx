import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import Entries from "./Entries";

function ActivityLog(){
    return (
        <div>
            <Header />
            <div className="px-4 py-20">
                <h1 className="text-6xl font-bold gradient-title mb-4">Activity Log</h1>
                <Entries />
            </div>
            <Footer />
        </div>
    );
}

export default ActivityLog;