import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { LayoutDashboard, Heart, Users, User, LogOut } from 'lucide-react'
import { motion } from 'framer-motion'
interface MainLayoutProps {
  children: React.ReactNode
}
const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const location = useLocation()
  const navItems = [
    {
      path: '/',
      icon: <LayoutDashboard size={20} />,
      label: 'Dashboard',
    },
    {
      path: '/donate',
      icon: <Heart size={20} />,
      label: 'Donate',
    },
    {
      path: '/ngo-selection',
      icon: <Users size={20} />,
      label: 'NGOs',
    },
    {
      path: '/profile',
      icon: <User size={20} />,
      label: 'Profile',
    },
  ]
  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <motion.div
        initial={{
          x: -100,
          opacity: 0,
        }}
        animate={{
          x: 0,
          opacity: 1,
        }}
        transition={{
          duration: 0.5,
        }}
        className="w-64 bg-white shadow-lg"
      >
        <div className="p-6">
          <h1 className="text-2xl font-bold text-indigo-600">DonorTrack</h1>
          <p className="text-sm text-gray-500">Transparent Giving</p>
        </div>
        <nav className="mt-6">
          <ul>
            {navItems.map((item) => (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={`flex items-center px-6 py-3 text-gray-600 hover:bg-indigo-50 hover:text-indigo-600 transition-colors ${location.pathname === item.path ? 'bg-indigo-50 text-indigo-600 border-r-4 border-indigo-600' : ''}`}
                >
                  {item.icon}
                  <span className="ml-3">{item.label}</span>
                </Link>
              </li>
            ))}
          </ul>
        </nav>
        <div className="absolute bottom-0 w-64 p-6">
          <button className="flex items-center text-gray-600 hover:text-indigo-600 transition-colors">
            <LogOut size={20} />
            <span className="ml-3">Logout</span>
          </button>
        </div>
      </motion.div>
      {/* Main content */}
      <div className="flex-1 overflow-auto">
        <div className="p-8">{children}</div>
      </div>
    </div>
  )
}
export default MainLayout
