'use client'
import React, { useState, useEffect } from 'react'
import { Wallet, Search, RefreshCw, CheckCircle2, XCircle, Clock, MoreVertical, User, Plus } from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import HostegoButton from '../HostegoButton'

const StatusBadge = ({ status }) => {
    const configs = {
        pending: { color: 'bg-yellow-100 text-yellow-700', label: 'Pending', icon: Clock },
        success: { color: 'bg-green-100 text-green-700', label: 'Completed', icon: CheckCircle2 },
        failed: { color: 'bg-red-100 text-red-700', label: 'Failed', icon: XCircle }
    }

    const config = configs[status] || configs.pending
    const Icon = config.icon

    return (
        <span className={`px-2 py-1 rounded-full text-xs font-medium flex items-center gap-1 ${config.color}`}>
            <Icon size={12} />
            {config.label}
        </span>
    )
}

const WithdrawalRequestCard = ({ request, onUpdate }) => {
    const [showActions, setShowActions] = useState(false)
    const [uniqueTransactionId, setUniqueTransactionId] = useState('')
    const [isVerifying, setIsVerifying] = useState(false)

    const handleStatusUpdate = async (newStatus) => {
        try {
            if (!uniqueTransactionId && newStatus === 'success') {
                alert('Please enter transaction ID to verify the payment')
                return
            }

            setIsVerifying(true)
            await axiosClient.patch(`/api/delivery-partner-wallet/withdrawal-requests/${request?.transaction_id}/verify`, {
                unique_transaction_id: uniqueTransactionId,
                transaction_status: newStatus
            })
            onUpdate()
        } catch (error) {
            console.error('Error updating withdrawal status:', error)
        } finally {
            setIsVerifying(false)
            setShowActions(false)
        }
    }

    return (
        <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100 hover:shadow-md transition-all">
            <div className="p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="bg-white/10 p-2 rounded-lg">
                            <Wallet className="w-5 h-5 text-white" />
                        </div>
                        <div>
                            <span className="text-white/70 text-sm">Request ID</span>
                            <p className="text-white font-medium">
                                #{request?.transaction_id}
                            </p>
                        </div>
                    </div>
                    <div className="relative">
                        <button
                            onClick={() => setShowActions(!showActions)}
                            className="p-2 hover:bg-white/10 rounded-lg transition-colors"
                        >
                            <MoreVertical className="w-5 h-5 text-white" />
                        </button>

                        {showActions && (
                            <div className="absolute right-0 mt-2 w-72 bg-white rounded-lg shadow-lg py-2 z-10">
                                <div className="px-4 py-2 border-b">
                                    <p className="text-sm font-medium">Verify Transaction</p>
                                    <p className="text-xs text-gray-500 mt-1">Enter transaction details to verify the withdrawal</p>
                                </div>
                                <div className="p-4">
                                    <div className="mb-4">
                                        <label className="text-xs text-gray-500 mb-1 block">
                                            Transaction ID *
                                        </label>
                                        <input
                                            type="text"
                                            disabled={request.transaction_status === 'success' || request.transaction_status === 'failed'}
                                            value={uniqueTransactionId}
                                            onChange={(e) => setUniqueTransactionId(e.target.value)}
                                            placeholder="Enter transaction ID"
                                            className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 focus:border-[var(--primary-color)]"
                                        />
                                    </div>

                                    <div className="flex flex-col gap-2">
                                        <button
                                            onClick={() => handleStatusUpdate('success')}
                                            className="w-full px-4 py-2 text-sm rounded-lg bg-green-500 text-white hover:bg-green-600 
                                                     disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
                                            disabled={request.transaction_status === 'success' || request.transaction_status === 'failed' || isVerifying}
                                        >
                                            {isVerifying ? (
                                                <>
                                                    <RefreshCw className="w-4 h-4 animate-spin" />
                                                    Verifying...
                                                </>
                                            ) : (
                                                'Verify & Complete'
                                            )}
                                        </button>
                                        <button
                                            onClick={() => handleStatusUpdate('failed')}
                                            className="w-full px-4 py-2 text-sm rounded-lg border border-red-200 text-red-600 
                                                     hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                            disabled={request.transaction_status === 'success' || request.transaction_status === 'failed' || isVerifying}
                                        >
                                            Mark as Failed
                                        </button>
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            <div className="p-6">
                <div className="grid grid-cols-2 gap-6">
                    <div className="bg-gray-50 p-4 rounded-lg">
                        <p className="text-sm text-gray-500 mb-1">Amount</p>
                        <p className="font-semibold text-xl">₹{request.amount}</p>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                        <p className="text-sm text-gray-500 mb-1">Status</p>
                        <StatusBadge status={request.transaction_status} />
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                        <p className="text-sm text-gray-500 mb-1">Partner </p>
                        <p className="font-medium text-gray-700 font-mono">
                            {request?.delivery_partner_id}-{request?.delivery_partner?.user?.first_name} {request?.delivery_partner?.user?.last_name}
                        </p>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                        <p className="text-sm text-gray-500 mb-1">Requested On</p>
                        <p className="font-medium text-gray-700">
                            {new Date(request.created_at).toLocaleDateString('en-US', {
                                day: 'numeric',
                                month: 'short',
                                year: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit'
                            })}
                        </p>
                    </div>
                </div>

                {request.payment_method?.unique_transaction_id && (
                    <div className="mt-6 p-4 bg-green-50 rounded-lg border border-green-100">
                        <div className="flex items-center gap-2 mb-2">
                            <CheckCircle2 className="w-4 h-4 text-green-500" />
                            <p className="text-sm font-medium text-green-700">Verified Transaction</p>
                        </div>
                        <p className="text-sm text-gray-500">Transaction ID</p>
                        <p className="font-medium text-green-700 font-mono">
                            {request.payment_method.unique_transaction_id}
                        </p>
                    </div>
                )}
            </div>
        </div>
    )
}

const DeliveryPartnerPaymentManager = () => {
    const [isLoading, setIsLoading] = useState(true)
    const [isRefreshing, setIsRefreshing] = useState(false)

    const [searchQuery, setSearchQuery] = useState('')
    const [withdrawalRequests, setWithdrawalRequests] = useState([])
    useEffect(() => {
        fetchWithdrawalRequests()
    }, [])


    const fetchWithdrawalRequests = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
            const { data } = await axiosClient.get('/api/delivery-partner-wallet/withdrawal-requests')
            setWithdrawalRequests(data)
        } catch (error) {
            console.error('Error fetching withdrawal requests:', error)
        } finally {
            setIsLoading(false)
            setIsRefreshing(false)
        }
    }

    const createWithdrawalRequest = async () => {
        try {
            const { data } = await axiosClient.post('/api/delivery-partner-wallet/withdrawal-requests')
            console.log(data)
            fetchWithdrawalRequests()
        } catch (error) {
            console.error('Error creating withdrawal request:', error)
        }
    }






    if (isLoading) return (
        <div className="flex items-center justify-center min-h-screen">
            <div className="flex flex-col items-center gap-4">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[var(--primary-color)]"></div>
                <p className="text-gray-500">Loading withdrawal requests...</p>
            </div>
        </div>
    )

    return (
        <div className="max-w-6xl mx-auto p-6">
            <div className="flex items-center justify-between mb-8">
                <div>
                    <h2 className="text-2xl font-bold text-gray-800">Withdrawal Requests</h2>
                    <p className="text-gray-500 mt-1">Manage and verify delivery partner withdrawal requests</p>
                </div>
                <div className="flex items-center gap-4">
                    <HostegoButton text='+ Withdrawal Request' onClick={() => createWithdrawalRequest()} />
                    <button
                        onClick={() => fetchWithdrawalRequests(true)}
                        disabled={isRefreshing}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white border border-gray-200
                                 text-gray-600 hover:bg-gray-50 transition-all duration-200 disabled:opacity-50"
                    >
                        <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
                        {isRefreshing ? 'Refreshing...' : 'Refresh'}
                    </button>
                </div>
            </div>

            <div className="grid grid-cols-1 gap-6">
                {withdrawalRequests.map(request => (
                    <WithdrawalRequestCard
                        key={request?.transaction_id}
                        request={request}
                        onUpdate={() => fetchWithdrawalRequests()}
                    />
                ))}

                {withdrawalRequests.length === 0 && (
                    <div className="text-center py-12 bg-white rounded-xl border border-gray-100">
                        <div className="w-16 h-16 mx-auto mb-4 flex items-center justify-center rounded-full bg-gray-100">
                            <Wallet className="w-8 h-8 text-gray-400" />
                        </div>
                        <h3 className="text-lg font-medium text-gray-900 mb-2">No withdrawal requests found</h3>
                        <p className="text-gray-500">
                            {searchQuery
                                ? "Try adjusting your search or filters"
                                : "No pending withdrawal requests at the moment"}
                        </p>
                    </div>
                )}
            </div>
        </div>
    )
}

export default DeliveryPartnerPaymentManager
