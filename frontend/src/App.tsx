import React from "react";
import {BrowserRouter as Router, Routes, Route} from "react-router-dom";
import './App.css';
import Home from "./components/routes/Home/Home";
import Dashboard from "./components/routes/Dashboard/Dashboard";
import Project from "./components/routes/Project/Project";
import Task from "./components/routes/Task/Task";
import AddEditTask from "./components/routes/AddEditTask/AddEditTask";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/project/:id" element={<Project />} />
        <Route path="/task/:id" element={<Task />} />
        <Route path="/addEdit" element={<AddEditTask />} />
      </Routes>
    </Router>
  );
}

export default App;