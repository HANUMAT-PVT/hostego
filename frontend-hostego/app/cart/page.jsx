"use client"

import React, { useState } from 'react'
import BackNavigationButton from '../components/BackNavigationButton'
import CartItem from '../components/Cart/CartItem'
import AddressList from '../components/Address/AddressList'
import { Home } from 'lucide-react'

const page = () => {
    const [openAddressList, setOpenAddressList] = useState(false);
    const [selectedAddress, setSelectedAddress] = useState(false)
    return (
        <div className='bg-[var(--bg-page-color)]'>
            <BackNavigationButton title={"Checkout"} />
            <div className='bg-white p-2'>
                <p className='font-semibold text-md pb-1 '>Delivery in 15 minutes</p>
                <p className='text-[12px] text-gray-600'>Shipment of 3 items</p>
            </div>
            <div className='p-2 flex flex-col rounded-lg '>
                <CartItem />
                <CartItem />
                <CartItem />


                {/* Bill Details */}
                <div className=' flex flex-col gap-2 bg-white p-2 mb-2 mt-2'>
                    <p className='font-semibold text-sm pb-1 border-b   '>Bill details</p>
                    <div className='flex justify-between  font-normal text-xs  text-gray-900'>
                        <p>Items total</p>
                        <p>₹230</p>
                    </div>
                    <div className='flex justify-between font-normal text-xs  text-gray-900'>
                        <p className=''>Delivery charge</p>
                        <p>₹23</p>
                    </div>
                    <div className='flex justify-between font-normal text-xs  text-gray-900 border-b pb-2'>
                        <p>Platform fee</p>
                        <p>+₹1</p>
                    </div>
                    <div className='flex justify-between text-xs font-light text-gray-900'>
                        <p className='font-semibold text-lg '>Grand total</p>
                        <p className='font-semibold text-lg '>₹254</p>
                    </div>

                </div>
                {/* Cancellation Policy */}
                <div className=' flex flex-col gap-2 bg-white p-2 mb-8 rounded-md'>
                    <p className='font-semibold text-sm pb-1 border-b   '>Cancellation Policy</p>
                    <p className='text-xs text-gray-500'>Orders cannot be cancelled once packed for delivery.
                        In case of unexpected delays, a refund will be provided, if applicable.</p>
                </div>
                <div className='h-[150px]'>

                </div>
                {<div className="fixed bottom-0 w-full p-2  bg-white z-2 flex justify-center items-center px-4 shadow-2xl  ">
                    {!selectedAddress && <button onClick={() => setOpenAddressList(!openAddressList)} className='bg-green-700 text-white w-[85vw] text-xs px-4 py-2 rounded-md font-normal '>Choose address at next step</button>}
                    {selectedAddress && <div className='w-full'>
                        <div onClick={() => { }} className='address-item flex items-center  gap-4 pb-2 rounded-md  cursor-pointer bg-white  '>
                            <div className='bg-[var(--bg-page-color)] p-2 w-[30px] h-[30px]  flex justify-center items-center rounded-full '>
                                <Home size={20} className='text-[var(--primary-color)]' />
                            </div>
                            <div className='w-full flex items-center justify-between mb-2 '>
                                <div>
                                    <p className='text-xs font-semibold'>Delivering to {selectedAddress?.heading}</p>
                                    <p className='text-[11px] w-[180px]  whitespace-nowrap overflow-hidden text-ellipsis  '>{selectedAddress?.street}</p>
                                </div>
                                <p className='text-[12px]  text-green-700' onClick={() => setOpenAddressList(!openAddressList)}>Change</p>
                            </div>
                        </div>
                        {/*Payment Option  */}
                        <div className='flex gap-2 justify-between items-center'>
                            <div>
                                <p className='text-[12px]  text-gray-500'>PAY USING</p>
                                <p className='text-[12px]'>Wallet</p>
                            </div>
                            <div className='bg-green-700 text-white min-w-[180px] px-3 py-1 h-fit flex justify-between rounded-md items-center '>
                                <div className='flex flex-col items-center  '>
                                    <p className='font-semibold text-[13px]'>₹254</p>
                                    <p className='font-normal text-[11px]'>TOTAL</p>
                                </div>
                                <p className='text-[16px]'>Place Order</p>
                            </div>
                        </div>
                    </div>}
                </div>}
            </div>
            <AddressList sendSelectedAddress={(e) => setSelectedAddress(e)} openAddressList={openAddressList} setOpenAddressList={(e) => setOpenAddressList(e)} />
        </div>
    )
}

export default page
