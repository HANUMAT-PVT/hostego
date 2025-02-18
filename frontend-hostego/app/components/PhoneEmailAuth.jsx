import { useEffect } from "react";

const PhoneEmailAuth = () => {
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

            alert(`Phone Verification Successful!\nFetch user data from:\n${userJsonUrl}`);
        };
    }, []);

    return (
        <div data-hide-name="true" className="pe_signin_button text-center rounded-lg bg-[#655df0] text-white text-center font-bold text-lg hover:bg-[#655df0] transition duration-200 ease-in-out shadow-md"
            data-client-id="13074759073192539427">
        </div>
    );
};

export default PhoneEmailAuth;
