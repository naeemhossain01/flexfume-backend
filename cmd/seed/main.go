package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
)

func main() {
	log.Println("Starting database seeding...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed product categories
	if err := seedProductCategories(); err != nil {
		log.Fatalf("Failed to seed product categories: %v", err)
	}

	// Seed products
	if err := seedProducts(); err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	// Seed discounts
	if err := seedDiscounts(); err != nil {
		log.Fatalf("Failed to seed discounts: %v", err)
	}

	// Seed delivery costs
	if err := seedDeliveryCosts(); err != nil {
		log.Fatalf("Failed to seed delivery costs: %v", err)
	}

	// Seed admin user
	if err := seedAdminUser(); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	// Seed regular user
	if err := seedRegularUser(); err != nil {
		log.Fatalf("Failed to seed regular user: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

func seedProductCategories() error {
	log.Println("Seeding product categories...")

	categories := []models.Category{
		{
			ID:          "74facc88-8de8-4a1d-94e2-694f730d9113",
			Name:        "Man",
			Description: "Men's fragrances",
			CreatedAt:   time.Date(2025, 10, 14, 16, 41, 7, 145395000, time.UTC),
			UpdatedAt:   time.Date(2025, 10, 14, 16, 41, 7, 145395000, time.UTC),
		},
		{
			ID:          "7b9a6f96-59ca-4f62-a48f-0f80dd079b5c",
			Name:        "Women",
			Description: "Women's fragrances",
			CreatedAt:   time.Date(2025, 10, 14, 16, 41, 7, 317154000, time.UTC),
			UpdatedAt:   time.Date(2025, 10, 14, 16, 41, 7, 317154000, time.UTC),
		},
		{
			ID:          "a73bcb25-561a-4fe3-ba1d-8a12f9c00663",
			Name:        "Man & Women",
			Description: "Unisex fragrances",
			CreatedAt:   time.Date(2025, 10, 14, 16, 41, 7, 428874000, time.UTC),
			UpdatedAt:   time.Date(2025, 10, 14, 16, 41, 7, 428874000, time.UTC),
		},
	}

	for _, category := range categories {
		// Check if category already exists
		var existingCategory models.Category
		result := database.DB.Where("id = ?", category.ID).First(&existingCategory)
		
		if result.Error == nil {
			log.Printf("Category '%s' already exists, skipping...", category.Name)
			continue
		}

		// Create category
		if err := database.DB.Create(&category).Error; err != nil {
			return err
		}

		log.Printf("✓ Category created successfully (ID: %s)", category.ID)
		log.Printf("  Name: %s", category.Name)
		log.Printf("  Description: %s", category.Description)
	}

	log.Println("Product categories seeding completed!")
	return nil
}

func seedProducts() error {
	log.Println("Seeding products...")

	products := []models.Product{
		{
			ID:          "5ae84616-6aa9-4654-a1d6-175291132400",
			ProductName: "Beauty",
			ProductCode: "beauty",
			Description: "An elegant and feminine fragrance that captures the essence of modern beauty. Beauty combines floral and fruity notes to create a delicate yet memorable scent that enhances your natural charm.",
			LongDescription: "Beauty is a sophisticated feminine fragrance designed for the modern woman who appreciates elegance and quality. This exquisite scent features a harmonious blend of fresh floral top notes with warm vanilla and musk base notes, creating a timeless and versatile fragrance. The carefully balanced composition ensures all-day wearability while maintaining its distinctive character. The fragrance opens with delicate rose and jasmine petals, complemented by sweet peach and bergamot, then evolves through creamy gardenia and lily of the valley, finally settling into a luxurious base of vanilla, sandalwood, and white musk. This carefully crafted blend provides 8-12 hours of elegant fragrance projection, perfect for both professional environments and romantic evenings.",
			ImageURL:    "/images/perfume-beauty.jpg",
			Price:       1490.00,
			Stock:       1000,
			CategoryID:  "7b9a6f96-59ca-4f62-a48f-0f80dd079b5c", // Women category
			KeyFeatures: pq.StringArray{
				"Long-lasting fragrance (8-12 hours)",
				"Premium quality ingredients sourced globally",
				"Two perfume containers included for convenience",
				"Compact and portable design (15ml each)",
				"Professional-grade packaging with gift box",
				"Suitable for all occasions and seasons",
				"Made with natural essential oils and synthetic compounds",
				"Cruelty-free and vegan-friendly formulation",
				"Alcohol-based for optimal longevity",
				"Feminine appeal with universal charm",
				"Made in Bangladesh with international standards",
				"Comes with detailed usage instructions",
			},
			WhyChooseBenefits: pq.StringArray{
				"Exceptional value with two containers (30ml total)",
				"Professional-grade quality at affordable price point",
				"Perfect for gifting or personal use",
				"Compact design fits in any pocket or travel bag",
				"Long-lasting scent that consistently receives compliments",
				"Suitable for both casual and formal occasions",
				"Made by experienced perfumers with 20+ years expertise",
				"Backed by our 100% satisfaction guarantee",
				"Free shipping on orders above ৳500",
				"Easy returns within 7 days",
				"Compatible with all skin types",
				"No harmful chemicals or allergens",
			},
			CreatedAt: time.Date(2025, 10, 14, 16, 41, 7, 59266000, time.UTC),
			UpdatedAt: time.Date(2025, 10, 14, 19, 26, 34, 662646000, time.UTC),
		},
		{
			ID:          "6dd2e13f-a8fb-48ca-91c6-152f0b3846b5",
			ProductName: "Beast",
			ProductCode: "beast",
			Description: "A bold and masculine fragrance designed for the modern man. Beast embodies strength, confidence, and sophistication with its rich, woody notes that create an unforgettable presence.",
			LongDescription: "Beast is a premium masculine fragrance crafted for the contemporary gentleman who values quality and distinction. This sophisticated scent combines deep woody undertones with subtle citrus top notes, creating a complex and alluring aroma that evolves throughout the day. Perfect for professional settings and evening occasions, Beast delivers long-lasting performance with exceptional sillage. The fragrance opens with fresh bergamot and lemon, transitions through spicy cardamom and lavender, and settles into a rich base of sandalwood, amber, and musk. Each application provides 8-12 hours of consistent fragrance projection, making it ideal for busy professionals and social occasions alike.",
			ImageURL:    "/images/perfume-beast.jpg",
			Price:       1490.00,
			Stock:       1000,
			CategoryID:  "74facc88-8de8-4a1d-94e2-694f730d9113", // Man category
			KeyFeatures: pq.StringArray{
				"Long-lasting fragrance (8-12 hours)",
				"Premium quality ingredients sourced globally",
				"Two perfume containers included for convenience",
				"Compact and portable design (15ml each)",
				"Professional-grade packaging with gift box",
				"Suitable for all occasions and seasons",
				"Made with natural essential oils and synthetic compounds",
				"Cruelty-free and vegan-friendly formulation",
				"Alcohol-based for optimal longevity",
				"Unisex appeal with masculine dominance",
				"Made in Bangladesh with international standards",
				"Comes with detailed usage instructions",
			},
			WhyChooseBenefits: pq.StringArray{
				"Exceptional value with two containers (30ml total)",
				"Professional-grade quality at affordable price point",
				"Perfect for gifting or personal use",
				"Compact design fits in any pocket or travel bag",
				"Long-lasting scent that consistently receives compliments",
				"Suitable for both casual and formal occasions",
				"Made by experienced perfumers with 20+ years expertise",
				"Backed by our 100% satisfaction guarantee",
				"Free shipping on orders above ৳500",
				"Easy returns within 7 days",
				"Compatible with all skin types",
				"No harmful chemicals or allergens",
			},
			CreatedAt: time.Date(2025, 10, 14, 16, 41, 8, 60659000, time.UTC),
			UpdatedAt: time.Date(2025, 10, 14, 16, 41, 8, 60659000, time.UTC),
		},
		{
			ID:          "ca488421-d1dc-4d73-8b50-97623a76dd99",
			ProductName: "Beauty & Beast",
			ProductCode: "beauty-beast",
			Description: "The perfect combination of masculine and feminine fragrances in one elegant package. Beauty & Beast offers both scents, allowing couples to share the experience or individuals to enjoy both sides of their personality.",
			LongDescription: "Beauty & Beast is our signature dual-fragrance collection that celebrates the harmony between masculine and feminine energies. This unique set includes both our acclaimed Beauty and Beast fragrances, each crafted with the same attention to detail and quality. Perfect for couples who want to share their fragrance journey or individuals who appreciate both sides of their personality. The complementary scents work beautifully together, creating a complete fragrance experience. The Beauty fragrance features delicate floral notes of rose, jasmine, and gardenia with warm vanilla and musk, while the Beast fragrance combines fresh citrus with woody undertones of sandalwood and amber. Together, they create a harmonious balance that appeals to all preferences and occasions.",
			ImageURL:    "/images/perfume-beauty-beast.jpg",
			Price:       1490.00,
			Stock:       1000,
			CategoryID:  "a73bcb25-561a-4fe3-ba1d-8a12f9c00663", // Man & Women category
			KeyFeatures: pq.StringArray{
				"Two complete fragrances in one package",
				"Long-lasting fragrance (8-12 hours each)",
				"Premium quality ingredients sourced globally",
				"Four perfume containers total (2 for each scent)",
				"Compact and portable design (15ml each)",
				"Professional-grade packaging with gift box",
				"Suitable for all occasions and seasons",
				"Made with natural essential oils and synthetic compounds",
				"Cruelty-free and vegan-friendly formulation",
				"Alcohol-based for optimal longevity",
				"Perfect for couples or individuals",
				"Made in Bangladesh with international standards",
				"Comes with detailed usage instructions",
				"Compatible with all skin types",
			},
			WhyChooseBenefits: pq.StringArray{
				"Best value with four containers total (60ml)",
				"Perfect for couples or individuals who love variety",
				"Professional-grade quality at affordable price point",
				"Ideal for gifting or personal use",
				"Compact design fits in any pocket or travel bag",
				"Long-lasting scents that consistently receive compliments",
				"Suitable for both casual and formal occasions",
				"Made by experienced perfumers with 20+ years expertise",
				"Backed by our 100% satisfaction guarantee",
				"Free shipping on orders above ৳500",
				"Easy returns within 7 days",
				"Complete fragrance wardrobe in one package",
				"No harmful chemicals or allergens",
			},
			CreatedAt: time.Date(2025, 10, 14, 16, 41, 8, 170464000, time.UTC),
			UpdatedAt: time.Date(2025, 10, 14, 16, 41, 8, 170464000, time.UTC),
		},
	}

	for _, product := range products {
		// Check if product already exists
		var existingProduct models.Product
		result := database.DB.Where("id = ?", product.ID).First(&existingProduct)
		
		if result.Error == nil {
			// Product exists, update it with all fields including longDescription
			log.Printf("Product '%s' already exists, updating...", product.ProductName)
			if err := database.DB.Model(&existingProduct).Updates(map[string]interface{}{
				"product_name":         product.ProductName,
				"product_code":         product.ProductCode,
				"description":          product.Description,
				"long_description":     product.LongDescription,
				"image_url":            product.ImageURL,
				"price":                product.Price,
				"stock":                product.Stock,
				"category_id":          product.CategoryID,
				"key_features":         product.KeyFeatures,
				"why_choose_benefits":  product.WhyChooseBenefits,
				"updated_at":           time.Now(),
			}).Error; err != nil {
				return err
			}
			continue
		}

		// Create product
		if err := database.DB.Create(&product).Error; err != nil {
			return err
		}

		log.Printf("✓ Product created successfully (ID: %s)", product.ID)
		log.Printf("  Name: %s", product.ProductName)
		log.Printf("  Code: %s", product.ProductCode)
		log.Printf("  Price: %.2f", product.Price)
		log.Printf("  Stock: %d", product.Stock)
		log.Printf("  Category ID: %s", product.CategoryID)
	}

	log.Println("Products seeding completed!")
	return nil
}

func seedDiscounts() error {
	log.Println("Seeding discounts...")

	// Product IDs from our seeded products
	productIDs := []string{
		"5ae84616-6aa9-4654-a1d6-175291132400", // Beauty
		"6dd2e13f-a8fb-48ca-91c6-152f0b3846b5", // Beast
		"ca488421-d1dc-4d73-8b50-97623a76dd99", // Beauty & Beast
	}

	for _, productID := range productIDs {
		// Check if discount already exists for this product
		var existingDiscount models.Discount
		result := database.DB.Where("product_id = ?", productID).First(&existingDiscount)
		
		if result.Error == nil {
			log.Printf("Discount for product '%s' already exists, skipping...", productID)
			continue
		}

		// Create discount with 33% off
		discount := models.Discount{
			ID:         uuid.New().String(),
			ProductID:  productID,
			Percentage: 33,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := database.DB.Create(&discount).Error; err != nil {
			return err
		}

		log.Printf("✓ Discount created successfully (ID: %s)", discount.ID)
		log.Printf("  Product ID: %s", discount.ProductID)
		log.Printf("  Percentage: %d%%", discount.Percentage)
	}

	log.Println("Discounts seeding completed!")
	return nil
}

func seedDeliveryCosts() error {
	log.Println("Seeding delivery costs...")

	deliveryCosts := []models.DeliveryCost{
		{
			ID:        "4dc3c0c2-330a-11f0-a05a-98e7f4669aeb",
			Location:  "Dhaka City",
			Service:   "Standard Delivery",
			Cost:      60.00,
			CreatedAt: time.Date(2025, 10, 14, 16, 41, 9, 412342000, time.UTC),
			UpdatedAt: time.Date(2025, 10, 14, 16, 41, 9, 412342000, time.UTC),
		},
		{
			ID:        "7d3aba03-c247-4039-b7e5-bd0f0fcbfb00",
			Location:  "Other Cities",
			Service:   "Standard Delivery",
			Cost:      150.00,
			CreatedAt: time.Date(2025, 10, 14, 16, 41, 9, 842228000, time.UTC),
			UpdatedAt: time.Date(2025, 10, 14, 16, 41, 9, 842228000, time.UTC),
		},
	}

	for _, deliveryCost := range deliveryCosts {
		// Check if delivery cost already exists
		var existingDeliveryCost models.DeliveryCost
		result := database.DB.Where("id = ?", deliveryCost.ID).First(&existingDeliveryCost)
		
		if result.Error == nil {
			log.Printf("Delivery cost for '%s' already exists, skipping...", deliveryCost.Location)
			continue
		}

		// Create delivery cost
		if err := database.DB.Create(&deliveryCost).Error; err != nil {
			return err
		}

		log.Printf("✓ Delivery cost created successfully (ID: %s)", deliveryCost.ID)
		log.Printf("  Location: %s", deliveryCost.Location)
		log.Printf("  Service: %s", deliveryCost.Service)
		log.Printf("  Cost: %.2f", deliveryCost.Cost)
	}

	log.Println("Delivery costs seeding completed!")
	return nil
}

func seedAdminUser() error {
	log.Println("Seeding admin user...")

	// Check if admin user already exists
	var existingAdmin models.User
	result := database.DB.Where("phone_number = ?", "+8801712402628").First(&existingAdmin)
	
	if result.Error == nil {
		log.Println("Admin user already exists, skipping...")
		return nil
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword("Welcome2Flexfume")
	if err != nil {
		return err
	}

	// Create admin user
	admin := models.User{
		Name:        "Flexfume Admin",
		PhoneNumber: "+8801712402628",
		Email:       "admin@flexfume.com",
		Password:    hashedPassword,
		Role:        "ADMIN",
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("✓ Admin user created successfully (ID: %s)", admin.ID)
	log.Println("  Name: Flexfume Admin")
	log.Println("  Phone: +8801712402628")
	log.Println("  Email: admin@flexfume.com")
	log.Println("  Role: ADMIN")
	log.Println("  Password: Welcome2Flexfume")

	return nil
}

func seedRegularUser() error {
	log.Println("Seeding regular user...")

	// Check if user already exists
	var existingUser models.User
	result := database.DB.Where("phone_number = ?", "+01712345678").First(&existingUser)
	
	if result.Error == nil {
		log.Println("Regular user already exists, skipping...")
		return nil
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword("Welcome2Flexfume")
	if err != nil {
		return err
	}

	// Create regular user
	user := models.User{
		Name:        "Test User",
		PhoneNumber: "+01712345678",
		Email:       "user@flexfume.com",
		Password:    hashedPassword,
		Role:        "USER",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	log.Printf("✓ Regular user created successfully (ID: %s)", user.ID)
	log.Println("  Name: Test User")
	log.Println("  Phone: +01712345678")
	log.Println("  Email: user@flexfume.com")
	log.Println("  Role: USER")
	log.Println("  Password: Welcome2Flexfume")

	return nil
}