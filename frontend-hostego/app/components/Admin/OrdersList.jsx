'use client'
import React, { useState, useEffect, useCallback } from 'react'
import { Package, Search, RefreshCw } from 'lucide-react'

import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '../HostegoLoader'
import debounce from 'lodash/debounce'

import LoadMoreData from "../LoadMoreData"
import OrderCard from '../Orders/OrderCard'



const OrdersList = () => {
    const [orders, setOrders] = useState([])
    const [isLoading, setIsLoading] = useState(true)
    const [isRefreshing, setIsRefreshing] = useState(false)
    const [searchTerm, setSearchTerm] = useState('')
    const [statusFilter, setStatusFilter] = useState('placed')
    const [debouncedSearchTerm, setDebouncedSearchTerm] = useState('')
    const [hasMore, setHasMore] = useState(true)
    const [page, setPage] = useState(1)

    const debouncedSearch = useCallback(
        debounce((searchValue) => {
            setDebouncedSearchTerm(searchValue)
        }, 500),
        []
    )

    useEffect(() => {
        fetchOrders()
    }, [statusFilter, debouncedSearchTerm])

    useEffect(() => {
        return () => {
            debouncedSearch.cancel()
        }
    }, [debouncedSearch])


    const handleSearchChange = (e) => {
        const value = e.target.value
        setSearchTerm(value)
        debouncedSearch(value)
    }

    const fetchOrders = async (showRefreshAnimation = false, gotNewPage = false, newPage) => {

        try {

            if (!newPage) {
                newPage = page;
            }
            if (showRefreshAnimation) {
                setIsRefreshing(true)
                newPage = 1;
            } else {
                setIsLoading(true)
            }
            const { data } = await axiosClient.get(`/api/order/all?filter=${statusFilter}&search=${debouncedSearchTerm}&page=${newPage}&limit=20`)

            if (gotNewPage) {
                setOrders([...orders, ...data])
            } else {
                setOrders(data)

            }

            setHasMore(data?.length < 20 ? false : true)
        } catch (error) {

        } finally {
            setIsRefreshing(false)
            setIsLoading(false)
        }
    }

    const handlePaginationOrders = () => {
        setPage(page + 1);
        fetchOrders(false, true, page + 1)

    }




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
                        <option value="">All Orders</option>
                        <option value="pending">Pending</option>
                        <option value="placed">Placed</option>
                        <option value="assigned">Assigned</option>
                        <option value="reached">Reached Shop</option>
                        <option value="on_the_way">On The Way</option>
                        <option value="reached_door">Reached Door</option>
                        <option value="delivered">Delivered</option>
                        <option value="cancelled">Cancelled</option>
                    </select>
                </div>
            </div>
            {isLoading && <HostegoLoader />}

            {/* Orders List */}
            <div className="space-y-4 overflow-y-auto max-h-[85vh]">
                {orders?.length > 0 ? orders?.map(order => (
                    <OrderCard
                        key={order.order_id}
                        order={order}
                        onRefresh={fetchOrders}
                    />
                )) : <div className="text-center py-12 bg-white rounded-xl">
                    <Package className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                    <h3 className="text-lg font-semibold text-gray-800 mb-2">No Orders Found</h3>
                    <p className="text-gray-600">
                        {searchTerm
                            ? 'Try adjusting your filters'
                            : 'There are no orders to display'}
                    </p>

                </div>}
            </div>
            {hasMore && <LoadMoreData loadMore={handlePaginationOrders} isLoading={isLoading} />}
        </div>
    )
}

export default OrdersList
