"use client"
import React from 'react'
import { formatDate } from '@/app/utils/helper'
import { Package, Clock, CheckCircle2, Truck, AlertCircle, UserRoundCheck, Check } from 'lucide-react'
import { useRouter } from 'next/navigation'


const OrderPreviewCard = ({ order }) => {

    const router = useRouter()
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
            icon: UserRoundCheck,
            color: 'text-blue-500',
            bgColor: 'bg-blue-50',
            label: 'Order Assigned'
        },
        on_the_way: {
            icon: Package,
            color: 'text-[var(--primary-color)]',
            bgColor: 'bg-[var(--primary-color)]/10',
            label: 'On the way'
        },
        picked: {
            icon: Truck,
            color: 'text-purple-500',
            bgColor: 'bg-purple-50',
            label: 'Order Picked'
        },
        reached: {
            icon: Check,
            color: 'text-purple-500',
            bgColor: 'bg-purple-50',
            label: 'Reached Shop'
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

    const status = statusConfig[order.order_status] || statusConfig.pending
    const StatusIcon = status.icon


    return (
        <div onClick={() => { router.push(`/orders/${order?.order_id}`) }} className="bg-white rounded-xl overflow-hidden shadow-sm m-2 border border-gray-100">
            {/* Order Status Header */}
            <div className={`p-4 ${status?.bgColor} border-b flex items-center justify-between`}>
                <div className="flex items-center gap-2">
                    <StatusIcon className={`w-5 h-5 ${status?.color}`} />
                    <span className={`font-medium text-sm ${status?.color}`}>{status?.label}</span>
                </div>
                <span className="text-sm text-gray-600">
                    {formatDate(order?.created_at)}
                </span>
            </div>

            {/* Order Items */}
            <div className="p-4">
                <div className="flex flex-wrap gap-2 mb-4">
                    {order?.order_items?.map((item) => (
                        <div key={item.cart_item_id} className="flex items-center gap-2 bg-gray-50 rounded-lg p-2">
                            <img
                                src={item?.product_item?.product_img_url}
                                alt={item?.product_item?.product_name}
                                className="w-12 h-12 rounded-md object-cover"
                            />
                            {/* <div>
                                <p className="text-sm font-medium line-clamp-1">{item.product_item.product_name}</p>
                                <div className="flex items-center gap-2 text-sm text-gray-600">
                                    <span>₹{item.product_item.food_price}</span>
                                    <span>×</span>
                                    <span>{item.quantity}</span>
                                </div>
                            </div> */}
                        </div>
                    ))}
                </div>

                {/* Order Summary */}
                <div className="space-y-2 text-sm border-t pt-4">
                    <div className="flex justify-between">
                        <span className="text-gray-600">Items Total</span>
                        <span>₹{order?.order_items?.reduce((acc, item) => acc + item?.sub_total, 0)}</span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-gray-600">Delivery Fee</span>
                        <span>₹{order?.shipping_fee}</span>
                    </div>

                    <div className="flex justify-between pt-2 border-t font-semibold text-lg">
                        <span>Total</span>
                        <span>₹{order?.final_order_value}</span>
                    </div>
                </div>
            </div>

            {/* Order Actions */}
            <div className="px-4 py-3 bg-gray-50 flex justify-between items-center">
                <div className="text-sm">
                    <span className="text-gray-600">Order ID: </span>
                    <span className="font-medium">{order?.order_id?.slice(0, 8)}</span>
                </div>
                <button
                    className="text-[var(--primary-color)] text-sm font-medium hover:underline"
                    onClick={() => { router.push(`/orders/${order?.order_id}`) }}
                >
                    View Details
                </button>
            </div>
        </div>
    )
}

export default OrderPreviewCard
