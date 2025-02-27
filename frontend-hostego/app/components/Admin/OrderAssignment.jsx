"use client";
import React, { useState, useEffect } from 'react'
import { Package, User, MapPin, IndianRupee, Clock, ChevronDown, ChevronUp, Phone, CheckCircle } from 'lucide-react'
import { formatDate } from '@/app/utils/helper'
import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '../HostegoLoader'
import HostegoButton from '../HostegoButton'

const OrderCard = ({ order, onAssign, selectedDeliveryPartner }) => {
    const [isExpanded, setIsExpanded] = useState(false)

    return (
        <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100">
            {/* Order Header */}
            <div className="p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <Package className="w-5 h-5 text-white" />
                        <span className="text-white font-medium">
                            #{order.order_id.slice(0, 8)}
                        </span>
                    </div>
                    <span className="px-3 py-1 bg-white/20 backdrop-blur-sm text-white rounded-full text-sm">
                        {formatDate(order.created_at)}
                    </span>
                </div>
            </div>

            {/* Order Summary */}
            <div className="p-4">
                <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-[var(--primary-color)]/10 flex items-center justify-center">
                            <IndianRupee className="w-5 h-5 text-[var(--primary-color)]" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Order Value</p>
                            <p className="text-xl font-semibold">₹{order.final_order_value}</p>
                        </div>
                    </div>
                    <button
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-gray-100 text-gray-700 hover:bg-gray-200 transition-colors"
                    >
                        {isExpanded ? 'Show Less' : 'View Items'}
                        {isExpanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
                    </button>
                </div>

                {/* Customer Details */}
                <div className="flex items-center gap-3 mb-4">
                    <div className="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center">
                        <User className="w-5 h-5 text-blue-600" />
                    </div>
                    <div>
                        <p className="text-sm text-gray-600">Customer</p>
                        <p className="font-medium">{order.user.mobile_number}</p>
                    </div>
                </div>

                {/* Delivery Address */}
                <div className="flex items-start gap-3 mb-4">
                    <div className="w-10 h-10 rounded-full bg-green-50 flex items-center justify-center">
                        <MapPin className="w-5 h-5 text-green-600" />
                    </div>
                    <div>
                        <p className="text-sm text-gray-600">Delivery Address</p>
                        <p className="font-medium">{order.address.address_line_1 || 'Address not available'}</p>
                    </div>
                </div>

                {/* Order Items */}
                {isExpanded && (
                    <div className="mt-4 pt-4 border-t space-y-3">
                        <p className="font-medium text-gray-700">Order Items</p>
                        {order.order_items.map((item, index) => (
                            <div key={index} className="flex items-center gap-3 bg-gray-50 p-3 rounded-lg">
                                <img
                                    src={item.product_item.product_img_url}
                                    alt={item.product_item.product_name}
                                    className="w-12 h-12 rounded-lg object-cover"
                                />
                                <div className="flex-1">
                                    <p className="font-medium">{item.product_item.product_name}</p>
                                    <p className="text-sm text-gray-600">
                                        {item.quantity} × ₹{item.product_item.food_price}
                                    </p>
                                </div>
                                <p className="font-medium">₹{item.sub_total}</p>
                            </div>
                        ))}
                    </div>
                )}

                {/* Assign Button */}
                {selectedDeliveryPartner && (
                    <HostegoButton
                        onClick={() => onAssign(order.order_id)}
                        text="Assign Order"
                        className="w-full mt-4"
                        icon={<CheckCircle className="w-4 h-4" />}
                    />
                )}
            </div>
        </div>
    )
}

const DeliveryPartnerCard = ({ partner, isSelected, onSelect }) => (
    <div
        onClick={() => onSelect(partner)}
        className={`p-4 rounded-xl border-2 cursor-pointer transition-all duration-200 
            ${isSelected
                ? 'border-[var(--primary-color)] bg-[var(--primary-color)]/5'
                : 'border-gray-100 bg-white hover:border-[var(--primary-color)]/30'}`}
    >
        <div className="flex items-center gap-3">
            <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center">
                <User className={`w-6 h-6 ${isSelected ? 'text-[var(--primary-color)]' : 'text-gray-500'}`} />
            </div>
            <div>
                <p className="font-medium">{partner.first_name} {partner.last_name}</p>
                <p className="text-sm text-gray-600">{partner.mobile_number}</p>
            </div>
            {isSelected && (
                <div className="ml-auto">
                    <CheckCircle className="w-5 h-5 text-[var(--primary-color)]" />
                </div>
            )}
        </div>
    </div>
)

const OrderAssignment = () => {
    const [orders, setOrders] = useState([])
    const [deliveryPartners, setDeliveryPartners] = useState([])
    const [selectedPartner, setSelectedPartner] = useState(null)
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        fetchData()
    }, [])

    const fetchData = async () => {
        try {
            setIsLoading(true)
            const [ordersRes, partnersRes] = await Promise.all([
                axiosClient.get('/api/order'),
                axiosClient.get('/api/users')
            ])
            setOrders(ordersRes?.data)
            setDeliveryPartners(partnersRes.data)
        } catch (error) {
            console.error('Error fetching data:', error)
        } finally {
            setIsLoading(false)
        }
    }

    const handleAssignOrder = async (orderId) => {
        try {
            await axiosClient.post('/api/assign-order', {
                order_id: orderId,
                delivery_partner_id: selectedPartner.user_id
            })
            // Refresh orders after assignment
            fetchData()
        } catch (error) {
            console.error('Error assigning order:', error)
        }
    }

    if (isLoading) {
        return <HostegoLoader />
    }

    return (
        <div className="max-w-6xl mx-auto p-4">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Orders Section */}
                <div className="lg:col-span-2 space-y-4">
                    <h2 className="text-xl font-semibold mb-4">Pending Orders</h2>
                    {orders.length === 0 ? (
                        <div className="text-center py-12 bg-white rounded-xl">
                            <Package className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                            <h3 className="text-lg font-semibold text-gray-800 mb-2">No Pending Orders</h3>
                            <p className="text-gray-600">There are no orders waiting for assignment</p>
                        </div>
                    ) : (
                        orders.map(order => (
                            <OrderCard
                                key={order.order_id}
                                order={order}
                                onAssign={handleAssignOrder}
                                selectedDeliveryPartner={selectedPartner}
                            />
                        ))
                    )}
                </div>

                {/* Delivery Partners Section */}
                <div>
                    <h2 className="text-xl font-semibold mb-4">Delivery Partners</h2>
                    <div className="space-y-3">
                        {deliveryPartners.map(partner => (
                            <DeliveryPartnerCard
                                key={partner.user_id}
                                partner={partner}
                                isSelected={selectedPartner?.user_id === partner.user_id}
                                onSelect={setSelectedPartner}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </div>
    )
}

export default OrderAssignment
