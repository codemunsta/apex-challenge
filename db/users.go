package db

type User struct {
	Name     string
	Password string
	Account  int64
}

type Session struct {
	IUser    User
	IsActive bool
	Number   int8
	Games    []Game
}

type Game struct {
	Roll1 int8
	Roll2 int8
}

type Transaction struct {
	IUser  string
	Type   string
	Amount int64
}
