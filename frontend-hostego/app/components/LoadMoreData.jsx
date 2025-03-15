import React from 'react'

const LoadMoreData = ({loadMore, isLoading}) => {
    return (
        <button
            onClick={loadMore}
            disabled={isLoading}
            className='mt-2 w-full py-3 bg-[var(--primary-color)] text-white rounded-lg hover:bg-[var(--primary-color)] disabled:bg-gray-300 disabled:cursor-not-allowed'
        >
            {isLoading ? (
                <div className='flex items-center justify-center gap-2'>
                    <div className='animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-white'></div>
                    Loading...
                </div>
            ) : (
                'Load More'
            )}
        </button>
    )
}

export default LoadMoreData
