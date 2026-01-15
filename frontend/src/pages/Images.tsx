import React, { useState } from 'react';
import { Plus, Download, Trash2, Database, Tag, Calendar } from 'lucide-react';

const Images: React.FC = () => {
    const [showAddImage, setShowAddImage] = useState(false);

    // Mock data - will be replaced with real API calls
    const images = [
        {
            id: '1',
            name: 'node',
            tag: '18-alpine',
            fullName: 'node:18-alpine',
            size: '45.2 MB',
            created: '2025-10-20T10:30:00Z',
            lastUsed: '2025-10-22T08:15:00Z',
            usageCount: 12
        },
        {
            id: '2',
            name: 'python',
            tag: '3.11-slim',
            fullName: 'python:3.11-slim',
            size: '123.8 MB',
            created: '2025-10-18T14:20:00Z',
            lastUsed: '2025-10-22T06:30:00Z',
            usageCount: 8
        },
        {
            id: '3',
            name: 'postgres',
            tag: '15',
            fullName: 'postgres:15',
            size: '156.4 MB',
            created: '2025-10-15T09:45:00Z',
            lastUsed: '2025-10-21T22:10:00Z',
            usageCount: 25
        },
        {
            id: '4',
            name: 'redis',
            tag: '7-alpine',
            fullName: 'redis:7-alpine',
            size: '28.9 MB',
            created: '2025-10-12T16:00:00Z',
            lastUsed: '2025-10-22T09:45:00Z',
            usageCount: 15
        }
    ];

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Docker Images</h1>
                    <p className="mt-2 text-gray-600">
                        Manage container images available for schedulers
                    </p>
                </div>
                <button
                    onClick={() => setShowAddImage(true)}
                    className="btn btn-primary flex items-center"
                >
                    <Plus className="w-4 h-4 mr-2" />
                    Add Image
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-blue-100 rounded-lg">
                            <Database className="w-6 h-6 text-blue-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Total Images</p>
                            <p className="text-2xl font-bold text-gray-900">{images.length}</p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-green-100 rounded-lg">
                            <Download className="w-6 h-6 text-green-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Total Size</p>
                            <p className="text-2xl font-bold text-gray-900">354 MB</p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-purple-100 rounded-lg">
                            <Tag className="w-6 h-6 text-purple-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Most Used</p>
                            <p className="text-2xl font-bold text-gray-900">postgres:15</p>
                        </div>
                    </div>
                </div>
                <div className="card">
                    <div className="flex items-center">
                        <div className="p-2 bg-orange-100 rounded-lg">
                            <Calendar className="w-6 h-6 text-orange-600" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-600">Last Added</p>
                            <p className="text-2xl font-bold text-gray-900">2 days ago</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Images List */}
            <div className="card">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="border-b border-gray-200">
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Image</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Tag</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Size</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Created</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Last Used</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Usage Count</th>
                                <th className="text-left py-3 px-4 font-medium text-gray-600">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {images.map((image) => (
                                <tr key={image.id} className="border-b border-gray-100 hover:bg-gray-50">
                                    <td className="py-4 px-4">
                                        <div className="flex items-center">
                                            <div className="p-2 bg-gray-100 rounded">
                                                <Database className="w-4 h-4 text-gray-600" />
                                            </div>
                                            <div className="ml-3">
                                                <div className="font-medium text-gray-900">{image.name}</div>
                                                <div className="text-sm text-gray-500">{image.fullName}</div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="py-4 px-4">
                                        <span className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded">
                                            {image.tag}
                                        </span>
                                    </td>
                                    <td className="py-4 px-4">
                                        <span className="text-sm text-gray-600">{image.size}</span>
                                    </td>
                                    <td className="py-4 px-4">
                                        <span className="text-sm text-gray-600">
                                            {new Date(image.created).toLocaleDateString()}
                                        </span>
                                    </td>
                                    <td className="py-4 px-4">
                                        <span className="text-sm text-gray-600">
                                            {new Date(image.lastUsed).toLocaleDateString()}
                                        </span>
                                    </td>
                                    <td className="py-4 px-4">
                                        <span className="text-sm text-gray-600">{image.usageCount} times</span>
                                    </td>
                                    <td className="py-4 px-4">
                                        <div className="flex items-center space-x-2">
                                            <button 
                                                className="btn btn-secondary text-sm"
                                                title="Pull Latest"
                                            >
                                                <Download className="w-4 h-4" />
                                            </button>
                                            <button 
                                                className="btn btn-secondary text-sm text-red-600 hover:bg-red-50"
                                                title="Remove Image"
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

            {/* Add Image Modal */}
            {showAddImage && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
                        <div className="p-6">
                            <h3 className="text-lg font-semibold text-gray-900 mb-4">Add Docker Image</h3>
                            <p className="text-gray-600 mb-4">Add a new Docker image to be available for schedulers.</p>
                            
                            <div className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Image Name
                                    </label>
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="e.g., node:18-alpine"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Description (Optional)
                                    </label>
                                    <textarea
                                        className="input"
                                        rows={3}
                                        placeholder="Brief description of the image..."
                                    />
                                </div>
                            </div>
                            
                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    onClick={() => setShowAddImage(false)}
                                    className="btn btn-secondary"
                                >
                                    Cancel
                                </button>
                                <button className="btn btn-primary">
                                    Add Image
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Images;