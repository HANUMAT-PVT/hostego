export const data = {
  order_id: "e3571ddc-3b87-45ae-a115-e793e4bf3e24",
  user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
  created_at: "2025-02-17 01:53:58.886302+05:30",
  updated_at: "2025-02-17 17:56:28.348765+05:30",
  order_items: [
    {
      user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
      quantity: 2,
      sub_total: 500,
      product_id: "f9fca8c1-dfb2-4ef0-b5de-9b388b3a3df3",
      cart_item_id: "e4c1e07f-b504-484c-a1cc-aa0ca246d647",
      product_item: {
        shop: {
          address: "123 MG Road, New Delhi, India",
          shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
          shop_img: "img_url",
          shop_name: "PUNJABI Delights",
          shop_status: 1,
          food_category: {
            is_veg: 0,
            is_cooked: 0,
          },
          preparation_time: "30 min",
        },
        tags: null,
        shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
        discount: {
          percentage: 10,
          is_available: 0,
        },
        created_at: "2025-02-13T18:48:33.882374+05:30",
        food_price: 250,
        product_id: "f9fca8c1-dfb2-4ef0-b5de-9b388b3a3df3",
        updated_at: "2025-02-13T18:48:33.882374+05:30",
        description:
          "A creamy and rich paneer dish cooked in butter with tomato-based gravy.",
        availability: 1,
        product_name: "Paneer Butter Masala",
        food_category: {
          is_veg: 0,
          is_cooked: 0,
        },
        product_img_url: "https://example.com/paneer.jpg",
        preparation_time: "30 min",
      },
    },
    {
      user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
      quantity: 7,
      sub_total: 1750,
      product_id: "e806a9c8-3136-4786-824c-1bc4e9b287e1",
      cart_item_id: "85cb1417-72d5-4b1d-bc1c-4edc9aae5e2c",
      product_item: {
        shop: {
          address: "123 MG Road, New Delhi, India",
          shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
          shop_img: "img_url",
          shop_name: "PUNJABI Delights",
          shop_status: 1,
          food_category: {
            is_veg: 0,
            is_cooked: 0,
          },
          preparation_time: "30 min",
        },
        tags: null,
        shop_id: "052dc329-3946-4b24-8055-19922d0aad6b",
        discount: {
          percentage: 10,
          is_available: 0,
        },
        created_at: "2025-02-13T18:58:38.914894+05:30",
        food_price: 250,
        product_id: "e806a9c8-3136-4786-824c-1bc4e9b287e1",
        updated_at: "2025-02-13T18:58:38.914894+05:30",
        description:
          "A creamy and rich paneer dish cooked in butter with tomato-based gravy.",
        availability: 1,
        product_name: "Paneer Butter Masala",
        food_category: {
          is_veg: 0,
          is_cooked: 0,
        },
        product_img_url: "https://example.com/paneer.jpg",
        preparation_time: "30 min",
      },
    },
  ],
  platform_fee: 1,
  shipping_fee: 30,
  final_order_value: 1781,
  delivery_partner_fee: 21,
  payment_transaction_id: "c4b66d17-d9a7-4a07-83a4-b5a135d5dff9",
  order_status: "pending",
  delivery_partner_id: "NULL",
  delivered_at: "0001-01-01 05:53:28+05:53:28",
  delivery_partner: {
    user: {
      email: "johndoe@example.com",
      user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
      last_name: "",
      created_at: "2025-02-17T00:39:39.655141+05:30",
      first_name: "John Doe",
      updated_at: "2025-02-17T00:39:39.655529+05:30",
      mobile_number: "+91-9876543211",
      last_login_timestamp: "2025-02-17T00:39:39.655141+05:30",
      firebase_otp_verified: 1,
    },
    address: "Zakir-A room number 1115",
    user_id: "a926b9cc-bc83-4c9b-a9bc-7ad985fc0d38",
    documents: {
      upi_id: "",
      aadhaar_back_img: "",
      bank_details_img: "",
      aadhaar_front_img: "",
    },
    account_status: 0,
    partner_img_url: "",
    availability_status: 0,
    delivery_partner_id: "5d8297ec-ad2a-443e-a9a1-a8e394107204",
  },
};

export const transformOrder = (order) => {
  const shopWiseOrders = {};

  order?.order_items?.forEach((item) => {
    const shop = item?.product_item?.shop;
    const shopId = shop?.shop_id;

    if (!shopWiseOrders[shopId]) {
      shopWiseOrders[shopId] = {
        shop_name: shop?.shop_name,
        shop_id: shopId,
        address: shop?.address,
        shop_img: shop?.shop_img,
        shop_status: shop?.shop_status,
        preparation_time: shop?.preparation_time,
        shop_products: [],
      };
    }

    shopWiseOrders[shopId].shop_products.push(item);
  });

  return {
    ...order,
    order_items: Object.values(shopWiseOrders), // Replace order_items with grouped shops
  };
};

export const transformDeliveryPartnerOrderEarnings = (orders) => {
  const earnings = orders.reduce((acc, order) => {
    return acc + order.final_order_value;
  }, 0);
  return earnings;
}



export const formatDate = (dateString) => {
  const date = new Date(dateString);
  const options = { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' };
  return date.toLocaleString('en-US', options);
};


export const transformOrdersByDate = (orders) => {
  return orders.reduce((acc, order) => {
    const orderDate = new Date(order.created_at).toISOString().split('T')[0]; // Extract YYYY-MM-DD
    
    if (!acc[orderDate]) {
      acc[orderDate] = [];
    }
    
    acc[orderDate].push(order);
    return acc;
  }, {});
};
