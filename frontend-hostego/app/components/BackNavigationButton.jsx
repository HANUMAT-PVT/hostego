"use client"

import { ArrowLeft } from "lucide-react";
import { useRouter } from "next/navigation";

const BackNavigationButton = () => {

    const router = useRouter()
    return (
        <div onClick={()=>router.back()} className=" cursor-pointer p-3 flex gap-2 text-center items-center shadow-md h-fit ">
            <ArrowLeft size={22} className="" />
            <span className="font-md font-medium text-base">Profile</span>
        </div>
    );
};

export default BackNavigationButton;
