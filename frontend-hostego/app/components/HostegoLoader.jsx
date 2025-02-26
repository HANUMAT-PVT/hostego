'use client'
import React from 'react'

const HostegoLoader = () => {
    return (
        <div className="min-h-screen bg-[var(--bg-page-color)]">
            {/* Top Bar Skeleton */}
            <div className="animate-pulse">
                <div className="h-14 bg-white shadow-sm flex items-center justify-between px-4">
                    <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-gray-200 rounded-full"></div>
                        <div className="w-24 h-4 bg-gray-200 rounded"></div>
                    </div>
                    <div className="w-8 h-8 bg-gray-200 rounded-full"></div>
                </div>
            </div>

            {/* Content Skeletons */}
            <div className="p-4 space-y-4">
                {/* Search Bar Skeleton */}
                <div className="w-full h-12 bg-white rounded-lg animate-pulse"></div>

                {/* Category Pills Skeleton */}
                <div className="flex gap-3 overflow-hidden py-2">
                    {[1, 2, 3, 4].map((i) => (
                        <div key={i} className="w-20 h-8 bg-white rounded-full animate-pulse"></div>
                    ))}
                </div>

                {/* Product Card Skeletons */}
                <div className="grid grid-cols-2 gap-4">
                    {[1, 2, 3, 4].map((i) => (
                        <div key={i} className="bg-white p-4 rounded-lg animate-pulse">
                            <div className="w-full h-32 bg-gray-200 rounded-lg mb-3"></div>
                            <div className="w-2/3 h-4 bg-gray-200 rounded mb-2"></div>
                            <div className="w-1/2 h-4 bg-gray-200 rounded"></div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Center Loader */}
            <div className="fixed inset-0 bg-white/60 backdrop-blur-xs flex items-center justify-center">
                <div className="relative flex flex-col items-center">
                    {/* Main spinner */}
                    <div className="w-16 h-16 relative">
                        {/* Outer spinning ring */}
                        <div className="absolute inset-0 rounded-full border-[3px] border-transparent border-t-[var(--primary-color)] animate-spin"></div>

                        {/* Middle ring with reverse spin */}
                        <div className="absolute inset-[3px] rounded-full border-[3px] border-transparent border-t-purple-500 animate-spin"
                            style={{ animationDirection: 'reverse', animationDuration: '0.6s' }}>
                        </div>

                        {/* Inner pulsing circle */}
                        <div className="absolute inset-[6px] rounded-full bg-gradient-to-tr from-[var(--primary-color)] to-purple-500 animate-pulse">
                            {/* Center dot */}
                            <div className="absolute inset-[35%] bg-white rounded-full shadow-lg"></div>
                        </div>
                    </div>

                    {/* Loading text */}
                    <div className="mt-6 bg-white/80 px-6 py-2 rounded-full shadow-lg">
                        <p className="text-[var(--primary-color)] font-medium tracking-wide">
                            Loading
                            <span className="inline-block w-1 animate-bounce delay-100">.</span>
                            <span className="inline-block w-1 animate-bounce delay-200">.</span>
                            <span className="inline-block w-1 animate-bounce delay-300">.</span>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default HostegoLoader
