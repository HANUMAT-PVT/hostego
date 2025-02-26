const ProductCardSkeleton = () => {
  return (
    <div className='product-item flex flex-col gap-2 w-[120px] h-[180px] animate-pulse'>
      <div className='relative p-2 flex flex-col items-center'>
        {/* Image skeleton */}
        <div className='w-[80px] h-[80px] bg-gray-200 rounded-lg' />
        {/* Add button skeleton */}
        
      </div>

      {/* Weight tag skeleton */}
      <div className='flex gap-2'>
        <div className='w-[30px] h-[20px] bg-gray-200 rounded-lg' />
      </div>

      {/* Product name skeleton */}
      <div className='w-full h-[15px] bg-gray-200 rounded' />

      {/* Price skeleton */}
      <div className='w-[40px] h-[15px] bg-gray-200 rounded' />
    </div>
  )
}

export default ProductCardSkeleton; 