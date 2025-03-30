"use client";

import React, { useEffect, useState } from "react";
import BackNavigationButton from "../components/BackNavigationButton";
import { Info, Landmark, ShoppingBag, Clock, Upload, User, ChevronDown, Shield, RefreshCw, Package, Wallet, ChevronUp } from "lucide-react";
import HostegoButton from "../components/HostegoButton"
import { uploadToS3Bucket } from '../lib/aws'
import axiosClient from "../utils/axiosClient"
import MaintainOrderStatusForDeliveryPartner from "../components/Delivery-Partner/MaintainOrderStatusForDeliveryPartner"
import { transformOrdersByDate, formatDate } from "../utils/helper";
import TransactionCard from '../components/TransactionCard';
import PersonalInfoAccordion from '../components/Delivery-Partner/PersonalInfoAccordion';
import VerificationStatus from '../components/Delivery-Partner/VerificationStatus';
import WalletCard from '../components/Delivery-Partner/WalletCard';

const Page = () => {
    const [isLoading, setIsLoading] = useState(true);
    const [isOnline, setIsOnline] = useState(false);
    const [deliveryPartnerOrders, setDeliveryPartnerOrders] = useState([]);
    const [deliveryPartner, setDeliveryPartner] = useState({});
    const [formSubmitingLoading, setFormSubmitingLoading] = useState(false);
    const [deliveryPartnerEarnings, setDeliveryPartnerEarnings] = useState(0);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [selectedFilter, setSelectedFilter] = useState("active")
    const [walletData, setWalletData] = useState(null);

    const [deliveryPartnerVerificationData, setDeliveryPartnerVerificationData] = useState({
        address: "",
        aadhaar_front_img: "",
        aadhaar_back_img: "",
        upi_id: "",
        bank_details_img: "",
    });

    const filterOptions = [
        { value: 'active', label: 'Active Orders' },
        { value: 'delivered', label: 'Delivered' },
        { value: '', label: 'All Orders' }
    ];

    useEffect(() => {
        fetchDeliveryPartner(selectedFilter);
    }, []);

    useEffect(() => {
        if (deliveryPartner?.delivery_partner_id) {
            fetchDeliveryPartnerOrders();
            fetchDeliveryPartnerEarnings();
            fetchWalletData();
        }
    }, [deliveryPartner]);

    const fetchDeliveryPartner = async () => {
        try {
            setIsLoading(true);
            let { data } = await axiosClient.get("/api/delivery-partner/find");
            setDeliveryPartner(data);
            setIsOnline(!!data?.availability_status)
        } catch (error) {
            console.log(error);
        } finally {
            setIsLoading(false);
        }
    };

    const updateDeliveryPartnerAvailabilityStatus = async (newStatus) => {
        try {
            let { data } = await axiosClient.patch(`/api/delivery-partner/${deliveryPartner?.delivery_partner_id}`, {
                availability_status: newStatus ? 1 : 0
            });

            setIsOnline(newStatus)
        } catch (error) {
            console.log(error)
        }
    }

    const updateOrderStatus = async (orderId, newStatus) => {
        try {
            let { data } = await axiosClient.patch(`/api/order/${orderId}`, {
                order_status: newStatus
            })
            fetchDeliveryPartnerOrders()
        } catch (error) {
            console.log(error)
        }
    }

    const fetchDeliveryPartnerEarnings = async () => {
        try {
            let { data } = await axiosClient.get(`/api/delivery-partner/earnings/${deliveryPartner?.delivery_partner_id}`)
            const earnings = transformOrdersByDate(data?.daily_earnings)

            setDeliveryPartnerEarnings({ ...data, earnings })
        } catch (error) {
            console.log(error)
        }
    }

    const fetchWalletData = async () => {
        try {
            const { data } = await axiosClient.get(`/api/delivery-partner-wallet/${deliveryPartner?.delivery_partner_id}`);
            const transactions = await axiosClient.get(`/api/delivery-partner-wallet/transactions/${deliveryPartner?.delivery_partner_id}?page=1&limit=10`)
            setWalletData({ ...data, recent_transactions: transactions.data });
        } catch (error) {
            console.error('Error fetching wallet data:', error);
        }
    };

    const handleDeliveryPartnerRegistration = async (e) => {
        e.preventDefault()
        setFormSubmitingLoading(true)
        const aadhaarImg = deliveryPartnerVerificationData?.aadhaar_front_img;
        const aadhaarBackImg = deliveryPartnerVerificationData?.aadhaar_back_img
        const bankImg = deliveryPartnerVerificationData?.bank_details_img
        try {
            if (!aadhaarImg || !aadhaarBackImg || !bankImg) {
                alert("Please upload all the images")
                return
            }
            const [a_front_img_url, a_back_img_url, b_details_img] = await Promise.all(
                [await uploadToS3Bucket(aadhaarImg),
                await uploadToS3Bucket(aadhaarBackImg),
                await uploadToS3Bucket(bankImg)]
            )
            const requestBody = {
                ...deliveryPartnerVerificationData,
                aadhaar_front_img: a_front_img_url,
                aadhaar_back_img: a_back_img_url,
                bank_details_img: b_details_img
            }

            let { data } = await axiosClient.post("/api/delivery-partner", {
                address: deliveryPartnerVerificationData.address,
                documents: requestBody
            });
            fetchDeliveryPartner(data)

        } catch (error) {
            console.log(error)
        } finally {
            setFormSubmitingLoading(false)
        }
    }

    const fetchDeliveryPartnerOrders = async () => {
        try {
            if (!deliveryPartner.delivery_partner_id) return
            let { data } = await axiosClient.get(`/api/order/delivery-partner/${deliveryPartner?.delivery_partner_id}?status=${selectedFilter}`);
            setDeliveryPartnerOrders(data?.orders);
        } catch (error) {
            console.error('Error fetching orders:', error);
        } finally {

        }
    };

    useEffect(() => {
        fetchDeliveryPartnerOrders();
    }, [selectedFilter]);

    const handleFilterChange = (value) => {
        setSelectedFilter(value);
    };

    const handleRefresh = async () => {
        try {
            setIsRefreshing(true);
            await Promise.all([
                fetchDeliveryPartnerOrders(),
                fetchDeliveryPartnerEarnings(),
                fetchWalletData()
            ]);
        } catch (error) {
            console.error('Error refreshing data:', error);
        } finally {
            setIsRefreshing(false);
        }
    };

    if (isLoading) {
        return (
            <div className="bg-[var(--bg-page-color)]">
                <BackNavigationButton title="Delivery Partner" />

                <div className="p-2 space-y-4">
                    {/* Verification Status Skeleton */}
                    <div className="bg-white p-4 rounded-md ">
                        <div className="h-8 bg-gray-200 rounded w-3/4 mb-3" />
                        <div className="h-6 bg-gray-200 rounded w-1/2" />
                    </div>

                    {/* Progress Skeleton */}
                    <div className="bg-white p-4 rounded-md ">
                        <div className="h-6 bg-gray-200 rounded w-1/3 mb-4" />
                        <div className="flex justify-between">
                            <div className="space-y-2">
                                <div className="h-8 bg-gray-200 rounded w-20" />
                                <div className="h-4 bg-gray-200 rounded w-24" />
                            </div>
                            <div className="space-y-2">
                                <div className="h-8 bg-gray-200 rounded w-20" />
                                <div className="h-4 bg-gray-200 rounded w-24" />
                            </div>
                        </div>
                    </div>

                    {/* Orders Skeleton */}
                    <div className="bg-white p-4 rounded-md ">
                        <div className="space-y-4">
                            {[1, 2, 3].map((item) => (
                                <div key={item} className="flex justify-between items-center">
                                    <div className="space-y-2">
                                        <div className="h-5 bg-gray-200 rounded w-32" />
                                        <div className="h-4 bg-gray-200 rounded w-24" />
                                    </div>
                                    <div className="h-6 bg-gray-200 rounded w-16" />
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[var(--bg-page-color)]">
            <BackNavigationButton title="Delivery Partner" />

            <div className="p-4 space-y-4">


                {deliveryPartner?.verification_status == 1 && <div className="bg-white p-4 rounded-xl shadow-sm mb-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <div className={`w-2 h-2 rounded-full ${isOnline ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`} />
                            <div className="flex flex-col">
                                <span className="font-medium text-gray-900">
                                    {isOnline ? 'Online' : 'Offline'}
                                </span>
                                <span className="text-xs text-gray-500">
                                    {isOnline ? 'You\'re accepting delivery orders' : 'Go online to start accepting orders'}
                                </span>
                            </div>
                        </div>

                        <label className="relative inline-flex items-center cursor-pointer">
                            <input
                                type="checkbox"
                                className="sr-only peer"
                                checked={isOnline}
                                onChange={() => updateDeliveryPartnerAvailabilityStatus(!isOnline)}
                            />
                            <div className={`w-[68px] h-[29px] p-1 bg-gray-200 peer-focus:outline-none rounded-full peer 
                                peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full 
                                peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] 
                                ${isOnline ? 'after:start-[16px]' : 'after:start-[2px]'}  after:bg-white after:border-gray-300 after:border after:rounded-full 
                                after:h-6 after:w-6 after:transition-all peer-checked:bg-[var(--primary-color)]`}>
                            </div>
                        </label>
                    </div>
                </div>}
                <VerificationStatus deliveryPartner={deliveryPartner} />
                <PersonalInfoAccordion deliveryPartner={deliveryPartner} />
                <WalletCard walletData={walletData} />

                {/* My Progress Section with Refresh Button */}
                <div className="bg-white p-2 m-4 rounded-md mt-4">
                    <div className="flex items-center justify-between border-b pb-2">
                        <p className="text-gray-600 font-normal text-md">MY PROGRESS</p>
                        <button
                            onClick={handleRefresh}
                            disabled={isRefreshing}
                            className={`p-2 rounded-full ${isRefreshing ? 'bg-gray-100' : 'hover:bg-gray-100'} 
                                transition-all duration-200 active:scale-95`}
                        >
                            <RefreshCw
                                className={`w-5 h-5 text-gray-600 ${isRefreshing ? 'animate-spin' : ''}`}
                            />
                        </button>
                    </div>
                    <div className="flex justify-between mt-3">
                        <div className="flex flex-col items-center gap-1 px-4">
                            <p className="font-semibold text-xl">₹ {deliveryPartnerEarnings?.summary?.total_earnings || 0}</p>
                            <div className="flex gap-2 items-center">
                                <Landmark size={14} />
                                <p className="text-xs">Total earnings</p>
                            </div>
                        </div>
                        <div className="flex flex-col items-center gap-1 px-4">
                            <p className="font-semibold text-xl">{deliveryPartnerEarnings?.summary?.total_orders || 0}</p>
                            <div className="flex gap-2 items-center">
                                <ShoppingBag size={14} />
                                <p className="text-xs">Orders</p>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Document Upload Form - show only if not verified */}
                {!deliveryPartner.documents?.upi_id && (
                    <div className="bg-white rounded-xl p-6 shadow-sm">
                        <h3 className="font-medium mb-4">Document Verification</h3>
                        <form className="space-y-4" onSubmit={handleDeliveryPartnerRegistration}>
                            {/* Amount Input */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Address
                                </label>
                                <div className="flex font-normal items-center border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">

                                    <input
                                        type="text"
                                        placeholder="Hostel-room-no"
                                        value={deliveryPartnerVerificationData?.address}
                                        onChange={(e) => setDeliveryPartnerVerificationData({ ...deliveryPartnerVerificationData, address: e.target.value })}
                                        className="ml-2 outline-none bg-transparent cursor-pointer w-full"
                                    />
                                </div>
                            </div>

                            {/* Unique Transaction ID Input */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    UPI ID
                                </label>
                                <div className="flex font-normal items-center border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <input
                                        type="text"
                                        placeholder="eg. xyz@ybl"
                                        value={deliveryPartnerVerificationData?.upi_id}
                                        onChange={(e) => setDeliveryPartnerVerificationData({ ...deliveryPartnerVerificationData, upi_id: e.target.value })}
                                        className="ml-2 outline-none bg-transparent cursor-pointer w-full"
                                    />
                                </div>
                            </div>

                            {/* Upload Adhaar Front Button */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Upload Aadhaar Front
                                </label>
                                <div className="flex items-center justify-between border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={(e) => setDeliveryPartnerVerificationData({ ...deliveryPartnerVerificationData, aadhaar_front_img: e.target.files[0] })}
                                        className="hidden "
                                        id="screenshot-upload-front"
                                        placeholder='Payment screenshot'
                                    />
                                    <label
                                        htmlFor="screenshot-upload-front"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500 truncate max-w-[200px] overflow-hidden text-ellipsis">
                                            {deliveryPartnerVerificationData?.aadhaar_front_img ? deliveryPartnerVerificationData?.aadhaar_front_img.name : "Aadhar front "}
                                        </span>
                                        <Upload className="text-[#655df0]" />
                                    </label>
                                </div>
                            </div>
                            {/* Upload Adhaar Front Button */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Upload Adhaar Back
                                </label>
                                <div className="flex items-center justify-between border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={(e) => setDeliveryPartnerVerificationData({ ...deliveryPartnerVerificationData, aadhaar_back_img: e.target.files[0] })}
                                        className="hidden"
                                        id="screenshot-upload-back"
                                        placeholder='Aadhar back '
                                    />
                                    <label
                                        htmlFor="screenshot-upload-back"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500 truncate max-w-[200px] overflow-hidden text-ellipsis">
                                            {deliveryPartnerVerificationData?.aadhaar_back_img ? deliveryPartnerVerificationData?.aadhaar_back_img.name : "Aadhar Back "}
                                        </span>
                                        <Upload className="text-[#655df0]" />
                                    </label>
                                </div>
                            </div>
                            {/* Upload Bank Details Button */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Upload Bank details
                                </label>
                                <div className="flex items-center justify-between border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={(e) => setDeliveryPartnerVerificationData({ ...deliveryPartnerVerificationData, bank_details_img: e.target.files[0] })}
                                        className="hidden"
                                        id="screenshot-upload-bank"
                                        placeholder='Payment screenshot'
                                    />
                                    <label
                                        htmlFor="screenshot-upload-bank"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500 truncate max-w-[200px] overflow-hidden text-ellipsis">
                                            {deliveryPartnerVerificationData?.bank_details_img ? deliveryPartnerVerificationData?.bank_details_img.name : "Bank Details "}
                                        </span>
                                        <Upload className="text-[#655df0]" />
                                    </label>
                                </div>
                            </div>

                            {/* Submit Button */}
                            <HostegoButton isLoading={formSubmitingLoading} type="submit" text={"Submit"} />
                        </form>
                    </div>
                )}
            </div>
            {/* Status Filter Section */}
            <div className="sticky top-0 z-10 bg-white border-b border-gray-100 p-4">
                <div className="flex items-center gap-2 sticky top-0 z-20 overflow-x-auto pb-2 scrollbar-hide">
                    {filterOptions.map((option) => (
                        <button
                            key={option.value}
                            onClick={() => handleFilterChange(option?.value)}
                            className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all
                            ${selectedFilter === option?.value
                                    ? 'bg-[var(--primary-color)] text-white'
                                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
                        >
                            {option?.label}
                        </button>
                    ))}
                </div>
            </div>
            <div className="flex flex-col gap-2 px-3">
                {deliveryPartnerOrders?.map((order) => (
                    <MaintainOrderStatusForDeliveryPartner onUpdateOrderStatus={updateOrderStatus} key={order?.order_id} order={order} />
                ))}
                {!isLoading && deliveryPartnerOrders.length === 0 && (
                    <div className="p-8 text-center">
                        <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                            <Package className="w-8 h-8 text-gray-400" />
                        </div>
                        <h3 className="text-lg font-medium text-gray-900 mb-1">No orders found</h3>
                        <p className="text-gray-500">
                            {selectedFilter === 'all'
                                ? "You don't have any orders yet"
                                : `No ${filterOptions.find(opt => opt.value === selectedFilter)?.label.toLowerCase()} orders`}
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Page;