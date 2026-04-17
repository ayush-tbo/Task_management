import React from "react";
import {BrowserRouter as Router, Routes, Route} from "react-router-dom";
import './App.css';
import Home from "./components/routes/Home/Home";
import Dashboard from "./components/routes/Dashboard/Dashboard";
import Project from "./components/routes/Project/Project";
import Task from "./components/routes/Task/Task";
import AddEditTask from "./components/routes/AddEditTask/AddEditTask";
import PageNotFound from "./components/routes/Authentication/PageNotFound";
import Login from "./components/routes/Authentication/Login";
import Register from "./components/routes/Authentication/Register";
import ActivityLog from "./components/routes/ActivityLog/ActivityLog";
import { AuthProvider } from "./context/AuthContext";
import { ProtectedRoute } from "./components/routes/Authentication/ProtectedRoutes";
import axios from "axios";

function App() {

  axios.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
  
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="*" element={<PageNotFound />} />

          <Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
          <Route path="/project/:id" element={<ProtectedRoute><Project /></ProtectedRoute>} />
          <Route path="/task/:id" element={<ProtectedRoute><Task /></ProtectedRoute>} />
          <Route path="/addEdit" element={<ProtectedRoute><AddEditTask /></ProtectedRoute>} />
          <Route path="/activity" element={<ProtectedRoute><ActivityLog /></ProtectedRoute>} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
