import React, { useState } from 'react';
import { Outlet, Link } from 'react-router-dom';
import { Calendar, Users, Settings, LogOut, Key } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { useToastContext } from '../contexts/ToastContext';
import ChangePasswordModal from './ChangePasswordModal';

const Layout: React.FC = () => {
    const { logout, user } = useAuth();
    const toast = useToastContext();
    const [showChangePassword, setShowChangePassword] = useState(false);

    const handleLogout = () => {
        logout();
    };

    const handlePasswordChangeSuccess = () => {
        toast.success('Password Changed', 'Your password has been updated successfully!');
    };

    const handlePasswordChangeError = (error: string) => {
        toast.error('Password Change Failed', error);
    };

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
                    <nav className="flex-1 px-4 py-6 space-y-2">
                        <Link
                            to="/"
                            className="flex items-center px-4 py-2 text-gray-700 rounded-lg hover:bg-gray-100"
                        >
                            <Calendar className="w-5 h-5 mr-3" />
                            Dashboard
                        </Link>
                        <Link
                            to="/schedulers/create"
                            className="flex items-center px-4 py-2 text-gray-700 rounded-lg hover:bg-gray-100"
                        >
                            <Settings className="w-5 h-5 mr-3" />
                            Create Scheduler
                        </Link>
                        <Link
                            to="/users"
                            className="flex items-center px-4 py-2 text-gray-700 rounded-lg hover:bg-gray-100"
                        >
                            <Users className="w-5 h-5 mr-3" />
                            Users
                        </Link>
                    </nav>

                    {/* User menu */}
                    <div className="px-4 py-4 border-t">
                        <div className="px-4 py-2 text-sm text-gray-600">
                            Logged in as: <span className="font-medium">{user?.username || 'User'}</span>
                        </div>
                        <button
                            onClick={() => setShowChangePassword(true)}
                            className="flex items-center w-full px-4 py-2 text-gray-700 rounded-lg hover:bg-gray-100 mb-2"
                        >
                            <Key className="w-5 h-5 mr-3" />
                            Change Password
                        </button>
                        <button
                            onClick={handleLogout}
                            className="flex items-center w-full px-4 py-2 text-gray-700 rounded-lg hover:bg-gray-100"
                        >
                            <LogOut className="w-5 h-5 mr-3" />
                            Logout
                        </button>
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