"use client"

import React, { useState, useEffect } from 'react'
import BackNavigationButton from '../components/BackNavigationButton'
import CartItem from '../components/Cart/CartItem'
import AddressList from '../components/Address/AddressList'
import { Home, Clock, Truck, CreditCard, MapPin, Timer, AlertCircle, TicketCheck, Wallet, ArrowRight, X, CheckCircle, Info } from 'lucide-react'
import HostegoButton from '../components/HostegoButton'
import axiosClient from '../utils/axiosClient'
import HostegoLoader from '../components/HostegoLoader'
import PaymentStatus from '../components/PaymentStatus'
import { useRouter } from 'next/navigation'
import HostegoToast from '../components/HostegoToast'
import { useDispatch, useSelector } from 'react-redux'
import { setFetchCartData, setFetchUserWalletBool, setUserAccountWallet } from '../lib/redux/features/user/userSlice'
import { subscribeToNotifications } from '../utils/webNotifications'
import { load } from "@cashfreepayments/cashfree-js";


const AddressSection = ({ selectedAddress, setOpenAddressList }) => {
    return (
        <div onClick={() => setOpenAddressList(true)} className={`bg-white mx-2 rounded-xl p-4 shadow-sm transition-all duration-200 
            ${!selectedAddress ? 'border-2 border-red-200 ' : 'border border-gray-100'}`}>
            <div className='flex items-center justify-between mb-3'>
                <div className='flex items-center gap-2'>
                    <Truck className={`w-5 h-5 ${!selectedAddress ? 'text-red-500' : 'text-[var(--primary-color)]'}`} />
                    <p className='font-medium'>Delivery Address</p>
                </div>
                <button

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


    let cashfree;
    var initializeSDK = async function () {
        cashfree = await load({
            mode: process.env.NODE_ENV == "production" ? "production" : "sandbox",
        });
    };
    initializeSDK();



    const [openAddressList, setOpenAddressList] = useState(false);
    const [selectedAddress, setSelectedAddress] = useState(false)
    const [isPageLoading, setIsPageLoading] = useState(true)
    const [paymentStatus, setPaymentStatus] = useState(null)
    const [orderTimer, setOrderTimer] = useState(10)
    const [isTimerRunning, setIsTimerRunning] = useState(false)
    const [isToastVisible, setIsToastVisible] = useState(false)
    const [isPageFirstLoad, setIsPageFirstLoad] = useState(true)
    const [cookingRequests, setCookingRequests] = useState('')
    const [showPaymentDrawer, setShowPaymentDrawer] = useState(false);

    const dispatch = useDispatch()
    const { cartData, userWallet } = useSelector((state) => state.user)
    const router = useRouter()

    // Calculate wallet status
    const hasInsufficientBalance = userWallet?.balance < cartData?.cart_value?.final_order_value;
    const amountNeeded = cartData?.cart_value?.final_order_value - (userWallet?.balance || 0);

    // Add new function to fetch wallet
    const fetchUserWallet = async () => {
        try {
            const { data } = await axiosClient.get("/api/wallet")
            dispatch(setUserAccountWallet(data))
        } catch (error) {
            console.error('Error fetching wallet:', error)
        }
    }

    useEffect(() => {
        // Fetch cart items and wallet data when page loads
        fetchCartItems()
        fetchUserWallet()
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

    const handleCashFreePayment = async () => {
        setPaymentStatus('processing')
        const { data } = await axiosClient.post('/api/order', {
            address_id: selectedAddress?.address_id,
            cooking_requests: cookingRequests
        })

        const response = await axiosClient.post(`/api/payment/cashfree`, {
            order_id: data?.order_id
        })

        await doPayment(response?.data?.payment_session_id, response?.data?.order_id, data?.order_id)

    }
    const doPayment = async (paymentSessionId, paymentSessionorderId, order_id) => {

        let checkoutOptions = {
            paymentSessionId: paymentSessionId,
            redirectTarget: "_modal",
        };
        let result = await cashfree.checkout(checkoutOptions)

        if (result.error) {
            // This will be true whenever user clicks on close icon inside the modal or any error happens during the payment
            console.log("User has closed the popup or there is some payment error, Check for Payment Status");
            setPaymentStatus('failed')
        }
        if (result.redirect) {

        }
        if (result.paymentDetails) {
            const tryPaymentStatus = (paymentSessionorderId, order_id, maxAttempts = 3, delay = 3000) => {
                return new Promise((resolve, reject) => {
                    let attempts = 0;

                    const interval = setInterval(async () => {
                        try {
                            attempts++;

                            const result = await axiosClient.post(`/api/payment/cashfree/${paymentSessionorderId}`, {
                                order_id: order_id
                            });

                            if (result?.data?.response?.order_status == "PAID") {
                                clearInterval(interval);
                                resolve(result);
                            } else if (attempts >= maxAttempts) {
                                clearInterval(interval);
                                reject(new Error("Max attempts reached. Payment details not available."));
                            }

                        } catch (err) {
                            clearInterval(interval);
                            reject(err);
                        }
                    }, delay);
                });
            };
            try {
                const paymentResult = await tryPaymentStatus(paymentSessionorderId, order_id);
                setPaymentStatus('success')
                router.push("/orders")
            } catch (error) {
                console.log(error)
            }
        }


    }
    const handleWalletCreateOrder = async () => {
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
                subscribeToNotifications("Payment Success", "Your order has been placed successfully")
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

    const handleCreateOrder = () => {
        setPaymentStatus('processing')
        if (!hasInsufficientBalance) {
            handleWalletCreateOrder()
        } else {
            handleCashFreePayment()
        }
    }

    const handlePlaceOrder = () => {
        if (cartData?.cart_items?.length === 0) {
            return;
        }
        if (!selectedAddress) {
            setIsToastVisible(true);
            return;
        }
        setShowPaymentDrawer(true);
    };



    if (isPageLoading) {
        return <HostegoLoader />
    }


    return (
        <div className='min-h-screen bg-[var(--bg-page-color)]'>
            <BackNavigationButton title="Checkout" />
            {/* <CartCheckout paymentSessionId={response?.data?.payment_session_id} /> */}
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
                                style={{ width: `${(orderTimer / 10) * 100}%` }}
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
                    <p className='text-sm opacity-90 mt-1'>Estimated delivery in 15-30 minutes</p>
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
                    {/* Rain extra charge Fee */}
                    <div className='flex justify-between font-normal items-start'>
                        <span className='text-gray-800 text-xs'>🌧️ Rain Chage </span>

                        {/* // Regular Delivery Fee Display */}
                        <div className='text-right'>
                            <div className='flex items-center gap-2'>
                                <span>+ ₹{cartData?.cart_value?.rain_extra_fee}</span>
                            </div>

                        </div>

                    </div>
                    {/* Delivery Fee */}
                    <div className='flex justify-between font-normal items-start'>
                        <span className='text-gray-800'>Delivery Fee <br /><span className='text-xs flex items-center'>   Includes delivery & payment fees</span></span>

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
                            +₹{cartData?.cart_value?.final_order_value}
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
                    onClick={handlePlaceOrder}
                    text={`Place Order • ₹${cartData?.cart_value?.final_order_value || 0}`}
                    className={`w-full py-3 rounded-xl font-medium transition-all duration-200
                        ${!selectedAddress
                            ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                            : 'bg-[var(--primary-color)] text-white hover:opacity-90'
                        }`}
                    disabled={!selectedAddress || isTimerRunning}
                />
            </div>

            {/* Payment Confirmation Drawer */}
            {showPaymentDrawer && (
                <div className="fixed inset-0 bg-black/50 z-50">
                    <div
                        className="fixed bottom-0 left-0 right-0 bg-white rounded-t-2xl p-4 animate-slide-up"
                        style={{ maxHeight: '80vh', overflowY: 'auto' }}
                    >
                        {/* Drawer Header */}
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-semibold">Payment Confirmation</h3>
                            <button
                                onClick={() => setShowPaymentDrawer(false)}
                                className="p-2 hover:bg-gray-100 rounded-full"
                            >
                                <X size={20} />
                            </button>
                        </div>

                        {/* Wallet Status */}
                        <div className="bg-gray-50 rounded-xl p-4 mb-4">
                            <div className="flex justify-between items-center mb-3">
                                <div className="flex items-center gap-2">
                                    <Wallet className="text-[var(--primary-color)]" size={20} />
                                    <span className="font-medium">Wallet Balance</span>
                                </div>
                                <span className="text-xl font-semibold">₹{(userWallet?.balance).toFixed(1) || 0}</span>
                            </div>
                            <div className="flex justify-between items-center text-sm text-gray-600">
                                <span>Order Amount</span>
                                <span>₹{cartData?.cart_value?.final_order_value}</span>
                            </div>
                        </div>

                        {/* Status Message */}
                        <div className={`rounded-xl p-4 mb-6 ${hasInsufficientBalance ? 'bg-red-50' : 'bg-green-50'}`}>
                            <div className="flex items-start gap-3">
                                {hasInsufficientBalance ? (
                                    <>
                                        <AlertCircle className="text-red-600 flex-shrink-0" size={20} />
                                        <div>
                                            <p className="font-medium text-red-600">Insufficient Balance</p>
                                            <p className="text-sm text-red-500">Add ₹{(amountNeeded).toFixed(1)} more to place order</p>
                                        </div>
                                    </>
                                ) : (
                                    <>
                                        <CheckCircle className="text-green-600 flex-shrink-0" size={20} />
                                        <div>
                                            <p className="font-medium text-green-600">Sufficient Balance</p>
                                            <p className="text-sm text-green-500">Your wallet has enough balance</p>
                                        </div>
                                    </>
                                )}
                            </div>
                        </div>

                        {/* Action Button */}
                        {hasInsufficientBalance ? (
                            <button
                                onClick={() => {
                                    setShowPaymentDrawer(false);
                                    handleCreateOrder()

                                }}
                                className="w-full bg-[var(--primary-color)] text-white py-3 rounded-xl font-medium  transition-colors flex items-center justify-center gap-2"
                            >
                                Pay Online
                                <ArrowRight size={18} />
                            </button>
                        ) : (
                            <button
                                onClick={() => {
                                    setShowPaymentDrawer(false);
                                    startOrderTimer()
                                    handleCreateOrder();
                                }}
                                className="w-full bg-[var(--primary-color)] text-white py-3 rounded-xl font-medium hover:opacity-90 transition-opacity"
                            >
                                Confirm Order
                            </button>
                        )}
                    </div>
                </div>
            )}

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


