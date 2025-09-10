import React from 'react'
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  CartesianGrid,
} from 'recharts'
import { motion } from 'framer-motion'
interface DashboardChartsProps {
  monthlyDonations: Array<{
    month: string
    amount: number
  }>
  categoryDistribution: Array<{
    name: string
    value: number
  }>
}
const COLORS = [
  '#6366F1',
  '#8B5CF6',
  '#EC4899',
  '#F43F5E',
  '#F97316',
  '#10B981',
  '#3B82F6',
]
const DashboardCharts: React.FC<DashboardChartsProps> = ({
  monthlyDonations,
  categoryDistribution,
}) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <motion.div
        className="bg-white p-6 rounded-lg shadow-md"
        initial={{
          opacity: 0,
          y: 20,
        }}
        animate={{
          opacity: 1,
          y: 0,
        }}
        transition={{
          duration: 0.5,
        }}
      >
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          Monthly Donations
        </h3>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={monthlyDonations}>
            <CartesianGrid strokeDasharray="3 3" vertical={false} />
            <XAxis dataKey="month" />
            <YAxis />
            <Tooltip
              formatter={(value) => [`₹${value}`, 'Amount']}
              contentStyle={{
                borderRadius: '8px',
              }}
            />
            <Bar
              dataKey="amount"
              fill="#6366F1"
              radius={[4, 4, 0, 0]}
              animationDuration={2000}
            />
          </BarChart>
        </ResponsiveContainer>
      </motion.div>
      <motion.div
        className="bg-white p-6 rounded-lg shadow-md"
        initial={{
          opacity: 0,
          y: 20,
        }}
        animate={{
          opacity: 1,
          y: 0,
        }}
        transition={{
          duration: 0.5,
          delay: 0.2,
        }}
      >
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          Donation Categories
        </h3>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={categoryDistribution}
              cx="50%"
              cy="50%"
              labelLine={false}
              outerRadius={100}
              fill="#8884d8"
              dataKey="value"
              animationDuration={2000}
              label={(entry: { name: string; percent: number }) =>
                `${entry.name}: ${(entry.percent * 100).toFixed(0)}%`
              }
            >
              {categoryDistribution.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={COLORS[index % COLORS.length]}
                />
              ))}
            </Pie>
            <Tooltip formatter={(value) => [`${value}%`, 'Percentage']} />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </motion.div>
    </div>
  )
}
export default DashboardCharts
