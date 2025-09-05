import React from 'react'
import { ResponsiveContainer } from 'recharts'

export default function ChartContainer({ title, children, height = "h-80" }) {
  return (
    <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6">
      <h2 className="text-xl font-semibold mb-4">{title}</h2>
      <div className={height}>
        <ResponsiveContainer width="100%" height="100%">
          {children}
        </ResponsiveContainer>
      </div>
    </div>
  )
}