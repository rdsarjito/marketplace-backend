package constants

// Payment status constants
const (
	PaymentStatusPendingPayment = "pending_payment"
	PaymentStatusPaid           = "paid"
	PaymentStatusExpired        = "expired"
	PaymentStatusFailed         = "failed"
	PaymentStatusCancelled      = "cancelled"
)

// Payment method constants
const (
	PaymentMethodCOD            = "COD"
	PaymentMethodVirtualAccount = "virtual_account"
	PaymentMethodEWallet        = "e_wallet"
	PaymentMethodBankTransfer   = "bank_transfer"
	PaymentMethodCreditCard     = "credit_card"
)

// Payment method aliases (for validation and mapping)
var (
	PaymentMethodAliases = map[string]string{
		"cod":              PaymentMethodCOD,
		"va":               PaymentMethodVirtualAccount,
		"virtual_account":  PaymentMethodVirtualAccount,
		"ewallet":          PaymentMethodEWallet,
		"e_wallet":         PaymentMethodEWallet,
		"gopay":            PaymentMethodEWallet,
		"ovo":              PaymentMethodEWallet,
		"dana":             PaymentMethodEWallet,
		"linkaja":          PaymentMethodEWallet,
		"bank_transfer":    PaymentMethodBankTransfer,
		"bank_transfer_bca": PaymentMethodBankTransfer,
		"bank_transfer_bni": PaymentMethodBankTransfer,
		"bank_transfer_mandiri": PaymentMethodBankTransfer,
		"cc":               PaymentMethodCreditCard,
		"credit_card":      PaymentMethodCreditCard,
	}
)


