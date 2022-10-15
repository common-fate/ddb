package ddb

// If EntityType is implemented, ddb will write a special
// 'ddb:type' when marshalling the object.
// This can be used to implement custom unmarshalling for
// queries which return multiple object types.
//
// For example, you may wish to save an invoice along with
// its line items as separate rows in DynamoDB.
// The `EntityType` of the invoice can be "invoice",
// and then `EntityType` of the line item can be "lineItem".
// When querying the database and unmarshalling these objects
// back into Go structs, you can check the type of them
// by looking at the value of 'ddb:type'.
type EntityTyper interface {
	EntityType() string
}
