package manager

import (
	"fmt"
	"log"

	"github.com/jtrotsky/spate/vend"
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
	products, err := manager.vend.Products()
	if err != nil {
		log.Fatalf("Failed to get products: %s", err)
	}

	fmt.Printf("\n\nGrabbing sales.\n")
	// Get all sales from the beginning of time.
	sales, err := manager.vend.Sales()
	if err != nil {
		log.Fatalf("Failed to get sales: %s", err)
	}

	fmt.Println("")
	log.Println("FIN.")

	fmt.Printf("\nGot %d sales.\n", len(sales))
	fmt.Println("Writing sales to CSV.")

	err = writer.SalesReport(registers, users, customers, products,
		sales, manager.vend.DomainPrefix, manager.vend.TimeZone)
	if err != nil {
		log.Fatalf("Failed writing sales to CSV: %s", err)
	}
}
