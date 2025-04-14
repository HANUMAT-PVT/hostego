"use client"
import React from 'react'

import axiosClient from '../utils/axiosClient'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { setUserAccount, setCartData, setFetchCartData, setUserAddresses, setUserRoles, setUserAccountWallet, setFetchUserAccount } from '../lib/redux/features/user/userSlice'
import { useDispatch, useSelector } from 'react-redux'
import { subscribeToNotifications } from '../utils/webNotifications'
import  SSEHandler  from '../components/SSEHandler'

const ClientComponent = ({ children }) => {

    const dispatch = useDispatch();
    const { fetchCartData, fetchUser,userAccount } = useSelector((state) => state.user)
    const router = useRouter()
   
    useEffect(() => {

        if (fetchUser) {
            fetchUserAccount();
        }


    }, [fetchUser]);

    useEffect(() => {
        if ("serviceWorker" in navigator) {
            navigator.serviceWorker
                .register("/sw.js")
                .catch((error) => {
                    console.error("Service Worker Registration Failed", error);
                });
        }
    }, []);

    useEffect(() => {
        fetchUserWallet()
    }, [])

    useEffect(() => {
        if (fetchCartData) {
            fetchCartItems();
        }
    }, [fetchCartData]);

    useEffect(() => {
        fetchUserRoles()
        fetchAddress()
    }, [])

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
            const { data } = await axiosClient.get("/api/users/me");

            dispatch(setUserAccount(data))
            fetchCartItems()
            dispatch(setFetchUserAccount(false))
        } catch (error) {
            console.log(error);
            router.push('/auth/sign-up')
        }
    };

    const fetchUserRoles = async () => {
        try {
            const { data } = await axiosClient.get('/api/user-roles')
            dispatch(setUserRoles(data))
        } catch (error) {
            console.error('Error fetching user roles:', error)
        }
    }

    const fetchUserWallet = async () => {
        try {
            const { data } = await axiosClient.get('/api/wallet')
            dispatch(setUserAccountWallet(data))
        } catch (error) {
            console.error('Error fetching user wallet:', error)
        }
    }


    const fetchAddress = async () => {
        try {
            const { data } = await axiosClient.get('/api/address')
            dispatch(setUserAddresses(data))
        } catch (error) {
            console.error('Error fetching address:', error)
        }
    }

    console.log(
        "%cSTOP!",
        "color: red; font-size: 50px; font-weight: bold; text-shadow: 2px 2px black;"
    );
    console.log(
        "%cThis is a browser feature intended for developers. If someone told you to copy-paste something here to enable an Hostego feature or \"hack\" someone's account, it is a scam and will give them access to your Hostego account",
        "color: white; font-size: 16px;"
    );

    return (

        <>
            <SSEHandler userId={userAccount?.user_id ||""} />
            {children}
        </>

    )
}

export default ClientComponent
