import { z } from "zod";

export const projectNameSchema = z.object({
    name: z.string().min(1, "Name is required")
});