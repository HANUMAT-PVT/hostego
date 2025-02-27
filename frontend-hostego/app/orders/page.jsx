"use client"
import React, { useEffect, useState } from 'react'
import BackNavigationButton from "../components/BackNavigationButton"
import OrderPreviewCard from "../components/Orders/OrderPreviewCard"
import axiosClient from '../utils/axiosClient'
import HostegoLoader from '../components/HostegoLoader'
import { PackageX } from 'lucide-react'
import { useRouter } from 'next/navigation'
import HostegoButton from '../components/HostegoButton'

const NoOrders = () => {
  const router = useRouter()
  return (
    <div className="flex flex-col items-center justify-center p-8 text-center">
      <div className="w-24 h-24 bg-gray-50 rounded-full flex items-center justify-center mb-6">
        <PackageX className="w-12 h-12 text-gray-400" />
      </div>
      <h3 className="text-lg font-semibold text-gray-800 mb-2">No Orders Yet</h3>
      <p className="text-gray-600 mb-6">
        Looks like you haven't placed any orders yet.
        Start shopping to see your orders here!
      </p>
      <HostegoButton
        onClick={() => router.push('/home')}
        text={"Start Shopping"}
      >

      </HostegoButton>
    </div>
  )
}

const OrdersPage = () => {
  const [orders, setOrders] = useState([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    fetchOrders()
  }, [])

  const fetchOrders = async () => {
    try {
      setIsLoading(true)
      const { data } = await axiosClient.get('/api/order')
      setOrders(data)
    } catch (error) {
      console.error('Error fetching orders:', error)
    } finally {
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return <HostegoLoader />
  }

  return (
    <div className="min-h-screen bg-[var(--bg-page-color)]">
      <BackNavigationButton title="My Orders" />

      <div className="pb-4">
        {orders.length === 0 ? (
          <NoOrders />
        ) : (
          orders.map((order) => (
            <OrderPreviewCard key={order.order_id} order={order} />
          ))
        )}
      </div>
    </div>
  )
}

export default OrdersPage
