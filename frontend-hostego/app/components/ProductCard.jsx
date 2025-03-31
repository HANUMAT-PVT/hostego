'use client'
import React, { useState } from 'react'
import axiosClient from '../utils/axiosClient'
import { Check, ShoppingCart, Loader2, Store, Leaf, Package, AlertCircle } from 'lucide-react'
import { useDispatch } from 'react-redux'
import { setFetchCartData } from '../lib/redux/features/user/userSlice'

const ProductCard = ({
    product_img_url,
    product_name,
    myKey,
    food_price,
    weight,
    product_id,
    isAlreadyInCart,
    food_category,
    shop
}) => {
    const [isInCart, setIsInCart] = useState(isAlreadyInCart);
    const [isLoading, setIsLoading] = useState(false)
    const dispatch = useDispatch()

    const isShopClosed = shop?.shop_status === 0

    const addProductInTheCart = async () => {
        try {
            setIsLoading(true)
            await axiosClient.post(`/api/cart/`, { product_id, quantity: 1 })
            setIsInCart(true)
            dispatch(setFetchCartData(true))
        } catch (error) {
            console.error("Error adding product in the cart:", error)
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div key={myKey} className='bg-white rounded-xl border border-gray-100 relative'>
            {/* Shop Closed Overlay */}
            {isShopClosed && (
                <div className="absolute right-3 top-3 z-6">
                    <div className="bg-red-500 text-white text-xs px-2 py-1 rounded-md shadow-sm">
                        Shop Closed
                    </div>
                </div>
            )}

            <div className='relative w-full flex justify-center p-3'>
                <div className='w-32 h-32 relative'>
                    <img
                        className={`rounded-md w-full h-full object-contain hover:scale-105 transition-transform duration-200 
                            ${isShopClosed ? 'opacity-60' : ''}`}
                        src={product_img_url}
                        alt={product_name}
                    />

                    {food_category?.is_veg === 1 ? <div className='absolute -top-1 -left-1 bg-green-100 p-1 rounded-full'>
                        <Leaf size={12} className={'text-green-600'} />
                    </div> : <div className='absolute -top-1 -left-1 bg-red-100 p-1 rounded-full'>
                        <Leaf size={12} className={'text-red-600'} />
                    </div>}
                </div>
            </div>

            <div className='p-3 pt-0'>
                {/* Product Name and Weight */}
                <div className='mb-2'>
                    <h3 className={`font-medium text-sm ${isShopClosed ? 'text-gray-400' : 'text-gray-800'} line-clamp-3`}>
                        {product_name}
                    </h3>
                    <div className='flex items-center gap-2 mt-1'>
                        <span className='text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 flex items-center gap-1'>
                            <Package size={10} />
                            {weight}
                        </span>
                    </div>
                </div>

                {/* Shop Info */}
                <div className='flex items-center justify-between mb-3'>
                    <div className='flex items-center gap-1.5'>
                        <Store size={12} className={isShopClosed ? 'text-red-400' : 'text-gray-500'} />
                        <span className={`text-xs ${isShopClosed ? 'text-red-400' : 'text-gray-600'}`}>
                            {shop?.shop_name}
                        </span>
                    </div>
                </div>

                {/* Price and Add to Cart */}
                <div className='flex items-center justify-between'>
                    <div>
                        <span className={`text-lg font-semibold ${isShopClosed ? 'text-gray-400' : 'text-gray-900'}`}>
                            ₹{food_price}
                        </span>
                    </div>

                    {isShopClosed ? (
                        <div className='flex items-center gap-1.5 text-red-500 text-sm'>
                            {/* <AlertCircle size={14} /> */}
                            {/* <span>Not Available</span> */}
                        </div>
                    ) : isInCart ? (
                        <div className='bg-green-100 text-sm text-green-600 px-3 py-1.5 rounded-lg font-medium border border-green-200 flex items-center gap-1.5'>
                            <Check size={14} className="stroke-2" />
                            <span>Added</span>
                        </div>
                    ) : (
                        <button
                            onClick={addProductInTheCart}
                            disabled={isLoading}
                            className={`text-sm px-3 py-1.5 rounded-lg font-medium 
                                transition-all duration-200 ease-in-out
                                ${isLoading
                                    ? 'bg-gray-100 text-gray-400'
                                    : 'bg-[var(--primary-color)] text-white hover:opacity-90 active:scale-95'
                                }
                                flex items-center gap-2`}
                        >
                            {isLoading ? (
                                <>
                                    <Loader2 size={14} className="animate-spin" />
                                    <span>Adding</span>
                                </>
                            ) : (
                                <>
                                    <ShoppingCart size={14} />
                                    <span>Add</span>
                                </>
                            )}
                        </button>
                    )}
                </div>
            </div>
        </div>
    )
}

export default ProductCard
