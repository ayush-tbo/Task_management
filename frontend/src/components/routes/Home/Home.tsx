import React from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import Lander from "./Lander";

function Home() {
    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="flex-1">
                <Lander />
            </div>
            <Footer />
        </div>
    )
}

export default Home;