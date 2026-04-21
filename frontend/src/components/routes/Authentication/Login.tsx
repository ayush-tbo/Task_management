import React, { useState } from "react";
import { Button } from "@/components/ui/button"
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema } from "@/lib/schema";
import axios from "axios";

function Login() {

    const [loginError, setLoginError] = useState<string | null>(null);

    const navigate = useNavigate();
    const { login } = useAuth();

    const { register, handleSubmit, formState:{errors}} = useForm({
        resolver:zodResolver(loginSchema),
        defaultValues:{
            email: "",
            password: "",
        },
    });

    const onSubmit = async (data : any) => {
        setLoginError(null);
        try{
            const res = await axios.post("http://localhost:8080/api/users/login", data);
            login(res.data.token, res.data.user);
            navigate(`/dashboard`)
        }
        catch(err: any){
            if (err.response && err.response.status === 401) {
                setLoginError("Invalid email or password. Please try again.");
            } else {
                setLoginError("Something went wrong. Please try again later.");
            }
            console.error("Login Error:", err);
        }
    };

    return (
        <div>
            <Header />
            <div className="pt-31 pb-12 px-105">
                <Card className="w-full max-w-sm">
                    <CardHeader>
                        <CardTitle>Login to your account</CardTitle>
                        <CardDescription>Enter your email below to login to your account</CardDescription>
                        <CardAction>
                        <Button variant="link" onClick={() => navigate(`/register`)}>Sign Up</Button>
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
                                <Button type="submit" className="w-full">Login</Button>
                            </CardFooter>
                        </form>
                    </CardContent>
                </Card>
            </div>
            <Footer />
        </div>
    )
}

export default Login;