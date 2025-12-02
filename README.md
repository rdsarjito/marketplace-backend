# Warung Budeh Ramah Backend

Backend API untuk aplikasi **Warung Budeh Ramah**, dibangun dengan Go dan Fiber framework, menerapkan prinsip Clean Architecture.

## Features

- **User Management**: Registration, login, profile management
- **Authentication**: JWT-based authentication with middleware
- **Address Management**: User address CRUD operations
- **Category Management**: Product category management
- **Shop Management**: Shop creation and management
- **Product Management**: Product CRUD with photo support
- **Transaction System**: Order processing with invoice generation
- **Payment Gateway Integration**: Midtrans payment gateway integration (Virtual Account, E-Wallet, Bank Transfer, Credit Card, COD)
- **Payment Status Tracking**: Real-time payment status updates via webhook
- **Email Notifications**: Payment success and expiration notifications
- **External API Integration**: Integration with API Wilayah Indonesia for province/city data

## Tech Stack

- **Backend**: Go with Fiber framework
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT (JSON Web Token)
- **Validation**: Go Playground Validator
- **Password Hashing**: bcrypt
- **External API**: API Wilayah Indonesia

## Project Structure

```
marketplace-backend/
├── config/                 # Configuration files
├── constants/              # Error and success messages
├── domain/                 # Domain layer
│   ├── dto/               # Data Transfer Objects
│   │   ├── request/       # Request DTOs
│   │   └── response/      # Response DTOs
│   └── model/             # Database models
├── handlers/              # HTTP handlers
├── helpers/               # Helper functions
├── middleware/            # Custom middleware
├── repositories/          # Data access layer
├── services/              # Business logic layer
├── utils/                 # Utility functions
├── main.go               # Application entry point
└── README.md             # This file
```

## Installation

### Prerequisites

- Go 1.19 or higher
- MySQL 8.0 or higher
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/rdsarjito/marketplace-backend.git
   cd marketplace-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup environment variables**
   Create a `.env` file in the root directory:
   ```env
   # APP CONFIG
   SECRET_KEY=your-secret-key-here-change-this-in-production
   APP_HOST=localhost
   APP_PORT=8080

   # DB CONFIG
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your-password
   DB_NAME=marketplace_backend

   # Integration
   API_LOCATION=https://emsifa.github.io/api-wilayah-indonesia/api

   # Storage (MinIO)
   MINIO_ENDPOINT=http://localhost:9000
   MINIO_ACCESS_KEY=admin
   MINIO_SECRET_KEY=admin123
   MINIO_BUCKET_NAME=product-media
   ASSET_BASE_URL=http://localhost:9000/product-media
   MINIO_USE_SSL=false

   # Payment Gateway (Midtrans)
   MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxxxxxxxxxx
   MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxxxxxxxxxx
   MIDTRANS_IS_PRODUCTION=false

   # Frontend URL (for payment redirect)
   FRONTEND_URL=http://localhost:5173
   ```

4. **Setup database**
   - Create a MySQL database named `marketplace_backend`
   - The application will automatically create tables on first run
   - For payment features, run the migration script (optional, GORM AutoMigrate will handle it):
     ```bash
     # Migration is handled automatically by GORM AutoMigrate
     # Manual migration script available at: migrations/001_add_payment_fields_to_trx.sql
     ```

