'use client'
import { ArrowRight } from "lucide-react";
import React, { useState, useRef,useEffect } from "react"

const SliderStatusTracker = ({ onConfirm, text }) => {
    const [startX, setStartX] = useState(0);
    const [currentX, setCurrentX] = useState(0);
    const [isDragging, setIsDragging] = useState(false);
    const sliderRef = useRef(null);

    const handleStart = (e) => {
        const clientX = e.type === 'mousedown' ? e.clientX : e.touches[0].clientX;
        setStartX(clientX);
        setIsDragging(true);
    };

    const handleMove = (e) => {
        if (!isDragging) return;
        e.preventDefault();

        const clientX = e.type === 'mousemove' ? e.clientX : e.touches[0].clientX;
        const sliderWidth = sliderRef.current.offsetWidth;
        const diff = clientX - startX;
        const newX = Math.max(0, Math.min(diff, sliderWidth - 60));
        setCurrentX(newX);

        if (newX > sliderWidth * 0.75) {
            setIsDragging(false);
            setCurrentX(0);
            onConfirm();
        }
    };

    useEffect(() => {
        if (isDragging) {
            window.addEventListener('mousemove', handleMove);
            window.addEventListener('touchmove', handleMove, { passive: false });
            window.addEventListener('mouseup', () => setIsDragging(false));
            window.addEventListener('touchend', () => setIsDragging(false));
        }
        return () => {
            window.removeEventListener('mousemove', handleMove);
            window.removeEventListener('touchmove', handleMove);
            window.removeEventListener('mouseup', () => setIsDragging(false));
            window.removeEventListener('touchend', () => setIsDragging(false));
        };
    }, [isDragging]);

    useEffect(() => {
        if (!isDragging) setCurrentX(0);
    }, [isDragging]);

    return (
        <div
            ref={sliderRef}
            className="   relative h-14 rounded-full overflow-hidden animate-slide-up bg-[var(--primary-color)]"
        >
            {/* Text */}
            <div className="absolute inset-0 flex items-center justify-center">
                <span className={`text-sm font-medium text-white transition-opacity ${isDragging ? 'opacity-0' : 'opacity-100'
                    }`}>
                    {text}
                </span>
            </div>

            {/* Slider Button */}
            <div
                className="absolute left-1 top-1 bottom-1 w-12 rounded-full bg-white flex items-center justify-center cursor-grab"
                style={{
                    transform: `translateX(${currentX}px)`,
                    transition: isDragging ? 'none' : 'all 0.4s ease'
                }}
                onMouseDown={handleStart}
                onTouchStart={handleStart}
            >
                <ArrowRight className="w-5 h-5 text-[var(--primary-color)]" />
            </div>
        </div>
    );
};

export default SliderStatusTracker