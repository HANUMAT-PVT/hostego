import { ChevronRight } from 'lucide-react'
import React from 'react'

const CartFloatingButton = () => {
    return (
        <div className='z-[1] min-w-[180px]  px-2 py-2 flex items-center  justify-between rounded-full fixed bottom-[80px] left-1/2 transform -translate-x-1/2  bg-green-600'>
            <div className=' flex flex-col items-center '>
                <img className='w-[40px] rounded-full border-2 '
                    src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                    alt={'Uncle chips'}
                />
               
            </div>
           <div>
           <p className='text-sm text-white font-semibold'>View cart</p>
           <p className='text-xs text-white font-normal'>1 ITEM</p>
           </div>
            <div className='p-2 bg-green-800 w-[40px] h-[40px] flex items-center rounded-full justify-center'>
                <ChevronRight color='white' size={24} />
            </div>
        </div>
    )
}

export default CartFloatingButton
