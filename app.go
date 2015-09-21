package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jtrotsky/go-vend/vce-dr/types"
	"github.com/vend/vce/tools"
)

func main() {
	// First, create a log file.
	lf, err := os.Create("./vend_sales_export.log")
	if err != nil {
		log.Fatalf("Error creating log file: %s", err)
	}
	defer lf.Close()
	log.SetOutput(lf)

	var domainPrefix, timeZone, token string
	flag.StringVar(&domainPrefix, "d", "tosca", "Retailer's Vend store name (x.vendhq.com).")
	flag.StringVar(&token, "token", "", "Retailer's personal API auth token.")
	flag.StringVar(&timeZone, "t", "Local", "Retailer TimeZone default is local.")
	flag.Parse()
	log.Printf("Arguments provided: -d %s -token %s -t %s", domainPrefix, token, timeZone)
	var fromDT, toDT string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter start date/time (YYYY-MM-DD HH:MM:SS): ")
	fromDT, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed reading user input for date from.")
	}
	fmt.Print("Enter end time (YYYY-MM-DD HH-MM-SS): ")
	toDT, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed reading user input for date from.")
	}

	log.Printf("Get sales from: %v", fromDT)
	log.Printf("To: %v", toDT)

	fromDateTime, err := time.Parse("2006-01-02 15:04:05", strings.Trim(fromDT, "\r\n"))
	if err != nil {
		log.Fatalf("Error parsing entered date/time: %s", err)
	}
	toDateTime, err := time.Parse("2006-01-02 15:04:05", strings.Trim(toDT, "\r\n"))
	if err != nil {
		log.Fatalf("Error parsing entered date/time: %s", err)
	}

	// Convert times to seconds-since-Unixtime for ease of comparison.
	fromDateTimeUnix := fromDateTime.Unix()
	toDateTimeUnix := toDateTime.Unix()

	DTfrom := fromDateTime.String()
	dateFromStr := DTfrom[0:10]

	log.Printf("Date from: %v", dateFromStr)

	fmt.Printf("\nGrabbing sales!\n")

	// First get all sales in the store.
	allSales, err := grabSales(dateFromStr, domainPrefix, token)
	if err != nil {
		log.Fatalf("Error getting sales: %s", err)
	}

	log.Printf("Returned %v sales in total. Now refining them down.", len(allSales))

	// Grab all registers and map them.
	allRegisters := getRegisters(domainPrefix, token)
	registerMap := buildRegisterMap(allRegisters)

	// Write sales to CSV file.
	fmt.Printf("\n\nWriting CSV :)\n")
	writeSales(allSales, timeZone, registerMap, fromDateTimeUnix, toDateTimeUnix)
}

