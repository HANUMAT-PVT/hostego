'use client'
import React, { useState } from 'react'
import { Upload } from "lucide-react"
import BackNavigationButton from '../components/BackNavigationButton'
import { uploadToS3Bucket } from '../lib/aws'

const Page = () => {
    const [amount, setamount] = useState(500);
    const [utrId, setUtrId] = useState("");
    const [selectedFile, setSelectedFile] = useState(null)
    const [screenshot, setScreenshot] = useState(null);

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        setScreenshot(file)


        if (file) {
            const url = URL.createObjectURL(file)
            setSelectedFile(url);
        }
    };


    const handleWalletTransactionSubmit = async (e) => {
        try {
            e.preventDefault()
            if (!selectedFile) return
            const imageUrl = await uploadToS3Bucket(screenshot);
            console.log(imageUrl, "image_url")
            alert("Upload successs")
            setScreenshot(imageUrl)
        } catch (error) {
            alert("error")
        }
    }

    return (
        <div className='bg-[var(--bg-page-color)]'>
            <BackNavigationButton title={"Wallet"} />
            <div className='p-4 flex flex-col gap-4'>
                <div className='bg-white rounded-md w-full p-2 '>
                    <p className='mb-2 text-sm'>Balance</p>
                    <p className='font-normal text-2xl'>₹2560</p>
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
                                        value={amount}
                                        onChange={(e) => setamount(e.target.value)}
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
                                        value={utrId}
                                        onChange={(e) => setUtrId(e.target.value)}
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
                                            {screenshot ? screenshot.name : "Payment screenshot "}
                                        </span>
                                        <Upload className="text-[#655df0]" />
                                    </label>
                                </div>
                            </div>

                            {/* Submit Button */}
                            <button type='submit' className=' font-normal  w-full bg-[var(--primary-color)] text-sm rounded-full p-2 text-white'>
                                Add ₹{amount}
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Page;
