'use client'
import React, { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import axiosClient from '@/app/utils/axiosClient'
import { Calendar } from 'lucide-react'
import OrderCard from "@/app/components/Orders/OrderCard"
import HostegoLoader from '@/app/components/HostegoLoader'

const DeliveryPartnerEarnings = () => {
    const { id } = useParams()
    const [earnings, setEarnings] = useState(null)
    const [isLoading, setIsLoading] = useState(true)
    const [dateRange, setDateRange] = useState({
        startDate: '',
        endDate: ''
    })

    useEffect(() => {
        fetchEarnings()
    }, [id, dateRange])

    const fetchEarnings = async () => {
        try {
            setIsLoading(true)
            const { startDate, endDate } = dateRange
            let url = `/api/delivery-partner/earnings/${id}`

            if (startDate && endDate) {
                url += `?start_date=${startDate}&end_date=${endDate}`
            }

            const { data } = await axiosClient.get(url)
            setEarnings(data)
        } catch (error) {
            console.error('Error fetching earnings:', error)
        } finally {
            setIsLoading(false)
        }
    }
    

    const handleDateChange = async (e) => {
        const { name, value } = e.target
        const newDateRange = {
            ...dateRange,
            [name]: value
        }

        setDateRange(newDateRange)

        // Only fetch if both dates are selected
        await fetchEarnings()
    }

    if (isLoading) {
        return <HostegoLoader />
    }

    return (
        <div className="max-w-4xl mx-auto p-6">
            {/* Header with Stats */}
            <div className="mb-8">
                <h1 className="text-2xl font-bold mb-6">Delivery Partner Earnings</h1>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                    <div className="bg-white p-4 rounded-xl shadow-sm">
                        <p className="text-sm text-gray-500">Total Earnings</p>
                        <h3 className="text-2xl font-bold">₹{earnings?.summary?.total_earnings || 0}</h3>
                    </div>
                    <div className="bg-white p-4 rounded-xl shadow-sm">
                        <p className="text-sm text-gray-500">Total Deliveries</p>
                        <h3 className="text-2xl font-bold">{earnings?.summary?.total_orders || 0}</h3>
                    </div>
                    <div className="bg-white p-4 rounded-xl shadow-sm">
                        <p className="text-sm text-gray-500">Average Per Delivery</p>
                        <h3 className="text-2xl font-bold">
                            ₹{earnings?.summary?.total_orders ? (earnings?.summary?.total_earnings / earnings?.summary?.total_orders).toFixed(2) : 0}
                        </h3>
                    </div>
                </div>

                {/* Date Filters */}
                <div className="flex items-center gap-4 mb-6">
                    <div className="relative">
                        <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <input
                            type="date"
                            name="startDate"
                            value={dateRange.startDate}
                            onChange={handleDateChange}
                            className="pl-10 pr-4 py-2 rounded-lg border border-gray-200 text-sm
                                     focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 
                                     focus:border-[var(--primary-color)]"
                        />
                    </div>
                    <div className="relative">
                        <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <input
                            type="date"
                            name="endDate"
                            value={dateRange.endDate}
                            onChange={handleDateChange}
                            min={dateRange.startDate}
                            className="pl-10 pr-4 py-2 rounded-lg border border-gray-200 text-sm
                                     focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 
                                     focus:border-[var(--primary-color)]"
                        />
                    </div>
                </div>
            </div>

            {/* Delivery Orders List */}
            <div className="space-y-4">
                <h2 className="text-xl font-semibold mb-4">Delivery Orders</h2>
                {earnings?.daily_earnings?.length > 0 ? (
                    earnings?.daily_earnings.map((order) => (
                        <OrderCard
                       
                            key={order.order_id}
                            order={order}
                            showDeliveryPartner={false}
                        />
                    ))
                ) : (
                    <div className="text-center py-8 bg-gray-50 rounded-lg">
                        <p className="text-gray-500">No orders found for the selected period</p>
                    </div>
                )}
            </div>
        </div>
    )
}

export default DeliveryPartnerEarnings 