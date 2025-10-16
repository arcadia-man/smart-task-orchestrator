import React from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery } from 'react-query'
import { jobsAPI } from '../api/jobs'

function JobDetails() {
  const { id } = useParams()
  
  const { data: job, isLoading, error } = useQuery(
    ['job', id],
    () => jobsAPI.getJob(id).then(res => res.data),
    { refetchInterval: 10000 }
  )

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

  if (isLoading) return <div className="text-center py-8">Loading job details...</div>
  if (error) return <div className="text-center py-8 text-red-600">Error loading job details</div>
  if (!job) return <div className="text-center py-8">Job not found</div>

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <Link to="/" className="text-blue-600 hover:text-blue-800 mb-4 inline-block">
          ← Back to Dashboard
        </Link>
        <h1 className="text-3xl font-bold text-gray-900">{job.name}</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Job Info */}
        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Job Information</h2>
          <div className="space-y-3">
            <div>
              <span className="font-medium text-gray-700">Status:</span>
              <span className={`ml-2 px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(job.status)}`}>
                {job.status}
              </span>
            </div>
            <div>
              <span className="font-medium text-gray-700">Type:</span>
              <span className="ml-2 text-gray-900">{job.type}</span>
            </div>
            <div>
              <span className="font-medium text-gray-700">Retry Count:</span>
              <span className="ml-2 text-gray-900">{job.retryCount} / {job.maxRetries}</span>
            </div>
            {job.cronExpr && (
              <div>
                <span className="font-medium text-gray-700">Cron Expression:</span>
                <span className="ml-2 text-gray-900 font-mono">{job.cronExpr}</span>
              </div>
            )}
            {job.nextRunAt && (
              <div>
                <span className="font-medium text-gray-700">Next Run:</span>
                <span className="ml-2 text-gray-900">{new Date(job.nextRunAt).toLocaleString()}</span>
              </div>
            )}
            <div>
              <span className="font-medium text-gray-700">Created:</span>
              <span className="ml-2 text-gray-900">{new Date(job.createdAt).toLocaleString()}</span>
            </div>
            <div>
              <span className="font-medium text-gray-700">Updated:</span>
              <span className="ml-2 text-gray-900">{new Date(job.updatedAt).toLocaleString()}</span>
            </div>
          </div>
        </div>

        {/* Payload */}
        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Payload</h2>
          <pre className="bg-gray-100 p-4 rounded text-sm overflow-auto">
            {JSON.stringify(job.payload, null, 2)}
          </pre>
        </div>
      </div>

      {/* History */}
      <div className="mt-6 bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Execution History</h2>
        <div className="space-y-3">
          {job.history?.map((event, index) => (
            <div key={index} className="border-l-4 border-blue-500 pl-4 py-2">
              <div className="flex justify-between items-start">
                <div>
                  <span className="font-medium text-gray-900">{event.event}</span>
                  <p className="text-gray-600 text-sm">{event.message}</p>
                </div>
                <span className="text-sm text-gray-500">
                  {new Date(event.timestamp).toLocaleString()}
                </span>
              </div>
            </div>
          ))}
          
          {(!job.history || job.history.length === 0) && (
            <p className="text-gray-500 text-center py-4">No history available</p>
          )}
        </div>
      </div>
    </div>
  )
}

export default JobDetails