package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Database handler
var db, err = gorm.Open("sqlite3", "test.db")

// Secret used for siging tokens
var secret = "secret"

// User model used for login and authorization
type user struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// Product model for database
type Product struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

// Login route handler function
// This method is used to generate new tokens for user "admin" & "user"
// Only "admin" has Admin rights for deleting, adding or updating products in database
func login(c echo.Context) (err error) {
	u := new(user)

	if err = c.Bind(u); err != nil {
		return echo.ErrUnsupportedMediaType
	}

	var user string
	var admin bool

	if u.Name == "admin" && u.Password == "admin" {
		user = "Admin"
		admin = true
	} else if u.Name == "user" && u.Password == "user" {
		user = "User"
		admin = false
	} else {
		return echo.ErrUnauthorized
	}

	fmt.Println(user)

	claims := &jwtCustomClaims{
		user,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 240).Unix(), // 10 days expiry time
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	// Returning token in a JSON format
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// Accessible route without JWT authentication
func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

// Get all products in database handler
func getAllProduct(c echo.Context) error {
	var products []Product

	if err := db.Find(&products).Error; err != nil {
		return c.JSON(404, "No products!!!")
	}

	return c.JSON(200, products)
}

// Get a specific product in databse by its id
func getProductByID(c echo.Context) error {
	product := new(Product)

	if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		return c.JSON(404, "Erorr!!!")
	}

	return c.JSON(200, product)
}

// Add a new product in datase
func putProduct(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)

	// If user is not admin return "forbidden" error message
	if claims.Admin != true {
		return echo.ErrForbidden
	}

	product := new(Product)
	c.Bind(product)

	db.Save(product)

	return c.JSON(200, product)
}

// Update a product in databasse
func updateProduct(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)

	// If user is not admin return "forbidden" error message
	if claims.Admin != true {
		return echo.ErrForbidden
	}

	product := new(Product)

	db.Where("id = ?", c.Param("id")).First(&product)

	c.Bind(product)

	db.Save(product)

	return c.JSON(200, product)
}

// Delete a product in database
func deleteProduct(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)

	// If user is not admin return "forbidden" error message
	if claims.Admin != true {
		return echo.ErrForbidden
	}

	db.Where("id = ?", c.Param("id")).Delete(&Product{})

	return c.JSON(200, "Delete success!!!"+c.Param("id"))
}

// Main function
func main() {
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Echo engine
	e := echo.New()

	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(secret),
	}

	r.Use(middleware.JWTWithConfig(config))

	r.GET("", getAllProduct)        // Get all products in database
	r.GET("/:id", getProductByID)   // Get a specific product in database
	r.POST("", putProduct)          // Add new product in database(only admin allowed)
	r.PATCH("/:id", updateProduct)  // Update product by id in database(only admin allowed)
	r.DELETE("/:id", deleteProduct) // Delete a product by id in database(only admin allowed)

	e.Logger.Fatal(e.Start(":8080"))
}
