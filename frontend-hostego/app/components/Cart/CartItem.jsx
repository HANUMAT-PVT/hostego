import React, { useState } from 'react'

const CartItem = () => {
    const [cartItemQuantity, setCartItemQuantity] = useState(1)
    return (

        <div className='flex  gap-5 p-2 items-center border-b  bg-white '>
            <div className='bg-[var(--bg-page-color)] rounded-md p-1'>
                <img className='min-w-[50px] max-w-[50px]'
                    src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                    alt={'Uncle chips'}
                />
            </div>
            <div className='flex flex-col gap-2'>
                <div className='flex gap-7 '>
                    <p className='text-xs font-normal w-[130px]'>Lay's India's Magic Masala Potato Chips</p>
                    <div className=' bg-green-700 h-[30px] flex gap-2 text-white text-xs font-semibold px-2 py-1 rounded-md '>
                        <button disabled={cartItemQuantity === 1} onClick={() => setCartItemQuantity(cartItemQuantity - 1)}>-</button>
                        <button className='w-[15px]'>{cartItemQuantity}</button>
                        <button onClick={() => setCartItemQuantity(cartItemQuantity + 1)}>+</button>
                    </div>
                </div>
                <div className='flex justify-between gap-2'>
                    <p className='text-[11px] text-gray-500 font-light'>48 g </p>
                    <p className='text-xs font-semibold'>₹{cartItemQuantity * 20}</p>
                </div>
            </div>
        </div>


    )
}

export default CartItem
