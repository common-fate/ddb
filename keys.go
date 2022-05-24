package ddb

// Keys are primary and Global Secondary Index (GSI)
// keys to be used when storing an item in DynamoDB.
// The ddb package is opinionated on the naming of these keys.
type Keys struct {
	PK     string
	SK     string
	GSI1PK string `json:",omitempty"`
	GSI1SK string `json:",omitempty"`
	GSI2PK string `json:",omitempty"`
	GSI2SK string `json:",omitempty"`
	GSI3PK string `json:",omitempty"`
	GSI3SK string `json:",omitempty"`
}

// Keyers give DynamoDB keys to be used when inserting an item.
type Keyer interface {
	DDBKeys() (Keys, error)
}
