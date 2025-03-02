'use client'
import React, { useState, useEffect, useCallback } from 'react'
import { Package, User, MapPin, IndianRupee, Clock, ChevronDown, ChevronUp, Phone, CheckCircle2, AlertCircle, Search, Filter, RefreshCw } from 'lucide-react'
import { formatDate } from '@/app/utils/helper'
import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '../HostegoLoader'
import debounce from 'lodash/debounce'

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
    const [selectedStatus, setSelectedStatus] = useState(order.order_status)

    const handleStatusUpdate = async (newStatus) => {
        try {
            if(newStatus=="cancelled"){
                setIsUpdating(true)
                await axiosClient.post(`/api/payment/refund`, {
                    order_id: order?.order_id
                })
                onRefresh(true) // Refresh the list after update
            }
            setIsUpdating(true)
            await axiosClient.patch(`/api/order/${order.order_id}`, {
                order_status: newStatus,
                delivery_partner_id: "" // Reset delivery partner ID
            })
            onRefresh(true) // Refresh the list after update
        } catch (error) {
            console.error('Error updating order status:', error)
        } finally {
            setIsUpdating(false)
        }
    }

    return (
        <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100 transition-all duration-200 hover:shadow-md">
            {/* Order Header */}
            <div className="p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <Package className="w-5 h-5 text-white" />
                        <span className="text-white font-medium">
                            #{order.order_id.slice(0, 8)}
                        </span>
                    </div>
                    <div className="flex items-center gap-2">
                        <OrderStatusBadge status={order.order_status} />
                        <select
                            value={selectedStatus}
                            onChange={(e) => {
                                setSelectedStatus(e.target.value)
                                handleStatusUpdate(e.target.value)
                            }}
                            disabled={isUpdating}
                            className="ml-2 px-3 py-1.5 rounded-lg border-2 border-white/20 
                                     bg-white/10 text-white text-sm font-medium
                                     focus:outline-none focus:border-white/40
                                     disabled:opacity-50"
                        >
                            <option value="pending">Pending</option>
                            <option value="placed">Placed</option>
                            <option value="assigned">Assigned</option>
                            <option value="picked">Picked</option>
                            <option value="reached">Reached </option>
                            <option value="on_the_way">On The Way</option>
                            <option value="delivered">Delivered</option>
                            <option value="cancelled">Cancelled</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Order Content */}
            <div className="p-4">
                {/* Order Summary */}
                <div className="grid grid-cols-2 gap-4 mb-4">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-[var(--primary-color)]/10 flex items-center justify-center">
                            <IndianRupee className="w-5 h-5 text-[var(--primary-color)]" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Order Value</p>
                            <p className="text-lg font-semibold">₹{order.final_order_value}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center">
                            <Clock className="w-5 h-5 text-blue-600" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Order Date</p>
                            <p className="font-medium">{formatDate(order.created_at)}</p>
                        </div>
                    </div>
                </div>

                {/* Customer & Delivery Info */}
                <div className="space-y-4">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-green-50 flex items-center justify-center">
                            <User className="w-5 h-5 text-green-600" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Customer</p>
                            <p className="font-medium">{order.user.mobile_number}</p>
                        </div>
                    </div>
                    <div className="flex items-start gap-3">
                        <div className="w-10 h-10 rounded-full bg-orange-50 flex items-center justify-center">
                            <MapPin className="w-5 h-5 text-orange-600" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Delivery Address</p>
                            <p className="font-medium">{order.address.address_line_1 || 'Address not available'}</p>
                        </div>
                    </div>
                </div>

                {/* Expand/Collapse Button */}
                <button
                    onClick={() => setIsExpanded(!isExpanded)}
                    className="w-full mt-4 flex items-center justify-center gap-2 py-2 text-[var(--primary-color)] hover:bg-[var(--primary-color)]/5 rounded-lg transition-colors"
                >
                    {isExpanded ? (
                        <>
                            <ChevronUp className="w-4 h-4" />
                            Show Less
                        </>
                    ) : (
                        <>
                            <ChevronDown className="w-4 h-4" />
                            View Items
                        </>
                    )}
                </button>

                {/* Order Items */}
                {isExpanded && (
                    <div className="mt-4 pt-4 border-t space-y-3">
                        {order.order_items.map((item, index) => (
                            <div key={index} className="flex items-center gap-3 bg-gray-50 p-3 rounded-lg">
                                <img
                                    src={item.product_item.product_img_url}
                                    alt={item.product_item.product_name}
                                    className="w-12 h-12 rounded-lg object-cover"
                                />
                                <div className="flex-1">
                                    <p className="font-medium">{item.product_item.product_name}</p>
                                    <div className="flex items-center gap-2 text-sm text-gray-600">
                                        <span>{item.quantity} × ₹{item.product_item.food_price}</span>
                                        <span className="text-[var(--primary-color)]">₹{item.sub_total}</span>
                                    </div>
                                </div>
                            </div>
                        ))}

                        {/* Order Summary */}
                        <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                            <div className="flex justify-between text-sm">
                                <span className="text-gray-600">Items Total</span>
                                <span>₹{order.order_items.reduce((acc, item) => acc + item.sub_total, 0)}</span>
                            </div>
                            
                            <div className="flex justify-between text-sm">
                                <span className="text-gray-600">Delivery Fee</span>
                                <span>₹{order.shipping_fee}</span>
                            </div>
                            <div className="flex justify-between font-medium pt-2 border-t">
                                <span>Total Amount</span>
                                <span>₹{order.final_order_value}</span>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}

const OrdersList = () => {
    const [orders, setOrders] = useState([])
    const [isLoading, setIsLoading] = useState(true)
    const [isRefreshing, setIsRefreshing] = useState(false)
    const [searchTerm, setSearchTerm] = useState('')
    const [statusFilter, setStatusFilter] = useState('placed')
    const [debouncedSearchTerm, setDebouncedSearchTerm] = useState('')

    const debouncedSearch = useCallback(
        debounce((searchValue) => {
            setDebouncedSearchTerm(searchValue)
        }, 500),
        []
    )

    const handleSearchChange = (e) => {
        const value = e.target.value
        setSearchTerm(value)
        debouncedSearch(value)
    }

    const fetchOrders = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
            const { data } = await axiosClient.get(`/api/order/all?filter=${statusFilter}&search=${debouncedSearchTerm}`)
            setOrders(data || [])
        } catch (error) {
            console.error('Error fetching orders:', error)
        } finally {
            setIsRefreshing(false)
            setIsLoading(false)
        }
    }

    useEffect(() => {
        fetchOrders()
    }, [statusFilter, debouncedSearchTerm])

    useEffect(() => {
        return () => {
            debouncedSearch.cancel()
        }
    }, [debouncedSearch])

    if (isLoading) {
        return <HostegoLoader />
    }

    const filteredOrders = orders
        .filter(order => {
            const matchesSearch = order?.order_id.toLowerCase()?.includes(debouncedSearchTerm.toLowerCase()) ||
                order.user.mobile_number.includes(debouncedSearchTerm)
            const matchesStatus = statusFilter === 'all' || order.order_status === statusFilter
            return matchesSearch && matchesStatus
        })
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))

    return (
        <div className="max-w-4xl mx-auto p-4">
            {/* Header */}
            <div className="mb-6">
                <div className="flex items-center justify-between mb-4">
                    <h1 className="text-2xl font-semibold">Orders Management</h1>
                    <button
                        onClick={() => fetchOrders(true)}
                        disabled={isRefreshing}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)]/10 
                                 text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/20 
                                 transition-all duration-200 disabled:opacity-50"
                    >
                        <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
                        {isRefreshing ? 'Refreshing...' : 'Refresh'}
                    </button>
                </div>

                {/* Search and Filter */}
                <div className="flex flex-col sm:flex-row gap-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Search by order ID or phone number"
                            value={searchTerm}
                            onChange={handleSearchChange}
                            className="w-full pl-10 pr-4 py-2 rounded-lg border-2 border-gray-100 focus:border-[var(--primary-color)] outline-none"
                        />
                    </div>
                    <select
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value)}
                        className="px-4 py-2 rounded-lg border-2 border-gray-100 focus:border-[var(--primary-color)] outline-none"
                    >
                        <option value="all">All Orders</option>
                        <option value="pending">Pending</option>
                        <option value="placed">Placed</option>
                        <option value="assigned">Assigned</option>
                        <option value="preparing">Preparing</option>
                        <option value="out_for_delivery">Out for Delivery</option>
                        <option value="delivered">Delivered</option>
                        <option value="cancelled">Cancelled</option>
                    </select>
                </div>
            </div>

            {/* Orders List */}
            <div className="space-y-4">
                {filteredOrders.length === 0 ? (
                    <div className="text-center py-12 bg-white rounded-xl">
                        <Package className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                        <h3 className="text-lg font-semibold text-gray-800 mb-2">No Orders Found</h3>
                        <p className="text-gray-600">
                            {searchTerm || statusFilter !== 'all'
                                ? 'Try adjusting your filters'
                                : 'There are no orders to display'}
                        </p>
                    </div>
                ) : (
                    filteredOrders.map(order => (
                        <OrderCard
                            key={order.order_id}
                            order={order}
                            onRefresh={fetchOrders}
                        />
                    ))
                )}
            </div>
        </div>
    )
}

export default OrdersList
