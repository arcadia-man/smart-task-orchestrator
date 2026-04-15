import React from 'react';
import { useNavigate } from 'react-router-dom';
import { FiExternalLink, FiGithub, FiLinkedin, FiMail } from 'react-icons/fi';
import { motion } from 'framer-motion';

const LandingPage = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-white">
      {/* Navbar */}
      <nav className="fixed top-0 left-0 right-0 z-50 bg-white border-b border-slate-100 h-16 flex items-center">
        <div className="max-w-6xl mx-auto px-6 w-full flex items-center justify-between">
          <span className="text-base font-extrabold tracking-tight cursor-pointer" onClick={() => navigate('/')}>
            SmartTask
          </span>
          
          <div className="hidden md:flex items-center gap-8">
            <a href="/" className="text-xs font-semibold text-black border-b border-black pb-0.5">Home</a>
            <a href="#about" className="text-xs font-medium text-slate-400 hover:text-black transition-colors">About</a>
            <a href="#contact" className="text-xs font-medium text-slate-400 hover:text-black transition-colors">Contact</a>
          </div>

          <div className="flex items-center gap-4">
            <button 
              onClick={() => navigate('/signin')}
              className="text-xs font-semibold text-slate-500 hover:text-black transition-colors"
            >
              Sign In
            </button>
            <button 
              onClick={() => navigate('/signup')}
              className="bg-black text-white px-4 py-2 rounded-lg text-xs font-bold hover:opacity-80 transition-opacity"
            >
              Get Started
            </button>
          </div>
        </div>
      </nav>

      <main className="pt-16">
        {/* Hero */}
        <section className="max-w-3xl mx-auto px-6 py-28 text-center">
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            {/* Eyebrow */}
            <span className="inline-block mb-6 text-[10px] font-black uppercase tracking-[0.3em] text-slate-400 border border-slate-200 rounded-full px-4 py-1.5">
              Open Source · Self-Hostable
            </span>
            <h1 className="text-4xl lg:text-5xl font-black text-black mb-6 leading-[1.1] tracking-tight">
              A modern sandbox for <br className="hidden md:block" />
              task orchestration.
            </h1>
            <p className="text-sm text-slate-500 mb-10 max-w-xl mx-auto leading-relaxed font-medium">
              Execute scripts, manage cron jobs, and scale workloads in a secure, isolated container environment — powered by Go, Kafka, and Docker.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
               <button 
                 onClick={() => navigate('/signup')}
                 className="bg-black text-white px-7 py-3 rounded-lg font-bold text-sm hover:opacity-85 transition-all w-full sm:w-auto"
               >
                 Start Building
               </button>
               <button className="text-xs font-bold flex items-center gap-1.5 text-slate-500 hover:text-black transition-colors">
                 View Docs <FiExternalLink size={14} />
               </button>
            </div>
          </motion.div>
        </section>

        {/* Feature Divider */}
        <div className="border-t border-slate-50" />

        {/* Features */}
        <section id="about" className="max-w-6xl mx-auto px-6 py-24 grid md:grid-cols-3 gap-16">
          {[
            {
              num: '01',
              title: 'Container Isolation',
              desc: 'Every task spins up a fresh container instance. Total isolation with enforced resource quotas.'
            },
            {
              num: '02',
              title: 'Event-Driven',
              desc: 'Powered by Kafka for high-throughput, sub-second job scheduling and real-time status streaming.'
            },
            {
              num: '03',
              title: 'Full Observability',
              desc: 'Capture stdout/stderr per execution. Store logs in MongoDB. Query any historic run instantly.'
            },
          ].map(f => (
            <div key={f.num}>
              <p className="text-[9px] font-black text-slate-300 uppercase tracking-[0.3em] mb-3">{f.num}</p>
              <h3 className="text-sm font-black mb-3 text-black">{f.title}</h3>
              <p className="text-xs text-slate-400 leading-relaxed font-medium">{f.desc}</p>
            </div>
          ))}
        </section>

        {/* CTA Banner */}
        <section className="bg-black mx-6 lg:mx-auto max-w-6xl rounded-2xl mb-20 p-12 flex flex-col md:flex-row items-center justify-between gap-6">
          <div>
            <h2 className="text-xl font-black text-white mb-1">Ready to orchestrate?</h2>
            <p className="text-xs text-slate-400 font-medium">Free to use. Open source. Deploy anywhere.</p>
          </div>
          <button 
            onClick={() => navigate('/signup')}
            className="bg-white text-black px-7 py-3 rounded-lg font-bold text-sm hover:opacity-90 transition-all shrink-0"
          >
            Create Free Account
          </button>
        </section>
      </main>

      {/* Footer */}
      <footer id="contact" className="py-12 border-t border-slate-50 bg-white">
        <div className="max-w-6xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
           <p className="text-[10px] font-bold text-slate-300 uppercase tracking-widest">
             © 2026 Pritam Kumar Maurya · Go · Kafka · Mongo · Docker
           </p>
           <div className="flex gap-6">
              <FiGithub className="text-slate-300 hover:text-black cursor-pointer transition-colors" size={18} />
              <FiLinkedin className="text-slate-300 hover:text-black cursor-pointer transition-colors" size={18} />
              <FiMail className="text-slate-300 hover:text-black cursor-pointer transition-colors" size={18} />
           </div>
        </div>
      </footer>
    </div>
  );
};

export default LandingPage;
