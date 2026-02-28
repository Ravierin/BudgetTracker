-- Add margin column back (for rollback)
ALTER TABLE "position" ADD COLUMN margin DECIMAL(10, 2);
