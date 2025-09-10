import React, { useEffect, useState } from 'react'
import { motion } from 'framer-motion'
import { Wallet, Building, FileCheck, CreditCard, Heart } from 'lucide-react'
interface TransactionJourneyMapProps {
  transactionId: string
  donorName: string
  ngoName: string
  amount: number
  currentStage: number // 0-4 representing the different stages
}
const TransactionJourneyMap: React.FC<TransactionJourneyMapProps> = ({
  transactionId,
  donorName,
  ngoName,
  amount,
  currentStage,
}) => {
  const [animationStage, setAnimationStage] = useState(-1)
  useEffect(() => {
    // Animate through the stages
    const timer = setTimeout(() => {
      setAnimationStage(0)
      const interval = setInterval(() => {
        setAnimationStage((prev) => {
          if (prev >= currentStage) {
            clearInterval(interval)
            return currentStage
          }
          return prev + 1
        })
      }, 1000)
      return () => clearInterval(interval)
    }, 500)
    return () => clearTimeout(timer)
  }, [currentStage])
  const stages = [
    {
      icon: <Wallet size={24} />,
      name: "Donor's Wallet",
      description: 'Funds leaving your account',
    },
    {
      icon: <Building size={24} />,
      name: 'Smart Contract',
      description: 'Funds held in escrow',
    },
    {
      icon: <FileCheck size={24} />,
      name: 'Proof Verification',
      description: 'NGO submits documentation',
    },
    {
      icon: <CreditCard size={24} />,
      name: 'Funds Released',
      description: 'Transfer to NGO account',
    },
    {
      icon: <Heart size={24} />,
      name: 'Impact Delivered',
      description: 'Real-world change happens',
    },
  ]
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-lg font-medium text-gray-900">Money Journey Map</h3>
        <span className="text-sm text-gray-500">ID: {transactionId}</span>
      </div>
      <div className="relative">
        {/* Journey path */}
        <div className="absolute top-10 left-0 right-0 h-1 bg-gray-200"></div>
        <div className="flex justify-between relative">
          {stages.map((stage, index) => (
            <div
              key={index}
              className="relative flex flex-col items-center w-1/5"
            >
              {/* Animated path fill */}
              {index < stages.length - 1 && animationStage >= index && (
                <motion.div
                  className="absolute top-10 left-[50%] h-1 bg-indigo-600 z-10"
                  initial={{
                    width: 0,
                  }}
                  animate={{
                    width: '100%',
                  }}
                  transition={{
                    duration: 0.8,
                  }}
                  style={{
                    transformOrigin: 'left',
                  }}
                />
              )}
              {/* Icon circle */}
              <motion.div
                className={`w-20 h-20 rounded-full flex items-center justify-center z-20 ${animationStage >= index ? 'bg-indigo-100 text-indigo-600' : 'bg-gray-100 text-gray-400'}`}
                initial={{
                  scale: 0.8,
                  opacity: 0.5,
                }}
                animate={{
                  scale: animationStage >= index ? 1 : 0.8,
                  opacity: animationStage >= index ? 1 : 0.5,
                }}
                transition={{
                  duration: 0.3,
                }}
              >
                {stage.icon}
              </motion.div>
              {/* Label */}
              <motion.div
                className="mt-3 text-center"
                initial={{
                  opacity: 0,
                }}
                animate={{
                  opacity: animationStage >= index ? 1 : 0.5,
                }}
                transition={{
                  delay: 0.2,
                  duration: 0.3,
                }}
              >
                <div
                  className={`font-medium ${animationStage >= index ? 'text-indigo-600' : 'text-gray-400'}`}
                >
                  {stage.name}
                </div>
                <div className="text-xs text-gray-500 mt-1 max-w-[100px] mx-auto">
                  {stage.description}
                </div>
              </motion.div>
              {/* Active pulse animation */}
              {animationStage === index && (
                <motion.div
                  className="absolute top-10 w-20 h-20 rounded-full bg-indigo-500 opacity-20"
                  initial={{
                    scale: 0.8,
                  }}
                  animate={{
                    scale: 1.2,
                    opacity: 0,
                  }}
                  transition={{
                    repeat: Infinity,
                    duration: 1.5,
                  }}
                />
              )}
            </div>
          ))}
        </div>
      </div>
      <div className="flex justify-between mt-8">
        <div className="text-center">
          <div className="text-sm text-gray-500">From</div>
          <div className="font-medium text-gray-900">{donorName}</div>
        </div>
        <div className="text-center">
          <div className="text-sm text-gray-500">Amount</div>
          <div className="font-medium text-green-600">
            ₹{amount.toLocaleString()}
          </div>
        </div>
        <div className="text-center">
          <div className="text-sm text-gray-500">To</div>
          <div className="font-medium text-gray-900">{ngoName}</div>
        </div>
      </div>
    </div>
  )
}
export default TransactionJourneyMap
