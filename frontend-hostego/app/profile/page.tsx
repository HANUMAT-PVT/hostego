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
import { useRouter } from "next/navigation";

const Profile = () => {
  const router = useRouter();
  return (
    <>
      <BackNavigationButton title={"Profile"} />

      <div className="p-4 mt-3 flex flex-col gap-5">
        <div className="flex flex-col gap-3">
          <p className="text-xl font-medium">My account </p>
          <p className="text-sm font-normal">8264121428</p>
        </div>
        {/* Suggestion Box done */}

        <div className=" bg-[#eae8ff]  flex gap-4 justify-between py-4 px-8 rounded-lg">
          <div
            onClick={() => router.push("/wallet")}
            className="flex flex-col gap-2 items-center text-center "
          >
            <Wallet size={20} />
            <p className="text-sm font-normal  ">Wallet</p>
          </div>
          <div   onClick={() => router.push("/transactions")}className="flex flex-col gap-2 items-center text-center ">
            <CreditCard size={20} />
            <p className="text-sm font-normal">Transactions</p>
          </div>
        </div>

        {/* Pesonal Information like orders addres  */}
        <div className="flex flex-col gap-3">
          <p className="text-sm font-normal text-gray-600">YOUR INFORMATION</p>
          {/* Your Orders */}
          <div
            onClick={() => router.push("/orders")}
            className="nav-account-bar flex items-center justify-between "
          >
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
            Wallet AND TRANSACTIONS
          </p>
          {/* Your Wallet */}
          <div  onClick={() => router.push("/wallet")} className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div
               
                className="nav-account-item-icon bg-gray-200 p-2 rounded-full "
              >
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal">Wallet</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
         
          {/* Transactions */}
          <div
            onClick={() => router.push("/transactions")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <ScrollText size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal"> Transactions</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>
        {/* SECTION 3 */}
        <div className="flex flex-col gap-3">
          <p className="text-sm font-normal text-gray-600">OTHERS</p>
          {/* Your Delivery Partner */}
          <div
            onClick={() => router.push("/delivery-partner")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div
              
                className="nav-account-item-icon bg-gray-200 p-2 rounded-full "
              >
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-sm font-normal">Join as Delivery Partner</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>
      </div>
    </>
  );
};

export default Profile;
