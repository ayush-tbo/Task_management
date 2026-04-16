<<<<<<< HEAD
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

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/project/:id" element={<Project />} />
        <Route path="/task/:id" element={<Task />} />
        <Route path="/addEdit" element={<AddEditTask />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/activity" element={<ActivityLog />} />
        <Route path="*" element={<PageNotFound />} />
      </Routes>
    </Router>
  );
}

export default App;
=======
export default function App() {
  return <h1>Hello from React frontend!</h1>;
}
>>>>>>> 3913f28a646f762fd92ac93f16be93d0ad6d3ceb
