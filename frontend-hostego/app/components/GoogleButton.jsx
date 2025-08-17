"use client";
import React from 'react';

const GoogleIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" className="w-5 h-5">
        <path fill="#FFC107" d="M43.611 20.083H42V20H24v8h11.303C33.819 32.658 29.277 36 24 36c-6.627 0-12-5.373-12-12s5.373-12 12-12c3.059 0 5.842 1.156 7.948 3.052l5.657-5.657C34.84 6.053 29.728 4 24 4 12.955 4 4 12.955 4 24s8.955 20 20 20 20-8.955 20-20c0-1.341-.138-2.651-.389-3.917z" />
        <path fill="#FF3D00" d="M6.306 14.691l6.571 4.819C14.601 16.108 18.961 12 24 12c3.059 0 5.842 1.156 7.948 3.052l5.657-5.657C34.84 6.053 29.728 4 24 4 16.318 4 9.656 8.337 6.306 14.691z" />
        <path fill="#4CAF50" d="M24 44c5.213 0 9.899-1.977 13.45-5.197l-6.207-5.238C29.277 36 24.735 36 24 36c-5.252 0-9.805-3.361-11.296-8.013l-6.55 5.046C9.466 39.552 16.167 44 24 44z" />
        <path fill="#1976D2" d="M43.611 20.083H42V20H24v8h11.303c-1.347 3.658-4.88 6.274-9.303 6.274-.735 0-1.277 0-1.277 0v7.565s.735 0 1.277 0c8.255 0 15.277-5.565 17.611-13.083L43.611 20.083z" />
    </svg>
);

const GoogleButton = ({ text = 'Sign in with Google', onClick, isLoading = false, className = '' }) => {
    return (
        <button
            type="button"
            onClick={!isLoading ? onClick : undefined}
            disabled={isLoading}
            className={`w-full flex items-center justify-center gap-3 px-4 py-2 bg-white border border-gray-200 rounded-lg shadow-sm hover:bg-gray-50 active:bg-gray-100 disabled:opacity-50 ${className}`}
        >
            <GoogleIcon />
            <span className="text-sm font-medium text-gray-700">{text}</span>
        </button>
    );
};

export default GoogleButton;


