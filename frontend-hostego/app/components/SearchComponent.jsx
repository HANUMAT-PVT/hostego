import { ArrowLeft, Router, Search } from 'lucide-react'
import { useRouter } from 'next/navigation'

import React, { useState } from 'react'

const SearchComponent = ({ viewOnly, sendSearchValue }) => {
    const router = useRouter();
    const [value, setInputValue] = useState(viewOnly ? 'Search "samosa"' : "")

    return (
        <div onClick={() => viewOnly ? router.push("/search-products") : {}} className="flex items-center w-[90vw] bg-white m-auto border-2 gray-400 rounded-lg px-2 py-2 transition-all ">
            {viewOnly ? <Search onChange={(e) => sendSearchValue(e?.target?.value)} className="text-gray-800 mr-2 text-bold" size={20} /> : <ArrowLeft onClick={() => router.back()} className="text-gray-800 mr-2 text-bold" size={20} />}
            <input
                type="text"
                value={value}
                onChange={(e) => { sendSearchValue(e?.target?.value), setInputValue(e.target.value) }}
                placeholder={"Search for samosa, noodlels, chips, coke and m..."}
                className="placeholder:text-gray-600 w-full bg-transparent outline-none text-sm  font-normal text-gray-600"
            />
        </div>
    )
}

export default SearchComponent
