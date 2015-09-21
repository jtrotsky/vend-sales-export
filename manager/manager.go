package manager

import (
	"fmt"
	"log"

	"github.com/jtrotsky/spate/vend"
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

	// Get all sales from the beginning of time.
	allSales, err := manager.vend.Sales()
	if err != nil {
		log.Fatalf("Failed to get sales: %s", err)
	}

	log.Println("FIN.")
	fmt.Printf("\n\n%v", allSales[0])
}
