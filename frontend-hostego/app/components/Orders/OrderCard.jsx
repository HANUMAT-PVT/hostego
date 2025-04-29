"use client"
import { useState } from "react"
import { Package, User, MapPin, IndianRupee, Clock, ChevronDown, ChevronUp, Phone, CheckCircle2, AlertCircle, Search, Filter, RefreshCw } from 'lucide-react'
import { formatDate } from '@/app/utils/helper'
import ConfirmationPopup from '../ConfirmationPopup'
import axiosClient from "@/app/utils/axiosClient"
import { useSelector } from "react-redux"


const OrderStatusBadge = ({ status }) => {

    const statusConfig = {
        pending: {
            color: 'bg-yellow-100 text-yellow-700',
            icon: Clock
        },
        confirmed: {
            color: 'bg-blue-100 text-blue-700',
            icon: CheckCircle2
        },
        preparing: {
            color: 'bg-purple-100 text-purple-700',
            icon: Package
        },
        out_for_delivery: {
            color: 'bg-indigo-100 text-indigo-700',
            icon: MapPin
        },
        delivered: {
            color: 'bg-green-100 text-green-700',
            icon: CheckCircle2
        },
        cancelled: {
            color: 'bg-red-100 text-red-700',
            icon: AlertCircle
        }
    }

    const config = statusConfig[status] || statusConfig.pending
    const Icon = config.icon

    return (
        <span className={`flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium ${config.color}`}>
            <Icon className="w-4 h-4" />
            {status.replace(/_/g, ' ').charAt(0).toUpperCase() + status.slice(1).replace(/_/g, ' ')}
        </span>
    )
}

