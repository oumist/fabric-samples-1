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
	ID           string `json:"ID"`
	Tipo         string `json:"tipo"`
	Title        string `json:"title"`
	CreationDate string `json:"creationdate"`
	Price        int    `json:"price"`
	Original     bool   `json:"original"`
}

type ItemCopy struct {
	IDcopy       string `json:"IDcopy"`
	IDoriginal   string `json:"Item"`
	Owner        string `json:"owner"`
	PurchaseDate string `json:"purchasedate"`
	Puntuation   int    `json:"puntuation"`
	Original     bool   `json:"original"`
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateItem(ctx contractapi.TransactionContextInterface, id string, tipo string, title string, date string, price int) error {
	exists, err := s.itemExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the item %s already exists", id)
	}

	item := Item{
		ID:           id,
		Tipo:         tipo,
		Title:        title,
		CreationDate: date,
		Price:        price,
		Original:     true,
	}
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, itemJSON)
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) Buy(ctx contractapi.TransactionContextInterface, idcopy string, id string, newOwner string, purchasedate string) error {
	//si no existe el item original, se genera un error
	exists, err := s.itemExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the item %s does not exist", id)
	}
	//hay que asegurarse de que lo que se compra es un item original y no una copia
	item, err3 := s.ReadItem(ctx, id)
	if err3 != nil {
		return err
	}
	if !item.Original {
		return fmt.Errorf("the item %s is not an original item", id)
	}
	//si ya existe el item copia, se genera un error
	exists2, err2 := s.itemExists(ctx, idcopy)
	if err2 != nil {
		return err
	}
	if exists2 {
		return fmt.Errorf("the copy %s already exists", idcopy)
	}

	itemCopy := ItemCopy{
		IDcopy:       idcopy,
		IDoriginal:   id,
		Owner:        newOwner,
		PurchaseDate: purchasedate,
		Original:     false,
	}

	itemJSON, err := json.Marshal(itemCopy)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(idcopy, itemJSON)
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) Rate(ctx contractapi.TransactionContextInterface, idcopy string, newPuntuation int) error {
	itemCopy, err := s.ReadCopy(ctx, idcopy)
	if err != nil {
		return err
	}

	if !itemCopy.Original {
		itemCopy.Puntuation = newPuntuation

	} else {
		return fmt.Errorf("the item %s is not a copy", idcopy)

	}

	itemCopyJSON, err := json.Marshal(itemCopy)
	if err != nil {
		return fmt.Errorf("failed to read the copy from world state: %v", err)
	}

	return ctx.GetStub().PutState(idcopy, itemCopyJSON)
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

	if !item.Original {
		return nil, fmt.Errorf("the item found is a copy, the requested must be original")
	}

	return &item, nil
}

// ReadItem returns the item stored in the world state with given id.
func (s *SmartContract) ReadCopy(ctx contractapi.TransactionContextInterface, idcopy string) (*ItemCopy, error) {
	itemCopyJSON, err := ctx.GetStub().GetState(idcopy)
	if err != nil {
		return nil, fmt.Errorf("failed to read the copy from world state: %v", err)
	}
	if itemCopyJSON == nil {
		return nil, fmt.Errorf("the copy of the item %s does not exist", idcopy)
	}

	var itemCopy ItemCopy
	err = json.Unmarshal(itemCopyJSON, &itemCopy)
	if err != nil {
		return nil, err
	}

	if itemCopy.Original {
		return nil, fmt.Errorf("the item %s found is original, the requested must be a copy", idcopy)
	}

	return &itemCopy, nil
}

// Return returns the item stored in the world state with given id.
func (s *SmartContract) Return(ctx contractapi.TransactionContextInterface, idcopy string, motive int) error {
	//como haces readcopy no se debería de comprobar si es orginal o copia
	itemCopy, err := s.ReadCopy(ctx, idcopy)
	if err != nil {
		return err
	}
	//suponiendo que las opciones 0 y 1 sean las válidas SE PUEDE CAMBIAR
	//no sé si se puede borrar de esa manera
	if !itemCopy.Original {
		if motive < 2 {
			fmt.Println("the item has been returned")
			return ctx.GetStub().DelState(idcopy)
		} else {
			return fmt.Errorf("it is not possible to return the item %s for that motive", idcopy)
		}
	} else {
		return fmt.Errorf("the item %s is not a copy", idcopy)
	}

}

// ItemExists returns true when asset with given ID exists in world state
func (s *SmartContract) itemExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return itemJSON != nil, nil
}

// ChangePrice updates the price of asset with given id in world state.
func (s *SmartContract) ChangePrice(ctx contractapi.TransactionContextInterface, id string, newPrice int) error {
	item, err := s.ReadItem(ctx, id)
	if err != nil {
		return err
	}
	if !item.Original {
		return fmt.Errorf("the item %s is a copy and should be original", id)
	}
	item.Price = newPrice
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
		if item.Original {
			items = append(items, &item)
		}

	}

	return items, nil
}

// GetAllCopies returns all copies found in world state
func (s *SmartContract) GetAllCopies(ctx contractapi.TransactionContextInterface) ([]*ItemCopy, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all items in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var itemsCopy []*ItemCopy
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var itemCopy ItemCopy
		err = json.Unmarshal(queryResponse.Value, &itemCopy)
		if err != nil {
			return nil, err
		}
		if !itemCopy.Original {
			itemsCopy = append(itemsCopy, &itemCopy)
		}
	}

	return itemsCopy, nil
}
