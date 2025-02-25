'use client'

import { useState } from 'react'
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle, TransitionChild } from '@headlessui/react'
import { CrossIcon, Home, X } from 'lucide-react'


export default function AddresList({ openAddressList, setOpenAddressList, sendSelectedAddress }) {
  const [open, setOpen] = useState(false)
  const userAddress = [
    {
      id: 1,
      heading: "Home",
      street: 'Room no. 1115,Zakir-A,Chandigarh University'
    },
    {
      id: 2,
      heading: "Friend",
      street: 'Room no. 1118,Zakir-A,Chandigarh University'
    }
  ]

  return (
    <Dialog open={openAddressList} onClose={() => setOpenAddressList(!openAddressList)} className="relative z-10 ">
      <DialogBackdrop
        transition
        className="fixed inset-0 bg-gray-500/75 transition-opacity duration-500 ease-in-out data-closed:opacity-0"
      />

      {/* Fullscreen container positioned at the bottom */}
      <div className="fixed inset-0 overflow-hidden">
        <div className="absolute inset-0 overflow-hidden">
          <div className="pointer-events-none fixed bottom-0 left-0 flex w-full max-h-screen">

            {/* Panel now opens from bottom to top */}
            <DialogPanel
              transition
              className="pointer-events-auto relative w-full max-h-[90vh] transform transition duration-500 ease-in-out data-closed:translate-y-full sm:duration-700 bg-white rounded-t-xl shadow-xl"
            >
              <TransitionChild>
                <div className="absolute top-4 right-4 duration-500 ease-in-out data-closed:opacity-0">
                  <button
                    type="button"
                    onClick={() => setOpenAddressList(!openAddressList)}
                    className="relative rounded-md text-gray-300 hover:text-black focus:ring-2 focus:ring-black focus:outline-hidden"
                  >
                    <span className="absolute -inset-2.5" />
                    <span className="sr-only">Close panel</span>
                    <X className="size-6" />
                  </button>
                </div>
              </TransitionChild>

              {/* Content inside the modal */}
              <div className="flex h-full flex-col overflow-y-auto py-6 bg-[var(--bg-page-color)]">
                <div className="px-4 sm:px-6">
                  <DialogTitle className="text-xl font-semibold text-gray-900">Select delivery location</DialogTitle>
                </div>
                <div className="relative mt-6 flex-1 px-4 ">
                  {/* Add your content here */}
                  <p className='text-md text-gray-500 '>Your saved addresses</p>
                  <div className="flex flex-col gap-4 mt-4">
                    {/* Address Ist */}
                    {
                      userAddress?.map((el) => <div key={el?.id} onClick={() => { sendSelectedAddress(el), setOpenAddressList(!openAddressList) }} className='address-item flex items-center rounded-md  gap-4 cursor-pointer p-2 bg-white '>
                        <div className='bg-[var(--bg-page-color)] p-2 w-[40px] h-[40px]  flex justify-center items-center rounded-full '>
                          <Home size={20} className='text-[var(--primary-color)]' />
                        </div>
                        <div className=''>
                          <p className='text-md font-semibold'>{el?.heading}</p>
                          <p className='text-sm '>{el?.street}</p>
                        </div>
                      </div>)
                    }


                  </div>

                </div>
              </div>
            </DialogPanel>
          </div>
        </div>
      </div>
    </Dialog>
  )
}
