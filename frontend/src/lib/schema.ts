import { z } from "zod";

export const projectNameSchema = z.object({
    name: z.string().min(1, "Name is required")
});

export const taskSchema = z.object({
    title: z.string().trim().min(1, "Title is required"),
    due_date: z.date({required_error: "Due Date is required"}),
    description: z.string().optional(),
    priority: z.enum(["p1", "p2", "p3", "p4"]),
    status: z.enum(["todo", "in_progress", "staging_review", "done"]),
    assignee_id: z.string().optional(),
});

export const registerSchema = z.object({
    name: z.string().trim().min(1, "Name is required"),
    email: z.string().min(1, "Email is required").email("Invalid email address"), 
    password: z.string().min(8, "Password must be at least 8 characters").max(50, "Password is too long"),
})

export const loginSchema = z.object({
    email: z.string().min(1, "Email is required").email("Invalid email address"), 
    password: z.string().min(8, "Password must be at least 8 characters").max(50, "Password is too long"),
})

export const editUserSchema = z.object({
    name: z.string().trim().min(1, "Name is required"),
    password: z.string().max(50, "Password is too long").refine(
        (val) => val.length === 0 || val.length >= 8,
        "Password must be empty or greater than 8 characters"
    ).optional().or(z.literal("")),
})

export const commentSchema = z.object({
    content: z.string().trim().min(1, "Content cannot be empty").max(5000, "Content is too long"),
})