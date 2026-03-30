import React from "react";
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button";
import { testingComments } from "@/lib/static";
import { Card, CardContent } from "@/components/ui/card";

function Comments() {
    return (
        <div className="mb-10 bg-[#ecf4f1]">
            <div className="container mx-auto px-4 space-y-1">
                <h2 className="text-3xl font-bold text-center pt-4 mb-5">Comments</h2>
                <Textarea placeholder="Type your comment here..." />
                <Button>Add Comment</Button>
            </div>
            <div className="container mx-auto py-4 px-4 space-y-2">
                {testingComments.map((comment) => (
                    <Card key={comment.id} className="bg-blue-100">
                        <CardContent>
                            <div className="flex justify-between text-sm text-gray-500">
                                <span>{comment.author}</span>
                                <span>{comment.createdAt.toLocaleDateString("en-GB")}</span>
                            </div>
                            <p className="mt-2 text-sm">{comment.text}</p>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}

export default Comments;