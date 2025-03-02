'use client'
import React, { useState, useEffect } from 'react'
import BackNavigationButton from '../components/BackNavigationButton'
import { Building2, Home, MapPin, Navigation } from 'lucide-react'
import axiosClient from '../utils/axiosClient'
import HostegoButton from '../components/HostegoButton'
import HostegoLoader from '../components/HostegoLoader'

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
    <div className="min-h-screen bg-[var(--bg-page-color)] ">
      <BackNavigationButton title="Address" />

      <form onSubmit={handleSubmit} className="p-2 space-y-4">


        {/* Address Details */}
        <div className="bg-white rounded-lg p-4 space-y-4 flex flex-col gap-4">
          <p className="text-xl font-normal ">
            Add New Address
          </p>
          {/* Address Line 1 */}
          <div className="relative">
            <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
              Address Type
            </label>
            <input
              value={addressData.address_type}
              onChange={(e) => setAddressData({ ...addressData, address_type: e.target.value })}
              className="w-full px-4 py-2 border-2 border-[var(--primary-color)] rounded-md outline-none  resize-none"
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

          {/* Submit Button */}
          <HostegoButton
            type="submit"
            text={isLoading ? 'Saving...' : 'Save Address'}
            className={`w-full bg-[var(--primary-color)] text-white py-3 rounded-lg font-medium
                        ${isLoading ? 'opacity-70 cursor-not-allowed' : 'hover:bg-opacity-90'}`}
          >

          </HostegoButton>
        </div>

      </form>
      {
        address?.map((el) =>
          <div key={el?.id} onClick={() => { }} className=' m-2 address-item flex items-center rounded-md   gap-4 cursor-pointer px-4 py-2 bg-white '>
            <div className='bg-[var(--bg-page-color)] p-2 w-[40px] h-[40px]  flex justify-center items-center rounded-full '>
              <Home size={20} className='text-[var(--primary-color)]' />
            </div>
            <div className=''>
              <p className='text-md font-semibold'>{el?.address_type}</p>
              <p className='text-sm '>{el?.address_line_1}</p>
            </div>
          </div>)
      }


    </div>
  )
}

export default page
