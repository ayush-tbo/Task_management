import React, { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Plus, Trash2, UserCircle } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import axios from "axios";

function MemberList({ projectId }: { projectId: string }) {
    const { user } = useAuth();
    const [members, setMembers] = useState<any[]>([]);
    const [allUsers, setAllUsers] = useState<any[]>([]);
    const [open, setOpen] = useState(false);
    const [selectedUserId, setSelectedUserId] = useState("");

    const fetchMembers = async () => {
        try {
            const res = await axios.get(`/api/projects/${projectId}/members`);
            setMembers(res.data.members || []);
        } catch (err) {
            console.error("Failed to fetch members:", err);
        }
    };

    const fetchAllUsers = async () => {
        try {
            const res = await axios.get("/api/users");
            setAllUsers(res.data.users || []);
        } catch (err) {
            console.error("Failed to fetch users:", err);
        }
    };

    useEffect(() => {
        fetchMembers();
    }, [projectId]);

    const handleOpenDialog = (val: boolean) => {
        setOpen(val);
        if (val) {
            fetchAllUsers();
            setSelectedUserId("");
        }
    };

    const availableUsers = allUsers.filter(
        (u) => !members.some((m) => m.user_id === u.id)
    );

    const handleAdd = async () => {
        if (!selectedUserId) return;
        try {
            await axios.post(`/api/projects/${projectId}/members`, { user_id: selectedUserId, role: "member" });
            setSelectedUserId("");
            setOpen(false);
            fetchMembers();
        } catch (err) {
            console.error("Failed to add member:", err);
        }
    };

    const handleRemove = async (memberId: string) => {
        if (!window.confirm("Remove this member?")) return;
        try {
            await axios.delete(`/api/projects/${projectId}/members/${memberId}`);
            fetchMembers();
        } catch (err) {
            console.error("Failed to remove member:", err);
        }
    };

    const roleColor: Record<string, string> = {
        owner: "bg-amber-100 text-amber-800",
        admin: "bg-blue-100 text-blue-700",
        member: "bg-slate-100 text-slate-600",
    };

    return (
        <div>
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-lg font-semibold">Members</h2>
                <Dialog open={open} onOpenChange={handleOpenDialog}>
                    <DialogTrigger asChild>
                        <Button variant="outline" size="sm">
                            <Plus size={16} className="mr-1" />Add Member
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-sm" aria-describedby={undefined}>
                        <DialogHeader>
                            <DialogTitle>Add Member</DialogTitle>
                        </DialogHeader>
                        <div>
                            <label className="text-sm font-medium">Select User</label>
                            {availableUsers.length === 0 ? (
                                <p className="text-sm text-muted-foreground mt-2">No users available to add.</p>
                            ) : (
                                <Select value={selectedUserId} onValueChange={setSelectedUserId}>
                                    <SelectTrigger className="mt-1">
                                        <SelectValue placeholder="Choose a user..." />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {availableUsers.map((u) => (
                                            <SelectItem key={u.id} value={u.id}>
                                                {u.name} ({u.email})
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            )}
                        </div>
                        <DialogFooter>
                            <DialogClose asChild>
                                <Button variant="outline">Cancel</Button>
                            </DialogClose>
                            <Button onClick={handleAdd} disabled={!selectedUserId}>Add</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            {members.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-8">No members found.</p>
            ) : (
                <div className="space-y-2">
                    {members.map((m) => (
                        <Card key={m.user_id}>
                            <CardContent className="p-3 flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                    {m.avatar_url ? (
                                        <img src={m.avatar_url} alt={m.name} className="w-8 h-8 rounded-full" />
                                    ) : (
                                        <UserCircle size={32} className="text-slate-300" />
                                    )}
                                    <div>
                                        <span className="text-sm font-medium">{m.name}</span>
                                        <p className="text-xs text-muted-foreground">{m.email}</p>
                                    </div>
                                    <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded ${roleColor[m.role] || roleColor.member}`}>
                                        {m.role?.toUpperCase()}
                                    </span>
                                </div>
                                {m.role !== "owner" && m.user_id !== user?.id && (
                                    <Button variant="ghost" size="icon" className="h-8 w-8 text-red-500 hover:text-red-700" onClick={() => handleRemove(m.user_id)}>
                                        <Trash2 size={14} />
                                    </Button>
                                )}
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}

export default MemberList;
