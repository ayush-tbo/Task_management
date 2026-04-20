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
import Notifications from "./components/routes/Notifications/Notifications";

function App() {

  axios.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });

  axios.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        window.location.href = "/login";
      }
      return Promise.reject(error);
    }
  )
  
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
          <Route path="/notifications" element={<ProtectedRoute><Notifications /></ProtectedRoute>} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
