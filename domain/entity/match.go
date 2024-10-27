package entity

type Match struct {
	TurnOwner PlayerName
	Players   map[PlayerName]Player
	Actions   []Action
	Board     Board
	Chat      []MessageAction
	Winner    *PlayerName
}
