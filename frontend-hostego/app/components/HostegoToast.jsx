'use client'
import React, { useEffect, useState } from 'react'
import { CheckCircle2, AlertCircle, Info, XCircle, X } from 'lucide-react'

const VARIANTS = {
    success: {
        icon: CheckCircle2,
        className: 'bg-green-50 border-green-200',
        iconClass: 'text-green-500',
        textClass: 'text-green-800',
        progressClass: 'bg-green-500'
    },
    error: {
        icon: XCircle,
        className: 'bg-red-50 border-red-200',
        iconClass: 'text-red-500',
        textClass: 'text-red-800',
        progressClass: 'bg-red-500'
    },
    warning: {
        icon: AlertCircle,
        className: 'bg-yellow-50 border-yellow-200',
        iconClass: 'text-yellow-500',
        textClass: 'text-yellow-800',
        progressClass: 'bg-yellow-500'
    },
    info: {
        icon: Info,
        className: 'bg-blue-50 border-blue-200',
        iconClass: 'text-blue-500',
        textClass: 'text-blue-800',
        progressClass: 'bg-blue-500'
    }
}

const HostegoToast = ({
    message = "This is a toast message",
    variant = 'info',
    show = false,
    onClose
}) => {
    const [isVisible, setIsVisible] = useState(show)
    const config = VARIANTS[variant]
    const Icon = config.icon

    useEffect(() => {
        setIsVisible(show)
        if (show) {
            const timer = setTimeout(() => {
                setIsVisible(false)
                onClose?.()
            }, 3000)
            return () => clearTimeout(timer)
        }
    }, [show, onClose])

    if (!isVisible) return null

    return (
        <div className="fixed top-4 right-4 z-50 animate-slide-up">
            <div className={`relative flex items-center gap-3 min-w-[320px] max-w-md p-4 
                         rounded-xl border shadow-lg ${config.className}`}
            >
                <Icon className={`w-5 h-5 ${config.iconClass}`} />
                <p className={`flex-1 text-sm font-medium ${config.textClass}`}>
                    {message}
                </p>
                <button
                    onClick={() => {
                        setIsVisible(false)
                        onClose?.()
                    }}
                    className="p-1 hover:bg-black/5 rounded-full transition-colors"
                >
                    <X className={`w-4 h-4 ${config.textClass}`} />
                </button>

                {/* Progress bar */}
                <div className="absolute bottom-0 left-0 right-0 h-1 overflow-hidden rounded-b-xl">
                    <div
                        className={`h-full ${config.progressClass} animate-progress`}
                    />
                </div>
            </div>
        </div>
    )
}

// Add these styles to your globals.css
const styles = `
@keyframes slide-up {
    from {
        opacity: 0;
        transform: translateY(16px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes progress {
    from {
        width: 100%;
    }
    to {
        width: 0%;
    }
}

.animate-slide-up {
    animation: slide-up 0.3s ease-out;
}

.animate-progress {
    animation: progress 3s linear forwards;
}
`

export default HostegoToast

// Usage Example:
/*
const [showToast, setShowToast] = useState(false)

<HostegoToast 
    show={showToast}
    variant="success" // or "error", "warning", "info"
    message="Order placed successfully!"
    onClose={() => setShowToast(false)}
/>

// Show toast
setShowToast(true)
*/
