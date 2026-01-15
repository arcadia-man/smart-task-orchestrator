import React, { useState } from 'react';
import { Search, Filter, Download, RefreshCw, AlertCircle, CheckCircle, XCircle, Clock } from 'lucide-react';

const Logs: React.FC = () => {
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedLevel, setSelectedLevel] = useState('all');
    const [autoRefresh, setAutoRefresh] = useState(false);

    // Mock data - will be replaced with real API calls
    const logs = [
        {
            id: '1',
            timestamp: '2025-10-22T10:30:15Z',
            level: 'info',
            source: 'scheduler-1',
            message: 'Daily ETL Pipeline started successfully',
            details: 'Job ID: job_123456, Duration: 0s'
        },
        {
            id: '2',
            timestamp: '2025-10-22T10:29:45Z',
            level: 'error',
            source: 'scheduler-3',
            message: 'Weekly Report Generator failed to connect to database',
            details: 'Error: connection timeout after 30s, Host: db.example.com:5432'
        },
        {
            id: '3',
            timestamp: '2025-10-22T10:25:30Z',
            level: 'success',
            source: 'scheduler-2',
            message: 'Health Check Monitor completed successfully',
            details: 'All systems operational, Response time: 150ms'
        },
        {
            id: '4',
            timestamp: '2025-10-22T10:20:12Z',
            level: 'warning',
            source: 'system',
            message: 'High memory usage detected',
            details: 'Memory usage: 85%, Available: 2.1GB'
        },
        {
            id: '5',
            timestamp: '2025-10-22T10:15:00Z',
            level: 'info',
            source: 'auth',
            message: 'User admin logged in successfully',
            details: 'IP: 192.168.1.100, User-Agent: Mozilla/5.0...'
        }
    ];

    const getLevelIcon = (level: string) => {
        switch (level) {
            case 'success':
                return <CheckCircle className="w-4 h-4 text-green-500" />;
            case 'error':
                return <XCircle className="w-4 h-4 text-red-500" />;
            case 'warning':
                return <AlertCircle className="w-4 h-4 text-yellow-500" />;
            case 'info':
                return <Clock className="w-4 h-4 text-blue-500" />;
            default:
                return <Clock className="w-4 h-4 text-gray-400" />;
        }
    };

    const getLevelBadge = (level: string) => {
        const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
        switch (level) {
            case 'success':
                return `${baseClasses} bg-green-100 text-green-800`;
            case 'error':
                return `${baseClasses} bg-red-100 text-red-800`;
            case 'warning':
                return `${baseClasses} bg-yellow-100 text-yellow-800`;
            case 'info':
                return `${baseClasses} bg-blue-100 text-blue-800`;
            default:
                return `${baseClasses} bg-gray-100 text-gray-800`;
        }
    };

    const filteredLogs = logs.filter(log => {
        const matchesSearch = log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
                            log.source.toLowerCase().includes(searchTerm.toLowerCase());
        const matchesLevel = selectedLevel === 'all' || log.level === selectedLevel;
        return matchesSearch && matchesLevel;
    });

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">System Logs</h1>
                    <p className="mt-2 text-gray-600">
                        Monitor system events and scheduler activities
                    </p>
                </div>
                <div className="flex items-center space-x-3">
                    <button
                        onClick={() => setAutoRefresh(!autoRefresh)}
                        className={`btn ${autoRefresh ? 'btn-primary' : 'btn-secondary'} flex items-center`}
                    >
                        <RefreshCw className={`w-4 h-4 mr-2 ${autoRefresh ? 'animate-spin' : ''}`} />
                        Auto Refresh
                    </button>
                    <button className="btn btn-secondary flex items-center">
                        <Download className="w-4 h-4 mr-2" />
                        Export
                    </button>
                </div>
            </div>

            {/* Filters */}
            <div className="card">
                <div className="flex flex-col sm:flex-row gap-4">
                    <div className="flex-1">
                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                            <input
                                type="text"
                                placeholder="Search logs..."
                                className="input pl-10"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>
                    </div>
                    <div className="sm:w-48">
                        <select
                            className="input"
                            value={selectedLevel}
                            onChange={(e) => setSelectedLevel(e.target.value)}
                        >
                            <option value="all">All Levels</option>
                            <option value="success">Success</option>
                            <option value="info">Info</option>
                            <option value="warning">Warning</option>
                            <option value="error">Error</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-green-100 rounded-lg">
                            <CheckCircle className="w-6 h-6 text-green-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Success</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {logs.filter(log => log.level === 'success').length}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-blue-100 rounded-lg">
                            <Clock className="w-6 h-6 text-blue-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Info</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {logs.filter(log => log.level === 'info').length}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-yellow-100 rounded-lg">
                            <AlertCircle className="w-6 h-6 text-yellow-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Warning</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {logs.filter(log => log.level === 'warning').length}
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
                            <p className="text-sm font-medium text-gray-600">Error</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {logs.filter(log => log.level === 'error').length}
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Logs List */}
            <div className="card">
                <div className="space-y-4">
                    {filteredLogs.map((log) => (
                        <div key={log.id} className="border-b border-gray-100 pb-4 last:border-b-0">
                            <div className="flex items-start justify-between">
                                <div className="flex items-start space-x-3 flex-1">
                                    {getLevelIcon(log.level)}
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center space-x-2 mb-1">
                                            <span className={getLevelBadge(log.level)}>
                                                {log.level}
                                            </span>
                                            <span className="text-sm text-gray-500">
                                                {log.source}
                                            </span>
                                            <span className="text-sm text-gray-400">
                                                {new Date(log.timestamp).toLocaleString()}
                                            </span>
                                        </div>
                                        <p className="text-gray-900 font-medium">{log.message}</p>
                                        {log.details && (
                                            <p className="text-sm text-gray-600 mt-1">{log.details}</p>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {filteredLogs.length === 0 && (
                    <div className="text-center py-8">
                        <p className="text-gray-500">No logs found matching your criteria.</p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Logs;