

export default function ShippingPolicy() {
  return (
    <>

        <title>Shipping & Delivery Policy | Hostego</title>
   


      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 py-10 text-gray-800">
        <h2 className="text-3xl font-semibold mb-2">Shipping & Delivery Policy</h2>
        <p className="text-sm text-gray-500 mb-6">
          Last updated on 11-04-2025
        </p>

        <p className="mb-4">
          At <strong>Hostego</strong>, we are committed to delivering food and essential goods from trusted
          local vendors directly to students living in hostels, right to their floors or rooms. 
          Our goal is to ensure fast, safe, and reliable delivery services that fit your daily needs.
        </p>

        <h3 className="text-xl font-semibold mt-6 mb-2">Delivery Areas</h3>
        <p className="mb-4">
          We currently operate within selected college and university hostel campuses. 
          Deliveries are made directly to the registered floor or room number provided during the order.
        </p>

        <h3 className="text-xl font-semibold mt-6 mb-2">Delivery Time</h3>
        <ul className="list-disc list-inside mb-4 space-y-2">
          <li>Delivery typically takes between <strong>20 to 45 minutes</strong>, depending on item availability and location.</li>
          <li>Peak hours or weather conditions may cause slight delays, but we’ll notify you if it does.</li>
          <li>Real-time order status and updates are available through your account dashboard.</li>
        </ul>

        <h3 className="text-xl font-semibold mt-6 mb-2">Shipping Charges</h3>
        <p className="mb-4">
          A small delivery fee may be applicable based on the distance or delivery urgency.
          Exact charges are shown at checkout before you confirm the order.
        </p>

        <h3 className="text-xl font-semibold mt-6 mb-2">Order Verification</h3>
        <p className="mb-4">
          All orders are verified before dispatch. In case the item is unavailable, you will be informed and refunded accordingly.
        </p>

        <h3 className="text-xl font-semibold mt-6 mb-2">Delivery Confirmation</h3>
        <p className="mb-4">
          Once your order is delivered, you will receive an in-app notification or SMS.
          Please ensure your contact details and room/floor number are accurate to avoid delivery issues.
        </p>

        <h3 className="text-xl font-semibold mt-6 mb-2">Failed Deliveries</h3>
        <ul className="list-disc list-inside mb-4 space-y-2">
          <li>If you are unavailable at the time of delivery, our delivery partner will attempt to contact you.</li>
          <li>If delivery is still unsuccessful, the order will be returned and no refund will be issued for perishable items.</li>
        </ul>

        <h3 className="text-xl font-semibold mt-6 mb-2">Razorpay Payment Gateway</h3>
        <p className="mb-4">
          We use <strong>Razorpay Payment Gateway</strong> to handle all online transactions securely. 
          Once payment is successful, the order is processed immediately for delivery. Refunds (if applicable) are processed 
          in 3–5 business days through the same payment method.
        </p>

        <h3 className="text-xl font-semibold mt-6 mb-2">Contact Us</h3>
        <p className="mb-2">
          For any questions regarding your delivery, feel free to contact our support team at:
        </p>
        <p className="mb-4">
          📧 support@hostego.in <br />
          📞 +91-8264121428
        </p>
      </main>

      {/* Footer */}
      <footer className="bg-gray-100 text-center py-4 text-sm text-gray-600 border-t">
        &copy; {new Date().getFullYear()} Hostego. All rights reserved.
      </footer>
    </>
  );
}
