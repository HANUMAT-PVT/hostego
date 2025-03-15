"use client"
import React, { useState, useEffect } from 'react'
import BottomNavigationBar from '../components/BottomNavigationBar'
import ProductCard from '../components/ProductCard'
import SearchComponent from '../components/SearchComponent'
import CartFloatingButton from "../components/Cart/CartFloatingButton"
import { ChevronDown, MapPin, ShoppingBagIcon, ShoppingBasket, User, UtensilsCrossed } from 'lucide-react'
import { useRouter } from 'next/navigation'
import axiosClient from '../utils/axiosClient'
import ProductCardSkeleton from '../components/ProductCardSkeleton'
import { useSelector } from 'react-redux'

const navItems = [
    { name: 'All', icon: ShoppingBasket, category: "" },
    { name: 'Food', icon: UtensilsCrossed, category: "food" },
    { name: 'Snacks', icon: ShoppingBagIcon, category: "snacks" },


]
const page = () => {
    const router = useRouter()
    const [activeIndex, setActiveIndex] = useState(0)
    const [products, setProducts] = useState([])
    const [isLoading, setIsLoading] = useState(true)
    const { cartData, useraddresses } = useSelector((state) => state.user)

    useEffect(() => {
        fetchProducts()
    }, [activeIndex])

    const fetchProducts = async () => {
        try {
            setIsLoading(true)
            const { data } = await axiosClient.get(`/api/products/all?page=1&limit=40&tags=${navItems[activeIndex]?.category}&admin=false`)
            setProducts(data)
        } catch (error) {
            console.error('Error fetching products:', error)
        } finally {
            setIsLoading(false)
        }
    }

    
    return (
        <div>
            <div className='gradient-background sticky top-0 z-10 '>
                <div className={` flex justify-between items-center   text-white p-4  items-left gap-1 `}>
                    <div>
                        <p className="text-xs font-bold">Hostego in</p>
                        <p className="text-2xl font-bold">few minutes </p>
                        <div onClick={()=>router.push("/address")} className='flex items-center gap-2'>
                            <p className="text-sm font-medium">{useraddresses[0]?.address_line_1 ||"Chandigarh University"}</p>
                            <ChevronDown className='w-4 h-4' />
                        </div>
                    </div>
                    <div onClick={() => router.push("/profile")} className='bg-white rounded-full p-1'>
                        <User color='black' className='rounded-full' size={22} />
                    </div>
                </div>
                <SearchComponent viewOnly={true} />
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
            <div className=' overflow-auto grid grid-cols-2 gap-2  justify-between mb-16 py-1 mb-16  '>
                {isLoading ? (
                    // Show 6 skeleton cards while loading
                    [...Array(6)].map((_, index) => (
                        <ProductCardSkeleton key={index} />
                    ))
                ) : (
                    products?.map((prd) => <ProductCard isAlreadyInCart={cartData?.cart_items?.some(item => item.product_id === prd?.product_id)} {...prd} key={prd?.product_id} />)
                )}

            </div>
            <BottomNavigationBar />
            {cartData?.cart_items?.length > 0 && <CartFloatingButton />}
        </div>
    )
}

export default page
