"use client"
import React, { useState, useEffect, useCallback } from 'react'
import SearchComponent from '../components/SearchComponent'
import CartFloatingButton from '../components/Cart/CartFloatingButton'
import ProductCard from '../components/ProductCard'
import { debounce } from 'lodash'
import axiosClient from '../utils/axiosClient'
import { useSelector } from 'react-redux'
import SearchProductCard, { SearchProductSkeleton, EmptyState } from '../components/SearchProductCard'

const SearchSkeletons = () => {
    return (
        <>
            {[1, 2, 3, 4, 5].map((item) => (
                <SearchProductSkeleton key={item} />
            ))}
        </>
    )
}

const page = () => {
    const [searchValue, setSearchValue] = useState("");
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [currentPage, setCurrentPage] = useState(1);

    const { cartData } = useSelector(state => state.user);

    // Initial fetch
    useEffect(() => {
        fetchProducts(searchValue);
    }, [currentPage]); // Only re-fetch when page changes


    const fetchProducts = useCallback(async (search) => {
        try {
            setLoading(true);
            const { data } = await axiosClient.get(`/api/products/all?page=${currentPage}&limit=50&search=${search}&admin=false`);
            setProducts(data?.products);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }, [currentPage]);

    // Create a debounced version of the search
    const debouncedSearch = useCallback(
        debounce((searchTerm) => {
            fetchProducts(searchTerm);
        }, 500),
        [fetchProducts]
    );

    // Handle search value changes
    const handleSearch = (value) => {
        setSearchValue(value);
        debouncedSearch(value);
    };

    return (
        <div>
            <div className='gradient-background py-4 sticky top-0 z-10'>
                <SearchComponent
                    viewOnly={false}
                    sendSearchValue={handleSearch}
                />
                {cartData?.cart_items?.length > 0 && <CartFloatingButton />}
            </div>

            <div className='p-4'>
                {loading ? (
                    <div className='flex flex-col gap-4'>
                        <SearchSkeletons />
                    </div>
                ) : error ? (
                    <div className="text-center text-red-500 py-4">
                        <p className="font-medium">Oops! Something went wrong</p>
                        <p className="text-sm mt-1">{error}</p>
                    </div>
                ) : (
                    <div className='flex flex-col gap-4'>
                        {products?.length > 0 ? (
                            products.map((prd) => (
                                <SearchProductCard
                                    key={prd?.product_id}
                                    {...prd}
                                    isAlreadyInCart={cartData?.cart_items?.some(
                                        item => item?.product_id === prd?.product_id
                                    )}
                                />
                            ))
                        ) : (
                            <EmptyState sendSearchValue={handleSearch} searchValue={searchValue} />
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default page;
