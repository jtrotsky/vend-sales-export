package writer

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jtrotsky/govend/client"
	"github.com/jtrotsky/govend/vend"
)

// SalesReport aims to mimic the report generated by
// exporting Vend sales history.
func SalesReport(registers *[]vend.Register, users *[]vend.User,
	customers *[]vend.Customer, products *[]vend.Product,
	sales *[]vend.Sale, domainPrefix, tz string) error {

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

	for _, sale := range *sales {

		// Prepare data to be written to CSV.

		// Do not include deleted sales in reports.
		if sale.DeletedAt != nil {
			continue
		}

		// Do not include sales with status of "OPEN"
		if *sale.Status == "OPEN" {
			continue
		}

		// Takes a Vend timestamp string as input and converts it to
		// a Go Time.time value.
		dtInLoc := client.ParseVendDT(*sale.SaleDate, tz)

		// Time string with timezone removed.
		dtStr := dtInLoc.String()[0:19]

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

		// TODO: Check how to see deleted registers.
		// should use deleted user names anyway, but with brackets maybe?
		var registerName string
		for _, register := range *registers {
			if *sale.RegisterID == *register.ID {
				registerName = *register.Name
			} else {
				registerName = "<Deleted Register>"
			}
		}

		var userName string
		for _, user := range *users {
			if *sale.UserID == *user.ID {
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
		record = append(record, dtStr)               // Date
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
			productRecord[0] = dtStr        // Sale Date
			productRecord[1] = ""           // Invoice Number
			productRecord[2] = "Sale Line"  // Line Type
			productRecord[3] = ""           // Customer Code
			productRecord[4] = ""           // Customer Company Name
			productRecord[5] = ""           // Customer Name
			productRecord[6] = ""           // TODO: line note from the product?
			productRecord[7] = quantity     // Quantity
			productRecord[8] = price        // Subtotal
			productRecord[9] = tax          // Sales Tax
			productRecord[10] = discount    // Discount
			productRecord[11] = loyalty     // Loyalty
			productRecord[12] = total       // Total
			productRecord[13] = ""          // Paid
			productRecord[14] = productName // Details
			productRecord[15] = ""          // Register
			productRecord[16] = ""          // User
			productRecord[17] = ""          // Status
			productRecord[18] = productSKU  // Sku
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
			paymentRecord[0] = dtStr     // Sale Date
			paymentRecord[1] = ""        // Invoice Number
			paymentRecord[2] = "Payment" // Line Type
			paymentRecord[3] = ""        // Customer Code
			paymentRecord[4] = ""        // Customer Company Name
			paymentRecord[5] = ""        // Customer Name
			paymentRecord[6] = ""        // TODO: line note
			paymentRecord[7] = ""        // Quantity
			paymentRecord[8] = ""        // Subtotal
			paymentRecord[9] = ""        // Sales Tax
			paymentRecord[10] = ""       // Discount
			paymentRecord[11] = ""       // Loyalty
			paymentRecord[12] = ""       // Total
			paymentRecord[13] = paid     // Paid
			paymentRecord[14] = name     //  Details
			paymentRecord[15] = ""       // Register
			paymentRecord[16] = ""       // User
			paymentRecord[17] = ""       // Status
			paymentRecord[18] = ""       // Sku
			// paymentRecord[19] = ""       // AccountCodeSale
			// paymentRecord[20] = ""       // AccountCodePurchase

			w.Write(paymentRecord)
		}
	}
	w.Flush()

	return err
}
