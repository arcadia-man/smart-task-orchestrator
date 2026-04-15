import React from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useQuery } from 'react-query';
import { jobsAPI } from '../api/jobs';
import { FiArrowLeft } from 'react-icons/fi';
import { motion } from 'framer-motion';

const STATUS_STYLE = {
  completed: 'text-green-600 bg-green-50',
  failed: 'text-red-600 bg-red-50',
  running: 'text-blue-600 bg-blue-50',
  pending: 'text-amber-600 bg-amber-50',
  queued: 'text-slate-600 bg-slate-100',
};

const JobDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  
  const { data: job, isLoading, error } = useQuery(
    ['job', id],
    () => jobsAPI.getJob(id).then(res => res.data),
    { refetchInterval: 10000, retry: 2 }
  );

  if (isLoading) return (
    <div className="min-h-screen flex items-center justify-center bg-white">
      <p className="text-[10px] font-black uppercase tracking-[0.4em] text-slate-200 animate-pulse">Retrieving execution logs...</p>
    </div>
  );

  if (error || !job) return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-white">
      <p className="text-red-500 font-bold mb-6 uppercase tracking-widest text-xs">Job not found or unauthorized</p>
      <button onClick={() => navigate('/dashboard')} className="bg-black text-white px-8 py-3 text-[10px] font-black uppercase tracking-widest hover:opacity-80">
        Return to Dashboard
      </button>
    </div>
  );

  return (
    <div className="min-h-screen bg-white font-sans text-black py-14 px-8 lg:px-20">
      <div className="max-w-5xl mx-auto">
        <button 
          onClick={() => navigate('/dashboard')}
          className="flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-slate-300 hover:text-black transition-colors mb-20"
        >
          <FiArrowLeft size={16} /> Back to Dashboard
        </button>

        {/* Job Header */}
        <header className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-24 pb-12 border-b border-slate-50">
           <div>
              <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Job ID: {job.id}</p>
              <h1 className="text-4xl lg:text-6xl font-black tracking-tight">{job.name}</h1>
           </div>
           <span className={`text-[10px] font-black uppercase tracking-widest px-4 py-2 rounded-full border ${STATUS_STYLE[job.status] || STATUS_STYLE.pending}`}>
              {job.status || 'pending'}
           </span>
        </header>

        <div className="grid lg:grid-cols-3 gap-20">
           {/* Main Info */}
           <div className="lg:col-span-2 space-y-16">
              <section>
                 <h2 className="text-[10px] font-black uppercase tracking-[0.2em] mb-10 pb-4 border-b-2 border-black inline-block">Execution Metadata</h2>
                 <div className="grid grid-cols-2 gap-10">
                    <div>
                       <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Job Type</p>
                       <p className="font-bold uppercase tracking-wider text-sm">{job.type}</p>
                    </div>
                    <div>
                       <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Docker Image</p>
                       <p className="font-mono text-sm font-bold">{job.image}</p>
                    </div>
                    <div>
                       <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Command</p>
                       <p className="font-mono text-xs font-bold bg-slate-50 px-3 py-2 rounded-lg">{job.command}</p>
                    </div>
                    {job.cron_expr && (
                       <div>
                          <p className="text-[10px] font-black text-orange-500 uppercase tracking-widest mb-2">Cron Schedule</p>
                          <p className="font-mono text-sm font-bold text-orange-500">{job.cron_expr}</p>
                       </div>
                    )}
                    <div>
                       <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Created</p>
                       <p className="font-bold text-sm">{new Date(job.created_at).toLocaleString()}</p>
                    </div>
                    <div>
                       <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Last Updated</p>
                       <p className="font-bold text-sm">{new Date(job.updated_at).toLocaleString()}</p>
                    </div>
                 </div>
              </section>

              {/* Sandbox Config */}
              {job.scaling && (
                 <section>
                    <h2 className="text-[10px] font-black uppercase tracking-[0.2em] mb-10 pb-4 border-b-2 border-black inline-block">Sandbox Scaling Config</h2>
                    <div className="grid grid-cols-2 gap-10">
                       <div>
                          <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Min Containers</p>
                          <p className="text-4xl font-black">{job.scaling.min_containers}</p>
                       </div>
                       <div>
                          <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest mb-2">Max Containers</p>
                          <p className="text-4xl font-black">{job.scaling.max_containers}</p>
                       </div>
                    </div>
                 </section>
              )}
           </div>

           {/* Sidebar */}
           <div className="space-y-10">
              <div className="border border-slate-50 rounded-2xl p-8 bg-slate-50/30">
                 <h3 className="text-[10px] font-black uppercase tracking-[0.2em] mb-8 text-slate-400">Quick Actions</h3>
                 <div className="space-y-4">
                    <button className="w-full py-4 bg-black text-white text-[10px] font-black uppercase tracking-widest hover:opacity-80 transition-all rounded-xl">
                       Trigger Manually
                    </button>
                    <button 
                       onClick={() => navigate('/create')}
                       className="w-full py-4 border border-slate-200 text-[10px] font-black uppercase tracking-widest text-slate-400 hover:text-black hover:border-black transition-all rounded-xl"
                    >
                       Clone Config
                    </button>
                 </div>
              </div>

              <div className="border border-dashed border-slate-100 rounded-2xl p-8">
                 <h3 className="text-[10px] font-black uppercase tracking-[0.2em] mb-4 text-slate-300">Kafka Metadata</h3>
                 <p className="text-[10px] font-bold text-slate-300 leading-relaxed">
                    Events for this job are dispatched via topic <span className="font-mono text-black">jobs.execute</span> on the cluster broker.
                 </p>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
};

export default JobDetails;