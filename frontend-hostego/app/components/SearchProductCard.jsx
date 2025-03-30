'use client'
import React, { useState } from 'react'
import { Check, ShoppingCart, Loader2, Store, Leaf, Package } from 'lucide-react'
import { useDispatch } from 'react-redux'
import { setFetchCartData } from '../lib/redux/features/user/userSlice'
import axiosClient from '../utils/axiosClient'

const SearchProductSkeleton = () => {
    return (
        <div className='bg-white rounded-xl border border-gray-100 p-3 flex gap-4 animate-pulse'>
            {/* Image Skeleton */}
            <div className='relative w-24 h-24 flex-shrink-0 bg-gray-200 rounded-lg'></div>

            {/* Content Skeleton */}
            <div className='flex-1 flex flex-col justify-between'>
                <div>
                    {/* Shop Name Skeleton */}
                    <div className='mb-1.5'>
                        <div className='h-5 bg-gray-200 rounded w-3/4 mb-2'></div>
                        <div className='h-5 w-20 bg-gray-200 rounded'></div>
                    </div>

                    {/* Product Name Skeleton */}
                    <div className='h-4 bg-gray-200 rounded w-1/2 mb-2'></div>
                </div>

                {/* Price and Button Skeleton */}
                <div className='flex items-center justify-between'>
                    <div className='h-6 bg-gray-200 rounded w-16'></div>
                    <div className='h-9 bg-gray-200 rounded w-24'></div>
                </div>
            </div>
        </div>
    )
}

const EmptyState = ({ searchValue, sendSearchValue }) => {
    return (
        <div className="flex flex-col items-center justify-center py-12">
            {/* Empty State Illustration */}
            <div className="w-48 h-48 mb-6 relative">
                <img
                    src="https://cdn-icons-png.flaticon.com/512/7486/7486744.png"
                    alt="No results"
                    className="w-full h-full object-contain opacity-50"
                />
            </div>

            {/* Empty State Message */}
            <div className="text-center space-y-2">
                <h3 className="text-lg font-semibold text-gray-800">
                    {searchValue
                        ? `No results found for "${searchValue}"`
                        : "Start searching for products"
                    }
                </h3>
                <p className="text-gray-500 text-sm max-w-md">
                    {searchValue
                        ? "Try searching with a different keyword or browse our categories"
                        : "Search for your favorite food items, snacks, and more"
                    }
                </p>
            </div>

            {/* Suggested Keywords (when no results) */}
            {searchValue && (
                <div className="mt-6">
                    <p className="text-sm text-gray-600 mb-3">Popular searches:</p>
                    <div className="flex flex-wrap gap-2 justify-center">
                        {['Samosa', 'Maggi', 'Chips', 'Thali', 'Snacks'].map((keyword) => (
                            <button
                                key={keyword}
                                className="px-4 py-2 rounded-full border border-gray-200 text-sm text-gray-600 hover:border-[var(--primary-color)] hover:text-[var(--primary-color)] transition-colors"
                                onClick={() => {
                                    // You can add click handler here to trigger search
                                    // handleSearch(keyword);
                                    sendSearchValue(keyword);
                                }}
                            >
                                {keyword}
                            </button>
                        ))}
                    </div>
                </div>
            )}
        </div>
    )
}

const SearchProductCard = ({
    product_img_url,
    product_name,
    food_price,
    weight,
    product_id,
    isAlreadyInCart,
    food_category,
    shop
}) => {
    const [isInCart, setIsInCart] = useState(isAlreadyInCart);
    const [isLoading, setIsLoading] = useState(false);
    const dispatch = useDispatch();

    const isShopClosed = shop?.shop_status === 0;

    const addProductInTheCart = async () => {
        try {
            setIsLoading(true);
            await axiosClient.post(`/api/cart/`, { product_id, quantity: 1 });
            setIsInCart(true);
            dispatch(setFetchCartData(true));
        } catch (error) {
            console.error("Error adding product in the cart:", error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className='bg-white rounded-xl border border-gray-100 p-3 flex gap-4'>
            {/* Product Image Section */}
            <div className='relative w-24 h-24 flex-shrink-0'>
                <img
                    className={`w-full h-full rounded-lg object-cover ${isShopClosed ? 'opacity-60' : ''}`}
                    src={product_img_url}
                    alt={product_name}
                />

                {food_category?.is_veg === 1 ? <div className='absolute -top-1 -left-1 bg-green-100 p-1 rounded-full'>
                    <Leaf size={12} className={'text-green-600'} />
                </div> : <div className='absolute -top-1 -left-1 bg-red-100 p-1 rounded-full'>
                    <Leaf size={12} className={'text-red-600'} />
                </div>}

                {isShopClosed && (
                    <div className="absolute -top-2 -right-2">
                        <div className="bg-red-500 text-white text-xs px-2 py-0.5 rounded-md">
                            Closed
                        </div>
                    </div>
                )}
            </div>

            {/* Product Details Section */}
            <div className='flex-1 flex flex-col justify-between'>
                <div>
                    {/* Product Name and Weight */}
                    <div className='mb-1.5'>
                        <h3 className={`text-gray-800 line-clamp-2 font-semibold ${isShopClosed ? 'text-gray-400' : ''}`}>
                            {shop?.shop_name}
                        </h3>
                        <div className='flex items-center gap-2 mt-1'>
                            <span className='text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 flex items-center gap-1'>
                                <Package size={10} />
                                {weight}
                            </span>
                        </div>
                    </div>

                    {/* Shop Info */}
                    <div className='flex items-center gap-1.5 mb-2'>
                        <span className={`text-sm ${isShopClosed ? 'text-gray-400' : 'text-gray-600'}`}>
                            {product_name}
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
                        <div className='text-sm text-red-500'>
                            Not Available
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

export default SearchProductCard

export { SearchProductSkeleton, EmptyState }
