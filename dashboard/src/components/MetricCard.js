import React from 'react'
import { TrendingUp, TrendingDown } from 'lucide-react'

export default function MetricCard({ title, value, icon: Icon, change, changeType }) {
  return (
    <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-400 text-sm">{title}</p>
          <p className="text-2xl font-bold text-white">{value}</p>
        </div>
        <Icon className="w-8 h-8 text-blue-400" />
      </div>
      {change && (
        <div className={`mt-4 flex items-center ${changeType === 'positive' ? 'text-green-400' : 'text-red-400'}`}>
          {changeType === 'positive' ? (
            <TrendingUp className="w-4 h-4 mr-1" />
          ) : (
            <TrendingDown className="w-4 h-4 mr-1" />
          )}
          <span className="text-sm">{change}</span>
        </div>
      )}
    </div>
  )
}