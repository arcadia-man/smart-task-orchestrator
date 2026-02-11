import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Create axios instance
const api = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
});

// Request interceptor to add auth token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Response interceptor to handle auth errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// Auth API
export const authAPI = {
    login: (username: string, password: string) =>
        api.post('/api/auth/login', { username, password }),
    
    refresh: (refreshToken: string) =>
        api.post('/api/auth/refresh', { refresh_token: refreshToken }),
    
    changePassword: (oldPassword: string, newPassword: string) =>
        api.post('/api/auth/change-password', { old_password: oldPassword, new_password: newPassword }),
    
    me: () => api.get('/api/me'),
};

// Dashboard API
export const dashboardAPI = {
    getStats: () => api.get('/api/dashboard/stats'),
};

// Schedulers API
export const schedulersAPI = {
    getAll: () => api.get('/api/schedulers'),
    
    getById: (id: string) => api.get(`/api/schedulers/${id}`),
    
    create: (scheduler: any) => api.post('/api/schedulers', scheduler),
    
    update: (id: string, scheduler: any) => api.put(`/api/schedulers/${id}`, scheduler),
    
    delete: (id: string) => api.delete(`/api/schedulers/${id}`),
    
    run: (id: string) => api.post(`/api/schedulers/${id}/run`),
    
    getHistory: (id: string) => api.get(`/api/schedulers/${id}/history`),
};

// Users API
export const usersAPI = {
    getAll: () => api.get('/api/users'),
    
    create: (user: any) => api.post('/api/users', user),
    
    update: (id: string, user: any) => api.put(`/api/users/${id}`, user),
    
    delete: (id: string) => api.delete(`/api/users/${id}`),
    
    resetPassword: (id: string, newPassword: string) => 
        api.post(`/api/users/${id}/reset-password`, { newPassword }),
};

// Roles API
export const rolesAPI = {
    getAll: () => api.get('/api/roles'),
    
    create: (role: any) => api.post('/api/roles', role),
    
    update: (id: string, role: any) => api.put(`/api/roles/${id}`, role),
    
    delete: (id: string) => api.delete(`/api/roles/${id}`),
    
    getPermissions: () => api.get('/api/roles/permissions'),
};

// Images API
export const imagesAPI = {
    getAll: () => api.get('/api/images'),
    
    create: (image: any) => api.post('/api/images', image),
    
    update: (id: string, image: any) => api.put(`/api/images/${id}`, image),
    
    delete: (id: string) => api.delete(`/api/images/${id}`),
};

// Logs API
export const logsAPI = {
    getAll: (params?: any) => api.get('/api/logs', { params }),
    
    getStats: () => api.get('/api/logs/stats'),
    
    getSources: () => api.get('/api/logs/sources'),
};

// Monitoring API
export const monitoringAPI = {
    getAll: () => api.get('/api/monitoring'),
    
    getMetrics: () => api.get('/api/monitoring/metrics'),
    
    getServices: () => api.get('/api/monitoring/services'),
    
    getAlerts: () => api.get('/api/monitoring/alerts'),
};

export default api;