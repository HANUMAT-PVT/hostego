import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  userAccount: {},
  cartData: {cart_items: []},
  fetchCartData: false,
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
  },
});

// Action creators are generated for each case reducer function
export const { setUserAccount, setCartData,setFetchCartData } = userSlice.actions;

export default userSlice.reducer;
