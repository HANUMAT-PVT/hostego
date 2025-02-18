
import { initializeApp } from "firebase/app";
import { getAuth, RecaptchaVerifier, PhoneAuthProvider } from "firebase/auth";


const firebaseConfig = {
  apiKey: "AIzaSyDjNkcSEk9Q_UU9YOh43A0YTPa11a1WI_c",

  authDomain: "hostego-3eccf.firebaseapp.com",

  projectId: "hostego-3eccf",

  storageBucket: "hostego-3eccf.firebasestorage.app",

  messagingSenderId: "650160587814",

  appId: "1:650160587814:web:6feb6435986e4a123452e7",

  measurementId: "G-CLC0FVZEFX",
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

export { auth, RecaptchaVerifier, PhoneAuthProvider };
