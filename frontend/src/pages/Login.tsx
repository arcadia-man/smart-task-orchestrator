import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Calendar, Lock, User } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { useToastContext } from '../contexts/ToastContext';
import ChangePasswordModal from '../components/ChangePasswordModal';

const Login: React.FC = () => {
    const [username, setUsername] = useState('admin');
    const [password, setPassword] = useState('12345');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [showChangePassword, setShowChangePassword] = useState(false);
    const [isInitialLogin, setIsInitialLogin] = useState(false);
    const navigate = useNavigate();
    const { login } = useAuth();
    const toast = useToastContext();



    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const result = await login(username, password);

            if (result.success) {
                // Check if password change is required
                if (result.user?.isInitialLogin) {
                    setIsInitialLogin(true);
                    setShowChangePassword(true);
                    toast.warning(
                        'Password Change Required',
                        'You must change your password before continuing.'
                    );
                } else {
                    toast.success('Login Successful', `Welcome back, ${result.user?.username}!`);
                    navigate('/');
                }
            } else {
                setError(result.error || 'Login failed');
                toast.error('Login Failed', result.error || 'Please check your credentials');
            }
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Connection failed';
            setError(errorMessage);
            toast.error('Connection Error', 'Unable to connect to server. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    const handlePasswordChangeSuccess = (newToken?: string) => {
        if (newToken) {
            // Update token in localStorage
            localStorage.setItem('token', newToken);
        }
        toast.success('Password Changed', 'Your password has been updated successfully!');
        setShowChangePassword(false);
        setIsInitialLogin(false);
        navigate('/');
    };

    const handlePasswordChangeError = (error: string) => {
        if (error.includes('contact') && error.includes('admin')) {
            toast.error('Contact Administrator', error);
        } else {
            toast.error('Password Change Failed', error);
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                <div>
                    <div className="flex justify-center">
                        <Calendar className="w-12 h-12 text-blue-600" />
                    </div>
                    <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                        Smart Task Orchestrator
                    </h2>
                    <p className="mt-2 text-center text-sm text-gray-600">
                        Sign in to your account
                    </p>
                </div>

                <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
                    <div className="space-y-4">
                        <div>
                            <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                                Username
                            </label>
                            <div className="mt-1 relative">
                                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <User className="h-5 w-5 text-gray-400" />
                                </div>
                                <input
                                    id="username"
                                    name="username"
                                    type="text"
                                    required
                                    className="input pl-10"
                                    placeholder="Enter your username"
                                    value={username}
                                    onChange={e => setUsername(e.target.value)}
                                />
                            </div>
                        </div>

                        <div>
                            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                                Password
                            </label>
                            <div className="mt-1 relative">
                                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <Lock className="h-5 w-5 text-gray-400" />
                                </div>
                                <input
                                    id="password"
                                    name="password"
                                    type="password"
                                    required
                                    className="input pl-10"
                                    placeholder="Enter your password"
                                    value={password}
                                    onChange={e => setPassword(e.target.value)}
                                />
                            </div>
                        </div>
                    </div>

                    {error && (
                        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                            {error}
                        </div>
                    )}

                    <div>
                        <button
                            type="submit"
                            disabled={loading}
                            className="w-full btn btn-primary"
                        >
                            {loading ? 'Signing in...' : 'Sign in'}
                        </button>
                    </div>

                    <div className="text-center">
                        <p className="text-sm text-gray-600">
                            Default credentials: <strong>admin / 12345</strong>
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                            You will be prompted to change the password on first login
                        </p>
                    </div>
                </form>

                <ChangePasswordModal
                    isOpen={showChangePassword}
                    onClose={() => !isInitialLogin && setShowChangePassword(false)}
                    onSuccess={handlePasswordChangeSuccess}
                    onError={handlePasswordChangeError}
                    isInitialLogin={isInitialLogin}
                />
            </div>
        </div>
    );
};

export default Login;