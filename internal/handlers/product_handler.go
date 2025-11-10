package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	productService services.ProductServiceInterface
	s3Service      *services.S3Service
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService services.ProductServiceInterface, s3Service *services.S3Service) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		s3Service:      s3Service,
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name              string   `json:"name" binding:"required"`
	Code              string   `json:"code"`
	Description       string   `json:"description"`
	LongDescription   string   `json:"longDescription"`
	Price             float64  `json:"price" binding:"required,gt=0"`
	CategoryID        string   `json:"categoryId" binding:"required"`
	Stock             int      `json:"stock"`
	ImageURL          string   `json:"imageUrl"`
	YouTubeVideoUrl   string   `json:"youtubeVideoUrl"`
	KeyFeatures       []string `json:"keyFeatures" binding:"omitempty"`
	WhyChooseBenefits []string `json:"whyChooseBenefits" binding:"omitempty"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name              string   `json:"name"`
	Code              string   `json:"code"`
	Description       string   `json:"description"`
	LongDescription   string   `json:"longDescription"`
	Price             float64  `json:"price"`
	CategoryID        string   `json:"categoryId"`
	Stock             int      `json:"stock"`
	ImageURL          string   `json:"imageUrl"`
	YouTubeVideoUrl   string   `json:"youtubeVideoUrl"`
	KeyFeatures       []string `json:"keyFeatures" binding:"omitempty"`
	WhyChooseBenefits []string `json:"whyChooseBenefits" binding:"omitempty"`
}

// CreateProduct creates a new product (Admin only)
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	product := &models.Product{
		ProductName:       req.Name,
		Description:       req.Description,
		LongDescription:   req.LongDescription,
		Price:             req.Price,
		Stock:             req.Stock,
		CategoryID:        req.CategoryID,
		ImageURL:          req.ImageURL,
		YouTubeVideoUrl:   req.YouTubeVideoUrl,
		KeyFeatures:       req.KeyFeatures,
		WhyChooseBenefits: req.WhyChooseBenefits,
	}

	createdProduct, err := h.productService.CreateProduct(product)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: createdProduct.ToResponse(),
	})
}

// UpdateProduct updates an existing product (Admin only)
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	product := &models.Product{
		ProductName:       req.Name,
		ProductCode:       req.Code,
		Description:       req.Description,
		LongDescription:   req.LongDescription,
		Price:             req.Price,
		Stock:             req.Stock,
		CategoryID:        req.CategoryID,
		ImageURL:          req.ImageURL,
		YouTubeVideoUrl:   req.YouTubeVideoUrl,
		KeyFeatures:       req.KeyFeatures,
		WhyChooseBenefits: req.WhyChooseBenefits,
	}

	updatedProduct, err := h.productService.UpdateProduct(productID, product)
	if err != nil {
		if err == services.ErrProductNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: updatedProduct.ToResponse(),
	})
}

// GetAllProducts retrieves all products
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	responses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = product.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// GetProductByID retrieves a product by ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")

	product, err := h.productService.GetProductByID(productID)
	if err != nil {
		if err == services.ErrProductNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: product.ToResponse(),
	})
}

// DeleteProduct deletes a product (Admin only)
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	err := h.productService.DeleteProduct(productID)
	if err != nil {
		if err == services.ErrProductNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "DELETED",
	})
}

// GetProductsByCategory retrieves all products in a category
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("id")

	products, err := h.productService.GetProductsByCategory(categoryID)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrNoProductsFoundByCategory {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	responses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = product.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// SearchProducts searches for products by name or description
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	searchValue := c.Query("value")
	if searchValue == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "search value is required",
		})
		return
	}

	products, err := h.productService.SearchProducts(searchValue)
	if err != nil {
		if err == services.ErrNoProductsFoundBySearch {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	responses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = product.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// UploadProductImage uploads a product image (Admin only)
func (h *ProductHandler) UploadProductImage(c *gin.Context) {
	productID := c.Param("id")

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "file is required",
		})
		return
	}

	// Validate image file
	if h.s3Service != nil {
		if err := h.s3Service.ValidateImageFile(fileHeader); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "failed to open uploaded file",
		})
		return
	}
	defer file.Close()

	var imageURL string

	// Upload to S3 if service is available
	if h.s3Service != nil {
		imageURL, err = h.s3Service.UploadFile(file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "failed to upload image: " + err.Error(),
			})
			return
		}
	} else {
		// Fallback to local path if S3 is not configured
		imageURL = "/uploads/" + fileHeader.Filename
	}

	// Update product with new image URL
	updatedProduct, err := h.productService.UploadProductImage(productID, imageURL)
	if err != nil {
		if err == services.ErrProductNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: updatedProduct.ToResponse(),
	})
}
