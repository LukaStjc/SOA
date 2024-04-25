package models

type GraphDBUser struct {
	ID       int64 // mozda bude izazvalo gresku
	Username string
}

type GraphDBUsers []*GraphDBUser
