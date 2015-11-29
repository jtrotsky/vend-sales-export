package writer

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jtrotsky/govend/vend"
)

// CreateReport ...
func CreateReport(domainPrefix string) (string, error) {

	// Create blank CSV file to be written to.
	// File name will be the current time in unixtime.
	fname := fmt.Sprintf("%s_sales_history_%v.csv", domainPrefix,
		time.Now().Unix())
	f, err := os.Create(fmt.Sprintf("./%s", fname))
	if err != nil {
		log.Fatalf("Error creating CSV file: %s", err)
	}
	// Ensure file closes at end.
	defer f.Close()

	w := csv.NewWriter(f)

	var headerLine []string
	headerLine = append(headerLine, "Sale Date")
	headerLine = append(headerLine, "Sale Time")
	headerLine = append(headerLine, "Invoice Number")
	headerLine = append(headerLine, "Line Type")
	headerLine = append(headerLine, "Customer Code")
	headerLine = append(headerLine, "Company Name")
	headerLine = append(headerLine, "Customer Name")
	headerLine = append(headerLine, "Sale Note")
	headerLine = append(headerLine, "Quantity")
	headerLine = append(headerLine, "Price")
	headerLine = append(headerLine, "Tax")
	headerLine = append(headerLine, "Discount")
	headerLine = append(headerLine, "Loyalty")
	headerLine = append(headerLine, "Total")
	headerLine = append(headerLine, "Paid")
	headerLine = append(headerLine, "Details")
	headerLine = append(headerLine, "Register")
	headerLine = append(headerLine, "User")
	headerLine = append(headerLine, "Status")
	headerLine = append(headerLine, "Product Sku")

	w.Write(headerLine)
	w.Flush()

	return fname, err
}

