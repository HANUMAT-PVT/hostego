"use client"
import React, { useState } from 'react'
import BottomNavigationBar from '../components/BottomNavigationBar'
import ProductCard from '../components/ProductCard'
import SearchComponent from '../components/SearchComponent'
import { CircleUserRound, House, Package, Search, ShoppingBag, ShoppingBagIcon, ShoppingBasket, User, UtensilsCrossed } from 'lucide-react'
import { useRouter } from 'next/navigation'

const navItems = [
    { name: 'All', icon: ShoppingBasket },
    { name: 'Food', icon: UtensilsCrossed }, // Replace with Settings icon
    { name: 'Snacks', icon: ShoppingBagIcon }, // Replace with Search icon
    // { name: 'Profile', icon: CircleUserRound },

]
const page = () => {
    const router=useRouter()
    const [activeIndex, setActiveIndex] = useState(0)

    return (
        <div>
            <div className='gradient-background sticky top-0 z-10 '>
                <div className={` flex justify-between items-center   text-white p-4  items-left gap-1 `}>
                    <div>
                        <p className="text-xs font-bold">Hostego in</p>
                        <p className="text-2xl font-bold">15 minutes </p>
                        <p className="text-sm font-medium">Zakir A,Changiarh University</p>
                    </div>
                    <div onClick={()=>router.push("/profile")} className='bg-white rounded-full p-1'>
                        <User color='black' className='rounded-full' size={22} />
                    </div>
                </div>
               <SearchComponent viewOnly={true}/>
                <div className="  w-full bg-white sticky top-0  flex gap-4 items-center px-4 mt-3 overflow-auto justify-between   ">
                    {navItems.map((item, index) => {
                        const Icon = item.icon
                        return (
                            <div
                                key={index}
                                className={`min-w-[60px] bottom-nav-item  rounded-sm  cursor-pointer gap-1 flex flex-col items-center py-2 text-gray-500 border-b-2  rounded-t-xs  ${activeIndex === index ? 'text-[var(--primary-color)] border-b-2 border-[var(--primary-color)] rounded-t-xs' : ''}`}
                                onClick={() => setActiveIndex(index)}
                            >
                                <Icon className={`${activeIndex === index ? 'text-[var(--primary-color)]' : ''}`} size={22} />
                                <p className={`text-xs  ${activeIndex === index ? 'text-xs text-[var(--primary-color)] font-semibold' : ''}`}>{item.name}</p>
                            </div>
                        )
                    })}
                </div>

            </div>
            <div className=' overflow-auto flex gap-3 flex-wrap justify-between p-4 '>
                 <ProductCard  myKey={1}/>
                <ProductCard myKey={2} />
                <ProductCard myKey={3} />
                <ProductCard myKey={4} />
                <ProductCard myKey={5} />
                </div>
            <BottomNavigationBar />
        </div>
    )
}

export default page
