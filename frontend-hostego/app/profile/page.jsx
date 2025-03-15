"use client";
import React, { useEffect, useState } from "react";
import {
  ChevronRight,
  CreditCard,
  NotebookTabs,
  Package,
  ScrollText,
  Wallet,
  MessageSquareText,
  LogOutIcon,
} from "lucide-react";
import BackNavigationButton from "../components/BackNavigationButton";
import { useRouter } from "next/navigation";
import axiosClient from "../utils/axiosClient";
import { useDispatch, useSelector } from "react-redux";
import { setUserAccount } from "../lib/redux/features/user/userSlice";
import { ChatIcon } from "lucide-react";
const Profile = () => {

  const router = useRouter();
  const dispatch = useDispatch();
  const { userAccount } = useSelector((state) => state.user)

  const handleUserLogout = () => {
    try {
      localStorage.removeItem("auth-response")
      router.push('/auth/sign-up')
      window.location.href = "/auth/sign-up"
    } catch (error) {
      console.log(error)
    }
  }
  const splitMobileNumberByCountryCode = (mobileNumber) => {
    if (!mobileNumber) return "+91 - XXXXXXXXXX";
    const countryCode = mobileNumber?.slice(0, 3);
    const mobileNumberWithoutCountryCode = mobileNumber?.slice(3);
    return `${countryCode} - ${mobileNumberWithoutCountryCode}`;
  }
  return (
    <>
      <BackNavigationButton title={"Profile"} />

      <div className="p-4 mt-3 flex flex-col gap-5">
        <div className="flex flex-col gap-1 font-semibold">
          <p className="text-2xl ">{((userAccount?.first_name||"Hostego") + " " + (userAccount?.last_name || "User" )).toUpperCase()} </p>
          <p className="text-sm  font-medium text-gray-500 flex items-center gap-2">{splitMobileNumberByCountryCode(userAccount?.mobile_number)} · {userAccount?.email}</p>
          <p onClick={() => router.push("/edit-account")} className="text-sm mt-1  text-[var(--primary-color)] flex items-center gap-1 cursor-pointer">Edit Profile <ChevronRight size={18} className="text-gray-600" /></p>
        </div>
        {/* Suggestion Box done */}

        <div className=" bg-[#eae8ff]  flex gap-4 justify-between py-4 px-8 rounded-lg">
          <div
            onClick={() => router.push("/wallet")}
            className="flex flex-col gap-2 items-center text-center cursor-pointer "
          >
            <Wallet size={20} />
            <p className="text-sm font-normal  ">Wallet</p>
          </div>
          <div
            onClick={() => router.push("/support")}
            className="flex flex-col gap-2 items-center text-center cursor-pointer "
          >
            <MessageSquareText size={20} />
            <p className="text-sm font-normal  ">Support</p>
          </div>
          <div
            onClick={() => router.push("/transactions")}
            className="flex flex-col gap-2 items-center text-center cursor-pointer "
          >
            <CreditCard size={20} />
            <p className="text-sm font-normal">Transactions</p>
          </div>
        </div>

        {/* Pesonal Information like orders addres  */}
        <div className="flex flex-col gap-3">
          <p className="text-md font-normal text-gray-600">YOUR INFORMATION</p>
          {/* Your Orders */}
          <div
            onClick={() => router.push("/orders")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Package size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Your orders</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Address */}
          <div onClick={() => router.push("/address")} className="nav-account-bar flex items-center justify-between ">
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <NotebookTabs size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Address</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>

        {/* SECTION 2 */}
        <div className="flex flex-col gap-3">
          <p className="text-md font-normal text-gray-600">
            Wallet AND TRANSACTIONS
          </p>
          {/* Your Wallet */}
          <div
            onClick={() => router.push("/wallet")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Wallet</p>
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
              <p className="text-md font-normal"> Transactions</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>
        {/* SECTION 3 */}
        <div className="flex flex-col gap-3">
          <p className="text-md font-normal text-gray-600">OTHERS</p>
          {/* Your Delivery Partner */}
          <div
            onClick={() => router.push("/delivery-partner")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Join as Delivery Partner</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* About us */}
          <div
            onClick={() => router.push("/about")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">About us</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Support us */}
          <div
            onClick={() => router.push("/support")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Support</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>

          {/* Terms & Conditions */}
          <div
            onClick={() => router.push("/terms-conditions")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Terms & Conditions</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Privacy Policy */}
          <div
            onClick={() => router.push("/privacy-policy")}
            className="nav-account-bar flex items-center justify-between "
          >
            <div className="flex items-center gap-3">
              <div className="nav-account-item-icon bg-gray-200 p-2 rounded-full ">
                <Wallet size={14} className="text-gray-500" />
              </div>
              <p className="text-md font-normal">Privacy Policy</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
          {/* Logout */}
          <div
            onClick={() => handleUserLogout()}
            className="nav-account-bar flex items-center justify-between  cursor-pointer"
          >
            <div className="flex items-center gap-3 ">
              <div className="nav-account-item-icon bg-red-500 p-2 rounded-full ">
                <LogOutIcon size={14} className="text-white" />
              </div>
              <p className="text-md font-normal">Logout</p>
            </div>
            <ChevronRight size={20} className="text-gray-400" />
          </div>
        </div>
      </div>
    </>
  );
};

export default Profile;
