'use client'
import React, { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import BackNavigationButton from '@/app/components/BackNavigationButton'
import { formatDate } from '@/app/utils/helper'
import { Package, MapPin, Clock, CheckCircle2, Truck, Check, AlertCircle, IndianRupee, RefreshCcw, Phone, Bike, MessageSquare, Award } from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '@/app/components/HostegoLoader'
import StatusTimeLine from '../../components/Orders/StatusTimeLine'
import { transformOrder } from '../../utils/helper'
import { ORDER_STATUSES } from '../../components/Delivery-Partner/MaintainOrderStatusForDeliveryPartner'

const OrderDetailsPage = () => {
    const { id } = useParams()
    const [order, setOrder] = useState(null)
    const [isLoading, setIsLoading] = useState(true)

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

    useEffect(() => {


        fetchOrder()

    }, [id])




    const fetchOrder = async () => {
        try {
            setIsLoading(true)
            const { data } = await axiosClient.get(`/api/order/${id}`)
            setOrder(data)
        } catch (error) {
            console.error('Error fetching order:', error)
        } finally {
            setIsLoading(false)
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
                        <RefreshCcw onClick={() => fetchOrder()} className="w-4 h-4 text-[var(--primary-color)] cursor-pointer" />
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
                                        {item?.product_item?.product_name}
                                    </h4>
                                    <p className="text-sm text-gray-600">
                                        {item?.product_item?.weight}
                                    </p>
                                    <div className="flex items-center gap-2 mt-1">
                                        <span className="text-sm">₹{item?.product_item?.food_price}</span>
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
                <div className="bg-white p-4 rounded-xl shadow-sm mb-4">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                            <div className="w-12 h-12 rounded-full bg-[var(--primary-color)]/10 flex items-center justify-center">
                                <Bike className="w-6 h-6 text-[var(--primary-color)]" />
                            </div>
                            <div>
                                <h3 className="font-medium text-gray-800">
                                    {order?.delivery_partner?.first_name}
                                </h3>
                                <p className="text-sm text-gray-500">Your Delivery Partner</p>
                            </div>
                        </div>
                       
                    </div>

                    <div className="flex gap-3">
                        {/* Call Button */}
                        <button
                            onClick={() => window.location.href = `tel:${order?.delivery_partner?.mobile_number}`}
                            className=" border-2 border-[var(--primary-color)] flex-1 py-3.5 px-4 rounded-xl bg-[var(--primary-color)]/10 hover:bg-[var(--primary-color)]/15 
                                     text-[var(--primary-color)] font-medium flex items-center justify-center gap-2 
                                     transition-all duration-200 group"
                        >
                            <div className="relative">
                                <Phone className="w-5 h-5" />
                                <span className="absolute -top-1 -right-1 w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                            </div>
                            <span>Call Partner</span>
                        </button>

                        {/* Chat Button - Optional */}
                        
                    </div>

                    {/* Delivery Info */}
                    <div className="mt-4 pt-4 border-t border-gray-100">
                        <div className="flex items-center gap-2 text-sm text-gray-600">
                            <Clock className="w-4 h-4" />
                            <span>Estimated delivery in {order?.estimated_delivery_time || '15-20'} minutes</span>
                        </div>
                        
                    </div>

                    {/* Safety Tip */}
                    <div className="mt-4 bg-yellow-50 rounded-lg p-3 flex items-start gap-2">
                        <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                        <p className="text-sm text-yellow-700">
                            For your safety Please don't share any sensitive information.
                        </p>
                    </div>
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
            </div>

        </div>
    )
}

export default OrderDetailsPage
