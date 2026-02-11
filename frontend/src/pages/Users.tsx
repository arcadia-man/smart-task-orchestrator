import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, User, Mail, Shield, Key } from 'lucide-react';
import { usersAPI, rolesAPI } from '../services/api';
import { useToastContext } from '../contexts/ToastContext';

const Users: React.FC = () => {
    const [users, setUsers] = useState<any[]>([]);
    const [roles, setRoles] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreateUser, setShowCreateUser] = useState(false);
    const [createUserData, setCreateUserData] = useState({
        username: '',
        email: '',
        roleId: '',
        password: ''
    });
    const toast = useToastContext();

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            setLoading(true);
            const [usersResponse, rolesResponse] = await Promise.all([
                usersAPI.getAll(),
                rolesAPI.getAll()
            ]);
            
            setUsers(usersResponse.data);
            setRoles(rolesResponse.data);
        } catch (error: any) {
            toast.error('Failed to load users', error.response?.data?.error || 'Please try again');
        } finally {
            setLoading(false);
        }
    };

    const handleCreateUser = async () => {
        try {
            await usersAPI.create(createUserData);
            toast.success('User created', 'User has been created successfully');
            setShowCreateUser(false);
            setCreateUserData({ username: '', email: '', roleId: '', password: '' });
            fetchData();
        } catch (error: any) {
            toast.error('Failed to create user', error.response?.data?.error || 'Please try again');
        }
    };

    const handleDeleteUser = async (userId: string, username: string) => {
        if (!window.confirm(`Are you sure you want to delete user "${username}"?`)) {
            return;
        }

        try {
            await usersAPI.delete(userId);
            toast.success('User deleted', 'User has been deleted successfully');
            fetchData();
        } catch (error: any) {
            toast.error('Failed to delete user', error.response?.data?.error || 'Please try again');
        }
    };

    const handleResetPassword = async (userId: string, username: string) => {
        const newPassword = window.prompt(`Enter new password for user "${username}":`);
        if (!newPassword || newPassword.length < 8) {
            toast.error('Invalid password', 'Password must be at least 8 characters long');
            return;
        }

        try {
            await usersAPI.resetPassword(userId, newPassword);
            toast.success('Password reset', 'User password has been reset successfully');
            fetchData();
        } catch (error: any) {
            toast.error('Failed to reset password', error.response?.data?.error || 'Please try again');
        }
    };

    const getStatusBadge = (status: string) => {
        const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
        switch (status) {
            case 'active':
                return `${baseClasses} bg-green-100 text-green-800`;
            case 'inactive':
                return `${baseClasses} bg-gray-100 text-gray-800`;
            default:
                return `${baseClasses} bg-gray-100 text-gray-800`;
        }
    };

    const getRoleBadge = (role: string) => {
        const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
        switch (role) {
            case 'Administrator':
                return `${baseClasses} bg-red-100 text-red-800`;
            case 'Scheduler Manager':
                return `${baseClasses} bg-blue-100 text-blue-800`;
            case 'Viewer':
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
                    <h1 className="text-3xl font-bold text-gray-900">User Management</h1>
                    <p className="mt-2 text-gray-600">
                        Manage users, roles, and permissions
                    </p>
                </div>
                <button
                    onClick={() => setShowCreateUser(true)}
                    className="btn btn-primary flex items-center"
                >
                    <Plus className="w-4 h-4 mr-2" />
                    Create User
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-blue-100 rounded-lg">
                            <User className="w-6 h-6 text-blue-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Total Users</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : users.length}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-green-100 rounded-lg">
                            <Shield className="w-6 h-6 text-green-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Active Users</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : users.filter(u => u.active).length}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-yellow-100 rounded-lg">
                            <Key className="w-6 h-6 text-yellow-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Pending Password Change</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : users.filter(u => u.isInitialLogin).length}
                            </p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-red-100 rounded-lg">
                            <Mail className="w-6 h-6 text-red-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Administrators</p>
                            <p className="text-2xl font-bold text-gray-900">
                                {loading ? '...' : users.filter(u => u.roleName === 'Administrator').length}
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Users Table */}
            <div className="card">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="border-b border-gray-200">
                                <th className="text-left py-3 px-4 font-medium text-gray-600">User</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Role</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Status</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Last Login</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Created</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {loading ? (
                                <tr>
                                    <td colSpan={6} className="py-8 text-center text-gray-500">
                                        Loading users...
                                    </td>
                                </tr>
                            ) : users.length === 0 ? (
                                <tr>
                                    <td colSpan={6} className="py-8 text-center text-gray-500">
                                        No users found. Create your first user to get started.
                                    </td>
                                </tr>
                            ) : (
                                users.map((user) => (
                                    <tr key={user.id} className="border-b border-gray-100 hover:bg-gray-50">
                                        <td className="py-4 px-4">
                                            <div className="flex items-center">
                                                <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
                                                    <User className="w-4 h-4 text-gray-600" />
                                                </div>
                                                <div className="ml-3">
                                                    <div className="font-medium text-gray-900 flex items-center">
                                                        {user.username}
                                                        {user.isInitialLogin && (
                                                            <span className="ml-2 px-2 py-1 text-xs bg-yellow-100 text-yellow-800 rounded">
                                                                Password Change Required
                                                            </span>
                                                        )}
                                                    </div>
                                                    <div className="text-sm text-gray-500">{user.email}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className={getRoleBadge(user.roleName || 'Unknown')}>
                                                {user.roleName || 'Unknown'}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className={getStatusBadge(user.active ? 'active' : 'inactive')}>
                                                {user.active ? 'active' : 'inactive'}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className="text-sm text-gray-600">
                                                {user.lastLoginAt 
                                                    ? new Date(user.lastLoginAt).toLocaleDateString()
                                                    : 'Never'}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className="text-sm text-gray-600">
                                                {new Date(user.createdAt).toLocaleDateString()}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="flex items-center space-x-2">
                                                <button 
                                                    className="btn btn-secondary text-sm"
                                                    title="Edit User"
                                                    onClick={() => toast.info('Edit User', 'Edit functionality coming soon')}
                                                >
                                                    <Edit className="w-4 h-4" />
                                                </button>
                                                <button 
                                                    className="btn btn-secondary text-sm"
                                                    title="Reset Password"
                                                    onClick={() => handleResetPassword(user.id, user.username)}
                                                >
                                                    <Key className="w-4 h-4" />
                                                </button>
                                                <button 
                                                    className="btn btn-secondary text-sm text-red-600 hover:bg-red-50"
                                                    title="Delete User"
                                                    onClick={() => handleDeleteUser(user.id, user.username)}
                                                >
                                                    <Trash2 className="w-4 h-4" />
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

            {/* Create User Modal */}
            {showCreateUser && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
                        <div className="p-6">
                            <h3 className="text-lg font-semibold text-gray-900 mb-4">Create New User</h3>
                            
                            <div className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Username
                                    </label>
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="Enter username"
                                        value={createUserData.username}
                                        onChange={(e) => setCreateUserData({...createUserData, username: e.target.value})}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Email
                                    </label>
                                    <input
                                        type="email"
                                        className="input"
                                        placeholder="Enter email address"
                                        value={createUserData.email}
                                        onChange={(e) => setCreateUserData({...createUserData, email: e.target.value})}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Role
                                    </label>
                                    <select 
                                        className="input"
                                        value={createUserData.roleId}
                                        onChange={(e) => setCreateUserData({...createUserData, roleId: e.target.value})}
                                    >
                                        <option value="">Select a role</option>
                                        {roles.map((role) => (
                                            <option key={role.id} value={role.id}>
                                                {role.roleName}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Initial Password
                                    </label>
                                    <input
                                        type="password"
                                        className="input"
                                        placeholder="Enter initial password"
                                        value={createUserData.password}
                                        onChange={(e) => setCreateUserData({...createUserData, password: e.target.value})}
                                    />
                                    <p className="text-xs text-gray-500 mt-1">
                                        User will be required to change this password on first login
                                    </p>
                                </div>
                            </div>
                            
                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    onClick={() => setShowCreateUser(false)}
                                    className="btn btn-secondary"
                                >
                                    Cancel
                                </button>
                                <button 
                                    className="btn btn-primary"
                                    onClick={handleCreateUser}
                                    disabled={!createUserData.username || !createUserData.email || !createUserData.roleId || !createUserData.password}
                                >
                                    Create User
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Users;