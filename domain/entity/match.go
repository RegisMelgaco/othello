package entity

type Match struct {
	TurnOwner PlayerName
	Players   []PlayerName
	Actions   []Action
	Board     *Board
	Chat      []MessageAction
	Winner    PlayerName
}
