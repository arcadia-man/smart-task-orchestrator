import React from 'react';
import { Link } from 'react-router-dom';
import { Plus, Play, Pause, Clock, CheckCircle, XCircle, Edit, Trash2 } from 'lucide-react';

const Schedulers: React.FC = () => {
    // Mock data - will be replaced with real API calls
    const schedulers = [
        {
            id: '1',
            name: 'Daily ETL Pipeline',
            description: 'Process daily data from multiple sources',
            status: 'active',
            jobType: 'cron',
            cronExpr: '0 2 * * *',
            lastRun: '2025-10-22T02:00:00Z',
            nextRun: '2025-10-23T02:00:00Z',
            lastStatus: 'success',
        },
        {
            id: '2',
            name: 'Health Check Monitor',
            description: 'Monitor system health every 5 minutes',
            status: 'active',
            jobType: 'interval',
            intervalSeconds: 300,
            lastRun: '2025-10-22T10:25:00Z',
            nextRun: '2025-10-22T10:30:00Z',
            lastStatus: 'running',
        },
        {
            id: '3',
            name: 'Weekly Report Generator',
            description: 'Generate weekly analytics report',
            status: 'paused',
            jobType: 'cron',
            cronExpr: '0 9 * * 1',
            lastRun: '2025-10-15T09:00:00Z',
            nextRun: null,
            lastStatus: 'failed',
        },
    ];

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'success':
                return <CheckCircle className="w-4 h-4 text-green-500" />;
            case 'failed':
                return <XCircle className="w-4 h-4 text-red-500" />;
            case 'running':
                return <Clock className="w-4 h-4 text-blue-500 animate-spin" />;
            default:
                return <Clock className="w-4 h-4 text-gray-400" />;
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
                    <h1 className="text-3xl font-bold text-gray-900">Schedulers</h1>
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

            {/* Schedulers List */}
            <div className="card">
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
                            {schedulers.map((scheduler) => (
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
                                                : `${scheduler.intervalSeconds}s`}
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
                                                title="Run Now"
                                            >
                                                <Play className="w-4 h-4" />
                                            </button>
                                            <button 
                                                className="btn btn-secondary text-sm"
                                                title="Pause"
                                            >
                                                <Pause className="w-4 h-4" />
                                            </button>
                                            <Link
                                                to={`/schedulers/${scheduler.id}/edit`}
                                                className="btn btn-secondary text-sm"
                                                title="Edit"
                                            >
                                                <Edit className="w-4 h-4" />
                                            </Link>
                                            <button 
                                                className="btn btn-secondary text-sm text-red-600 hover:bg-red-50"
                                                title="Delete"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default Schedulers;