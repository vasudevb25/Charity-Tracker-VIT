// Mock data for demonstration purposes
export const ngoData = [
  { 
    id: 1, 
    name: "Child Education Foundation", 
    category: "Education", 
    trustScore: 92,
    location: "Bangalore",
    description: "Providing quality education to underprivileged children",
    image: "https://images.unsplash.com/photo-1488521787991-ed7bbaae773c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
  },
  { 
    id: 2, 
    name: "Clean Water Initiative", 
    category: "Environment", 
    trustScore: 87,
    location: "Mumbai",
    description: "Making clean water accessible to rural communities",
    image: "https://images.unsplash.com/photo-1581088382991-83640a2f2c40?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
  },
  { 
    id: 3, 
    name: "Disaster Relief Network", 
    category: "Emergency", 
    trustScore: 95,
    location: "Delhi",
    description: "Rapid response to natural disasters and emergencies",
    image: "https://images.unsplash.com/photo-1587653263995-422546a7a559?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
  },
  { 
    id: 4, 
    name: "Healthcare for All", 
    category: "Healthcare", 
    trustScore: 89,
    location: "Chennai",
    description: "Providing medical services to underserved communities",
    image: "https://images.unsplash.com/photo-1532938911079-1b06ac7ceec7?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1332&q=80",
  },
  { 
    id: 5, 
    name: "Women Empowerment Trust", 
    category: "Social Welfare", 
    trustScore: 91,
    location: "Kolkata",
    description: "Empowering women through education and entrepreneurship",
    image: "https://images.unsplash.com/photo-1537799943037-f5da89a65689?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
  }
];
export const donationData = [
  { month: 'Jan', amount: 12000 },
  { month: 'Feb', amount: 19000 },
  { month: 'Mar', amount: 15000 },
  { month: 'Apr', amount: 25000 },
  { month: 'May', amount: 22000 },
  { month: 'Jun', amount: 30000 },
  { month: 'Jul', amount: 28000 },
];
export const categoryDistribution = [
  { name: 'Education', value: 35 },
  { name: 'Healthcare', value: 25 },
  { name: 'Environment', value: 15 },
  { name: 'Emergency', value: 10 },
  { name: 'Social Welfare', value: 15 },
];
export const donorLeaderboard = [
  { id: 1, name: "Raj Sharma", amount: 50000, badges: ["Education Hero", "Platinum Donor"] },
  { id: 2, name: "Priya Patel", amount: 45000, badges: ["Healthcare Champion"] },
  { id: 3, name: "Amit Kumar", amount: 38000, badges: ["Environment Warrior"] },
  { id: 4, name: "Neha Singh", amount: 32000, badges: ["First-time Donor"] },
  { id: 5, name: "Vikram Reddy", amount: 28000, badges: ["Regular Contributor"] },
];
export const ngoLeaderboard = [
  { id: 1, name: "Child Education Foundation", transparencyScore: 98, impactScore: 95 },
  { id: 2, name: "Disaster Relief Network", transparencyScore: 96, impactScore: 97 },
  { id: 3, name: "Healthcare for All", transparencyScore: 94, impactScore: 92 },
  { id: 4, name: "Clean Water Initiative", transparencyScore: 92, impactScore: 90 },
  { id: 5, name: "Women Empowerment Trust", transparencyScore: 90, impactScore: 93 },
];
export const userTransactions = [
  { 
    id: "TX123456",
    ngo: "Child Education Foundation",
    amount: 5000,
    date: "2023-06-15",
    status: "completed",
    stages: [
      { name: "Donation Sent", completed: true, date: "2023-06-15" },
      { name: "Locked in Escrow", completed: true, date: "2023-06-15" },
      { name: "Proof Verified", completed: true, date: "2023-06-18" },
      { name: "Funds Released", completed: true, date: "2023-06-19" },
      { name: "Impact Delivered", completed: true, date: "2023-06-25" },
    ]
  },
  { 
    id: "TX123457",
    ngo: "Clean Water Initiative",
    amount: 3000,
    date: "2023-06-20",
    status: "in-progress",
    stages: [
      { name: "Donation Sent", completed: true, date: "2023-06-20" },
      { name: "Locked in Escrow", completed: true, date: "2023-06-20" },
      { name: "Proof Verified", completed: true, date: "2023-06-22" },
      { name: "Funds Released", completed: false, date: null },
      { name: "Impact Delivered", completed: false, date: null },
    ]
  },
  { 
    id: "TX123458",
    ngo: "Healthcare for All",
    amount: 2000,
    date: "2023-06-25",
    status: "early-stage",
    stages: [
      { name: "Donation Sent", completed: true, date: "2023-06-25" },
      { name: "Locked in Escrow", completed: true, date: "2023-06-25" },
      { name: "Proof Verified", completed: false, date: null },
      { name: "Funds Released", completed: false, date: null },
      { name: "Impact Delivered", completed: false, date: null },
    ]
  }
];
export const userProfile = {
  name: "Aditya Sharma",
  email: "aditya.sharma@example.com",
  totalDonated: 125000,
  donationCount: 15,
  joinedDate: "2022-11-10",
  badges: [
    { name: "Education Hero", icon: "🎓", description: "Donated over ₹50,000 to education causes" },
    { name: "Healthcare Champion", icon: "⚕️", description: "Supported 5+ healthcare initiatives" },
    { name: "Platinum Donor", icon: "🏆", description: "Among top 10% of donors" },
    { name: "First Responder", icon: "🚨", description: "Quick response to emergency appeals" },
  ],
  impactStats: {
    childrenEducated: 50,
    medicalCheckups: 125,
    treesPlanted: 30,
    mealsProvided: 500,
  }
};