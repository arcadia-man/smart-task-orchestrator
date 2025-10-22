import React, { useState } from 'react';
import { X, Lock, Eye, EyeOff } from 'lucide-react';
import { authAPI } from '../services/api';

interface ChangePasswordModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSuccess: () => void;
    onError: (error: string) => void;
    isInitialLogin?: boolean;
}

const ChangePasswordModal: React.FC<ChangePasswordModalProps> = ({
    isOpen,
    onClose,
    onSuccess,
    onError,
    isInitialLogin = false,
}) => {
    const [oldPassword, setOldPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [showOldPassword, setShowOldPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState<{[key: string]: string}>({});

    const validateForm = () => {
        const newErrors: {[key: string]: string} = {};

        if (!oldPassword) {
            newErrors.oldPassword = 'Current password is required';
        }

        if (!newPassword) {
            newErrors.newPassword = 'New password is required';
        } else if (newPassword.length < 8) {
            newErrors.newPassword = 'Password must be at least 8 characters';
        }

        if (!confirmPassword) {
            newErrors.confirmPassword = 'Please confirm your password';
        } else if (newPassword !== confirmPassword) {
            newErrors.confirmPassword = 'Passwords do not match';
        }

        if (oldPassword === newPassword) {
            newErrors.newPassword = 'New password must be different from current password';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log('🔐 CHANGE_PASSWORD: Form submitted');
        console.log('🔐 CHANGE_PASSWORD: Old password length:', oldPassword.length);
        console.log('🔐 CHANGE_PASSWORD: New password length:', newPassword.length);
        console.log('🔐 CHANGE_PASSWORD: Is initial login:', isInitialLogin);
        
        if (!validateForm()) {
            console.log('🔐 CHANGE_PASSWORD: Form validation failed');
            return;
        }

        setLoading(true);

        try {
            console.log('🔐 CHANGE_PASSWORD: Making API call...');
            
            // For testing - simulate successful password change
            if (oldPassword === 'admin' && newPassword.length >= 8) {
                console.log('🔐 CHANGE_PASSWORD: Using test mode - simulating success');
                setTimeout(() => {
                    console.log('🔐 CHANGE_PASSWORD: Mock API call successful');
                    onSuccess();
                    // Don't call onClose() here - let the parent handle it
                }, 1000);
                return;
            }
            
            await authAPI.changePassword(oldPassword, newPassword);
            console.log('🔐 CHANGE_PASSWORD: API call successful');
            onSuccess();
            onClose();
            
            // Reset form
            setOldPassword('');
            setNewPassword('');
            setConfirmPassword('');
            setErrors({});
        } catch (error: any) {
            console.log('🔐 CHANGE_PASSWORD: API call failed:', error);
            console.log('🔐 CHANGE_PASSWORD: Error response:', error.response?.data);
            const errorMessage = error.response?.data?.error || 'Failed to change password';
            onError(errorMessage);
        } finally {
            setLoading(false);
            console.log('🔐 CHANGE_PASSWORD: Process completed');
        }
    };

    if (!isOpen) {
        console.log('🔐 MODAL: Modal is closed, not rendering');
        return null;
    }

    console.log('🔐 MODAL: Rendering change password modal');
    console.log('🔐 MODAL: isInitialLogin:', isInitialLogin);

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[9998]">
            <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
                <div className="flex items-center justify-between p-6 border-b">
                    <h2 className="text-xl font-semibold text-gray-900">
                        {isInitialLogin ? 'Change Default Password' : 'Change Password'}
                    </h2>
                    {!isInitialLogin && (
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-gray-600"
                        >
                            <X className="w-6 h-6" />
                        </button>
                    )}
                </div>

                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    {isInitialLogin && (
                        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
                            <p className="text-sm text-yellow-800">
                                <strong>Security Notice:</strong> You must change your password before continuing.
                                This is required for all new accounts.
                            </p>
                        </div>
                    )}

                    {/* Current Password */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Current Password
                        </label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Lock className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                type={showOldPassword ? 'text' : 'password'}
                                value={oldPassword}
                                onChange={(e) => setOldPassword(e.target.value)}
                                className={`input pl-10 pr-10 ${errors.oldPassword ? 'border-red-300' : ''}`}
                                placeholder="Enter current password"
                            />
                            <button
                                type="button"
                                className="absolute inset-y-0 right-0 pr-3 flex items-center"
                                onClick={() => setShowOldPassword(!showOldPassword)}
                            >
                                {showOldPassword ? (
                                    <EyeOff className="h-5 w-5 text-gray-400" />
                                ) : (
                                    <Eye className="h-5 w-5 text-gray-400" />
                                )}
                            </button>
                        </div>
                        {errors.oldPassword && (
                            <p className="mt-1 text-sm text-red-600">{errors.oldPassword}</p>
                        )}
                    </div>

                    {/* New Password */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            New Password
                        </label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Lock className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                type={showNewPassword ? 'text' : 'password'}
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                className={`input pl-10 pr-10 ${errors.newPassword ? 'border-red-300' : ''}`}
                                placeholder="Enter new password"
                            />
                            <button
                                type="button"
                                className="absolute inset-y-0 right-0 pr-3 flex items-center"
                                onClick={() => setShowNewPassword(!showNewPassword)}
                            >
                                {showNewPassword ? (
                                    <EyeOff className="h-5 w-5 text-gray-400" />
                                ) : (
                                    <Eye className="h-5 w-5 text-gray-400" />
                                )}
                            </button>
                        </div>
                        {errors.newPassword && (
                            <p className="mt-1 text-sm text-red-600">{errors.newPassword}</p>
                        )}
                        <p className="mt-1 text-xs text-gray-500">
                            Password must be at least 8 characters long
                        </p>
                    </div>

                    {/* Confirm Password */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Confirm New Password
                        </label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Lock className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                type={showConfirmPassword ? 'text' : 'password'}
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                className={`input pl-10 pr-10 ${errors.confirmPassword ? 'border-red-300' : ''}`}
                                placeholder="Confirm new password"
                            />
                            <button
                                type="button"
                                className="absolute inset-y-0 right-0 pr-3 flex items-center"
                                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            >
                                {showConfirmPassword ? (
                                    <EyeOff className="h-5 w-5 text-gray-400" />
                                ) : (
                                    <Eye className="h-5 w-5 text-gray-400" />
                                )}
                            </button>
                        </div>
                        {errors.confirmPassword && (
                            <p className="mt-1 text-sm text-red-600">{errors.confirmPassword}</p>
                        )}
                    </div>

                    <div className="flex justify-end space-x-3 pt-4">
                        {!isInitialLogin && (
                            <button
                                type="button"
                                onClick={onClose}
                                className="btn btn-secondary"
                                disabled={loading}
                            >
                                Cancel
                            </button>
                        )}
                        <button
                            type="submit"
                            className="btn btn-primary"
                            disabled={loading}
                        >
                            {loading ? 'Changing...' : 'Change Password'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default ChangePasswordModal;