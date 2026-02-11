
import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import LoginRoute from './components/LoginRoute';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Schedulers from './pages/Schedulers';
import SchedulerDetail from './pages/SchedulerDetail';
import CreateScheduler from './pages/CreateScheduler';
import Users from './pages/Users';
import Roles from './pages/Roles';
import Images from './pages/Images';
import Logs from './pages/Logs';
import Monitoring from './pages/Monitoring';

function App() {
    return (
        <Routes>
            <Route path="/login" element={
                <LoginRoute>
                    <Login />
                </LoginRoute>
            } />
            <Route path="/" element={
                <ProtectedRoute>
                    <Layout />
                </ProtectedRoute>
            }>
                <Route index element={<Dashboard />} />
                <Route path="schedulers" element={<Schedulers />} />
                <Route path="schedulers/:id" element={<SchedulerDetail />} />
                <Route path="schedulers/create" element={<CreateScheduler />} />
                <Route path="users" element={<Users />} />
                <Route path="roles" element={<Roles />} />
                <Route path="images" element={<Images />} />
                <Route path="logs" element={<Logs />} />
                <Route path="monitoring" element={<Monitoring />} />
            </Route>
        </Routes>
    );
}

export default App;