'use client'
import React, { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import BackNavigationButton from '@/app/components/BackNavigationButton'
import { formatDate, tryPaymentStatus } from '@/app/utils/helper'
import { Package, MapPin, Clock, CheckCircle2, Truck, Check, AlertCircle, IndianRupee, RefreshCcw, Phone, Bike, MessageSquare, Award } from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '@/app/components/HostegoLoader'
import StatusTimeLine from '../../components/Orders/StatusTimeLine'
import { transformOrder } from '../../utils/helper'
import { ORDER_STATUSES } from '../../components/Delivery-Partner/MaintainOrderStatusForDeliveryPartner'
import { subscribeToNotifications } from '../../utils/webNotifications'
import ConfirmationPopup from '@/app/components/ConfirmationPopup'

const DeliveryPartnerSection = ({ order, isActiveOrder }) => {
    const [timeElapsed, setTimeElapsed] = useState('');

    // Calculate time elapsed since order creation
    useEffect(() => {
        const calculateTimeElapsed = () => {
            const orderTime = new Date(order?.created_at);
            const now = new Date();
            const diff = Math.floor((now - orderTime) / 1000 / 60); // minutes
            setTimeElapsed(
                diff < 60
                    ? `${diff} min ago`
                    : `${Math.floor(diff / 60)}h ${diff % 60}m ago`
            );
        };

        calculateTimeElapsed();
        const timer = setInterval(calculateTimeElapsed, 60000); // Update every minute

        return () => clearInterval(timer);
    }, [order?.created_at]);

    return (
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
            {/* Status Banner */}
            <div className="bg-gradient-to-r from-[var(--primary-color)] to-purple-600 p-4 text-white">
                <div className="flex items-center justify-between">
                    <h3 className="font-medium">Delivery Partner</h3>
                    {isActiveOrder && (
                        <span className="flex items-center gap-2 text-sm">
                            <span className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
                            Active Order
                        </span>
                    )}
                </div>
                <p className="text-sm opacity-80 mt-1">{timeElapsed}</p>
            </div>

            {/* Partner Details */}
            <div className="p-4">
                <div className="flex items-center gap-4 mb-6">
                    <div className="w-16 h-16 rounded-full bg-[var(--primary-color)]/10 flex items-center justify-center flex-shrink-0">
                        <Bike className="w-8 h-8 text-[var(--primary-color)]" />
                    </div>
                    <div>
                        <h3 className="font-medium text-lg text-gray-900">
                            {order?.delivery_partner?.user?.first_name} {order?.delivery_partner?.user?.last_name}
                        </h3>
                        <p className="text-gray-500">Your Delivery Partner</p>
                    </div>
                </div>

                {/* Contact and Info */}
                <div className="space-y-4">
                    {/* Call Button */}
                    <button
                        onClick={() => window.location.href = `tel:${order?.delivery_partner?.user?.mobile_number}`}
                        className="w-full py-3 px-4 rounded-xl border-2 border-[var(--primary-color)] 
                                 bg-[var(--primary-color)]/5 hover:bg-[var(--primary-color)]/10
                                 text-[var(--primary-color)] font-medium flex items-center justify-center gap-3
                                 transition-all duration-200"
                    >
                        <div className="relative">
                            <Phone className="w-5 h-5" />
                            <span className="absolute -top-1 -right-1 w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                        </div>
                        <span>Call Partner</span>
                    </button>

                    {/* Delivery Info */}
                    <div className="bg-gray-50 rounded-xl p-4">
                        <div className="flex items-center gap-2 text-gray-600">
                            <Clock className="w-4 h-4" />
                            <span>Estimated delivery in {order?.estimated_delivery_time || '20-30'} minutes</span>
                        </div>
                    </div>

                    {/* Safety Tip */}
                    <div className="bg-yellow-50 rounded-xl p-4 flex items-start gap-3">
                        <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                        <div>
                            <p className="font-medium text-yellow-800 mb-1">Safety First</p>
                            <p className="text-sm text-yellow-700">
                                For your safety, please don't share any sensitive information with the delivery partner.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

const OrderDetailsPage = () => {
    const { id } = useParams()
    const [order, setOrder] = useState(null)
    const [isLoading, setIsLoading] = useState(true)
    const [isConfirming, setIsConfirming] = useState(false)
    const [isConfirmationPopupOpen, setIsConfirmationPopupOpen] = useState(false)

    const statusConfig = {
        pending: {
            icon: Clock,
            color: 'text-orange-500',
            bgColor: 'bg-orange-50',
            label: 'Payment Pending'
        },
        placed: {
            icon: Package,
            color: 'text-blue-500',
            bgColor: 'bg-blue-50',
            label: 'Order Placed'
        },
        assigned: {
            icon: Package,
            color: 'text-[var(--primary-color)]',
            bgColor: 'bg-[var(--primary-color)]/10',
            label: 'Asssigned'
        },

        reached: {
            icon: Check,
            color: 'text-purple-500',
            bgColor: 'bg-purple-50',
            label: 'Reached Shop'
        },
        picked: {
            icon: Check,
            color: 'text-purple-500',
            bgColor: 'bg-purple-50',
            label: 'Picked Up'
        },
        on_the_way: {
            icon: Check,
            color: 'text-purple-500',
            bgColor: 'bg-purple-50',
            label: 'On The Way'
        },
        reached_door: {
            icon: Check,
            color: 'text-purple-500',
            bgColor: 'bg-purple-50',
            label: 'Reached Door'
        },
        delivered: {
            icon: CheckCircle2,
            color: 'text-green-500',
            bgColor: 'bg-green-50',
            label: 'Delivered'
        },
        cancelled: {
            icon: AlertCircle,
            color: 'text-red-500',
            bgColor: 'bg-red-50',
            label: 'Cancelled'
        }
    }

    // Check if order is in final state
    const isOrderActive = (status) => {
        return !['delivered', 'cancelled'].includes(status?.toLowerCase());
    };

    const verifythePendingOrder = async (order) => {

        if (order && order.order_status == "pending") {
            let {data} = await tryPaymentStatus(order.order_id);
            if (data?.response?.order_status == "PAID") {
                window.location.reload()
            }
        }
    }

    useEffect(() => {
        fetchOrder()
    }, [id])

    const fetchOrder = async () => {
        try {
            setIsLoading(true)
            const { data } = await axiosClient.get(`/api/order/${id}`)
            setOrder(data)

            await verifythePendingOrder(data)

        } catch (error) {
            console.error('Error fetching order:', error)
        } finally {
            setIsLoading(false)
        }
    }

    const handleOrderDelivered = async () => {
        try {
            setIsConfirming(true)
            await axiosClient.patch(`/api/order/${id}`, {
                order_status: 'delivered'
            })
            await fetchOrder() // Refresh order data
        } catch (error) {
            console.error('Error confirming delivery:', error)
        } finally {
            setIsConfirming(false)
            setIsConfirmationPopupOpen(false)
        }
    }

  

    if (isLoading) return <HostegoLoader />
    if (!order) return <div>Order not found</div>

    const status = statusConfig[order?.order_status] || statusConfig.pending
    const StatusIcon = status.icon

    return (
        <div className="min-h-screen bg-[var(--bg-page-color)]">
            <BackNavigationButton title="Order Details" />

            {/* Order Status Card */}
            <div className={`mx-2 mt-2 rounded-xl overflow-hidden bg-white shadow-sm`}>
                <div className={`p-4 ${status.bgColor} border-b flex justify-between items-center`}>
                    <div className='flex flex-col '>
                        <div className="flex items-center gap-3 mb-2">
                            <StatusIcon className={`w-6 h-6 ${status?.color}`} />
                            <h2 className={`text-lg font-semibold ${status?.color}`}>
                                {status?.label}
                            </h2>
                        </div>

                        <p className="text-sm text-gray-600">
                            Ordered on {formatDate(order?.created_at)}
                        </p>
                    </div>
                    <div className='flex items-center justify-end'>
                        <RefreshCcw onClick={()=>fetchOrder()} className="w-4 h-4 text-[var(--primary-color)] cursor-pointer" />
                    </div>
                </div>


                {/* Order Items */}
                <div className="p-4">
                    <h3 className="font-medium mb-3">Order Items</h3>
                    <div className="space-y-3">
                        {order?.order_items?.map((item) => (
                            <div key={item?.cart_item_id}
                                className="flex items-center gap-3 p-2 bg-gray-50 rounded-lg">
                                <img
                                    src={item?.product_item?.product_img_url}
                                    alt={item?.product_item?.product_name}
                                    className="w-16 h-16 rounded-lg object-cover"
                                />
                                <div className="flex-1">
                                    <h4 className="font-medium text-sm">
                                        {item?.product_item?.product_name} <span className='text-gray-600 text-xs'>( {item?.product_item?.shop?.shop_name} )</span>
                                    </h4>

                                    <p className="text-sm text-gray-600">
                                        {item?.product_item?.weight}
                                    </p>
                                    <div className="flex items-center gap-2 mt-1">
                                        <span className="text-sm">₹{item?.product_item?.selling_price}</span>
                                        <span className="text-gray-400">×</span>
                                        <span className="text-sm">{item?.quantity}</span>
                                    </div>
                                </div>
                                <div className="text-right">
                                    <p className="font-medium">₹{item?.sub_total}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
            {/* Delivery Partner Contact Section */}
            {order?.delivery_partner && (
                <div className="mx-4 my-4">
                    <DeliveryPartnerSection
                        order={order}
                        isActiveOrder={isOrderActive(order.order_status)}
                    />
                </div>
            )}
            {/* Delivery Address */}
            <div className="bg-white mx-2 mt-3 rounded-xl p-4 shadow-sm">
                <div className="flex items-center gap-2 mb-3">
                    <MapPin className="w-5 h-5 text-[var(--primary-color)]" />
                    <h3 className="font-medium">Delivery Address</h3>
                </div>
                <div className="text-sm text-gray-600">
                    <p className="font-medium text-gray-900 ml-6">{order?.address?.address_type}</p>
                    <p className='ml-6'>{order?.address?.address_line_1}</p>
                </div>
            </div>

            {/* Bill Details */}
            <div className="bg-white mx-2 mt-3 rounded-xl p-4 shadow-sm">
                <div className="flex items-center gap-2 mb-3">
                    <IndianRupee className="w-5 h-5 text-[var(--primary-color)]" />
                    <h3 className="font-medium">Bill Details</h3>
                </div>
                <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                        <span className="text-gray-600">Item Total</span>
                        <span>₹{order?.order_items?.reduce((acc, item) => acc + item?.sub_total, 0)}</span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-gray-600">Delivery Fee</span>
                        <span>₹{order?.shipping_fee}</span>
                    </div>

                    <div className="flex justify-between pt-2 border-t font-medium">
                        <span>Total Amount</span>
                        <span>₹{order?.final_order_value}</span>
                    </div>
                </div>
            </div>

            {/* Order Info */}
            <div className="bg-white mx-2 mt-3 mb-4 rounded-xl p-4 shadow-sm">
                <h3 className="font-medium mb-3">Order Information</h3>
                <div className="space-y-3 text-sm">
                    <div>
                        <p className="text-gray-600">Order ID</p>
                        <p className="font-medium">{order?.order_id}</p>
                    </div>
                    <div>
                        <p className="text-gray-600">Payment Method</p>
                        <p className="font-medium">
                            {order?.payment_transaction?.payment_method || 'Online Payment'}
                        </p>
                    </div>
                    {order?.delivered_at && (
                        <div>
                            <p className="text-gray-600">Delivered On</p>
                            <p className="font-medium">{formatDate(order?.delivered_at)}</p>
                        </div>
                    )}
                </div>
            </div>
            {/* Maintaining order status */}
            <div className="bg-white mx-2 mt-3 mb-4 rounded-xl p-4 shadow-sm">
                <StatusTimeLine activeOrder={transformOrder(order)} ORDER_STATUSES={ORDER_STATUSES} />

                {/* Confirmation Button - Only show when status is reached_door */}
                {order?.order_status === 'reached_door' && (
                    <div className="mt-4 space-y-3">
                        <div className="bg-yellow-50 rounded-lg p-3 flex items-start gap-2">
                            <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                            <p className="text-sm text-yellow-700">
                                Please confirm only after receiving your order from the delivery partner.
                            </p>
                        </div>
                        <button
                            onClick={() => setIsConfirmationPopupOpen(true)}
                            disabled={isConfirming}
                            className={`w-full py-3 px-4 rounded-xl font-medium 
                                ${isConfirming
                                    ? 'bg-gray-100 text-gray-400'
                                    : 'bg-[var(--primary-color)] text-white hover:bg-[var(--primary-color)]/90'
                                } transition-all duration-200 flex items-center justify-center gap-2`}
                        >
                            {isConfirming ? (
                                <>
                                    <RefreshCcw className="w-5 h-5 animate-spin" />
                                    <span>Confirming...</span>
                                </>
                            ) : (
                                <>
                                    <Check className="w-5 h-5" />
                                    <span>Confirm Order Received</span>
                                </>
                            )}
                        </button>
                    </div>
                )}
            </div>

            {/* Confirmation Popup */}
            <ConfirmationPopup
                variant="info"
                title="Confirm Order Delivery"
                isOpen={isConfirmationPopupOpen}
                message="Have you received your order from the delivery partner?"
                onConfirm={handleOrderDelivered}
                onCancel={() => setIsConfirmationPopupOpen(false)}
                confirmText="Yes, I've Received It"
                cancelText="No, Not Yet"
            />
        </div>
    )
}

export default OrderDetailsPage
