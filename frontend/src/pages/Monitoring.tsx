import React from 'react';
import { Activity, Cpu, HardDrive, MemoryStick, Network, Server, AlertTriangle, CheckCircle } from 'lucide-react';

const Monitoring: React.FC = () => {
 // Mock data - will be replaced with real API calls
 const systemMetrics = {
  cpu: { usage: 45, cores: 8 },
  memory: { used: 6.2, total: 16, percentage: 39 },
  disk: { used: 120, total: 500, percentage: 24 },
  network: { inbound: 1.2, outbound: 0.8 }
 };

 const services = [
  { name: 'API Server', status: 'healthy', uptime: '99.9%', responseTime: '45ms' },
  { name: 'Database', status: 'healthy', uptime: '99.8%', responseTime: '12ms' },
  { name: 'Redis Cache', status: 'healthy', uptime: '100%', responseTime: '3ms' },
  { name: 'Scheduler Engine', status: 'warning', uptime: '98.5%', responseTime: '120ms' },
  { name: 'Log Aggregator', status: 'healthy', uptime: '99.7%', responseTime: '8ms' }
 ];

 const recentAlerts = [
  {
   id: '1',
   type: 'warning',
   message: 'High CPU usage detected on scheduler engine',
   timestamp: '2025-10-22T10:25:00Z',
   resolved: false
  },
  {
   id: '2',
   type: 'info',
   message: 'Database backup completed successfully',
   timestamp: '2025-10-22T02:00:00Z',
   resolved: true
  },
  {
   id: '3',
   type: 'error',
   message: 'Failed to connect to external API',
   timestamp: '2025-10-21T18:30:00Z',
   resolved: true
  }
 ];

 const getStatusIcon = (status: string) => {
  switch (status) {
   case 'healthy':
    return <CheckCircle className="w-4 h-4 text-green-500" />;
   case 'warning':
    return <AlertTriangle className="w-4 h-4 text-yellow-500" />;
   case 'error':
    return <AlertTriangle className="w-4 h-4 text-red-500" />;
   default:
    return <CheckCircle className="w-4 h-4 text-gray-400" />;
  }
 };

 const getStatusBadge = (status: string) => {
  const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
  switch (status) {
   case 'healthy':
    return `${baseClasses} bg-green-100 text-green-800`;
   case 'warning':
    return `${baseClasses} bg-yellow-100 text-yellow-800`;
   case 'error':
    return `${baseClasses} bg-red-100 text-red-800`;
   default:
    return `${baseClasses} bg-gray-100 text-gray-800`;
  }
 };

 const getProgressBarColor = (percentage: number) => {
  if (percentage >= 80) return 'bg-red-500';
  if (percentage >= 60) return 'bg-yellow-500';
  return 'bg-green-500';
 };

 return (
  <div className="space-y-6">
   {/* Header */}
   <div>
    <h1 className="text-3xl font-bold text-gray-900">System Monitoring</h1>
    <p className="mt-2 text-gray-600">
     Monitor system health and performance metrics
    </p>
   </div>

   {/* System Metrics */}
   <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    <div className="card">
     <div className="flex items-center justify-between mb-4">
      <div className="flex items-center">
       <div className="p-2 bg-blue-100 rounded-lg">
        <Cpu className="w-6 h-6 text-blue-600" />
       </div>
       <div className="ml-3">
        <p className="text-sm font-medium text-gray-600">CPU Usage</p>
        <p className="text-2xl font-bold text-gray-900">{systemMetrics.cpu.usage}%</p>
       </div>
      </div>
     </div>
     <div className="w-full bg-gray-200 rounded-full h-2">
      <div
       className={`h-2 rounded-full ${getProgressBarColor(systemMetrics.cpu.usage)}`}
       style={{ width: `${systemMetrics.cpu.usage}%` }}
      ></div>
     </div>
     <p className="text-xs text-gray-500 mt-2">{systemMetrics.cpu.cores} cores available</p>
    </div>

    <div className="card">
     <div className="flex items-center justify-between mb-4">
      <div className="flex items-center">
       <div className="p-2 bg-green-100 rounded-lg">
        <MemoryStick className="w-6 h-6 text-green-600" />
       </div>
       <div className="ml-3">
        <p className="text-sm font-medium text-gray-600">Memory</p>
        <p className="text-2xl font-bold text-gray-900">{systemMetrics.memory.percentage}%</p>
       </div>
      </div>
     </div>
     <div className="w-full bg-gray-200 rounded-full h-2">
      <div
       className={`h-2 rounded-full ${getProgressBarColor(systemMetrics.memory.percentage)}`}
       style={{ width: `${systemMetrics.memory.percentage}%` }}
      ></div>
     </div>
     <p className="text-xs text-gray-500 mt-2">
      {systemMetrics.memory.used}GB / {systemMetrics.memory.total}GB
     </p>
    </div>

    <div className="card">
     <div className="flex items-center justify-between mb-4">
      <div className="flex items-center">
       <div className="p-2 bg-purple-100 rounded-lg">
        <HardDrive className="w-6 h-6 text-purple-600" />
       </div>
       <div className="ml-3">
        <p className="text-sm font-medium text-gray-600">Disk Usage</p>
        <p className="text-2xl font-bold text-gray-900">{systemMetrics.disk.percentage}%</p>
       </div>
      </div>
     </div>
     <div className="w-full bg-gray-200 rounded-full h-2">
      <div
       className={`h-2 rounded-full ${getProgressBarColor(systemMetrics.disk.percentage)}`}
       style={{ width: `${systemMetrics.disk.percentage}%` }}
      ></div>
     </div>
     <p className="text-xs text-gray-500 mt-2">
      {systemMetrics.disk.used}GB / {systemMetrics.disk.total}GB
     </p>
    </div>

    <div className="card">
     <div className="flex items-center justify-between mb-4">
      <div className="flex items-center">
       <div className="p-2 bg-orange-100 rounded-lg">
        <Network className="w-6 h-6 text-orange-600" />
       </div>
       <div className="ml-3">
        <p className="text-sm font-medium text-gray-600">Network</p>
        <p className="text-2xl font-bold text-gray-900">
         {systemMetrics.network.inbound + systemMetrics.network.outbound} MB/s
        </p>
       </div>
      </div>
     </div>
     <div className="flex justify-between text-xs text-gray-500">
      <span>↓ {systemMetrics.network.inbound} MB/s</span>
      <span>↑ {systemMetrics.network.outbound} MB/s</span>
     </div>
    </div>
   </div>

   {/* Services Status */}
   <div className="card">
    <h2 className="text-xl font-semibold text-gray-900 mb-4">Service Health</h2>
    <div className="overflow-x-auto">
     <table className="w-full">
      <thead>
       <tr className="border-b border-gray-200">
        <th className="text-left py-3 px-4 font-medium text-gray-600">Service</th>
        <th className="text-left py-3 px-4 font-medium text-gray-600">Status</th>
        <th className="text-left py-3 px-4 font-medium text-gray-600">Uptime</th>
        <th className="text-left py-3 px-4 font-medium text-gray-600">Response Time</th>
       </tr>
      </thead>
      <tbody>
       {services.map((service, index) => (
        <tr key={index} className="border-b border-gray-100 hover:bg-gray-50">
         <td className="py-4 px-4">
          <div className="flex items-center">
           <Server className="w-4 h-4 text-gray-400 mr-3" />
           <span className="font-medium text-gray-900">{service.name}</span>
          </div>
         </td>
         <td className="py-4 px-4">
          <div className="flex items-center">
           {getStatusIcon(service.status)}
           <span className={`ml-2 ${getStatusBadge(service.status)}`}>
            {service.status}
           </span>
          </div>
         </td>
         <td className="py-4 px-4">
          <span className="text-sm text-gray-600">{service.uptime}</span>
         </td>
         <td className="py-4 px-4">
          <span className="text-sm text-gray-600">{service.responseTime}</span>
         </td>
        </tr>
       ))}
      </tbody>
     </table>
    </div>
   </div>

   {/* Recent Alerts */}
   <div className="card">
    <h2 className="text-xl font-semibold text-gray-900 mb-4">Recent Alerts</h2>
    <div className="space-y-4">
     {recentAlerts.map((alert) => (
      <div key={alert.id} className="flex items-start space-x-3 p-3 bg-gray-50 rounded-lg">
       {alert.type === 'warning' && <AlertTriangle className="w-5 h-5 text-yellow-500 mt-0.5" />}
       {alert.type === 'error' && <AlertTriangle className="w-5 h-5 text-red-500 mt-0.5" />}
       {alert.type === 'info' && <CheckCircle className="w-5 h-5 text-blue-500 mt-0.5" />}
       <div className="flex-1">
        <p className="text-gray-900 font-medium">{alert.message}</p>
        <div className="flex items-center space-x-2 mt-1">
         <span className="text-sm text-gray-500">
          {new Date(alert.timestamp).toLocaleString()}
         </span>
         {alert.resolved && (
          <span className="px-2 py-1 text-xs bg-green-100 text-green-800 rounded">
           Resolved
          </span>
         )}
        </div>
       </div>
      </div>
     ))}
    </div>
   </div>
  </div>
 );
};

export default Monitoring;