// WriteReport aims to mimic the report generated by
// exporting Vend sales history.
func WriteReport(fname string, registers *[]vend.Register, users *[]vend.User,
	customers *[]vend.Customer, products *[]vend.Product,
	sales []vend.Sale, domainPrefix, tz string) {

	fmt.Println("Writing")
	fmt.Println(sales)

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Error opening CSV file: %s", err)
	}
	// Ensure file closes at end.
	defer f.Close()
	w := csv.NewWriter(f)

	for _, sale := range sales {

		// Prepare data to be written to CSV.

		// Do not include deleted sales in reports.
		if sale.DeletedAt != nil {
			continue
		}

		// Do not include sales with status of "OPEN"
		if sale.Status != nil && *sale.Status == "OPEN" {
			continue
		}

		// Takes a Vend timestamp string as input and converts it to
		// a Go Time.time value.
		dtInLoc := vend.ParseVendDT(*sale.SaleDate, tz)

		// Time string with timezone removed.
		dtStr := dtInLoc.String()[0:19]

		// Split time and date on space.
		// Example date/time string: 2015-07-01 07:03:22
		var dateStr, timeStr string
		dateStr = dtStr[0:10]
		timeStr = dtStr[10:19]

		var invoiceNumber string
		if sale.InvoiceNumber != nil {
			invoiceNumber = *sale.InvoiceNumber
		}

		// TODO: Clean up.
		var customerName, customerFirstName, customerLastName,
			customerCompanyName, customerCode string
		var customerFullName []string
		for _, customer := range *customers {
			// Make sure we only use info from customer on our sale.
			if *customer.ID == *sale.CustomerID {
				if customer.FirstName != nil {
					customerFirstName = *customer.FirstName
					customerFullName = append(customerFullName, customerFirstName)
				}
				if customer.LastName != nil {
					customerLastName = *customer.LastName
					customerFullName = append(customerFullName, customerLastName)
				}
				if customer.Code != nil {
					customerCode = *customer.Code
				}
				if customer.CompanyName != nil {
					customerCompanyName = *customer.CompanyName
				}
				customerName = strings.Join(customerFullName, " ")
				break
			}
		}

		// Sale not wrapped in quote marks.
		var saleNote string
		if sale.Note != nil {
			saleNote = fmt.Sprintf("%q", *sale.Note)
		}

		// Add up the total quantities of each product line item.
		// TODO: Confirm sum of line items plus sale total..
		var totalQuantity, totalDiscount float64
		var saleItems []string
		for _, lineitem := range *sale.LineItems {
			totalQuantity += *lineitem.Quantity
			totalDiscount += *lineitem.DiscountTotal

			for _, product := range *products {
				if *product.ID == *lineitem.ProductID {
					var productItems []string
					productItems = append(productItems, fmt.Sprintf("%v", *lineitem.Quantity))
					productItems = append(productItems, fmt.Sprintf("%s", *product.Name))

					prodItem := strings.Join(productItems, " X ")
					saleItems = append(saleItems, fmt.Sprintf("%v", prodItem))
					break
				}
			}
		}
		totalQuantityStr := strconv.FormatFloat(totalQuantity, 'f', -1, 64)
		totalDiscountStr := strconv.FormatFloat(totalDiscount, 'f', -1, 64)
		// Show items sold separated by + sign.
		saleDetails := strings.Join(saleItems, " + ")

		// Sale subtotal.
		totalPrice := strconv.FormatFloat(*sale.TotalPrice, 'f', -1, 64)
		// Sale tax.
		totalTax := strconv.FormatFloat(*sale.TotalTax, 'f', -1, 64)
		// Sale total (subtotal plus tax).
		total := strconv.FormatFloat((*sale.TotalPrice + *sale.TotalTax), 'f', -1, 64)

		// Total loyalty on sale.
		totalLoyaltyStr := strconv.FormatFloat(*sale.TotalLoyalty, 'f', -1, 64)

		// saleDetails
		// TODO: Confirm what this is.

		var registerName string
		for _, register := range *registers {
			if *sale.RegisterID == *register.ID {
				registerName = *register.Name
				// Append (Deleted) to name if register is deleted.
				if register.DeletedAt != nil {
					registerName += " (Deleted)"
				}
				break
			} else {
				registerName = "<Deleted Register>"
			}
		}

		var userName string
		for _, user := range *users {
			if sale.UserID != nil && *sale.UserID == *user.ID {
				userName = *user.DisplayName
				break
			} else {
				userName = ""
			}
		}

		var saleStatus string
		if sale.Status != nil {
			saleStatus = *sale.Status
		}

		var record []string
		record = append(record, dateStr)             // Date
		record = append(record, timeStr)             // Time
		record = append(record, invoiceNumber)       // Receipt Number
		record = append(record, "Sale")              // Line Type
		record = append(record, customerCode)        // Customer Code
		record = append(record, customerCompanyName) // Customer Company Name
		record = append(record, customerName)        // Customer Name
		record = append(record, saleNote)            // Note
		record = append(record, totalQuantityStr)    // Quantity
		record = append(record, totalPrice)          // Subtotal
		record = append(record, totalTax)            // Sales Tax
		record = append(record, totalDiscountStr)    // Discount
		record = append(record, totalLoyaltyStr)     // Loyalty
		record = append(record, total)               // Sale total
		record = append(record, "")                  // Paid
		record = append(record, saleDetails)         // Details
		record = append(record, registerName)        // Register
		record = append(record, userName)            // User
		record = append(record, saleStatus)          // Status
		record = append(record, "")                  // Sku
		// record = append(record, "")                  // AccountCodeSale
		// record = append(record, "")                  // AccountCodePurchase
		w.Write(record)

		for _, lineitem := range *sale.LineItems {

			quantity := strconv.FormatFloat(*lineitem.Quantity, 'f', -1, 64)
			price := strconv.FormatFloat(*lineitem.Price, 'f', -1, 64)
			tax := strconv.FormatFloat(*lineitem.Tax, 'f', -1, 64)
			discount := strconv.FormatFloat(*lineitem.Discount, 'f', -1, 64)
			loyalty := strconv.FormatFloat(*lineitem.LoyaltyValue, 'f', -1, 64)
			total := strconv.FormatFloat(((*lineitem.Price + *lineitem.Tax) * *lineitem.Quantity), 'f', -1, 64)

			// TODO: Add supplier code?
			var productName, productSKU string
			for _, product := range *products {
				if *product.ID == *lineitem.ProductID {
					productName = *product.VariantName
					productSKU = *product.SKU
				}
			}

			productRecord := record
			productRecord[0] = dateStr      // Sale Date
			productRecord[1] = timeStr      // Sale Time
			productRecord[2] = ""           // Invoice Number
			productRecord[3] = "Sale Line"  // Line Type
			productRecord[4] = ""           // Customer Code
			productRecord[5] = ""           // Customer Company Name
			productRecord[6] = ""           // Customer Name
			productRecord[7] = ""           // TODO: line note from the product?
			productRecord[8] = quantity     // Quantity
			productRecord[9] = price        // Subtotal
			productRecord[10] = tax         // Sales Tax
			productRecord[11] = discount    // Discount
			productRecord[12] = loyalty     // Loyalty
			productRecord[13] = total       // Total
			productRecord[14] = ""          // Paid
			productRecord[15] = productName // Details
			productRecord[16] = ""          // Register
			productRecord[17] = ""          // User
			productRecord[18] = ""          // Status
			productRecord[19] = productSKU  // Sku
			// productRecord[19] = ""          // AccountCodeSale
			// productRecord[20] = ""          // AccountCodePurchase
			w.Write(productRecord)
		}

		payments := *sale.Payments
		for _, payment := range payments {

			paid := strconv.FormatFloat(*payment.Amount, 'f', -1, 64)
			name := fmt.Sprintf("%s", *payment.Name)
			// label := *payment.Label

			paymentRecord := record
			paymentRecord[0] = dateStr   // Sale Date
			paymentRecord[1] = timeStr   // Sale Time
			paymentRecord[2] = ""        // Invoice Number
			paymentRecord[3] = "Payment" // Line Type
			paymentRecord[4] = ""        // Customer Code
			paymentRecord[5] = ""        // Customer Company Name
			paymentRecord[6] = ""        // Customer Name
			paymentRecord[7] = ""        // TODO: line note
			paymentRecord[8] = ""        // Quantity
			paymentRecord[9] = ""        // Subtotal
			paymentRecord[10] = ""       // Sales Tax
			paymentRecord[11] = ""       // Discount
			paymentRecord[12] = ""       // Loyalty
			paymentRecord[13] = ""       // Total
			paymentRecord[14] = paid     // Paid
			paymentRecord[15] = name     //  Details
			paymentRecord[16] = ""       // Register
			paymentRecord[17] = ""       // User
			paymentRecord[18] = ""       // Status
			paymentRecord[19] = ""       // Sku
			// paymentRecord[19] = ""       // AccountCodeSale
			// paymentRecord[20] = ""       // AccountCodePurchase

			w.Write(paymentRecord)
		}
	}
	w.Flush()
}
