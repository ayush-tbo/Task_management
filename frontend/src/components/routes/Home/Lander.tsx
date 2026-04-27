import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button"

function Lander(){

    return (
        <div className="pb-20 sm:pb-40 px-4 pt-24 sm:pt-35">
            <div className="container mx-auto text-center">
                <h1 className="text-3xl sm:text-5xl md:text-8xl lg:text-[105px] pb-6 gradient-title">
                    Task Management System <br /> Easy to Use
                </h1>
                <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
                    TMS is a platform for Smart Task Management<br /> 
                    Track, Plan, and Build your Tasks Effortlessly.
                </p>
                <Link to="/dashboard">
                    <Button size="lg" className="px-8">Get Started</Button>
                </Link>
            </div>
        </div>
    );
}

export default Lander;