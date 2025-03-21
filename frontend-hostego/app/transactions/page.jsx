"use client"
import React, { useEffect, useState } from 'react';
import { ArrowUpRight, ArrowDownRight, CheckCircle, Clipboard, Clock, XCircle, RefreshCw } from 'lucide-react';
import BackNavigationButton from '../components/BackNavigationButton';
import axiosClient from "../utils/axiosClient"
import LoadMoreData from '../components/LoadMoreData';
import TransactionCard from '../components/TransactionCard';

const Transactions = () => {
  const [paymentTransactions, setPaymentTransactions] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)

  const [hasMore, setHasMore] = useState(true)
  const ITEMS_PER_PAGE = 10

  useEffect(() => {
    fetchUserAccountTransactions()
  }, [currentPage])

  const fetchUserAccountTransactions = async () => {
    try {

      let { data } = await axiosClient.get(`/api/wallet/transactions?page=${currentPage}&limit=${ITEMS_PER_PAGE}`)
      setPaymentTransactions(prev => currentPage === 1 ? data : [...prev, ...data])

      setHasMore(data.length < ITEMS_PER_PAGE ? false : true)
    } catch (error) {
      console.error("Error fetching transactions:", error)
    } finally {
      setIsLoading(false)
    }
  }

  const loadMore = () => {
    if (!isLoading && hasMore) {
      setCurrentPage(prev => prev + 1)
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

  const getTransactionIcon = (type) => {
    switch (type?.toLowerCase()) {
      case 'credit':
        return <ArrowDownRight size={20} className="text-green-600" />;
      case 'debit':
        return <ArrowUpRight size={20} className="text-red-500" />;
      case 'refund':
        return <RefreshCw size={20} className="text-blue-500" />;
      default:
        return <ArrowDownRight size={20} />;
    }
  };

  const getTransactionColor = (type) => {
    switch (type?.toLowerCase()) {
      case 'credit':
        return {
          bg: 'bg-green-100',
          text: 'text-green-600',
          label: 'Money Added'
        };
      case 'debit':
        return {
          bg: 'bg-red-100',
          text: 'text-red-500',
          label: 'Order Payment'
        };
      case 'refund':
        return {
          bg: 'bg-blue-100',
          text: 'text-blue-500',
          label: 'Refund'
        };
      default:
        return {
          bg: 'bg-gray-100',
          text: 'text-gray-600',
          label: 'Transaction'
        };
    }
  };

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
        {paymentTransactions.map((transaction) => (
          <TransactionCard
            key={transaction.transaction_id}
            transaction={transaction}
          />
        ))}

        {hasMore && (
          <LoadMoreData loadMore={loadMore} isLoading={isLoading} />
        )}
      </div>
    </div>
  );
};

export default Transactions;
