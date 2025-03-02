"use client"

import React, { useState, useEffect } from 'react'
import BackNavigationButton from '../components/BackNavigationButton'
import CartItem from '../components/Cart/CartItem'
import AddressList from '../components/Address/AddressList'
import { Home, Clock, Truck, CreditCard, MapPin, Timer, AlertCircle, TicketCheck } from 'lucide-react'
import HostegoButton from '../components/HostegoButton'
import axiosClient from '../utils/axiosClient'
import HostegoLoader from '../components/HostegoLoader'
import PaymentStatus from '../components/PaymentStatus'
import { useRouter } from 'next/navigation'
import HostegoToast from '../components/HostegoToast'
import { useDispatch, useSelector } from 'react-redux'
import { setFetchCartData } from '../lib/redux/features/user/userSlice'

const AddressSection = ({ selectedAddress, setOpenAddressList }) => {
    return (
        <div className={`bg-white mx-2 rounded-xl p-4 shadow-sm transition-all duration-200 
            ${!selectedAddress ? 'border-2 border-red-200 ' : 'border border-gray-100'}`}>
            <div className='flex items-center justify-between mb-3'>
                <div className='flex items-center gap-2'>
                    <Truck className={`w-5 h-5 ${!selectedAddress ? 'text-red-500' : 'text-[var(--primary-color)]'}`} />
                    <p className='font-medium'>Delivery Address</p>
                </div>
                <button
                    onClick={() => setOpenAddressList(true)}
                    className={`px-3 py-1.5 rounded-full text-sm font-medium transition-all
                        ${!selectedAddress
                            ? 'bg-red-50 text-red-500 hover:bg-red-100'
                            : 'bg-[var(--primary-color)]/10 text-[var(--primary-color)] hover:bg-[var(--primary-color)]/20'
                        }`}
                >
                    {selectedAddress ? 'Change' : 'Add Address'}
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
                <div className='flex items-center gap-3 p-3 bg-red-50 rounded-lg border border-red-100'>
                    <div className='w-10 h-10 rounded-full bg-red-100 flex items-center justify-center'>
                        <MapPin className='w-5 h-5 text-red-500' />
                    </div>
                    <div>
                        <p className='font-medium text-red-600'>Delivery Address Required</p>
                        <p className='text-sm text-red-500'>Please select a delivery address to continue</p>
                    </div>
                </div>
            )}
        </div>
    )
}

