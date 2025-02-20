import React from 'react'

const ProductCard = ({ img, src,myKey }) => {
    const productTags = [{ id: 1, title: "48g" }, { id: 2, title: "Spicy Treat" }]
    return (
        <div key={myKey} className='product-item flex flex-col gap-2  w-[140px] '>
            <div className='relative  p-2 flex flex-col items-center'>
                <img className='w-[90px]'
                    src={"https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"}
                    alt={'Uncle chips'}
                />
                <button className='absolute bottom-0 right-0 bg-white text-sm text-green-600 px-2 py-1 border-2 border-green-600 rounded-md font-semibold'>
                    ADD
                </button>
            </div>


            <div className='flex gap-2'>
                {productTags?.map((el) => <p key={el?.id} className='text-[10px] px-2 py-1 rounded-lg bg-gray-200 font-medium text-gray-800'>{el?.title}</p>)}
            </div>
            <p className='product-title text-sm font-medium'>Uncle Chipps Spicy Treat Flavour Potato Chips</p>
            <p className='product-price text-left text-sm font-semibold'>₹ 20</p>
        </div>
    )
}

export default ProductCard
