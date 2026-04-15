import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from 'react-query';
import { Link } from 'react-router-dom';
import { jobsAPI } from '../api/jobs';
import { FiPlus, FiLogOut, FiUser, FiSettings } from 'react-icons/fi';
import { motion } from 'framer-motion';

const STATUS_STYLE = {
  completed: 'text-green-600 bg-green-50',
  failed: 'text-red-600 bg-red-50',
  running: 'text-blue-600 bg-blue-50',
  pending: 'text-amber-600 bg-amber-50',
  queued: 'text-slate-600 bg-slate-100',
};

const Dashboard = () => {
  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  const [activeTab, setActiveTab] = useState('jobs');

  const { data: jobs, isLoading, error } = useQuery(
    'jobs',
    () => jobsAPI.getAllJobs().then(res => res.data),
    { refetchInterval: 15000, refetchOnWindowFocus: false, retry: 2 }
  );

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/');
  };

  const stats = {
    total: jobs?.length || 0,
    running: jobs?.filter(j => j.status === 'running').length || 0,
    completed: jobs?.filter(j => j.status === 'completed').length || 0,
    failed: jobs?.filter(j => j.status === 'failed').length || 0,
  };

  return (
    <div className="min-h-screen bg-white font-sans text-black flex">
      {/* Sidebar */}
      <aside className="hidden lg:flex flex-col w-64 border-r border-slate-50 sticky top-0 h-screen p-8">
        <div className="font-black text-xl tracking-tighter mb-16 cursor-pointer" onClick={() => navigate('/dashboard')}>
          SmartTask
        </div>
        <nav className="flex-1 space-y-2">
          {[
            { label: 'Dashboard', tab: 'jobs', icon: FiUser },
            { label: 'Configuration', tab: 'config', icon: FiSettings },
          ].map(({ label, tab, icon: Icon }) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-xl text-xs font-black uppercase tracking-wider transition-all ${
                activeTab === tab 
                ? 'bg-black text-white'
                : 'text-slate-300 hover:text-black hover:bg-slate-50'
              }`}
            >
              <Icon size={16} />
              {label}
            </button>
          ))}
        </nav>
        <div className="space-y-4">
          <div className="border-t border-slate-50 pt-6">
            <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-1">Signed in as</p>
            <p className="text-sm font-bold truncate">{user.name || user.email || 'User'}</p>
          </div>
          <button onClick={handleLogout} className="flex items-center gap-2 text-[10px] font-black text-slate-300 hover:text-red-500 uppercase tracking-widest transition-colors">
            <FiLogOut size={14} /> Log Out
          </button>
        </div>
      </aside>

      {/* Main */}
      <main className="flex-1 overflow-auto">
        {/* Header */}
        <header className="h-16 border-b border-slate-50 px-8 flex items-center justify-between sticky top-0 bg-white/95 backdrop-blur-sm z-30">
          <div>
            <h1 className="font-black text-lg">
              {activeTab === 'jobs' ? 'Execution Monitor' : 'Configuration'}
            </h1>
          </div>
          <Link
            to="/create"
            className="bg-black text-white px-5 py-2 rounded-lg text-[10px] font-black uppercase tracking-widest hover:opacity-80 transition-opacity flex items-center gap-2"
          >
            <FiPlus size={14} /> Deploy Task
          </Link>
        </header>

        <div className="p-8 lg:p-14">
          {activeTab === 'jobs' && (
            <>
              {/* Stats Row */}
              <div className="grid grid-cols-2 md:grid-cols-4 gap-10 mb-20">
                {[
                  { label: 'Total Jobs', val: stats.total },
                  { label: 'Running', val: stats.running },
                  { label: 'Completed', val: stats.completed },
                  { label: 'Failed', val: stats.failed },
                ].map((s, i) => (
                  <motion.div key={i} initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: i * 0.05 }}>
                    <p className="text-[9px] font-black text-slate-300 uppercase tracking-[0.2em] mb-2">{s.label}</p>
                    <h3 className="text-4xl font-black leading-none">{s.val}</h3>
                  </motion.div>
                ))}
              </div>

              {/* Jobs Table */}
              {isLoading && (
                <p className="text-[10px] font-black text-slate-200 uppercase tracking-widest animate-pulse">Synchronizing cluster...</p>
              )}

              {error && (
                <div className="py-10 text-center">
                  <p className="text-red-500 text-xs font-bold mb-2 uppercase tracking-widest">API Connection Failed</p>
                  <p className="text-slate-400 text-xs">Make sure the backend is running on port 8080.</p>
                </div>
              )}

              {jobs && (
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b-2 border-slate-900">
                        <th className="py-4 text-left text-[10px] font-black uppercase tracking-widest text-slate-400">Job Name</th>
                        <th className="py-4 text-left text-[10px] font-black uppercase tracking-widest text-slate-400">Type</th>
                        <th className="py-4 text-left text-[10px] font-black uppercase tracking-widest text-slate-400">Image</th>
                        <th className="py-4 text-left text-[10px] font-black uppercase tracking-widest text-slate-400">Status</th>
                        <th className="py-4 text-left text-[10px] font-black uppercase tracking-widest text-slate-400">Created</th>
                        <th className="py-4"></th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-50">
                      {jobs.map((job) => (
                        <tr key={job.id} className="group hover:bg-slate-50/50 transition-colors">
                          <td className="py-5">
                            <Link to={`/jobs/${job.id}`} className="font-bold text-sm hover:text-orange-500 transition-colors">
                              {job.name}
                            </Link>
                          </td>
                          <td className="py-5">
                            <span className="text-[10px] font-black uppercase tracking-widest text-slate-400">{job.type}</span>
                          </td>
                          <td className="py-5">
                            <span className="font-mono text-xs text-slate-500">{job.image}</span>
                          </td>
                          <td className="py-5">
                            <span className={`text-[10px] font-black uppercase tracking-widest px-2 py-1 rounded ${STATUS_STYLE[job.status] || STATUS_STYLE.pending}`}>
                              {job.status || 'pending'}
                            </span>
                          </td>
                          <td className="py-5 text-[11px] font-bold text-slate-300">
                            {new Date(job.created_at).toLocaleDateString()}
                          </td>
                          <td className="py-5 text-right">
                            <Link to={`/jobs/${job.id}`} className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-200 hover:text-black transition-colors">
                              Inspect →
                            </Link>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>

                  {jobs.length === 0 && (
                    <div className="py-32 text-center">
                      <p className="text-[10px] font-black text-slate-200 uppercase tracking-[0.3em] mb-8">No tasks have been deployed yet.</p>
                      <Link to="/create" className="bg-black text-white px-8 py-4 text-[10px] font-black uppercase tracking-widest hover:opacity-80 transition-opacity">
                        Deploy First Task
                      </Link>
                    </div>
                  )}
                </div>
              )}
            </>
          )}

          {activeTab === 'config' && (
            <div className="max-w-lg">
              <h2 className="text-2xl font-black mb-2">API Configuration</h2>
              <p className="text-sm text-slate-400 font-bold mb-12">Your API key can be used to trigger jobs programmatically.</p>
              <div className="space-y-3">
                <label className="text-[10px] font-black text-slate-300 uppercase tracking-widest">Your API Key</label>
                <div className="flex gap-4">
                  <input
                    readOnly
                    value={user.api_key || '—'}
                    className="flex-1 border-b-2 border-slate-100 py-3 font-mono text-sm font-bold outline-none bg-transparent text-slate-600"
                  />
                  <button 
                    onClick={() => navigator.clipboard.writeText(user.api_key || '')}
                    className="text-[10px] font-black uppercase tracking-widest text-slate-400 hover:text-black transition-colors"
                  >
                    Copy
                  </button>
                </div>
                <p className="text-[10px] text-slate-300 font-bold uppercase tracking-widest">Pass as <span className="font-mono">X-API-Key: your_key</span> header</p>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;