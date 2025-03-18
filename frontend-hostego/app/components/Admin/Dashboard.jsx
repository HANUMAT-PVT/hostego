'use client'
import React from 'react'
import { BarChart3, DollarSign, ShoppingBag, TrendingUp, Package } from 'lucide-react'

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
    console.log(product)
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
                <p className="text-medium font-bold text-gray-800 mt-2 ">Revenue: ₹{product?.last_week_revenue}</p>
            </div>
        </div>
    </div>
}

const Dashboard = ({ data }) => {
    const {
        overall_stats,
        product_stats
    } = data

    // Sort products by order count
    const topProducts = [...product_stats]
        .sort((a, b) => b.order_count - a.order_count)
        .slice(0, 5)

    return (
        <div className="p-6 max-w-7xl mx-auto">
            <h1 className="text-2xl font-bold mb-6">Dashboard Overview</h1>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                <StatCard
                    title="Total Revenue"
                    value={overall_stats.total_revenue}
                    icon={DollarSign}
                    trend={overall_stats.last_month_revenue}
                />
                <StatCard
                    title="Total Orders"
                    value={overall_stats.total_orders}
                    icon={ShoppingBag}
                    trend={overall_stats.last_month_orders}
                />
                <StatCard
                    title="Last Week Revenue"
                    value={overall_stats.last_week_revenue}
                    icon={BarChart3}
                    trend={overall_stats.last_week_revenue - overall_stats.last_month_revenue}
                />
                <StatCard
                    title="Last Month Revenue"
                    value={overall_stats.last_month_revenue}
                    icon={Package}
                    trend={overall_stats.last_month_revenue}
                />
            </div>

            {/* Top Products */}
            <div className="mb-8">
                <h2 className="text-xl font-semibold mb-4">Top Performing Products</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {topProducts.map(product => (
                        <ProductCard key={product.product_id} product={product} />
                    ))}
                </div>
            </div>

            {/* Product Inventory */}
            <div>
                <h2 className="text-xl font-semibold mb-4">Product Inventory</h2>
                <div className="bg-white rounded-xl shadow-sm overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm text-left">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3">Product Name</th>
                                    <th className="px-6 py-3">Price</th>
                                    <th className="px-6 py-3">Orders</th>
                                    <th className="px-6 py-3">Revenue</th>
                                    <th className="px-6 py-3">Stock</th>
                                    <th className="px-6 py-3">Status</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100">
                                {product_stats.map(product => (
                                    <tr key={product.product_id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 font-medium">{product.product_name}</td>
                                        <td className="px-6 py-4">₹{product.current_price}</td>
                                        <td className="px-6 py-4">{product.order_count}</td>
                                        <td className="px-6 py-4">₹{product.total_revenue}</td>
                                        <td className="px-6 py-4">{product.stock_quantity}</td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2 py-1 rounded-full text-xs ${product.availability ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                                                }`}>
                                                {product.availability ? 'Active' : 'Inactive'}
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
