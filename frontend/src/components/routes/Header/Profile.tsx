import React from "react";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

function Profile() {
    return (
        <div>
            <Link to="/login">
                <Button variant="outline">Log in</Button>
            </Link>
        </div>
    );
}

export default Profile;