package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserRole enum type
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
	RoleGuest UserRole = "guest"
	RoleSuperAdmin UserRole = "superadmin"
)

// User models
type UserBase struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	FullName string   `json:"full_name,omitempty"`
	Role     UserRole `json:"role" binding:"required,oneof=admin user guest superadmin"`
}

type UserCreate struct {
	UserBase
	Password string `json:"password" binding:"required,min=8"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name,omitempty"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// Item models
type ItemBase struct {
	Name        string   `json:"name" binding:"required,min=1,max=100"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Tax         *float64 `json:"tax,omitempty" binding:"omitempty,gte=0"`
}

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Tax         float64 `json:"tax,omitempty"`
	OwnerID     int     `json:"owner_id"`
}

// In-memory databases
var (
	usersDB       []User
	itemsDB       []Item
	userIDCounter = 1
	itemIDCounter = 1
	validate      = validator.New()
)

// Error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success response
type MessageResponse struct {
	Message string `json:"message"`
}

func main() {
	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Root endpoints
	router.GET("/", rootHandler)
	router.GET("/health", healthCheckHandler)

	// User routes
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", createUserHandler)
		userRoutes.GET("/", getUsersHandler)
		userRoutes.GET("/:id", getUserHandler)
		userRoutes.PUT("/:id", updateUserHandler)
		userRoutes.DELETE("/:id", deleteUserHandler)
		userRoutes.POST("/:id/items", createItemHandler)
		userRoutes.GET("/:id/items", getUserItemsHandler)
	}

	// Item routes
	itemRoutes := router.Group("/items")
	{
		itemRoutes.GET("/", getAllItemsHandler)
		itemRoutes.GET("/:id", getItemHandler)
	}

	// Start server
	println("🚀 Server starting on http://localhost:8080")
	router.Run(":8080")
}

// Root handler
func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Sample Go/Gin Application",
		"version": "1.0.0",
	})
}

// Health check handler
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Create user handler
func createUserHandler(c *gin.Context) {
	var userCreate UserCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if username exists
	for _, u := range usersDB {
		if u.Username == userCreate.Username {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Username already exists. Please log in."})
			return
		}
	}

	// Check if email exists
	for _, u := range usersDB {
		if u.Email == userCreate.Email {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email already exists"})
			return
		}
	}

	// Create new user
	newUser := User{
		ID:        userIDCounter,
		Username:  userCreate.Username,
		Email:     userCreate.Email,
		FullName:  userCreate.FullName,
		Role:      userCreate.Role,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	usersDB = append(usersDB, newUser)
	userIDCounter++

	c.JSON(http.StatusCreated, newUser)
}

// Get users handler with pagination and filtering
func getUsersHandler(c *gin.Context) {
	// Parse query parameters
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	roleFilter := c.Query("role")

	// Validate parameters
	if skip < 0 {
		skip = 0
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Filter users
	filteredUsers := usersDB
	if roleFilter != "" {
		var filtered []User
		for _, u := range usersDB {
			if string(u.Role) == roleFilter {
				filtered = append(filtered, u)
			}
		}
		filteredUsers = filtered
	}

	// Apply pagination
	start := skip
	end := skip + limit

	if start >= len(filteredUsers) {
		c.JSON(http.StatusOK, []User{})
		return
	}

	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	c.JSON(http.StatusOK, filteredUsers[start:end])
}

// Get user by ID handler
func getUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Find user
	for _, u := range usersDB {
		if u.ID == id {
			c.JSON(http.StatusOK, u)
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
}

// Update user handler
func updateUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	var userUpdate UserBase
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Find and update user
	for i, u := range usersDB {
		if u.ID == id {
			usersDB[i].Username = userUpdate.Username
			usersDB[i].Email = userUpdate.Email
			usersDB[i].FullName = userUpdate.FullName
			usersDB[i].Role = userUpdate.Role

			c.JSON(http.StatusOK, usersDB[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
}

// Delete user handler
func deleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Find and delete user
	for i, u := range usersDB {
		if u.ID == id {
			usersDB = append(usersDB[:i], usersDB[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
}

// Create item handler
func createItemHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)

	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Check if user exists
	userExists := false
	for _, u := range usersDB {
		if u.ID == userID {
			userExists = true
			break
		}
	}

	if !userExists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	var itemBase ItemBase
	if err := c.ShouldBindJSON(&itemBase); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Create new item
	tax := 0.0
	if itemBase.Tax != nil {
		tax = *itemBase.Tax
	}

	newItem := Item{
		ID:          itemIDCounter,
		Name:        itemBase.Name,
		Description: itemBase.Description,
		Price:       itemBase.Price,
		Tax:         tax,
		OwnerID:     userID,
	}

	itemsDB = append(itemsDB, newItem)
	itemIDCounter++

	c.JSON(http.StatusCreated, newItem)
}

// Get user items handler
func getUserItemsHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)

	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Check if user exists
	userExists := false
	for _, u := range usersDB {
		if u.ID == userID {
			userExists = true
			break
		}
	}

	if !userExists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Get user's items
	var userItems []Item
	for _, item := range itemsDB {
		if item.OwnerID == userID {
			userItems = append(userItems, item)
		}
	}

	if userItems == nil {
		userItems = []Item{}
	}

	c.JSON(http.StatusOK, userItems)
}

// Get all items handler
func getAllItemsHandler(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	// Validate parameters
	if skip < 0 {
		skip = 0
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Filter items
	filteredItems := itemsDB

	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			var filtered []Item
			for _, item := range filteredItems {
				if item.Price >= minPrice {
					filtered = append(filtered, item)
				}
			}
			filteredItems = filtered
		}
	}

	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			var filtered []Item
			for _, item := range filteredItems {
				if item.Price <= maxPrice {
					filtered = append(filtered, item)
				}
			}
			filteredItems = filtered
		}
	}

	// Apply pagination
	start := skip
	end := skip + limit

	if start >= len(filteredItems) {
		c.JSON(http.StatusOK, []Item{})
		return
	}

	if end > len(filteredItems) {
		end = len(filteredItems)
	}

	c.JSON(http.StatusOK, filteredItems[start:end])
}

// Get item by ID handler
func getItemHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid item ID"})
		return
	}

	// Find item
	for _, item := range itemsDB {
		if item.ID == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{Error: "Item not found"})
}
