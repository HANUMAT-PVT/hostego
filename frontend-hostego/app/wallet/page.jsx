'use client'
import React, { useEffect, useState } from 'react'
import { Upload } from "lucide-react"
import BackNavigationButton from '../components/BackNavigationButton'
import { uploadToS3Bucket } from '../lib/aws'
import HostegoButton from "../components/HostegoButton"

import axiosClient from '../utils/axiosClient'


const Page = () => {
    const defaultWalletDetails = {
        amount: 0,
        unique_transaction_id: "",
        payment_screenshot_img_url: ""
    }
    const [walletDetails, setWalletDetails] = useState(defaultWalletDetails)
    const [userWallet, setUserWallet] = useState({ amount: 0 });
    const [walletTransactionCreationLoading, setWalletTransactionCreationLoading] = useState(false)
    const [paymentScreenShotImgUrl, setPaymentScreenShotImgUrl] = useState(null);

    useEffect(() => {
        fetchUserWallet()
    }, [])

    const fetchUserWallet = async () => {
        try {
            let { data } = await axiosClient.get("/api/wallet")
            setUserWallet(wallet)
        } catch (error) {

        }
    }
    const handleFileChange = (e) => {
        const file = e.target.files[0];
        setPaymentScreenShotImgUrl(file)
    };


    const handleWalletTransactionSubmit = async (e) => {
        try {
            setWalletTransactionCreationLoading(true)
            e.preventDefault()
            if (walletDetails?.amount < 100) {
                alert("Minimum amount to add is 100");
                return;
            }
            if (!paymentScreenShotImgUrl) return
            const imageUrl = await uploadToS3Bucket(paymentScreenShotImgUrl);

            await axiosClient.post('/api/wallet/credit', {
                ...walletDetails,
                payment_screenshot_img_url: imageUrl,
            })
            setPaymentScreenShotImgUrl(null)
            setWalletDetails(defaultWalletDetails)
            alert("Money add request added successfully")
        } catch (error) {
            console.log(error)
            alert("error")
        }
        finally {
            setWalletTransactionCreationLoading(false)
        }
    }

    return (
        <div className='bg-[var(--bg-page-color)]'>
            <BackNavigationButton title={"Wallet"} />
            <div className='p-4 flex flex-col gap-4'>
                <div className='bg-white rounded-md w-full p-2 '>
                    <p className='mb-2 text-sm'>Balance</p>
                    <p className='font-normal text-2xl'>₹ {userWallet?.amount}</p>
                </div>
                {/* QR Code Section */}
                <div className="bg-white p-4 rounded-md flex flex-col items-center">
                    <p className="text-black font-normal text-lg mb-2">Scan QR to Pay</p>
                    <img
                        src="/hostego_payment_qr.jpg"
                        alt="Payment QR Code"
                        className="w-[300px] h-[300px] object-cover"
                    />
                    <p className="text-sm text-gray-500 mt-2">Scan the QR code to make a payment</p>
                </div>

                <div className='bg-white p-2 rounded-md'>
                    <div className="flex flex-col gap-6">
                        <p className="text-black font-normal text-lg">Add Money</p>
                        <form className='flex flex-col gap-6' onSubmit={handleWalletTransactionSubmit} >
                            {/* Amount Input */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Amount
                                </label>
                                <div className="flex font-normal items-center border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <span className="text-gray-700 ml-2">₹</span>
                                    <input
                                        type="number"
                                        placeholder="500"
                                        value={walletDetails?.amount}
                                        onChange={(e) => {
                                            const value = e.target.value ? Number(e.target.value) : "";
                                            setWalletDetails({ ...walletDetails, amount: value });
                                        }}
                                        className="ml-2 outline-none bg-transparent cursor-pointer w-full"
                                    />
                                </div>
                            </div>

                            {/* Unique Transaction ID Input */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Unique Transaction Id
                                </label>
                                <div className="flex font-normal items-center border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <input
                                        type="text"
                                        placeholder="Enter UTR e.g. #12121Ddfs"
                                        value={walletDetails?.unique_transaction_id}
                                        onChange={(e) => setWalletDetails({ ...walletDetails, unique_transaction_id: e.target.value })}
                                        className="ml-2 outline-none bg-transparent cursor-pointer w-full"
                                    />
                                </div>
                            </div>

                            {/* Upload Screenshot Button */}
                            <div className="relative">
                                <label className="absolute text-[#655df0] font-normal -top-3 left-3 bg-white px-1 text-sm">
                                    Upload Payment Screenshot
                                </label>
                                <div className="flex items-center justify-between border-2 border-[#655df0] max-w-[400px] rounded-md px-4 py-3 w-full">
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={handleFileChange}
                                        className="hidden"
                                        id="screenshot-upload"
                                        placeholder='Payment screenshot'
                                    />
                                    <label
                                        htmlFor="screenshot-upload"
                                        className="flex items-center justify-between w-full cursor-pointer"
                                    >
                                        <span className="text-gray-500">
                                            {paymentScreenShotImgUrl ? paymentScreenShotImgUrl.name : "Payment screenshot "}
                                        </span>
                                        <Upload className="text-[#655df0]" />
                                    </label>
                                </div>
                            </div>

                            {/* Submit Button */}
                            <HostegoButton isLoading={walletTransactionCreationLoading} onClick={handleWalletTransactionSubmit} text={`Add ₹${walletDetails?.amount}`} />
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Page;
