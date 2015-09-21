/*
   "customer": {
     "id": "02fb2229-0adc-11e3-a415-bc764e10976c",
     "name": "",
     "customer_code": "WALKIN",
     "customer_group_id": "02fa2d73-0adc-11e3-a415-bc764e10976c",
     "customer_group_name": "All Customers",
     "updated_at": "2013-11-14 20:41:04",
     "deleted_at": "",
     "balance": "0.000",
     "year_to_date": "5884.37012",
     "date_of_birth": "",
     "sex": "",
     "custom_field_1": "",
     "custom_field_2": "",
     "custom_field_3": "",
     "custom_field_4": "",
     "contact": {}
   },

*/

package types

type Customer struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"customer_code,omitempty"`
}