const page = () => {
    // Hostego – Simplify Your Hostel Life"
    const [openAddressList, setOpenAddressList] = useState(false);
    const [selectedAddress, setSelectedAddress] = useState(false)
    const [isPageLoading, setIsPageLoading] = useState(true)
    const [paymentStatus, setPaymentStatus] = useState(null)
    const [orderTimer, setOrderTimer] = useState(10)
    const [isTimerRunning, setIsTimerRunning] = useState(false)
    const [isToastVisible, setIsToastVisible] = useState(false)
    const [isPageFirstLoad, setIsPageFirstLoad] = useState(true)
    const [cookingRequests, setCookingRequests] = useState('')

    const dispatch = useDispatch()
    const { cartData } = useSelector((state) => state.user)
    const router = useRouter()


    useEffect(() => {
        fetchCartItems()
    }, [])

    useEffect(() => {
        if (isPageFirstLoad) {
            setIsPageFirstLoad(false)
            return
        }
    }, [isPageFirstLoad])

    useEffect(() => {
        let interval
        if (isTimerRunning && orderTimer > 0) {
            interval = setInterval(() => {
                setOrderTimer((prev) => prev - 1)
            }, 1000)
        } else if (orderTimer === 0) {
            handleCreateOrder()
        }

        return () => clearInterval(interval)
    }, [isTimerRunning, orderTimer])

    const fetchCartItems = async () => {
        try {
            if (!isPageFirstLoad) setIsPageLoading(true)
            dispatch(setFetchCartData(true))
        } catch (error) {
            console.error('Error fetching cart:', error)
        } finally {
            setIsPageLoading(false)

        }
    }

    const startOrderTimer = () => {
        try {

            if (!selectedAddress) {
                setIsToastVisible(true)
                return
            }
            setIsTimerRunning(true)
        } catch (error) {
            console.error('Error starting order timer:', error)
        }
    }

    const cancelOrder = () => {
        setIsTimerRunning(false)
        setOrderTimer(10)
    }

    const handleCreateOrder = async () => {
        try {

            setPaymentStatus('processing')
            const { data } = await axiosClient.post('/api/order', {
                address_id: selectedAddress?.address_id,
                cooking_requests: cookingRequests
            })

            const response = await axiosClient.post(`/api/payment`, {
                order_id: data?.order_id
            })

            if (response.data) {
                setPaymentStatus('success')
                dispatch(setFetchCartData(true))
                setTimeout(() => {
                    router.push('/orders')
                }, 2000)

            } else {
                setPaymentStatus('failed')
            }
        } catch (error) {
            console.error('Error processing order:', error)
            setPaymentStatus('failed')
        } finally {
            setIsTimerRunning(false)
            setOrderTimer(10)
        }
    }

    if (isPageLoading) {
        return <HostegoLoader />
    }


    return (
        <div className='min-h-screen bg-[var(--bg-page-color)]'>
            <BackNavigationButton title="Checkout" />
            <HostegoToast message="Please select a delivery address" variant="error" show={isToastVisible} onClose={() => setIsToastVisible(false)} />
            {/* Timer Banner - shows when timer is running */}
            {isTimerRunning && (
                <div className="fixed top-16 left-0 right-0 bg-white shadow-md z-10">
                    <div className="w-fit mx-auto p-4">
                        <div className="flex items-center flex-col justify-between mb-3 gap-2">
                            <div className="flex items-center gap-2">
                                <Timer className="w-5 h-5 text-[var(--primary-color)] animate-pulse" />
                                <span className="font-medium">
                                    Order will be placed in {orderTimer} seconds
                                </span>
                            </div>
                            <div className="flex items-center gap-2">
                                <button
                                    onClick={handleCreateOrder}
                                    className="px-4 py-1.5 bg-green-600 text-white rounded-full text-sm font-medium
                                             hover:bg-green-700 transition-colors text-nowrap"
                                >
                                    Place Now
                                </button>
                                <button
                                    onClick={cancelOrder}
                                    className="px-3 py-1.5 bg-red-50 text-red-600 rounded-full text-sm font-medium
                                             hover:bg-red-100 transition-colors"
                                >
                                    Cancel
                                </button>
                            </div>
                        </div>
                        {/* Progress bar */}
                        <div className="h-1 bg-gray-100 rounded-full overflow-hidden">
                            <div
                                className="h-full bg-[var(--primary-color)] transition-all duration-1000"
                                style={{ width: `${(orderTimer / 30) * 100}%` }}
                            />
                        </div>
                    </div>
                </div>
            )}

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
            <AddressSection
                selectedAddress={selectedAddress}
                setOpenAddressList={setOpenAddressList}
            />
            {/* Extra Infor */}
            <div className='bg-white mx-2 mt-4 rounded-xl p-4 shadow-sm'>
                <div className="relative">
                    <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
                        Type Cooking Requests
                    </label>
                    <textarea
                        value={cookingRequests}
                        onChange={(e) => setCookingRequests(e.target.value)}
                        className="w-full px-4 py-3 border-2 border-[var(--primary-color)] rounded-md outline-none min-h-[70px] resize-none"
                        placeholder="Enter your complete cooking requests"
                        required
                    />
                </div>
            </div>
            {/* Bill Details */}
            <div className='bg-white mx-2 mt-4 rounded-xl p-4 shadow-sm'>
                <div className='flex items-center gap-2 mb-4'>
                    <CreditCard className='w-5 h-5 text-[var(--primary-color)]' />
                    <p className='font-medium text-md'>Bill Details</p>
                </div>

                <div className='space-y-2 text-md'>
                    {/* Item Total */}
                    <div className='flex justify-between font-normal'>
                        <span className='text-gray-800'>Item Total</span>
                        <span>₹{cartData?.cart_value?.subtotal}</span>
                    </div>

                    {/* Delivery Fee */}
                    <div className='flex justify-between font-normal items-start'>
                        <span className='text-gray-800'>Delivery Fee</span>
                        {cartData?.free_delivery ? (

                            <div className='flex items-center gap-2'>
                                <p className='line-through'>₹{cartData?.cart_value?.actual_shipping_fee}</p>
                                <div className='bg-gradient-to-r from-[#655df0] to-[#9333ea] text-white px-3  rounded-md'>
                                    <span className='font-bold tracking-wide'>FREE</span>
                                </div>

                            </div>
                        ) : (
                            // Regular Delivery Fee Display
                            <div className='text-right'>
                                <div className='flex items-center gap-2'>
                                    <span>₹{cartData?.cart_value?.actual_shipping_fee}</span>
                                </div>

                            </div>
                        )}
                    </div>

                    {/* Additional Savings Banner for Free Delivery */}
                    {cartData?.free_delivery && (
                        <div className='bg-gradient-to-r from-[#655df0] to-[#9333ea] p-0.5 rounded-lg mt-2'>
                            <div className='bg-white rounded-[7px] p-3 flex items-center gap-3'>
                                <div className='w-10 h-10 rounded-full bg-gradient-to-r from-[#655df0]/10 to-[#9333ea]/10 flex items-center justify-center'>
                                    <span className='text-[#655df0] font-bold'>₹</span>
                                </div>
                                <div>
                                    <p className='font-medium text-gray-800'>Welcome To HOSTEGO !</p>
                                    <p className='text-sm text-gray-600'>
                                        You saved <span className='font-bold'>₹{cartData?.cart_value.actual_shipping_fee}</span> with <span className='font-bold'>FREE DELIVERY </span>
                                    </p>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Total Amount */}
                    <div className='flex justify-between pt-3 border-t mt-2'>
                        <div>
                            <span className='font-semibold text-xl'>Total Amount</span>
                            <p className='text-xs text-gray-500'>Inclusive of all taxes</p>
                        </div>
                        <span className='font-semibold text-xl'>
                            ₹{cartData?.cart_value?.final_order_value}
                        </span>
                    </div>
                </div>
            </div>

            {/* Place Order Button */}
            <div className='fixed bottom-0 left-0 right-0 p-4 bg-white border-t'>
                {!selectedAddress && (
                    <div className='text-sm text-red-500 text-center mb-2'>
                        ⚠️ Please select a delivery address
                    </div>
                )}
                <HostegoButton
                    onClick={startOrderTimer}
                    text={`Place Order • ₹${cartData?.cart_value?.final_order_value}`}
                    className={`w-full py-3 rounded-xl font-medium transition-all duration-200
                        ${!selectedAddress
                            ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                            : 'bg-[var(--primary-color)] text-white hover:opacity-90'
                        }`}
                    disabled={!selectedAddress || isTimerRunning}
                />
            </div>

            <AddressList
                showAddressButton={false}
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


