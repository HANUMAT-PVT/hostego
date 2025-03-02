'use client'
import React, { useState } from 'react'
import axiosClient from '../utils/axiosClient'
import { Check, ShoppingCart, Loader2 } from 'lucide-react'
import { useDispatch } from 'react-redux'
import { setFetchCartData } from '../lib/redux/features/user/userSlice'

const ProductCard = ({ product_img_url, product_name, myKey, tags, food_price, weight, product_id }) => {
    const [isInCart, setIsInCart] = useState(false)
    const [isLoading, setIsLoading] = useState(false)
    const dispatch = useDispatch()
    
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
        <div key={myKey} className='product-item flex flex-col gap-2  p-2'>
            <div className='relative p-2 flex flex-col items-center group'>
                <img className='w-[90px] h-[90px]'
                    src={product_img_url}
                    alt={product_name}
                />

                {isInCart ? (
                    <div className='absolute bottom-0 right-0 bg-green-100 text-sm text-green-600 px-2 py-1 rounded-md font-medium border border-green-200 flex items-center gap-1'>
                        <ShoppingCart size={12} className="stroke-2" />
                        <span>In Cart</span>
                    </div>
                ) : (
                    <button
                        onClick={addProductInTheCart}
                        disabled={isLoading}
                        className={`absolute bottom-0 right-0 bg-white text-sm text-green-600 px-2 py-1 border-2 border-green-600 rounded-md font-semibold 
                            transition-all duration-200 ease-in-out
                            hover:bg-green-600 hover:text-white hover:scale-105
                            active:scale-95
                            disabled:opacity-50 disabled:cursor-not-allowed
                            flex items-center gap-1 min-w-[60px] justify-center`}
                    >
                        {isLoading ? (
                            <>
                                <Loader2 size={14} className="animate-spin" />
                                <span>Adding</span>
                            </>
                        ) : (
                            'ADD'
                        )}
                    </button>
                )}
            </div>

            <div className='flex gap-2'>
                <p className='text-[10px] px-2 py-1 rounded-lg bg-gray-200 font-medium text-gray-800'>{weight}</p>
            </div>
            <p className='product-title text-sm font-medium'>{product_name}</p>
            <p className='product-price text-left text-sm font-semibold'>₹ {food_price}</p>
        </div>
    )
}

export default ProductCard
