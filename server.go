package escpos

import (
	"escpos/controller"
	"escpos/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	// controller.Printer()

	// Start Gin GO
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(CORS)

	routerGroup := router.Group("/api")

	NewRoutes(routerGroup)
	router.Run(":8000")
}

func NewRoutes(app *gin.RouterGroup) {

	app.GET("/", getAPI())
	app.POST("/print", controller.Printer())

}

func getAPI() func(c *gin.Context) {
	return func(c *gin.Context) {
		service.PrintTextWithSpaceTest(nil, "Subtotal: ", 412880.00, 35)
		service.PrintTextWithSpaceTest(nil, "Discount: ", 0.00, 35)
		service.PrintTextWithSpaceTest(nil, fmt.Sprintf("Service Charge(%d%%): ", 10), (41288), 35)

		service.PrintTextWithSpaceTest(nil, "Before VAT: ", 412880+41288, 35)
		service.PrintTextWithSpaceTest(nil, "VAT(7%): ", 31791.760000000002, 35)
		service.PrintTextWithSpaceTest(nil, "Rounding: ", 0.23999999999068677, 35)

		c.JSON(200, gin.H{
			"message": "API-ESCPOS",
		})
	}
}

func CORS(c *gin.Context) {
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
		return
	} else {
		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS
		// c.AbortWithStatusJSON(401, gin.H{
		// 	"code":    "ERROR",
		// 	"message": "Unauthorized",
		// })
		c.AbortWithStatus(http.StatusOK)
		return
	}
}
