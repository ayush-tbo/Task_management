import { z } from "zod";

export const projectNameSchema = z.object({
    name: z.string().min(1, "Name is required")
});

export const taskSchema = z.object({
    title: z.string().min(1, "Title is required"),
    dueDate: z.date({required_error: "Due Date is required"}),
    description: z.string().optional(),
    priority: z.string().min(1, "Priority Number is required"),
    status: z.enum(["todo", "inProgress", "review", "completed"]),
    assignedTo: z.string().optional(),
});