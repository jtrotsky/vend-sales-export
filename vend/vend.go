package vend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Client contains API authentication details.
type Client struct {
	Token        string
	DomainPrefix string
	TimeZone     string
}

// NewClient is called to pass authentication details to the manager.
func NewClient(token, domainPrefix, tz string) Client {
	return Client{token, domainPrefix, tz}
}

// Sales grabs and collates all sales in pages of 10,000.
func (c Client) Sales() (*[]Sale, error) {

	sales := []Sale{}
	s := []Sale{}
	data := []byte{}
	var v int64

	// v is a version that is used to objects by page.
	// Here we get the first page.
	data, v, err := resourcePage(1, c.DomainPrefix, c.Token, "sales")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &s)

	sales = append(sales, s...)

	for len(s) > 0 {
<<<<<<< HEAD
		// Continue  pages until we receive an empty one.
		s, v, err = salePage(v, c.DomainPrefix, c.Token, "sales")
=======
		// Continue grabbing pages until we receive an empty one.
		data, v, err = resourcePage(v, c.DomainPrefix, c.Token, "sales")
>>>>>>> consolidate-getting-resource-pages
		if err != nil {
			return nil, err
		}

<<<<<<< HEAD
=======
		// Unmarshal payload into sales object.
		err = json.Unmarshal(data, &s)

>>>>>>> consolidate-getting-resource-pages
		// Append sale page to list of sales.
		sales = append(sales, s...)
	}

	return &sales, err
}

func resourcePage(version int64, domainPrefix, key,
	resource string) ([]byte, int64, error) {

	// Build the URL for the product page.
	url := urlFactory(version, domainPrefix, resource)

	body, err := urlGet(key, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

<<<<<<< HEAD
	// Decode the JSON into our defined object.
	response := SalePayload{}
=======
	// Decode the raw JSON.
	response := Payload{}
>>>>>>> consolidate-getting-resource-pages
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling payload: %s", err)
		return nil, 0, err
	}

	// Data is an array of product objects.
	data := response.Data

	// The product version is a sequence number on each product object. Knowing
	// the highest number means we can continue grabbing results that are
	// after that number until we have all of the products.
	version = response.Version["max"]

	return data, version, err
}

