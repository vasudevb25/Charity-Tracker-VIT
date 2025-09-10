// import React, { useState } from 'react'
// import { motion } from 'framer-motion'
// import { Search, Filter, MapPin, TrendingUp } from 'lucide-react'
// import TrustScoreCard from '../components/trustscore'
// // import { ngoData } from '../utils/MockData'
// const NGOSelection: React.FC = () => {
//   const [selectedNgo, setSelectedNgo] = useState<(typeof ngoData)[0] | null>(
//     null,
//   )
//   const [searchTerm, setSearchTerm] = useState('')
//   const [categoryFilter, setCategoryFilter] = useState('All')
//   const categories = [
//     'All',
//     'Education',
//     'Healthcare',
//     'Environment',
//     'Emergency',
//     'Social Welfare',
//   ]
//   const filteredNgos = ngoData.filter((ngo) => {
//     const matchesSearch =
//       ngo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
//       ngo.description.toLowerCase().includes(searchTerm.toLowerCase())
//     const matchesCategory =
//       categoryFilter === 'All' || ngo.category === categoryFilter
//     return matchesSearch && matchesCategory
//   })
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
//         <h1 className="text-3xl font-bold text-gray-900">NGO Directory</h1>
//         <p className="text-gray-500 mt-1">
//           Find and support verified organizations
//         </p>
//       </motion.div>
//       {/* Search and Filter */}
//       <div className="flex flex-col md:flex-row gap-4">
//         <div className="relative flex-1">
//           <Search
//             size={18}
//             className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
//           />
//           <input
//             type="text"
//             placeholder="Search NGOs by name or cause..."
//             value={searchTerm}
//             onChange={(e) => setSearchTerm(e.target.value)}
//             className="w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
//           />
//         </div>
//         <div className="flex items-center gap-2 overflow-x-auto pb-2">
//           <Filter size={18} className="text-gray-400 flex-shrink-0" />
//           {categories.map((category) => (
//             <button
//               key={category}
//               className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap ${categoryFilter === category ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
//               onClick={() => setCategoryFilter(category)}
//             >
//               {category}
//             </button>
//           ))}
//         </div>
//       </div>
//       {/* NGO List */}
//       <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
//         {filteredNgos.map((ngo) => (
//           <motion.div
//             key={ngo.id}
//             className={`bg-white rounded-lg overflow-hidden shadow-md cursor-pointer ${selectedNgo?.id === ngo.id ? 'ring-2 ring-indigo-500' : ''}`}
//             whileHover={{
//               y: -5,
//             }}
//             onClick={() => setSelectedNgo(ngo)}
//             layoutId={`ngo-card-${ngo.id}`}
//           >
//             <div className="h-40 overflow-hidden">
//               <img
//                 src={ngo.image}
//                 alt={ngo.name}
//                 className="w-full h-full object-cover"
//               />
//             </div>
//             <div className="p-4">
//               <div className="flex justify-between items-start">
//                 <h3 className="font-bold text-gray-900">{ngo.name}</h3>
//                 <div
//                   className={`px-2 py-1 rounded text-xs font-medium ${ngo.trustScore >= 90 ? 'bg-green-100 text-green-600' : ngo.trustScore >= 80 ? 'bg-yellow-100 text-yellow-600' : 'bg-red-100 text-red-600'}`}
//                 >
//                   {ngo.trustScore}/100
//                 </div>
//               </div>
//               <div className="flex items-center mt-2 text-sm text-gray-500">
//                 <MapPin size={14} className="mr-1" />
//                 {ngo.location}
//               </div>
//               <p className="text-sm text-gray-600 mt-2 line-clamp-2">
//                 {ngo.description}
//               </p>
//               <div className="flex justify-between items-center mt-4">
//                 <span className="inline-flex items-center bg-indigo-50 text-indigo-700 text-xs px-2 py-1 rounded-full">
//                   {ngo.category}
//                 </span>
//                 <button className="text-indigo-600 text-sm font-medium flex items-center">
//                   View Details
//                   <TrendingUp size={14} className="ml-1" />
//                 </button>
//               </div>
//             </div>
//           </motion.div>
//         ))}
//       </div>
//       {/* Selected NGO Details */}
//       {selectedNgo && (
//         <motion.div
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           transition={{
//             duration: 0.5,
//           }}
//           className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
//           onClick={() => setSelectedNgo(null)}
//         >
//           <motion.div
//             className="bg-white rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto"
//             onClick={(e) => e.stopPropagation()}
//             layoutId={`ngo-card-${selectedNgo.id}`}
//           >
//             <div className="h-64 relative">
//               <img
//                 src={selectedNgo.image}
//                 alt={selectedNgo.name}
//                 className="w-full h-full object-cover"
//               />
//               <button
//                 className="absolute top-4 right-4 w-8 h-8 bg-black bg-opacity-50 rounded-full flex items-center justify-center text-white"
//                 onClick={() => setSelectedNgo(null)}
//               >
//                 ✕
//               </button>
//             </div>
//             <div className="p-6">
//               <div className="flex justify-between items-start">
//                 <div>
//                   <h2 className="text-2xl font-bold text-gray-900">
//                     {selectedNgo.name}
//                   </h2>
//                   <div className="flex items-center mt-1 text-sm text-gray-500">
//                     <MapPin size={14} className="mr-1" />
//                     {selectedNgo.location}
//                   </div>
//                 </div>
//                 <span className="inline-flex items-center bg-indigo-50 text-indigo-700 text-sm px-3 py-1 rounded-full">
//                   {selectedNgo.category}
//                 </span>
//               </div>
//               <p className="text-gray-600 mt-4">
//                 {selectedNgo.description} Lorem ipsum dolor sit amet,
//                 consectetur adipiscing elit. Nullam at vestibulum eros, ac
//                 finibus justo. Cras interdum eros at ipsum sollicitudin, at
//                 consequat lectus venenatis. Mauris tempor metus vitae tortor
//                 commodo, nec tempus mauris aliquet.
//               </p>
//               <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
//                 <TrustScoreCard
//                   ngoName={selectedNgo.name}
//                   trustScore={selectedNgo.trustScore}
//                 />
//                 <div className="bg-white p-6 rounded-lg shadow-md">
//                   <h3 className="text-lg font-medium text-gray-900 mb-4">
//                     Impact Metrics
//                   </h3>
//                   <div className="space-y-4">
//                     <div>
//                       <div className="flex justify-between text-sm">
//                         <span className="text-gray-600">People Helped</span>
//                         <span className="font-medium">2,500+</span>
//                       </div>
//                       <div className="h-2 bg-gray-200 rounded-full mt-1">
//                         <div
//                           className="h-2 bg-indigo-600 rounded-full"
//                           style={{
//                             width: '80%',
//                           }}
//                         ></div>
//                       </div>
//                     </div>
//                     <div>
//                       <div className="flex justify-between text-sm">
//                         <span className="text-gray-600">
//                           Projects Completed
//                         </span>
//                         <span className="font-medium">36</span>
//                       </div>
//                       <div className="h-2 bg-gray-200 rounded-full mt-1">
//                         <div
//                           className="h-2 bg-indigo-600 rounded-full"
//                           style={{
//                             width: '65%',
//                           }}
//                         ></div>
//                       </div>
//                     </div>
//                     <div>
//                       <div className="flex justify-between text-sm">
//                         <span className="text-gray-600">Funds Utilized</span>
//                         <span className="font-medium">95%</span>
//                       </div>
//                       <div className="h-2 bg-gray-200 rounded-full mt-1">
//                         <div
//                           className="h-2 bg-indigo-600 rounded-full"
//                           style={{
//                             width: '95%',
//                           }}
//                         ></div>
//                       </div>
//                     </div>
//                   </div>
//                 </div>
//               </div>
//               <div className="mt-8 flex gap-4">
//                 <button className="flex-1 bg-indigo-600 text-white py-3 rounded-lg font-medium hover:bg-indigo-700 transition-colors">
//                   Donate Now
//                 </button>
//                 <button className="flex-1 border border-indigo-600 text-indigo-600 py-3 rounded-lg font-medium hover:bg-indigo-50 transition-colors">
//                   View Full Profile
//                 </button>
//               </div>
//             </div>
//           </motion.div>
//         </motion.div>
//       )}
//     </div>
//   )
// }
// export default NGOSelection
