'use client'
import { CircleUserRound, House, Package, Search } from 'lucide-react'
import React, { useState } from 'react'

const BottomNavigationBar = () => {
  const [activeIndex, setActiveIndex] = useState(0)

  const navItems = [
    { name: 'Home', icon: House },
    { name: 'Orders', icon: Package }, // Replace with Settings icon
    { name: 'Search', icon: Search }, // Replace with Search icon
    { name: 'Profile', icon: CircleUserRound }, // Replace with Profile icon
  ]

  return (
    <div className="fixed bottom-0 w-full bg-white z-2 flex justify-between items-center px-4 shadow-2xl  ">
      {navItems.map((item, index) => {
        const Icon = item.icon
        return (
          <div
            key={index}
            className={`w-[60px] bottom-nav-item cursor-pointer gap-1 flex flex-col items-center py-2 text-gray-500 border-t-2  rounded-t-xs  ${activeIndex === index ? 'text-[var(--primary-color)] border-t-2 border-[var(--primary-color)] rounded-t-xs' : ''}`}
            onClick={() => setActiveIndex(index)}
          >
            <Icon className={`${activeIndex === index ? 'text-[var(--primary-color)]' : ''}`} size={22} />
            <p className={`text-xs  ${activeIndex === index ? 'text-xs text-[var(--primary-color)] font-semibold' : ''}`}>{item.name}</p>
          </div>
        )
      })}
    </div>
  )
}

export default BottomNavigationBar
