import React from 'react';
import { Loader2 } from 'lucide-react';

const HostegoButton = ({ text, onClick, isLoading, type = "button" }) => {
  return (
    <button 
      type={type}
      onClick={!isLoading ? onClick : undefined} 
      className={`w-full text-md font-normal bg-[var(--primary-color)] text-white rounded-md p-2 flex justify-center items-center transition ${
        isLoading ? "opacity-50 cursor-not-allowed" : "hover:bg-opacity-90"
      }`}
      disabled={isLoading}
    >
      {isLoading ? <Loader2 className="animate-spin" size={22} /> : text}
    </button>
  );
}

export default HostegoButton;
