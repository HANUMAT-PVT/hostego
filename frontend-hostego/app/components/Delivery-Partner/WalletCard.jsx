import { useState } from 'react';
import { Wallet, ChevronUp, ChevronDown } from 'lucide-react';
import TransactionCard from '../TransactionCard';

const WalletCard = ({ walletData }) => {
    const [showTransactions, setShowTransactions] = useState(false);

    return (
        <div className="bg-white rounded-xl shadow-sm p-6">
            <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                    <div className="p-3 rounded-full bg-blue-100">
                        <Wallet className="w-6 h-6 text-blue-600" />
                    </div>
                    <div>
                        <h3 className="font-medium">Wallet Balance</h3>
                        <p className="text-2xl font-bold text-blue-600">₹{(walletData?.balance || 0).toFixed(1)}</p>
                    </div>
                </div>
            </div>

            <div className="mt-4 border-t pt-4 ">
                <div onClick={() => setShowTransactions(!showTransactions)} className="flex items-center justify-between">
                    <h4 className="text-sm font-medium mb-3">Recent Transactions</h4>
                    <button >
                        {showTransactions ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
                    </button>
                </div>
                {showTransactions && (
                    <div className="space-y-4 overflow-y-auto max-h-[300px]">
                        {walletData?.recent_transactions?.map((transaction) => (
                            <TransactionCard
                                label={transaction?.transaction_type === "credit" ? "Order Earning" : "Withdrawal Request"}
                                key={transaction.transaction_id}
                                transaction={transaction}
                            />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default WalletCard; 