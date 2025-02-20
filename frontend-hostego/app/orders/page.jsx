"use client"
import React from 'react'
import BackNavigationButton from "../components/BackNavigationButton"
import OrderPreviewCard from "../components/Orders/OrderPreviewCard"
const page = () => {
  return (
    <div className='bg-[#F4F6FB]'>
      <BackNavigationButton title={"Your orders"} />
      <div className='p-2 flex flex-col gap-2'>
        <OrderPreviewCard myKey={1} />
        <OrderPreviewCard myKey={2} />
        <OrderPreviewCard myKey={3} />
      </div>

    </div>
  )
}

export default page
