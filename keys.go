package ddb

// Keys are primary and Global Secondary Index (GSI)
// keys to be used when storing an item in DynamoDB.
// The ddb package is opinionated on the naming of these keys.
type Keys struct {
	PK     string
	SK     string
	GSI1PK string
	GSI1SK string
	GSI2PK string
	GSI2SK string
	GSI3PK string
	GSI3SK string
}

// Keyers give DynamoDB keys to be used when inserting an item.
type Keyer interface {
	DDBKeys() (Keys, error)
}
