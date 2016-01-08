package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/spate/writer"
)

// Manager contains the Vend client.
type Manager struct {
	vend vend.Client
}

// NewManager creates an instance of manager.
func NewManager(vend vend.Client) *Manager {
	return &Manager{vend}
}

// Run executes the process of grabbing sales then writing them to CSV.
func (manager *Manager) Run() {
	// Using log gives us an opening timestamp.
	log.Printf("BEGIN\n")

	fmt.Printf("\nGrabbing registers.\n")
	// Get registers.
	registers, err := manager.vend.Registers()
	if err != nil {
		log.Fatalf("Failed to get registers: %s", err)
	}

	fmt.Printf("\n\nGrabbing users.\n")
	// Get users.
	users, err := manager.vend.Users()
	if err != nil {
		log.Fatalf("Failed to get users: %s", err)
	}

	fmt.Printf("\n\nGrabbing customers.\n")
	// Get customers.
	customers, err := manager.vend.Customers()
	if err != nil {
		log.Fatalf("Failed to get customers: %s", err)
	}

	fmt.Printf("\n\nGrabbing products.\n")
	// Get all products from the beginning of time.
	products, _, err := manager.vend.Products()
	if err != nil {
		log.Fatalf("Failed to get products: %s", err)
	}

	// Create template report to be written to.
	file, err := writer.CreateReport(manager.vend.DomainPrefix)
	if err != nil {
		log.Fatalf("Failed writing sales to CSV: %s", err)
	}
	// Make sure file is closed at end.
	defer file.Close()

	fmt.Printf("\n\nGrabbing and writing sales.\n")
	// Version number, to paginate.
	var v int64
	// Sale object to unmarshal raw JSON into.
	sales := []vend.Sale{}
	// Get first page of sales
	rawSalePage, v, err := vend.ResourcePage(v, manager.vend.DomainPrefix, manager.vend.Token,
		"sales")
	// Unmarshal payload into sales object.
	if err = json.Unmarshal(rawSalePage, &sales); err != nil {
		fmt.Printf("Error unmarshelling sale JSON: %v", err)
	}

	// Get and write remaining pages if we got any sales from the first page.
	if len(rawSalePage) > 2 {
		fmt.Println("Got page, writing.")
		// Write first sale page.
		file = writer.WriteReport(file, registers, users, customers, products, sales,
			manager.vend.DomainPrefix, manager.vend.TimeZone)
		// Get and write remaining pages.
		for len(rawSalePage) > 2 {
			sales = []vend.Sale{}

			// Continue grabbing pages until we receive an empty one.
			rawSalePage, v, err = vend.ResourcePage(v, manager.vend.DomainPrefix,
				manager.vend.Token, "sales")
			if err != nil {
				fmt.Printf("Error getting sale page: %v", err)
			}

			// Unmarshal payload into sales object.
			if err = json.Unmarshal(rawSalePage, &sales); err != nil {
				fmt.Printf("Error unmarshelling sale JSON: %v", err)
			}

			// No point trying to write when response is empty.
			if len(rawSalePage) > 2 {

				fmt.Println("Got page, writing.")
				file = writer.WriteReport(file, registers, users, customers, products, sales,
					manager.vend.DomainPrefix, manager.vend.TimeZone)
			} else {
				fmt.Println("No results back.")
				break
			}
		}
	} else {
		fmt.Println("No results back.")
		// Remove template CSV file as it would be empty anyway.
		file.Close()
		os.Remove(file.Name())
	}
	// Using log gives us a closing timestamp.
	log.Println("")
	log.Printf("FIN\n")
}
