import { ChevronDown, ChevronUp, ShoppingBag } from 'lucide-react';
import React,{useState} from 'react'

const OrderItemAccordian = ({ order }) => {
    const [isItemsExpanded, setIsItemsExpanded] = useState(false);
    return (
        <div className="border-b">
            {/* Order Items Accordion */}
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

            {isItemsExpanded && (
                <div className="px-4 pb-4">
                    {order?.order_items?.map((item, index) => (
                        <div key={index} className="mb-4 last:mb-0">
                            <div className="flex items-start justify-between">
                                <div className="flex items-start gap-2">
                                    <div className={`w-2 h-2 rounded-full mt-2 ${item?.product_item?.food_category?.is_veg ? 'bg-green-500' : 'bg-red-500'
                                        }`} />
                                    <div>
                                        <p className="font-medium">{item?.product_item?.product_name}</p>
                                        <p className="text-sm text-gray-500">{item?.product_item?.description}</p>
                                        <p className="text-sm text-gray-600 mt-1">
                                            ₹{item?.product_item?.food_price} x {item?.quantity}
                                        </p>
                                    </div>
                                </div>
                                <span className="font-medium">₹{item?.sub_total}</span>
                            </div>
                        </div>
                    ))}

                    <div className="mt-4 pt-4 border-t space-y-2">
                        <div className="flex justify-between text-sm text-gray-600">
                            <span>Items Total</span>
                            <span>₹{order?.order_items?.reduce((acc, item) => acc + item?.sub_total, 0)}</span>
                        </div>
                        <div className="flex justify-between text-sm text-gray-600">
                            <span>Platform Fee</span>
                            <span>₹{order?.platform_fee}</span>
                        </div>
                        <div className="flex justify-between text-sm text-gray-600">
                            <span>Delivery Fee</span>
                            <span>₹{order?.shipping_fee}</span>
                        </div>
                        <div className="flex justify-between font-medium pt-2 border-t">
                            <span>Total Amount</span>
                            <span>₹{order?.final_order_value}</span>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default OrderItemAccordian
