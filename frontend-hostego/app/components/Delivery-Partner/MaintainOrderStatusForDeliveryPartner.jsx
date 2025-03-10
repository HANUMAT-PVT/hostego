'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Phone, Navigation, Clock, Check, ChevronDown, ChevronUp, ShoppingBag, ArrowRight, IndianRupee, User, Package } from 'lucide-react';
import { transformOrder } from '../../utils/helper'
import SliderStatusTracker from "./SliderStatusTracker"
import StatusTimeLine from '../Orders/StatusTimeLine';
import ConfirmationPopup from '../ConfirmationPopup';
import { formatDate } from '../../utils/helper';

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
    id: 'reached_door',
    label: 'Reached Door',
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

const AccordionSection = ({ title, icon: Icon, children, defaultOpen = false, count, status }) => {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  return (
    <div className="border-b last:border-b-0">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full px-4 py-3 flex items-center justify-between hover:bg-gray-50"
      >
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-[var(--primary-color)]/10 rounded-lg flex items-center justify-center">
            <Icon className="w-4 h-4 text-[var(--primary-color)]" />
          </div>
          <span className="font-medium text-gray-900">{title}</span>
        </div>
        <ChevronDown className={`w-5 h-5 text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      {isOpen && (
        <div className="p-4 bg-gray-50 animate-fade-in">
          {children}
        </div>
      )}
    </div>
  );
};

const MaintainOrderStatusForDeliveryPartner = ({ order, onUpdateOrderStatus }) => {
  const [activeOrder, setActiveOrder] = useState(transformOrder(order));
  const [isConfirmationPopupOpen, setIsConfirmationPopupOpen] = useState(false);

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

  const shouldShowSlider = (status) => {
    return !['reached_door', 'delivered', 'cancelled'].includes(status?.toLowerCase());
  };

  return (
    <AccordionSection title="Order Details" icon={Package} defaultOpen={true}>
      <div className="bg-white rounded-xl shadow-sm overflow-hidden mb-4">
        {/* Order Header - Always Visible */}
        <div className="p-4 bg-gradient-to-r from-[var(--primary-color)]   to-purple-600 text-white">
          <div className="flex flex-col items-start gap-1">
            <span className="bg-white/20 px-3 py-1 rounded-full text-sm">
              Order #{activeOrder?.order_id?.slice(-6)}
            </span>
            <span className="text-sm flex items-center gap-1">
              <Clock size={14} />
              {formatDate(activeOrder?.created_at)}
            </span>
          </div>
        </div>

        {/* Earnings Section */}

        <div className="bg-[var(--primary-color)]/5 rounded-xl p-4 text-center">
          <div className="inline-block bg-white px-6 py-1 rounded-full border-2 border-[var(--primary-color)] text-[var(--primary-color)] font-medium mb-3">
            Active Order
          </div>
          <h3 className="text-gray-600 font-medium mb-2">Expected Earnings</h3>
          <p className="text-3xl font-bold text-[var(--primary-color)]">
            ₹{activeOrder?.delivery_partner_fee}
          </p>
        </div>

        {/* Customer Section */}
        <AccordionSection title="Customer Details" icon={User} defaultOpen={false}>
          <div className="space-y-4">
            <div className="flex flex-col items-start gap-2">
              <div>
                <p className="text-sm text-gray-500">Customer Name</p>
                <p className="font-medium">{activeOrder?.user?.first_name} {activeOrder?.user?.last_name}</p>
              </div>
              <button
                onClick={() => window.location.href = `tel:${activeOrder?.user?.mobile_number}`}
                className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)] text-white font-medium hover:opacity-90 transition-opacity"
              >
                <Phone size={16} />
                Call {activeOrder?.user?.first_name}
              </button>
            </div>
            <div className="bg-gray-100 p-3 rounded-lg">
              <p className="text-sm text-gray-600 mb-2">Delivery Address:</p>
              <p className="font-medium">{activeOrder?.address?.address_line_1}</p>
            </div>
          </div>
        </AccordionSection>

        {/* Order Status Section */}
        <AccordionSection title="Order Status" icon={Clock} defaultOpen={false}>
          <StatusTimeLine ORDER_STATUSES={ORDER_STATUSES} activeOrder={activeOrder} />
          <div className="p-4 bg-gray-50 border-t">
            {shouldShowSlider(activeOrder?.order_status) && (
              <SliderStatusTracker
                text={`Slide to ${ORDER_STATUSES[getStatusStep(activeOrder.order_status) + 1].label}`}
                onConfirm={() => setIsConfirmationPopupOpen(true)}
              />
            )}
          </div>
        </AccordionSection>

        {/* Order Items Section */}
        {activeOrder?.order_items?.map((shop, index) => (
          <AccordionSection
            key={shop?.shop_id}
            title={`${shop?.shop_name} (${shop?.shop_products?.length} items)`}
            icon={ShoppingBag}
          >
            <div className="space-y-4">
              {shop?.shop_products?.map((product) => (
                <div
                  key={product?.product_id}
                  className="flex items-center gap-3 bg-white p-3 rounded-lg shadow-sm"
                >
                  <img
                    src={product?.product_item?.product_img_url}
                    alt={product?.product_item?.product_name}
                    className="w-12 h-12 rounded-lg object-cover"
                  />
                  <div className="flex-1">
                    <h4 className="font-medium">{product?.product_item?.product_name}</h4>
                    <p className="text-sm font-medium text-gray-500">
                      ₹{product?.product_item?.food_price} ×  {product?.quantity}
                    </p>

                  </div>
                  <span className="font-medium">₹{product?.sub_total}</span>
                </div>
              ))}

              {activeOrder?.cooking_requests && (
                <div className="bg-yellow-50 p-3 rounded-lg border border-yellow-100">
                  <p className="text-sm text-yellow-800">
                    {activeOrder?.cooking_requests}
                  </p>
                </div>
              )}
            </div>

          </AccordionSection>
        ))}

        {/* Bottom Actions */}


        <ConfirmationPopup
          variant="info"
          title="Confirm Order Status"
          isOpen={isConfirmationPopupOpen}
          message="Are you sure you want to update the order status?"
          onConfirm={() => {
            const nextStatus = getNextStatus(activeOrder);
            if (nextStatus) onUpdateOrderStatus(activeOrder?.order_id, nextStatus?.id);
            setIsConfirmationPopupOpen(false);
          }}
          onCancel={() => setIsConfirmationPopupOpen(false)}
        />
      </div>
    </AccordionSection>
  );
};

export default MaintainOrderStatusForDeliveryPartner;
