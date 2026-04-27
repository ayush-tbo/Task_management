import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PenBox, Loader2 } from "lucide-react";
import axios from "axios";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { editUserSchema } from "@/lib/schema";

function Profile() {

    const { isAuthenticated, user, login, logout } = useAuth();
    const navigate = useNavigate();

    const [open, setOpen] = useState(false);
    const [updateError, setUpdateError] = useState<string | null>(null);

    const { register, handleSubmit, formState:{errors}, reset} = useForm({
        resolver:zodResolver(editUserSchema),
    });

    useEffect(() => {
        if (user) {
            reset({
                name: user.name,
                password: "",
            });
        }
    }, [user, reset, open]);

    const handleLogout = () => {
        logout();
        navigate("/login");
    };

    const handleUpdateProfile = async (data: any) => {
        setUpdateError(null);
        try{
            const res = await axios.patch(`/api/users/${user.id}`, data);
            login(localStorage.getItem("token") || "", res.data.user);
            setOpen(false);
        }
        catch(err){
            setUpdateError("Failed to update profile. Please try again later.");
            console.error("Update failed:", err);
        }
    };

    const getInitials = (name?: string) => {
        if (!name) return "";
        return name
        .split(" ")
        .map((n) => n[0])
        .join("")
        .toUpperCase()
        .slice(0, 2);
    };

    return (
        <div>
            {isAuthenticated && user ? (
                <div>
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="relative h-10 w-10 rounded-full">
                            <Avatar className="h-10 w-10">
                                <AvatarImage src={user?.avatarURL} alt={user?.name} />
                                <AvatarFallback>{getInitials(user?.name)}</AvatarFallback>
                            </Avatar>
                        </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56" align="end" forceMount>
                        <DropdownMenuLabel className="font-normal">
                            <div className="flex flex-col space-y-1">
                                <p className="text-sm font-medium leading-none">{user?.name}</p>
                                <p className="text-xs leading-none text-muted-foreground">{user?.email}</p>
                            </div>
                        </DropdownMenuLabel>
                        <DropdownMenuItem onClick={() => setOpen(true)} className="cursor-pointer">
                            <PenBox size={16} className="mr-2" /> Edit Profile
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem onClick={handleLogout} className="text-red-600 focus:text-red-600">
                            Log out
                        </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    <Dialog open={open} onOpenChange={setOpen}>
                        <DialogContent aria-describedby={undefined}>
                            <DialogHeader>
                                <DialogTitle>Edit Profile</DialogTitle>
                            </DialogHeader>
                            <form onSubmit={handleSubmit(handleUpdateProfile)}>
                                <div className="grid gap-4 py-4">
                                    {updateError && (
                                        <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-2 rounded-md text-sm">
                                            {updateError}
                                        </div>
                                    )}

                                    <div className="grid gap-2">
                                        <Label htmlFor="name">Enter Your Name</Label>
                                        <Input id="name" {...register("name")}/>
                                        {errors.name && (
                                            <p className="text-sm text-red-500">{errors.name.message as string}</p>
                                        )}
                                    </div>
                                    <div className="grid gap-2">
                                        <Label htmlFor="password">New Password (leave blank to keep current)</Label>
                                        <Input id="password" {...register("password")} />
                                        {errors.password && (
                                            <p className="text-sm text-red-500">{errors.password.message as string}</p>
                                        )}
                                    </div>
                                </div>
                                <DialogFooter>
                                    <Button type="button" variant="outline" onClick={() => setOpen(false)}>Cancel</Button>
                                    <Button type="submit">Save Changes</Button>
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>
                </div>
            ) : (
                <Link to="/login">
                    <Button variant="outline">Log in</Button>
                </Link>
            )}
        </div>
    );
}

export default Profile;