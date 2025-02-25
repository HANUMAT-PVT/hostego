"use client"
import React, { useEffect, useState } from 'react';
import { ArrowUpRight, ArrowDownRight, CheckCircle, Clipboard, Clock, XCircle } from 'lucide-react';
import BackNavigationButton from '../components/BackNavigationButton';
import axiosClient from "../utils/axiosClient"

const Transactions = () => {
  const [paymentTransactions, setPaymentTransactions] = useState([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    fetchUserAccountTransactions()
  }, [])

  const fetchUserAccountTransactions = async () => {
    try {
      setIsLoading(true)
      let { data } = await axiosClient.get("/api/wallet/transactions")
      setPaymentTransactions(data)
    } catch (error) {
      console.log(error, "error")
    } finally {
      setIsLoading(false)
    }
  }

  const getStatusIcon = (status) => {
    switch (status?.toLowerCase()) {
      case 'success':
        return <CheckCircle size={14} className='text-green-500' />
      case 'pending':
        return <Clock size={14} className='text-orange-500' />
      case 'failed':
        return <XCircle size={14} className='text-red-500' />
      default:
        return <Clock size={14} className='text-gray-500' />
    }
  }

  if (isLoading) {
    return (
      <div className='bg-[var(--bg-page-color)] min-h-screen'>
        <BackNavigationButton title={"Transactions"} />
        <div className='flex flex-col items-center justify-center h-[80vh] px-4'>
          <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500'></div>
          <p className='text-gray-500 mt-4'>Loading transactions...</p>
        </div>
      </div>
    )
  }

  if (!paymentTransactions.length) {
    return (
      <div className='bg-[var(--bg-page-color)] min-h-screen'>
        <BackNavigationButton title={"Transactions"} />
        <div className='flex flex-col items-center justify-center h-[80vh] px-4'>
          <div className='bg-gray-100 rounded-full p-6 mb-4'>
            <ArrowUpRight size={40} className='text-gray-400' />
          </div>
          <h3 className='text-xl font-semibold text-gray-700 mb-2'>No Transactions Yet</h3>
          <p className='text-gray-500 text-center'>Your transaction history will appear here once you make your first transaction.</p>
        </div>
      </div>
    )
  }

  return (
    <div className='bg-[var(--bg-page-color)] min-h-screen'>
      <BackNavigationButton title={"Transactions"} />
      <div className='flex flex-col gap-4 px-4 py-4'>
        {paymentTransactions.map((el) => (
          <div key={el?.transaction_id} className='flex flex-col gap-2 bg-white px-4 py-3 rounded-lg shadow-md hover:shadow-lg transition-shadow'>
            <div className='flex flex-col gap-1 mb-2 '>
              <p className='text-gray-500 text-xs'>Transaction ID:</p>
              <div className='flex items-center gap-2'>
                <p className='font-normal text-xs'>{el?.transaction_id}</p>
                <Clipboard
                  size={16}
                  className='text-gray-400 cursor-pointer text-right text-now-wrap hover:text-gray-600'
                  onClick={() => {
                    navigator.clipboard.writeText(el.transaction_id)
                    // Optional: Add a toast notification here
                  }}
                />
              </div>
            </div>
            <div className='flex justify-between items-center'>
              <div className='flex items-center gap-3'>
                <div className={`p-2 rounded-full ${el?.transaction_type === "credit" ? "bg-green-100 text-green-600" : "bg-red-100 text-red-500"}`}>
                  {el.transaction_type === "credit" ? <ArrowUpRight size={20} /> : <ArrowDownRight size={20} />}
                </div>
                <div>
                  <p className='text-sm font-semibold'>{el.transaction_type.toUpperCase()}</p>
                  <p className='text-xs text-gray-500'>{new Date(el.created_at)?.toLocaleString()}</p>
                </div>
              </div>
              <div className='text-right'>
                <p className={`text-lg font-bold ${el.transaction_type === "credit" ? "text-green-600" : "text-red-500"}`}>₹{el.amount}</p>
                <p className='text-xs text-gray-500 flex items-center gap-1 justify-end'>
                  {getStatusIcon(el?.transaction_status)} {el?.transaction_status}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Transactions;
