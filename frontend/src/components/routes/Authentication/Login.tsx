import React from "react";
import { Button } from "@/components/ui/button"
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Header from "../Header/Header";
import Footer from "../Header/Footer";
import { useNavigate } from "react-router-dom";

function Login() {

    const navigate = useNavigate();

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
                        <form>
                            <div className="flex flex-col gap-6">
                                <div className="grid gap-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="m@example.com"
                                    required
                                />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="password">Password</Label>
                                    <Input id="password" type="password" required />
                                </div>
                            </div>
                        </form>
                    </CardContent>
                    <CardFooter className="flex-col gap-2">
                        <Button type="submit" className="w-full">Login</Button>
                    </CardFooter>
                </Card>
            </div>
            <Footer />
        </div>
    )
}

export default Login;