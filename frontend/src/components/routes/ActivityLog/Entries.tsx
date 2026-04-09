import React from "react";
import { testingActivity } from "@/lib/static";
import { Card, CardContent } from "@/components/ui/card";

function Entries(){
    return (
        <div className="container mx-auto py-4 px-4 space-y-2">
            {testingActivity.map((activity) => (
                <Card key={activity.id}>
                    <CardContent>
                        <p className="mt-2">
                            <span className="font-bold">{activity.user}</span>
                            {" "}
                            <span className="font-semibold underline">{activity.activity}</span>
                            {" "}
                            the task
                            {" "}
                            <span className="font-semibold">{activity.taskName}</span>
                            {" "}
                            in project
                            {" "}
                            <span className="font-bold">{activity.projectName}</span>
                            {" "}
                            on
                            {" "}
                            <span className="font-semibold">{activity.date.toLocaleDateString("en-GB")}</span>
                        </p>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}

export default Entries;