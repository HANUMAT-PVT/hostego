"use client"
import React, { useEffect, useState } from 'react'
import BackNavigationButton from "../components/BackNavigationButton"
import OrderPreviewCard from "../components/Orders/OrderPreviewCard"
import axiosClient from '../utils/axiosClient'
const page = () => {

  const [orders, setOrders] = useState([])

  useEffect(() => {
    fetchOrders()
  }, [])

  const fetchOrders = async () => {
    const { data } = await axiosClient.get('/api/order')
    console.log(data, "orders")
    setOrders(data)
  }

  return (
    <div className='bg-[#F4F6FB]'>
      <BackNavigationButton title={"Your orders"} />
      <div className='p-2 flex flex-col gap-2'>
        {orders?.map((el) => (
          <OrderPreviewCard key={el?.order_id} order={el} />
        ))}
        
      </div>

    </div>
  )
}

export default page
