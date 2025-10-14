package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DIDRegistry represents the smart contract for managing DIDs
type DIDRegistry struct {
	contractapi.Contract
}

// DIDDocument represents a W3C compliant DID Document
type DIDDocument struct {
	Context           []string    `json:"@context"`
	ID               string      `json:"id"`
	Controller       string      `json:"controller"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
	Created          string      `json:"created"`
	Updated          string      `json:"updated"`
	Status           string      `json:"status"`
}

// VerificationMethod represents a public key entry in the DID Document
type VerificationMethod struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Controller   string `json:"controller"`
	PublicKeyJwk JWK    `json:"publicKeyJwk"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

// CreateDID creates a new DID Document for an IoT device
func (d *DIDRegistry) CreateDID(ctx contractapi.TransactionContextInterface, did string, controller string, publicKeyJwk string) error {
	// Check if DID already exists
	existing, err := ctx.GetStub().GetState(did)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("the DID %s already exists", did)
	}

	// Parse the public key JWK
	var jwk JWK
	err = json.Unmarshal([]byte(publicKeyJwk), &jwk)
	if err != nil {
		return fmt.Errorf("invalid JWK format: %v", err)
	}

	// Create verification method
	verMethod := VerificationMethod{
		ID:           did + "#keys-1",
		Type:         "JsonWebKey2020",
		Controller:   controller,
		PublicKeyJwk: jwk,
	}

	// Create DID Document
	now := time.Now().UTC().Format(time.RFC3339)
	didDoc := DIDDocument{
		Context:           []string{"https://www.w3.org/ns/did/v1"},
		ID:               did,
		Controller:       controller,
		VerificationMethod: []VerificationMethod{verMethod},
		Created:          now,
		Updated:          now,
	}

	// Convert to JSON and store
	didDocJSON, err := json.Marshal(didDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal DID Document: %v", err)
	}

	err = ctx.GetStub().PutState(did, didDocJSON)
	if err != nil {
		return fmt.Errorf("failed to put DID Document in world state: %v", err)
	}

	return nil
}

// ResolveDID retrieves a DID Document by its DID
func (d *DIDRegistry) ResolveDID(ctx contractapi.TransactionContextInterface, did string) (*DIDDocument, error) {
	didDocJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if didDocJSON == nil {
		return nil, fmt.Errorf("the DID %s does not exist", did)
	}

	var didDoc DIDDocument
	err = json.Unmarshal(didDocJSON, &didDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DID Document: %v", err)
	}

	return &didDoc, nil
}

// UpdateDID updates an existing DID Document
func (d *DIDRegistry) UpdateDID(ctx contractapi.TransactionContextInterface, did string, controller string, publicKeyJwk string) error {
	// Check if DID exists
	existing, err := ctx.GetStub().GetState(did)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("the DID %s does not exist", did)
	}

	var existingDoc DIDDocument
	err = json.Unmarshal(existing, &existingDoc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal existing DID Document: %v", err)
	}

	// Parse the new public key JWK
	var jwk JWK
	err = json.Unmarshal([]byte(publicKeyJwk), &jwk)
	if err != nil {
		return fmt.Errorf("invalid JWK format: %v", err)
	}

	// Update the DID Document
	verMethod := VerificationMethod{
		ID:           did + "#keys-1",
		Type:         "JsonWebKey2020",
		Controller:   controller,
		PublicKeyJwk: jwk,
	}

	existingDoc.Controller = controller
	existingDoc.VerificationMethod = []VerificationMethod{verMethod}
	existingDoc.Updated = time.Now().UTC().Format(time.RFC3339)

	// Convert to JSON and store
	didDocJSON, err := json.Marshal(existingDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal DID Document: %v", err)
	}

	err = ctx.GetStub().PutState(did, didDocJSON)
	if err != nil {
		return fmt.Errorf("failed to update DID Document in world state: %v", err)
	}

	return nil
}

// RevokeDID revokes a DID by setting its status to "revoked"
func (d *DIDRegistry) RevokeDID(ctx contractapi.TransactionContextInterface, did string) error {
	// Check if DID exists
	existing, err := ctx.GetStub().GetState(did)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("the DID %s does not exist", did)
	}

	var didDoc DIDDocument
	err = json.Unmarshal(existing, &didDoc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal DID Document: %v", err)
	}

	// Check if already revoked
	if didDoc.Status == "revoked" {
		return fmt.Errorf("the DID %s is already revoked", did)
	}

	// Update status to revoked
	didDoc.Status = "revoked"
	didDoc.Updated = time.Now().UTC().Format(time.RFC3339)

	// Store updated DID Document
	didDocJSON, err := json.Marshal(didDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal DID Document: %v", err)
	}

	err = ctx.GetStub().PutState(did, didDocJSON)
	if err != nil {
		return fmt.Errorf("failed to revoke DID Document in world state: %v", err)
	}

	return nil
}

// InitLedger initializes the ledger with sample DID data
func (d *DIDRegistry) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Sample DID for testing
	sampleDID := "did:fabric:test:device1"
	sampleController := "did:fabric:test:manufacturer1"
	sampleJWK := JWK{
		Kty: "EC",
		Crv: "P-256",
		X:   "sample-x-coordinate",
		Y:   "sample-y-coordinate",
	}

	verMethod := VerificationMethod{
		ID:           sampleDID + "#keys-1",
		Type:         "JsonWebKey2020",
		Controller:   sampleController,
		PublicKeyJwk: sampleJWK,
	}

	now := time.Now().UTC().Format(time.RFC3339)
	didDoc := DIDDocument{
		Context:           []string{"https://www.w3.org/ns/did/v1"},
		ID:               sampleDID,
		Controller:       sampleController,
		VerificationMethod: []VerificationMethod{verMethod},
		Created:          now,
		Updated:          now,
		Status:           "active",
	}

	didDocJSON, err := json.Marshal(didDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal sample DID Document: %v", err)
	}

	err = ctx.GetStub().PutState(sampleDID, didDocJSON)
	if err != nil {
		return fmt.Errorf("failed to put sample DID Document in world state: %v", err)
	}

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&DIDRegistry{})
	if err != nil {
		fmt.Printf("Error creating DID registry chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting DID registry chaincode: %v", err)
	}
}
