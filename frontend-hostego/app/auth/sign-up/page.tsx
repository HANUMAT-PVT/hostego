import React from "react";
import HostegoButton from "../../components/HostegoButton";
const Page = () => {
  return (
    <div className="relative flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 text-white ">
      {/* Background Blurs */}
      <div className="flex flex-col justify-between bg-white w-full h-[80vh] absolute bottom-0 px-4 py-6 w-full  rounded-t-2xl">
        <div className="flex flex-col gap-6 w-full">
          <p className="text-black font-semibold text-xl">Enter your number</p>
          <div className="relative w-full max-w-sm">
            <label className="absolute text-[#655df0]  font-semibold -top-3 left-3 bg-white px-1 text-sm">
              Mobile Number
            </label>
            <div className="flex font-semibold items-center text-center border border-[#655df0] w-full rounded-xl px-4 py-3">
              <span className="text-gray-700  text-center ">
                +91
              </span>
              <span className="text-gray-300 ml-2">|</span>
              <input
                type="number"
                placeholder="Enter your number"
                className="flex-1 ml-2 outline-none bg-transparent cursor-pointer"
              />
            </div>
          </div>
        </div>
        <div className="flex flex-col gap-3">
          <HostegoButton text={"Continue"}/>
          <p className="text-gray-600">By clicking, I accept the <a className="font-bold underline" href="">terms of service</a> and <a className="font-bold underline" href="">privacy policy</a></p>
        </div>
      </div>
    </div>
  );
};

export default Page;
