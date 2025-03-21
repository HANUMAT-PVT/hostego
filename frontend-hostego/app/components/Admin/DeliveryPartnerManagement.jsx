'use client'
import React, { useState, useEffect } from 'react'
import { Shield, User, MapPin, CheckCircle2, X, ExternalLink, Image as ImageIcon, RefreshCw, Phone, Mail, AlertCircle, MoreVertical, Search } from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import HostegoLoader from '../HostegoLoader'
import HostegoButton from '../HostegoButton'

const StatusBadge = ({ status, type }) => {
  const configs = {
    verification: {
      0: { color: 'bg-yellow-100 text-yellow-700', label: 'Pending Verification' },
      1: { color: 'bg-green-100 text-green-700', label: 'Verified' }
    },
    account: {
      0: { color: 'bg-red-100 text-red-700', label: 'Inactive' },
      1: { color: 'bg-green-100 text-green-700', label: 'Active' }
    },
    availability: {
      0: { color: 'bg-gray-100 text-gray-700', label: 'Offline' },
      1: { color: 'bg-green-100 text-green-700', label: 'Online' }
    }
  }

  const config = configs[type][status] || configs[type][0]

  return (
    <span className={`px-2 py-1 rounded-full text-xs font-medium ${config.color}`}>
      {config.label}
    </span>
  )
}

const DeliveryPartnerProfileCard = ({ partner, onUpdate }) => {
  const [showDeliveryPartnerDetails, setShowDeliveryPartnerDetails] = useState(false)
  const [showActions, setShowActions] = useState(false)

  const handleStatusUpdate = async (field, value) => {
    try {

      await axiosClient.patch(`/api/delivery-partner/${partner?.delivery_partner_id}`, {
        [field]: value
      })
      onUpdate()
    } catch (error) {
      console.error('Error updating partner:', error)
    } finally {

      setShowActions(false)
    }
  }

  return (
    <div className="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100">
      {/* Partner Header */}
      <div onClick={() => setShowDeliveryPartnerDetails(!showDeliveryPartnerDetails)} className="p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <User className="w-5 h-5 text-white" />
            <span className="text-white font-medium">
              {partner?.user?.first_name} {partner?.user?.last_name} #{partner?.delivery_partner_id?.slice(0, 8)}
            </span>
          </div>
          <div className="relative">
            <button
              onClick={(e) => {
                e.stopPropagation()
                setShowActions(!showActions)
              }}
              className="p-1 hover:bg-white/10 rounded-full transition-colors"
            >
              <MoreVertical className="w-5 h-5 text-white" />
            </button>

            {/* Actions Dropdown */}
            {showActions && (
              <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-1 z-10">
                <div className="px-4 py-2 border-b">
                  <p className="text-sm font-medium">Update Status</p>
                </div>

                {/* Verification Status */}
                <div className="px-4 py-2 border-b">
                  <p className="text-xs text-gray-500 mb-2">Verification Status</p>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleStatusUpdate('verification_status', 1)}
                      className="px-2 py-1 text-xs rounded bg-green-50 text-green-600 hover:bg-green-100"
                      disabled={partner.verification_status === 1}
                    >
                      Verify
                    </button>
                    <button
                      onClick={() => handleStatusUpdate('verification_status', 0)}
                      className="px-2 py-1 text-xs rounded bg-red-50 text-red-600 hover:bg-red-100"
                      disabled={partner.verification_status === 0}
                    >
                      Unverify
                    </button>
                  </div>
                </div>

                {/* Account Status */}
                <div className="px-4 py-2 border-b">
                  <p className="text-xs text-gray-500 mb-2">Account Status</p>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleStatusUpdate('account_status', 1)}
                      className="px-2 py-1 text-xs rounded bg-green-50 text-green-600 hover:bg-green-100"
                      disabled={partner.account_status === 1}
                    >
                      Activate
                    </button>
                    <button
                      onClick={() => handleStatusUpdate('account_status', 0)}
                      className="px-2 py-1 text-xs rounded bg-red-50 text-red-600 hover:bg-red-100"
                      disabled={partner.account_status === 0}
                    >
                      Deactivate
                    </button>
                  </div>
                </div>
                <div className="px-4 py-2 border-b">
                  <p className="text-xs text-gray-500 mb-2">Availability Status</p>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleStatusUpdate('availability_status', 1)}
                      className="px-2 py-1 text-xs rounded bg-green-50 text-green-600 hover:bg-green-100"
                      disabled={partner.availability_status === 1}
                    >
                      Online
                    </button>

                    <button
                      onClick={() => handleStatusUpdate('availability_status', 0)}
                      className="px-2 py-1 text-xs rounded bg-red-50 text-red-600 hover:bg-red-100"
                      disabled={partner.availability_status === 0}
                    >
                      Offline
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Partner Details */}
      {showDeliveryPartnerDetails && <div className="p-4">
        {/* Status Badges */}
        <div className="flex flex-wrap gap-2 mb-4">
          <StatusBadge status={partner.verification_status} type="verification" />
          <StatusBadge status={partner.account_status} type="account" />
          <StatusBadge status={partner.availability_status} type="availability" />
        </div>

        {/* Personal Information */}
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-gray-500">Name</p>
              <p className="font-medium">{partner.user?.first_name || 'N/A'}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Phone</p>
              <p className="font-medium">{partner.user?.mobile_number}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Email</p>
              <p className="font-medium">{partner.user?.email || 'N/A'}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Address</p>
              <p className="font-medium">{partner?.address || 'N/A'}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">UPI ID</p>
              <p className="font-medium">{partner?.documents?.upi_id || 'N/A'}</p>
            </div>
          </div>
        </div>

        {/* Documents Section */}
        <div className="mt-4 pt-4 border-t">
          <h3 className="font-medium text-gray-800 mb-3">Documents</h3>
          <div className="grid grid-cols-2 gap-4">
            {['aadhaar_front_img', 'aadhaar_back_img', 'bank_details_img'].map((doc) => (
              partner.documents?.[doc] && (
                <a
                  key={doc}
                  href={partner.documents[doc]}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg hover:bg-gray-100"
                >
                  <ImageIcon className="w-4 h-4 text-gray-500" />
                  <span className="text-sm font-medium">
                    {doc.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')}
                  </span>
                  <ExternalLink className="w-4 h-4 ml-auto text-gray-400" />
                </a>
              )
            ))}
          </div>
        </div>
      </div>}
    </div>
  )
}

