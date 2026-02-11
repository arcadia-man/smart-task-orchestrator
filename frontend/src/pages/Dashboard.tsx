import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Plus, Play, Pause, Clock, CheckCircle, XCircle, Calendar } from 'lucide-react';
import { dashboardAPI, schedulersAPI } from '../services/api';
import { useToastContext } from '../contexts/ToastContext';

const Dashboard: React.FC = () => {
    const [schedulers, setSchedulers] = useState<any[]>([]);
    const [stats, setStats] = useState<any>({});
    const [loading, setLoading] = useState(true);
    const toast = useToastContext();

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            setLoading(true);
            const [schedulersResponse, statsResponse] = await Promise.all([
                schedulersAPI.getAll(),
                dashboardAPI.getStats()
            ]);
            
            setSchedulers(schedulersResponse.data.slice(0, 5)); // Show only first 5 for dashboard
            setStats(statsResponse.data);
        } catch (error: any) {
            toast.error('Failed to load dashboard data', error.response?.data?.error || 'Please try again');
        } finally {
            setLoading(false);
        }
    };

    const handleRunScheduler = async (schedulerId: string) => {
        try {
            await schedulersAPI.run(schedulerId);
            toast.success('Scheduler started', 'The scheduler has been queued for execution');
            fetchData(); // Refresh data
        } catch (error: any) {
            toast.error('Failed to run scheduler', error.response?.data?.error || 'Please try again');
        }
    };

    const handleToggleScheduler = async (schedulerId: string, currentStatus: string) => {
        try {
            const newStatus = currentStatus === 'active' ? 'paused' : 'active';
            await schedulersAPI.update(schedulerId, { status: newStatus });
            toast.success('Scheduler updated', `Scheduler has been ${newStatus}`);
            fetchData(); // Refresh data
        } catch (error: any) {
            toast.error('Failed to update scheduler', error.response?.data?.error || 'Please try again');
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'success':
                return <CheckCircle className="w-5 h-5 text-green-500" />;
            case 'failed':
                return <XCircle className="w-5 h-5 text-red-500" />;
            case 'running':
                return <Clock className="w-5 h-5 text-blue-500 animate-spin" />;
            default:
                return <Clock className="w-5 h-5 text-gray-400" />;
        }
    };

    const getStatusBadge = (status: string) => {
        const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
        switch (status) {
            case 'active':
                return `${baseClasses} bg-green-100 text-green-800`;
            case 'paused':
                return `${baseClasses} bg-yellow-100 text-yellow-800`;
            case 'inactive':
                return `${baseClasses} bg-gray-100 text-gray-800`;
            default:
                return `${baseClasses} bg-gray-100 text-gray-800`;
        }
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                    <p className="mt-2 text-gray-600">
                        Manage and monitor your scheduled tasks
                    </p>
                </div>
                <Link
                    to="/schedulers/create"
                    className="btn btn-primary flex items-center"
                >
                    <Plus className="w-4 h-4 mr-2" />
                    Create Scheduler
                </Link>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-blue-100 rounded-lg">
                            <Calendar className="w-6 h-6 text-blue-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Total Schedulers</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : stats.total || 0}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-green-100 rounded-lg">
                            <Play className="w-6 h-6 text-green-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Active</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : stats.active || 0}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-yellow-100 rounded-lg">
                            <Pause className="w-6 h-6 text-yellow-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Paused</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : stats.paused || 0}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-red-100 rounded-lg">
                            <XCircle className="w-6 h-6 text-red-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Inactive</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : stats.inactive || 0}
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Schedulers List */}
            <div className="card">
                <div className="flex items-center justify-between mb-6">
                    <h2 className="text-xl font-semibold text-gray-900">Schedulers</h2>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="border-b border-gray-200">
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Name</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Status</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Type</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Schedule</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Last Run</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Next Run</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {loading ? (
                                <tr>
                                    <td colSpan={7} className="py-8 text-center text-gray-500">
                                        Loading schedulers...
                                    </td>
                                </tr>
                            ) : schedulers.length === 0 ? (
                                <tr>
                                    <td colSpan={7} className="py-8 text-center text-gray-500">
                                        No schedulers found. <Link to="/schedulers/create" className="text-blue-600 hover:text-blue-800">Create your first scheduler</Link>
                                    </td>
                                </tr>
                            ) : (
                                schedulers.map((scheduler) => (
                                    <tr key={scheduler.id} className="border-b border-gray-100 hover:bg-gray-50">
                                        <td className="py-4 px-4">
                                            <div>
                                                <Link
                                                    to={`/schedulers/${scheduler.id}`}
                                                    className="font-medium text-blue-600 hover:text-blue-800"
                                                >
                                                    {scheduler.name}
                                                </Link>
                                                <p className="text-sm text-gray-500">{scheduler.description}</p>
                                            </div>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className={getStatusBadge(scheduler.status)}>
                                                {scheduler.status}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className="text-sm text-gray-600">{scheduler.jobType}</span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className="text-sm text-gray-600">
                                                {scheduler.jobType === 'cron' 
                                                    ? scheduler.cronExpr 
                                                    : scheduler.jobType === 'interval'
                                                    ? `${scheduler.intervalSeconds}s`
                                                    : 'Immediate'}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="flex items-center">
                                                {getStatusIcon(scheduler.lastStatus)}
                                                <span className="ml-2 text-sm text-gray-600">
                                                    {scheduler.lastRun 
                                                        ? new Date(scheduler.lastRun).toLocaleString()
                                                        : 'Never'}
                                                </span>
                                            </div>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className="text-sm text-gray-600">
                                                {scheduler.nextRun 
                                                    ? new Date(scheduler.nextRun).toLocaleString()
                                                    : 'N/A'}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="flex items-center space-x-2">
                                                <button 
                                                    className="btn btn-secondary text-sm"
                                                    onClick={() => handleRunScheduler(scheduler.id)}
                                                    title="Run Now"
                                                >
                                                    <Play className="w-4 h-4" />
                                                </button>
                                                <button 
                                                    className="btn btn-secondary text-sm"
                                                    onClick={() => handleToggleScheduler(scheduler.id, scheduler.status)}
                                                    title={scheduler.status === 'active' ? 'Pause' : 'Resume'}
                                                >
                                                    {scheduler.status === 'active' ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;