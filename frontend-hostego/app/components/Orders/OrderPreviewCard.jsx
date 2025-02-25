"use client"
import { Check } from 'lucide-react'
import { useRouter } from 'next/navigation'
import React from 'react'

const OrderPreviewCard = ({ myKey }) => {
    const router = useRouter()
    return (
        <div onClick={() => router.push("/orders/test")} key={myKey} className='rounded-lg bg-white p-2 rounded-md border  '>
            <div className='flex gap-4 items-center   py-2 p-2 border-b mb-2  '>
                <div className='bg-green-100 p-2 flex items-center justify-center w-[40px] h-[40px] rounded-md'>
                    <Check className='text-green-900 font-bold' size={20} />
                </div>
                <div className=''>
                    <p className='font-semibold text-lg'>Arrived in 23 Minutes</p>
                    <p className='text-sm text-gray-600'>₹112 {" "} 14 Feb, 8:12 pm</p>
                </div>
            </div>
            <div className='flex gap-4 p-2'>
                <div className='bg-[var(--bg-page-color)] p-1 rounded-md'>
                    <img className='w-[40px]'
                        src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                        alt={'Uncle chips'}
                    />
                </div>
                <div className='bg-[var(--bg-page-color)] p-1 rounded-md'>
                    <img className='w-[40px]'
                        src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                        alt={'Uncle chips'}
                    />
                </div>
                <div className='bg-[var(--bg-page-color)] p-1 rounded-md'>
                    <img className='w-[40px]'
                        src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                        alt={'Uncle chips'}
                    />
                </div>
            </div>
        </div>
    )
}

export default OrderPreviewCard
