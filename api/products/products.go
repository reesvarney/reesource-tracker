package products

import (
	"apollo-sample-tracker/lib/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(route *gin.RouterGroup) {
	route.GET("/products", getProducts)
	route.POST("/product", createProduct)
	route.GET("/product/[product_id]", getProduct)
	route.POST("/product/[product_id]", updateProduct)
}

func createProduct(c *gin.Context) {
	var req struct {
		Name            string  `json:"name"`
		ParentProductID *string `json:"parent_product_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// UpsertProduct expects UpsertProductParams struct
	params := database.UpsertProductParams{
		Name:            req.Name,
		ParentProductID: req.ParentProductID,
	}
	err := database.Connection.UpsertProduct(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
		return
	}
	product, err := database.Connection.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func updateProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
		return
	}

	var req struct {
		Name            string  `json:"name"`
		ParentProductID *string `json:"parent_product_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// UpsertProduct expects UpsertProductParams struct
	params := database.UpsertProductParams{
		ID:              productID,
		Name:            req.Name,
		ParentProductID: req.ParentProductID,
	}
	err := database.Connection.UpsertProduct(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getProducts(c *gin.Context) {
	res, err := database.Connection.GetProducts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
