/*

{
  "id": "a2961c07-c48b-93c0-11e4-f4386fa3a759",
  "register_id": "b8ca3a65-0125-11e4-fbb5-a67748404883",
  "market_id": "1",
  "customer_id": "e52b2846-e925-11e4-f98b-f47bd1b36da8",
  "customer_name": "Brucy Dubbo",
customer []
  "user_id": "b8ca3a65-0125-11e4-fbb5-a5f35dbcb9e8",
  "user_name": "vend",
  "sale_date": "2015-05-07 05:42:19",
  "created_at": "2015-05-07 05:42:22",
  "updated_at": "2015-05-07 05:46:43",
  "total_price": 39.08,
  "total_cost": 5,
  "total_tax": 2.74,
  "tax_name": "VAT",
  "note": "",
  "status": "LAYBY_CLOSED",
  "short_code": "k0g60z",
  "invoice_number": "90",
  "return_for": "",
  "register_sale_products": [
    {
      "id": "a2961c07-c48b-ac15-11e4-f46e8469293f",
      "product_id": "b8ca3a65-0125-11e4-fbb5-9aad09f44068",
      "register_id": "b8ca3a65-0125-11e4-fbb5-a67748404883",
      "sequence": "0",
      "handle": "name3-2",
      "sku": "20057",
      "name": "Dark Mulch / KG",
      "quantity": 1,
      "price": 39.08174,
      "cost": 5,
      "price_set": 0,
      "discount": 0,
      "loyalty_value": 41.81746,
      "tax": 2.73572,
      "tax_id": "b8ca3a65-0125-11e4-fbb5-af1fcdf1b2ac",
      "tax_name": "VAT",
      "tax_rate": 0.07,
      "tax_total": 2.73572,
      "price_total": 39.08174,
      "display_retail_price_tax_inclusive": "1",
      "status": "CONFIRMED",
      "attributes": [
        {
          "name": "line_note",
          "value": ""
        }
      ]
    }
  ],
  "totals": {
    "total_tax": 2.74,
    "total_price": 39.08,
    "total_payment": 41.82,
    "total_to_pay": 0
  },
  "register_sale_payments": [
    {
      "id": "e52b2846-e925-11e4-f98b-f47c706bd5b7",
      "payment_type_id": "1",
      "register_id": "b8ca3a65-0125-11e4-fbb5-a67748404883",
      "retailer_payment_type_id": "b8ca3a65-011c-11e4-fbb5-5973b0e3e4ba",
      "name": "Cash",
      "label": "Cash",
      "payment_date": "2015-05-06 00:00:00",
      "amount": 41.82
    }
  ],
  "taxes": [
    {
      "id": "cdf19579-af1f-11e4-9bb5-b8ca3a65011c",
      "tax": 2.73572,
      "name": "VAT",
      "rate": 0.07
    }
  ]
}
*/

// Package types ...
package types

// Sale ...
type Sale struct {
	ID            string     `json:"id"`
	RegisterID    string     `json:"register_id"`
	MarketID      string     `json:"market_id,omitempty"`
	CustomerID    string     `json:"customer_id,omitempty"`
	CustomerName  string     `json:"customer_name,omitempty"`
	UserID        string     `json:"user_id,omitempty"`
	UserName      string     `json:"user_name,omitempty"`
	SaleDate      *string    `json:"sale_date,omitempty"`  // String here but converted to time.Time later.
	CreatedAt     string     `json:"created_at"`           // String here but converted to time.Time later.
	UpdatedAt     string     `json:"updated_at,omitempty"` // String here but converted to time.Time later.
	TotalPrice    float64    `json:"total_price,omitempty"`
	TotalCost     float64    `json:"total_cost,omitempty"`
	TotalTax      float64    `json:"total_tax,omitempty"`
	TaxName       string     `json:"tax_name,omitempty"`
	Note          string     `json:"note,omitempty"`
	Status        string     `json:"status,omitempty"`
	ShortCode     string     `json:"short_code,omitempty"`
	InvoiceNumber string     `json:"invoice_number,omitempty"`
	ReturnFor     string     `json:"return_for,omitempty"`
	Products      []Products `json:"register_sale_products,omitempty"`
	Payments      []Payments `json:"register_sale_payments,omitempty"`
	Customer      Customer   `json:"customer,omitempty"`
	Totals        Totals     `json:"totals,omitempty"`
	// Taxes                []Taxes
}

/*
"customer": {
	"id": "e52b2846-e925-11e4-f98b-f47bd1b36da8",
	"name": "Brucy Dubbo",
	"customer_code": "Brucy-T589",
	"customer_group_id": "b8ca3a65-011c-11e4-fbb5-5973b0cb3205",
	"customer_group_name": "All Customers",
	"first_name": "Brucy",
	"last_name": "Dubbo",
	"company_name": "",
	"phone": "",
	"mobile": "",
	"fax": "",
	"email": "",
	"twitter": "",
	"website": "",
	"physical_address1": "",
	"physical_address2": "",
	"physical_suburb": "",
	"physical_city": "",
	"physical_postcode": "",
	"physical_state": "",
	"physical_country_id": "NZ",
	"postal_address1": "",
	"postal_address2": "",
	"postal_suburb": "",
	"postal_city": "",
	"postal_postcode": "",
	"postal_state": "",
	"postal_country_id": "",
	"enable_loyalty": 1,
	"loyalty_balance": "41.81746",
	"updated_at": "2015-05-07 05:46:43",
	"deleted_at": "",
	"balance": "0.000",
	"year_to_date": "41.82000",
	"date_of_birth": "",
	"sex": "",
	"custom_field_1": "",
	"custom_field_2": "",
	"custom_field_3": "",
	"custom_field_4": "",
	"note": "",
	"contact": {
		"company_name": "",
		"phone": "",
		"email": ""
	}
}
*/