// Registers gets all registers from a store.
func (c Client) Registers() (*[]Register, error) {

	// Build the URL for the register page.
	url := urlFactory(0, c.DomainPrefix, "registers")

	body, err := urlGet(c.Token, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

	// Decode the JSON into our defined product object.
	response := RegisterPayload{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling Vend sale payload: %s", err)
		return &[]Register{}, err
	}

	// Data is an array of product objects.
	data := response.Data

	// Do not expect more than one page of registers.
	// TODO: Consider including check for multiple pages.
	// version = response.Version["max"]

	return &data, err
}

// Users gets all users from a store.
func (c Client) Users() (*[]User, error) {

	// Build the URL for the register page.
	url := urlFactory(0, c.DomainPrefix, "users")

	body, err := urlGet(c.Token, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

	// Decode the JSON into our defined product object.
	response := UserPayload{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling Vend sale payload: %s", err)
		return &[]User{}, err
	}

	// Data is an array of product objects.
	data := response.Data

	// Do not expect more than one page of registers.
	// TODO: Consider including check for multiple pages.
	// version = response.Version["max"]

	return &data, err
}

// Customers grabs and collates all customers in pages of 10,000.
func (c Client) Customers() (*[]Customer, error) {

	customers := []Customer{}
	cp := []Customer{}
<<<<<<< HEAD
=======
	data := []byte{}
>>>>>>> consolidate-getting-resource-pages
	var v int64

	// v is a version that is used to get customers by page.
	// Here we get the first page.
<<<<<<< HEAD
	cp, v, err := customerPage(1, c.DomainPrefix, c.Token, "customers")
	customers = append(customers, cp...)

	for len(cp) > 0 {
		// Continue grabbing pages until we receive an empty one.
		cp, v, err = customerPage(v, c.DomainPrefix, c.Token, "customers")
=======
	data, v, err := resourcePage(1, c.DomainPrefix, c.Token, "customers")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &cp)

	customers = append(customers, cp...)

	for len(cp) > 0 {
		// Continue grabbing pages until we receive an empty one.
		data, v, err = resourcePage(v, c.DomainPrefix, c.Token, "customers")
>>>>>>> consolidate-getting-resource-pages
		if err != nil {
			return nil, err
		}

<<<<<<< HEAD
=======
		// Unmarshal payload into sales object.
		err = json.Unmarshal(data, &cp)

>>>>>>> consolidate-getting-resource-pages
		// Append customer page to list of customers.
		customers = append(customers, cp...)
	}

	return &customers, err
}

<<<<<<< HEAD
func customerPage(version int64, domainPrefix, key,
	resource string) ([]Customer, int64, error) {

	// Build the URL for the customer page.
	url := urlFactory(version, domainPrefix, resource)

	body, err := urlGet(key, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

	// Decode the JSON into our defined product object.
	response := CustomerPayload{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling Vend customer payload: %s", err)
		return nil, 0, err
	}

	// Data is an array of product objects.
	data := response.Data

	// The customer version is a sequence number on each customer object. Knowing
	// the highest number means we can continue grabbing results that are
	// after that number until we have all of the customers.
	version = response.Version["max"]

	return data, version, err
}

// Products grabs and collates all products in pages of 10,000.
func (c Client) Products() (*[]Product, error) {

	products := []Product{}
	p := []Product{}
	var v int64

	// v is a version that is used to get products by page.
	// Here we get the first page.
	p, v, err := productPage(1, c.DomainPrefix, c.Token, "products")
	products = append(products, p...)

	for len(p) > 0 {
		// Continue grabbing pages until we receive an empty one.
		p, v, err = productPage(v, c.DomainPrefix, c.Token, "products")
		if err != nil {
			return nil, err
		}

		// Append page to list.
		products = append(products, p...)
	}

	return &products, err
}

func productPage(version int64, domainPrefix, key,
	resource string) ([]Product, int64, error) {

	// Build the URL for the product page.
	url := urlFactory(version, domainPrefix, resource)

	body, err := urlGet(key, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

	// Decode the JSON into our defined object.
	response := ProductPayload{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling payload: %s", err)
		return nil, 0, err
	}

	// Data is an array of product objects.
	data := response.Data

	// The product version is a sequence number on each product object. Knowing
	// the highest number means we can continue grabbing results that are
	// after that number until we have all of the products.
	version = response.Version["max"]

	return data, version, err
}

=======
>>>>>>> consolidate-getting-resource-pages
// urlGet performs a basic get request on a url with Vend API authentication.
func urlGet(key, url string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("\nError creating http request: %s", err)
		return nil, err
	}

	// Using personal token authentication for the Vend API.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	req.Header.Set("User-Agent", "Support-tool: spate")

	log.Printf("Grabbing: %s\n", url)
	// Doing the request.
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nError performing request: %s", err)
		return nil, err
	}
	// Make sure response body is closed at end.
	defer res.Body.Close()

	ResponseCheck(res.StatusCode)

	// Read what we got back.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("\nError while reading response body: %s\n", err)
		return nil, err
	}

	return body, err
}

// ResponseCheck checks the HTTP status codes returned from Vend.
func ResponseCheck(statusCode int) {

	// Check HTTP response status codes.
	switch statusCode {
	case 200:
	// 	Response is bueno.
	case 401:
		fmt.Printf("\nAccess denied - check personal API token. Status: %d",
			statusCode)
		os.Exit(0)
	case 404:
		fmt.Printf("\nURL not found - check domain prefix. Status: %d",
			statusCode)
		os.Exit(0)
	case 429:
		fmt.Printf("\nRate limited by the Vend API :S Status: %d",
			statusCode)
	default:
		fmt.Printf("\nGot an unknown status code - Google it. Status: %d",
			statusCode)
		os.Exit(0)
	}
}

// urlFactory creates a Vend API 2.0 URL based on a resource.
func urlFactory(version int64, domainPrefix, resource string) string {
	// Page size is capped at ten thousand.
	// TODO: check if deleted is working for 2.0 yet.
	const (
		pageSize = 10000
		deleted  = false
	)

	// Using 2.x Endpoint.
	address := fmt.Sprintf("https://%s.vendhq.com/api/2.0/", domainPrefix)
	query := url.Values{}
	query.Add("deleted", fmt.Sprintf("%t", deleted))
	query.Add("page_size", fmt.Sprintf("%d", pageSize))
	query.Add("after", fmt.Sprintf("%d", version))

	address += fmt.Sprintf("%s?%s", resource, query.Encode())
	return address
}
