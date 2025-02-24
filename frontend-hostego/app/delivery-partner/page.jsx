"use client";

import React, { useState } from "react";
import BackNavigationButton from "../components/BackNavigationButton";
import { Info, Landmark, ShoppingBag, Clock, Upload } from "lucide-react";
import HostegoButton from "../components/HostegoButton"
import { uploadToS3Bucket } from '../lib/aws'
import axiosClient from "../utils/axiosClient"
import MaintainOrderStatusForDeliveryPartner from "../components/Delivery-Partner/MaintainOrderStatusForDeliveryPartner"

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

const page = () => {

    const defaultDeliveryParnter = {
        address: "",
        aadhaar_front_img: "",
        aadhaar_back_img: "",
        upi_id: "",
        bank_details_img: "",
        address: ""
    }


    const [isOnline, setIsOnline] = useState(false);
    const [deliveryParnterDetails, setDeliveryParnterDetails] = useState(defaultDeliveryParnter)

    const totalEarnings = ordersData.reduce(
        (sum, day) => sum + day?.orders?.reduce((daySum, order) => daySum + order?.earning, 0),
        0
    );



    const handleDeliveryPartnerRegistration = async (e) => {
        e.preventDefault()
        try {

            const [a_front_img_url, a_back_img_url, b_details_img] = await Promise.all(
                await uploadToS3Bucket(deliveryParnterDetails?.aadhaar_front_img),
                await uploadToS3Bucket(deliveryParnterDetails?.aadhaar_back_img),
                await uploadToS3Bucket(deliveryParnterDetails?.bank_details_img)
            )
            console.log(a_front_img_url)
            return
            let { data } = await axiosClient.post("/users",
                {
                    ...deliveryParnterDetails,
                    aadhaar_front_img: a_front_img_url,
                    aadhaar_back_img: a_back_img_url,
                    bank_details_img: b_details_img
                });
            console.log(data, "data from the api response")

        } catch (error) {
            console.log(error)
        }
    }

    return (
        <div className="bg-[var(--bg-page-color)] ">
            <div className="sticky top-0 z-30">
                <BackNavigationButton title="Delivery Partner" />
                {/* Online/Offline Toggle */}
                <div className="bg-white  px-4 py-4 rounded-md flex justify-between items-center shadow-md ">
                    <div
                        className={`relative w-24 h-8 rounded-full cursor-pointer transition flex items-center ${isOnline ? "bg-green-500" : "bg-gray-400"}`}
                        onClick={() => setIsOnline(!isOnline)}
                    >
                        <span
                            className={`absolute w-full text-center text-xs ${isOnline ? "-ml-3" : "ml-3"} font-bold text-white`}
                        >
                            {isOnline ? "ONLINE" : "OFFLINE"}
                        </span>
                        <div
                            className={`absolute top-1 left-1 w-6 h-6 bg-white rounded-full shadow-md transition-transform ${isOnline ? "translate-x-[60px]" : "translate-x-0"}`}
                        />
                    </div>
                    <Info color="var(--primary-color)" size={24} />
                </div>
            </div>

            {/* Delivery Partner Form */}
            <div className='p-4 flex flex-col gap-4'>
                <div className='bg-white p-2 rounded-md'>
                    <div className="flex flex-col gap-6">
                        <p className="text-black font-normal text-lg">Personal Information</p>
                        <form className='flex flex-col gap-6' onSubmit={handleDeliveryPartnerRegistration}  >
                            {/* Amount Input */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Address
                                </label>
                                <div className="flex font-normal items-center border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">

                                    <input
                                        type="text"
                                        placeholder="Hostel-room-no"
                                        value={deliveryParnterDetails?.address}
                                        onChange={(e) => setDeliveryParnterDetails({ ...deliveryParnterDetails, address: e.target.value })}
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
                                        value={deliveryParnterDetails?.upi_id}
                                        onChange={(e) => setDeliveryParnterDetails({ ...deliveryParnterDetails, upi_id: e.target.value })}
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
                                        onChange={(e) => setDeliveryParnterDetails({ ...deliveryParnterDetails, aadhaar_front_img: e.target.files[0] })}
                                        className="hidden"
                                        id="screenshot-upload-front"
                                        placeholder='Payment screenshot'
                                    />
                                    <label
                                        htmlFor="screenshot-upload-front"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500">
                                            {deliveryParnterDetails?.aadhaar_front_img ? deliveryParnterDetails?.aadhaar_front_img.name : "Aadhar front "}
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
                                        onChange={(e) => setDeliveryParnterDetails({ ...deliveryParnterDetails, aadhaar_back_img: e.target.files[0] })}
                                        className="hidden"
                                        id="screenshot-upload-back"
                                        placeholder='Aadhar back '
                                    />
                                    <label
                                        htmlFor="screenshot-upload-back"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500">
                                            {deliveryParnterDetails?.aadhaar_back_img ? deliveryParnterDetails?.aadhaar_back_img.name : "Aadhar Back "}
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
                                        onChange={(e) => setDeliveryParnterDetails({ ...deliveryParnterDetails, bank_details_img: e.target.files[0] })}
                                        className="hidden"
                                        id="screenshot-upload-bank"
                                        placeholder='Payment screenshot'
                                    />
                                    <label
                                        htmlFor="screenshot-upload-bank"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500">
                                            {deliveryParnterDetails?.bank_details_img ? deliveryParnterDetails?.bank_details_img.name : "Bank Details "}
                                        </span>
                                        <Upload className="text-[#655df0]" />
                                    </label>
                                </div>
                            </div>

                            {/* Submit Button */}
                            <HostegoButton type="submit" text={"Submit"} />
                        </form>
                    </div>
                </div>
            </div>



            {/* My Progress Section */}
            <div className="bg-white p-2 m-4 rounded-md mt-4">
                <p className="text-gray-600 font-normal text-md border-b pb-2">MY PROGRESS</p>
                <div className="flex justify-between mt-3">
                    <div className="flex flex-col items-center gap-1 px-4">
                        <p className="font-semibold text-xl">₹ {totalEarnings}</p>
                        <div className="flex gap-2 items-center">
                            <Landmark size={14} />
                            <p className="text-xs">Total earnings</p>
                        </div>
                    </div>
                    <div className="flex flex-col items-center gap-1 px-4">
                        <p className="font-semibold text-xl">5</p>
                        <div className="flex gap-2 items-center">
                            <ShoppingBag size={14} />
                            <p className="text-xs">Orders</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Order History Section */}
            <div className="mt-6 px-4 ">
                {ordersData.map((day, index) => (
                    <div key={index} className="mb-6 mt-6">
                        <p className="text-lg font-semibold mb-2">{day?.date}</p>
                        <div className="bg-white p-4 rounded-md shadow-md">
                            {day.orders.map((order, idx) => (
                                <div
                                    key={idx}
                                    className="flex justify-between items-center border-b last:border-none py-2"
                                >
                                    <div className="flex flex-col gap-2">
                                        <p className="text-md font-normal">Order ID: {order?.id}</p>
                                        <p className="text-sm text-gray-500 flex items-center gap-1">
                                            <Clock size={14} />
                                            {order.time}
                                        </p>
                                    </div>
                                    <p className="font-semibold text-green-600">₹ {order.earning}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>

            <MaintainOrderStatusForDeliveryPartner />
        </div>
    );
};

export default page;