const DeliveryPartnerManagement = () => {
  const [partners, setPartners] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [filter, setFilter] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')

  const fetchPartners = async (showRefreshAnimation = false) => {
    try {
      showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
      const { data } = await axiosClient.get('/api/delivery-partner/all')
      setPartners(data)
    } catch (error) {
      console.error('Error fetching partners:', error)
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  useEffect(() => {
    fetchPartners()
  }, [])

  const filteredPartners = partners.filter(partner => {
    // First apply status filters
    const statusFilter = () => {
      if (filter === 'verified') return partner.verification_status === 1
      if (filter === 'unverified') return partner.verification_status === 0
      if (filter === 'active') return partner.account_status === 1
      if (filter === 'inactive') return partner.account_status === 0
      return true
    }

    // Then apply search filter
    const searchFilter = () => {
      if (!searchQuery) return true
      const query = searchQuery.toLowerCase()
      return (
        partner.delivery_partner_id?.toLowerCase().includes(query) ||
        partner.user?.first_name?.toLowerCase().includes(query) ||
        partner.user?.mobile_number?.includes(query) ||
        partner.user?.email?.toLowerCase().includes(query)
      )
    }

    return statusFilter() && searchFilter()
  })

  if (isLoading) return <HostegoLoader />

  return (
    <div className="max-w-6xl mx-auto p-4">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold">Delivery Partners</h2>
        <div className="flex items-center gap-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
            <input
              type="text"
              placeholder="Search by name, ID, phone..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-64 pl-10 pr-4 py-2 rounded-lg border-2 border-gray-100 
                       focus:border-[var(--primary-color)] outline-none transition-all"
            />
          </div>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="px-4 py-2 rounded-lg border-2 border-gray-100 
                     focus:border-[var(--primary-color)] outline-none"
          >
            <option value="all">All Partners</option>
            <option value="verified">Verified</option>
            <option value="unverified">Unverified</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
          <button
            onClick={() => fetchPartners(true)}
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

      {/* Results count */}
      <div className="mb-4 text-sm text-gray-500">
        Found {filteredPartners?.length} delivery partner{filteredPartners?.length !== 1 ? 's' : ''}
      </div>

      {/* No results message */}
      {filteredPartners?.length === 0 && (
        <div className="text-center py-8 bg-white rounded-xl">
          <div className="w-16 h-16 mx-auto mb-4 flex items-center justify-center rounded-full bg-gray-100">
            <User className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-1">No delivery partners found</h3>
          <p className="text-gray-500">
            {searchQuery
              ? "Try adjusting your search or filters"
              : "No delivery partners match the selected filters"}
          </p>
        </div>
      )}

      <div className="grid grid-cols-1 gap-6 overflow-y-auto max-h-[85vh]">
        {filteredPartners?.map(partner => (
          <DeliveryPartnerProfileCard
            key={partner?.delivery_partner_id}
            partner={partner}
            onUpdate={fetchPartners}
          />
        ))}
      </div>
    </div>
  )
}

export default DeliveryPartnerManagement
