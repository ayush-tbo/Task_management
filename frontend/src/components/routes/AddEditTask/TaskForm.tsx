import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { format } from 'date-fns';
import { CalendarIcon, ChevronLeftIcon, ChevronRightIcon, Loader2 } from "lucide-react";
import { Slider } from "@/components/ui/slider";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod/src/zod.js";
import { taskSchema } from "@/lib/schema";
import axios from "axios";
import { useSearchParams } from "react-router-dom";
import { Textarea } from "@/components/ui/textarea";
import { useNavigate } from "react-router-dom";

function TaskForm(){
    
    const [loading, setLoading] = useState(false);
    const [editMode, setEditMode] = useState(false);
    const [initialData, setInitialData] = useState(undefined);

    const [searchParams] = useSearchParams();
    const taskId = searchParams.get("id");

    const navigate = useNavigate();

    useEffect(() => {
        if(taskId){
            // getTask();
            setEditMode(true);
        }
    }, [taskId]);
    
    const { register, handleSubmit, formState:{errors}, setValue, watch, reset } = useForm({
        resolver:zodResolver(taskSchema),
        defaultValues:{
            title: "",
            dueDate: new Date(),
            description: "",
            priority: "6",
            status: "todo",
            assignedTo: "",
        },
    });

    const dueDate = watch("dueDate");
    const priority = watch("priority");

    const onSubmit = async (data: any) => {
        try{
            console.log("Form data:", data); // data => {title:'Project_name', dueDate:Date, description: 'desc', priority: '2', status: 'todo', assignedTo:'Meet K.'}
            const res = await axios.post("http://localhost:8080/proj/new", data);
        }
        catch(err){
            console.error("Failed to create task:", err);
        }
    };

    const getTask = async () => {
        try{
            const res = await axios.post("http://localhost:8080/task/get", {taskId : taskId});
            setInitialData(res.data.transaction);
        }
        catch(err){
            console.error("Failed to get task:", err);
        }
    };

    return (
        <div className="max-w-3xl mx-auto px-5 pt-20">
            <h1 className="text-5xl gradient-title">{editMode ? "Edit" : "Add"} task</h1>
            <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
                <div className="space-y-2">
                    <label htmlFor="title" className="text-sm font-medium">Title</label>
                    <Input id="title" placeholder="e.g., Initialize monorepo repository" {...register("title")} />
                    {errors.title && (
                        <p className="text-sm text-red-500">{errors.title.message}</p>
                    )}
                </div>

                <div className="grid gap-6 md:grid-cols-2">
                    <div className="space-y-2">
                        <label htmlFor="status" className="text-sm font-medium">Task Status</label>
                        <Select onValueChange={(value) => setValue("status", value as "todo" | "inProgress" | "review" | "completed")} defaultValue={watch("status")}>
                            <SelectTrigger className="w-88.75" id="status">
                                <SelectValue placeholder="Select Status" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="todo">To Do</SelectItem>
                                <SelectItem value="inProgress">In Progress</SelectItem>
                                <SelectItem value="review">Review</SelectItem>
                                <SelectItem value="completed">Completed</SelectItem>
                            </SelectContent>
                        </Select>
                        {errors.status && (
                            <p className="text-sm text-red-500">{errors.status.message}</p>
                        )}
                    </div>
                    <div className="space-y-2">
                        <label htmlFor="assignedTo" className="text-sm font-medium">Assigned To</label>
                        <Input id="assignedTo" placeholder="e.g., Meek K." {...register("assignedTo")} />
                        {errors.assignedTo && (
                            <p className="text-sm text-red-500">{errors.assignedTo.message}</p>
                        )}
                    </div>
                </div>

                <div className="space-y-2">
                    <label htmlFor="description" className="text-sm font-medium">Description</label>
                    <Textarea id="description" placeholder="e.g., Set up the base Go backend, React TypeScript frontend, and Docker compose configurations." {...register("description")} />
                    {errors.description && (
                        <p className="text-sm text-red-500">{errors.description.message}</p>
                    )}
                </div>

                <div className="space-y-2">
                    <label className="text-sm font-medium">Due Date</label>
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button variant="outline" className="w-full pl-3 text-left font-normal">
                                {dueDate ? format(dueDate, "PPP") : <span>Pick a Date</span>}
                                <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                            <Calendar
                                mode="single"
                                selected={dueDate}
                                onSelect={(date) => date && setValue("dueDate", date)}
                                autoFocus
                            />
                        </PopoverContent>
                    </Popover>
                    {errors.dueDate && (
                        <p className="text-sm text-red-500">{errors.dueDate.message}</p>
                    )}
                </div>

                <div className="space-y-2">
                    <div className="flex flex-row items-center justify-start gap-2">
                        <label htmlFor="priority" className="text-sm font-medium">Priority</label>
                        <span className="flex flex-row items-center text-sm text-muted-foreground font-bold">
                            <ChevronLeftIcon />
                            {priority}
                            <ChevronRightIcon />
                        </span>
                    </div>
                    <Slider 
                        defaultValue={[6]}
                        max={10}
                        step={1}
                        className="w-full"
                        onValueChange={(value) => setValue("priority", value[0].toString())}
                    />
                    {errors.priority && (
                        <p className="text-sm text-red-500">{errors.priority.message}</p>
                    )}
                </div>

                <div className="flex gap-4 pb-20">
                    <Button type="button" variant="outline" className="w-88.75" onClick={() => navigate(-1)}>Cancel</Button>
                    <Button type="submit" className="w-88.75" disabled={loading}>
                        {loading ? (
                            <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving...</>
                        ) : editMode ? "Update Transaction" : "Create Transaction"}
                    </Button>
                </div>
            </form>
        </div>
    );
}

export default TaskForm;