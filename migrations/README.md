# Database Migrations

This directory contains SQL migration scripts for the marketplace backend database.

## Migration Files

### 001_add_payment_fields_to_trx.sql
Adds payment gateway fields to the `trx` (transaction) table to support Midtrans payment gateway integration.

**Fields added:**
- `payment_status` (VARCHAR(50), default: 'pending_payment')
- `payment_token` (VARCHAR(255), nullable)
- `payment_url` (TEXT, nullable)
- `midtrans_order_id` (VARCHAR(255), nullable, indexed)
- `payment_expired_at` (TIMESTAMP, nullable)

## Running Migrations

### Option 1: Using GORM AutoMigrate (Development)
The application automatically runs migrations using GORM AutoMigrate when it starts. This is suitable for development environments.

### Option 2: Manual SQL Migration (Production)
For production environments, it's recommended to run migrations manually:

```bash
mysql -u [username] -p [database_name] < migrations/001_add_payment_fields_to_trx.sql
```

Or using MySQL client:
```sql
SOURCE migrations/001_add_payment_fields_to_trx.sql;
```

## Rollback

To rollback this migration, you can run:

```sql
ALTER TABLE trx DROP INDEX idx_midtrans_order_id;
ALTER TABLE trx DROP COLUMN payment_expired_at;
ALTER TABLE trx DROP COLUMN midtrans_order_id;
ALTER TABLE trx DROP COLUMN payment_url;
ALTER TABLE trx DROP COLUMN payment_token;
ALTER TABLE trx DROP COLUMN payment_status;
```

## Notes

- GORM AutoMigrate will automatically add these columns if they don't exist
- For existing production databases, run the SQL migration script first
- The migration is idempotent - running it multiple times won't cause errors (except for index creation)


