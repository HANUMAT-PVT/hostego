import React from 'react'
import BackNavigationButton from '../../components/BackNavigationButton'
const page = () => {
    return (
        <div className='bg-[var(--bg-page-color)]' >
            <BackNavigationButton />
            <div className='px-4 py-2 bg-white mb-2 '>
                <div className='mb-4 '>
                    <p className='font-bold text-lg'>Order Summary</p>
                    <p className='text-gray-500 text-xs font-normal'>Arrived at 8:35 pm</p>
                </div>
                <p className='font-semibold text-sm mb-4 '>3 items in this order</p>
                <div className='flex flex-col gap-3'>
                    <div className='flex gap-5 justify-between items-center '>
                        <div className='flex gap-5 '>
                            <div className='bg-[var(--bg-page-color)] rounded-md p-1'>
                                <img className=' min-w-[50px] max-w-[50px] '
                                    src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                                    alt={'Uncle chips'}
                                />
                            </div>
                            <div>
                                <p className='text-xs'>Lay's India's Magic Masala Potato Chips</p>
                                <p className='text-[11px] text-gray-500 font-light'>48 g x 1</p>
                            </div>
                        </div>
                        <div className='flex flex-col gap-2'>

                            <p className='text-xs font-semibold'>₹20</p>

                        </div>
                    </div>
                    <div className='flex gap-5 justify-between items-center '>
                        <div className='flex gap-5 '>
                            <div className='bg-[var(--bg-page-color)] rounded-md p-1'>
                                <img className=' min-w-[50px] max-w-[50px] '
                                    src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                                    alt={'Uncle chips'}
                                />
                            </div>
                            <div>
                                <p className='text-xs'>Lay's India's Magic Masala Potato Chips</p>
                                <p className='text-[11px] text-gray-500 font-light'>48 g x 1</p>
                            </div>
                        </div>
                        <div className='flex flex-col gap-2'>

                            <p className='text-xs font-semibold'>₹20</p>

                        </div>
                    </div> <div className='flex gap-5 justify-between items-center '>
                        <div className='flex gap-5 '>
                            <div className='bg-[var(--bg-page-color)] rounded-md p-1'>
                                <img className=' min-w-[50px] max-w-[50px]'
                                    src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                                    alt={'Uncle chips'}
                                />
                            </div>
                            <div>
                                <p className='text-xs'>Lay's India's Magic Masala Potato Chips</p>
                                <p className='text-[11px] text-gray-500 font-light'>48 g x 1</p>
                            </div>
                        </div>
                        <div className='flex flex-col gap-2'>

                            <p className='text-xs font-semibold'>₹20</p>

                        </div>
                    </div>
                </div>
            </div>
            {/* Bill Details */}
            <div className=' flex flex-col gap-2 bg-white px-4 py-2 mb-2'>
                <p className='font-semibold text-sm pb-1 border-b   '>Bill details</p>
                <div className='flex justify-between text-xs font-light text-gray-900'>
                    <p>MRP</p>
                    <p>₹230</p>
                </div>
                <div className='flex justify-between text-xs font-light text-gray-900'>
                    <p>Delivery charge</p>
                    <p>₹23</p>
                </div>
                <div className='flex justify-between text-xs font-light text-gray-900'>
                    <p>Platform fee</p>
                    <p>+₹1</p>
                </div>
                <div className='flex justify-between text-xs font-light text-gray-900'>
                    <p className='font-semibold'>Bill total</p>
                    <p className='font-semibold'>₹254</p>
                </div>

            </div>
            {/* ORrder details */}
            <div className=' flex flex-col gap-2 bg-white px-4 py-2'>
                <p className='font-semibold text-sm pb-1 border-b   '>Order details</p>
                <div className='flex flex-col gap-2 text-xs font-light text-gray-900'>
                    <p>order id</p>
                    <p className='font-semibold '>ORD33424213232</p>
                </div>
                <div className='flex flex-col gap-2 text-xs font-light text-gray-900'>
                    <p>Payment</p>
                    <p className='font-semibold'>Paid Online</p>
                </div>
                <div className='flex flex-col gap-2 text-xs font-light text-gray-900'>
                    <p>Deliver to</p>
                    <p className='font-semibold'>Room no. 1115, Zakir-A, Chandigarh University</p>
                </div>
                <div className='flex flex-col flex-col gap-2 text-xs font-light text-gray-900'>
                    <p>Order placed</p>
                    <p className='font-semibold'>placed on Fri, 14 Feb'25, 8:12 PM</p>
                </div>


            </div>
        </div>
    )
}

export default page
