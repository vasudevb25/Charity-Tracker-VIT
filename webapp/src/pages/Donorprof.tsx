// import React, { Children } from 'react'
// import { motion } from 'framer-motion'
// import BadgeDisplay from '../components/BadgeDisplay'
// import DonationRoadmap from '../components/DonationRoadmap'
// // import { userProfile, userTransactions } from '../layouts/MockData.js'
// import { Clock, Heart, Award, TrendingUp, BookOpen, Users } from 'lucide-react'
// const DonorProfile: React.FC = () => {
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
//         className="flex flex-col md:flex-row md:items-center md:justify-between gap-6"
//       >
//         <div className="flex items-center gap-4">
//           <div className="w-20 h-20 bg-indigo-100 rounded-full flex items-center justify-center text-2xl font-bold text-indigo-600">
//             {userProfile.name.charAt(0)}
//           </div>
//           <div>
//             <h1 className="text-3xl font-bold text-gray-900">
//               {userProfile.name}
//             </h1>
//             <p className="text-gray-500">{userProfile.email}</p>
//           </div>
//         </div>
//         <div className="flex items-center gap-2">
//           <span className="text-gray-500 text-sm">
//             Member since {new Date(userProfile.joinedDate).toLocaleDateString()}
//           </span>
//           <button className="bg-indigo-600 text-white py-2 px-4 rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors">
//             Edit Profile
//           </button>
//         </div>
//       </motion.div>
//       {/* Stats Cards */}
//       <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
//         <motion.div
//           className="bg-white p-6 rounded-lg shadow-md"
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           transition={{
//             delay: 0.1,
//             duration: 0.5,
//           }}
//         >
//           <div className="flex justify-between items-start">
//             <div>
//               <p className="text-sm text-gray-500">Total Donated</p>
//               <h3 className="text-2xl font-bold text-gray-900 mt-1">
//                 ₹{userProfile.totalDonated.toLocaleString()}
//               </h3>
//             </div>
//             <div className="p-2 rounded-lg bg-indigo-100 text-indigo-600">
//               <Heart size={20} />
//             </div>
//           </div>
//           <div className="text-xs text-gray-500 mt-4">
//             {userProfile.donationCount} donations made
//           </div>
//         </motion.div>
//         <motion.div
//           className="bg-white p-6 rounded-lg shadow-md"
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           transition={{
//             delay: 0.2,
//             duration: 0.5,
//           }}
//         >
//           <div className="flex justify-between items-start">
//             <div>
//               <p className="text-sm text-gray-500">Active Donations</p>
//               <h3 className="text-2xl font-bold text-gray-900 mt-1">2</h3>
//             </div>
//             <div className="p-2 rounded-lg bg-green-100 text-green-600">
//               <Clock size={20} />
//             </div>
//           </div>
//           <div className="text-xs text-gray-500 mt-4">
//             Last activity 2 days ago
//           </div>
//         </motion.div>
//         <motion.div
//           className="bg-white p-6 rounded-lg shadow-md"
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           transition={{
//             delay: 0.3,
//             duration: 0.5,
//           }}
//         >
//           <div className="flex justify-between items-start">
//             <div>
//               <p className="text-sm text-gray-500">Badges Earned</p>
//               <h3 className="text-2xl font-bold text-gray-900 mt-1">
//                 {userProfile.badges.length}
//               </h3>
//             </div>
//             <div className="p-2 rounded-lg bg-yellow-100 text-yellow-600">
//               <Award size={20} />
//             </div>
//           </div>
//           <div className="text-xs text-gray-500 mt-4">Top 10% of donors</div>
//         </motion.div>
//         <motion.div
//           className="bg-white p-6 rounded-lg shadow-md"
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           transition={{
//             delay: 0.4,
//             duration: 0.5,
//           }}
//         >
//           <div className="flex justify-between items-start">
//             <div>
//               <p className="text-sm text-gray-500">Impact Score</p>
//               <h3 className="text-2xl font-bold text-gray-900 mt-1">87</h3>
//             </div>
//             <div className="p-2 rounded-lg bg-blue-100 text-blue-600">
//               <TrendingUp size={20} />
//             </div>
//           </div>
//           <div className="text-xs text-gray-500 mt-4">
//             +12 points this month
//           </div>
//         </motion.div>
//       </div>
//       {/* Impact Stats */}
//       <motion.div
//         className="bg-white p-6 rounded-lg shadow-md"
//         initial={{
//           opacity: 0,
//           y: 20,
//         }}
//         animate={{
//           opacity: 1,
//           y: 0,
//         }}
//         transition={{
//           delay: 0.5,
//           duration: 0.5,
//         }}
//       >
//         <h3 className="text-lg font-medium text-gray-900 mb-6">Your Impact</h3>
//         <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
//           <div className="text-center">
//             <div className="w-16 h-16 rounded-full bg-indigo-100 flex items-center justify-center mx-auto mb-3">
//               <BookOpen size={24} className="text-indigo-600" />
//             </div>
//             <div className="text-2xl font-bold text-gray-900">
//               {userProfile.impactStats.childrenEducated}
//             </div>
//             <div className="text-sm text-gray-500">Children Educated</div>
//           </div>
//           <div className="text-center">
//             <div className="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-3">
//               <svg
//                 viewBox="0 0 24 24"
//                 width="24"
//                 height="24"
//                 className="text-green-600"
//               >
//                 <path
//                   fill="currentColor"
//                   d="M17,8C8,10 5.9,16.17 3.82,21.34L5.71,22L6.66,19.7C7.14,19.87 7.64,20 8,20C19,20 22,3 22,3C21,5 14,5.25 9,6.25C4,7.25 2,11.5 2,13.5C2,15.5 3.75,17.25 3.75,17.25C7,8 17,8 17,8Z"
//                 />
//               </svg>
//             </div>
//             <div className="text-2xl font-bold text-gray-900">
//               {userProfile.impactStats.treesPlanted}
//             </div>
//             <div className="text-sm text-gray-500">Trees Planted</div>
//           </div>
//           <div className="text-center">
//             <div className="w-16 h-16 rounded-full bg-red-100 flex items-center justify-center mx-auto mb-3">
//               <svg
//                 viewBox="0 0 24 24"
//                 width="24"
//                 height="24"
//                 className="text-red-600"
//               >
//                 <path
//                   fill="currentColor"
//                   d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z"
//                 />
//               </svg>
//             </div>
//             <div className="text-2xl font-bold text-gray-900">
//               {userProfile.impactStats.medicalCheckups}
//             </div>
//             <div className="text-sm text-gray-500">Medical Checkups</div>
//           </div>
//           <div className="text-center">
//             <div className="w-16 h-16 rounded-full bg-yellow-100 flex items-center justify-center mx-auto mb-3">
//               <Users size={24} className="text-yellow-600" />
//             </div>
//             <div className="text-2xl font-bold text-gray-900">
//               {userProfile.impactStats.mealsProvided}
//             </div>
//             <div className="text-sm text-gray-500">Meals Provided</div>
//           </div>
//         </div>
//       </motion.div>
//       {/* Badges */}
//       <BadgeDisplay badges={userProfile.badges} />
//       {/* Recent Transactions */}
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
//           delay: 0.6,
//           duration: 0.5,
//         }}
//       >
//         <h3 className="text-lg font-medium text-gray-900 mb-6">
//           Recent Donations
//         </h3>
//         <div className="space-y-6">
//           {userTransactions.map((transaction) => (
//             <DonationRoadmap
//               key={transaction.id}
//               stages={transaction.stages}
//               transactionId={transaction.id}
//             />
//           ))}
//         </div>
//       </motion.div>
//     </div>
//   )
// }
// export default DonorProfile
