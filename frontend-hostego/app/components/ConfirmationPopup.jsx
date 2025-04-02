'use client'
import React from 'react'
import { AlertCircle, X } from 'lucide-react'

const ConfirmationPopup = ({
    title = "Are you sure?",
    message = "This action cannot be undone.",
    confirmText = "Confirm",
    cancelText = "Cancel",
    variant = 'danger', // or 'warning', 'info'
    isOpen = false,
    onConfirm,
    onCancel
}) => {
    if (!isOpen) return null

    const variants = {
        danger: {
            icon: 'text-red-500',
            confirmButton: 'bg-red-500 hover:bg-red-600 focus:ring-red-200',
            title: 'text-red-500'
        },
        warning: {
            icon: 'text-yellow-500',
            confirmButton: 'bg-yellow-500 hover:bg-yellow-600 focus:ring-yellow-200',
            title: 'text-yellow-500'
        },
        info: {
            icon: 'text-blue-500',
            confirmButton: 'bg-blue-500 hover:bg-blue-600 focus:ring-blue-200',
            title: 'text-blue-500'
        }
    }

    const currentVariant = variants[variant]

    return (
        <div className="fixed inset-0 z-50 overflow-y-auto">
            {/* Backdrop */}
            <div className="fixed inset-0 bg-black/50 backdrop-blur-sm transition-opacity" />

            {/* Dialog */}
            <div className="flex min-h-full items-center justify-center p-4">
                <div className="relative transform overflow-hidden rounded-2xl bg-white shadow-xl transition-all max-w-md w-full mx-4 animate-scale-up">
                    {/* Close button */}
                    <button
                        onClick={onCancel}
                        className="absolute right-4 top-4 p-1 rounded-full hover:bg-gray-100 transition-colors"
                    >
                        <X className="w-5 h-5 text-gray-500" />
                    </button>

                    <div className="p-6">
                        {/* Icon */}
                        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-50">
                            <AlertCircle className={`h-6 w-6 ${currentVariant?.icon}`} />
                        </div>

                        {/* Content */}
                        <div className="mt-4 text-center">
                            <h3 className={`text-sm font-semibold ${currentVariant?.title}`}>
                                {title}
                            </h3>
                            <p className="mt-2 text-xs text-gray-600">
                                {message}
                            </p>
                        </div>

                        {/* Buttons */}
                        <div className="mt-6 flex gap-3">
                            <button
                                onClick={onCancel}
                                className="flex-1 px-4 py-2.5 rounded-lg border-2 border-gray-200 
                                         text-gray-700 font-medium hover:bg-gray-50 text-sm 
                                         focus:outline-none focus:ring-2 focus:ring-gray-200 
                                         transition-colors"
                            >
                                {cancelText}
                            </button>
                            <button
                                onClick={onConfirm}
                                className={`flex-1 px-4 py-2.5 rounded-lg text-white font-medium text-sm 
                                          focus:outline-none focus:ring-2 transition-colors
                                          ${currentVariant?.confirmButton}`}
                            >
                                {confirmText}
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default ConfirmationPopup

// Usage Example:
/*
const [showConfirmation, setShowConfirmation] = useState(false)

<ConfirmationPopup 
    isOpen={showConfirmation}
    title="Delete Item"
    message="Are you sure you want to delete this item? This action cannot be undone."
    confirmText="Delete"
    cancelText="Cancel"
    variant="danger"
    onConfirm={() => {
        // Handle confirmation
        setShowConfirmation(false)
    }}
    onCancel={() => setShowConfirmation(false)}
/>
*/
