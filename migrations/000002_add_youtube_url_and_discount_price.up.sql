-- Add youtube_video_url column to products table
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS youtube_video_url TEXT;

-- Add discount_price column to discounts table
ALTER TABLE discounts 
ADD COLUMN IF NOT EXISTS discount_price DECIMAL(10,2);

-- Update existing discounts to have discount_price calculated from percentage
-- This will calculate discount_price for existing records based on the product price and percentage
UPDATE discounts d
SET discount_price = (
    SELECT p.price * (1 - d.percentage / 100.0)
    FROM products p
    WHERE p.id = d.product_id
)
WHERE discount_price IS NULL;

-- Make discount_price NOT NULL after backfilling
ALTER TABLE discounts 
ALTER COLUMN discount_price SET NOT NULL;

-- Make percentage column nullable (since it's now optional)
ALTER TABLE discounts 
ALTER COLUMN percentage DROP NOT NULL;
