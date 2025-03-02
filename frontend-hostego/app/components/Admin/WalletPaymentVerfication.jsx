'use client'
import React, { useState, useEffect } from 'react'
import { formatDate } from '@/app/utils/helper'
import { CheckCircle2, Clock, IndianRupee, X, ExternalLink, Image as ImageIcon, RefreshCw } from 'lucide-react'
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
                            : 'bg-green-50 text-green-600'
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

    const fetchPaymentTransactions = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
            const response = await axiosClient.get('/api/wallet/all-transactions?status=pending')
            setPaymentTransactions(response?.data)
        } catch (error) {
            console.error('Error fetching transactions:', error)
        } finally {
            setIsLoading(false)
            setIsRefreshing(false)
        }
    }

    useEffect(() => {
        fetchPaymentTransactions()
    }, [])

    const handleRefresh = () => {
        fetchPaymentTransactions(true)
    }

    if (isLoading) {
        return <HostegoLoader />
    }

    return (
        <div className="max-w-2xl mx-auto p-4 space-y-4">
            <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-semibold">Payment Verifications</h2>
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
                        There are no payment transactions to verify at the moment.
                        New transactions will appear here when users add money to their wallet.
                    </p>
                </div>
            ) : (
                <div className="space-y-4">
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