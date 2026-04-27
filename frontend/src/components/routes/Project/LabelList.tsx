import React, { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Plus, Trash2, Pencil } from "lucide-react";
import axios from "axios";

const defaultColors = ["#ef4444", "#f97316", "#eab308", "#22c55e", "#3b82f6", "#8b5cf6", "#ec4899", "#6b7280"];

function LabelList({ projectId }: { projectId: string }) {
    const [labels, setLabels] = useState<any[]>([]);
    const [open, setOpen] = useState(false);
    const [editingLabel, setEditingLabel] = useState<any>(null);
    const [name, setName] = useState("");
    const [color, setColor] = useState("#3b82f6");

    const fetchLabels = async () => {
        try {
            const res = await axios.get(`/api/projects/${projectId}/labels`);
            setLabels(res.data.labels || []);
        } catch (err) {
            console.error("Failed to fetch labels:", err);
        }
    };

    useEffect(() => {
        fetchLabels();
    }, [projectId]);

    const resetForm = () => {
        setName("");
        setColor("#3b82f6");
        setEditingLabel(null);
    };

    const handleSave = async () => {
        if (!name.trim()) return;
        try {
            if (editingLabel) {
                await axios.put(`/api/labels/${editingLabel.id}`, { name, color });
            } else {
                await axios.post(`/api/projects/${projectId}/labels`, { name, color });
            }
            resetForm();
            setOpen(false);
            fetchLabels();
        } catch (err) {
            console.error("Failed to save label:", err);
        }
    };

    const handleEdit = (label: any) => {
        setEditingLabel(label);
        setName(label.name);
        setColor(label.color || "#3b82f6");
        setOpen(true);
    };

    const handleDelete = async (id: string) => {
        if (!window.confirm("Delete this label?")) return;
        try {
            await axios.delete(`/api/labels/${id}`);
            fetchLabels();
        } catch (err) {
            console.error("Failed to delete label:", err);
        }
    };

    return (
        <div>
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-lg font-semibold">Labels</h2>
                <Dialog open={open} onOpenChange={(val) => { setOpen(val); if (!val) resetForm(); }}>
                    <DialogTrigger asChild>
                        <Button variant="outline" size="sm">
                            <Plus size={16} className="mr-1" />New Label
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-sm" aria-describedby={undefined}>
                        <DialogHeader>
                            <DialogTitle>{editingLabel ? "Edit Label" : "Create Label"}</DialogTitle>
                        </DialogHeader>
                        <div className="space-y-3">
                            <div>
                                <label className="text-sm font-medium">Name</label>
                                <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Bug, Feature, etc." />
                            </div>
                            <div>
                                <label className="text-sm font-medium">Color</label>
                                <div className="flex gap-2 mt-1.5 flex-wrap">
                                    {defaultColors.map((c) => (
                                        <button
                                            key={c}
                                            type="button"
                                            onClick={() => setColor(c)}
                                            className={`w-7 h-7 rounded-full border-2 transition-all ${
                                                color === c ? "border-slate-900 scale-110" : "border-transparent"
                                            }`}
                                            style={{ backgroundColor: c }}
                                        />
                                    ))}
                                </div>
                            </div>
                        </div>
                        <DialogFooter>
                            <DialogClose asChild>
                                <Button variant="outline">Cancel</Button>
                            </DialogClose>
                            <Button onClick={handleSave}>{editingLabel ? "Update" : "Create"}</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            {labels.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-8">No labels yet. Create one to categorize tasks.</p>
            ) : (
                <div className="flex flex-wrap gap-3">
                    {labels.map((label) => (
                        <Card key={label.id} className="w-fit">
                            <CardContent className="p-3 flex items-center gap-3">
                                <div
                                    className="w-4 h-4 rounded-full shrink-0"
                                    style={{ backgroundColor: label.color || "#6b7280" }}
                                />
                                <span className="font-medium text-sm">{label.name}</span>
                                <div className="flex gap-0.5 ml-2">
                                    <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => handleEdit(label)}>
                                        <Pencil size={13} />
                                    </Button>
                                    <Button variant="ghost" size="icon" className="h-7 w-7 text-red-500 hover:text-red-700" onClick={() => handleDelete(label.id)}>
                                        <Trash2 size={13} />
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}

export default LabelList;
