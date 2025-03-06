'use client'

import { useState } from 'react'
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle, TransitionChild } from '@headlessui/react'
import { Home, X, Plus, Edit2, Check } from 'lucide-react'
import axiosClient from '../../utils/axiosClient'
import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

const AddressItem = ({ address, onSelect, onEdit }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editedAddress, setEditedAddress] = useState(address);
  const [isUpdating, setIsUpdating] = useState(false);

  const handleUpdate = async () => {
    try {
      setIsUpdating(true);
      await axiosClient.patch(`/api/address/${address.id}`, editedAddress);
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
      <div className='bg-white rounded-lg p-4 shadow-sm border-2 border-[var(--primary-color)] animate-fade-in'>
        <div className='space-y-3'>
          <div>
            <label className="text-xs text-[var(--primary-color)]">Address Type</label>
            <input
              type="text"
              value={editedAddress.address_type}
              onChange={(e) => setEditedAddress({...editedAddress, address_type: e.target.value})}
              className="w-full mt-1 p-2 rounded-md border border-gray-200 focus:outline-none focus:border-[var(--primary-color)]"
            />
          </div>
          <div>
            <label className="text-xs text-[var(--primary-color)]">Address</label>
            <input
              type="text"
              value={editedAddress.address_line_1}
              onChange={(e) => setEditedAddress({...editedAddress, address_line_1: e.target.value})}
              className="w-full mt-1 p-2 rounded-md border border-gray-200 focus:outline-none focus:border-[var(--primary-color)]"
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
    <div className='address-item flex items-center rounded-md gap-4 cursor-pointer p-4 bg-white hover:bg-gray-50 transition-colors group'>
      <div className='bg-[var(--bg-page-color)] p-2 w-[40px] h-[40px] flex justify-center items-center rounded-full'>
        <Home size={20} className='text-[var(--primary-color)]' />
      </div>
      <div className='flex-1' onClick={() => onSelect(address)}>
        <p className='text-md font-semibold'>{address.address_type}</p>
        <p className='text-sm text-gray-600'>{address.address_line_1}</p>
      </div>
      <button
        onClick={() => setIsEditing(true)}
        className="p-2 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-gray-100 rounded-full"
      >
        <Edit2 size={16} className="text-[var(--primary-color)]" />
      </button>
    </div>
  );
};

export default function AddressList({ openAddressList, setOpenAddressList, sendSelectedAddress }) {
  const router = useRouter()
  const [address, setAddress] = useState([])

  useEffect(() => {
    fetchAddress()
  }, [])

  const fetchAddress = async () => {
    const { data } = await axiosClient.get('/api/address')
    setAddress(data)
  }

  return (
    <Dialog open={openAddressList} onClose={() => setOpenAddressList(!openAddressList)} className="relative z-10">
      <DialogBackdrop
        transition
        className="fixed inset-0 bg-gray-500/75 transition-opacity duration-500 ease-in-out data-closed:opacity-0"
      />

      <div className="fixed inset-0 overflow-hidden animate-slide-up">
        <div className="absolute inset-0 overflow-hidden">
          <div className="pointer-events-none fixed bottom-0 left-0 flex w-full max-h-screen">
            <DialogPanel className="pointer-events-auto relative w-full max-h-[90vh] transform transition duration-500 ease-in-out data-closed:translate-y-full sm:duration-700 bg-white rounded-t-xl shadow-xl">
              <TransitionChild>
                <div className="absolute top-4 right-4 duration-500 ease-in-out data-closed:opacity-0">
                  <button
                    onClick={() => setOpenAddressList(!openAddressList)}
                    className="relative rounded-md text-gray-300 hover:text-black focus:outline-none"
                  >
                    <X className="size-6" />
                  </button>
                </div>
              </TransitionChild>

              <div className="flex flex-col h-full overflow-y-auto py-6 bg-[var(--bg-page-color)]">
                <div className="px-4 sm:px-6">
                  <DialogTitle className="text-xl font-semibold text-gray-900">
                    Select delivery location
                  </DialogTitle>
                </div>

                <div className="relative mt-6 flex-1 px-4">
                  {address?.length > 0 ? (
                    <>
                      <p className='text-md text-gray-500'>Your saved addresses</p>
                      <div className="flex flex-col gap-4 mt-4 max-h-[50vh] overflow-y-auto">
                        {address?.map((addr) => (
                          <AddressItem 
                            key={addr.id}
                            address={addr}
                            onSelect={() => {
                              sendSelectedAddress(addr);
                              setOpenAddressList(false);
                            }}
                            onEdit={fetchAddress}
                          />
                        ))}
                      </div>
                    </>
                  ) : (
                    <div className="text-center py-8">
                      <p className="text-gray-500">No addresses saved yet</p>
                    </div>
                  )}

                  {/* Add New Address Button */}
                  <button
                    onClick={() => {
                      router.push('/address')
                      setOpenAddressList(false)
                    }}
                    className="mt-4 w-full flex items-center justify-center gap-2 p-4 bg-white rounded-lg border-2 border-dashed border-[var(--primary-color)] text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/5 transition-colors"
                  >
                    <Plus className="w-5 h-5" />
                    Add New Address
                  </button>
                </div>
              </div>
            </DialogPanel>
          </div>
        </div>
      </div>
    </Dialog>
  )
}
