package entity

type PlayerName string

type Player struct {
	PlayerName PlayerName
	Address    string
	Color      string
}
