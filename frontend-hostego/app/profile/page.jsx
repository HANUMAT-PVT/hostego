"use client";
import React, { useEffect, useState } from "react";
import {
  ChevronRight,
  CreditCard,
  NotebookTabs,
  Package,
  ScrollText,
  Wallet,
  LogOutIcon
} from "lucide-react";
import BackNavigationButton from "../components/BackNavigationButton";
import { useRouter } from "next/navigation";
import axiosClient from "../utils/axiosClient";

const Profile = () => {
  const [userAccount, setUserAccount] = useState({});
  const router = useRouter();

  useEffect(() => {
    fetchUserAccount();
  }, []);

  const fetchUserAccount = async () => {
    try {
      const { data } = await axiosClient.get("/api/user/me");

      setUserAccount(data)
    } catch (error) {
      console.log(error);
    }
  };

  const handleUserLogout = () => {
    try {
      localStorage.removeItem("auth-response")
      router.push('/auth/sign-up')
      window.location.href = "/auth/sign-up"
    } catch (error) {
      console.log(error)
    }
  }

  return (
    <>
      <BackNavigationButton title={"Profile"} />

      <div className="p-4 mt-3 flex flex-col gap-5">
        <div className="flex flex-col gap-3">
          <p className="text-2xl font-medium">My account </p>
          <p className="text-md font-normal">{userAccount?.mobile_number}</p>
        </div>
        {/* Suggestion Box done */}

        <div className=" bg-[#eae8ff]  flex gap-4 justify-between py-4 px-8 rounded-lg">
          <div
            onClick={() => router.push("/wallet")}
            className="flex flex-col gap-2 items-center text-center "
          >
            <Wallet size={20} />
            <p className="text-md font-normal  ">Wallet</p>
          </div>
          <div
            onClick={() => router.push("/transactions")}
            className="flex flex-col gap-2 items-center text-center "
          >
            <CreditCard size={20} />
            <p className="text-md font-normal">Transactions</p>
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
            onClick={() => router.push("/delivery-partner")}
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
