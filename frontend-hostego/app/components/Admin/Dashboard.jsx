'use client'
import React, { useState, useEffect } from 'react'
import { BarChart3, DollarSign, ShoppingBag, TrendingUp, Package, Calendar, RefreshCw } from 'lucide-react'
import axiosClient from "@/app/utils/axiosClient"

const StatCard = ({ title, value, icon: Icon, trend }) => (
    <div className="bg-white p-6 rounded-xl shadow-sm">
        <div className="flex items-center justify-between mb-4">
            <div>
                <p className="text-sm text-gray-500">{title}</p>
                <h3 className="text-2xl font-bold">₹{value}</h3>
            </div>
            <div className={`p-3 rounded-full ${trend > 0 ? 'bg-green-100' : 'bg-blue-100'}`}>
                <Icon className={`w-6 h-6 ${trend > 0 ? 'text-green-600' : 'text-blue-600'}`} />
            </div>
        </div>
        <div className="flex items-center gap-2">
            <TrendingUp className={`w-4 h-4 ${trend > 0 ? 'text-green-500' : 'text-gray-400'}`} />
            <span className="text-sm text-gray-600">vs. last month</span>
        </div>
    </div>
)

const ProductCard = ({ product }) => {

    return <div className="bg-white p-4 rounded-lg shadow-sm hover:shadow-md transition-all">
        <div className="flex items-center gap-4">
            <img
                src={product?.product_img_url}
                alt={product?.product_name}
                className="w-16 h-16 rounded-lg object-cover"
            />
            <div className="flex-1">
                <h4 className="font-medium text-gray-900">{product?.product_name}</h4>

                <div className="flex items-center justify-between mt-2">
                    <span className="text-sm font-medium text-gray-900">₹{product?.current_price}</span>
                    <span className={`px-2 py-1 rounded-full text-xs ${product?.availability ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                        }`}>
                        {product.availability ? 'In Stock' : 'Out of Stock'}
                    </span>
                </div>
                <p className="text-semibold font-semibold text-gray-800 mt-2 ">  ₹{product?.last_week_revenue} </p>
                <p className="text-sm font-semibold text-gray-600 ">{product?.total_quantity} X ₹{product?.current_price} </p>
            </div>
        </div>
    </div>
}

const Dashboard = ({ dashboardStats }) => {
    
    const [currentData, setCurrentData] = useState(dashboardStats)
    const [isLoading, setIsLoading] = useState(false)
    const [dateRange, setDateRange] = useState({
        startDate: '',
        endDate: ''
    })
    useEffect(() => {
        fetchDashboardData()
    }, [])
    const fetchDashboardData = async (filters = {}) => {
        try {
            setIsLoading(true)
            const { startDate, endDate } = filters
            let url = '/api/order/order-items'

            if (startDate && endDate) {
                url += `?start_date=${startDate}&end_date=${endDate}`
            }

            const { data } = await axiosClient.get(url)
            setCurrentData(data)
        } catch (error) {
            console.error('Error fetching dashboard data:', error)
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
        if (newDateRange.startDate && newDateRange.endDate) {
            await fetchDashboardData(newDateRange)
        }
    }

    const handleRefresh = () => {
        fetchDashboardData(dateRange)
    }

    const {
        overall_stats,
        product_stats
    } = currentData

    // Sort products by order count
    const topProducts = [...product_stats]
        .sort((a, b) => b.order_count - a.order_count)
        .slice(0, 5)

    return (
        <div className="p-6 max-w-7xl mx-auto">
            <div className="flex items-center justify-between mb-6">
                <h1 className="text-2xl font-bold">Dashboard Overview</h1>

                <div className="flex items-center gap-4">
                    {/* Date Filters */}
                    <div className="flex items-center gap-4">
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

                    {/* Refresh Button */}
                    <button
                        onClick={handleRefresh}
                        disabled={isLoading}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)]/10 
                                 text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/20 
                                 transition-all duration-200 disabled:opacity-50"
                    >
                        <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
                        {isLoading ? 'Refreshing...' : 'Refresh'}
                    </button>
                </div>
            </div>

            {/* Loading Overlay */}
            {isLoading && (
                <div className="fixed inset-0 bg-white/50 flex items-center justify-center z-50">
                    <div className="flex flex-col items-center gap-4">
                        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[var(--primary-color)]"></div>
                        <p className="text-gray-500">Loading dashboard data...</p>
                    </div>
                </div>
            )}

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                <StatCard
                    title="Total Revenue"
                    value={overall_stats?.total_revenue}
                    icon={DollarSign}
                    trend={overall_stats?.last_month_revenue}
                />
                <StatCard
                    title="Total Orders"
                    value={overall_stats?.total_orders}
                    icon={ShoppingBag}
                    trend={overall_stats?.last_month_orders}
                />
                <StatCard
                    title="Last Week Revenue"
                    value={overall_stats?.last_week_revenue}
                    icon={BarChart3}
                    trend={overall_stats?.last_week_revenue - overall_stats?.last_month_revenue}
                />
                <StatCard
                    title="Last Month Revenue"
                    value={overall_stats?.last_month_revenue}
                    icon={Package}
                    trend={overall_stats?.last_month_revenue}
                />
            </div>

            {/* Top Products */}
            <div className="mb-8">
                <h2 className="text-xl font-semibold mb-4">Top Performing Products</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {topProducts.map(product => (
                        <ProductCard key={product?.product_id} product={product} />
                    ))}
                </div>
            </div>

            {/* Product Inventory */}
            <div>
                <h2 className="text-xl font-semibold mb-4">Product Inventory</h2>
                <div className="bg-white rounded-xl shadow-sm overflow-hidden">
                    <div className="overflow-x-auto overflow-y-auto max-h-[600px]">
                        <table className="w-full text-sm text-left">
                            <thead className="bg-gray-50 sticky top-0">
                                <tr>
                                    <th className="px-6 py-3">Product Name</th>
                                    <th className="px-6 py-3">Price</th>
                                    <th className="px-6 py-3">Orders</th>
                                    <th className="px-6 py-3">Revenue</th>
                                    <th className="px-6 py-3">Stock</th>
                                    <th className="px-6 py-3">Status</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 ">
                                {product_stats?.map(product => (
                                    <tr key={product.product_id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 font-medium">{product.product_name}</td>
                                        <td className="px-6 py-4">₹{product.current_price}</td>
                                        <td className="px-6 py-4">{product.order_count}</td>
                                        <td className="px-6 py-4">₹{product.total_revenue}</td>
                                        <td className="px-6 py-4">{product.stock_quantity}</td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2 py-1 rounded-full text-xs ${product.availability ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                                                }`}>
                                                {product?.availability ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Dashboard
