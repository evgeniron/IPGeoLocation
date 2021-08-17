package db

type Location struct {
	Country, City string
}

type DB interface {
	GetLocation(ip string) Location
}
