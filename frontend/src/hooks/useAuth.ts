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
        try {
            const response = await authAPI.login(username, password);
            const { access_token, user: userData, must_change_password } = response.data;
            
            localStorage.setItem('token', access_token);
            setToken(access_token);
            
            // Create user object with initial login flag
            const userWithInitialLogin = {
                ...userData,
                isInitialLogin: must_change_password || userData?.isInitialLogin || false
            };
            
            setUser(userWithInitialLogin);
            
            return { success: true, user: userWithInitialLogin };
        } catch (error: any) {
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
        window.location.href = '/login';
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