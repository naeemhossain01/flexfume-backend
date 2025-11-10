-- Revert percentage column to NOT NULL
ALTER TABLE discounts 
ALTER COLUMN percentage SET NOT NULL;

-- Remove discount_price column from discounts table
ALTER TABLE discounts 
DROP COLUMN IF EXISTS discount_price;

-- Remove youtube_video_url column from products table
ALTER TABLE products 
DROP COLUMN IF EXISTS youtube_video_url;
