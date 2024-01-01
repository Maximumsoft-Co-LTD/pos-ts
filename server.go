package escpos

import (
	"context"
	"encoding/json"
	"escpos/controller"
	"escpos/model"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func StartServer() {
	// controller.Printer()

	// // Start Gin GO
	// router := gin.Default()
	// router.Use(gin.Recovery())
	// // router.Use(CORS)

	// routerGroup := router.Group("/api")

	// NewRoutes(routerGroup)
	// router.Run(":8000")
	fmt.Println("START Bill Application")
	app, err := initFirebaseWithStringCredentials()
	if err != nil {
		fmt.Println("ERROR Init Firebase")
		return
	}

	client, err := app.Database(context.Background())
	if err != nil {
		// Handle error
		fmt.Println("ERROR CONNECT", err)
		return
	}
	refAllBill := client.NewRef("bills/")
	// Reference the desired data location

	// read from user_scores using ref
	// type UserScore struct {
	// 	Score int `json:"score"`
	// }
	// for  {

	// }
	for {

		fmt.Println("Start Get Bills")
		var bills map[string]interface{}

		if err := refAllBill.Get(context.TODO(), &bills); err != nil {
			log.Fatalln("error in reading from firebase DB: ", err)
			continue
		}
		fmt.Println("Total Bills", len(bills))
		for i, v := range bills {
			fmt.Println("Key", i)
			var bill model.Bill
			billbytes, _ := json.Marshal(v)
			json.Unmarshal(billbytes, &bill)

			controller.Printer("printer-maxclub.maxpos.io", bill)
			refEachBill := client.NewRef("bills/" + i + "/")

			if err := refEachBill.Delete(context.TODO()); err != nil {
				log.Fatalln("error in deleting ref: ", err)
				continue
			}
			fmt.Println("deleted successfully:)")
			time.Sleep(500 * time.Millisecond)
		}
		time.Sleep(1 * time.Second)
	}

	// ref := client.NewRef("user_scores/1")

	// Listen for data changes
	// listener := ref.Listen(context.Background())
	// defer listener.Close() // Close the listener when finished

	// go func() {
	// 	for {
	// 		snap, err := listener.Next()
	// 		if err != nil {
	// 			// Handle error
	// 		}

	// 		// Create a worker process for each new data change
	// 		go func() {
	// 			// Process the new data within this worker process
	// 			fmt.Println("New data:", snap.Val())
	// 			// Perform the necessary actions based on the data
	// 		}()
	// 	}
	// }()
}

func initFirebaseWithStringCredentials() (*firebase.App, error) {
	credentials := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  "maxwallet",
		"private_key_id":              "ac44791f21c56a053423c47490224c79226f2793",
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC+H7Igend6GcBD\nT1BSILnSmP1e53so6XZTLw0k+m9605DFDGX07Ez73PTLJgqR5CUDLK9oszasW2h3\nnhkqlLFw82wHFVmJ8HlkarxDQNLkwgqUBL6T6IxaHDmVrjUZTCY0QXLsjGBEwBgG\nwIFjhO2x95EYQNUVhcdjqLfEioigiYE9VTL4rQX5QFSWOTelTTHAW30fDlJzgrpy\npXNqX0bs7F/4sNMOqC78KpozwCVnHxQ5DTt3MPibNNsu7+wEGx+Pi0K3zr20O+Zd\nVcyUs0hlk+1Ga/ErUu+A5oJyAi800EqrVNUkxJFp+5xxpAYQSjInF1oKZgCAKo5D\nND+yEHKXAgMBAAECggEAMrH66kQdhZZGOOx2K2AEjI40yVCJ+4+8+FNMeva4MRUm\nvhQfP56erG3vbBfZq6sc32rT3uQiiPTat0KVzU/WvJp7zKu5s1lG1SrMxlYLpenD\nrlJRitjwS7rI9At+px8x9K8a+M45gTuDbJF6LtIsG087Bi+LYfDJCN7dGXAjHEL3\nbORbkWTilIDgqThdK7JtDzDvTUUkzghzzncM6AcyEx2wFbLiLy+6FZjJvaobTk93\nDwChSK7jM0Rj/e1gPGjgRe/mZh/3qDIN3mborO8tRcZ31yMbEBAt+OeRDfm8oLDB\nsJXU4HlpynLcmv2UmtHQbpzgvjI5LW9uMZOgRjIqcQKBgQD67ShpVRq/3kc2DPVo\n5OR9Jma5xJFv+fXeEegSj3nKLhfmdotRBW57zfWLrL9la/D1udIZCmWQVAgkD/xC\nldJBCfUMD0KKK4zP4IX7okEiv53Y7kIt25JBREEVMvdWRvkf0XAS39huqNLD2Nqc\nrCUTfUsvNOwmWWzZ3QQxYDqW7wKBgQDB98/2MRQzXOsgaO7dNbLk5EmqFIAqQTgZ\n6DDtmzasXcuQlyAEH2BxokbwV77x3wq0tMPLowKFPOHPAnQA8zZqXmHWLSoJx8xX\nig7lU6mBJlavFAOzWxK7YFMRXw35OZZhuSEyKFoEjJ3skPYRvMSxbitmBEvY5X92\n4sULJLGe2QKBgDtNy6xBWeYY4ZmrrGCTIFFXvxWOmJTvbaWDc+bXFACtriZgxAJt\nFzSZc7wEIuQUg8l/lmEmrORUh+wF/ye5gwyuDsU/4gkHy+rhdKkJKv3MbcD0Zp2x\n9DoKqJsbBYvVkFFtzWAYmW1l5xI0cU5v/P9DMH7CskFKB0jiTHhi9tXvAoGBAME2\nmt8/8EFhs41rOVT84qCCjqZvGWP952ZXFjX5QLLeE6KKB4hTwPwwi71pinjglod+\n8PJuBFq4VK6iYO95VaELNyXjg1aOwYwJp+DkP5q4l+x6YV1NwREJWHWbXA4AQT5C\n7UBLVa/maoF1vMfaY4vilDRg3zTFFv1T0rfQ62WxAoGAa4UK2NaX1qJzTOb0+IdM\nlnpDOtI5tmgxsS8X/Uz/AxzRjfdmDSmNUI5IM+AYphv1Jy4vdGHk3ST3OiKJAL8x\nq8nq+GLVRhvgblSQuBWja+v7to/W5drG7ndaYVyniX9W0BvC8NgBEHHxomlfnZR6\nZxYqpOHrQJQcysG1DN86ToQ=\n-----END PRIVATE KEY-----\n",
		"client_email":                "firebase-adminsdk-l1ved@maxwallet.iam.gserviceaccount.com",
		"client_id":                   "108462313032826368896",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-l1ved%40maxwallet.iam.gserviceaccount.com",
		"universe_domain":             "googleapis.com",
	}

	credentialsString, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("Error marshaling credentials:", err)
		return nil, err
	}
	conf := &firebase.Config{
		DatabaseURL: "https://maxwallet-default-rtdb.asia-southeast1.firebasedatabase.app",
	}

	opt := option.WithCredentialsJSON(credentialsString)
	return firebase.NewApp(context.Background(), conf, opt)
}

func NewRoutes(app *gin.RouterGroup) {

	app.GET("/", getAPI())
	// app.POST("/print", controller.Printer())

}

func getAPI() func(c *gin.Context) {
	return func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "API-ESCPOS",
		})
	}
}

// func CORS(c *gin.Context) {
// 	// First, we add the headers with need to enable CORS
// 	// Make sure to adjust these headers to your needs
// 	c.Header("Access-Control-Allow-Origin", "*")
// 	c.Header("Access-Control-Allow-Methods", "*")
// 	c.Header("Access-Control-Allow-Headers", "*")
// 	c.Header("Content-Type", "application/json")

// 	// Second, we handle the OPTIONS problem
// 	if c.Request.Method != "OPTIONS" {
// 		c.Next()
// 		return
// 	} else {
// 		// Everytime we receive an OPTIONS request,
// 		// we just return an HTTP 200 Status Code
// 		// Like this, Angular can now do the real
// 		// request using any other method than OPTIONS
// 		// c.AbortWithStatusJSON(401, gin.H{
// 		// 	"code":    "ERROR",
// 		// 	"message": "Unauthorized",
// 		// })
// 		c.AbortWithStatus(http.StatusOK)
// 		return
// 	}
// }
