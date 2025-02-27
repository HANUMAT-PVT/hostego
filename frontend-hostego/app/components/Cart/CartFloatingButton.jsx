
"use client"
import axiosClient from '@/app/utils/axiosClient'
import { ChevronRight } from 'lucide-react'
import { useRouter } from 'next/navigation'
import React, { useEffect, useState } from 'react'

const CartFloatingButton = () => {
    const router = useRouter()
    const [cartTotalItems, setCartTotalItems] = useState(0)
    const [cartData, setCartData] = useState(null)
    useEffect(() => {

        fetchCartItems()
    }, [])


    const fetchCartItems = async () => {
        const { data } = await axiosClient.get('/api/cart/')
      
        setCartData(data)
        let totalItems = 0
        data?.cart_items?.forEach(item => {
            totalItems += item?.quantity
        })
        setCartTotalItems(totalItems)
    }
    return (
        <div className='fixed bottom-[80px] w-fit m-auto translate-x-1/2 left-0 right-0 flex justify-center animate-slide-up'>
            <div
                onClick={() => router.push("/cart")}
                className='z-[1] min-w-[180px] px-2 py-2 flex items-center justify-between rounded-full bg-[var(--primary-color)]'
            >
                <div className='flex flex-col items-center'>
                    <img
                        className='w-[40px] rounded-full border-2'
                        src={cartData?.cart_items[0]?.product_item?.product_img_url}
                        alt={cartData?.cart_items[0]?.product_item?.product_name}
                    />
                </div>
                <div>
                    <p className='text-sm text-white font-semibold'>View cart</p>
                    <p className='text-xs text-white font-normal'>{cartTotalItems} ITEM</p>
                </div>
                <div className='p-2 bg-white w-[40px] h-[40px] flex items-center rounded-full justify-center'>
                    <ChevronRight color='var(--primary-color)' size={24} />
                </div>
            </div>
        </div>
    )
}

export default CartFloatingButton
