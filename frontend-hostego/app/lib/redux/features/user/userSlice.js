import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  userAccount: {},
  cartData: { cart_items: [] },
  fetchCartData: false,
  useraddresses: [],
};

export const userSlice = createSlice({
  name: "user-slice",
  initialState,
  reducers: {
    setUserAccount: (state, { payload }) => {
      state.userAccount = payload;
    },
    setCartData: (state, { payload }) => {
      state.cartData = payload;
    },
    setFetchCartData: (state, { payload }) => {
      state.fetchCartData = payload;
    },
    setUserAddresses: (state, { payload }) => {
      state.useraddresses = payload;
    },
  },
});

// Action creators are generated for each case reducer function
export const { setUserAccount, setCartData, setFetchCartData, setUserAddresses } =
  userSlice.actions;

export default userSlice.reducer;
