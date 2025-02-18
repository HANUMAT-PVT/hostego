"use client";
// Import auth service this error is coming  in this file "use client";
import React, { useState } from "react";


import PhoneEmailAuthButton from "../../components/PhoneEmailAuth"
const Page = () => {
  const [phoneNumber, setPhoneNumber] = useState("");


  const handlePhoneNumberChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPhoneNumber(e.target.value);
  };

  return (
    <div className="relative min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 text-white flex items-center justify-center">
      <div className="flex flex-col absolute justify-between bg-white h-[75vh] bottom-0 max-w-[400px] px-6 py-6 rounded-t-2xl">
        <div className="flex flex-col gap-6">
          <p className="text-black font-semibold text-xl">Enter your number</p>
          <div className="relative">
            <label className="absolute text-[#655df0] font-semibold -top-3 left-3 bg-white px-1 text-sm">
              Mobile Number
            </label>
            <div className="flex font-semibold items-center border border-[#655df0] max-w-[400px] rounded-xl px-4 py-3 w-full">
              <span className="text-gray-700">+91</span>
              <span className="text-gray-300 ml-2">|</span>
              <input
                type="number"
                placeholder="Enter your number"
                value={phoneNumber}
                onChange={handlePhoneNumberChange}
                className="ml-2 outline-none bg-transparent cursor-pointer w-[200px]"
              />
            </div>
          </div>
        </div>
        <div className="flex flex-col gap-3">
          <PhoneEmailAuthButton />
          

          <p className="text-gray-600">
            By clicking, I accept the{" "}
            <a className="font-bold underline" href="#">
              terms of service
            </a>{" "}
            and{" "}
            <a className="font-bold underline" href="#">
              privacy policy
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Page;
