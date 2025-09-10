import React, { useEffect, useState } from 'react'
import { motion } from 'framer-motion'
import { Check, Clock, AlertTriangle } from 'lucide-react'
interface Stage {
  name: string
  completed: boolean
  date: string | null
}
interface DonationRoadmapProps {
  stages: Stage[]
  transactionId: string
}
const DonationRoadmap: React.FC<DonationRoadmapProps> = ({
  stages,
  transactionId,
}) => {
  const [activeStage, setActiveStage] = useState(0)
  // Find the current active stage
  useEffect(() => {
    const completedStages = stages.filter((stage) => stage.completed)
    setActiveStage(completedStages.length - 1)
  }, [stages])
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-medium text-gray-900">
          Transaction Roadmap
        </h3>
        <span className="text-sm text-gray-500">ID: {transactionId}</span>
      </div>
      <div className="relative">
        {/* Progress line */}
        <div className="absolute top-5 left-5 h-full w-0.5 bg-gray-200"></div>
        {stages.map((stage, index) => (
          <motion.div
            key={stage.name}
            initial={{
              opacity: 0,
              y: 20,
            }}
            animate={{
              opacity: 1,
              y: 0,
            }}
            transition={{
              delay: index * 0.2,
            }}
            className="flex items-start mb-6 relative"
          >
            {/* Status icon */}
            <div
              className={`flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center z-10 ${stage.completed ? 'bg-green-100 text-green-600' : index === activeStage + 1 ? 'bg-yellow-100 text-yellow-600' : 'bg-gray-100 text-gray-400'}`}
            >
              {stage.completed ? (
                <Check size={18} />
              ) : index === activeStage + 1 ? (
                <Clock size={18} />
              ) : (
                <div className="w-3 h-3 bg-gray-300 rounded-full"></div>
              )}
            </div>
            {/* Stage details */}
            <div className="ml-4">
              <h4
                className={`font-medium ${stage.completed ? 'text-green-600' : index === activeStage + 1 ? 'text-yellow-600' : 'text-gray-400'}`}
              >
                {stage.name}
              </h4>
              {stage.date && (
                <p className="text-sm text-gray-500">
                  {new Date(stage.date).toLocaleDateString()}
                </p>
              )}
              {index === activeStage + 1 && (
                <p className="text-xs text-yellow-600 flex items-center mt-1">
                  <Clock size={12} className="mr-1" />
                  In progress
                </p>
              )}
            </div>
            {/* Animated highlight for active stage */}
            {index === activeStage && (
              <motion.div
                initial={{
                  scale: 0.8,
                  opacity: 0,
                }}
                animate={{
                  scale: 1,
                  opacity: 1,
                }}
                transition={{
                  repeat: Infinity,
                  repeatType: 'reverse',
                  duration: 1.5,
                }}
                className="absolute left-5 top-5 w-10 h-10 -ml-5 -mt-5 rounded-full bg-green-500 opacity-20"
              />
            )}
          </motion.div>
        ))}
      </div>
    </div>
  )
}
export default DonationRoadmap
