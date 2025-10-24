package models

type Todo struct {
	ID          int    `json:"id"`
	Item        string `json:"item"`
	IsCompleted bool   `json:"isCompleted"`
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
