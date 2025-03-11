'use client'
import React, { useState, useEffect } from 'react'
import { formatDate } from '@/app/utils/helper'
import { CheckCircle2, Clock, IndianRupee, X, ExternalLink, Image as ImageIcon, RefreshCw, Filter, Search } from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import Image from 'next/image'
import HostegoLoader from '../HostegoLoader'



const PaymentCard = ({ transaction, onVerify, onReject }) => {
    const [isLoading, setIsLoading] = useState(false)
    const [isImageModalOpen, setIsImageModalOpen] = useState(false)

    const handleVerify = async (status, transactionId) => {
        try {
            setIsLoading(true)
            await axiosClient.post(`/api/wallet/verifiy-wallet-transaction/${transactionId}`, { "transaction_status": status })
            onVerify(transactionId)
        } catch (error) {
            console.error('Error verifying payment:', error)
        } finally {
            setIsLoading(false)
        }
    }

    if (isLoading) {
        return <HostegoLoader />
    }


    return (
        <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100 transition-all duration-200 hover:shadow-md">
            {/* Header */}
            <div className="p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <Clock className="w-5 h-5 text-white" />
                        <span className="text-white font-medium">
                            {formatDate(transaction?.created_at)}
                        </span>
                    </div>
                    <span className="px-3 py-1 bg-white/20 backdrop-blur-sm text-white rounded-full text-sm">
                        #{transaction?.transaction_id?.slice(0, 8)}
                    </span>
                </div>
            </div>

            {/* Content */}
            <div className="p-4">
                {/* Amount */}
                <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-2">
                        <div className="w-10 h-10 rounded-full bg-[var(--primary-color)]/10 flex items-center justify-center">
                            <IndianRupee className="w-5 h-5 text-[var(--primary-color)]" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Amount</p>
                            <p className="text-xl font-semibold">₹{transaction?.amount}</p>
                        </div>
                    </div>
                    <span className={`px-3 py-1 rounded-full text-sm 
                        ${transaction?.transaction_status === 'pending'
                            ? 'bg-orange-50 text-orange-600'
                            : transaction?.transaction_status === 'success'
                                ? 'bg-green-50 text-green-600'
                                : transaction?.transaction_status === 'failed'
                                    ? 'bg-red-50 text-red-600'
                                    : 'bg-gray-50 text-gray-600'
                        }`}>
                        {transaction?.transaction_status?.toUpperCase()}
                    </span>
                </div>

                {/* Transaction Details */}
                <div className="space-y-3 mb-4">
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Name</span>
                        <span className="font-medium capitalize">{transaction?.user?.first_name}{" "} {transaction?.user?.last_name}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Phone Number</span>
                        <span className="font-medium capitalize">{transaction?.user?.mobile_number}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Transaction Type</span>
                        <span className="font-medium capitalize">{transaction?.transaction_type}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-600">UPI Transaction ID</span>
                        <span className="font-medium">{transaction?.payment_method?.unique_transaction_id}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Image Url</span>
                        <span className="font-medium">{transaction?.payment_method?.payment_screenshot_img_url}</span>
                    </div>
                </div>

                {/* Screenshot Preview */}
                {transaction?.payment_method?.payment_screenshot_img_url && (
                    <div className="mt-4">
                        <button
                            onClick={() => setIsImageModalOpen(true)}
                            className="w-full p-3 border-2 border-dashed border-gray-200 rounded-lg flex items-center justify-center gap-2 hover:border-[var(--primary-color)] transition-colors"
                        >
                            <ImageIcon className="w-5 h-5 text-gray-500" />
                            <span className="text-sm font-medium text-gray-600">View Payment Screenshot</span>
                        </button>
                    </div>
                )}

                {/* Action Buttons */}
                <div className="flex gap-3 mt-4">
                    <button
                        onClick={() => handleVerify("success", transaction?.transaction_id)}
                        disabled={isLoading || transaction?.transaction_status !== 'pending'}
                        className="flex-1 py-2 px-4 bg-[var(--primary-color)] text-white rounded-lg font-medium 
                                 hover:opacity-90 transition-all duration-200 disabled:opacity-50 
                                 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                        <CheckCircle2 className="w-4 h-4" />
                        {isLoading ? 'Verifying...' : 'Verify Payment'}
                    </button>
                    <button
                        onClick={() => handleVerify("failed", transaction?.transaction_id)}
                        disabled={isLoading || transaction?.transaction_status !== 'pending'}
                        className="flex-1 py-2 px-4 border-2 border-red-500 text-red-500 rounded-lg font-medium 
                                 hover:bg-red-50 transition-all duration-200 disabled:opacity-50 
                                 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                        <X className="w-4 h-4" />
                        Reject
                    </button>
                </div>
            </div>

            {/* Image Modal */}
            {isImageModalOpen && (
                <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
                    <div className="bg-white rounded-xl p-4 max-w-2xl w-[90%] max-h-[90vh] overflow-auto">
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-semibold">Payment Screenshot</h3>
                            <button
                                onClick={() => setIsImageModalOpen(false)}
                                className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                            >
                                <X className="w-5 h-5" />
                            </button>
                        </div>
                        <img
                            src={transaction?.payment_method?.payment_screenshot_img_url}
                            alt="Payment Screenshot"
                            className="w-full rounded-lg"
                        />
                        <div className="mt-4 flex justify-end">
                            <a
                                href={transaction?.payment_method?.payment_screenshot_img_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex items-center gap-2 text-[var(--primary-color)] hover:underline"
                            >
                                <ExternalLink className="w-4 h-4" />
                                Open Original
                            </a>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

const WalletPaymentVerfication = () => {
    const [paymentTransactions, setPaymentTransactions] = useState([])
    const [isLoading, setIsLoading] = useState(true)
    const [isRefreshing, setIsRefreshing] = useState(false)
    const [transactionType, setTransactionType] = useState('credit') // Default to credit
    const [transactionStatus, setTransactionStatus] = useState('pending') // Default to pending
    const [searchTerm, setSearchTerm] = useState('')
    
    // Use array destructuring for useDebounce ho

    const fetchPaymentTransactions = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
            const queryParams = new URLSearchParams()

            // Only add params if they have values
            if (transactionStatus) {
                queryParams.append('transaction_status', transactionStatus)
            }
            if (transactionType) {
                queryParams.append('transaction_type', transactionType)
            }
            if (searchTerm) {
                queryParams.append('search', searchTerm)
            }

            const response = await axiosClient.get(`/api/wallet/all-transactions?${queryParams}&search=${searchTerm}`)
            setPaymentTransactions(response?.data)
        } catch (error) {
            console.error('Error fetching transactions:', error)
        } finally {
            setIsLoading(false)
            setIsRefreshing(false)
        }
    }

    // Add debouncedSearchTerm to dependency array
    useEffect(() => {
        fetchPaymentTransactions()
    }, [transactionType, transactionStatus, searchTerm])

    const handleRefresh = () => {
        fetchPaymentTransactions(true)
    }

    

    return (
        <div className="max-w-2xl mx-auto p-4 space-y-4">
            <div className="mb-6">
                <h2 className="text-xl font-semibold mb-4">Payment Verifications</h2>

                {/* Search and Filters */}
                <div className="flex flex-col sm:flex-row gap-4">
                    {/* Search Input */}
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
                        <input
                            type="text"
                            placeholder="Search by ID, name, phone or amount..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="w-full pl-10 pr-4 py-2 rounded-lg border-2 border-gray-100 
                                    focus:border-[var(--primary-color)] outline-none"
                        />
                        {searchTerm && (
                            <button
                                onClick={() => setSearchTerm('')}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 
                                         hover:text-gray-600 transition-colors"
                            >
                                <X size={16} />
                            </button>
                        )}
                    </div>

                    {/* Filters */}
                    <div className="flex items-center gap-3">
                        {/* Transaction Type Filter */}
                        <select
                            value={transactionType}
                            onChange={(e) => setTransactionType(e.target.value)}
                            className="px-4 py-2 rounded-lg border-2 border-gray-100 
                                    focus:border-[var(--primary-color)] outline-none text-sm"
                        >
                            <option value="">All Types</option>
                            <option value="credit">Credit</option>
                            <option value="debit">Debit</option>
                            <option value="refund">Refund</option>
                        </select>

                        {/* Status Filter */}
                        <select
                            value={transactionStatus}
                            onChange={(e) => setTransactionStatus(e.target.value)}
                            className="px-4 py-2 rounded-lg border-2 border-gray-100 
                                    focus:border-[var(--primary-color)] outline-none text-sm"
                        >
                            <option value="">All Status</option>
                            <option value="pending">Pending</option>
                            <option value="success">Success</option>
                            <option value="failed">Failed</option>
                        </select>

                        {/* Refresh Button */}
                        <button
                            onClick={handleRefresh}
                            disabled={isRefreshing}
                            className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)]/10 
                                    text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/20 
                                    transition-all duration-200 disabled:opacity-50"
                        >
                            <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
                            {isRefreshing ? 'Refreshing...' : 'Refresh'}
                        </button>
                    </div>
                </div>

                {/* Active Filters Display */}
                {(searchTerm || transactionType || transactionStatus) && (
                    <div className="flex flex-wrap gap-2 mt-4">
                        {searchTerm && (
                            <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm 
                                         bg-[var(--primary-color)]/10 text-[var(--primary-color)]">
                                <Search size={14} />
                                {searchTerm}
                                <button
                                    onClick={() => setSearchTerm('')}
                                    className="ml-1 hover:text-[var(--primary-color)]/70"
                                >
                                    <X size={14} />
                                </button>
                            </span>
                        )}
                        {transactionType && (
                            <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm 
                                         bg-blue-50 text-blue-600">
                                <Filter size={14} />
                                {transactionType}
                                <button
                                    onClick={() => setTransactionType('')}
                                    className="ml-1 hover:text-blue-400"
                                >
                                    <X size={14} />
                                </button>
                            </span>
                        )}
                        {transactionStatus && (
                            <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm 
                                         bg-green-50 text-green-600">
                                <Clock size={14} />
                                {transactionStatus}
                                <button
                                    onClick={() => setTransactionStatus('')}
                                    className="ml-1 hover:text-green-400"
                                >
                                    <X size={14} />
                                </button>
                            </span>
                        )}
                    </div>
                )}
            </div>

            {/* Results count */}
            <div className="text-sm text-gray-500 mb-4">
                Found {paymentTransactions.length} transaction{paymentTransactions.length !== 1 ? 's' : ''}
            </div>
            {isLoading && <HostegoLoader />}
            {isRefreshing ? (
                <div className="flex items-center justify-center py-12">
                    <HostegoLoader />
                </div>
            ) : paymentTransactions.length === 0 ? (
                <div className="text-center py-12 bg-white rounded-xl shadow-sm">
                    <div className="w-20 h-20 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-6">
                        <IndianRupee className="w-10 h-10 text-gray-400" />
                    </div>
                    <h3 className="text-xl font-semibold text-gray-800 mb-3">No Transactions</h3>
                    <p className="text-gray-600 max-w-sm mx-auto">
                        {status || transactionType
                            ? "No transactions match the selected filters."
                            : "There are no payment transactions to verify at the moment."}
                    </p>
                </div>
            ) : (
                <div className="space-y-4 overflow-y-auto max-h-[85vh]">
                    {paymentTransactions.map(transaction => (
                        <PaymentCard
                            key={transaction?.transaction_id}
                            transaction={transaction}
                            onVerify={handleRefresh}
                            onReject={handleRefresh}
                        />
                    ))}
                </div>
            )}
        </div>
    )
}

export default WalletPaymentVerfication