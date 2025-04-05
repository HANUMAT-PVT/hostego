"use client"
import React, { useState, useEffect, useRef } from 'react'
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
    const [page, setPage] = useState(1)
    const [hasMore, setHasMore] = useState(true)
    const { cartData, useraddresses } = useSelector((state) => state.user)
    const productsWrapperRef = useRef(null)

    useEffect(() => {
        // Reset state when category changes
        setProducts([])
        setPage(1)
        setHasMore(true)
        fetchProducts(1, true)
    }, [activeIndex])

    const fetchProducts = async (pageNum, isNewCategory = false) => {
        try {
            setIsLoading(true)
            const { data } = await axiosClient.get(
                `/api/products/all?page=${pageNum}&limit=15&tags=${navItems[activeIndex]?.category}&admin=false`
            )

            if (data.length < 15) {
                setHasMore(false)
            }

            setProducts(prev => isNewCategory ? data : [...prev, ...data])
        } catch (error) {
            console.error('Error fetching products:', error)
        } finally {
            setIsLoading(false)
        }
    }

    // Handle scroll event for infinite scrolling
    const handleScroll = () => {
        if (!productsWrapperRef.current || isLoading || !hasMore) return;

        const { scrollTop, scrollHeight, clientHeight } = productsWrapperRef.current;

        // If we're near the bottom (within 100px), load more
        if (scrollHeight - scrollTop - clientHeight < 300) {
            setPage(prevPage => {
                const nextPage = prevPage + 1;
                fetchProducts(nextPage);
                return nextPage;
            });
        }
    };

    useEffect(() => {
        const wrapper = productsWrapperRef.current;
        if (wrapper) {
            wrapper.addEventListener('scroll', handleScroll);
        }

        return () => {
            if (wrapper) {
                wrapper.removeEventListener('scroll', handleScroll);
            }
        };
    }, [hasMore, isLoading]);

    return (
        <div>
            <div className='gradient-background sticky top-0 z-10 '>
                <div className={` flex justify-between items-center   text-white p-4  items-left gap-1 `}>
                    <div>
                        <p className="text-xs font-bold">Hostego in</p>
                        <p className="text-2xl font-bold">few minutes </p>
                        <div onClick={() => router.push("/address")} className='flex items-center gap-2'>
                            <p className="text-sm font-medium">{useraddresses[0]?.address_line_1 || "Chandigarh University"}</p>
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
            <div
                ref={productsWrapperRef}
                className='h-[calc(100vh-200px)] overflow-y-auto'
            >
                <div className='grid grid-cols-2 gap-2 p-2'>
                    {products?.map((prd) => (
                        <ProductCard
                            isAlreadyInCart={cartData?.cart_items?.some(item => item.product_id === prd?.product_id)}
                            {...prd}
                            key={prd?.product_id}
                        />
                    ))}

                    {/* Loading state */}
                    {isLoading && (
                        [...Array(4)].map((_, index) => (
                            <ProductCardSkeleton key={`skeleton-${index}`} />
                        ))
                    )}

                    {/* Loading indicator at bottom */}
                    {hasMore && !isLoading && products.length > 0 && (
                        <div className="col-span-2 flex justify-center py-4">
                            <div className="w-6 h-6 border-2 border-[var(--primary-color)] border-t-transparent rounded-full animate-spin"></div>
                        </div>
                    )}

                    {/* No more products message */}
                    {!hasMore && products.length > 0 && (
                        <div className="col-span-2 text-center py-4 text-gray-500">
                            No more products to load
                        </div>
                    )}
                </div>
            </div>

            <BottomNavigationBar />
            {cartData?.cart_items?.length > 0 && <CartFloatingButton />}
        </div>
    )
}

export default page
