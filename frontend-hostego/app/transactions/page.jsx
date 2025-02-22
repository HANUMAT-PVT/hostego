"use client"
import React, { useState } from 'react';
import { ArrowUpRight, ArrowDownRight, CheckCircle, Clipboard, ChevronsRight } from 'lucide-react';
import BackNavigationButton from '../components/BackNavigationButton';

const transactions = [
  {
    transaction_id: "a99a1830-63b4-425e-9ab7-778959f73f07",
    amount: 100,
    transaction_type: "credit",
    created_at: "2025-02-17 02:41:53",
    transaction_status: "success",
  },
  {
    transaction_id: "df7ec025-a0fd-499b-845d-bb761cc4f6bd",
    amount: 2000,
    transaction_type: "credit",
    created_at: "2025-02-17 02:48:16",
    transaction_status: "success",
  },
  {
    transaction_id: "6eb9a1d9-2269-47c7-b960-7d9442400f75",
    amount: 1781,
    transaction_type: "debit",
    created_at: "2025-02-17 02:50:45",
    transaction_status: "success",
  }
];

const Transactions = () => {
  const [accepted, setAccepted] = useState({});
  const [positions, setPositions] = useState({});

  const handleTouchMove = (e, id) => {
    const touch = e.touches[0];
    const newPos = Math.min(Math.max(touch.clientX - 50, 0), 250);
    setPositions((prev) => ({ ...prev, [id]: newPos }));
  };

  const handleTouchEnd = (id) => {
    if (positions[id] > 200) {
      setAccepted((prev) => ({ ...prev, [id]: true }));
    }
    setPositions((prev) => ({ ...prev, [id]: 0 }));
  };

  return (
    <div className='bg-[var(--bg-page-color)]'>
      <BackNavigationButton title={"Transactions"} />
      <div className='flex flex-col gap-4 px-4 py-4'>
        {transactions.map((el) => (
          <div key={el.transaction_id} className='flex flex-col gap-2 bg-white px-4 py-3 rounded-lg shadow-md'>
            <div className='flex justify-between items-center'>
              <p className='text-gray-500 text-xs'>Transaction ID:</p>
              <div className='flex items-center gap-2'>
                <p className='font-normal text-xs'>{el.transaction_id}</p>
                <Clipboard size={16} className='text-gray-400 cursor-pointer text-right text-now-wrap hover:text-gray-600' onClick={() => navigator.clipboard.writeText(el.transaction_id)} />
              </div>
            </div>
            <div className='flex justify-between items-center'>
              <div className='flex items-center gap-3'>
                <div className={`p-2 rounded-full ${el.transaction_type === "credit" ? "bg-green-100 text-green-600" : "bg-red-100 text-red-500"}`}> 
                  {el.transaction_type === "credit" ? <ArrowUpRight size={20} /> : <ArrowDownRight size={20} />}
                </div>
                <div>
                  <p className='text-sm font-semibold'>{el.transaction_type.toUpperCase()}</p>
                  <p className='text-xs text-gray-500'>{new Date(el.created_at).toLocaleString()}</p>
                </div>
              </div>
              <div className='text-right'>
                <p className={`text-lg font-bold ${el.transaction_type === "credit" ? "text-green-600" : "text-red-500"}`}>₹{el.amount}</p>
                <p className='text-xs text-gray-500 flex items-center gap-1'>
                  <CheckCircle size={14} className='text-green-500' /> {el.transaction_status}
                </p>
              </div>
            </div>
            {!accepted[el.transaction_id] && (
              <div className='relative w-full bg-gray-200 rounded-full h-10 mt-3 overflow-hidden'>
                <div
                  className='absolute w-10 h-10 bg-green-500 text-white rounded-full flex items-center justify-center transition-transform duration-300 cursor-pointer active:scale-95'
                  style={{ transform: `translateX(${positions[el.transaction_id] || 0}px)` }}
                  onTouchMove={(e) => handleTouchMove(e, el.transaction_id)}
                  onTouchEnd={() => handleTouchEnd(el.transaction_id)}
                >
                  <ChevronsRight size={20} />
                </div>
                <p className='text-sm text-gray-500 text-center leading-10'>Swipe to Accept</p>
              </div>
            )}
            {accepted[el.transaction_id] && (
              <p className='text-center text-green-600 text-sm font-semibold'>Accepted</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Transactions;
