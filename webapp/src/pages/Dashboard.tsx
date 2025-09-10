// import React, { useState } from 'react'
// import { motion } from 'framer-motion'
// import DashboardCharts from '../components/Dashboardchart'
// //import Leaderboard from '../components/Leaderboard'
// import TransactionJourneyMap from '../components/Transaction'
// // import {
// //   donationData,
// //   categoryDistribution,
// //   donorLeaderboard,
// //   ngoLeaderboard,
// // } from '../utils/mockData'
// import { ArrowUpRight, BarChart3, TrendingUp, Users } from 'lucide-react'
// const Dashboard: React.FC = () => {
//   const [selectedTransaction, setSelectedTransaction] = useState({
//     transactionId: 'TX123457',
//     donorName: 'Aditya Sharma',
//     ngoName: 'Clean Water Initiative',
//     amount: 3000,
//     currentStage: 2, // 0-4 for the different stages
//   })
//   // Stats data
//   const stats = [
//     {
//       name: 'Total Donations',
//       value: '₹2.4M',
//       icon: <BarChart3 size={20} />,
//       change: '+12%',
//     },
//     {
//       name: 'Active NGOs',
//       value: '126',
//       icon: <Users size={20} />,
//       change: '+8%',
//     },
//     {
//       name: 'Avg. Trust Score',
//       value: '89',
//       icon: <TrendingUp size={20} />,
//       change: '+5%',
//     },
//     {
//       name: 'Success Rate',
//       value: '98%',
//       icon: <ArrowUpRight size={20} />,
//       change: '+2%',
//     },
//   ]
//   return (
//     <div className="space-y-8">
//       <motion.div
//         initial={{
//           opacity: 0,
//           y: -20,
//         }}
//         animate={{
//           opacity: 1,
//           y: 0,
//         }}
//         transition={{
//           duration: 0.5,
//         }}
//       >
//         <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
//         <p className="text-gray-500 mt-1">
//           Track your donations and their impact
//         </p>
//       </motion.div>
//       {/* Stats */}
//       <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
//         {stats.map((stat, index) => (
//           <motion.div
//             key={stat.name}
//             className="bg-white p-6 rounded-lg shadow-md"
//             initial={{
//               opacity: 0,
//               y: 20,
//             }}
//             animate={{
//               opacity: 1,
//               y: 0,
//             }}
//             transition={{
//               delay: index * 0.1,
//               duration: 0.5,
//             }}
//           >
//             <div className="flex justify-between items-start">
//               <div>
//                 <p className="text-sm text-gray-500">{stat.name}</p>
//                 <h3 className="text-2xl font-bold text-gray-900 mt-1">
//                   {stat.value}
//                 </h3>
//               </div>
//               <div
//                 className={`p-2 rounded-lg ${stat.name === 'Total Donations' ? 'bg-blue-100 text-blue-600' : stat.name === 'Active NGOs' ? 'bg-purple-100 text-purple-600' : stat.name === 'Avg. Trust Score' ? 'bg-green-100 text-green-600' : 'bg-amber-100 text-amber-600'}`}
//               >
//                 {stat.icon}
//               </div>
//             </div>
//             <div className="flex items-center mt-4">
//               <span className="text-green-500 text-sm font-medium">
//                 {stat.change}
//               </span>
//               <span className="text-gray-500 text-sm ml-2">vs last month</span>
//             </div>
//           </motion.div>
//         ))}
//       </div>
//       {/* Charts */}
//       <DashboardCharts
//         monthlyDonations={donationData}
//         categoryDistribution={categoryDistribution}
//       />
//       {/* Transaction Journey Map */}
//       <motion.div
//         initial={{
//           opacity: 0,
//           y: 20,
//         }}
//         animate={{
//           opacity: 1,
//           y: 0,
//         }}
//         transition={{
//           delay: 0.3,
//           duration: 0.5,
//         }}
//       >
//         <TransactionJourneyMap
//           transactionId={selectedTransaction.transactionId}
//           donorName={selectedTransaction.donorName}
//           ngoName={selectedTransaction.ngoName}
//           amount={selectedTransaction.amount}
//           currentStage={selectedTransaction.currentStage}
//         />
//       </motion.div>
//       {/* Leaderboard */}
//       <motion.div
//         initial={{
//           opacity: 0,
//           y: 20,
//         }}
//         animate={{
//           opacity: 1,
//           y: 0,
//         }}
//         transition={{
//           delay: 0.4,
//           duration: 0.5,
//         }}
//       >
//         <Leaderboard donorData={donorLeaderboard} ngoData={ngoLeaderboard} />
//       </motion.div>
//     </div>
//   )
// }
// export default Dashboard
