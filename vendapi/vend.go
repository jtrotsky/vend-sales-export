package vendapi

import (
	"encoding/json"
	"fmt"

	"github.com/jtrotsky/govend/vend"
)

// SalePage grabs and collates all sales.
func SalePage(v int64, domainPrefix, token string) ([]vend.Sale, int64, error) {

	var s, sales []vend.Sale

	// v is a version that is used to objects by page.
	// Here we get the first page.
	data, v, err := vend.ResourcePage(v, domainPrefix, token, "sales")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &s)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s", err)
	}

	// Append page to list.
	sales = append(sales, s...)

	return sales, v, err
}
