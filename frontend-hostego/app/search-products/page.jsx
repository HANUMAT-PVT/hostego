"use client"
import React, { useState, useEffect, useCallback } from 'react'
import SearchComponent from '../components/SearchComponent'
import CartFloatingButton from '../components/Cart/CartFloatingButton'
import ProductCard from '../components/ProductCard'
import { debounce } from 'lodash'
import axiosClient from '../utils/axiosClient'
import { useSelector } from 'react-redux'

const page = () => {
    const [searchValue, setSearchValue] = useState("");
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalItems, setTotalItems] = useState(0);

    const { cartData } = useSelector(state => state.user);


    // Initial fetch
    useEffect(() => {
        fetchProducts(searchValue);
    }, [currentPage]); // Only re-fetch when page changes

    // Create a memoized fetch function
    // what will be use of memoization here
    // memoization is used to prevent the re-rendering of the component when the same function is called again and again
    // it will only call the function when the search value changes
    const fetchProducts = useCallback(async (search) => {
        try {
            setLoading(true);
            const { data } = await axiosClient.get(`/api/products/all?page=${currentPage}&limit=10&search=${search}`);
            setProducts(data);
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
                {/* <CartFloatingButton /> */}
            </div>

            <div className='p-4'>
                {loading ? (
                    <div className="text-center py-4">Loading...</div>
                ) : error ? (
                    <div className="text-center text-red-500 py-4">{error}</div>
                ) : (
                    <div className='grid grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-4'>
                        {products?.map((prd) => (
                            <ProductCard
                                isAlreadyInCart={cartData?.cart_items?.some(
                                    item => item?.product_id === prd?.product_id
                                )}
                                {...prd}
                                key={prd?.product_id}
                            />
                        ))}
                        {products?.length === 0 && (
                            <div className="col-span-full text-center py-8 text-gray-500">
                                No products found
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default page;
