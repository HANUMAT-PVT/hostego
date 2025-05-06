package dto

type AuthSignUpUserRequest struct {
	Gender              string `json:"gender""`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	MobileNumber        string `json:"mobile_number" binding:"required"`
	FirebaseOTPVerified int    `json:"firebase_otp_verified"`
	RefferalCode        string `json:"refferal_code"`
	IsVerifiedEmail     bool   `json:"is_verified_email"`
	DateOfBirth         string `json:"date_of_birth"`
	AdvertisementId     int    `json:"advertisement_id"`
	IsVerifiedMobile    bool   `json:"is_verified_mobile"`
	IsUserBlocked       bool   `json:"is_user_blocked"`
	IsWalletDisabled    bool   `json:"is_wallet_disabled"`
	IsDeliveryPerson    bool   `json:"is_delivery_person"`
}
