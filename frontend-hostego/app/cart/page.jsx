"use client"

import React, { useState, useEffect } from 'react'
import BackNavigationButton from '../components/BackNavigationButton'
import CartItem from '../components/Cart/CartItem'
import AddressList from '../components/Address/AddressList'
import { Home, Clock, Truck, CreditCard } from 'lucide-react'
import HostegoButton from '../components/HostegoButton'
import axiosClient from '../utils/axiosClient'
import HostegoLoader from '../components/HostegoLoader'
import PaymentStatus from '../components/PaymentStatus'
import { useRouter } from 'next/navigation'

const page = () => {
    // Hostego – Simplify Your Hostel Life"
    const [openAddressList, setOpenAddressList] = useState(false);
    const [selectedAddress, setSelectedAddress] = useState(false)
    const [cartData, setCartData] = useState({})
    const [isPageLoading, setIsPageLoading] = useState(true)
    const [paymentStatus, setPaymentStatus] = useState(null)
    const router = useRouter()

    useEffect(() => {
        fetchCartItems()
    }, [])

    const fetchCartItems = async () => {
        try {
            setIsPageLoading(true)
            const { data } = await axiosClient.get('/api/cart/')
            setCartData(data)
            console.log(data)
        } catch (error) {
            console.error('Error fetching cart:', error)
        } finally {
            setIsPageLoading(false)
        }
    }

    if (isPageLoading) {
        return <HostegoLoader />
    }

    const handleCreateOrder = async () => {
        try {
            setPaymentStatus('processing')
            const { data } = await axiosClient.post('/api/order', {
                address_id: selectedAddress?.address_id
            })

            const response = await axiosClient.get(`/api/payment`, {
                order_id: data?.order_id
            })

            // Simulate payment processing time
            await new Promise(resolve => setTimeout(resolve, 2000))

            setPaymentStatus('success')
            // Redirect after success
            setTimeout(() => {
                router.push('/orders')
            }, 2000)

        } catch (error) {
            console.error('Error creating order:', error)
            setPaymentStatus('failed')
            // Reset status after error
            setTimeout(() => {
                setPaymentStatus(null)
            }, 3000)
        }
    }
    return (
        <div className='min-h-screen bg-[var(--bg-page-color)]'>
            <BackNavigationButton title="Checkout" />

            {/* Delivery Info Card */}
            <div className='bg-white m-2 rounded-xl overflow-hidden shadow-sm'>
                <div className='bg-gradient-to-r from-[var(--primary-color)] to-purple-500 px-4 py-3 text-white'>
                    <div className='flex items-center gap-2'>
                        <Clock className="w-5 h-5" />
                        <p className='font-medium'>Express Delivery</p>
                    </div>
                    <p className='text-sm opacity-90 mt-1'>Estimated delivery in 15-20 minutes</p>
                </div>
                <div className='p-4'>
                    <p className='text-sm text-gray-600'>Order Summary • {cartData?.cart_items?.length || 0} items</p>
                </div>
            </div>

            {/* Cart Items */}
            <div className='space-y-2 mb-4'>
                {cartData?.cart_items?.map((el) => (
                    <CartItem
                        fetchCartAgain={fetchCartItems}
                        {...el}
                        key={el?.cart_item_id}
                    />
                ))}
            </div>

            {/* Delivery Address */}
            <div className='bg-white mx-2 rounded-xl p-4 shadow-sm'>
                <div className='flex items-center justify-between mb-3'>
                    <div className='flex items-center gap-2'>
                        <Truck className='w-5 h-5 text-[var(--primary-color)]' />
                        <p className='font-medium'>Delivery Address</p>
                    </div>
                    <button
                        onClick={() => setOpenAddressList(true)}
                        className='text-sm text-[var(--primary-color)] font-medium'
                    >
                        {selectedAddress ? 'Change' : 'Add'}
                    </button>
                </div>

                {selectedAddress ? (
                    <div className='flex items-start gap-3'>
                        <div className='bg-[var(--bg-page-color)] p-2 rounded-full'>
                            <Home className='w-5 h-5 text-[var(--primary-color)]' />
                        </div>
                        <div>
                            <p className='font-medium'>{selectedAddress?.address_type}</p>
                            <p className='text-sm text-gray-600'>{selectedAddress?.address_line_1}</p>
                        </div>
                    </div>
                ) : (
                    <p className='text-sm text-gray-500'>Please select a delivery address</p>
                )}
            </div>

            {/* Bill Details */}
            <div className='bg-white mx-2 mt-4 rounded-xl p-4 shadow-sm'>
                <div className='flex items-center gap-2 mb-3'>
                    <CreditCard className='w-5 h-5 text-[var(--primary-color)]' />
                    <p className='font-medium text-md'>Bill Details</p>
                </div>
                <div className='space-y-2 text-md '>
                    <div className='flex justify-between font-normal'>
                        <span className='text-gray-800 '>Item Total</span>
                        <span>+ ₹{cartData?.cart_value.subtotal}</span>
                    </div>
                    <div className='flex justify-between  font-normal'>
                        <span className='text-gray-800'>Delivery Fee</span>
                        <span>+ ₹{cartData?.cart_value.shipping_fee}</span>
                    </div>

                    <div className='flex justify-between pt-2 border-t mt-2 font-semibold text-xl'>
                        <span>Total Amount</span>
                        <span>₹{(cartData?.cart_value.final_order_value)}</span>
                    </div>
                </div>
            </div>

            {/* Place Order Button */}
            <div className='fixed bottom-0 left-0 right-0 p-4 bg-white border-t'>
                <HostegoButton
                    onClick={handleCreateOrder}
                    text={`Place Order • ₹${cartData?.cart_value?.final_order_value}`}
                    className='w-full bg-[var(--primary-color)] text-white py-3 rounded-xl font-medium'
                    disabled={!selectedAddress}
                />
            </div>

            <AddressList
                sendSelectedAddress={setSelectedAddress}
                openAddressList={openAddressList}
                setOpenAddressList={setOpenAddressList}
            />

            {/* Bottom Spacing */}
            <div className='h-24'></div>

            {paymentStatus && <PaymentStatus status={paymentStatus} />}
        </div>
    )
}

export default page
