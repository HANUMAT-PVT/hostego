"use client"
import React from 'react'

import axiosClient from '../utils/axiosClient'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { setUserAccount, setCartData, setFetchCartData, setUserAddresses } from '../lib/redux/features/user/userSlice'
import { useDispatch, useSelector } from 'react-redux'


const ClientComponent = ({ children }) => {

    const dispatch = useDispatch();
    const { fetchCartData } = useSelector((state) => state.user)
    const router = useRouter()

    useEffect(() => {
        fetchUserAccount();
    }, []);

    useEffect(() => {
        if (fetchCartData) {
            fetchCartItems();
        }
    }, [fetchCartData]);

    const fetchCartItems = async () => {
        try {

            const { data } = await axiosClient.get('/api/cart/')
            dispatch(setCartData(data))
            dispatch(setFetchCartData(false))
        } catch (error) {
            console.error('Error fetching cart:', error)
        } finally {

        }
    }

    const fetchUserAccount = async () => {
        try {
            const { data } = await axiosClient.get("/api/user/me");

            dispatch(setUserAccount(data))
            fetchCartItems()
        } catch (error) {
            console.log(error);
            router.push('/auth/sign-up')
        }
    };

    useEffect(() => {
        fetchAddress()
    }, [])

    const fetchAddress = async () => {
        try {
            const { data } = await axiosClient.get('/api/address')
            dispatch(setUserAddresses(data))
        } catch (error) {
            console.error('Error fetching address:', error)
        }
    }

    return (

        <>
            {children}
        </>

    )
}

export default ClientComponent