func writeSales(sales []types.Sale, tz string, registerMap map[string]types.Register, fromDateTime, toDateTime int64) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	monthStr := strconv.Itoa(month)
	dayStr := strconv.Itoa(day)
	hourStr := strconv.Itoa(hour)
	minuteStr := strconv.Itoa(minute)
	secondStr := strconv.Itoa(second)

	monthStr = tools.PadLeft(monthStr, "0", 2)
	dayStr = tools.PadLeft(dayStr, "0", 2)
	hourStr = tools.PadLeft(hourStr, "0", 2)
	minuteStr = tools.PadLeft(minuteStr, "0", 2)
	secondStr = tools.PadLeft(secondStr, "0", 2)

	fname := fmt.Sprintf("vend_sales_history_%v_%v_%v_%v_%v_%v.csv", year, monthStr, dayStr, hourStr, minuteStr, secondStr)
	f, err := os.Create(fmt.Sprintf("./%s", fname))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer f.Close()

	loc, err := time.LoadLocation(tz)
	log.Printf("Time zone loaded from the argument: %s", loc)
	if err != nil {
		log.Fatalf("Error while loading time zone location: %v", err)
	}

	w := csv.NewWriter(f)

	var headerLine []string
	headerLine = append(headerLine, "Date")
	headerLine = append(headerLine, "Receipt Number")
	headerLine = append(headerLine, "Line Type")
	headerLine = append(headerLine, "Customer Code")
	headerLine = append(headerLine, "Customer Name")
	headerLine = append(headerLine, "Note")
	headerLine = append(headerLine, "Quantity")
	headerLine = append(headerLine, "Subtotal")
	headerLine = append(headerLine, "Sales Tax")
	headerLine = append(headerLine, "Discount")
	headerLine = append(headerLine, "Loyalty")
	headerLine = append(headerLine, "Total")
	headerLine = append(headerLine, "Paid")
	headerLine = append(headerLine, "Details")
	headerLine = append(headerLine, "Register")
	headerLine = append(headerLine, "User")
	headerLine = append(headerLine, "Status")
	headerLine = append(headerLine, "Sku")
	headerLine = append(headerLine, "AccountCodeSale")
	headerLine = append(headerLine, "AccountCodePurchase")

	w.Write(headerLine)

	for _, sale := range sales {

		const longForm = "2006-01-02 15:04:05"
		t, err := time.Parse(longForm, *sale.SaleDate)
		if err != nil {
			log.Fatalf("Error parsing time into deafult timestamp: %s", err)
			return
		}

		// Time in retailer's location.
		saleTimeTZ := t.In(loc)

		// Sale date in seconds since Jan 1st 1970 UTC.
		saleTime := saleTimeTZ.Unix()

		if fromDateTime > saleTime {
			continue
		}
		if toDateTime < saleTime {
			continue
		}

		// Time string with timezone correction.
		timeStr := saleTimeTZ.String()[0:19]

		invoiceNumber := sale.InvoiceNumber

		var customerCode string
		var customerName string
		customerName = sale.CustomerName
		customerCode = sale.Customer.Code

		// Format with quotes around text.
		saleNote := fmt.Sprintf("%q", sale.Note)

		var totalDiscount float64
		var totalLoyalty float64
		var saleItems []string
		var totalQuantity float64

		for _, product := range sale.Products {
			// TODO: Fix up discount total.
			totalDiscount += (product.Discount * product.Quantity)
			totalLoyalty += (product.Loyalty * product.Quantity)

			var productItems []string
			productItems = append(productItems, fmt.Sprintf("%v", product.Quantity))
			productItems = append(productItems, fmt.Sprintf("%v", product.Name))

			prodItem := strings.Join(productItems, " X ")
			saleItems = append(saleItems, fmt.Sprintf("%v", prodItem))
			totalQuantity += product.Quantity
		}
		totalQuantityStr := formatFloat(totalQuantity)
		totalDiscountStr := formatFloat(totalDiscount)
		totalLoyaltyStr := formatFloat(totalLoyalty)
		price := formatFloat(sale.Totals.Price)
		tax := formatFloat(sale.Totals.Tax)
		total := formatFloat(sale.Totals.Price + sale.Totals.Tax)

		// Show items sold separated by + sign.
		saleDetails := strings.Join(saleItems, " + ")

		var registerName string
		if rg, ok := registerMap[sale.RegisterID]; ok {
			registerName = rg.Name
		} else {
			registerName = fmt.Sprintf("Deleted register (%v)", sale.RegisterID)
		}

		saleUser := sale.UserName
		saleStatus := sale.Status

		var record []string
		record = append(record, timeStr)          // Date
		record = append(record, invoiceNumber)    // Receipt Number
		record = append(record, "Sale")           // Line Type
		record = append(record, customerCode)     // Customer Code
		record = append(record, customerName)     // Customer Name
		record = append(record, saleNote)         // Note
		record = append(record, totalQuantityStr) // Quantity
		record = append(record, price)            // Subtotal
		record = append(record, tax)              // Sales Tax
		record = append(record, totalDiscountStr) // Discount
		record = append(record, totalLoyaltyStr)  // Loyalty
		record = append(record, total)            // Total
		record = append(record, "")               // Paid
		record = append(record, saleDetails)      // Details
		record = append(record, registerName)     // Register
		record = append(record, saleUser)         // User
		record = append(record, saleStatus)       // Status
		record = append(record, "")               // Sku
		record = append(record, "")               // AccountCodeSale
		record = append(record, "")               // AccountCodePurchase

		w.Write(record)

		products := sale.Products
		for _, product := range products {

			quantity := formatFloat(product.Quantity)
			price := formatFloat(product.Price)
			tax := formatFloat(product.Tax)
			discount := formatFloat(product.Discount)
			loyalty := formatFloat(product.Loyalty)
			total := formatFloat((product.Price + product.Tax) * product.Quantity)
			name := product.Name
			sku := product.Sku

			productRecord := record
			productRecord[2] = "Sale Line" // Line Type
			productRecord[3] = ""          // Customer Code Code
			productRecord[4] = ""          // Customer Name Name
			productRecord[5] = ""          // Note TODO: line note from the product?
			productRecord[6] = quantity    // Quantity
			productRecord[7] = price       // Subtotal
			productRecord[8] = tax         // Sales Tax
			productRecord[9] = discount    // Discount
			productRecord[10] = loyalty    // Loyalty
			productRecord[11] = total      // Total
			productRecord[12] = ""         // Paid
			productRecord[13] = name       // Details
			productRecord[14] = ""         // Register
			productRecord[15] = ""         // User
			productRecord[16] = ""         // Status
			productRecord[17] = sku        // Sku
			productRecord[18] = ""         // AccountCodeSale
			productRecord[19] = ""         // AccountCodePurchase

			w.Write(productRecord)
		}

		payments := sale.Payments
		for _, payment := range payments {

			paid := formatFloat(payment.Amount)
			// name := fmt.Sprintf("%q", *payment.Name)
			label := payment.Label

			paymentRecord := record
			paymentRecord[2] = "Payment" // Line Type
			paymentRecord[3] = ""        // Customer Code Code
			paymentRecord[4] = ""        // Customer Name Name
			paymentRecord[5] = ""        // Note TODO: line note
			paymentRecord[6] = ""        // Quantity
			paymentRecord[7] = ""        // Subtotal
			paymentRecord[8] = ""        // Sales Tax
			paymentRecord[9] = ""        // Discount
			paymentRecord[10] = ""       // Loyalty
			paymentRecord[11] = ""       // Total
			paymentRecord[12] = paid     // Paid
			paymentRecord[13] = label    //  Details
			paymentRecord[14] = ""       // Register
			paymentRecord[15] = ""       // User
			paymentRecord[16] = ""       // Status
			paymentRecord[17] = ""       // Sku
			paymentRecord[18] = ""       // AccountCodeSale
			paymentRecord[19] = ""       // AccountCodePurchase

			w.Write(paymentRecord)
		}
	}

	w.Flush()
}

