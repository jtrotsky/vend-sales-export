package vend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client contains API authentication details.
type Client struct {
	Token        string
	DomainPrefix string
}

// NewClient is called to pass authentication details to the manager.
func NewClient(token, domainPrefix string) Client {
	return Client{token, domainPrefix}
}

// Sales grabs and collates all sales in pages of 10,000.
func (c Client) Sales() ([]Sale, error) {

	sales := []Sale{}
	s := []Sale{}
	var v int64

	// v is a version that is used to get products by page.
	// Here we get the first page.
	s, v, err := salePage(1, c.DomainPrefix, c.Token, "sales")
	sales = append(sales, s...)

	for len(s) > 0 {
		// Continue grabbing pages until we receive an empty one.
		s, v, err = salePage(v, c.DomainPrefix, c.Token, "sales")
		if err != nil {
			return nil, err
		}

		// Each period is a page of 10,000 sales.
		fmt.Printf(".")

		// Append sale page to list of sales.
		sales = append(sales, s...)
	}

	// Got no sales back, all done.
	fmt.Printf("\nGot no sales back, must have em' all.\n")

	return sales, err
}

func salePage(version int64, domainPrefix, key,
	resource string) ([]Sale, int64, error) {

	// Build the URL for the product page.
	url := urlFactory(version, domainPrefix, resource)

	body, err := urlGet(key, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

	// Decode the JSON into our defined product object.
	response := SalePayload{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling Vend sale payload: %s", err)
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

	fmt.Printf("\nGrabbing: %s", url)
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
	case 404:
		fmt.Printf("\nURL not found - check domain prefix. Status: %d",
			statusCode)
	case 429:
		fmt.Printf("\nRate limited by the Vend API :S Status: %d",
			statusCode)
	default:
		fmt.Printf("\nGot an unknown status code - Google it. Status: %d",
			statusCode)
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
