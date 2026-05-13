package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the smart contract structure for Compliance
type SmartContract struct {
	contractapi.Contract
}

// ComplianceRecord represents a traceability or quality check record during distribution
type ComplianceRecord struct {
	RecordID     string `json:"recordId"`
	BatchID      string `json:"batchId"`
	ReporterID   string `json:"reporterId"`
	Stage        string `json:"stage"` // "COLLECTION", "PROCESSING", "PACKAGING", "SHIPPING"
	Status       string `json:"status"` // "COMPLIANT", "NON_COMPLIANT"
	Details      string `json:"details"`
	ReportedAt   string `json:"reportedAt"`
}

// ===================================================================================
// COMPLIANCE MANAGEMENT
// ===================================================================================

// RecordCompliance creates a new compliance/traceability record for a batch
func (s *SmartContract) RecordCompliance(ctx contractapi.TransactionContextInterface, recordId string, batchId string, reporterId string, stage string, status string, details string) error {
	// Security Check: Only AggregatorMSP or ProcessorMSP can record compliance
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "AggregatorMSP" && clientMSPID != "ProcessorMSP" {
		return fmt.Errorf("unauthorized: only Aggregator or Processor can record compliance. Caller MSP: %s", clientMSPID)
	}

	exists, err := s.AssetExists(ctx, recordId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the compliance record %s already exists", recordId)
	}

	record := ComplianceRecord{
		RecordID:   recordId,
		BatchID:    batchId,
		ReporterID: reporterId,
		Stage:      stage,
		Status:     status,
		Details:    details,
		ReportedAt: time.Now().UTC().Format(time.RFC3339),
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(recordId, recordJSON)
}

// GetComplianceRecord retrieves a record by ID
func (s *SmartContract) GetComplianceRecord(ctx contractapi.TransactionContextInterface, recordId string) (*ComplianceRecord, error) {
	recordJSON, err := ctx.GetStub().GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if recordJSON == nil {
		return nil, fmt.Errorf("the compliance record %s does not exist", recordId)
	}

	var record ComplianceRecord
	err = json.Unmarshal(recordJSON, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
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
		log.Panicf("Error creating AgriTrace compliance chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting AgriTrace compliance chaincode: %v", err)
	}
}
