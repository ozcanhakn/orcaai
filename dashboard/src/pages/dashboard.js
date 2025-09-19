import {
    Activity,
    AlertTriangle,
    Brain,
    CheckCircle,
    Clock, DollarSign, Shield
} from 'lucide-react';
import { useEffect, useState } from 'react';
import {
    Bar,
    BarChart,
    CartesianGrid,
    Cell,
    Legend,
    Line,
    LineChart,
    Pie,
    PieChart,
    Tooltip,
    XAxis, YAxis
} from 'recharts';
import MetricCard from '../components/MetricCard';

export default function Dashboard() {
  const [metrics, setMetrics] = useState({
    requests: [],
    providers: [],
    costData: [],
    latencyData: [],
    cacheData: [],
    keyMetrics: {
      totalRequests: 0,
      avgLatency: 0,
      costSavings: 0,
      uptime: 0,
      activeProviders: 0,
      totalCost: 0,
      cacheHitRate: 0,
      errorRate: 0
    }
  });

  const [selectedProvider, setSelectedProvider] = useState('all');
  const [timeRange, setTimeRange] = useState('24h');
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulated real-time metrics - replace with actual API calls
    const interval = setInterval(() => {
      // Generate mock data
      const now = new Date()
      const requestsData = []
      const costData = []
      const latencyData = []
      
      for (let i = 0; i < 24; i++) {
        const time = new Date(now.getTime() - (23 - i) * 60 * 60 * 1000)
        requestsData.push({
          time: time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          requests: Math.floor(Math.random() * 1000) + 500,
          errors: Math.floor(Math.random() * 50)
        })
        
        costData.push({
          time: time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          cost: (Math.random() * 10).toFixed(2)
        })
        
        latencyData.push({
          time: time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          latency: Math.floor(Math.random() * 500) + 500
        })
      }
      
      const providers = [
        { name: 'OpenAI GPT-4', requests: 12450, errors: 12, latency: 850, cost: 245.67, status: 'active' },
        { name: 'Claude 3 Opus', requests: 9870, errors: 8, latency: 1200, cost: 189.45, status: 'active' },
        { name: 'Gemini Pro', requests: 7650, errors: 22, latency: 950, cost: 76.32, status: 'warning' },
        { name: 'OpenAI GPT-3.5', requests: 15430, errors: 45, latency: 650, cost: 89.21, status: 'active' }
      ]
      
      const cacheData = [
        { name: 'Cache Hits', value: 65 },
        { name: 'Cache Misses', value: 35 }
      ]
      
      const keyMetrics = {
        totalRequests: 45300,
        avgLatency: 842,
        costSavings: 2456.78,
        uptime: 99.98
      }
      
      setMetrics({
        requests: requestsData,
        providers,
        costData,
        latencyData,
        cacheData,
        keyMetrics
      })
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042']

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 text-white">
      {/* Navigation */}
      <nav className="border-b border-white/10 backdrop-blur-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Brain className="w-5 h-5 text-white" />
              </div>
              <span className="text-2xl font-bold text-white">OrcaAI</span>
            </div>
            <div className="flex items-center space-x-6">
              <Link href="/" className="text-gray-300 hover:text-white transition-colors">
                Home
              </Link>
              <Link href="/dashboard" className="text-white border-b-2 border-blue-500 pb-1">
                Dashboard
              </Link>
              <Link href="/docs" className="text-gray-300 hover:text-white transition-colors">
                Docs
              </Link>
              <button className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-4 py-2 rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all">
                API Keys
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Analytics Dashboard</h1>
          <div className="flex space-x-2">
            <button 
              className={`px-4 py-2 rounded-lg ${timeRange === '1h' ? 'bg-blue-600' : 'bg-white/10'}`}
              onClick={() => setTimeRange('1h')}
            >
              1H
            </button>
            <button 
              className={`px-4 py-2 rounded-lg ${timeRange === '24h' ? 'bg-blue-600' : 'bg-white/10'}`}
              onClick={() => setTimeRange('24h')}
            >
              24H
            </button>
            <button 
              className={`px-4 py-2 rounded-lg ${timeRange === '7d' ? 'bg-blue-600' : 'bg-white/10'}`}
              onClick={() => setTimeRange('7d')}
            >
              7D
            </button>
            <button 
              className={`px-4 py-2 rounded-lg ${timeRange === '30d' ? 'bg-blue-600' : 'bg-white/10'}`}
              onClick={() => setTimeRange('30d')}
            >
              30D
            </button>
          </div>
        </div>

        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <MetricCard
            title="Total Requests"
            value={metrics.keyMetrics.totalRequests.toLocaleString()}
            icon={Activity}
            change="12.5% from last week"
            changeType="positive"
          />
          
          <MetricCard
            title="Avg Latency"
            value={`${metrics.keyMetrics.avgLatency}ms`}
            icon={Clock}
            change="8.3% improvement"
            changeType="positive"
          />
          
          <MetricCard
            title="Cost Savings"
            value={`$${metrics.keyMetrics.costSavings.toLocaleString()}`}
            icon={DollarSign}
            change="24.7% savings"
            changeType="positive"
          />
          
          <MetricCard
            title="Uptime"
            value={`${metrics.keyMetrics.uptime}%`}
            icon={Shield}
            change="0.02% improvement"
            changeType="positive"
          />
        </div>

        {/* Charts Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Requests Chart */}
          <ChartContainer title="Requests Over Time">
            <LineChart data={metrics.requests}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis dataKey="time" stroke="#9CA3AF" />
              <YAxis stroke="#9CA3AF" />
              <Tooltip 
                contentStyle={{ backgroundColor: '#1F2937', borderColor: '#374151' }}
                itemStyle={{ color: 'white' }}
              />
              <Legend />
              <Line 
                type="monotone" 
                dataKey="requests" 
                stroke="#3B82F6" 
                strokeWidth={2}
                dot={{ strokeWidth: 2, r: 2 }}
                activeDot={{ r: 6 }}
                name="Requests"
              />
              <Line 
                type="monotone" 
                dataKey="errors" 
                stroke="#EF4444" 
                strokeWidth={2}
                dot={{ strokeWidth: 2, r: 2 }}
                activeDot={{ r: 6 }}
                name="Errors"
              />
            </LineChart>
          </ChartContainer>

          {/* Latency Chart */}
          <ChartContainer title="Latency Over Time">
            <LineChart data={metrics.latencyData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis dataKey="time" stroke="#9CA3AF" />
              <YAxis stroke="#9CA3AF" />
              <Tooltip 
                contentStyle={{ backgroundColor: '#1F2937', borderColor: '#374151' }}
                itemStyle={{ color: 'white' }}
              />
              <Legend />
              <Line 
                type="monotone" 
                dataKey="latency" 
                stroke="#10B981" 
                strokeWidth={2}
                dot={{ strokeWidth: 2, r: 2 }}
                activeDot={{ r: 6 }}
                name="Latency (ms)"
              />
            </LineChart>
          </ChartContainer>

          {/* Cost Chart */}
          <ChartContainer title="Cost Over Time">
            <BarChart data={metrics.costData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis dataKey="time" stroke="#9CA3AF" />
              <YAxis stroke="#9CA3AF" />
              <Tooltip 
                contentStyle={{ backgroundColor: '#1F2937', borderColor: '#374151' }}
                itemStyle={{ color: 'white' }}
                formatter={(value) => [`$${value}`, 'Cost']}
              />
              <Legend />
              <Bar 
                dataKey="cost" 
                fill="#8B5CF6" 
                name="Cost ($)"
              />
            </BarChart>
          </ChartContainer>

          {/* Cache Hit Rate */}
          <ChartContainer title="Cache Hit Rate">
            <PieChart>
              <Pie
                data={metrics.cacheData}
                cx="50%"
                cy="50%"
                labelLine={false}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
                label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
              >
                {metrics.cacheData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip 
                contentStyle={{ backgroundColor: '#1F2937', borderColor: '#374151' }}
                formatter={(value) => [`${value}%`, 'Rate']}
              />
              <Legend />
            </PieChart>
          </ChartContainer>
        </div>

        {/* Provider Performance */}
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Provider Performance</h2>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-700">
              <thead>
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Provider</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Requests</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Errors</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Avg Latency</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Cost</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-700">
                {metrics.providers.map((provider, index) => (
                  <tr key={index}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-white">{provider.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{provider.requests.toLocaleString()}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{provider.errors}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{provider.latency}ms</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">${provider.cost}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {provider.status === 'active' ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                          <CheckCircle className="w-3 h-3 mr-1" />
                          Active
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                          <AlertTriangle className="w-3 h-3 mr-1" />
                          Warning
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
          <h2 className="text-xl font-semibold mb-4">Recent Activity</h2>
          <div className="space-y-4">
            <div className="flex items-start">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center">
                  <Activity className="w-4 h-4 text-white" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm text-white">New request routed to OpenAI GPT-4</p>
                <p className="text-xs text-gray-400">2 minutes ago</p>
              </div>
            </div>
            
            <div className="flex items-start">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 rounded-full bg-green-500 flex items-center justify-center">
                  <CheckCircle className="w-4 h-4 text-white" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm text-white">Cache hit for summarization request</p>
                <p className="text-xs text-gray-400">5 minutes ago</p>
              </div>
            </div>
            
            <div className="flex items-start">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 rounded-full bg-yellow-500 flex items-center justify-center">
                  <AlertTriangle className="w-4 h-4 text-white" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm text-white">Fallback triggered for Gemini Pro</p>
                <p className="text-xs text-gray-400">12 minutes ago</p>
              </div>
            </div>
            
            <div className="flex items-start">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 rounded-full bg-purple-500 flex items-center justify-center">
                  <DollarSign className="w-4 h-4 text-white" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm text-white">Cost savings of $0.024 achieved</p>
                <p className="text-xs text-gray-400">18 minutes ago</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}