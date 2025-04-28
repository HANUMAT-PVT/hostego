import React, { useState, useCallback } from 'react'
import axiosClient from '../../utils/axiosClient.js'
import { debounce } from 'lodash'
import { Minus, Plus } from 'lucide-react'

const CartItem = ({ product_item, quantity, cart_item_id, fetchCartAgain }) => {
    const [cartItemQuantity, setCartItemQuantity] = useState(quantity || 1)

    const debouncedUpdateQuantity = useCallback(

        debounce(async (cartItemId, quantity) => {
            try {
                await axiosClient.patch(`/api/cart/${cartItemId}`, { quantity })
                fetchCartAgain()
            } catch (error) {
                console.error('Error updating cart item quantity:', error)
                fetchCartAgain()
            }
        }, 500),
        []
    )

    const updateCartItemQuantity = (cart_item_id, quantity) => {
        if (quantity < 0 || quantity > 20) return

        setCartItemQuantity(quantity)
        debouncedUpdateQuantity(cart_item_id, quantity)
    }

    return (
        <div className='mx-2 bg-white rounded-xl overflow-hidden shadow-sm'>
            <div className='p-3 flex items-center justify-between  gap-2'>
                {/* Product Info */}
                <div className='flex gap-3 items-center'>
                    {/* Image Container */}
                    <div className='relative group'>
                        <div className='w-[55px] h-[55px] rounded-lg overflow-hidden bg-[var(--bg-page-color)]'>
                            <img
                                className='w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-200'
                                src={product_item?.product_img_url}
                                alt={product_item?.product_name}
                            />
                        </div>
                    </div>

                    {/* Product Details */}
                    <div className='flex flex-col justify-between py-1 '>
                        <div>
                            <h3 className='font-medium text-gray-800 leading-snug mb-1 text-sm'>
                                {product_item?.product_name} <span className='text-gray-600 text-xs'>( {product_item?.shop?.shop_name} )</span>
                            </h3>
                            <p className='text-sm font-semibold text-gray-500'>
                                ₹{product_item?.selling_price} × {quantity}
                            </p>
                        </div>
                    </div>
                </div>

                {/* Price and Quantity Controls */}
                <div className='flex flex-col items-end gap-2'>
                    {/* Quantity Controls */}
                    <div className='flex items-center gap-1 bg-[var(--bg-page-color)] rounded-lg p-1'>
                        <button
                            onClick={() => updateCartItemQuantity(cart_item_id, cartItemQuantity - 1)}
                            className='w-7 h-7 flex items-center justify-center rounded-md  bg-[var(--primary-color)] 
                                      hover:text-white transition-colors duration-200 text-white
                                     active:scale-95 transform'

                        >
                            <Minus size={16} />
                        </button>
                        <span className='w-8 text-center font-medium'>
                            {cartItemQuantity}
                        </span>
                        <button
                            onClick={() => updateCartItemQuantity(cart_item_id, cartItemQuantity + 1)}
                            className='w-7 h-7 flex items-center justify-center rounded-md  text-white bg-[var(--primary-color)]
                                     hover:bg-[var(--primary-color)] hover:text-white transition-colors duration-200
                                     active:scale-95 transform'
                        >
                            <Plus size={16} />
                        </button>
                    </div>

                    {/* Total Price */}
                    <p className='font-semibold text-[var(--primary-color)]'>
                        ₹{product_item?.selling_price * cartItemQuantity}
                    </p>
                </div>
            </div>
        </div>
    )
}

export default CartItem
