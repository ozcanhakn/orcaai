import { motion } from 'framer-motion'
import {
    Activity,
    ArrowRight,
    Brain,
    Clock,
    DollarSign,
    Shield,
    Sparkles,
    TrendingUp,
    Zap
} from 'lucide-react'
import Link from 'next/link'
import { useEffect, useState } from 'react'

export default function Home() {
  const [stats, setStats] = useState({
    totalRequests: 0,
    avgLatency: 0,
    costSavings: 0,
    uptime: 99.9
  })

  useEffect(() => {
    // Simulated real-time stats - replace with actual API call
    const interval = setInterval(() => {
      setStats({
        totalRequests: Math.floor(Math.random() * 10000) + 50000,
        avgLatency: Math.floor(Math.random() * 500) + 800,
        costSavings: (Math.random() * 30 + 20).toFixed(1),
        uptime: (99.5 + Math.random() * 0.5).toFixed(1)
      })
    }, 3000)

    return () => clearInterval(interval)
  }, [])

  const features = [
    {
      icon: Brain,
      title: "Smart AI Routing",
      description: "Automatically routes requests to the best AI provider based on cost, latency, and quality."
    },
    {
      icon: Zap,
      title: "Lightning Fast",
      description: "Advanced caching and optimization deliver responses in milliseconds."
    },
    {
      icon: TrendingUp,
      title: "Cost Optimization",
      description: "Save up to 40% on AI costs with intelligent provider selection."
    },
    {
      icon: Shield,
      title: "Enterprise Ready",
      description: "Built for scale with fallback systems and 99.9% uptime guarantee."
    }
  ]

  const providers = [
    { name: "OpenAI", status: "active", requests: "45K", latency: "850ms" },
    { name: "Claude", status: "active", requests: "38K", latency: "1.2s" },
    { name: "Gemini", status: "active", requests: "22K", latency: "950ms" }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
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
                            <Link
                href="/login"
                className="bg-blue-600 text-white px-6 py-3 rounded-lg inline-flex items-center gap-2 hover:bg-blue-700"
              >
                Login
                <ArrowRight className="w-4 h-4" />
              </Link>
              <Link
                href="/docs"
                className="text-gray-300 px-6 py-3 rounded-lg inline-flex items-center gap-2 hover:text-white"
              >
                View API Docs
                <ArrowRight className="w-4 h-4" />
              </Link>
              <button className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-4 py-2 rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all">
                Get Started
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-center"
        >
          <h1 className="text-5xl md:text-7xl font-bold text-white mb-6">
            AI <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-400">Orchestration</span>
            <br />Made Simple
          </h1>
          <p className="text-xl text-gray-300 mb-8 max-w-3xl mx-auto">
            Route AI requests intelligently across multiple providers. Optimize for cost, speed, and quality with automatic failover and caching.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link href="/dashboard">
              <button className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-4 rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all flex items-center justify-center space-x-2 text-lg font-semibold">
                <span>Open Dashboard</span>
                <ArrowRight className="w-5 h-5" />
              </button>
            </Link>
            <button className="border border-white/20 text-white px-8 py-4 rounded-xl hover:bg-white/10 transition-all text-lg font-semibold">
              View API Docs
            </button>
          </div>
        </motion.div>

        {/* Real-time Stats */}
        <motion.div 
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.3 }}
          className="mt-20 grid grid-cols-1 md:grid-cols-4 gap-6"
        >
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Total Requests</p>
                <p className="text-2xl font-bold text-white">{stats.totalRequests.toLocaleString()}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-400" />
            </div>
          </div>
          
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Avg Latency</p>
                <p className="text-2xl font-bold text-white">{stats.avgLatency}ms</p>
              </div>
              <Clock className="w-8 h-8 text-green-400" />
            </div>
          </div>
          
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Cost Savings</p>
                <p className="text-2xl font-bold text-white">{stats.costSavings}%</p>
              </div>
              <DollarSign className="w-8 h-8 text-purple-400" />
            </div>
          </div>
          
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-400 text-sm">Uptime</p>
                <p className="text-2xl font-bold text-white">{stats.uptime}%</p>
              </div>
              <Shield className="w-8 h-8 text-emerald-400" />
            </div>
          </div>
        </motion.div>

        {/* Features Grid */}
        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.6 }}
          className="mt-32"
        >
          <h2 className="text-4xl font-bold text-white text-center mb-16">
            Why Choose <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-400">OrcaAI</span>?
          </h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={feature.title}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.8 + index * 0.1 }}
                className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6 hover:bg-white/10 transition-all group"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                  <feature.icon className="w-6 h-6 text-white" />
                </div>
                <h3 className="text-xl font-semibold text-white mb-3">{feature.title}</h3>
                <p className="text-gray-400">{feature.description}</p>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Provider Status */}
        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 1.0 }}
          className="mt-32"
        >
          <h2 className="text-3xl font-bold text-white text-center mb-12">Provider Status</h2>
          
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {providers.map((provider, index) => (
                <div key={provider.name} className="flex items-center justify-between p-4 bg-white/5 rounded-xl">
                  <div className="flex items-center space-x-3">
                    <div className={`w-3 h-3 rounded-full ${provider.status === 'active' ? 'bg-green-400' : 'bg-red-400'}`}></div>
                    <span className="text-white font-semibold">{provider.name}</span>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-400">Requests</p>
                    <p className="text-white font-bold">{provider.requests}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-400">Avg Latency</p>
                    <p className="text-white font-bold">{provider.latency}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </motion.div>

        {/* CTA Section */}
        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 1.2 }}
          className="mt-32 text-center"
        >
          <div className="bg-gradient-to-r from-blue-600/20 to-purple-600/20 border border-blue-500/20 rounded-3xl p-12">
            <Sparkles className="w-16 h-16 text-blue-400 mx-auto mb-6" />
            <h2 className="text-4xl font-bold text-white mb-6">
              Ready to Optimize Your AI Stack?
            </h2>
            <p className="text-xl text-gray-300 mb-8 max-w-2xl mx-auto">
              Join thousands of developers who trust OrcaAI to handle their AI orchestration needs.
            </p>
            <Link href="/dashboard">
              <button className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-12 py-4 rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all text-lg font-semibold">
                Get Started Now
              </button>
            </Link>
          </div>
        </motion.div>
      </div>

      {/* Footer */}
      <footer className="border-t border-white/10 mt-32">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center space-x-2 mb-4">
                <Brain className="w-6 h-6 text-blue-400" />
                <span className="text-xl font-bold text-white">OrcaAI</span>
              </div>
              <p className="text-gray-400">
                Intelligent AI orchestration for modern applications.
              </p>
            </div>
            <div>
              <h3 className="text-white font-semibold mb-4">Product</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Dashboard</a></li>
                <li><a href="#" className="hover:text-white transition-colors">API</a></li>
                <li><a href="#" className="hover:text-white transition-colors">CLI</a></li>
              </ul>
            </div>
            <div>
              <h3 className="text-white font-semibold mb-4">Resources</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Documentation</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Examples</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Support</a></li>
              </ul>
            </div>
            <div>
              <h3 className="text-white font-semibold mb-4">Company</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">About</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Blog</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contact</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-white/10 mt-8 pt-8 text-center text-gray-400">
            <p>&copy; 2024 OrcaAI. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}