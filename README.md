# Marketplace Backend

A comprehensive marketplace backend API built with Go and Fiber framework, implementing Clean Architecture principles.

## Features

- **User Management**: Registration, login, profile management
- **Authentication**: JWT-based authentication with middleware
- **Address Management**: User address CRUD operations
- **Category Management**: Product category management
- **Shop Management**: Shop creation and management
- **Product Management**: Product CRUD with photo support
- **Transaction System**: Order processing with invoice generation
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
   ```

4. **Setup database**
   - Create a MySQL database named `marketplace_backend`
   - The application will automatically create tables on first run

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

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

### Health Check
- `GET /health` - Server health check

## Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Database Schema

The application uses the following main entities:

- **Users**: User accounts with profile information
- **Shops**: User shops (automatically created on registration)
- **Categories**: Product categories
- **Products**: Products with photos and stock management
- **Addresses**: User delivery addresses
- **Transactions**: Orders with detailed line items

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
