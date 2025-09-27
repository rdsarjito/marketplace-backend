package constants

const (
	// Authentication errors
	ErrInvalidCredentials = "Invalid email or password"
	ErrUserNotFound       = "User not found"
	ErrUserAlreadyExists  = "User already exists"
	ErrInvalidToken       = "Invalid or expired token"
	ErrUnauthorized       = "Unauthorized access"
	ErrForbidden          = "Forbidden access"

	// Validation errors
	ErrInvalidInput       = "Invalid input data"
	ErrRequiredField      = "Required field is missing"
	ErrInvalidEmail       = "Invalid email format"
	ErrInvalidPhone       = "Invalid phone number format"
	ErrPasswordTooShort   = "Password must be at least 6 characters"

	// Database errors
	ErrDatabaseConnection = "Database connection failed"
	ErrRecordNotFound     = "Record not found"
	ErrDuplicateEntry     = "Duplicate entry"

	// Business logic errors
	ErrInsufficientStock  = "Insufficient stock"
	ErrInvalidQuantity    = "Invalid quantity"
	ErrShopNotFound       = "Shop not found"
	ErrProductNotFound    = "Product not found"
	ErrCategoryNotFound   = "Category not found"
	ErrAddressNotFound    = "Address not found"

	// External API errors
	ErrExternalAPI        = "External API error"
	ErrProvinceNotFound   = "Province not found"
	ErrCityNotFound       = "City not found"
)
