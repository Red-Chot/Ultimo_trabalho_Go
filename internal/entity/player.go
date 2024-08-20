package entity

import "github.com/google/uuid"

type Player struct {
	ID       string
	Nickname string
	Life     int
	Attack   int
	Armor    int // Alterado de "Defesa" para "Armor"
}

func NewPlayer(nickname string, life, attack, armor int) *Player { // Alterado "defesa" para "armor"
	return &Player{
		ID:       uuid.New().String(),
		Nickname: nickname,
		Life:     life,
		Attack:   attack,
		Armor:    armor, // Alterado de "defesa" para "armor"
	}
}
