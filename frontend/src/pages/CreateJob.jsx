import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jobsAPI } from '../api/jobs';
import { FiArrowLeft, FiPlay } from 'react-icons/fi';
import { motion } from 'framer-motion';

const CreateJob = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    type: 'one-time',
    image: 'alpine:latest',
    command: 'echo "Hello World"',
    cron_expr: '',
    scaling: { min_containers: 1, max_containers: 5 },
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');
    try {
      const payload = {
        name: formData.name,
        type: formData.type,
        image: formData.image,
        command: formData.command,
        ...(formData.type === 'cron' && { cron_expr: formData.cron_expr }),
        ...(formData.type === 'sandbox' && { scaling: formData.scaling }),
      };
      await jobsAPI.createJob(payload);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create job. Check all fields and try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleChange = (e) => setFormData({ ...formData, [e.target.name]: e.target.value });
  const handleScalingChange = (e) => setFormData({ ...formData, scaling: { ...formData.scaling, [e.target.name]: parseInt(e.target.value) } });

  return (
    <div className="min-h-screen bg-white font-sans text-black p-8 lg:p-20">
      <div className="max-w-2xl mx-auto">
        <button 
          onClick={() => navigate('/dashboard')}
          className="flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-slate-300 hover:text-black transition-colors mb-20"
        >
          <FiArrowLeft size={16} /> Back to Dashboard
        </button>

        <div className="mb-20">
          <h1 className="text-4xl font-black mb-2">Deploy New Task</h1>
          <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Configure your execution environment</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-14">
          {/* Name */}
          <div className="space-y-3">
            <label className="text-[10px] font-black text-slate-300 uppercase tracking-widest">Job Name</label>
            <input
              type="text"
              name="name"
              required
              placeholder="data-pipeline-v1"
              value={formData.name}
              onChange={handleChange}
              className="w-full border-b-2 border-slate-100 py-4 text-xl font-bold outline-none focus:border-black transition-colors bg-transparent"
            />
          </div>

          {/* Type */}
          <div className="space-y-3">
            <label className="text-[10px] font-black text-slate-300 uppercase tracking-widest">Execution Type</label>
            <div className="grid grid-cols-3 gap-3">
              {[
                { val: 'one-time', label: 'One-Time' },
                { val: 'cron', label: 'Scheduled' },
                { val: 'sandbox', label: 'Sandbox Pool' },
              ].map(opt => (
                <button
                  key={opt.val}
                  type="button"
                  onClick={() => setFormData({ ...formData, type: opt.val })}
                  className={`py-3 px-2 border-2 text-[10px] font-black uppercase tracking-wider transition-all ${
                    formData.type === opt.val ? 'border-black bg-black text-white' : 'border-slate-100 text-slate-400 hover:border-black hover:text-black'
                  }`}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          </div>

          {/* Docker Image */}
          <div className="space-y-3">
            <label className="text-[10px] font-black text-slate-300 uppercase tracking-widest">Docker Image</label>
            <input
              type="text"
              name="image"
              required
              placeholder="alpine:latest"
              value={formData.image}
              onChange={handleChange}
              className="w-full border-b-2 border-slate-100 py-4 font-mono font-bold outline-none focus:border-black transition-colors bg-transparent"
            />
            <p className="text-[10px] text-slate-300 font-bold uppercase tracking-wider">Any public Docker Hub image or private registry URI</p>
          </div>

          {/* Shell Command */}
          <div className="space-y-3">
            <label className="text-[10px] font-black text-slate-300 uppercase tracking-widest">Shell Command</label>
            <textarea
              name="command"
              required
              rows={3}
              placeholder='sh -c "echo hello && python script.py"'
              value={formData.command}
              onChange={handleChange}
              className="w-full border border-slate-100 rounded-xl p-5 font-mono text-sm font-bold outline-none focus:border-black transition-colors bg-slate-50/50"
            />
          </div>

          {/* Cron Expression (conditional) */}
          {formData.type === 'cron' && (
            <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} className="space-y-3">
              <label className="text-[10px] font-black text-orange-500 uppercase tracking-widest">Cron Expression</label>
              <input
                type="text"
                name="cron_expr"
                placeholder="0 */5 * * *"
                value={formData.cron_expr}
                onChange={handleChange}
                className="w-full border-b-2 border-orange-200 py-4 font-mono font-bold outline-none focus:border-orange-500 transition-colors bg-transparent"
              />
              <p className="text-[10px] text-slate-300 font-bold uppercase tracking-wider">Standard 5-field cron: min hour day month weekday</p>
            </motion.div>
          )}

          {/* Sandbox Scaling config (conditional) */}
          {formData.type === 'sandbox' && (
            <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} className="space-y-3">
              <label className="text-[10px] font-black text-blue-500 uppercase tracking-widest">Container Scaling</label>
              <div className="grid grid-cols-2 gap-8">
                <div className="space-y-2">
                  <label className="text-[9px] font-black text-slate-300 uppercase tracking-widest">Min Containers</label>
                  <input type="number" name="min_containers" min={1} max={100} value={formData.scaling.min_containers} onChange={handleScalingChange}
                    className="w-full border-b-2 border-blue-100 py-3 font-bold outline-none focus:border-blue-500 bg-transparent" />
                </div>
                <div className="space-y-2">
                  <label className="text-[9px] font-black text-slate-300 uppercase tracking-widest">Max Containers</label>
                  <input type="number" name="max_containers" min={1} max={1000} value={formData.scaling.max_containers} onChange={handleScalingChange}
                    className="w-full border-b-2 border-blue-100 py-3 font-bold outline-none focus:border-blue-500 bg-transparent" />
                </div>
              </div>
            </motion.div>
          )}

          {error && <p className="text-red-500 text-[10px] font-black uppercase tracking-widest">{error}</p>}

          <div className="flex gap-6 pt-6">
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 bg-black text-white py-5 rounded-xl font-bold uppercase tracking-[0.2em] text-[10px] hover:opacity-85 transition-all shadow-2xl shadow-black/10 flex items-center justify-center gap-3 disabled:opacity-50"
            >
              <FiPlay size={16} />
              {isSubmitting ? 'Deploying...' : 'Confirm Deployment'}
            </button>
            <button
              type="button"
              onClick={() => navigate('/dashboard')}
              className="flex-1 border-2 border-slate-100 py-5 rounded-xl font-bold uppercase tracking-[0.2em] text-[10px] text-slate-400 hover:text-black hover:border-black transition-all"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateJob;