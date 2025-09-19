import { Activity, Clock, DollarSign, Zap } from 'lucide-react';
import { useEffect, useState } from 'react';
import {
    Bar,
    BarChart,
    CartesianGrid,
    Legend,
    Line,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis, YAxis
} from 'recharts';
import Layout from '../../components/Layout';
import MetricCard from '../../components/MetricCard';
import { apiFetch } from '../../components/apiClient';

export default function Dashboard() {
  const [metrics, setMetrics] = useState({
    totalRequests: 0,
    avgLatency: 0,
    costSavings: 0,
    activeProviders: 0,
    requests: [],
    costs: [],
  });

  useEffect(() => {
    const load = async () => {
      try {
        const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined;
        const data = await apiFetch('/metrics', { token });
        setMetrics({
          totalRequests: data.key_metrics?.total_requests || 0,
          avgLatency: data.key_metrics?.avg_latency || 0,
          costSavings: data.key_metrics?.cost_savings || 0,
          activeProviders: (data.providers || []).length,
          requests: (data.requests || []).map(r => ({ time: r.time, requests: r.requests })),
          costs: (data.cost_data || []).map(d => ({ day: d.time, cost: d.cost })),
        });
      } catch (e) {
        // fallback mock if needed
      }
    };
    load();
  }, []);

  return (
    <Layout>
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-6">Dashboard</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <MetricCard 
            title="Total Requests"
            value={metrics.totalRequests.toLocaleString()}
            icon={Activity}
            change="+12.3%"
            changeType="positive"
          />
          <MetricCard 
            title="Average Latency"
            value={`${metrics.avgLatency}ms`}
            icon={Clock}
            change="-5.2%"
            changeType="positive"
          />
          <MetricCard 
            title="Cost Savings"
            value={`$${metrics.costSavings.toFixed(2)}`}
            icon={DollarSign}
            change="+8.1%"
            changeType="positive"
          />
          <MetricCard 
            title="Active Providers"
            value={metrics.activeProviders}
            icon={Zap}
            change="0"
            changeType="neutral"
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-gray-800 rounded-lg p-6">
            <h2 className="text-xl font-bold mb-4">Request Volume</h2>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={metrics.requests}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis dataKey="time" stroke="#9CA3AF" />
                  <YAxis stroke="#9CA3AF" />
                  <Tooltip 
                    contentStyle={{ 
                      backgroundColor: '#1F2937',
                      border: 'none',
                      borderRadius: '0.5rem'
                    }} 
                  />
                  <Legend />
                  <Line 
                    type="monotone" 
                    dataKey="requests" 
                    stroke="#3B82F6" 
                    strokeWidth={2}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="bg-gray-800 rounded-lg p-6">
            <h2 className="text-xl font-bold mb-4">Daily Costs</h2>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={metrics.costs}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis dataKey="day" stroke="#9CA3AF" />
                  <YAxis stroke="#9CA3AF" />
                  <Tooltip 
                    contentStyle={{ 
                      backgroundColor: '#1F2937',
                      border: 'none',
                      borderRadius: '0.5rem'
                    }} 
                  />
                  <Legend />
                  <Bar dataKey="cost" fill="#3B82F6" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}
