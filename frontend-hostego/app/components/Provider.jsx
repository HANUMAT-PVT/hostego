"use client"; // ✅ Ensure this is a Client Component

import { Provider } from "react-redux";
import { store } from "../lib/redux/store";

const Providers = ({ children }) => {
  return <Provider store={store}>{children}</Provider>;
};

export default Providers;
