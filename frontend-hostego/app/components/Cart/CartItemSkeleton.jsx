const CartItemSkeleton = () => {
  return (
    <div className="flex items-center gap-4 bg-white p-4 rounded-lg animate-pulse">
      {/* Product image skeleton */}
      <div className="w-[60px] h-[60px] bg-gray-200 rounded-md" />
      
      <div className="flex-1 space-y-2">
        {/* Product name skeleton */}
        <div className="h-4 bg-gray-200 rounded w-3/4" />
        
        {/* Weight skeleton */}
        <div className="h-3 bg-gray-200 rounded w-16" />
        
        <div className="flex items-center justify-between mt-2">
          {/* Price skeleton */}
          <div className="h-4 bg-gray-200 rounded w-20" />
          
          {/* Quantity controls skeleton */}
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gray-200 rounded-full" />
            <div className="w-6 h-6 bg-gray-200 rounded" />
            <div className="w-8 h-8 bg-gray-200 rounded-full" />
          </div>
        </div>
      </div>
    </div>
  );
};

export default CartItemSkeleton; 