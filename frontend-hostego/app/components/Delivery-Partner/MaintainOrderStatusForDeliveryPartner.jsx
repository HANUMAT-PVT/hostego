'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Phone, Navigation, Clock, Check, ChevronDown, ChevronUp, ShoppingBag, ArrowRight, IndianRupee, User, Package, MapPin, AlertCircle } from 'lucide-react';
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

const AccordionSection = ({ title, icon: Icon, children, defaultOpen = false, count, badge }) => {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  return (
    <div className="bg-white">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full px-4 py-4 flex items-center justify-between hover:bg-gray-50/50 transition-colors"
      >
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-[var(--primary-color)]/10 rounded-xl flex items-center justify-center">
            <Icon className="w-5 h-5 text-[var(--primary-color)]" />
          </div>
          <div className="flex items-center gap-2">
            <span className="font-medium text-gray-900">{title}</span>
            {count && (
              <span className="px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 text-sm">
                {count}
              </span>
            )}
            {badge && (
              <span className="px-2 py-0.5 rounded-full bg-[var(--primary-color)]/10 text-[var(--primary-color)] text-sm font-medium">
                {badge}
              </span>
            )}
          </div>
        </div>
        <ChevronDown
          className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`}
        />
      </button>
      {isOpen && (
        <div className="p-4 animate-fade-in">
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

  // Add status filter options


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

  const findOutCorrectDeliveryPartnerCost = (orderItems) => {
    let totalCost = 0;
  
    for (let shop of orderItems) {
      for (let product of shop.shop_products) {
        const quantity = product.quantity || 0;
        const foodPrice = product.product_item?.food_price || 0;
  
        totalCost += quantity * foodPrice;
      }
    }
  
    return totalCost;
  };
  

  return (
    <div className="bg-white rounded-xl shadow-sm overflow-hidden">



      {/* Orders List */}


      {/* Sticky Header with Order Status */}
      <div className="sticky top-0 z-10 bg-gradient-to-r from-[var(--primary-color)] to-purple-600 p-4">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Package className="w-5 h-5 text-white" />
            <span className="text-white font-medium">#{activeOrder?.order_id}</span>
          </div>
          <span className="text-white/80 text-sm">{formatDate(activeOrder?.created_at)}</span>
        </div>

        {/* Order Status Badge */}
        <div className="flex items-center gap-2 bg-white/10 backdrop-blur-sm rounded-xl p-3">
          <div className="flex-1">
            <p className="text-white/80 text-sm">Current Status</p>
            <p className="text-white font-medium capitalize">
              {activeOrder?.order_status?.replace(/_/g, ' ')}
            </p>
          </div>
          <div className="h-10 w-[2px] bg-white/20"></div>
          <div className="flex-1 text-right">
            <p className="text-white/80 text-sm">Expected Earnings</p>
            <p className="text-white font-bold text-xl">₹{activeOrder?.delivery_partner_fee}</p>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-2 gap-2 p-3 bg-gray-50">
        <button
          onClick={() => window.location.href = `tel:${activeOrder?.user?.mobile_number}`}
          className="flex items-center justify-center gap-2 bg-[var(--primary-color)] text-white p-2 rounded-xl text-sm font-medium hover:opacity-90 active:scale-95 transition-all"
        >
          <Phone size={18} />
          Call Customer
        </button>

      </div>

      {/* Main Content Accordions */}
      <div className="divide-y">
        {/* Status Update Section */}
        <AccordionSection
          title="Order Status"
          icon={Clock}
          defaultOpen={false}
          badge={shouldShowSlider(activeOrder?.order_status) ? "Action Required" : null}
        >
          <div className="space-y-4">
            <StatusTimeLine ORDER_STATUSES={ORDER_STATUSES} activeOrder={activeOrder} />

            {shouldShowSlider(activeOrder?.order_status) && (
              <SliderStatusTracker
                text={`Slide to ${ORDER_STATUSES[getStatusStep(activeOrder.order_status) + 1].label}`}
                onConfirm={() => setIsConfirmationPopupOpen(true)}
              />
            )}
          </div>
        </AccordionSection>

        {/* Customer Details Section */}
        <AccordionSection
          title="Customer Details"
          icon={User}
          defaultOpen={false}
        >
          <div className="space-y-4">
            <div className="flex items-start gap-4 bg-blue-50 p-4 rounded-xl">
              <div className="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
                <User className="w-6 h-6 text-blue-600" />
              </div>
              <div className="flex-1">
                <h4 className="font-medium text-gray-900">
                  {activeOrder?.user?.first_name} {activeOrder?.user?.last_name}
                </h4>
                <p className="text-gray-600 mt-1">{activeOrder?.user?.mobile_number}</p>
                <div className="mt-3 bg-white p-3 rounded-lg border border-blue-100">
                  <p className="text-sm text-gray-600">Delivery Address</p>
                  <p className="text-gray-900 mt-1">{activeOrder?.address?.address_line_1}</p>
                </div>
              </div>
            </div>
          </div>
        </AccordionSection>

        {/* Order Items Section */}
        {activeOrder?.order_items?.map((shop) => (
          <AccordionSection
            key={shop?.shop_id}
            title={shop?.shop_name}
            icon={ShoppingBag}
            count={shop?.shop_products?.length}
            defaultOpen={false}
          >
            <div className="space-y-3">
              {shop?.shop_products?.map((product) => (
                <div
                  key={product?.product_id}
                  className="flex items-center gap-4 bg-white p-3 rounded-xl border border-gray-100"
                >
                  <img
                    src={product?.product_item?.product_img_url}
                    alt={product?.product_item?.product_name}
                    className="w-16 h-16 rounded-lg object-cover"
                  />
                  <div className="flex-1 min-w-0">
                    <h4 className="font-medium text-gray-900 truncate">
                      {product?.product_item?.product_name}
                    </h4>
                    <div className="flex items-center gap-2 mt-1">
                      <span className="text-sm text-gray-600">
                        {product?.quantity} × ₹{product?.product_item?.food_price}
                      </span>
                      <span className="text-sm font-medium text-[var(--primary-color)]">
                        ₹{product?.quantity*product?.product_item?.food_price}
                      </span>
                    </div>
                  </div>
                </div>
              ))}

              {activeOrder?.cooking_requests && (
                <div className="bg-yellow-50 p-4 rounded-xl border border-yellow-100">
                  <div className="flex items-start gap-3">
                    <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                    <p className="text-sm text-yellow-800 flex-1">
                      {activeOrder?.cooking_requests}
                    </p>
                  </div>
                </div>
              )}
            </div>
          </AccordionSection>

        ))}
        {/* Total Order value Section */}
        <AccordionSection

          title={"Total Order Amount"}
          icon={ShoppingBag}
          count={`₹ ${findOutCorrectDeliveryPartnerCost(activeOrder?.order_items||[])}`}

          defaultOpen={false}
        >
        
         
        </AccordionSection>
      </div>

      <ConfirmationPopup
        variant="info"
        title="Confirm Status Update"
        isOpen={isConfirmationPopupOpen}
        message={`Are you sure you want to update the order status to ${getNextStatus(activeOrder)?.label}?`}
        onConfirm={() => {
          const nextStatus = getNextStatus(activeOrder);
          if (nextStatus) onUpdateOrderStatus(activeOrder?.order_id, nextStatus?.id);
          setIsConfirmationPopupOpen(false);
        }}
        onCancel={() => setIsConfirmationPopupOpen(false)}
      />
    </div >
  );
};

export default MaintainOrderStatusForDeliveryPartner;
