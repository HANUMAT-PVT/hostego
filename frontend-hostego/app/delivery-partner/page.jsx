"use client";

import React, { useEffect, useState } from "react";
import BackNavigationButton from "../components/BackNavigationButton";
import { Info, Landmark, ShoppingBag, Clock, Upload, User, ChevronDown, Shield, RefreshCw } from "lucide-react";
import HostegoButton from "../components/HostegoButton"
import { uploadToS3Bucket } from '../lib/aws'
import axiosClient from "../utils/axiosClient"
import MaintainOrderStatusForDeliveryPartner from "../components/Delivery-Partner/MaintainOrderStatusForDeliveryPartner"
import { formatDate, transformDeliveryPartnerOrderEarnings, transformOrdersByDate } from "../utils/helper";

const ordersData = [
    {
        date: "22 Feb 2025",
        orders: [
            { id: "ORD1234", earning: 15, time: "10:30 AM" },
            { id: "ORD1235", earning: 23, time: "11:15 AM" },
            { id: "ORD1236", earning: 19, time: "01:45 PM" },
        ],
    },
    {
        date: "21 Feb 2025",
        orders: [
            { id: "ORD1229", earning: 21, time: "09:00 AM" },
            { id: "ORD1230", earning: 29, time: "12:30 PM" },
        ],
    },
];

const PersonalInfoAccordion = ({ deliveryPartner }) => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="w-full px-4 py-3 flex items-center justify-between"
            >
                <div className="flex items-center gap-3">
                    <User className="w-5 h-5 text-[var(--primary-color)]" />
                    <span className="font-medium">Personal Information</span>
                </div>
                <ChevronDown className={`w-5 h-5 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
            </button>

            {isOpen && (
                <div className="px-4 pb-4 space-y-4 border-t animate-fade-in">
                    <div className="flex items-center justify-between">
                        <span className="text-gray-500">Name</span>
                        <span className="font-medium">{deliveryPartner.user?.first_name}</span>
                    </div>
                    <div className="flex items-center justify-between">
                        <span className="text-gray-500">Phone</span>
                        <span className="font-medium">{deliveryPartner.user?.mobile_number}</span>
                    </div>
                    <div className="flex items-center justify-between">
                        <span className="text-gray-500">Email</span>
                        <span className="font-medium">{deliveryPartner.user?.email}</span>
                    </div>
                    <div className="flex items-center justify-between">
                        <span className="text-gray-500">Address</span>
                        <span className="font-medium">{deliveryPartner.address}</span>
                    </div>
                </div>
            )}
        </div>
    );
};

const VerificationStatus = ({ deliveryPartner }) => {
    const status = deliveryPartner?.verification_status
    return (
        <div className="bg-white rounded-xl p-6 shadow-sm">
            <div className="flex items-center gap-3 mb-4">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${status ? 'bg-green-50' : 'bg-yellow-50'
                    }`}>
                    <Shield className={`w-6 h-6 ${status ? 'text-green-500' : 'text-yellow-500'
                        }`} />
                </div>
                <div>
                    <h3 className="font-medium">Verification Status</h3>
                    <p className={`text-sm ${status ? 'text-green-500' : 'text-yellow-500'
                        }`}>
                        {status ? 'Verified Partner' : 'Verification Pending'}
                    </p>
                </div>
            </div>

            {!status && (
                <div className="bg-yellow-50 border border-yellow-100 rounded-lg p-3">
                    <p className="text-sm text-yellow-700">
                        {deliveryPartner?.documents?.upi_id
                            ? "We have received your details. Sit back and relax while we verify your documents."
                            : "Please complete your verification to start accepting orders. Upload the required documents below."}


                    </p>
                </div>
            )}
        </div>
    );
};



const Page = () => {
    const [isLoading, setIsLoading] = useState(true);
    const [isOnline, setIsOnline] = useState(false);
    const [deliveryPartnerOrders, setDeliveryPartnerOrders] = useState([]);
    const [deliveryPartner, setDeliveryPartner] = useState({});
    const [formSubmitingLoading, setFormSubmitingLoading] = useState(false);
    const [deliveryPartnerEarnings, setDeliveryPartnerEarnings] = useState(0);
    const [isRefreshing, setIsRefreshing] = useState(false);

    const [deliveryPartnerVerificationData, setDeliveryPartnerVerificationData] = useState({
        address: "",
        aadhaar_front_img: "",
        aadhaar_back_img: "",
        upi_id: "",
        bank_details_img: "",
    });

    useEffect(() => {
        fetchDeliveryPartner();
    }, []);

    useEffect(() => {
        if (deliveryPartner?.delivery_partner_id) {
            fetchDeliveryPartnerOrders();
            fetchDeliveryPartnerEarnings();
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
            let { data } = await axiosClient.get(`/api/order/delivery-partner/${deliveryPartner?.delivery_partner_id}`)
            setDeliveryPartnerOrders(data?.orders)
        } catch (error) {
            console.log(error)
        }
    }

    const handleRefresh = async () => {
        try {
            setIsRefreshing(true);
            await Promise.all([
                fetchDeliveryPartnerOrders(),
                fetchDeliveryPartnerEarnings()
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
                {/* Verification Status */}
                <VerificationStatus deliveryPartner={deliveryPartner} />

                {/* Personal Information Accordion */}
                <PersonalInfoAccordion deliveryPartner={deliveryPartner} />

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
                            <p className="font-semibold text-xl">₹ {deliveryPartnerEarnings?.summary?.total_earnings}</p>
                            <div className="flex gap-2 items-center">
                                <Landmark size={14} />
                                <p className="text-xs">Total earnings</p>
                            </div>
                        </div>
                        <div className="flex flex-col items-center gap-1 px-4">
                            <p className="font-semibold text-xl">{deliveryPartnerEarnings?.summary?.total_orders}</p>
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

                {/* Order History Section */}
                {/* <div className="mt-6 px-4 ">
                    {deliveryPartnerEarnings?.earnings?.map((day, index) => (
                        <div key={index} className="mb-6 mt-6">
                            <p className="text-lg font-semibold mb-2">{day?.date}</p>
                            <div className="bg-white p-4 rounded-md shadow-md flex flex-col gap-2">
                                {day?.orders?.map((order, idx) => (
                                    <div
                                        key={idx}
                                        className="flex justify-between items-center border-b last:border-none py-2"
                                    >
                                        <div className="flex flex-col gap-2">
                                            <div className="flex items-center gap-2">
                                                <span className="bg-blue-100 text-blue-700 px-3 py-1 rounded-full text-sm font-medium">
                                                    Order #{order?.order_id?.slice(-6)}
                                                </span>

                                            </div>
                                            <p className="text-sm text-gray-500 flex items-center gap-1">
                                                <Clock size={14} />
                                                {formatDate(order?.created_at)}
                                            </p>
                                        </div>
                                        <p className="font-semibold text-green-600">₹ {(order?.delivery_partner_fee).toFixed(1)}</p>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div> */}
            </div>

            {deliveryPartnerOrders?.map((order) => (
                <MaintainOrderStatusForDeliveryPartner onUpdateOrderStatus={updateOrderStatus} key={order?.order_id} order={order} />
            ))}
        </div>
    );
};

export default Page;