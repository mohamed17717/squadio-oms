package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ViewSet[T any, C any, U any] struct {
	DB                   *gorm.DB
	PerformCreateFunc    func(c *gin.Context, obj *T) error
	InputOfCreateToModel func(n *C) T
	InputOfUpdateToModel func(n *U) T
}

func (v ViewSet[T, C, U]) Retrieve(c *gin.Context) {
	var obj T
	id := c.Param("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := v.DB.First(&obj, uuidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	c.JSON(http.StatusOK, obj)
}

func (v ViewSet[T, C, U]) List(c *gin.Context) {
	var objs []T

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	active := c.Query("active")
	includeVariants := c.Query("include_variants") == "true"
	offset := (page - 1) * limit

	var model T
	query := v.DB.Model(&model)

	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if active != "" {
		query = query.Where("is_active = ?", active == "true")
	}

	// Get total count
	var total int64
	query.Count(&total)

	if includeVariants {
		query = query.Preload("Variants")
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&objs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch objects"})
		return
	}

	// Convert to response format
	responses := make([]T, len(objs))
	copy(responses, objs)

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (v ViewSet[T, C, U]) Create(c *gin.Context) {
	var input C
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var obj T = v.InputOfCreateToModel(&input)

	// Call the injected custom create logic
	if v.PerformCreateFunc != nil {
		if err := v.PerformCreateFunc(c, &obj); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Save the object after performing custom logic
	if err := v.DB.Create(&obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create object"})
		return
	}

	c.JSON(http.StatusOK, obj)
}

func (v ViewSet[T, C, U]) Update(c *gin.Context) {
	var obj T
	id := c.Param("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := v.DB.First(&obj, uuidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	var input U
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build the updates struct from input without overwriting the loaded object's primary key
	updates := v.InputOfUpdateToModel(&input)

	// Apply updates onto the existing row using its bound primary key (obj)
	// Omit immutable fields like ID (and optionally CreatedAt if present on the model)
	if err := v.DB.Model(&obj).Omit("id").Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update object"})
		return
	}

	// Re-fetch to return the latest state after update
	if err := v.DB.First(&obj, uuidID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load updated object"})
		return
	}

	c.JSON(http.StatusOK, obj)
}

func (v ViewSet[T, C, U]) Delete(c *gin.Context) {
	var obj T
	id := c.Param("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := v.DB.First(&obj, uuidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	if err := v.DB.Delete(&obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete object"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object deleted"})
}
