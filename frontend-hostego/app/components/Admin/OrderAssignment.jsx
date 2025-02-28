"use client";
import React, { useState, useEffect } from 'react'
import { Package, User, MapPin, IndianRupee, Clock, ChevronDown, ChevronUp, Phone, CheckCircle } from 'lucide-react'

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


    useEffect(() => {
        fetchData()
    }, [])

    const fetchData = async () => {
        try {
            setIsLoading(true)
            const [ordersRes, partnersRes] = await Promise.all([
                axiosClient.get('/api/order/all?filter=placed'),
                axiosClient.get('/api/delivery-partner/all?availability=1')
            ])
            setOrders(ordersRes?.data)
            setDeliveryPartners(partnersRes.data)
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
                <div className="lg:col-span-2 space-y-4 h-[90vh] overflow-y-auto">
                    <h2 className="text-xl font-semibold mb-4">Pending Orders</h2>
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

                {/* Delivery Partners Section */}
                <div>
                    <h2 className="text-xl font-semibold mb-4">Delivery Partners</h2>
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
                <div>
                    {selectedOrderItem?.order_id && selectedPartner?.delivery_partner_id && (
                        <div className="absolute bottom-[40px] w-[300px]  left-1/2 -translate-x-1/2  ">
                            <HostegoButton
                               
                                onClick={handleAssignOrder}
                                text="Assign Order"
                                className="w-[300px] mt-4 "
                                icon={<CheckCircle className="w-4 h-4" />}
                            />
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

export default OrderAssignment
