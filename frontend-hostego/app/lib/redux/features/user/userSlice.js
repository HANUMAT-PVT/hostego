import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  userAccount: {},
  cartData: { cart_items: [] },
  fetchCartData: false,
  useraddresses: [],
  userRoles: [],
  fetchUser: true,
  userWallet: {},
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
    setUserRoles: (state, { payload }) => {
      state.userRoles = payload;
    },
    setFetchUserAccount: (state, { payload }) => {
      state.fetchUser = payload;
    },
    setUserAccountWallet: (state, { payload }) => {
      state.userWallet = payload;
    },
  },
});

// Action creators are generated for each case reducer function
export const { setUserAccount, setCartData, setFetchCartData, setUserAddresses, setUserRoles, setFetchUserAccount, setUserAccountWallet } =
  userSlice.actions;

export default userSlice.reducer;
