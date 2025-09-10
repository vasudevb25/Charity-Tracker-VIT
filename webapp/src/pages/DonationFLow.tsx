// import React, { useState } from 'react'
// import { motion } from 'framer-motion'
// import { ArrowRight, Check, CreditCard, Wallet } from 'lucide-react'
// import TrustScoreCard from '../components/trustscore'
// // import { ngoData } from '../utils/mockData'
// const DonationFlow: React.FC = () => {
//   const [step, setStep] = useState(1)
//   const [selectedNgo, setSelectedNgo] = useState<(typeof ngoData)[0] | null>(
//     null,
//   )
//   const [amount, setAmount] = useState(1000)
//   const [paymentMethod, setPaymentMethod] = useState('upi')
//   const handleNgoSelect = (ngo: (typeof ngoData)[0]) => {
//     setSelectedNgo(ngo)
//     setStep(2)
//   }
//   const handlePaymentSubmit = () => {
//     setStep(3)
//     // In a real app, this would handle the payment processing
//   }
//   const handleComplete = () => {
//     setStep(4)
//     // In a real app, this would redirect to a confirmation page
//   }
//   return (
//     <div className="max-w-4xl mx-auto">
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
//         className="mb-8"
//       >
//         <h1 className="text-3xl font-bold text-gray-900">Make a Donation</h1>
//         <p className="text-gray-500 mt-1">Support causes that matter to you</p>
//       </motion.div>
//       {/* Progress Steps */}
//       <div className="flex justify-between mb-8 relative">
//         <div className="absolute top-4 left-0 right-0 h-1 bg-gray-200"></div>
//         <div
//           className="absolute top-4 left-0 h-1 bg-indigo-600"
//           style={{
//             width: `${(step - 1) * 33.3}%`,
//           }}
//         ></div>
//         {[1, 2, 3].map((s) => (
//           <div key={s} className="relative z-10 flex flex-col items-center">
//             <motion.div
//               className={`w-8 h-8 rounded-full flex items-center justify-center ${step > s ? 'bg-indigo-600 text-white' : step === s ? 'bg-indigo-600 text-white' : 'bg-gray-200 text-gray-500'}`}
//               animate={{
//                 scale: step === s ? [1, 1.1, 1] : 1,
//               }}
//               transition={{
//                 duration: 0.5,
//                 repeat: step === s ? Infinity : 0,
//                 repeatDelay: 2,
//               }}
//             >
//               {step > s ? <Check size={16} /> : s}
//             </motion.div>
//             <div className="text-sm font-medium mt-2">
//               {s === 1 ? 'Select NGO' : s === 2 ? 'Payment' : 'Confirm'}
//             </div>
//           </div>
//         ))}
//       </div>
//       {/* Step 1: Select NGO */}
//       {step === 1 && (
//         <motion.div
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           exit={{
//             opacity: 0,
//             y: -20,
//           }}
//           transition={{
//             duration: 0.5,
//           }}
//         >
//           <h2 className="text-xl font-bold text-gray-900 mb-6">
//             Select an NGO to support
//           </h2>
//           <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
//             {ngoData.map((ngo) => (
//               <motion.div
//                 key={ngo.id}
//                 className="border rounded-lg overflow-hidden hover:shadow-md cursor-pointer transition-shadow"
//                 whileHover={{
//                   y: -5,
//                 }}
//                 onClick={() => handleNgoSelect(ngo)}
//               >
//                 <div className="h-40 overflow-hidden">
//                   <img
//                     src={ngo.image}
//                     alt={ngo.name}
//                     className="w-full h-full object-cover"
//                   />
//                 </div>
//                 <div className="p-4">
//                   <div className="flex justify-between items-start">
//                     <div>
//                       <h3 className="font-bold text-gray-900">{ngo.name}</h3>
//                       <p className="text-sm text-gray-500">{ngo.category}</p>
//                     </div>
//                     <div className="bg-green-100 text-green-600 text-sm font-bold px-2 py-1 rounded">
//                       {ngo.trustScore}/100
//                     </div>
//                   </div>
//                   <p className="text-sm text-gray-600 mt-2">
//                     {ngo.description}
//                   </p>
//                   <div className="flex justify-between items-center mt-4">
//                     <span className="text-xs text-gray-500">
//                       {ngo.location}
//                     </span>
//                     <button className="text-indigo-600 text-sm font-medium flex items-center">
//                       View Details
//                       <ArrowRight size={14} className="ml-1" />
//                     </button>
//                   </div>
//                 </div>
//               </motion.div>
//             ))}
//           </div>
//         </motion.div>
//       )}
//       {/* Step 2: Payment */}
//       {step === 2 && selectedNgo && (
//         <motion.div
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           exit={{
//             opacity: 0,
//             y: -20,
//           }}
//           transition={{
//             duration: 0.5,
//           }}
//           className="grid grid-cols-1 md:grid-cols-2 gap-8"
//         >
//           <div>
//             <h2 className="text-xl font-bold text-gray-900 mb-6">
//               Donation Details
//             </h2>
//             <div className="bg-white p-6 rounded-lg shadow-md mb-6">
//               <div className="flex items-center justify-between mb-6">
//                 <div>
//                   <h3 className="font-bold text-gray-900">
//                     {selectedNgo.name}
//                   </h3>
//                   <p className="text-sm text-gray-500">
//                     {selectedNgo.category}
//                   </p>
//                 </div>
//                 <img
//                   src={selectedNgo.image}
//                   alt={selectedNgo.name}
//                   className="w-16 h-16 object-cover rounded"
//                 />
//               </div>
//               <div className="mb-6">
//                 <label className="block text-sm font-medium text-gray-700 mb-2">
//                   Donation Amount (₹)
//                 </label>
//                 <div className="flex items-center">
//                   <button
//                     className="w-10 h-10 rounded-l bg-gray-100 flex items-center justify-center"
//                     onClick={() => setAmount(Math.max(100, amount - 100))}
//                   >
//                     -
//                   </button>
//                   <input
//                     type="number"
//                     value={amount}
//                     onChange={(e) => setAmount(Number(e.target.value))}
//                     className="w-full h-10 border-y text-center outline-none"
//                   />
//                   <button
//                     className="w-10 h-10 rounded-r bg-gray-100 flex items-center justify-center"
//                     onClick={() => setAmount(amount + 100)}
//                   >
//                     +
//                   </button>
//                 </div>
//                 <div className="flex justify-between mt-2">
//                   <button
//                     className="bg-gray-100 px-3 py-1 rounded text-sm"
//                     onClick={() => setAmount(500)}
//                   >
//                     ₹500
//                   </button>
//                   <button
//                     className="bg-gray-100 px-3 py-1 rounded text-sm"
//                     onClick={() => setAmount(1000)}
//                   >
//                     ₹1000
//                   </button>
//                   <button
//                     className="bg-gray-100 px-3 py-1 rounded text-sm"
//                     onClick={() => setAmount(2000)}
//                   >
//                     ₹2000
//                   </button>
//                   <button
//                     className="bg-gray-100 px-3 py-1 rounded text-sm"
//                     onClick={() => setAmount(5000)}
//                   >
//                     ₹5000
//                   </button>
//                 </div>
//               </div>
//               <div>
//                 <label className="block text-sm font-medium text-gray-700 mb-2">
//                   Payment Method
//                 </label>
//                 <div className="grid grid-cols-3 gap-3">
//                   <button
//                     className={`p-3 border rounded-lg flex flex-col items-center justify-center ${paymentMethod === 'upi' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200'}`}
//                     onClick={() => setPaymentMethod('upi')}
//                   >
//                     <Wallet
//                       size={24}
//                       className={
//                         paymentMethod === 'upi'
//                           ? 'text-indigo-600'
//                           : 'text-gray-400'
//                       }
//                     />
//                     <span className="text-sm mt-1">UPI</span>
//                   </button>
//                   <button
//                     className={`p-3 border rounded-lg flex flex-col items-center justify-center ${paymentMethod === 'card' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200'}`}
//                     onClick={() => setPaymentMethod('card')}
//                   >
//                     <CreditCard
//                       size={24}
//                       className={
//                         paymentMethod === 'card'
//                           ? 'text-indigo-600'
//                           : 'text-gray-400'
//                       }
//                     />
//                     <span className="text-sm mt-1">Card</span>
//                   </button>
//                   <button
//                     className={`p-3 border rounded-lg flex flex-col items-center justify-center ${paymentMethod === 'crypto' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200'}`}
//                     onClick={() => setPaymentMethod('crypto')}
//                   >
//                     <svg
//                       viewBox="0 0 24 24"
//                       width="24"
//                       height="24"
//                       className={
//                         paymentMethod === 'crypto'
//                           ? 'text-indigo-600'
//                           : 'text-gray-400'
//                       }
//                     >
//                       <path
//                         fill="currentColor"
//                         d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13v2h2v-2h1v2h2v2h-2v2h-2v-2h-2v-2h2V7h-1z"
//                       />
//                     </svg>
//                     <span className="text-sm mt-1">Crypto</span>
//                   </button>
//                 </div>
//               </div>
//               <button
//                 className="w-full mt-6 bg-indigo-600 text-white py-3 rounded-lg font-medium hover:bg-indigo-700 transition-colors"
//                 onClick={handlePaymentSubmit}
//               >
//                 Proceed to Pay ₹{amount}
//               </button>
//             </div>
//           </div>
//           <div>
//             <h2 className="text-xl font-bold text-gray-900 mb-6">
//               Trust Score
//             </h2>
//             <TrustScoreCard
//               ngoName={selectedNgo.name}
//               trustScore={selectedNgo.trustScore}
//             />
//           </div>
//         </motion.div>
//       )}
//       {/* Step 3: Confirmation */}
//       {step === 3 && selectedNgo && (
//         <motion.div
//           initial={{
//             opacity: 0,
//             y: 20,
//           }}
//           animate={{
//             opacity: 1,
//             y: 0,
//           }}
//           exit={{
//             opacity: 0,
//             y: -20,
//           }}
//           transition={{
//             duration: 0.5,
//           }}
//           className="text-center"
//         >
//           <div className="mb-6">
//             <motion.div
//               className="w-24 h-24 bg-green-100 rounded-full flex items-center justify-center mx-auto"
//               initial={{
//                 scale: 0,
//               }}
//               animate={{
//                 scale: 1,
//               }}
//               transition={{
//                 type: 'spring',
//                 stiffness: 200,
//                 delay: 0.2,
//               }}
//             >
//               <Check size={48} className="text-green-600" />
//             </motion.div>
//           </div>
//           <h2 className="text-2xl font-bold text-gray-900 mb-2">
//             Payment Successful!
//           </h2>
//           <p className="text-gray-600 mb-8">
//             Your donation of ₹{amount} to {selectedNgo.name} has been processed.
//           </p>
//           <div className="bg-white p-6 rounded-lg shadow-md mb-8 max-w-md mx-auto">
//             <h3 className="font-bold text-gray-900 mb-4">What happens next?</h3>
//             <ul className="space-y-3 text-left">
//               <li className="flex items-start">
//                 <div className="bg-indigo-100 rounded-full p-1 mr-3 mt-0.5">
//                   <Check size={12} className="text-indigo-600" />
//                 </div>
//                 <span className="text-gray-600 text-sm">
//                   Your donation is now locked in a smart contract escrow
//                 </span>
//               </li>
//               <li className="flex items-start">
//                 <div className="bg-gray-100 rounded-full p-1 mr-3 mt-0.5">
//                   <Check size={12} className="text-gray-400" />
//                 </div>
//                 <span className="text-gray-600 text-sm">
//                   {selectedNgo.name} will submit proof of impact
//                 </span>
//               </li>
//               <li className="flex items-start">
//                 <div className="bg-gray-100 rounded-full p-1 mr-3 mt-0.5">
//                   <Check size={12} className="text-gray-400" />
//                 </div>
//                 <span className="text-gray-600 text-sm">
//                   Funds will be released upon verification
//                 </span>
//               </li>
//               <li className="flex items-start">
//                 <div className="bg-gray-100 rounded-full p-1 mr-3 mt-0.5">
//                   <Check size={12} className="text-gray-400" />
//                 </div>
//                 <span className="text-gray-600 text-sm">
//                   You'll receive updates on the impact of your donation
//                 </span>
//               </li>
//             </ul>
//           </div>
//           <button
//             className="bg-indigo-600 text-white py-3 px-8 rounded-lg font-medium hover:bg-indigo-700 transition-colors"
//             onClick={handleComplete}
//           >
//             View Your Dashboard
//           </button>
//         </motion.div>
//       )}
//       {/* Step 4: Success Animation */}
//       {step === 4 && (
//         <motion.div
//           initial={{
//             opacity: 0,
//           }}
//           animate={{
//             opacity: 1,
//           }}
//           transition={{
//             duration: 0.5,
//           }}
//           className="text-center"
//         >
//           <motion.div
//             initial={{
//               scale: 0.8,
//               opacity: 0,
//             }}
//             animate={{
//               scale: 1,
//               opacity: 1,
//             }}
//             transition={{
//               delay: 0.3,
//               duration: 0.5,
//             }}
//             className="mb-8"
//           >
//             <div className="w-32 h-32 bg-green-100 rounded-full flex items-center justify-center mx-auto">
//               <Check size={64} className="text-green-600" />
//             </div>
//           </motion.div>
//           <h2 className="text-3xl font-bold text-gray-900 mb-4">Thank You!</h2>
//           <p className="text-gray-600 mb-8 text-lg">
//             Your donation will make a real difference.
//           </p>
//           <div className="flex justify-center space-x-4">
//             <button className="bg-indigo-600 text-white py-3 px-6 rounded-lg font-medium hover:bg-indigo-700 transition-colors">
//               Track Your Impact
//             </button>
//             <button className="border border-indigo-600 text-indigo-600 py-3 px-6 rounded-lg font-medium hover:bg-indigo-50 transition-colors">
//               Share Your Support
//             </button>
//           </div>
//         </motion.div>
//       )}
//     </div>
//   )
// }
// export default DonationFlow
