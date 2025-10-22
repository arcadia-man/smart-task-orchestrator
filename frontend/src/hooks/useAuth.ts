import { useState, useEffect } from 'react';
import { authAPI } from '../services/api';

interface User {
    id: string;
    username: string;
    email: string;
    roleId: string;
    isInitialLogin: boolean;
}

export const useAuth = () => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [token, setToken] = useState<string | null>(localStorage.getItem('token'));

    useEffect(() => {
        if (token) {
            // Verify token and get user info
            authAPI.me()
                .then((response) => {
                    setUser(response.data);
                })
                .catch(() => {
                    // Token is invalid
                    localStorage.removeItem('token');
                    setToken(null);
                })
                .finally(() => {
                    setLoading(false);
                });
        } else {
            setLoading(false);
        }
    }, [token]);

    const login = async (username: string, password: string) => {
        console.log('🔐 AUTH: Starting login process');
        console.log('🔐 AUTH: Username:', username);
        
        try {
            console.log('🔐 AUTH: Making API call to /api/auth/login');
            const response = await authAPI.login(username, password);
            console.log('🔐 AUTH: API response:', response.data);
            
            // Handle different response formats
            if (response.data.message) {
                console.log('🔐 AUTH: Received placeholder response from backend');
                // Placeholder response from backend
                return { 
                    success: false, 
                    error: 'Backend API not fully implemented yet. Please implement the login handler.' 
                };
            }
            
            const { access_token, user: userData, must_change_password } = response.data;
            console.log('🔐 AUTH: Extracted data - token:', !!access_token, 'user:', userData, 'must_change_password:', must_change_password);
            
            localStorage.setItem('token', access_token);
            setToken(access_token);
            
            // Create user object with initial login flag
            const userWithInitialLogin = {
                ...userData,
                isInitialLogin: must_change_password || userData?.isInitialLogin || false
            };
            
            console.log('🔐 AUTH: Final user object:', userWithInitialLogin);
            setUser(userWithInitialLogin);
            
            return { success: true, user: userWithInitialLogin };
        } catch (error: any) {
            console.log('🔐 AUTH: Login error:', error);
            console.log('🔐 AUTH: Error response:', error.response?.data);
            return { 
                success: false, 
                error: error.response?.data?.error || 'Login failed' 
            };
        }
    };

    const logout = () => {
        localStorage.removeItem('token');
        setToken(null);
        setUser(null);
    };

    const isAuthenticated = !!token && !!user;

    return {
        user,
        loading,
        isAuthenticated,
        login,
        logout,
    };
};