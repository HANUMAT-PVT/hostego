'use client'
import React, { useState, useEffect } from 'react'
import {
    Store,
    Search,
    RefreshCw,
    CheckCircle2,
    XCircle,
    Clock,
    MoreVertical,
    ChevronDown,
    MapPin,
    Phone,
    Mail,
    IndianRupee,
    Building2,
    Calendar,
    CreditCard,
    User,
    StarIcon
} from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import HostegoButton from '../HostegoButton'

const StatusBadge = ({ status }) => {
    const configs = {
        pending: { color: 'bg-yellow-100 text-yellow-700', label: 'Pending', icon: Clock },
        processing: { color: 'bg-blue-100 text-blue-700', label: 'Processing', icon: RefreshCw },
        paid: { color: 'bg-green-100 text-green-700', label: 'Paid', icon: CheckCircle2 },
        failed: { color: 'bg-red-100 text-red-700', label: 'Failed', icon: XCircle }
    }

    const config = configs[status] || configs.pending
    const Icon = config.icon

    return (
        <span className={`px-3 py-1 rounded-full text-xs font-medium flex items-center gap-1 ${config.color}`}>
            <Icon size={12} />
            {config.label}
        </span>
    )
}

const RestaurantPayoutCard = ({ payout, onUpdate }) => {
    const [showActions, setShowActions] = useState(false)
    const [paymentRefId, setPaymentRefId] = useState('')
    const [paymentMethod, setPaymentMethod] = useState('')
    const [isVerifying, setIsVerifying] = useState(false)
    const [isOpen, setIsOpen] = useState(false)

    const handleVerifyPayout = async () => {
        try {
            if (!paymentRefId || !paymentMethod) {
                alert('Please enter both Payment Reference ID and Payment Method')
                return
            }

            setIsVerifying(true)
            await axiosClient.patch(`/api/restaurant-payout/verify/${payout.payout_id}`, {
                payout_id: payout.payout_id.toString(),
                payment_ref_id: paymentRefId,
                payment_method: paymentMethod
            })
            onUpdate()
        } catch (error) {
            console.error('Error verifying payout:', error)
            alert('Error verifying payout. Please try again.')
        } finally {
            setIsVerifying(false)
            setShowActions(false)
        }
    }

    return (
        <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100 hover:shadow-md transition-all">
            {/* Header with Restaurant Info (accordion trigger) */}
            <button
                type="button"
                onClick={() => setIsOpen(!isOpen)}
                className="w-full text-left p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600"
            >
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <div className="bg-white/10 p-3 rounded-lg">
                            <Store className="w-6 h-6 text-white" />
                        </div>
                        <div>
                            <h3 className="text-white font-semibold text-lg">{payout.shop.shop_name}</h3>
                            <div className="flex items-center gap-2 text-white/80 text-sm">
                                <Building2 className="w-4 h-4" />
                                <span>ID: {payout.shop.shop_id}</span>
                                <span className="mx-2">•</span>
                                <span className="flex items-center gap-1">
                                    <StarIcon className="w-3 h-3 fill-current" />
                                    {payout.shop.shop_rating || 0}
                                </span>
                            </div>
                        </div>
                    </div>
                    <div className="flex items-center gap-4">
                        <div className="text-right">
                            <span className="text-white/70 text-sm">Payout ID</span>
                            <p className="text-white font-medium">#{payout.payout_id}</p>
                        </div>
                        <div className="hidden sm:flex items-center gap-2 text-white/90">
                            <IndianRupee className="w-4 h-4" />
                            <span className="font-semibold">₹{payout.total_amount.toLocaleString()}</span>
                        </div>
                        <div className="sm:hidden text-white/90 text-sm font-semibold">₹{payout.total_amount.toLocaleString()}</div>
                        <div className="hidden xs:block">
                            <StatusBadge status={payout.status} />
                        </div>
                        <ChevronDown className={`w-5 h-5 text-white transition-transform ${isOpen ? 'rotate-180' : ''}`} />
                        {payout.status === 'pending' && (
                            <div className="relative" onClick={(e) => e.stopPropagation()}>
                                <button
                                    onClick={() => setShowActions(!showActions)}
                                    className="p-2 hover:bg-white/10 rounded-lg transition-colors"
                                >
                                    <MoreVertical className="w-5 h-5 text-white" />
                                </button>

                                {showActions && (
                                    <div className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg py-2 z-10">
                                        <div className="px-4 py-2 border-b">
                                            <p className="text-sm font-medium">Verify Payout</p>
                                            <p className="text-xs text-gray-500 mt-1">Enter payment details to verify the payout</p>
                                        </div>
                                        <div className="p-4">
                                            <div className="mb-4">
                                                <label className="text-xs text-gray-500 mb-1 block">Payment Reference ID *</label>
                                                <input
                                                    type="text"
                                                    value={paymentRefId}
                                                    onChange={(e) => setPaymentRefId(e.target.value)}
                                                    placeholder="Enter payment reference ID"
                                                    className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 focus:border-[var(--primary-color)]"
                                                />
                                            </div>
                                            <div className="mb-4">
                                                <label className="text-xs text-gray-500 mb-1 block">Payment Method *</label>
                                                <select
                                                    value={paymentMethod}
                                                    onChange={(e) => setPaymentMethod(e.target.value)}
                                                    className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 focus:border-[var(--primary-color)]"
                                                >
                                                    <option value="">Select payment method</option>
                                                    <option value="bank_transfer">Bank Transfer</option>
                                                    <option value="upi">UPI</option>
                                                    <option value="razorpay">Razorpay</option>
                                                    <option value="paytm">Paytm</option>
                                                    <option value="phonepe">PhonePe</option>
                                                    <option value="googlepay">Google Pay</option>
                                                </select>
                                            </div>
                                            <button
                                                onClick={handleVerifyPayout}
                                                className="w-full px-4 py-2 text-sm rounded-lg bg-green-500 text-white hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
                                                disabled={isVerifying}
                                            >
                                                {isVerifying ? (
                                                    <>
                                                        <RefreshCw className="w-4 h-4 animate-spin" />
                                                        Verifying...
                                                    </>
                                                ) : (
                                                    <>
                                                        <CheckCircle2 className="w-4 h-4" />
                                                        Verify & Complete Payout
                                                    </>
                                                )}
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            </button>

            {/* Main Content (accordion body) */}
            {isOpen && (
                <div className="p-6">
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        {/* Left Column - Restaurant Details */}
                        <div className="lg:col-span-2 space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div className="bg-gradient-to-br from-green-50 to-emerald-50 p-4 rounded-lg border border-green-100">
                                    <div className="flex items-center gap-2 mb-2">
                                        <IndianRupee className="w-5 h-5 text-green-600" />
                                        <p className="text-sm text-green-600 font-medium">Payout Amount</p>
                                    </div>
                                    <p className="font-bold text-2xl text-green-700">₹{payout.total_amount.toLocaleString()}</p>
                                </div>
                                <div className="bg-gray-50 p-4 rounded-lg">
                                    <p className="text-sm text-gray-500 mb-2">Status</p>
                                    <StatusBadge status={payout.status} />
                                </div>
                            </div>

                            <div className="bg-gray-50 p-4 rounded-lg">
                                <div className="flex items-start gap-3">
                                    <MapPin className="w-5 h-5 text-gray-500 mt-0.5" />
                                    <div>
                                        <p className="text-sm text-gray-500 mb-1">Restaurant Address</p>
                                        <p className="font-medium text-gray-700">{payout.shop.address}</p>
                                        <p className="text-sm text-gray-500 mt-1">{payout.shop.shop_type}</p>
                                    </div>
                                </div>
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div className="bg-gray-50 p-4 rounded-lg">
                                    <div className="flex items-center gap-2 mb-2">
                                        <User className="w-4 h-4 text-gray-500" />
                                        <p className="text-sm text-gray-500">Owner Details</p>
                                    </div>
                                    <p className="font-medium text-gray-700">{payout.shop.owner_name}</p>
                                    <div className="flex items-center gap-2 mt-2 text-sm text-gray-500">
                                        <Phone className="w-3 h-3" />
                                        <span>{payout.shop.owner_phone}</span>
                                    </div>
                                    <div className="flex items-center gap-2 mt-1 text-sm text-gray-500">
                                        <Mail className="w-3 h-3" />
                                        <span>{payout.shop.owner_email}</span>
                                    </div>
                                </div>
                                <div className="bg-gray-50 p-4 rounded-lg">
                                    <div className="flex items-center gap-2 mb-2">
                                        <Building2 className="w-4 h-4 text-gray-500" />
                                        <p className="text-sm text-gray-500">Bank Details</p>
                                    </div>
                                    <p className="font-medium text-gray-700">{payout.shop.bank_name}</p>
                                    <p className="text-sm text-gray-500 mt-1">
                                        A/c: ****{payout.shop.bank_account_number?.slice(-4)}
                                    </p>
                                    <p className="text-sm text-gray-500">
                                        IFSC: {payout.shop.bank_ifsc_code}
                                    </p>
                                </div>
                            </div>
                        </div>

                        {/* Right Column - Timeline */}
                        <div className="space-y-4">
                            <div className="bg-gray-50 p-4 rounded-lg">
                                <div className="flex items-center gap-2 mb-3">
                                    <Calendar className="w-4 h-4 text-gray-500" />
                                    <p className="text-sm text-gray-500 font-medium">Timeline</p>
                                </div>
                                <div className="space-y-3">
                                    <div className="flex items-center gap-3">
                                        <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                                        <div>
                                            <p className="text-sm font-medium">Created</p>
                                            <p className="text-xs text-gray-500">
                                                {new Date(payout.created_at).toLocaleDateString('en-US', {
                                                    day: 'numeric',
                                                    month: 'short',
                                                    year: 'numeric',
                                                    hour: '2-digit',
                                                    minute: '2-digit'
                                                })}
                                            </p>
                                        </div>
                                    </div>
                                    {payout.paid_at && (
                                        <div className="flex items-center gap-3">
                                            <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                                            <div>
                                                <p className="text-sm font-medium">Paid</p>
                                                <p className="text-xs text-gray-500">
                                                    {new Date(payout.paid_at).toLocaleDateString('en-US', {
                                                        day: 'numeric',
                                                        month: 'short',
                                                        year: 'numeric',
                                                        hour: '2-digit',
                                                        minute: '2-digit'
                                                    })}
                                                </p>
                                            </div>
                                        </div>
                                    )}
                                </div>
                            </div>

                            {payout.payment_ref_id && (
                                <div className="bg-green-50 p-4 rounded-lg border border-green-100">
                                    <div className="flex items-center gap-2 mb-2">
                                        <CreditCard className="w-4 h-4 text-green-600" />
                                        <p className="text-sm font-medium text-green-700">Payment Details</p>
                                    </div>
                                    <div className="space-y-2">
                                        <div>
                                            <p className="text-xs text-green-600">Reference ID</p>
                                            <p className="font-mono text-sm text-green-700">{payout.payment_ref_id}</p>
                                        </div>
                                        <div>
                                            <p className="text-xs text-green-600">Method</p>
                                            <p className="text-sm text-green-700 capitalize">{payout.payment_method}</p>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

const RestaurantPayoutManger = () => {
    const [isLoading, setIsLoading] = useState(true)
    const [isRefreshing, setIsRefreshing] = useState(false)
    const [isStartingPayout, setIsStartingPayout] = useState(false)
    const [searchQuery, setSearchQuery] = useState('')
    const [payouts, setPayouts] = useState([])
    const [filteredPayouts, setFilteredPayouts] = useState([])
    const [statusFilter, setStatusFilter] = useState('all')

    useEffect(() => {
        fetchPayouts()
    }, [])

    useEffect(() => {
        filterPayouts()
    }, [payouts, searchQuery, statusFilter])

    const fetchPayouts = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
            const { data } = await axiosClient.get('/api/restaurant-payout')
            setPayouts(data)
        } catch (error) {
            console.error('Error fetching restaurant payouts:', error)
        } finally {
            setIsLoading(false)
            setIsRefreshing(false)
        }
    }

    const startRestaurantPayout = async () => {
        try {
            setIsStartingPayout(true)
            const { data } = await axiosClient.post('/api/restaurant-payout/initiate')
            console.log('Payout started:', data)
            fetchPayouts()
        } catch (error) {
            console.error('Error starting restaurant payout:', error)
            alert('Error starting restaurant payout. Please try again.')
        } finally {
            setIsStartingPayout(false)
        }
    }

    const filterPayouts = () => {
        let filtered = payouts

        // Filter by search query
        if (searchQuery) {
            filtered = filtered.filter(payout =>
                payout.shop.shop_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                payout.shop.owner_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                payout.payout_id.toString().includes(searchQuery)
            )
        }

        // Filter by status
        if (statusFilter !== 'all') {
            filtered = filtered.filter(payout => payout.status === statusFilter)
        }

        setFilteredPayouts(filtered)
    }

    if (isLoading) return (
        <div className="flex items-center justify-center min-h-screen">
            <div className="flex flex-col items-center gap-4">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[var(--primary-color)]"></div>
                <p className="text-gray-500">Loading restaurant payouts...</p>
            </div>
        </div>
    )

    return (
        <div className="max-w-7xl mx-auto p-6">
            {/* Header */}
            <div className="flex flex-col lg:flex-row lg:items-center justify-between mb-8 gap-4">
                <div>
                    <h2 className="text-3xl font-bold text-gray-800">Restaurant Payout Manager</h2>
                    <p className="text-gray-500 mt-1">Manage and verify restaurant payouts efficiently</p>
                </div>
                <div className="flex items-center gap-4">
                    <HostegoButton
                        text={isStartingPayout ? 'Starting Payout...' : '🚀 Start Restaurant Payout'}
                        onClick={startRestaurantPayout}
                        isLoading={isStartingPayout}
                    />
                    <button
                        onClick={() => fetchPayouts(true)}
                        disabled={isRefreshing}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white border border-gray-200
                                 text-gray-600 hover:bg-gray-50 transition-all duration-200 disabled:opacity-50"
                    >
                        <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
                        {isRefreshing ? 'Refreshing...' : 'Refresh'}
                    </button>
                </div>
            </div>

            {/* Filters */}
            <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-100 mb-6">
                <div className="flex flex-col lg:flex-row gap-4">
                    <div className="flex-1">
                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                            <input
                                type="text"
                                placeholder="Search by restaurant name, owner, or payout ID..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 focus:border-[var(--primary-color)]"
                            />
                        </div>
                    </div>
                    <div>
                        <select
                            value={statusFilter}
                            onChange={(e) => setStatusFilter(e.target.value)}
                            className="px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]/20 focus:border-[var(--primary-color)]"
                        >
                            <option value="all">All Status</option>
                            <option value="pending">Pending</option>
                            <option value="processing">Processing</option>
                            <option value="paid">Paid</option>
                            <option value="failed">Failed</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                <div className="bg-gradient-to-br from-blue-50 to-indigo-50 p-6 rounded-lg border border-blue-100">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-blue-600 font-medium">Total Payouts</p>
                            <p className="text-2xl font-bold text-blue-700">{payouts.length}</p>
                        </div>
                        <Store className="w-8 h-8 text-blue-500" />
                    </div>
                </div>
                <div className="bg-gradient-to-br from-yellow-50 to-orange-50 p-6 rounded-lg border border-yellow-100">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-yellow-600 font-medium">Pending</p>
                            <p className="text-2xl font-bold text-yellow-700">
                                {payouts.filter(p => p.status === 'pending').length}
                            </p>
                        </div>
                        <Clock className="w-8 h-8 text-yellow-500" />
                    </div>
                </div>
                <div className="bg-gradient-to-br from-green-50 to-emerald-50 p-6 rounded-lg border border-green-100">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-green-600 font-medium">Completed</p>
                            <p className="text-2xl font-bold text-green-700">
                                {payouts.filter(p => p.status === 'paid').length}
                            </p>
                        </div>
                        <CheckCircle2 className="w-8 h-8 text-green-500" />
                    </div>
                </div>
                <div className="bg-gradient-to-br from-purple-50 to-violet-50 p-6 rounded-lg border border-purple-100">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-purple-600 font-medium">Total Amount</p>
                            <p className="text-2xl font-bold text-purple-700">
                                ₹{payouts.reduce((sum, p) => sum + p.total_amount, 0).toLocaleString()}
                            </p>
                        </div>
                        <IndianRupee className="w-8 h-8 text-purple-500" />
                    </div>
                </div>
            </div>

            {/* Payouts List */}
            <div className="space-y-6">
                {filteredPayouts.map(payout => (
                    <RestaurantPayoutCard
                        key={payout.payout_id}
                        payout={payout}
                        onUpdate={() => fetchPayouts()}
                    />
                ))}

                {filteredPayouts.length === 0 && payouts.length > 0 && (
                    <div className="text-center py-12 bg-white rounded-xl border border-gray-100">
                        <div className="w-16 h-16 mx-auto mb-4 flex items-center justify-center rounded-full bg-gray-100">
                            <Search className="w-8 h-8 text-gray-400" />
                        </div>
                        <h3 className="text-lg font-medium text-gray-900 mb-2">No payouts found</h3>
                        <p className="text-gray-500">Try adjusting your search or filters</p>
                    </div>
                )}

                {payouts.length === 0 && (
                    <div className="text-center py-12 bg-white rounded-xl border border-gray-100">
                        <div className="w-16 h-16 mx-auto mb-4 flex items-center justify-center rounded-full bg-gray-100">
                            <Store className="w-8 h-8 text-gray-400" />
                        </div>
                        <h3 className="text-lg font-medium text-gray-900 mb-2">No restaurant payouts found</h3>
                        <p className="text-gray-500 mb-4">
                            Get started by initiating restaurant payouts
                        </p>

                    </div>
                )}
            </div>
        </div>
    )
}

export default RestaurantPayoutManger
