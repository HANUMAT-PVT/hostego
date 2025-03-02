'use client'
import React, { useEffect, useState } from 'react'
import { Upload, Wallet, IndianRupee, History, AlertCircle, CheckCircle2 } from "lucide-react"
import BackNavigationButton from '../components/BackNavigationButton'
import { uploadToS3Bucket } from '../lib/aws'
import HostegoButton from "../components/HostegoButton"
import axiosClient from '../utils/axiosClient'
import HostegoLoader from '../components/HostegoLoader'

const Page = () => {
    const defaultWalletDetails = {
        amount: 0,
        unique_transaction_id: "",
        payment_screenshot_img_url: ""
    }
    const [walletDetails, setWalletDetails] = useState(defaultWalletDetails)
    const [userWallet, setUserWallet] = useState({ amount: 0 })
    const [walletTransactionCreationLoading, setWalletTransactionCreationLoading] = useState(false)
    const [paymentScreenShotImgUrl, setPaymentScreenShotImgUrl] = useState(null)
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        fetchUserWallet()
    }, [])

    const fetchUserWallet = async () => {
        try {
            setIsLoading(true)
            let { data } = await axiosClient.get("/api/wallet")
            setUserWallet(data)
        } catch (error) {
            console.error('Error fetching wallet:', error)
        } finally {
            setIsLoading(false)
        }
    }

    const handleFileChange = (e) => {
        const file = e.target.files[0]
        setPaymentScreenShotImgUrl(file)
    }

    const handleWalletTransactionSubmit = async (e) => {
        e.preventDefault()
        try {
            if (walletDetails?.amount < 100) {
                alert("Minimum amount to add is 100")
                return
            }
            if (!paymentScreenShotImgUrl) return

            setWalletTransactionCreationLoading(true)
            const imageUrl = await uploadToS3Bucket(paymentScreenShotImgUrl)

            await axiosClient.post('/api/wallet/credit', {
                ...walletDetails,
                payment_screenshot_img_url: imageUrl,
            })
            setPaymentScreenShotImgUrl(null)
            setWalletDetails(defaultWalletDetails)
            alert("Money add request added successfully")
            fetchUserWallet() // Refresh wallet balance
        } catch (error) {
            console.error(error)
            alert("Error adding money to wallet")
        } finally {
            setWalletTransactionCreationLoading(false)
        }
    }

    if (isLoading) {
        return <HostegoLoader />
    }
    const handlePayment = () => {
        const upiID = "8264121428@superyes"; // Replace with actual UPI ID
        const amount = "100"; // Replace with the desired amount
        const transactionNote = "Payment for services";
        
        const upiURL = `upi://pay?pa=${upiID}&pn=YourName&tn=${transactionNote}&am=${amount}&cu=INR`;
      
        window.location.href = upiURL; // Redirect to UPI apps
      };
    return (
        <div className='min-h-screen bg-[var(--bg-page-color)]'>
            <BackNavigationButton title="Wallet" />
        <div className='p-2 border-2 border-[var(--primary-color)] rounded-lg' onClick={handlePayment}>DO Payment</div>
            {/* Balance Card */}
            <div className='p-4'>
                <div className='bg-gradient-to-r from-[var(--primary-color)] to-purple-600 rounded-xl p-6 text-white mb-6'>
                    <div className='flex items-center gap-2 mb-2'>
                        <Wallet className="w-5 h-5 text-white/80" />
                        <p className='text-white/80'>Available Balance</p>
                    </div>
                    <div className="flex items-baseline gap-1">
                        <IndianRupee className="w-6 h-6" />
                        <span className="text-3xl font-semibold">{userWallet?.balance}</span>
                    </div>
                </div>

                {/* QR Code Section */}
                <div className="bg-white rounded-xl overflow-hidden mb-6">
                    <div className="bg-[var(--primary-color)]/5 p-4 border-b border-[var(--primary-color)]/10">
                        <h3 className="font-medium text-[var(--primary-color)]">Scan QR to Pay</h3>
                    </div>
                    <div className="p-6 flex flex-col items-center">
                        <div className="bg-white p-3 rounded-xl shadow-lg mb-4">
                            <img
                                src="/hostego_payment_qr.jpg"
                                alt="Payment QR Code"
                                className="w-[250px] h-[250px] object-cover rounded-lg"
                            />
                        </div>
                        <p className="text-sm text-gray-600">Scan to make payment via UPI</p>
                    </div>
                </div>

                {/* Add Money Form */}
                <div className='bg-white rounded-xl overflow-hidden'>
                    <div className="bg-[var(--primary-color)]/5 p-4 border-b border-[var(--primary-color)]/10">
                        <h3 className="font-medium text-[var(--primary-color)]">Add Money to Wallet</h3>
                    </div>
                    <form className='p-6 space-y-6' onSubmit={handleWalletTransactionSubmit}>
                        {/* Amount Input */}
                        <div className="relative">
                            <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
                                Amount
                            </label>
                            <div className="flex items-center border-2 border-[var(--primary-color)] rounded-lg px-4 py-3">
                                <span className="text-gray-700">₹</span>
                                <input
                                    type="number"
                                    placeholder="Enter amount (min ₹100)"
                                    value={walletDetails?.amount || ''}
                                    onChange={(e) => {
                                        const value = e.target.value ? Number(e.target.value) : ""
                                        setWalletDetails({ ...walletDetails, amount: value })
                                    }}
                                    className="ml-2 outline-none bg-transparent w-full"
                                    min="100"
                                />
                            </div>
                        </div>

                        {/* Transaction ID Input */}
                        <div className="relative">
                            <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
                                UPI Transaction ID
                            </label>
                            <input
                                type="text"
                                placeholder="Enter UTR number"
                                value={walletDetails?.unique_transaction_id}
                                onChange={(e) => setWalletDetails({ ...walletDetails, unique_transaction_id: e.target.value })}
                                className="w-full border-2 border-[var(--primary-color)] rounded-lg px-4 py-3 outline-none"
                            />
                        </div>

                        {/* Screenshot Upload */}
                        <div className="relative">
                            <label className="absolute text-[var(--primary-color)] text-sm -top-3 left-3 bg-white px-1">
                                Payment Screenshot
                            </label>
                            <div className="border-2 border-[var(--primary-color)] rounded-lg overflow-hidden">
                                <label
                                    htmlFor="screenshot-upload"
                                    className="flex items-center justify-between px-4 py-3 cursor-pointer hover:bg-[var(--primary-color)]/5 transition-colors"
                                >
                                    <span className="text-gray-600 truncate max-w-[200px] overflow-hidden text-ellipsis">
                                        {paymentScreenShotImgUrl ? paymentScreenShotImgUrl.name : "Upload screenshot"}
                                    </span>
                                    <Upload className="text-[var(--primary-color)]" />
                                </label>
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleFileChange}
                                    className="hidden"
                                    id="screenshot-upload"
                                />
                            </div>
                        </div>

                        {/* Submit Button */}
                        <HostegoButton
                            isLoading={walletTransactionCreationLoading}
                            text={`Add ₹${walletDetails?.amount || '0'} to Wallet`}
                            className="w-full"
                            onClick={handleWalletTransactionSubmit}
                        />
                    </form>
                </div>
            </div>
        </div>
    )
}

export default Page