5. **Setup Midtrans Payment Gateway** (Optional)
   - Register at [Midtrans Dashboard](https://dashboard.midtrans.com/)
   - Get your **Server Key** and **Client Key** from Settings → Access Keys
   - Add the keys to your `.env` file (see above)
   - For testing, use Sandbox keys (set `MIDTRANS_IS_PRODUCTION=false`)
   - See [PAYMENT_TESTING.md](./PAYMENT_TESTING.md) for detailed testing guide

6. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

### MinIO Setup (Local)

1. **Start MinIO via Docker**
   ```bash
   docker run -d --name minio \
     -p 9000:9000 -p 9090:9090 \
     -e MINIO_ROOT_USER=admin \
     -e MINIO_ROOT_PASSWORD=admin123 \
     -v ~/minio-data:/data \
     minio/minio server /data --console-address ":9090"
   ```
2. **Create bucket**
   - Open `http://localhost:9090`, log in with the credentials above.
   - Create a bucket named `product-media` (or match `MINIO_BUCKET_NAME`).
3. **Update `.env`**
   - Ensure the values under the *Storage (MinIO)* block match the server/bucket.
   - `ASSET_BASE_URL` should point to the public path of the bucket (path-style by default).

For production deployments, front MinIO with HTTPS (reverse proxy or load balancer) and rotate `MINIO_ACCESS_KEY`/`MINIO_SECRET_KEY` to strong secrets. Enable versioning & regular backups to protect media assets.

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### User Management
- `GET /api/v1/user` - Get user profile
- `PUT /api/v1/user` - Update user profile
- `GET /api/v1/user/alamat` - Get user addresses
- `GET /api/v1/user/alamat/:id` - Get address detail
- `POST /api/v1/user/alamat` - Create address
- `PUT /api/v1/user/alamat/:id` - Update address
- `DELETE /api/v1/user/alamat/:id` - Delete address

### Province & City (Public)
- `GET /api/v1/provcity/listprovincies` - Get provinces list
- `GET /api/v1/provcity/detailprovince/:prov_id` - Get province detail
- `GET /api/v1/provcity/listcities/:prov_id` - Get cities by province
- `GET /api/v1/provcity/detailcity/:city_id` - Get city detail

### Category Management
- `GET /api/v1/category` - Get categories list
- `GET /api/v1/category/:id` - Get category detail
- `POST /api/v1/category` - Create category
- `PUT /api/v1/category/:id` - Update category
- `DELETE /api/v1/category/:id` - Delete category

### Shop Management
- `GET /api/v1/toko/my` - Get my shop
- `GET /api/v1/toko` - Get shops list
- `GET /api/v1/toko/:id_toko` - Get shop detail
- `PUT /api/v1/toko/:id_toko` - Update shop profile

### Product Management
- `GET /api/v1/product` - Get products list
- `GET /api/v1/product/:id` - Get product detail
- `POST /api/v1/product` - Create product
- `PUT /api/v1/product/:id` - Update product
- `DELETE /api/v1/product/:id` - Delete product

### Transaction Management
- `GET /api/v1/trx` - Get transactions list
- `GET /api/v1/trx/:id` - Get transaction detail
- `POST /api/v1/trx` - Create transaction
- `POST /api/v1/trx/:id/check-payment` - Check payment status manually

### Payment Gateway
- `POST /api/v1/payment/webhook` - Midtrans payment webhook endpoint (public)

### Health Check
- `GET /health` - Server health check

## Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Payment Gateway Integration

This application integrates with **Midtrans** payment gateway to support multiple payment methods:

### Supported Payment Methods

- **COD (Cash on Delivery)**: Direct payment on delivery
- **Virtual Account**: Bank transfer via Virtual Account (BCA, BNI, Mandiri)
- **E-Wallet**: GoPay, OVO, DANA, LinkAja
- **Bank Transfer**: Direct bank transfer (BCA, BNI, Mandiri)
- **Credit Card**: Credit card payment with 3DS support

### Payment Flow

1. **Create Transaction**: User creates transaction with selected payment method
2. **Payment Creation**: For non-COD methods, payment is created via Midtrans API
3. **Payment URL**: User is redirected to Midtrans payment page
4. **Payment Processing**: User completes payment on Midtrans
5. **Webhook Notification**: Midtrans sends webhook to update payment status
6. **Status Update**: Transaction status is updated automatically

### Payment Status

- `pending_payment`: Payment is pending
- `paid`: Payment completed successfully
- `expired`: Payment expired
- `failed`: Payment failed
- `cancelled`: Payment cancelled

### Setup & Testing

For detailed setup instructions and testing guide, see [PAYMENT_TESTING.md](./PAYMENT_TESTING.md)

## Database Schema

The application uses the following main entities:

- **Users**: User accounts with profile information
- **Shops**: User shops (automatically created on registration)
- **Categories**: Product categories
- **Products**: Products with photos and stock management
- **Addresses**: User delivery addresses
- **Transactions**: Orders with detailed line items and payment information

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o marketplace-backend main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