const OrderCard = ({ order, onRefresh }) => {
    const [isExpanded, setIsExpanded] = useState(false)
    const [isUpdating, setIsUpdating] = useState(false)
    const [selectedStatus, setSelectedStatus] = useState(order?.order_status)
    const [isConfirmationPopupOpen, setIsConfirmationPopupOpen] = useState(false)
    const [isEditable, setIsEditable] = useState({})

    const { userRoles } = useSelector(state => state.user)


    const checkIsSuperAdmin = () => {
        for (let el of userRoles) {
            if (el?.role?.role_id == 1) return true
        }
    }
    const handleStatusUpdate = async (newStatus) => {
        try {
            let delivery_partner_id = 0;
            if (newStatus == "") {

                return
            }
            let cancelNoRefund = false;

            if (newStatus === "cancelled-no-refund") {
                setIsUpdating(true)
                await axiosClient.post(`/api/order/cancel-no-refund`, {
                    order_id: order?.order_id
                })
                cancelNoRefund = true
                onRefresh(true)
            }
            if (newStatus == "cancelled") {
                setIsUpdating(true)
                await axiosClient.post(`/api/payment/refund`, {
                    order_id: order?.order_id
                })
                onRefresh(true) // Refresh the list after update
            }
            setIsUpdating(true)
            if (newStatus == "delivered") {
                delivery_partner_id = order?.delivery_partner?.delivery_partner_id;
            }
            await axiosClient.patch(`/api/order/${order.order_id}`, {
                order_status: cancelNoRefund ? "cancelled" : newStatus,
                delivery_partner_id // Reset delivery partner ID
            })
            setSelectedStatus("")
            onRefresh(true) // Refresh the list after update
        } catch (error) {
            console.error('Error updating order status:', error)
        } finally {
            setIsUpdating(false)
        }
    }

    const handleOrderItemRefund = async () => {
        try {
            let data = await axiosClient.post("/api/order-item/refund", isEditable)
            alert("Refund Proccessed")
            setIsEditable({})
        } catch (error) {
            console.log(error)
        }
    }





    return (
        <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100 transition-all duration-200">
            {/* Order Header - Improved contrast and spacing */}
            <div className="p-4 bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="bg-white/10 p-2 rounded-lg">
                            <Package className="w-5 h-5 text-white" />
                        </div>
                        <div>
                            <span className="text-white/70 text-sm">Order ID</span>
                            <p className="text-white font-medium">#{order?.order_id}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3">
                        <OrderStatusBadge status={order?.order_status} />
                        {order?.order_status !== "delivered" && <select
                            value={selectedStatus}
                            onChange={(e) => {
                                setSelectedStatus(e.target.value)
                                setIsConfirmationPopupOpen(true)
                            }}
                            disabled={isUpdating}
                            className="px-3 py-2 rounded-lg border-2 border-white/20 
                                      text-black text-sm font-medium
                                     focus:outline-none focus:ring-2 focus:ring-white/40
                                     disabled:opacity-50 cursor-pointer"
                        >
                            <option value="">Update Status</option>
                            <option value="placed">Placed</option>
                            <option value="cancelled">Cancelled</option>
                            <option value="cancelled-no-refund">Cancel(No Refund)</option>
                            <option value="delivered">Delivered</option>
                        </select>}
                    </div>
                </div>
            </div>

            <div className="p-4">
                {/* Quick Actions - New section */}
                <div className="flex gap-3 mb-6">
                    <button
                        onClick={() => window.location.href = `tel:${order?.user?.mobile_number}`}
                        className="flex-1 py-3 px-4 rounded-xl bg-[var(--primary-color)]/5 
                                 hover:bg-[var(--primary-color)]/10 text-[var(--primary-color)]
                                 font-medium flex items-center justify-center gap-2 transition-all"
                    >
                        <Phone className="w-4 h-4" />
                        Call Customer
                    </button>
                    <button
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="flex-1 py-3 px-4 rounded-xl bg-gray-50 hover:bg-gray-100
                                 text-gray-700 font-medium flex items-center justify-center gap-2
                                 transition-all"
                    >
                        {isExpanded ? 'Hide Details' : 'View Details'}
                    </button>
                </div>

                {/* Order Summary Cards - Improved layout */}
                <div className="grid grid-cols-2 gap-4 mb-6">
                    <div className="p-4 rounded-xl bg-[var(--primary-color)]/5 space-y-1">
                        <p className="text-sm text-gray-600">Total Amount</p>
                        <p className="text-2xl font-semibold text-[var(--primary-color)]">
                            ₹{order.final_order_value}
                        </p>
                    </div>
                    <div className="p-4 rounded-xl bg-blue-50 space-y-1">
                        <p className="text-sm text-gray-600">Order Placed</p>
                        <p className="text-lg font-medium text-blue-700">
                            {formatDate(order.created_at)}
                        </p>
                    </div>
                </div>

                {/* Add Delivery Partner Info if assigned */}
                {order.delivery_partner && (
                    <div className="bg-blue-50 rounded-xl p-4 mb-6">
                        <div className="flex items-center gap-3 mb-4 pb-4 border-b border-blue-100">
                            <div className="w-12 h-12 rounded-full bg-white flex items-center justify-center shadow-sm">
                                <Package className="w-6 h-6 text-blue-600" />
                            </div>
                            <div>
                                <p className="text-sm text-blue-600">Delivery Partner</p>
                                <p className="font-medium text-lg">{order.delivery_partner.user?.first_name}</p>
                                <div className="flex items-center gap-2 mt-1">
                                    <a
                                        href={`tel:${order.delivery_partner.user?.mobile_number}`}
                                        className="text-sm text-blue-600 flex items-center gap-1 hover:text-blue-700"
                                    >
                                        <Phone className="w-4 h-4" />
                                        {order.delivery_partner.user?.mobile_number}
                                    </a>
                                    <span className="text-sm text-blue-600">
                                        (ID: #{order.delivery_partner.delivery_partner_id})
                                    </span>
                                </div>
                            </div>
                        </div>

                        <div className="flex items-center justify-between text-sm">
                            <span className="text-blue-600">Current Status</span>
                            <OrderStatusBadge status={order.order_status} />
                        </div>
                    </div>
                )}

                {/* Delivery Info - Better organized */}
                <div className="bg-gray-50 rounded-xl p-4 mb-6">
                    <div className="flex items-center gap-3 mb-4 pb-4 border-b border-gray-200">
                        <div className="w-12 h-12 rounded-full bg-white flex items-center justify-center shadow-sm">
                            <User className="w-6 h-6 text-[var(--primary-color)]" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Customer</p>
                            <p className="font-medium text-lg">{order.user.first_name} {order.user.last_name}</p>
                            <p className="text-sm text-gray-600">{order.user.mobile_number}</p>
                        </div>
                    </div>

                    <div className="flex items-start gap-3">
                        <div className="w-12 h-12 rounded-full bg-white flex items-center justify-center shadow-sm">
                            <MapPin className="w-6 h-6 text-orange-600" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Delivery Address</p>
                            <p className="font-medium">{order.address.address_line_1 || 'Address not available'}</p>
                        </div>
                    </div>
                </div>

                {/* Order Items - Collapsible section with improved visibility */}
                {isExpanded && (
                    <div className="space-y-4 animate-fade-in">
                        <h3 className="font-medium text-gray-900 mb-3">Order Items</h3>
                        {order.order_items.map((item, index) => (
                            <div key={index} className="flex items-center gap-4 bg-gray-50 p-3 rounded-xl">
                                <img
                                    src={item.product_item.product_img_url}
                                    alt={item.product_item.product_name}
                                    className="w-16 h-16 rounded-xl object-cover"
                                />

                                <div className="flex-1">
                                    <p className="font-medium text-gray-900">{item?.product_item?.product_name}</p>
                                    <p className="font-medium text-xs text-gray-600">{item?.product_item?.shop?.shop_name}</p>
                                    <div className="flex items-center gap-2 mt-1">
                                        <span eid className="text-sm text-gray-600">
                                            {item?.quantity} × ₹{item?.product_item?.food_price}
                                        </span>
                                        <span className="text-sm font-medium text-[var(--primary-color)]">
                                            ₹{item?.sub_total}
                                        </span>
                                    </div>
                                </div>
                                {checkIsSuperAdmin() && <> <button onClick={() => setIsEditable({ product_id: item?.product_id, order_id: order?.order_id, quantity: item?.quantity })} >Edit</button>
                                    {isEditable?.product_id == item.product_id && <input type="number" value={isEditable?.quantity || 0} onChange={(e) => setIsEditable({ ...isEditable, quantity: Number(e.target.value) })} />}
                                    <button onClick={handleOrderItemRefund}>Submit</button></>}
                            </div>


                        ))}

                        {/* Order Summary - Clear breakdown */}
                        <div className="bg-gray-50 p-4 rounded-xl space-y-3 mt-4">
                            <div className="flex justify-between text-sm">
                                <span className="text-gray-600">Items Total</span>
                                <span className="font-medium">₹{order?.order_items?.reduce((acc, item) => acc + item.sub_total, 0)}</span>
                            </div>
                            <div className="flex justify-between text-sm">
                                <span className="text-gray-600">Delivery Fee</span>
                                <span className="font-medium">₹{order?.shipping_fee}</span>
                            </div>
                            <div className="flex justify-between font-medium pt-3 border-t">
                                <span className="text-gray-900">Total Amount</span>
                                <span className="text-[var(--primary-color)]">₹{order?.final_order_value}</span>
                            </div>
                        </div>
                    </div>
                )}
            </div>

            <ConfirmationPopup
                variant="info"
                title={`Confirm Order ${selectedStatus}`}
                isOpen={isConfirmationPopupOpen}
                message={`Are you sure you want to update the order status to ${selectedStatus}?`}
                onConfirm={() => {
                    handleStatusUpdate(selectedStatus)
                    setIsConfirmationPopupOpen(false)
                }}
                onCancel={() => {
                    setIsConfirmationPopupOpen(false)
                    setSelectedStatus("")
                }}
            />
        </div >
    )
}

export default OrderCard;