import React from "react";
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Card, CardContent } from "@/components/ui/card";
import { Plus } from "lucide-react";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { projectNameSchema } from "@/lib/schema";
import axios from "axios";

function CreateProject({open, setOpen, onCreated}: any) {

    const { register, handleSubmit, formState:{errors}, reset } = useForm({
        resolver:zodResolver(projectNameSchema),
        defaultValues:{
            name: "",
            description: "",
        },
    });

    const onSubmit = async (data: any) => {
        try{
            console.log("Form data:", data);
            const res = await axios.post("/api/projects", data);
            reset();
            setOpen(false);
            onCreated?.();
        }
        catch(err){
            console.error("Failed to create project:", err);
        }
    };

    return (
        <div>
            <Dialog open={open} onOpenChange={setOpen}>
                <DialogTrigger asChild>
                    <Card className="h-39 hover:shadow-md transition-shadow cursor-pointer border-dashed">
                        <CardContent className="flex flex-col items-center justify-center text-muted-foreground h-full pt-5">
                            <Plus className="h-10 w-10 mb-2"/>
                            <p className="text-sm font-medium">Add New Project</p>
                        </CardContent>
                    </Card>
                </DialogTrigger>
                <DialogContent className="sm:max-w-sm" aria-describedby={undefined}>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <DialogHeader>
                            <DialogTitle className="pb-4">Create New Project</DialogTitle>
                        </DialogHeader>
                        <div className="space-y-2">
                            <label htmlFor="name" className="block text-sm font-medium">Project Name</label>
                            <Input id="name" {...register("name")} />
                            {errors.name && (
                                <p className="text-sm text-red-500">{errors.name.message}</p>
                            )}
                        </div>
                        <div className="space-y-2">
                            <label htmlFor="description" className="block text-sm font-medium">Description</label>
                            <Textarea id="description" placeholder="What is this project about?" {...register("description")} />
                        </div>
                        <DialogFooter>
                            <DialogClose asChild>
                                <Button variant="outline">Cancel</Button>
                            </DialogClose>
                            <Button type="submit">Submit</Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>
        </div>
    );
}

export default CreateProject;