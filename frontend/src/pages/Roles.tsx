import React, { useState } from 'react';
import { Plus, Edit, Trash2, Shield, Users } from 'lucide-react';

const Roles: React.FC = () => {
    const [showCreateRole, setShowCreateRole] = useState(false);

    // Mock data - will be replaced with real API calls
    const roles = [
        {
            id: '1',
            name: 'Administrator',
            description: 'Full system access with all permissions',
            userCount: 2,
            permissions: ['read', 'write', 'delete', 'admin'],
            isSystem: true
        },
        {
            id: '2',
            name: 'Scheduler Manager',
            description: 'Can create and manage schedulers',
            userCount: 5,
            permissions: ['read', 'write'],
            isSystem: false
        },
        {
            id: '3',
            name: 'Viewer',
            description: 'Read-only access to view schedulers and logs',
            userCount: 10,
            permissions: ['read'],
            isSystem: false
        }
    ];

    const permissions = [
        { id: 'read', name: 'Read', description: 'View schedulers and system information' },
        { id: 'write', name: 'Write', description: 'Create and modify schedulers' },
        { id: 'delete', name: 'Delete', description: 'Delete schedulers and data' },
        { id: 'admin', name: 'Admin', description: 'Full administrative access' }
    ];

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Roles & Permissions</h1>
                    <p className="mt-2 text-gray-600">
                        Manage user roles and their permissions
                    </p>
                </div>
                <button
                    onClick={() => setShowCreateRole(true)}
                    className="btn btn-primary flex items-center"
                >
                    <Plus className="w-4 h-4 mr-2" />
                    Create Role
                </button>
            </div>

            {/* Roles List */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {roles.map((role) => (
                    <div key={role.id} className="card">
                        <div className="flex items-start justify-between mb-4">
                            <div className="flex items-center">
                                <div className="p-2 bg-blue-100 rounded-lg">
                                    <Shield className="w-5 h-5 text-blue-600" />
                                </div>
                                <div className="ml-3">
                                    <h3 className="text-lg font-semibold text-gray-900">
                                        {role.name}
                                        {role.isSystem && (
                                            <span className="ml-2 px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded">
                                                System
                                            </span>
                                        )}
                                    </h3>
                                    <p className="text-sm text-gray-600">{role.description}</p>
                                </div>
                            </div>
                            {!role.isSystem && (
                                <div className="flex items-center space-x-2">
                                    <button className="btn btn-secondary text-sm">
                                        <Edit className="w-4 h-4" />
                                    </button>
                                    <button className="btn btn-secondary text-sm text-red-600 hover:bg-red-50">
                                        <Trash2 className="w-4 h-4" />
                                    </button>
                                </div>
                            )}
                        </div>

                        <div className="mb-4">
                            <div className="flex items-center text-sm text-gray-600 mb-2">
                                <Users className="w-4 h-4 mr-1" />
                                {role.userCount} users assigned
                            </div>
                            <div className="flex flex-wrap gap-2">
                                {role.permissions.map((permission) => {
                                    const perm = permissions.find(p => p.id === permission);
                                    return (
                                        <span
                                            key={permission}
                                            className="px-2 py-1 text-xs bg-green-100 text-green-800 rounded-full"
                                        >
                                            {perm?.name || permission}
                                        </span>
                                    );
                                })}
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            {/* Permissions Reference */}
            <div className="card">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Permission Reference</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {permissions.map((permission) => (
                        <div key={permission.id} className="flex items-start p-3 bg-gray-50 rounded-lg">
                            <div className="p-1 bg-white rounded">
                                <Shield className="w-4 h-4 text-gray-600" />
                            </div>
                            <div className="ml-3">
                                <div className="font-medium text-gray-900">{permission.name}</div>
                                <div className="text-sm text-gray-600">{permission.description}</div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Create Role Modal would go here */}
            {showCreateRole && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
                        <div className="p-6">
                            <h3 className="text-lg font-semibold text-gray-900 mb-4">Create New Role</h3>
                            <p className="text-gray-600">Role creation form will be implemented here.</p>
                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    onClick={() => setShowCreateRole(false)}
                                    className="btn btn-secondary"
                                >
                                    Cancel
                                </button>
                                <button className="btn btn-primary">
                                    Create Role
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Roles;