import React from 'react';
import { ArrowUpRight, ArrowDownRight, CheckCircle, Clipboard, Clock, XCircle, RefreshCw } from 'lucide-react';

const TransactionCard = ({ transaction,label }) => {
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

  const transactionStyle = getTransactionColor(transaction.transaction_type);

  return (
    <div className='flex flex-col gap-2 bg-white px-4 py-3 rounded-lg shadow-sm hover:shadow-md transition-all'>
      <div className='flex flex-col gap-1 mb-2'>
        <p className='text-gray-500 text-xs'>Transaction ID:</p>
        <div className='flex items-center gap-2'>
          <p className='font-normal text-xs'>{transaction?.transaction_id}</p>
          <Clipboard
            size={16}
            className='text-gray-400 cursor-pointer hover:text-gray-600'
            onClick={(e) => {
              e.stopPropagation();
              navigator.clipboard.writeText(transaction.transaction_id);
            }}
          />
        </div>
      </div>

      <div className='flex justify-between items-center'>
        <div className='flex items-center gap-3'>
          <div className={`p-2 rounded-full ${transactionStyle.bg}`}>
            {getTransactionIcon(transaction.transaction_type)}
          </div>
          <div>
            <p className={`text-sm font-semibold ${transactionStyle.text}`}>
              {label || transactionStyle?.label}
            </p>
            <p className='text-xs text-gray-500'>
              {new Date(transaction.created_at)?.toLocaleString()}
            </p>
          </div>
        </div>

        <div className='text-right'>
          <p className={`text-lg font-bold ${transactionStyle.text}`}>
            {transaction.transaction_type.toLowerCase() === 'debit' ? '-' : '+'}₹{transaction?.amount}
          </p>
          <p className='text-xs text-gray-500 flex items-center gap-1 justify-end'>
            {getStatusIcon(transaction?.transaction_status)} {transaction?.transaction_status}
          </p>
        </div>
      </div>

      {/* Optional: Transaction Note/Description */}
      {transaction?.description && (
        <p className='text-sm text-gray-500 mt-2 border-t pt-2'>
          {transaction?.description}
        </p>
      )}
    </div>
  );
};

export default TransactionCard; 