import React, { useState } from 'react'

const CartItem = () => {
    const [cartItemQuantity, setCartItemQuantity] = useState(1)
    return (

        <div className='w-full flex  gap-5 p-2 items-center  border-b  bg-white  justify-between'>

            <div className='flex gap-4'>
                <div className='bg-[var(--bg-page-color)] rounded-md p-1 flex items-center justify-center'>
                    <img className='min-w-[50px] max-w-[50px] max-h-[50px] min-h-[50px]'
                        src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                        alt={'Uncle chips'}
                    />
                </div>
                <div className='flex flex-col gap-2'>
                    <p className='text-sm font-normal w-[130px]'>Lay's India's Magic Masala Potato Chips</p>
                    <p className='text-xs text-gray-500 font-semibold'>48 g </p>
                </div>
            </div>



            <div className='flex flex-col gap-2 text-right '>

                <div className=' bg-green-700 h-[30px]  text-sm flex gap-2 text-white text-sm font-semibold px-3 py-1 rounded-md '>
                    <button disabled={cartItemQuantity === 1} onClick={() => setCartItemQuantity(cartItemQuantity - 1)}>-</button>
                    <button className='w-[15px] text-sm'>{cartItemQuantity}</button>
                    <button onClick={() => setCartItemQuantity(cartItemQuantity + 1)}>+</button>
                </div>
                <p className='text-md font-semibold mt-1'>₹{cartItemQuantity * 20}</p>
            </div>

        </div>


    )
}

export default CartItem
