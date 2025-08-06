package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DEALERID    string `json:"DEALERID"`
	MSISDN      string `json:"MSISDN"`
	MPIN        string `json:"MPIN"`
	BALANCE     string `json:"BALANCE"`
	STATUS      string `json:"STATUS"`
	TRANSAMOUNT string `json:"TRANSAMOUNT"`
	TRANSTYPE   string `json:"TRANSTYPE"`
	REMARKS     string `json:"REMARKS"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DEALERID: "dealer1", MSISDN: "1234567890", MPIN: "1234", BALANCE: "1000", STATUS: "active", TRANSAMOUNT: "0", TRANSTYPE: "init", REMARKS: "Initial asset"},
		{DEALERID: "dealer2", MSISDN: "0987654321", MPIN: "4321", BALANCE: "2000", STATUS: "active", TRANSAMOUNT: "0", TRANSTYPE: "init", REMARKS: "Initial asset"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.DEALERID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance string, status string, transAmount string, transType string, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", dealerID)
	}

	asset := Asset{
		DEALERID:    dealerID,
		MSISDN:      msisdn,
		MPIN:        mpin,
		BALANCE:     balance,
		STATUS:      status,
		TRANSAMOUNT: transAmount,
		TRANSTYPE:   transType,
		REMARKS:     remarks,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerID, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerID string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", dealerID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance string, status string, transAmount string, transType string, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dealerID)
	}

	// overwriting original asset with new asset
	asset := Asset{
		DEALERID:    dealerID,
		MSISDN:      msisdn,
		MPIN:        mpin,
		BALANCE:     balance,
		STATUS:      status,
		TRANSAMOUNT: transAmount,
		TRANSTYPE:   transType,
		REMARKS:     remarks,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerID, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, dealerID string) error {
	exists, err := s.AssetExists(ctx, dealerID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dealerID)
	}

	return ctx.GetStub().DelState(dealerID)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dealerID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-tracker chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-tracker chaincode: %v", err)
	}
}
