import React, { useEffect, useState } from 'react'
import { motion } from 'framer-motion'
import { Shield, AlertCircle, CheckCircle } from 'lucide-react'
interface TrustScoreCardProps {
  ngoName: string
  trustScore: number
  details?: {
    documentScore: number
    socialProofScore: number
    reviewScore: number
    financialScore: number
  }
}
const TrustScoreCard: React.FC<TrustScoreCardProps> = ({
  ngoName,
  trustScore,
  details = {
    documentScore: 90,
    socialProofScore: 85,
    reviewScore: 95,
    financialScore: 88,
  },
}) => {
  const [score, setScore] = useState(0)
  useEffect(() => {
    // Animate the score from 0 to the actual value
    const timer = setTimeout(() => {
      setScore(trustScore)
    }, 500)
    return () => clearTimeout(timer)
  }, [trustScore])
  // Determine the color based on the score
  const getScoreColor = (value: number) => {
    if (value >= 90) return 'text-green-500'
    if (value >= 70) return 'text-yellow-500'
    return 'text-red-500'
  }
  const getScoreBackground = (value: number) => {
    if (value >= 90) return 'bg-green-100'
    if (value >= 70) return 'bg-yellow-100'
    return 'bg-red-100'
  }
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-medium text-gray-900">{ngoName}</h3>
        <Shield className="text-indigo-600" size={20} />
      </div>
      <div className="flex items-center justify-center my-6">
        <motion.div
          className={`relative w-32 h-32 rounded-full flex items-center justify-center ${getScoreBackground(trustScore)}`}
          initial={{
            scale: 0.8,
            opacity: 0,
          }}
          animate={{
            scale: 1,
            opacity: 1,
          }}
          transition={{
            duration: 0.5,
          }}
        >
          <motion.div
            className={`text-4xl font-bold ${getScoreColor(trustScore)}`}
            initial={{
              opacity: 0,
            }}
            animate={{
              opacity: 1,
            }}
            transition={{
              delay: 0.3,
              duration: 0.5,
            }}
          >
            {score}
          </motion.div>
          <div className="absolute bottom-1 text-xs text-gray-500">
            Trust Score
          </div>
        </motion.div>
      </div>
      <div className="space-y-3">
        {Object.entries(details).map(([key, value], index) => (
          <motion.div
            key={key}
            className="flex items-center justify-between"
            initial={{
              x: -20,
              opacity: 0,
            }}
            animate={{
              x: 0,
              opacity: 1,
            }}
            transition={{
              delay: 0.2 * index,
              duration: 0.5,
            }}
          >
            <div className="flex items-center">
              {value >= 90 ? (
                <CheckCircle size={16} className="text-green-500 mr-2" />
              ) : value >= 70 ? (
                <Shield size={16} className="text-yellow-500 mr-2" />
              ) : (
                <AlertCircle size={16} className="text-red-500 mr-2" />
              )}
              <span className="text-sm text-gray-600">
                {key
                  .replace(/([A-Z])/g, ' $1')
                  .replace(/^./, (str) => str.toUpperCase())}
              </span>
            </div>
            <span className={`text-sm font-medium ${getScoreColor(value)}`}>
              {value}/100
            </span>
          </motion.div>
        ))}
      </div>
    </div>
  )
}
export default TrustScoreCard
