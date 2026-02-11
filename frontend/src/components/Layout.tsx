import React, { useState } from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { 
    Calendar, 
    Users, 
    Settings, 
    LogOut, 
    Key, 
    Plus, 
    Shield, 
    Database, 
    Activity,
    FileText,
    Clock,
    BarChart3
} from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { useToastContext } from '../contexts/ToastContext';
import ChangePasswordModal from './ChangePasswordModal';

const Layout: React.FC = () => {
    const { logout, user } = useAuth();
    const toast = useToastContext();
    const location = useLocation();
    const [showChangePassword, setShowChangePassword] = useState(false);

    const handleLogout = () => {
        logout();
        toast.success('Logged Out', 'You have been successfully logged out.');
    };

    const handlePasswordChangeSuccess = (newToken?: string) => {
        if (newToken) {
            // Update token in localStorage and reload user info
            localStorage.setItem('token', newToken);
            window.location.reload(); // Reload to refresh user context
        } else {
            toast.success('Password Changed', 'Your password has been updated successfully!');
            setShowChangePassword(false);
        }
    };

    const handlePasswordChangeError = (error: string) => {
        toast.error('Password Change Failed', error);
    };

    const isActiveRoute = (path: string) => {
        return location.pathname === path;
    };

    const navItems = [
        {
            path: '/',
            label: 'Dashboard',
            icon: BarChart3,
            description: 'Overview and statistics'
        },
        {
            path: '/schedulers',
            label: 'Schedulers',
            icon: Clock,
            description: 'Manage scheduled tasks'
        },
        {
            path: '/schedulers/create',
            label: 'Create Scheduler',
            icon: Plus,
            description: 'Create new scheduled task'
        },
        {
            path: '/users',
            label: 'User Management',
            icon: Users,
            description: 'Manage users and roles'
        },
        {
            path: '/roles',
            label: 'Roles & Permissions',
            icon: Shield,
            description: 'Manage roles and permissions'
        },
        {
            path: '/images',
            label: 'Docker Images',
            icon: Database,
            description: 'Manage container images'
        },
        {
            path: '/logs',
            label: 'System Logs',
            icon: FileText,
            description: 'View system logs'
        },
        {
            path: '/monitoring',
            label: 'Monitoring',
            icon: Activity,
            description: 'System health monitoring'
        }
    ];

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Sidebar */}
            <div className="fixed inset-y-0 left-0 w-64 bg-white shadow-lg">
                <div className="flex flex-col h-full">
                    {/* Logo */}
                    <div className="flex items-center px-6 py-4 border-b">
                        <Calendar className="w-8 h-8 text-blue-600" />
                        <span className="ml-2 text-xl font-bold text-gray-900">
                            Task Orchestrator
                        </span>
                    </div>

                    {/* Navigation */}
                    <nav className="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
                        {navItems.map((item) => {
                            const Icon = item.icon;
                            const isActive = isActiveRoute(item.path);
                            
                            return (
                                <Link
                                    key={item.path}
                                    to={item.path}
                                    className={`flex items-center px-4 py-3 rounded-lg transition-colors group ${
                                        isActive
                                            ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-600'
                                            : 'text-gray-700 hover:bg-gray-100'
                                    }`}
                                    title={item.description}
                                >
                                    <Icon className={`w-5 h-5 mr-3 ${isActive ? 'text-blue-600' : 'text-gray-500'}`} />
                                    <div className="flex-1">
                                        <div className={`font-medium ${isActive ? 'text-blue-700' : 'text-gray-900'}`}>
                                            {item.label}
                                        </div>
                                        <div className="text-xs text-gray-500 mt-0.5">
                                            {item.description}
                                        </div>
                                    </div>
                                </Link>
                            );
                        })}
                    </nav>

                    {/* User menu */}
                    <div className="px-4 py-4 border-t bg-gray-50">
                        {/* Password Change Warning Badge */}
                        {user?.isInitialLogin && (
                            <div className="px-3 py-2 mb-3 bg-red-50 border border-red-200 rounded-lg">
                                <div className="flex items-center">
                                    <Key className="w-4 h-4 text-red-600 mr-2" />
                                    <p className="text-xs text-red-800 font-medium">
                                        Change password to secure your account
                                    </p>
                                </div>
                            </div>
                        )}
                        
                        <div className="px-4 py-3 mb-3 bg-white rounded-lg border">
                            <div className="flex items-center">
                                <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                                    <span className="text-blue-600 font-semibold text-sm">
                                        {user?.username?.charAt(0).toUpperCase() || 'U'}
                                    </span>
                                </div>
                                <div className="ml-3">
                                    <div className="text-sm font-medium text-gray-900">
                                        {user?.username || 'User'}
                                    </div>
                                    <div className="text-xs text-gray-500">
                                        {user?.email || 'No email'}
                                    </div>
                                </div>
                            </div>
                        </div>
                        
                        <div className="space-y-1">
                            <button
                                onClick={() => setShowChangePassword(true)}
                                className="flex items-center w-full px-4 py-2 text-gray-700 rounded-lg hover:bg-white hover:shadow-sm transition-all"
                            >
                                <Key className="w-4 h-4 mr-3 text-gray-500" />
                                <span className="text-sm">Change Password</span>
                            </button>
                            <button
                                onClick={handleLogout}
                                className="flex items-center w-full px-4 py-2 text-red-700 rounded-lg hover:bg-red-50 transition-all"
                            >
                                <LogOut className="w-4 h-4 mr-3 text-red-500" />
                                <span className="text-sm">Logout</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Main content */}
            <div className="ml-64">
                <main className="p-8">
                    <Outlet />
                </main>
            </div>

            <ChangePasswordModal
                isOpen={showChangePassword}
                onClose={() => setShowChangePassword(false)}
                onSuccess={handlePasswordChangeSuccess}
                onError={handlePasswordChangeError}
                isInitialLogin={false}
            />
        </div>
    );
};

export default Layout;