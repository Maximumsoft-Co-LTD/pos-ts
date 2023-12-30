package controller

import (
	"encoding/json"

	"escpos/model"
	"escpos/service"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kenshaw/escpos"
)

func Printer() func(c *gin.Context) {
	type Body struct {
		IP   string     `form:"ip" json:"ip" binding:"required"`
		Bill model.Bill `bson:"bill_data" form:"bill_data" json:"bill_data" binding:"required"`
	}
	return func(c *gin.Context) {
		var Body Body
		if err := c.ShouldBind(&Body); err != nil {
			c.JSON(400, gin.H{
				"code":  400,
				"msg":   "ERROR BINDING",
				"error": err.Error(),
			})
			return
		}
		printerIP := Body.IP
		printerPort := 9100

		// Create a connection to the network printer
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", printerIP, printerPort))
		if err != nil {
			c.JSON(400, gin.H{
				"code":  400,
				"msg":   "Error Printer Connection",
				"error": err.Error(),
			})

			return
		}

		defer conn.Close()

		// Create a new ESC/POS printer
		printer := escpos.New(conn)

		// Initialize the printer

		// Load your JSON data here (data from test.json)
		data := Body.Bill
		// return
		orders := data.Order
		item := 0

		// Initialize the formatter for currency
		// formatter := NewCurrencyFormatter("th-TH", "THB")

		// Get the current date and time
		dataTime := data.Datetime
		dateTime := dataTime.Format("02/01/2006 15:04")

		fmt.Println(dateTime)

		// Print the receipt
		printer.Init()
		// if err != nil {
		// 	c.JSON(400, gin.H{
		// 		"code":  400,
		// 		"msg":   "ERROR CALL API DOMAIN",
		// 		"error": err.Error(),
		// 	})

		// 	return
		// }
		defer printer.End()
		defer printer.Cut()

		// printer.SetSmooth(128)
		printer.SetLang("en")
		printer.SetFont("B")
		printer.SetAlign("center")

		printer.SetFontSize(3, 2)

		// Print the header
		printer.SetEmphasize(1)
		printer.Write("MAX WALLET\n\n")

		printer.SetFontSize(1, 1)
		printer.SetFont("A")
		printer.Write("Receipt\n")
		printer.SetAlign("center")
		printer.Write("----------------------------------------------\n")

		// // Print table, cashier, customer, and date/time
		printer.SetAlign("left")
		printer.Write(fmt.Sprintf("TABLE: %s\n", data.Table))
		printer.Write(fmt.Sprintf("CASHIER: %s\n", data.Operator))

		var CustomerName string
		if len(data.Customer.Aka) == 0 {
			CustomerName = data.Customer.MemberID
		} else {
			CustomerName = data.Customer.Aka
		}

		printer.Write(fmt.Sprintf("CUSTOMER: %s\n", CustomerName))

		dateLine := fmt.Sprintf("DATE: %s", dataTime.Format("02/01/2006"))
		timeLine := fmt.Sprintf("TIME: %s", dataTime.Format("15:04"))
		totalLine := len(dateLine) + len(fmt.Sprint(timeLine))
		spacesNeeded := service.Max(0, 48-totalLine)
		spaces := service.Repeat(" ", spacesNeeded)
		dataText := dateLine + spaces + timeLine

		printer.SetAlign("left")
		printer.Write(dataText + "\n")
		printer.SetAlign("center")
		printer.Write("----------------------------------------------\n")

		// Print order details
		for _, order := range orders {

			var moreText string
			var defaultText string
			var orderLine string
			var PercentDecision float64 = 0.70
			var countTextOder int
			if len(order.Name) >= 20 {
				splitText := strings.Split(order.Name, " ")

				for i, v := range splitText {
					if countTextOder >= 24 {
						PercentDecision = 0.10
					}
					percentage := float64(len(splitText)) * PercentDecision
					count := int(percentage)

					countTextOder += len(v)

					if len(splitText) <= 2 && countTextOder >= 24 {
						totalCharacter := countTextOder - len(v)
						for _, eachCharacter := range v {
							characterAsString := string(eachCharacter)

							if totalCharacter <= 24 {
								totalCharacter += 1
								defaultText += characterAsString
							} else {
								moreText += characterAsString
							}
						}
					} else {
						if i >= count {
							moreText += v + " "
						} else {
							defaultText += v + " "
						}
					}

					orderLine = fmt.Sprintf("%d   %s ", order.Quantity, defaultText)
				}
			} else {
				orderLine = fmt.Sprintf("%d   %s", order.Quantity, order.Name)
			}

			// price := formatter.Format(order.Price)
			price := service.FormatCurrency(order.Price)
			totalLine := len(orderLine) + len(fmt.Sprint(price))
			spacesNeeded := service.Max(0, 48-totalLine)
			spaces := service.Repeat(" ", spacesNeeded)
			orderLineData := orderLine + spaces + price
			orderLineData = service.ReplaceThbSymbol(orderLineData)
			fmt.Println(orderLineData)

			if len(order.Name) >= 20 {
				fmt.Println()
				printer.Write(orderLineData + "\n")
				moreSpace := service.Max(0, 48-len(moreText))
				spaces := service.Repeat(" ", moreSpace)
				MoreOrderLineData := service.ReplaceThbSymbol(fmt.Sprintf("    %s", moreText) + spaces)
				printer.Write(MoreOrderLineData + "\n")
			} else {
				printer.Write(orderLineData + "\n\n")
			}
			item += order.Quantity
		}

		// // Calculate and print the summary
		printer.SetAlign("center")
		printer.Write("----------------------------------------------\n")
		printer.SetAlign("left")
		printer.Write(fmt.Sprintf("ITEMS: %d\n", item))

		grandTotal := data.PriceWithDiscount + data.ServiceCharge

		maxLengthAmount := 0
		for _, eachAmount := range []float64{data.Price, data.Discount, data.ServiceCharge, grandTotal, data.Vat, data.Rounding} {
			stringAmount := service.FormatCurrency(eachAmount)
			if len(stringAmount) > maxLengthAmount {

				maxLengthAmount = len(stringAmount)
				fmt.Println(stringAmount, maxLengthAmount)
			}
		}

		fmt.Println("maxLengthAmount", maxLengthAmount)
		// var maxLengthPriceEndSlip int
		// // Print subtotal, discount, service charge, and other details
		service.PrintTextWithSpace(printer, "Subtotal: ", data.Price, 35, maxLengthAmount)
		service.PrintTextWithSpace(printer, "Discount: ", data.Discount, 35, maxLengthAmount)
		service.PrintTextWithSpace(printer, fmt.Sprintf("Service Charge(%d%%): ", data.PercentServiceCharge), (data.ServiceCharge), 35, maxLengthAmount)

		service.PrintTextWithSpace(printer, "Before VAT: ", grandTotal, 35, maxLengthAmount)
		service.PrintTextWithSpace(printer, "VAT(7%): ", data.Vat, 35, maxLengthAmount)
		service.PrintTextWithSpace(printer, "Rounding: ", data.Rounding, 35, maxLengthAmount)

		printer.SetAlign("right")
		printer.Write("=========================\n")
		service.PrintTextWithSpace(printer, "Total: ", (data.Total), 24, 0)
		printer.SetAlign("right")
		printer.Write("=========================\n")

		printer.SetAlign("center")
		printer.Write("----------------------------------------------\n")
		printer.SetAlign("center")
		printer.Write("Thank you for your order!\n")

	}
}

func loadData() model.Bill {
	// Read the JSON file
	fileContent, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Create an instance of the struct to hold the parsed data
	var bill model.Bill

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(fileContent, &bill)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Access the parsed data
	fmt.Printf("Name: %s\n", bill.Customer.Name)
	return bill
}
