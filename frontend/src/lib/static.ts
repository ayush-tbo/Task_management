export const testingProjects = [
    {
        id: 1,
        name: "Personal Work",
        todoTasks: 3,
        inProgressTasks: 1,
        reviewTasks: 0,
        completedTasks: 5
    },
    {
        id: 2,
        name: "Job",
        todoTasks: 2,
        inProgressTasks: 4,
        reviewTasks: 2,
        completedTasks: 2
    },
    {
        id: 3,
        name: "TMS",
        todoTasks: 10,
        inProgressTasks: 2,
        reviewTasks: 2,
        completedTasks: 0
    },
    {
        id: 4,
        name: "Finance Platform",
        todoTasks: 7,
        inProgressTasks: 0,
        reviewTasks: 0,
        completedTasks: 25
    },
]

export const testingTasks = [
    {
        id: 1,
        title: "Initialize monorepo repository",
        dueDate: new Date("2026-03-20T10:00:00Z"),
        description: "Set up the base Go backend, React TypeScript frontend, and Docker compose configurations.",
        priority: 1,
        status: "completed",
        assignedTo: "Meet K."
    },
    {
        id: 2,
        title: "Database schema migration",
        dueDate: new Date("2026-03-22T14:30:00Z"),
        description: "Write initial SQL migrations for the user and billing tables.",
        priority: 2,
        status: "inProgress",
        assignedTo: "Meet K."
    },
    {
        id: 3,
        title: "Implement User Authentication",
        dueDate: new Date("2026-03-27T17:00:00Z"),
        description: "Add JWT based auth in the Go backend and set up protected routes in React.",
        priority: 3,
        status: "review",
        assignedTo: "Adarsh K."
    },
    {
        id: 4,
        title: "Create Wallet Integration UI",
        dueDate: new Date("2026-03-29T12:00:00Z"),
        description: "Build reusable shadcn/ui components for the wallet connection modal.",
        priority: 4,
        status: "review",
        assignedTo: "Adarsh K."
    },
    {
        id: 5,
        title: "Develop Transaction History API",
        dueDate: new Date("2026-04-03T09:00:00Z"),
        description: "Create the GET endpoints in Go to fetch paginated transaction records.",
        priority: 5,
        status: "inProgress",
        assignedTo: "Ayush R."
    },
    {
        id: 6,
        title: "Optimize Docker Builds",
        dueDate: new Date("2026-04-05T15:00:00Z"),
        description: "Implement multi-stage builds to reduce the image size for the React frontend container.",
        priority: 6,
        status: "inProgress",
        assignedTo: "Ayush R."
    },
    {
        id: 7,
        title: "Setup CI/CD Pipeline",
        dueDate: new Date("2026-04-08T11:00:00Z"),
        description: "Configure GitHub Actions for automated testing and linting on pull requests.",
        priority: 7,
        status: "inProgress",
        assignedTo: "Meet K."
    },
    {
        id: 8,
        title: "Finalize API Documentation",
        dueDate: new Date("2026-04-12T16:00:00Z"),
        description: "Generate Swagger/OpenAPI documentation for all billing endpoints.",
        priority: 8,
        status: "todo",
        assignedTo: "Unassigned"
    },
    {
        id: 9,
        title: "Set Up Webhook Notifications",
        dueDate: new Date("2026-04-15T10:00:00Z"),
        description: "Implement event listeners in Go for failed billing transactions and send alerts.",
        priority: 9,
        status: "todo",
        assignedTo: "Meet K."
    },
    {
        id: 10,
        title: "Post-Launch Performance Review",
        dueDate: new Date("2026-04-25T14:00:00Z"),
        description: "Analyze system performance, check container health, and gather user feedback.",
        priority: 10,
        status: "todo",
        assignedTo: "Unassigned"
    }
]