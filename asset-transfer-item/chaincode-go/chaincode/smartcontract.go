package chaincode

import (
	"encoding/json"
	"fmt"


	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Item struct {
	ID             string `json:"ID"`
	Tipo		   string `json:"tipo"`
	Title          string `json:"title"`
	//Owner          string `json:"owner"`
	Date           string `json:"date"`
	Price		   int    `json:"price"`
	//Validate       bool   `json:"validate"`
}

type ItemCopy struct {
	Item             Item `json:"Item"`
	Owner          string `json:"owner"`
	Validate       bool   `json:"validate"`
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateItem(ctx contractapi.TransactionContextInterface, id string, tipo string, title string, owner string, date string, price int, validate bool) error {
	exists, err := s.itemExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the item %s already exists", id)
	}

	item := Item{
		ID:             id,
		Tipo:			tipo,
		Title:			title,
		Owner:			owner,
		Date:			date,
		Price:			price,
		Validate:		validate,
	}
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}

// ReadItem returns the item stored in the world state with given id.
func (s *SmartContract) ReadItem(ctx contractapi.TransactionContextInterface, id string) (*Item, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if itemJSON == nil {
		return nil, fmt.Errorf("the item %s does not exist", id)
	}

	var item Item
	err = json.Unmarshal(itemJSON, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}


// ItemExists returns true when asset with given ID exists in world state
func (s *SmartContract) itemExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return itemJSON != nil, nil
}

func (s *SmartContract) itemValidated(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	item, err := s.ReadItem(ctx, id)
	if err != nil {
		return false, err
	}

	/*itemJSON, err := json.Marshal(item)

	if err != nil {
		return false, err
	} */

	return item.Validate, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	item, err := s.ReadItem(ctx, id)
	if err != nil {
		return err
	}

	item.Owner = newOwner
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}

// ChangePrice updates the price of asset with given id in world state.
func (s *SmartContract) ChangePrice(ctx contractapi.TransactionContextInterface, id string, newPrice int) error {
	item, err := s.ReadItem(ctx, id)
	if err != nil {
		return err
	}
	/*
	var valid
	valid, _ = s.itemValidated(id)
	if !valid {
		return fmt.Errorf("The item %s has been deleted", id)
	}
	*/

	item.Price = newPrice
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}


// EliminateItem invelidates the asset with given id in world state.
func (s *SmartContract) EliminateItem(ctx contractapi.TransactionContextInterface, id string) error {
	item, err := s.ReadItem(ctx, id)
	if err != nil {
		return err
	}

	item.Validate = false
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}


// GetAllItems returns all items found in world state
func (s *SmartContract) GetAllItems(ctx contractapi.TransactionContextInterface) ([]*Item, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all items in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var items []*Item
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var item Item
		err = json.Unmarshal(queryResponse.Value, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}
