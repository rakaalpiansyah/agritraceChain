package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the smart contract structure for LC Settlement
type SmartContract struct {
	contractapi.Contract
}

// LetterOfCredit represents a financial contract/settlement record between Buyer and Supplier
type LetterOfCredit struct {
	LCID         string `json:"lcId"`
	BuyerID      string `json:"buyerId"`
	SupplierID   string `json:"supplierId"` // Farmer, Aggregator, or Processor
	BatchID      string `json:"batchId"`
	Amount       int    `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"` // "ISSUED", "ACCEPTED", "SETTLED", "CANCELLED"
	IssueDate    string `json:"issueDate"`
	SettledDate  string `json:"settledDate"`
}

// ===================================================================================
// LC SETTLEMENT MANAGEMENT
// ===================================================================================

// IssueLC creates a new Letter of Credit (only Buyer can issue)
func (s *SmartContract) IssueLC(ctx contractapi.TransactionContextInterface, lcId string, buyerId string, supplierId string, batchId string, amount int, currency string) error {
	// Security Check: Only BuyerMSP can issue LC
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "BuyerMSP" {
		return fmt.Errorf("unauthorized: only Buyer can issue LC. Caller MSP: %s", clientMSPID)
	}

	exists, err := s.AssetExists(ctx, lcId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the LC %s already exists", lcId)
	}

	lc := LetterOfCredit{
		LCID:        lcId,
		BuyerID:     buyerId,
		SupplierID:  supplierId,
		BatchID:     batchId,
		Amount:      amount,
		Currency:    currency,
		Status:      "ISSUED",
		IssueDate:   time.Now().UTC().Format(time.RFC3339),
		SettledDate: "",
	}

	lcJSON, err := json.Marshal(lc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(lcId, lcJSON)
}

// SettleLC marks the LC as settled
func (s *SmartContract) SettleLC(ctx contractapi.TransactionContextInterface, lcId string) error {
	// Buyer or Regulator can settle/verify it
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "BuyerMSP" && clientMSPID != "RegulatorMSP" {
		return fmt.Errorf("unauthorized: only Buyer or Regulator can settle LC")
	}

	lcJSON, err := ctx.GetStub().GetState(lcId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if lcJSON == nil {
		return fmt.Errorf("the LC %s does not exist", lcId)
	}

	var lc LetterOfCredit
	err = json.Unmarshal(lcJSON, &lc)
	if err != nil {
		return err
	}

	if lc.Status == "SETTLED" {
		return fmt.Errorf("the LC %s is already settled", lcId)
	}

	lc.Status = "SETTLED"
	lc.SettledDate = time.Now().UTC().Format(time.RFC3339)

	lcJSON, err = json.Marshal(lc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(lcId, lcJSON)
}

// GetLC retrieves an LC by ID
func (s *SmartContract) GetLC(ctx contractapi.TransactionContextInterface, lcId string) (*LetterOfCredit, error) {
	lcJSON, err := ctx.GetStub().GetState(lcId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if lcJSON == nil {
		return nil, fmt.Errorf("the LC %s does not exist", lcId)
	}

	var lc LetterOfCredit
	err = json.Unmarshal(lcJSON, &lc)
	if err != nil {
		return nil, err
	}

	return &lc, nil
}

// ===================================================================================
// UTILITY FUNCTIONS
// ===================================================================================

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating AgriTrace LC settlement chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting AgriTrace LC settlement chaincode: %v", err)
	}
}
