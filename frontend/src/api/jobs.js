import axios from 'axios'

const API_BASE = '/api'

export const jobsAPI = {
  getAllJobs: () => axios.get(`${API_BASE}/jobs`),
  getJob: (id) => axios.get(`${API_BASE}/jobs/${id}`),
  createJob: (data) => axios.post(`${API_BASE}/jobs`, data),
  retryJob: (id) => axios.post(`${API_BASE}/jobs/${id}/retry`),
  getJobStatus: (id) => axios.get(`${API_BASE}/jobs/${id}/status`)
}