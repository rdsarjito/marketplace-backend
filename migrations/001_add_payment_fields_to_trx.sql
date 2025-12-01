-- Migration: Add payment gateway fields to trx table
-- Date: 2024
-- Description: Add payment_status, payment_token, payment_url, midtrans_order_id, payment_expired_at columns to support Midtrans payment gateway integration

-- Add payment_status column
ALTER TABLE trx 
ADD COLUMN payment_status VARCHAR(50) DEFAULT 'pending_payment' AFTER method_bayar;

-- Add payment_token column (nullable)
ALTER TABLE trx 
ADD COLUMN payment_token VARCHAR(255) NULL AFTER payment_status;

-- Add payment_url column (nullable)
ALTER TABLE trx 
ADD COLUMN payment_url TEXT NULL AFTER payment_token;

-- Add midtrans_order_id column (nullable) with index
ALTER TABLE trx 
ADD COLUMN midtrans_order_id VARCHAR(255) NULL AFTER payment_url;

-- Add index for midtrans_order_id for faster queries
CREATE INDEX idx_midtrans_order_id ON trx(midtrans_order_id);

-- Add payment_expired_at column (nullable)
ALTER TABLE trx 
ADD COLUMN payment_expired_at TIMESTAMP NULL AFTER midtrans_order_id;

-- Add comment to columns for documentation
ALTER TABLE trx 
MODIFY COLUMN payment_status VARCHAR(50) DEFAULT 'pending_payment' COMMENT 'Payment status: pending_payment, paid, expired, failed, cancelled';

ALTER TABLE trx 
MODIFY COLUMN payment_token VARCHAR(255) NULL COMMENT 'Payment token from Midtrans';

ALTER TABLE trx 
MODIFY COLUMN payment_url TEXT NULL COMMENT 'Redirect URL for Midtrans payment page';

ALTER TABLE trx 
MODIFY COLUMN midtrans_order_id VARCHAR(255) NULL COMMENT 'Order ID from Midtrans (invoice code)';

ALTER TABLE trx 
MODIFY COLUMN payment_expired_at TIMESTAMP NULL COMMENT 'Payment expiration timestamp';


