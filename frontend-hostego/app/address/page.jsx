'use client'
import React, { useState, useEffect } from 'react'
import BackNavigationButton from '../components/BackNavigationButton'
import { Building2, Home, MapPin, Navigation, Edit2 } from 'lucide-react'
import axiosClient from '../utils/axiosClient'
import HostegoButton from '../components/HostegoButton'
import HostegoLoader from '../components/HostegoLoader'

const AddressItem = ({ address, onEdit }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editedAddress, setEditedAddress] = useState(address);
  const [isUpdating, setIsUpdating] = useState(false);

  const handleUpdate = async () => {
    try {
      setIsUpdating(true);
      await axiosClient.patch(`/api/address/${address?.address_id}`, editedAddress);
      setIsEditing(false);
      onEdit(); // Refresh address list
    } catch (error) {
      console.error('Error updating address:', error);
    } finally {
      setIsUpdating(false);
    }
  };

  if (isEditing) {
    return (
      <div className='bg-white rounded-lg p-4 shadow-sm border-2 border-[var(--primary-color)] animate-fade-in m-2'>
        <div className='space-y-3'>
          <div>
            <label className="text-xs text-[var(--primary-color)]">Address Type</label>
            <input
              type="text"
              value={editedAddress.address_type}
              onChange={(e) => setEditedAddress({ ...editedAddress, address_type: e.target.value })}
              className="w-full mt-1 p-2 rounded-md border border-gray-200 focus:outline-none focus:border-[var(--primary-color)]"
            />
          </div>
          <div>
            <label className="text-xs text-[var(--primary-color)]">Complete Address</label>
            <textarea
              value={editedAddress.address_line_1}
              onChange={(e) => setEditedAddress({ ...editedAddress, address_line_1: e.target.value })}
              className="w-full mt-1 p-2 rounded-md border border-gray-200 focus:outline-none focus:border-[var(--primary-color)] min-h-[100px] resize-none"
            />
          </div>
          <div className='flex gap-2 mt-3'>
            <button
              onClick={handleUpdate}
              disabled={isUpdating}
              className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium 
                ${isUpdating
                  ? 'bg-gray-100 text-gray-400'
                  : 'bg-[var(--primary-color)] text-white hover:opacity-90'}`}
            >
              {isUpdating ? 'Updating...' : 'Update'}
            </button>
            <button
              onClick={() => setIsEditing(false)}
              className="flex-1 py-2 px-4 rounded-lg text-sm font-medium border-2 border-gray-200 
                       text-gray-700 hover:bg-gray-50 transition-all"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className='m-2 address-item flex items-center rounded-md gap-4 cursor-pointer px-4 py-2 bg-white hover:bg-gray-50 transition-colors group'>
      <div className='bg-[var(--bg-page-color)] p-2 w-[40px] h-[40px] flex justify-center items-center rounded-full'>
        <Home size={20} className='text-[var(--primary-color)]' />
      </div>
      <div className='flex-1'>
        <p className='text-md font-semibold'>{address.address_type}</p>
        <p className='text-sm'>{address.address_line_1}</p>
      </div>
      <button
        onClick={() => setIsEditing(true)}
        className="p-2 hover:bg-gray-100 rounded-full"
      >
        <Edit2 size={16} className="text-[var(--primary-color)]" />
      </button>
    </div>
  );
};

const page = () => {
  const [isLoading, setIsLoading] = useState(false)
  const [isPageLoading, setIsPageLoading] = useState(true)
  const [address, setAddress] = useState([])
  const [addressData, setAddressData] = useState({
    address_type: '',
    address_line_1: '',
  })

  useEffect(() => {
    fetchAddress()
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      setIsLoading(true)
      await axiosClient.post('/api/address', addressData)
      fetchAddress()
      setAddressData({
        address_type: '',
        address_line_1: ''
      })
    } catch (error) {
      console.error('Error saving address:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const fetchAddress = async () => {
    try {
      setIsPageLoading(true)
      const { data } = await axiosClient.get('/api/address')
      setAddress(data)
    } catch (error) {
      console.error('Error fetching addresses:', error)
    } finally {
      setIsPageLoading(false)
    }
  }

  if (isPageLoading) {
    return <HostegoLoader />
  }

  return (
    <div className="min-h-screen bg-[var(--bg-page-color)]">
      <BackNavigationButton title="Address" />

      <form onSubmit={handleSubmit} className="p-2 space-y-4">
        <div className="bg-white rounded-lg p-4 space-y-4 flex flex-col gap-4">
          <p className="text-xl font-normal">Add New Address</p>

          <div className="relative">
            <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
              Address Type
            </label>
            <input
              value={addressData.address_type}
              onChange={(e) => setAddressData({ ...addressData, address_type: e.target.value })}
              className="w-full px-4 py-2 border-2 border-[var(--primary-color)] rounded-md outline-none resize-none"
              placeholder="Address type eg, Home, Friend, etc."
              required
            />
          </div>

          <div className="relative">
            <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
              Complete Address
            </label>
            <textarea
              value={addressData?.address_line_1}
              onChange={(e) => setAddressData({ ...addressData, address_line_1: e.target.value })}
              className="w-full px-4 py-3 border-2 border-[var(--primary-color)] rounded-md outline-none min-h-[100px] resize-none"
              placeholder="Enter your complete address eg, Room No 1115, Zakir-A"
              required
            />
          </div>

          <HostegoButton
            type="submit"
            text={isLoading ? 'Saving...' : 'Save Address'}
            isLoading={isLoading}
          />
        </div>
      </form>

      {/* Saved Addresses */}
      <div className="mt-4">
        <h2 className="px-4 text-lg font-medium text-gray-800 mb-2">Saved Addresses</h2>
        {address?.map((addr) => (
          <AddressItem
            key={addr?.address_id}
            address={addr}
            onEdit={fetchAddress}
          />
        ))}
      </div>
    </div>
  )
}

export default page
