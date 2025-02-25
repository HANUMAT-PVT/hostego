'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Phone, Navigation, Clock, Check, ChevronDown, ChevronUp, ShoppingBag, ArrowRight } from 'react-feather';
import { transformOrder } from '../../utils/helper'
import SliderStatusTracker from "./SliderStatusTracker"
import StatusTimeLine from '../Orders/StatusTimeLine';

export const ORDER_STATUSES = [
  {
    id: 'pending',
    label: 'Order Placed',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'reached_shop',
    label: 'Reached Restaurant',
    icon: Check,
    color: 'var(--primary-color)'
  },
  {
    id: 'picked_up',
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

const DUMMY_ORDERS = {
  order_id: "e3571ddc-3b87-45ae-a115-e793e4bf3e24",
  user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
  created_at: "2025-02-17 01:53:58.886302+05:30",
  updated_at: "2025-02-17 17:56:28.348765+05:30",
  order_items: [
    {
      user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
      quantity: 2,
      sub_total: 500,
      product_id: "f9fca8c1-dfb2-4ef0-b5de-9b388b3a3df3",
      cart_item_id: "e4c1e07f-b504-484c-a1cc-aa0ca246d647",
      product_item: {
        shop: {
          address: "123 MG Road, New Delhi, India",
          shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
          shop_img: "img_url",
          shop_name: "PUNJABI Delights",
          shop_status: 1,
          food_category: {
            is_veg: 0,
            is_cooked: 0,
          },
          preparation_time: "30 min",
        },
        tags: null,
        shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
        discount: {
          percentage: 10,
          is_available: 0,
        },
        created_at: "2025-02-13T18:48:33.882374+05:30",
        food_price: 250,
        product_id: "f9fca8c1-dfb2-4ef0-b5de-9b388b3a3df3",
        updated_at: "2025-02-13T18:48:33.882374+05:30",
        description:
          "A creamy and rich paneer dish cooked in butter with tomato-based gravy.",
        availability: 1,
        product_name: "Paneer Butter Masala",
        food_category: {
          is_veg: 0,
          is_cooked: 0,
        },
        product_img_url: "https://example.com/paneer.jpg",
        preparation_time: "30 min",
      },
    },
    {
      user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
      quantity: 7,
      sub_total: 1750,
      product_id: "e806a9c8-3136-4786-824c-1bc4e9b287e1",
      cart_item_id: "85cb1417-72d5-4b1d-bc1c-4edc9aae5e2c",
      product_item: {
        shop: {
          address: "123 MG Road, New Delhi, India",
          shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
          shop_img: "img_url",
          shop_name: "PUNJABI Delights",
          shop_status: 1,
          food_category: {
            is_veg: 0,
            is_cooked: 0,
          },
          preparation_time: "30 min",
        },
        tags: null,
        shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
        discount: {
          percentage: 10,
          is_available: 0,
        },
        created_at: "2025-02-13T18:58:38.914894+05:30",
        food_price: 250,
        product_id: "e806a9c8-3136-4786-824c-1bc4e9b287e1",
        updated_at: "2025-02-13T18:58:38.914894+05:30",
        description:
          "A creamy and rich paneer dish cooked in butter with tomato-based gravy.",
        availability: 1,
        product_name: "Paneer Butter Masala",
        food_category: {
          is_veg: 0,
          is_cooked: 0,
        },
        product_img_url: "https://example.com/paneer.jpg",
        preparation_time: "30 min",
      },
    },
  ],
  platform_fee: 1,
  shipping_fee: 30,
  final_order_value: 2281,
  delivery_partner_fee: 21,
  payment_transaction_id: "c4b66d17-d9a7-4a07-83a4-b5a135d5dff9",
  order_status: "pending",
  delivery_partner_id: "NULL",
  delivered_at: "0001-01-01 05:53:28+05:53:28",
  delivery_partner: {
    user: {
      email: "johndoe@example.com",
      user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
      last_name: "",
      created_at: "2025-02-17T00:39:39.655141+05:30",
      first_name: "John Doe",
      updated_at: "2025-02-17T00:39:39.655529+05:30",
      mobile_number: "+91-9876543211",
      last_login_timestamp: "2025-02-17T00:39:39.655141+05:30",
      firebase_otp_verified: 1,
    },
    address: "Zakir-A room number 1115",
    user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
    documents: {
      upi_id: "",
      aadhaar_back_img: "",
      bank_details_img: "",
      aadhaar_front_img: "",
    },
    account_status: 0,
    partner_img_url: "",
    availability_status: 0,
    delivery_partner_id: "5d8297ec-ad2a-443e-a9a1-a8e394107204",
  },
};



const MaintainOrderStatusForDeliveryPartner = () => {
  const [orders, setOrders] = useState(transformOrder(DUMMY_ORDERS));

  const [activeOrder, setActiveOrder] = useState(orders);
  const [isItemsExpanded, setIsItemsExpanded] = useState(false);

  const updateOrderStatus = async (orderId, newStatus) => {

    const updatedOrders = { ...activeOrder }
    updatedOrders.order_status = newStatus

    setOrders(updatedOrders);
    setActiveOrder(updatedOrders)
  };

  const getStatusStep = (status) => {
    const index = ORDER_STATUSES.findIndex(s => s.id === status);
    return index === -1 ? 0 : index;
  };

  const getNextStatus = (order) => {
    return ORDER_STATUSES[getStatusStep(order.order_status) + 1];
  };

 

  return (
    <div className="min-h-screen bg-[#F4F6FB]">
      {/* Main Content */}
      <div className="max-w-md mx-auto p-4">
        {/* Earnings Card */}
        <div className='rounded-xl px-4 py-4 text-center bg-white mb-2 flex flex-col gap-2'>
          <p className='font-semibold px-6 text-md  py-1 rounded-full border-4 border-[var(--primary-color)] w-fit m-auto'>Order !</p>
          <p className='font-semibold text-xl'>Expected earnings</p>
          <p className='font-bold text-3xl'>₹{activeOrder?.delivery_partner_fee}</p>

        </div>

        {/* Restaurant Details */}

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
                </div>
              </div >
            )}
          </div>)}



          <div className="p-4 flex gap-3">
            <button className="flex-1 py-3 px-4 rounded-lg border-2 border-[var(--primary-color)] text-[var(--primary-color)] font-medium flex items-center justify-center gap-2">
              <Phone className="w-4 h-4" />
              Call
            </button>
            <button className="flex-1 py-3 px-4 rounded-lg bg-[var(--primary-color)] text-white font-medium flex items-center justify-center gap-2">
              <Navigation className="w-4 h-4" />
              Navigate
            </button>
          </div>
        </div >

        {/* Status Timeline */}
        <StatusTimeLine ORDER_STATUSES={ORDER_STATUSES} activeOrder={activeOrder} />

      </div>

      {/* Bottom Action Bar */}
      <div className="z-10 fixed bottom-0 left-0 right-0 bg-white border-t p-4">
        <div className="max-w-md mx-auto">
          {getStatusStep(activeOrder.order_status) < ORDER_STATUSES.length - 1 && (
            <SliderStatusTracker
              text={`Slide to ${ORDER_STATUSES[getStatusStep(activeOrder.order_status) + 1].label}`}
              onConfirm={() => {
                const nextStatus = getNextStatus(activeOrder);
                if (nextStatus) updateOrderStatus(activeOrder.order_id, nextStatus.id);
              }}
            />
          )}
        </div>
      </div>
    </div >
  );
};

export default MaintainOrderStatusForDeliveryPartner;
