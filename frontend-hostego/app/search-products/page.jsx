"use client"
import React from 'react'
import SearchComponent from '../components/SearchComponent'
import CartFloatingButton from '../components/CartFloatingButton'

const page = () => {
    return (
        <div>
            <div className='gradient-background py-4'>
                <SearchComponent sendSearchValue={(e)=>console.log(e)}  />
                    <CartFloatingButton/>
            </div>
        </div>
    )
}

export default page
