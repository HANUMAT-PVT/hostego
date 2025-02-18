"use client";
import React from "react";
import {
  ChevronRight,
  CreditCard,
  NotebookTabs,
  Package,
  ScrollText,
  Wallet,
} from "lucide-react";
import BackNavigationButton from "../components/BackNavigationButton";


const page = () => {
  return (
    <>
      <BackNavigationButton />
      <div className="home-page-top-bar-section text-white p-4 flex flex-col items-left gap-1">
      <p className="text-xs font-bold">Hostego in</p>
      <p className="text-2xl font-bold">15 minutes </p>
      <p className="text-sm font-medium">Zakir A,Changiarh University</p>
      </div>
      <div className="p-4 mt-3 flex flex-col gap-5">
        <div className="flex flex-col gap-3">
          <p className="text-xl font-medium">My account </p>
          <p className="text-sm font-normal">8264121428</p>
        </div>
        {/* Suggestion Box */}
        <div className=" bg-[#eae8ff]  flex gap-4 justify-between py-4 px-8 rounded-lg">
          <div className="flex flex-col gap-2 items-center text-center ">
            <Wallet size={20} />
            <p className="text-sm font-normal  ">Wallet</p>
          </div>
          <div className="flex flex-col gap-2 items-center text-center ">
            <CreditCard size={20} />
            <p className="text-sm font-normal">Payments</p>
          </div>
        </div>
        {/* Pesonal Information like orders addres  */}
        <div className="flex flex-col gap-3">
          <p className="text-sm font-normal text-gray-600">YOUR INFORMATION</p>
          {/* Your Orders */}
          <div className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Package size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal">Your orders</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Address */}
          <div className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <NotebookTabs size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal">Address</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>

        {/* SECTION 2 */}
        <div className="flex flex-col gap-3">
          <p className="text-sm font-normal text-gray-600">
            PAYMENTS AND TRANSACTIONS
          </p>
          {/* Your Wallet */}
          <div className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal">Wallet</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Payments */}
          <div className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <CreditCard size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal">Payments</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Transactions */}
          <div className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <ScrollText size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal"> Transactions</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>
      </div>
    </>
  );
};

export default page;
