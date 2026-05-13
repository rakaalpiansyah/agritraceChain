package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the smart contract structure for Certification
type SmartContract struct {
	contractapi.Contract
}

// Certificate represents a quality/compliance certificate issued to a batch
type Certificate struct {
	CertID      string `json:"certId"`
	IssuerID    string `json:"issuerId"`
	BatchID     string `json:"batchId"`
	CertType    string `json:"certType"` // e.g., "Organic", "FairTrade", "GAP"
	IssueDate   string `json:"issueDate"`
	ValidUntil  string `json:"validUntil"`
	Status      string `json:"status"` // "ACTIVE", "REVOKED", "EXPIRED"
}

// ===================================================================================
// CERTIFICATION MANAGEMENT
// ===================================================================================

// IssueCertificate creates a new certificate for a specific agricultural batch
func (s *SmartContract) IssueCertificate(ctx contractapi.TransactionContextInterface, certId string, issuerId string, batchId string, certType string, validUntil string) error {
	// Security Check: Ensure only RegulatorMSP can issue certificates
	// This makes the Smart Contract highly secure and realistic for the paper
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "RegulatorMSP" {
		return fmt.Errorf("unauthorized: only Regulator can issue certificates. Caller MSP: %s", clientMSPID)
	}

	exists, err := s.AssetExists(ctx, certId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the certificate %s already exists", certId)
	}

	cert := Certificate{
		CertID:     certId,
		IssuerID:   issuerId,
		BatchID:    batchId,
		CertType:   certType,
		IssueDate:  time.Now().UTC().Format(time.RFC3339),
		ValidUntil: validUntil,
		Status:     "ACTIVE",
	}

	certJSON, err := json.Marshal(cert)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(certId, certJSON)
}

// GetCertificate retrieves a certificate by ID
func (s *SmartContract) GetCertificate(ctx contractapi.TransactionContextInterface, certId string) (*Certificate, error) {
	certJSON, err := ctx.GetStub().GetState(certId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if certJSON == nil {
		return nil, fmt.Errorf("the certificate %s does not exist", certId)
	}

	var cert Certificate
	err = json.Unmarshal(certJSON, &cert)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

// RevokeCertificate changes the status of a certificate to REVOKED
func (s *SmartContract) RevokeCertificate(ctx contractapi.TransactionContextInterface, certId string) error {
	// Security Check: Only RegulatorMSP can revoke
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "RegulatorMSP" {
		return fmt.Errorf("unauthorized: only Regulator can revoke certificates")
	}

	cert, err := s.GetCertificate(ctx, certId)
	if err != nil {
		return err
	}

	cert.Status = "REVOKED"

	certJSON, err := json.Marshal(cert)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(certId, certJSON)
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
		log.Panicf("Error creating AgriTrace certification chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting AgriTrace certification chaincode: %v", err)
	}
}
