package manager

import (
	"fmt"
	"log"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/spate/vendapi"
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

	fmt.Printf("\n\nGrabbing sales.\n")
	// Version, to paginate.
	var v int64
	// Get first page.
	sales, v, err := vendapi.SalePage(v, manager.vend.DomainPrefix, manager.vend.Token)
	fname, err := writer.CreateReport(manager.vend.DomainPrefix)
	if err != nil {
		log.Fatalf("Failed writing sales to CSV: %s", err)
	}

	// Get and write remaining pages.
	for len(sales) > 0 {
		sales, v, err = vendapi.SalePage(v, manager.vend.DomainPrefix, manager.vend.Token)

		writer.WriteReport(fname, registers, users, customers, products,
			sales, manager.vend.DomainPrefix, manager.vend.TimeZone)
		if err != nil {
			log.Fatalf("Failed writing sales to CSV: %s", err)
		}
	}
}
