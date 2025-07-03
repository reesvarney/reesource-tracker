package api

import (
	"apollo-sample-tracker/api/locations"
	"apollo-sample-tracker/api/products"
	"apollo-sample-tracker/api/samples"

	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {
	api_routes := route.Group("/api")
	samples.Routes(api_routes)
	products.Routes(api_routes)
	locations.Routes(api_routes)
}
