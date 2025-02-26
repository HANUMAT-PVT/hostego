'use client'
import React from 'react'
import { CheckCircle2, Loader2, XCircle } from 'lucide-react'

const PaymentStatus = ({ status }) => {
    const statusConfig = {
        processing: {
            title: 'Processing Payment',
            subtitle: 'Please wait while we process your payment',
            icon: <Loader2 className="w-16 h-16 text-[var(--primary-color)] animate-spin" />,
            bgColor: 'bg-white'
        },
        success: {
            title: 'Payment Successful!',
            subtitle: 'Your order has been placed successfully',
            icon: <CheckCircle2 className="w-16 h-16 text-green-500" />,
            bgColor: 'bg-green-50'
        },
        failed: {
            title: 'Payment Failed',
            subtitle: 'Something went wrong. Please try again',
            icon: <XCircle className="w-16 h-16 text-red-500" />,
            bgColor: 'bg-red-50'
        }
    }

    const currentStatus = statusConfig[status]

    return (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
            <div className={`w-[90%] max-w-md ${currentStatus.bgColor} rounded-2xl p-6 mx-4 animate-scale-up`}>
                <div className="flex flex-col items-center text-center">
                    {/* Icon */}
                    <div className="mb-4">
                        {currentStatus.icon}
                    </div>

                    {/* Title */}
                    <h2 className="text-xl font-semibold mb-2">
                        {currentStatus.title}
                    </h2>

                    {/* Subtitle */}
                    <p className="text-gray-600 mb-6">
                        {currentStatus.subtitle}
                    </p>

                    {/* Progress Bar for Processing State */}
                    {status === 'processing' && (
                        <div className="w-full bg-gray-200 rounded-full h-2 mb-4 overflow-hidden">
                            <div className="h-full bg-[var(--primary-color)] animate-progress-indeterminate"></div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

export default PaymentStatus 