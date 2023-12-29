package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/kenshaw/escpos"
)

func main() {
	printerIP := "192.168.1.148"
	printerPort := 9100

	// Create a connection to the network printer
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", printerIP, printerPort))
	if err != nil {
		log.Fatalf("Error connecting to the printer: %v", err)
	}
	defer conn.Close()

	// Create a new ESC/POS printer
	printer := escpos.New(conn)

	// Initialize the printer

	// Load your JSON data here (data from test.json)
	data := loadData()
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
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer printer.End()
	defer printer.Cut()

	// p.SetSmooth(1)
	// p.SetFontSize(2, 3)
	// p.SetFont("A")
	// p.Write("test ")
	// p.SetFont("B")
	// p.Write("test2 ")
	// p.SetFont("C")
	// p.Write("test3 ")
	// p.Formfeed()

	// p.SetFont("B")
	// p.SetFontSize(1, 1)

	// p.SetEmphasize(1)
	// p.Write("halle")
	// p.Formfeed()

	// p.SetUnderline(1)
	// p.SetFontSize(4, 4)
	// p.Write("halle")

	// p.SetReverse(1)
	// p.SetFontSize(2, 4)
	// p.Write("halle")
	// p.Formfeed()

	// p.SetFont("C")
	// p.SetFontSize(8, 8)
	// p.Write("halle")
	// p.FormfeedN(5)

	// Set font, style, and size

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
	printer.Write(fmt.Sprintf("CUSTOMER: %s\n", data.Customer.Aka))

	dateLine := fmt.Sprintf("DATE: %s", dataTime.Format("02/01/2006"))
	timeLine := fmt.Sprintf("TIME: %s", dataTime.Format("15:04"))
	totalLine := len(dateLine) + len(fmt.Sprint(timeLine))
	spacesNeeded := Max(0, 48-totalLine)
	spaces := repeat(" ", spacesNeeded)
	dataText := dateLine + spaces + timeLine

	printer.SetAlign("left")
	printer.Write(dataText + "\n")
	printer.SetAlign("center")
	printer.Write("----------------------------------------------\n")

	// Print order details
	for _, order := range orders {
		orderLine := fmt.Sprintf("%d %s (%s)", order.Quantity, order.Name, order.Size)
		// price := formatter.Format(order.Price)
		price := formatCurrency(order.Price)
		totalLine := len(orderLine) + len(fmt.Sprint(price))
		spacesNeeded := Max(0, 48-totalLine)
		spaces := repeat(" ", spacesNeeded)
		orderLineData := orderLine + spaces + price
		orderLineData = replaceThbSymbol(orderLineData)
		fmt.Println(orderLineData)
		printer.Write(orderLineData + "\n\n")
		item += order.Quantity
	}

	// // Calculate and print the summary
	printer.SetAlign("center")
	printer.Write("----------------------------------------------\n")
	printer.SetAlign("left")
	printer.Write(fmt.Sprintf("ITEMS: %d\n", item))

	// // Print subtotal, discount, service charge, and other details
	printTextWithSpace(printer, "Subtotal: ", data.Price, 24)
	printTextWithSpace(printer, "Discount: ", data.Discount, 24)
	printTextWithSpace(printer, fmt.Sprintf("Service Charge(%d%%): ", data.PercentServiceCharge), (data.ServiceCharge), 24)

	grandTotal := data.PriceWithDiscount + data.ServiceCharge
	printTextWithSpace(printer, "Before VAT: ", grandTotal, 24)
	printTextWithSpace(printer, "VAT(7%): ", data.Vat, 24)
	printTextWithSpace(printer, "Rounding: ", data.Rounding, 24)

	printer.SetAlign("right")
	printer.Write("=========================\n")
	printTextWithSpace(printer, "Total: ", (data.Total), 24)
	printer.SetAlign("right")
	printer.Write("=========================\n")

	printer.SetAlign("center")
	printer.Write("----------------------------------------------\n")
	printer.SetAlign("center")
	printer.Write("Thank you for your order!\n")
}

func loadData() Bill {
	// Read the JSON file
	fileContent, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Create an instance of the struct to hold the parsed data
	var bill Bill

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(fileContent, &bill)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Access the parsed data
	fmt.Printf("Name: %s\n", bill.Customer.Name)
	return bill
}

