'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Phone, Navigation, Clock, Check, ChevronDown, ChevronUp, ShoppingBag, ArrowRight } from 'react-feather';
import { transformOrder } from '../../utils/helper'
import SliderStatusTracker from "./SliderStatusTracker"
import StatusTimeLine from '../Orders/StatusTimeLine';
import ConfirmationPopup from '../ConfirmationPopup';

export const ORDER_STATUSES = [
  {
    id: 'placed',
    label: 'Order Placed',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'assigned',
    label: 'Order Assigned',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'reached',
    label: 'Reached Shop',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'picked',
    label: 'Picked Up',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'on_the_way',
    label: 'On The Way',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'delivered',
    label: 'Delivered',
    icon: Check,
    color: 'var(--primary-color)'
  }
];





const MaintainOrderStatusForDeliveryPartner = ({ order, onUpdateOrderStatus }) => {
  const [activeOrder, setActiveOrder] = useState(transformOrder(order));
  const [isConfirmationPopupOpen, setIsConfirmationPopupOpen] = useState(false);
  const [isItemsExpanded, setIsItemsExpanded] = useState(false);



  useEffect(() => {
    setActiveOrder(transformOrder(order))
  }, [order])

  const getStatusStep = (status) => {
    const index = ORDER_STATUSES.findIndex(s => s.id === status);
    return index === -1 ? 0 : index;
  };

  const getNextStatus = (order) => {
    return ORDER_STATUSES[getStatusStep(order?.order_status) + 1];
  };


  return (
    <div className=" bg-[#F4F6FB]">
      {/* Main Content */}
      <div className="max-w-md mx-auto p-4" onClick={() => {
       
      }}>
        {/* Earnings Card */}
        <div className='rounded-xl px-4 py-4 text-center bg-white mb-2 flex flex-col gap-2'>
          <p className='font-semibold px-6 text-md  py-1 rounded-full border-4 border-[var(--primary-color)] w-fit m-auto'>Order !</p>
          <p className='font-semibold text-xl'>Expected earnings</p>
          <p className='font-bold text-3xl'>₹{activeOrder?.delivery_partner_fee}</p>

        </div>

        {/* Restaurant Details */}
        { <>
          <div className="bg-white rounded-xl shadow-sm mb-6">
            {activeOrder?.order_items?.map((shop) => <div key={shop?.shop_id}>


              <div className="p-4 border-b" >
                <div className="flex items-center justify-between mb-3">
                  <span className="px-3 py-1 bg-black text-white text-sm rounded-full">
                    #{activeOrder?.order_id?.slice(0, 8)}
                  </span>
                  <span className="text-sm font-medium text-[var(--primary-color)]">
                    {shop?.shop_products?.length} items
                  </span>
                </div>
                <h4 className="font-semibold text-lg mb-1">
                  {shop?.shop_name}
                </h4>
                <p className="text-gray-600 text-sm">
                  {shop?.address}
                </p>
                <div className="mt-2 text-sm text-gray-500">
                  <Clock className="w-4 h-4 inline mr-1" />
                  Prep Time: {shop?.preparation_time}
                </div>
              </div>
              {/* Order Items Accordion */}
              <div className="border-b">
                <button
                  className="w-full p-4 flex items-center justify-between text-left"
                  onClick={() => setIsItemsExpanded(!isItemsExpanded)}
                >
                  <div className="flex items-center gap-2">
                    <ShoppingBag className="w-5 h-5 text-[var(--primary-color)]" />
                    <span className="font-medium">Order Items</span>
                  </div>
                  {isItemsExpanded ? (
                    <ChevronUp className="w-5 h-5 text-gray-500" />
                  ) : (
                    <ChevronDown className="w-5 h-5 text-gray-500" />
                  )}
                </button>


              </div>

              {isItemsExpanded && (
                <div className="px-4 pb-4 mt-2">
                  {shop.shop_products.map((item, index) => (
                    <div key={index} className="mb-4 last:mb-0">
                      <div className="flex items-start justify-between">
                        <div className="flex items-start gap-2">
                          <div className={`w-2 h-2 rounded-full mt-2 ${item.product_item.food_category?.is_veg ? 'bg-green-500' : 'bg-red-500'
                            }`} />
                          <div>
                            <p className="font-medium">{item.product_item.product_name}</p>
                            <p className="text-sm text-gray-500">{item.product_item.description}</p>
                            <p className="text-sm text-gray-600 mt-1">
                              ₹{item.product_item.food_price} x {item.quantity}
                            </p>
                          </div>
                        </div>
                        <span className="font-medium">₹{item.sub_total}</span>
                      </div>
                    </div>


                  ))}

                  <div className="mt-4 pt-4 border-t space-y-2">
                    <div className="flex justify-between text-sm text-gray-600">
                      <span>Items Total</span>
                      <span>₹{shop?.shop_products?.reduce((acc, item) => acc + item.sub_total, 0)}</span>
                    </div>

                    <div className="flex justify-between text-sm text-gray-600">
                      <span>Platform Fee</span>
                      <span>₹{activeOrder?.platform_fee}</span>
                    </div>
                    <div className="flex justify-between text-sm text-gray-600">
                      <span>Delivery Fee</span>
                      <span>₹{activeOrder?.shipping_fee}</span>
                    </div>
                    <div className="flex justify-between font-medium pt-2 border-t">
                      <span>Total Amount</span>
                      <span>₹{activeOrder?.final_order_value}</span>
                    </div>
                    <div className="flex justify-between text-sm mt-2 text-gray-600 pt-2 border-2 border-[var(--primary-color)] rounded-md p-2">
                      <span>Cooking Requests</span>
                      <span className='text-[var(--primary-color)] '>{activeOrder?.cooking_requests}</span>
                    </div>
                  </div>
                </div >
              )}
            </div>)}



            <div className="p-4 flex gap-3">
              <button onClick={() => window.location.href = `tel:${activeOrder?.user?.mobile_number}`} className="flex-1 py-3 px-4 rounded-lg border-2 border-[var(--primary-color)] text-[var(--primary-color)] font-medium flex items-center justify-center gap-2">
                <Phone className="w-4 h-4" />
                Call
              </button>

            </div>

          </div >

          {/* Status Timeline */}
          <StatusTimeLine ORDER_STATUSES={ORDER_STATUSES} activeOrder={activeOrder} />
        </>}
      </div>

      {/* Bottom Action Bar */}
      <div className="z-10  left-0 right-0   p-4">
        <div className="max-w-md mx-auto">
          {getStatusStep(activeOrder?.order_status) < ORDER_STATUSES.length - 1 && (
            <SliderStatusTracker
              text={`Slide to ${ORDER_STATUSES[getStatusStep(activeOrder.order_status) + 1].label}`}
              onConfirm={() => {
                setIsConfirmationPopupOpen(true)
              }}
            />
          )}
        </div>
      </div>
      <ConfirmationPopup
        variant="info"
        title="Confirm Order Status"
        isOpen={isConfirmationPopupOpen}
        message="Are you sure you want to update the order status?"
        onConfirm={() => {
          const nextStatus = getNextStatus(activeOrder);
          if (nextStatus) onUpdateOrderStatus(activeOrder?.order_id, nextStatus?.id);
          setIsConfirmationPopupOpen(false)
        }}
        onCancel={() => {
          setIsConfirmationPopupOpen(false)
        }}
      />
    </div >
  );
};

export default MaintainOrderStatusForDeliveryPartner;
