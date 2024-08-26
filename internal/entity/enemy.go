package entity

import "github.com/google/uuid"

type Enemy struct {
	ID       string
	Nickname string
	Life     int
	Attack   int
	Armor    int // Alterado de "Defesa" para "Armor"
}

func NewEnemy(nickname string, life, attack, armor int) *Enemy { // Alterado "defesa" para "armor"
	return &Enemy{
		ID:       uuid.New().String(),
		Nickname: nickname,
		Life:     life,
		Attack:   attack,
		Armor:    armor, 
	}
}
