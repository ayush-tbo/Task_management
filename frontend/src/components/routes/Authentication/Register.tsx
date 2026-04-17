import React, { useState } from "react";
import { Button } from "@/components/ui/button"
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { useNavigate } from "react-router-dom";
import { registerSchema } from "@/lib/schema";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod/src/zod.js";
import axios from "axios";
import { useAuth } from "@/context/AuthContext";

function Register() {

    const [loginError, setLoginError] = useState<string | null>(null);

    const navigate = useNavigate();
    const { login } = useAuth();

    const { register, handleSubmit, formState:{errors}} = useForm({
        resolver:zodResolver(registerSchema),
        defaultValues:{
            name: "",
            email: "",
            password: "",
        },
    });

    const onSubmit = async (data : any) => {
        try{
            const res = await axios.post("http://localhost:8080/api/users/register", data);
            login(res.data.token, res.data.user);
            navigate(`/dashboard`)
        }
        catch(err: any){
            if (err.response && err.response.status === 409){
                setLoginError("Email already exist. Please try another email.");
            }
            else{
                setLoginError("Something went wrong. Please try again later.");
            }
            console.error("Register Error:", err);
        }
    };

    return (
        <div>
            <Header />
            <div className="pt-31 pb-12 px-105">
                <Card className="w-full max-w-sm">
                    <CardHeader>
                        <CardTitle>Create New Account</CardTitle>
                        <CardDescription>Enter your email below to create your account</CardDescription>
                        <CardAction>
                            <Button variant="link" onClick={() => navigate(`/login`)}>Sign In</Button>
                        </CardAction>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <div className="flex flex-col gap-6">
                                {loginError && (
                                    <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-2 rounded-md text-sm">
                                        {loginError}
                                    </div>
                                )}

                                <div className="grid gap-2">
                                    <Label htmlFor="name">Enter Your Name</Label>
                                    <Input id="name" placeholder="Meet Kotadiya" {...register("name")} />
                                    {errors.name && (
                                        <p className="text-sm text-red-500">{errors.name.message}</p>
                                    )}
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="email">Enter Your Email</Label>
                                    <Input id="email" placeholder="m@example.com" {...register("email")} />
                                    {errors.email && (
                                        <p className="text-sm text-red-500">{errors.email.message}</p>
                                    )}
                                </div>
                                <div className="grid gap-2 pb-2">
                                    <Label htmlFor="password">Create Your Password</Label>
                                    <Input id="password" {...register("password")} />
                                    {errors.password && (
                                        <p className="text-sm text-red-500">{errors.password.message}</p>
                                    )}
                                </div>
                            </div>
                            <CardFooter className="flex-col gap-2">
                                <Button type="submit" className="w-full">Register</Button>
                            </CardFooter>
                        </form>
                    </CardContent>
                </Card>
            </div>
            <Footer />
        </div>
    );
}

export default Register;