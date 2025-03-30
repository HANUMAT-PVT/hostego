"use client"
import axios from "axios";
import { useEffect } from "react";
import axiosClient from "../utils/axiosClient"
import { useRouter } from "next/navigation"
import { useDispatch, useSelector } from "react-redux"
import { setFetchUserAccount } from "../lib/redux/features/user/userSlice"

const PhoneEmailAuth = () => {
    const router = useRouter()
    const dispatch = useDispatch()
    const { userAccount } = useSelector((state) => state.user)
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

    useEffect(() => {

        if (userAccount?.user_id) {
            router.push("/home")
        }
    }, [userAccount])

    // Function that gets triggered on successful authentication
    useEffect(() => {
        window.phoneEmailListener = (userObj) => {
            const userJsonUrl = userObj.user_json_url;


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
                first_name: data?.user_first_name || "Hostego",
                last_name: data?.user_last_name || "User",
                email: data?.user_email || `test${Math.random() * 2323 + Date.now()}@hostego.in`

            })
            console.log(response.data, "response")
            localStorage.setItem("auth-response", JSON.stringify(response.data))

            dispatch(setFetchUserAccount(true))
            router.push("/home")

            setTimeout(() => {
                window.location.reload()
            }, 1500)

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
