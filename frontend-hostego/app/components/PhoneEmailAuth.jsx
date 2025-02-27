"use client"
import axios from "axios";
import { useEffect } from "react";
import axiosClient from "../utils/axiosClient"

import { useRouter } from "next/navigation"

const PhoneEmailAuth = () => {
    const router = useRouter()

    useEffect(() => {
        // Dynamically load the Phone.Email script
        const script = document.createElement("script");
        script.src = "https://www.phone.email/sign_in_button_v1.js";
        script.async = true;
        document.body.appendChild(script);

        // Cleanup function
        return () => {
            document.body.removeChild(script);
        };
    }, []);

    // Function that gets triggered on successful authentication
    useEffect(() => {
        window.phoneEmailListener = (userObj) => {
            const userJsonUrl = userObj.user_json_url;

            console.log("Authenticated User JSON URL:", userJsonUrl);
            localStorage.setItem("userJsonUrl", userJsonUrl);

            handleUserAccountCreation(userJsonUrl)

            // alert(`Phone Verification Successful!\nFetch user data from:\n${userJsonUrl}`);
        };
    }, []);


    const handleUserAccountCreation = async (jsonUrl) => {

        try {
            let { data } = await axios.get(jsonUrl);

            let response = await axiosClient.post("/api/auth/signup", {
                mobile_number: data?.user_country_code + data?.user_phone_number,
                first_name: "Test",
                last_name: "User",
                email: "test@test.com"

            })

            localStorage.setItem("auth-response", JSON.stringify(response.data))
            // alert("Signup successfull")
            router.push("/home")

        } catch (error) {
            console.log(error, "error")
        }
    }

    return (
        <div data-hide-name="true" className="pe_signin_button flex items-center text-center rounded-lg bg-[#655df0] text-white text-center font-bold text-lg hover:bg-[#655df0] transition duration-200 ease-in-out shadow-md"
            data-client-id="13074759073192539427">
        </div>
    );
};

export default PhoneEmailAuth;
