package locations

import (
	"apollo-sample-tracker/lib/database"

	"github.com/gin-gonic/gin"
)

func Routes(route *gin.RouterGroup) {
	route.GET("/locations", getLocations)
	route.GET("/location", getLocation)
	route.POST("/location", updateLocation)
}

func getLocations(c *gin.Context) {
	res, err := database.Connection.ListLocations(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve locations"})
		return
	}
	c.JSON(200, res)
}

func getLocation(c *gin.Context) {
}

func updateLocation(c *gin.Context) {
}
