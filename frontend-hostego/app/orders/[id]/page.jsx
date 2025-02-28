'use client'
import React, { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import BackNavigationButton from '@/app/components/BackNavigationButton'
import { formatDate } from '@/app/utils/helper'
import { Package, MapPin, Clock, CheckCircle2, Truck, Check, AlertCircle, IndianRupee } from 'lucide-react'
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
            label: 'Order Pending'
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

        console.log("fetching order")
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
                <div className={`p-4 ${status.bgColor} border-b`}>
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
                    <div className="flex justify-between">
                        <span className="text-gray-600">Platform Fee</span>
                        <span>₹{order?.platform_fee}</span>
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
