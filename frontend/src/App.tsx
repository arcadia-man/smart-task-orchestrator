import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import SchedulerDetail from './pages/SchedulerDetail';
import CreateScheduler from './pages/CreateScheduler';
import Users from './pages/Users';

function App() {
    return (
        <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={
                <ProtectedRoute>
                    <Layout />
                </ProtectedRoute>
            }>
                <Route index element={<Dashboard />} />
                <Route path="schedulers/:id" element={<SchedulerDetail />} />
                <Route path="schedulers/create" element={<CreateScheduler />} />
                <Route path="users" element={<Users />} />
            </Route>
        </Routes>
    );
}

export default App;