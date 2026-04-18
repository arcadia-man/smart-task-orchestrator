import axios from 'axios';

const API_BASE = '/api';

// Attach JWT token to every request automatically
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  signup: (data) => axios.post(`${API_BASE}/signup`, data),
  login: (data) => axios.post(`${API_BASE}/login`, data),
};

export const jobsAPI = {
  getAllJobs: () => axios.get(`${API_BASE}/jobs`),
  getJob: (id) => axios.get(`${API_BASE}/jobs/${id}`),
  createJob: (data) => axios.post(`${API_BASE}/jobs`, data),
  getExecutions: (jobId) => axios.get(`${API_BASE}/jobs/${jobId}/executions`),
};

export const userAPI = {
  getProfile: () => axios.get(`${API_BASE}/profile`),
  updateProfile: (data) => axios.put(`${API_BASE}/profile`, data),
  resetPassword: (data) => axios.post(`${API_BASE}/reset-password`, data),
};