import React from 'react'
import { motion } from 'framer-motion'
interface Badge {
  name: string
  icon: string
  description: string
}
interface BadgeDisplayProps {
  badges: Badge[]
}
const BadgeDisplay: React.FC<BadgeDisplayProps> = ({ badges }) => {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h3 className="text-lg font-medium text-gray-900 mb-4">
        Your Achievement Badges
      </h3>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {badges.map((badge, index) => (
          <motion.div
            key={badge.name}
            className="flex flex-col items-center p-4 border border-gray-100 rounded-lg hover:border-indigo-200 hover:bg-indigo-50 transition-colors"
            initial={{
              opacity: 0,
              scale: 0.8,
            }}
            animate={{
              opacity: 1,
              scale: 1,
            }}
            transition={{
              delay: index * 0.1,
              duration: 0.3,
            }}
            whileHover={{
              y: -5,
            }}
          >
            <div className="text-4xl mb-2">{badge.icon}</div>
            <h4 className="font-medium text-gray-900 text-center">
              {badge.name}
            </h4>
            <p className="text-xs text-gray-500 text-center mt-1">
              {badge.description}
            </p>
          </motion.div>
        ))}
      </div>
    </div>
  )
}
export default BadgeDisplay
