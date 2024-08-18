package service

import (
	"errors"

	"Ultimo_trabalho_Go/internal/entity"
	"Ultimo_trabalho_Go/internal/repository"
)

type BattleService struct {
	PlayerRepository repository.PlayerRepository
	EnemyRepository  repository.EnemyRepository
	BattleRepository repository.BattleRepository
}

func NewBattleService(playerRepo repository.PlayerRepository, enemyRepo repository.EnemyRepository, battleRepo repository.BattleRepository) *BattleService {
	return &BattleService{
		PlayerRepository: playerRepo,
		EnemyRepository:  enemyRepo,
		BattleRepository: battleRepo,
	}
}

func (bs *BattleService) CreateBattle(playerNickname, enemyNickname string) (*entity.Battle, error) {
	player, err := bs.PlayerRepository.LoadPlayerByNickname(playerNickname)
	if err != nil || player == nil {
		return nil, errors.New("player not found")
	}

	enemy, err := bs.EnemyRepository.LoadEnemyByNickname(enemyNickname)
	if err != nil || enemy == nil {
		return nil, errors.New("enemy not found")
	}

	if player.Life <= 0 || enemy.Life <= 0 {
		return nil, errors.New("both player and enemy must have life > 0 to battle")
	}

	// Gera um valor para DiceThrown (número do dado)
	battle := entity.NewBattle(player.ID, enemy.ID, player.Nickname, enemy.Nickname)
	dice := battle.DiceThrown

	if dice <= 3 {
		damage := enemy.Attack - player.Defesa
		if damage < 0 {
			damage = 0
		}
		player.Life -= damage
		if player.Life < 0 {
			player.Life = 0
		}
		if err := bs.PlayerRepository.SavePlayer(player.ID, player); err != nil {
			return nil, errors.New("failed to update player life")
		}
		battle.Result = "Enemy won"
	} else {
		damage := player.Attack - enemy.Defesa
		if damage < 0 {
			damage = 0
		}
		enemy.Life -= damage
		if enemy.Life < 0 {
			enemy.Life = 0
		}
		if err := bs.EnemyRepository.SaveEnemy(enemy.ID, enemy); err != nil {
			return nil, errors.New("failed to update enemy life")
		}
		battle.Result = "Player won"
	}

	// Insere a batalha no banco de dados
	if _, err := bs.BattleRepository.AddBattle(battle); err != nil {
		return nil, err
	}

	return battle, nil
}

// Função para carregar batalhas do banco de dados
func (bs *BattleService) LoadBattles() ([]*entity.Battle, error) {
	return bs.BattleRepository.LoadBattles()
}
