import React from 'react'

const HostegoButton = ({ text, onClick }) => {
    return (
        <button onClick={()=>onClick()} className='bg-[#655df0] w-full p-4 rounded-lg' >
            <p className="font-bold text-xl">{text}</p>
        </button>
    )
}

export default HostegoButton
