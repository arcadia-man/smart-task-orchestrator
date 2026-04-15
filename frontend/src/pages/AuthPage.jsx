import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../api/jobs';
import { FiArrowRight, FiMail, FiLock, FiUser, FiPhone } from 'react-icons/fi';
import { motion } from 'framer-motion';

const AuthPage = ({ mode = 'signin' }) => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({ name: '', email: '', phone: '', password: '' });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e) => setFormData({ ...formData, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      if (mode === 'signup') {
        await authAPI.signup({ name: formData.name, email: formData.email, phone: formData.phone, password: formData.password });
        // After signup, auto login
        const res = await authAPI.login({ email: formData.email, password: formData.password });
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('user', JSON.stringify(res.data.user));
      } else {
        const res = await authAPI.login({ email: formData.email, password: formData.password });
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('user', JSON.stringify(res.data.user));
      }
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || 'Something went wrong. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center p-6 font-sans">
      <div className="w-full max-w-[360px]">
        <div className="mb-20 text-center">
          <div 
             onClick={() => navigate('/')}
             className="text-xl font-black tracking-tighter cursor-pointer mb-12 inline-block"
          >
            SmartTask
          </div>
          <h1 className="text-4xl font-black text-black mb-3">
            {mode === 'signin' ? 'Welcome back.' : 'Get started.'}
          </h1>
          <p className="text-[11px] font-black text-slate-300 uppercase tracking-[0.2em]">
            {mode === 'signin' ? "Sign in to your workspace" : "Create your account"}
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-10">
          {mode === 'signup' && (
            <>
              <div className="space-y-3">
                <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none">Full Name</label>
                <input 
                  type="text" 
                  name="name"
                  required
                  placeholder="Pritam Kumar" 
                  value={formData.name}
                  onChange={handleChange}
                  className="w-full border-b-2 border-slate-100 py-3 text-black font-bold placeholder:text-slate-200 outline-none focus:border-black transition-colors bg-transparent"
                />
              </div>
              <div className="space-y-3">
                <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none">Phone Number</label>
                <input 
                  type="tel" 
                  name="phone"
                  placeholder="+91 98765 43210" 
                  value={formData.phone}
                  onChange={handleChange}
                  className="w-full border-b-2 border-slate-100 py-3 text-black font-bold placeholder:text-slate-200 outline-none focus:border-black transition-colors bg-transparent"
                />
              </div>
            </>
          )}
          
          <div className="space-y-3">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none">Email Address</label>
            <input 
              type="email" 
              name="email"
              required
              placeholder="operator@system.io" 
              value={formData.email}
              onChange={handleChange}
              className="w-full border-b-2 border-slate-100 py-3 text-black font-bold placeholder:text-slate-200 outline-none focus:border-black transition-colors bg-transparent"
            />
          </div>

          <div className="space-y-3">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none">Password</label>
            <input 
              type="password" 
              name="password"
              required
              placeholder="••••••••" 
              value={formData.password}
              onChange={handleChange}
              className="w-full border-b-2 border-slate-100 py-3 text-black font-bold placeholder:text-slate-200 outline-none focus:border-black transition-colors bg-transparent"
            />
          </div>

          {error && (
            <p className="text-red-500 text-[10px] font-black uppercase tracking-widest">{error}</p>
          )}

          <button 
            type="submit"
            disabled={loading}
            className="w-full bg-black text-white py-5 rounded-xl font-bold text-sm hover:opacity-90 transition-all mt-6 shadow-2xl shadow-black/10 flex items-center justify-center gap-3 uppercase tracking-widest disabled:opacity-50"
          >
            {loading ? 'Processing...' : (mode === 'signin' ? 'Sign In' : 'Create Account')}
            <FiArrowRight size={18} />
          </button>
        </form>

        <div className="mt-16 text-center">
           <button 
             onClick={() => navigate(mode === 'signin' ? '/signup' : '/signin')}
             className="text-[10px] font-black text-slate-300 hover:text-black transition-colors uppercase tracking-[0.2em] leading-none"
           >
             {mode === 'signin' ? "New here? Create Account" : "Have an account? Sign In"}
           </button>
        </div>
      </div>
    </div>
  );
};

export default AuthPage;
