import React from 'react'
import { useQuery } from 'react-query'
import { Link } from 'react-router-dom'
import { jobsAPI } from '../api/jobs'

function Dashboard() {
  const { data: jobs, isLoading, error, refetch } = useQuery(
    'jobs',
    () => jobsAPI.getAllJobs().then(res => res.data),
    { 
      refetchInterval: 30000, // Auto-refresh every 30 seconds (reduced frequency)
      retry: 3, // Only retry 3 times
      retryDelay: 5000, // Wait 5 seconds between retries
      refetchOnWindowFocus: false // Don't refetch when window gains focus
    }
  )

  const handleRetry = async (jobId) => {
    try {
      await jobsAPI.retryJob(jobId)
      refetch()
    } catch (error) {
      console.error('Failed to retry job:', error)
    }
  }

  const getStatusColor = (status) => {
    const colors = {
      scheduled: 'bg-yellow-100 text-yellow-800',
      queued: 'bg-blue-100 text-blue-800',
      running: 'bg-purple-100 text-purple-800',
      completed: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800'
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (isLoading) return <div className="text-center py-8">Loading jobs...</div>
  if (error) return <div className="text-center py-8 text-red-600">Error loading jobs</div>

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Job Dashboard</h1>
        <Link 
          to="/create"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors"
        >
          Create New Job
        </Link>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Type
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Retry Count
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Created At
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {jobs?.map((job) => (
              <tr key={job.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <Link 
                    to={`/jobs/${job.id}`}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    {job.name}
                  </Link>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(job.status)}`}>
                    {job.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {job.type}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {job.retryCount} / {job.maxRetries}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(job.createdAt).toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  {job.status === 'failed' && job.retryCount < job.maxRetries && (
                    <button
                      onClick={() => handleRetry(job.id)}
                      className="text-blue-600 hover:text-blue-900 mr-4"
                    >
                      Retry
                    </button>
                  )}
                  <Link 
                    to={`/jobs/${job.id}`}
                    className="text-indigo-600 hover:text-indigo-900"
                  >
                    View Details
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        
        {jobs?.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            No jobs found. <Link to="/create" className="text-blue-600">Create your first job</Link>
          </div>
        )}
      </div>
    </div>
  )
}

export default Dashboard