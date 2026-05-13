package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the smart contract structure
type SmartContract struct {
	contractapi.Contract
}

// Actor represents a participant in the AgriTrace network
type Actor struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"` // Farmer, Aggregator, Processor, Regulator, Buyer
	Location  string `json:"location"`
	CreatedAt string `json:"createdAt"`
}

// Batch represents an agricultural product batch registered in the network
type Batch struct {
	BatchID   string `json:"batchId"`
	OwnerID   string `json:"ownerId"`
	CropType  string `json:"cropType"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"` // e.g., "REGISTERED", "HARVESTED"
	CreatedAt string `json:"createdAt"`
}

// ===================================================================================
// ACTOR MANAGEMENT
// ===================================================================================

// RegisterActor registers a new actor in the network
func (s *SmartContract) RegisterActor(ctx contractapi.TransactionContextInterface, id string, name string, role string, location string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the actor %s already exists", id)
	}

	actor := Actor{
		ID:        id,
		Name:      name,
		Role:      role,
		Location:  location,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	actorJSON, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	// Save to state
	return ctx.GetStub().PutState(id, actorJSON)
}

// GetActor retrieves an actor by ID
func (s *SmartContract) GetActor(ctx contractapi.TransactionContextInterface, id string) (*Actor, error) {
	actorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if actorJSON == nil {
		return nil, fmt.Errorf("the actor %s does not exist", id)
	}

	var actor Actor
	err = json.Unmarshal(actorJSON, &actor)
	if err != nil {
		return nil, err
	}

	return &actor, nil
}

// ===================================================================================
// BATCH MANAGEMENT
// ===================================================================================

// RegisterBatch registers a new batch of agricultural product
func (s *SmartContract) RegisterBatch(ctx contractapi.TransactionContextInterface, batchId string, ownerId string, cropType string, quantity int) error {
	exists, err := s.AssetExists(ctx, batchId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the batch %s already exists", batchId)
	}

	// Verify the owner exists
	ownerExists, err := s.AssetExists(ctx, ownerId)
	if err != nil || !ownerExists {
		return fmt.Errorf("owner actor %s does not exist", ownerId)
	}

	batch := Batch{
		BatchID:   batchId,
		OwnerID:   ownerId,
		CropType:  cropType,
		Quantity:  quantity,
		Status:    "REGISTERED",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	batchJSON, err := json.Marshal(batch)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(batchId, batchJSON)
}

// GetBatch retrieves a batch by ID
func (s *SmartContract) GetBatch(ctx contractapi.TransactionContextInterface, batchId string) (*Batch, error) {
	batchJSON, err := ctx.GetStub().GetState(batchId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if batchJSON == nil {
		return nil, fmt.Errorf("the batch %s does not exist", batchId)
	}

	var batch Batch
	err = json.Unmarshal(batchJSON, &batch)
	if err != nil {
		return nil, err
	}

	return &batch, nil
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
		log.Panicf("Error creating AgriTrace registration chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting AgriTrace registration chaincode: %v", err)
	}
}