package sales

import (
	"encoding/json"

	"github.com/jtrotsky/govend/vend"
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

// SalePage grabs and collates all sales.
func (c Client) SalePage(v int64) (*[]vend.Sale, int64, error) {

	var s, sales []vend.Sale

	// v is a version that is used to objects by page.
	// Here we get the first page.
	data, v, err := vend.ResourcePage(v, c.DomainPrefix, c.Token, "sales")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &s)

	// Append page to list.
	sales = append(sales, s...)

	return &sales, v, err
}
