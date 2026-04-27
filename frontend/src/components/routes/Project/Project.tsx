import React, { useEffect, useState } from "react";
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import TaskGrid from "./TaskGrid";
import SprintList from "./SprintList";
import LabelList from "./LabelList";
import MemberList from "./MemberList";
import ProjectCharts from "./ProjectCharts";
import { Link, useNavigate, useParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Logs, Plus, Users, Tag, Zap, BarChart3, Settings, Trash2, Pencil } from "lucide-react";
import axios from "axios";

type Tab = "board" | "sprints" | "labels" | "members" | "charts";

function Project() {
    const { id: projectId } = useParams<{ id: string }>();
    const [project, setProject] = useState<any>(null);
    const [tab, setTab] = useState<Tab>("board");
    const [editOpen, setEditOpen] = useState(false);
    const [editName, setEditName] = useState("");
    const [editDesc, setEditDesc] = useState("");
    const navigate = useNavigate();

    const fetchProject = async () => {
        try {
            const res = await axios.get(`/api/projects/${projectId}`);
            setProject(res.data.project);
        } catch (err) {
            console.error("Failed to fetch project:", err);
        }
    };

    useEffect(() => {
        if (projectId) fetchProject();
    }, [projectId]);

    const handleEditOpen = () => {
        setEditName(project?.name || "");
        setEditDesc(project?.description || "");
        setEditOpen(true);
    };

    const handleEditSave = async () => {
        if (!editName.trim()) return;
        try {
            const res = await axios.put(`/api/projects/${projectId}`, {
                name: editName,
                description: editDesc,
            });
            setProject(res.data.project);
            setEditOpen(false);
        } catch (err) {
            console.error("Failed to update project:", err);
        }
    };

    const handleDelete = async () => {
        if (!window.confirm("Delete this project and all its data? This cannot be undone.")) return;
        try {
            await axios.delete(`/api/projects/${projectId}`);
            navigate("/dashboard");
        } catch (err) {
            console.error("Failed to delete project:", err);
        }
    };

    const tabs: { key: Tab; label: string; icon: React.ReactNode }[] = [
        { key: "board", label: "Board", icon: <Logs size={16} /> },
        { key: "sprints", label: "Sprints", icon: <Zap size={16} /> },
        { key: "labels", label: "Labels", icon: <Tag size={16} /> },
        { key: "members", label: "Members", icon: <Users size={16} /> },
        { key: "charts", label: "Charts", icon: <BarChart3 size={16} /> },
    ];

    return (
        <div className="min-h-screen flex flex-col">
            <Header />
            <div className="px-4 pt-20 pb-5 flex-1">
                <div className="flex flex-col sm:flex-row justify-between items-start gap-3">
                    <div>
                        <h1 className="text-3xl sm:text-5xl font-bold gradient-title mb-1">{project?.name ?? ""}</h1>
                        {project?.description && (
                            <p className="text-sm text-muted-foreground mb-2">{project.description}</p>
                        )}
                    </div>
                    <div className="flex flex-wrap gap-2">
                        {tab === "board" && (
                            <Link to={`/addEdit?projectId=${projectId}`}>
                                <Button variant="outline" size="sm">
                                    <Plus size={16} className="mr-1" />Add Task
                                </Button>
                            </Link>
                        )}
                        <Link to={`/activity?projectId=${projectId}`}>
                            <Button size="sm" variant="outline">
                                <Logs size={16} className="mr-1" />Activity
                            </Button>
                        </Link>
                        <Button size="sm" variant="outline" onClick={handleEditOpen}>
                            <Pencil size={14} className="mr-1" />Edit
                        </Button>
                        <Button size="sm" variant="outline" className="text-red-600 hover:text-red-700" onClick={handleDelete}>
                            <Trash2 size={14} className="mr-1" />Delete
                        </Button>
                    </div>
                </div>

                <div className="flex gap-1 border-b mb-4 overflow-x-auto">
                    {tabs.map((t) => (
                        <button
                            key={t.key}
                            onClick={() => setTab(t.key)}
                            className={`flex items-center gap-1.5 px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                                tab === t.key
                                    ? "border-blue-600 text-blue-600"
                                    : "border-transparent text-muted-foreground hover:text-foreground"
                            }`}
                        >
                            {t.icon}{t.label}
                        </button>
                    ))}
                </div>

                {tab === "board" && <TaskGrid projectId={projectId!} />}
                {tab === "sprints" && <SprintList projectId={projectId!} />}
                {tab === "labels" && <LabelList projectId={projectId!} />}
                {tab === "members" && <MemberList projectId={projectId!} />}
                {tab === "charts" && <ProjectCharts projectId={projectId!} />}
            </div>

            {/* Edit Project Dialog */}
            <Dialog open={editOpen} onOpenChange={setEditOpen}>
                <DialogContent className="sm:max-w-md" aria-describedby={undefined}>
                    <DialogHeader>
                        <DialogTitle>Edit Project</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-3">
                        <div>
                            <label className="text-sm font-medium">Name</label>
                            <Input value={editName} onChange={(e) => setEditName(e.target.value)} />
                        </div>
                        <div>
                            <label className="text-sm font-medium">Description</label>
                            <Textarea value={editDesc} onChange={(e) => setEditDesc(e.target.value)} placeholder="Project description..." />
                        </div>
                    </div>
                    <DialogFooter>
                        <DialogClose asChild>
                            <Button variant="outline">Cancel</Button>
                        </DialogClose>
                        <Button onClick={handleEditSave}>Save</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Footer />
        </div>
    );
}

export default Project;