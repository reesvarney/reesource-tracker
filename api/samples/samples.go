package samples

import (
	"apollo-sample-tracker/lib/database"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SampleData struct {
	database.ListSamplesRow
	Mods []database.SampleMod `json:"mods"`
}

func Routes(route *gin.RouterGroup) {
	route.GET("/samples", getSamples)
	route.GET("/sample/:sample_id", getSample)
	route.POST("/sample/:sample_id", updateSample)
}

func getSample(c *gin.Context) {
	sampleID := c.Param("sample_id")
	if sampleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sample ID is required"})
		return
	}
	// Expecting sampleID in the format "xx-xx-xx" (6 base36 chars, 3 pairs)
	parts := strings.Split(sampleID, "-")
	if len(parts) != 3 || len(parts[0]) != 2 || len(parts[1]) != 2 || len(parts[2]) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID format"})
		return
	}
	// Convert base36 pairs to bytes
	var rawID [4]byte
	for i, part := range parts {
		val, err := strconv.ParseUint(part, 36, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID format"})
			return
		}
		rawID[i] = byte(val)
	}
	RawSampleID := rawID[:]
	res, err := database.Connection.GetSampleById(c, RawSampleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	mod_data, err := database.Connection.ListSampleMods(c, RawSampleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sample": res, "mods": mod_data})
}

func updateSample(c *gin.Context) {
	// Update a specific sample by ID, or create a new sample with the provided uuid if it doesn't exist
	// These IDs will be provided by a QR code that opens the sample page
	sampleID := c.Param("sample_id")
	if sampleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sample ID is required"})
		return
	}
	// Expecting sampleID in the format "xx-xx-xx" (6 base36 chars, 3 pairs)
	parts := strings.Split(sampleID, "-")
	if len(parts) != 3 || len(parts[0]) != 2 || len(parts[1]) != 2 || len(parts[2]) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID format"})
		return
	}
	// Convert base36 pairs to bytes
	var rawID [4]byte
	for i, part := range parts {
		val, err := strconv.ParseUint(part, 36, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID format"})
			return
		}
		rawID[i] = byte(val)
	}
	RawSampleID := rawID[:]
	current_time := time.Now()
	res, err := database.Connection.UpdateOrCreateSample(c, database.UpdateOrCreateSampleParams{
		ID:             RawSampleID,
		LocationID:     c.PostForm("location_id"),
		ProductID:      c.PostForm("product_id"),
		TimeRegistered: current_time,
		LastUpdate:     current_time,
		State:          c.PostForm("state"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func getSamples(c *gin.Context) {
	res, err := database.Connection.ListSamples(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var samples []SampleData = []SampleData{}
	for _, sample := range res {
		// Convert time_registered to a string
		mod_data, err := database.Connection.ListSampleMods(c, sample.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		samples = append(samples, SampleData{sample, mod_data})
	}
	c.JSON(http.StatusOK, samples)
}
