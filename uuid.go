package konnex

type IDprovider interface {
	ID() (string, error)
}
