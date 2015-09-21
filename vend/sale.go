// Package vend interacts with the Vend API.
package vend

// SalePayload contains sales data and versioning info.
type SalePayload struct {
	Data    []Sale           `json:"data,omitempty"`
	Version map[string]int64 `json:"version,omitempty"`
}

// Sale is a basic sale object.
type Sale struct {
	ID        string     `json:"id,omitempty"`
	LineItems []LineItem `json:"line_items,omitempty"`
	Payments  []Payment  `json:"payments,omitempty"`
	Taxes     []Tax      `json:"taxes,omitempty"`
}

// LineItem is a product on a sale.
type LineItem struct {
	ID string `json:"id,omitempty"`
}

// Payment is a payment on a sale.
type Payment struct {
	ID string `json:"id,omitempty"`
}

// Tax is tax on a sale.
type Tax struct {
	ID string `json:"id,omitempty"`
}

/*
ENDPOINT:
.vendhq.com/api/2.0/sales

EXAMPLE PAYLOAD:
{
  "id": "b8ca3a65-0125-11e4-fbb5-71004ed35970",
  "outlet_id": "b8ca3a65-011c-11e4-fbb5-5973b0e19f1a",
  "register_id": "b8ca3a65-011c-11e4-fbb5-5973b0e372df",
  "user_id": "b8ca3a65-011c-11e4-fbb5-5973b0ee0c0d",
  "customer_id": "b8ca3a65-0125-11e4-fbb5-6fe038ec33f8",
  "invoice_number": "25",
  "receipt_number": "25",
  "invoice_sequence": 99,
  "receipt_sequence": 99,
  "status": "CLOSED",
  "note": null,
  "short_code": "oqnyct",
  "return_for": null,
  "created_at": "2014-11-20T21:58:06+00:00",
  "total_price": 1006.27,
  "total_loyalty": 0,
  "total_tax": 100.63,
  "updated_at": "2014-11-20T21:58:06+00:00",
  "sale_date": "2014-11-19T11:36:37+00:00",
  "deleted_at": null,
  "line_items": [
    {
      "id": "b8ca3a65-0125-11e4-fbb5-71004ee998b9",
      "product_id": "b8ca3a65-0125-11e4-fbb5-6fdfe30d4dd4",
      "quantity": 60,
      "price": 10.90909,
      "unit_price": 10.90909,
      "price_total": 654.5454,
      "total_price": 654.5454,
      "discount": 0,
      "unit_discount": 0,
      "discount_total": 0,
      "total_discount": 0,
      "loyalty_value": 0,
      "unit_loyalty_value": 0,
      "total_loyalty_value": 0,
      "cost": 1,
      "unit_cost": 1,
      "cost_total": 60,
      "total_cost": 60,
      "tax": 1.09091,
      "unit_tax": 1.09091,
      "tax_total": 65.4546,
      "total_tax": 65.4546,
      "tax_id": "b8ca3a65-0125-11e4-fbb5-5d6b28d9407e",
      "tax_components": [
        {
          "rate_id": "28d927ec-5d6b-11e4-9bb5-b8ca3a65011c",
          "total_tax": 65.45727
        }
      ],
      "price_set": false,
      "sequence": 0,
      "note": null,
      "status": "CONFIRMED",
      "is_return": false
    }
  ],
  "payments": [
    {
      "id": "b8ca3a65-0125-11e4-fbb5-71004eeb098b",
      "register_id": "b8ca3a65-011c-11e4-fbb5-5973b0e372df",
      "retailer_payment_type_id": "b8ca3a65-011c-11e4-fbb5-5973b0e3e4ba",
      "payment_type_id": "1",
      "name": "Cash",
      "payment_date": "2014-11-19T11:36:30+00:00",
      "amount": 300
    },
  ],
  "taxes": [
    {
      "id": "28d927ec-5d6b-11e4-9bb5-b8ca3a65011c",
      "amount": 100.62733
    }
  ],
  "version": 95313988
}
*/