func buildRegisterMap(registers []types.Register) map[string]types.Register {
	registerMap := make(map[string]types.Register)

	for _, register := range registers {
		registerMap[register.ID] = register
	}

	return registerMap
}

func getRegisters(domainPrefix, token string) []types.Register {
	var registers []types.Register
	var payload types.RegisterPayload
	url := buildURL(domainPrefix, "registers")

	body, err := getResource(url, token)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	registers = payload.Registers

	return registers
}

func grabSales(dateFrom, domainPrefix, token string) ([]types.Sale, error) {

	sales := []types.Sale{}
	s := []types.Sale{}
	var page int64

	// seq is sequence/version, so that we can get products by page.
	// this gets the first page.
	s, page, err := grabSalePage(dateFrom, domainPrefix, token, 1)
	sales = append(sales, s...)

	for {
		// Now continue grabbing pages until we get an empty one.
		s, page, err = grabSalePage(dateFrom, domainPrefix, token, page)
		fmt.Printf(".")
		if err != nil {
			log.Fatalf("Error getting sales: %s", err)
			return nil, err
		}

		// If our payload isn't empty, keep going!
		if len(s) > 0 {
			log.Println("Appending page to list.")
			sales = append(sales, s...)
		} else {
			// Got no sales back, stopping.
			log.Println("Got no sales back, must have em' all.")
			log.Printf("Empty body: %s", s)
			break
		}
	}

	return sales, err
}

func grabSalePage(dateFrom, domainPrefix, token string,
	page int64) ([]types.Sale, int64, error) {

	client := &http.Client{}
	url := buildSalePageURL(page, dateFrom, domainPrefix)

	log.Printf("Getting sale page: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "Support-tool: vce-dr - one of JOEYM8's tools.")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request for sale page: %s", err)
		return nil, 0, err
	}
	log.Printf("Response status: %v", resp.StatusCode)
	if resp.StatusCode == 429 {
		log.Println("Got rate limited :S")
	}
	if resp.StatusCode == 401 {
		log.Println("Not Authenticated!")
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading response body: %s\n", err)
	}
	response := new(types.SalePayload)
	err = json.Unmarshal(body, response)
	if err != nil {
		log.Fatalf("Error unmarshalling response body: %s", err)
		return nil, 0, err
	}

	page = page + 1
	sales := response.Sales

	return sales, page, err
}

func getResource(url, token string) ([]byte, error) {
	log.Printf("Request url: %s\n", url)
	client := &http.Client{}
	var emptyByteSlice []byte

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Support-tool: vce-dr - one of JOEYM8's tools.")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatalf("Error: %s", err)
		return emptyByteSlice, err
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal("Response status: ", res.Status)
		return emptyByteSlice, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return emptyByteSlice, err
	}

	return body, nil
}

// Returns strings from float64s.
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 5, 32)
}

func buildURL(store, resource string) string {
	return fmt.Sprintf("https://%s.vendhq.com/api/%s", store, resource)
}

func buildSalePageURL(page int64, dateFrom, domainPrefix string) string {
	const pageSize = 200

	address := fmt.Sprintf("https://%s.vendhq.com/api/", domainPrefix)
	query := url.Values{}
	query.Add("page_size", fmt.Sprintf("%v", pageSize))
	query.Add("page", fmt.Sprintf("%v", page))
	query.Add("since", fmt.Sprintf("%s", dateFrom))
	address += fmt.Sprintf("register_sales?%s", query.Encode())

	return address
}
