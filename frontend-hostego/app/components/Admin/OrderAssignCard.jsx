import { CheckCircle, ChevronDown, ChevronUp, IndianRupee, MapPin, User, Package } from 'lucide-react'
import { useState } from 'react'
import { formatDate } from '@/app/utils/helper'
import HostegoButton from '../HostegoButton'

const OrderAssignCard = ({ order, selectOrderItem, selectedOrderItem }) => {

    const [isExpanded, setIsExpanded] = useState(false)


    return (
        <div onClick={() => selectOrderItem(order)} className={`bg-white rounded-xl overflow-hidden shadow-sm border border-2  ${order.order_id == selectedOrderItem?.order_id
            ? 'border-[var(--primary-color)] bg-[var(--primary-color)]/5'
            : 'border-gray-100 bg-white hover:border-[var(--primary-color)]/30'}`}>
            {/* Order Header */}
            <div className="p-4 border-b bg-gradient-to-r from-[var(--primary-color)] to-purple-600">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <Package className="w-5 h-5 text-white" />
                        <span className="text-white font-medium">
                            #{order?.order_id}
                        </span>
                    </div>
                    <span className="px-3 py-1 bg-white/20 backdrop-blur-sm text-white rounded-full text-sm">
                        {formatDate(order?.created_at)}
                    </span>
                </div>
            </div>

            {/* Order Summary */}
            <div className="p-4">
                <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-[var(--primary-color)]/10 flex items-center justify-center">
                            <IndianRupee className="w-5 h-5 text-[var(--primary-color)]" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-600">Order Value</p>
                            <p className="text-xl font-semibold">₹{order?.final_order_value}</p>
                        </div>
                    </div>
                    <button
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-gray-100 text-gray-700 hover:bg-gray-200 transition-colors"
                    >
                        {isExpanded ? 'Show Less' : 'View Items'}
                        {isExpanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
                    </button>
                </div>

                {/* Customer Details */}
                <div className="flex items-center gap-3 mb-4">
                    <div className="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center">
                        <User className="w-5 h-5 text-blue-600" />
                    </div>
                    <div>
                        <p className="text-sm text-gray-600">Customer Name</p>
                        <p className="font-medium">{order?.user?.first_name} {order?.user?.last_name}</p>
                    </div>

                    <div>
                        <p className="text-sm text-gray-600">Mobile Number</p>
                        <p className="font-medium">{order?.user?.mobile_number}</p>
                    </div>
                </div>

                {/* Delivery Address */}
                <div className="flex items-start gap-3 mb-4">
                    <div className="w-10 h-10 rounded-full bg-green-50 flex items-center justify-center">
                        <MapPin className="w-5 h-5 text-green-600" />
                    </div>
                    <div>
                        <p className="text-sm text-gray-600">Delivery Address</p>
                        <p className="font-medium"> {order?.address?.address_type || ""} {order?.address?.address_line_1 || 'Address not available'}</p>
                    </div>
                </div>

                {/* Order Items */}
                {isExpanded && (
                    <div className="mt-4 pt-4 border-t space-y-3">
                        <p className="font-medium text-gray-700">Order Items</p>
                        {order?.order_items?.map((item, index) => (
                            <div key={index} className="flex items-center gap-3 bg-gray-50 p-3 rounded-lg">
                                <img
                                    src={item?.product_item?.product_img_url}
                                    alt={item?.product_item?.product_name}
                                    className="w-12 h-12 rounded-lg object-cover"
                                />
                                <div className="flex-1">
                                    <p className="font-medium">{item?.product_item?.product_name} <span className='text-gray-600 text-xs'>( {item?.product_item?.shop?.shop_name} )</span></p>
                                    <p className="text-sm text-gray-600">
                                        {item?.quantity} × ₹{item?.product_item?.selling_price}
                                    </p>
                                </div>
                                <p className="font-medium">₹{item?.sub_total}</p>
                            </div>
                        ))}
                    </div>
                )}



            </div>
        </div>
    )
}


export default OrderAssignCard