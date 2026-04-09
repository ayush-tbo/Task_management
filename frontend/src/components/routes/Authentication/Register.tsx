import React from "react";
import { Button } from "@/components/ui/button"
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { useNavigate } from "react-router-dom";

function Register() {

    const navigate = useNavigate();

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
                        <form>
                            <div className="flex flex-col gap-6">
                                <div className="grid gap-2">
                                <Label htmlFor="email">Enter Your Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="m@example.com"
                                    required
                                />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="password">Create Your Password</Label>
                                    <Input id="password" type="password" required />
                                </div>
                            </div>
                        </form>
                    </CardContent>
                    <CardFooter className="flex-col gap-2">
                        <Button type="submit" className="w-full">Register</Button>
                    </CardFooter>
                </Card>
            </div>
            <Footer />
        </div>
    );
}

export default Register;