// Format formats a currency value.
func formatCurrency(value float64) string {
	// Format the number with a thousands separator and two decimal places
	formatted := strconv.FormatFloat(value, 'f', 2, 64)

	// Add thousands separator (',') to the integer part of the number
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	var formattedWithComma string

	// Add a comma as a thousands separator
	for i, digit := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			formattedWithComma += ","
		}
		formattedWithComma += string(digit)
	}

	// Combine the integer and decimal parts
	if decimalPart != "" {
		formattedWithComma += "." + decimalPart
	}

	return formattedWithComma
}

// Max returns the maximum of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Repeat returns a string consisting of repeated copies of s.
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// ReplaceThbSymbol replaces the THB currency symbol with an empty string.
func replaceThbSymbol(s string) string {
	return strings.ReplaceAll(s, "฿", "")
}

// PrintTextWithSpace prints text with spaces to reach a specified length.
func printTextWithSpace(printer *escpos.Escpos, text string, data float64, length int) {
	dataCurrency := formatCurrency(data)
	totalLine := len(text) + len(fmt.Sprint(dataCurrency))
	spacesNeeded := Max(0, length-totalLine)
	spaces := repeat(" ", spacesNeeded)
	textData := text + spaces + fmt.Sprint(dataCurrency)
	textData = replaceThbSymbol(textData)
	printer.SetAlign("right")
	if text == "Total: " {
		printer.Write(textData + "\n")
	} else {
		printer.Write(textData + "\n\n")
	}

}

type Bill struct {
	ID       string    `json:"id"`
	Datetime time.Time `json:"datetime"`
	BillDate time.Time `json:"bill_date"`
	BillID   string    `json:"bill_id"`
	BillType struct {
		Status int    `json:"status"`
		Label  string `json:"label"`
	} `json:"bill_type"`
	Table  string `json:"table"`
	Status struct {
		Datetime    time.Time `json:"datetime"`
		Status      int       `json:"status"`
		Label       string    `json:"label"`
		Description string    `json:"description"`
		By          string    `json:"by"`
	} `json:"status"`
	StatusHistory []struct {
		Datetime    time.Time `json:"datetime"`
		Status      int       `json:"status"`
		Label       string    `json:"label"`
		Description string    `json:"description"`
		By          string    `json:"by"`
	} `json:"status_history"`
	Customer struct {
		Name         string `json:"name"`
		MemberID     string `json:"member_id"`
		Aka          string `json:"aka"`
		Rank         int    `json:"rank"`
		CustomerType string `json:"customer_type"`
	} `json:"customer"`
	Order []struct {
		ID            string    `json:"id"`
		Datetime      time.Time `json:"datetime"`
		Name          string    `json:"name"`
		Image         string    `json:"image"`
		Price         float64   `json:"price"`
		Category      int       `json:"category"`
		Subcategory   int       `json:"subcategory"`
		Size          string    `json:"size"`
		Quantity      int       `json:"quantity"`
		OrderTimeline []struct {
			Datetime    time.Time `json:"datetime"`
			Status      int       `json:"status"`
			Label       string    `json:"label"`
			Description string    `json:"description"`
			By          string    `json:"by"`
		} `json:"order_timeline"`
		Status struct {
			Datetime    time.Time `json:"datetime"`
			Status      int       `json:"status"`
			Label       string    `json:"label"`
			Description string    `json:"description"`
			By          string    `json:"by"`
		} `json:"status"`
		Seller         string `json:"seller"`
		SellerID       string `json:"seller_id"`
		SellerName     string `json:"seller_name"`
		SellerUsername string `json:"seller_username"`
		SellerAka      string `json:"seller_aka"`
		Package        bool   `json:"package"`
		Points         int    `json:"points"`
		Commission     bool   `json:"commission"`
	} `json:"order"`
	Operator             string        `json:"operator"`
	DiscountList         []interface{} `json:"discount_list"`
	Discount             float64       `json:"discount"`
	Price                float64       `json:"price"`
	Points               float64       `json:"points"`
	PriceWithDiscount    float64       `json:"price_with_discount"`
	PaymentProof         string        `json:"payment_proof"`
	ServiceCharge        float64       `json:"service_charge"`
	Vat                  float64       `json:"vat"`
	Rounding             float64       `json:"rounding"`
	Total                float64       `json:"total"`
	Payment              string        `json:"payment"`
	PercentServiceCharge int8          `json:"percent_service_charge"`
	IsServiceCharge      bool          `json:"is_service_charge"`
	IsVat                bool          `json:"is_vat"`
}