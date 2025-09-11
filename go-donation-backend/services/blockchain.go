package services

import (
	"fmt"
	"log"
)

func SendToBlockchain(donorID string, organizationID string, amount float64, currency string) (string, error) {
	// Simulate sending data to the blockchain
	log.Printf("Sending to blockchain: DonorID=%s, OrganizationID=%s, Amount=%.2f, Currency=%s", donorID, organizationID, amount, currency)
	transactionHash := fmt.Sprintf("0x%X", donorID)

	return transactionHash, nil
}
