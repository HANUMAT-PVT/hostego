"use client"
import { ArrowLeft, Search } from 'lucide-react'
import { useRouter } from 'next/navigation'
import React, { useState } from 'react'

const SearchComponent = ({ viewOnly, sendSearchValue }) => {
    const router = useRouter();
    const [value, setInputValue] = useState("");

    const handleInputChange = (e) => {
        const newValue = e.target.value;
        setInputValue(newValue);
        sendSearchValue(newValue);
    };

    return (
        <div
            className="flex items-center w-[90vw] bg-white m-auto border-2 border-gray-400 rounded-lg px-2 py-2 transition-all"
        >
            {viewOnly ? (
                <div
                    onClick={() => router.push("/search-products")}
                    className="flex items-center w-full cursor-pointer"
                >
                    <Search className="text-gray-800 mr-2 text-bold" size={20} />
                    <span className="text-gray-600 text-sm">
                       Search 'samosa'
                    </span>
                </div>
            ) : (
                <>
                    <ArrowLeft
                        onClick={() => router.back()}
                        className="text-gray-800 mr-2 text-bold cursor-pointer"
                        size={20}
                    />
                    <input
                        type="text"
                        value={value}
                        onChange={handleInputChange}
                        placeholder="Search for samosa, noodles, chips, coke and more..."
                        className="placeholder:text-gray-600 w-full bg-transparent outline-none text-sm font-normal text-gray-600"
                    />
                </>
            )}
        </div>
    );
}

export default SearchComponent;
