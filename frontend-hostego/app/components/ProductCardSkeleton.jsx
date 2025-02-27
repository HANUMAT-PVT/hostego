const ProductCardSkeleton = () => {
  return (
    <div className='product-item flex flex-col gap-2   h-[180px] animate-pulse p-2 rounded-md'>
      <div className='relative p-2 flex flex-col items-center'>
        {/* Image skeleton */}
        <div className='w-[60px] h-[60px] bg-gray-200 rounded-md' />
        {/* Add button skeleton */}
        
      </div>

      {/* Weight tag skeleton */}
      <div className='flex gap-2'>
        <div className='w-[30px] h-[20px] bg-gray-200 rounded-lg' />
      </div>

      {/* Product name skeleton */}
      <div className='w-full h-[15px] bg-gray-200 rounded' />

      {/* Price skeleton */}
      <div className='w-[30px] h-[10px] bg-gray-200 rounded' />
    </div>
  )
}

export default ProductCardSkeleton; 