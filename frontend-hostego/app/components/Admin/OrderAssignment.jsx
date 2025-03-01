"use client";
import React, { useState, useEffect } from 'react'
import { Package, User, MapPin, IndianRupee, Clock, ChevronDown, ChevronUp, Phone, CheckCircle, RefreshCw } from 'lucide-react'

import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '../HostegoLoader'
import HostegoButton from '../HostegoButton'
import OrderAssignCard from './OrderAssignCard'
import DeliveryPartnerCard from './DeliveryPartnerCard';


const OrderAssignment = () => {
    const [orders, setOrders] = useState([])
    const [deliveryPartners, setDeliveryPartners] = useState([])
    const [selectedPartner, setSelectedPartner] = useState(null)
    const [isLoading, setIsLoading] = useState(true)
    const [selectedOrderItem, setSelectedOrderItem] = useState({})
    const [isRefreshingOrders, setIsRefreshingOrders] = useState(false)
    const [isRefreshingPartners, setIsRefreshingPartners] = useState(false)


    useEffect(() => {
        fetchData()
    }, [])

    const fetchOrders = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshingOrders(true) : setIsLoading(true)
            const { data } = await axiosClient.get('/api/order/all?filter=placed')
            setOrders(data)
        } catch (error) {
            console.error('Error fetching orders:', error)
        } finally {
            setIsRefreshingOrders(false)
            setIsLoading(false)
        }
    }

    const fetchPartners = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshingPartners(true) : setIsLoading(true)
            const { data } = await axiosClient.get('/api/delivery-partner/all?availability=1')
            setDeliveryPartners(data)
        } catch (error) {
            console.error('Error fetching partners:', error)
        } finally {
            setIsRefreshingPartners(false)
            setIsLoading(false)
        }
    }

    const fetchData = async () => {
        try {
            setIsLoading(true)
            await Promise.all([fetchOrders(), fetchPartners()])
        } catch (error) {
            console.error('Error fetching data:', error)
        } finally {
            setIsLoading(false)
        }
    }

    const handleAssignOrder = async () => {
        try {
            await axiosClient.post('/api/order/assign-order-delivery', {
                order_id: selectedOrderItem?.order_id,
                delivery_partner_id: selectedPartner?.delivery_partner_id
            })
            // Refresh both lists after assignment
            fetchData()
            // Reset selections
            setSelectedOrderItem({})
            setSelectedPartner(null)
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
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-xl font-semibold">Pending Orders</h2>
                        <button
                            onClick={() => fetchOrders(true)}
                            disabled={isRefreshingOrders}
                            className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)]/10 
                                     text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/20 
                                     transition-all duration-200 disabled:opacity-50"
                        >
                            <RefreshCw className={`w-4 h-4 ${isRefreshingOrders ? 'animate-spin' : ''}`} />
                            {isRefreshingOrders ? 'Refreshing...' : 'Refresh '}
                        </button>
                    </div>

                    <div className="h-[90vh] overflow-y-auto">
                        {orders.length === 0 ? (
                            <div className="text-center py-12 bg-white rounded-xl">
                                <Package className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                                <h3 className="text-lg font-semibold text-gray-800 mb-2">No Pending Orders</h3>
                                <p className="text-gray-600">There are no orders waiting for assignment</p>
                            </div>
                        ) : (
                            orders?.map(order => (
                                <OrderAssignCard
                                    key={order?.order_id}
                                    order={order}
                                    selectedOrderItem={selectedOrderItem}
                                    selectOrderItem={(e) => setSelectedOrderItem(e)}
                                />
                            ))
                        )}
                    </div>
                </div>

                {/* Delivery Partners Section */}
                <div>
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-xl font-semibold">Delivery Partners</h2>
                        <button
                            onClick={() => fetchPartners(true)}
                            disabled={isRefreshingPartners}
                            className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)]/10 
                                     text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/20 
                                     transition-all duration-200 disabled:opacity-50"
                        >
                            <RefreshCw className={`w-4 h-4 ${isRefreshingPartners ? 'animate-spin' : ''}`} />
                            {isRefreshingPartners ? 'Refreshing...' : 'Refresh '}
                        </button>
                    </div>

                    <div className="space-y-3">
                        {deliveryPartners?.map(partner => (
                            <DeliveryPartnerCard
                                key={partner?.user_id}
                                partner={partner}
                                isSelected={selectedPartner?.user_id === partner?.user_id}
                                onSelect={setSelectedPartner}
                            />
                        ))}
                    </div>
                </div>
            </div>

            {/* Assign Order Button */}
            {selectedOrderItem?.order_id && selectedPartner?.delivery_partner_id && (
                <div className="fixed bottom-[40px] w-[300px] left-1/2 -translate-x-1/2">
                    <HostegoButton
                        onClick={handleAssignOrder}
                        text="Assign Order"
                        className="w-[300px] mt-4"
                        icon={<CheckCircle className="w-4 h-4" />}
                    />
                </div>
            )}
        </div>
    )
}

export default OrderAssignment
