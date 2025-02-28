"use client"

import { ArrowLeft } from "lucide-react";
import { useRouter } from "next/navigation";

const BackNavigationButton = ({title}) => {

    const router = useRouter()
    return (
        <div onClick={()=>router.back()} className="sticky top-0 z-10 bg-white cursor-pointer p-3 flex gap-2 text-center items-center shadow-md h-fit ">
            <ArrowLeft size={22} className="" />
            <span className="font-md font-medium text-lg">{title}</span>
        </div>
    );
};

export default BackNavigationButton;
