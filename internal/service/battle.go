package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

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

// Critico
func (bs *BattleService) isCriticalHit() bool {
	return rand.Intn(100) < 20
}

func (bs *BattleService) calculateDamage(attack, armor int) int {
	baseDamage := attack - armor/2
	if baseDamage < 1 {
		baseDamage = 1
	}
	return baseDamage
}

func (bs *BattleService) CreateBattle(playerNickname, enemyNickname string) (*entity.Battle, string, error) {
	player, err := bs.PlayerRepository.LoadPlayerByNickname(playerNickname)
	if err != nil || player == nil {
		return nil, "", errors.New("jogador não encontrado")
	}

	enemy, err := bs.EnemyRepository.LoadEnemyByNickname(enemyNickname)
	if err != nil || enemy == nil {
		return nil, "", errors.New("inimigo não encontrado")
	}

	// Logs para diagnóstico
	fmt.Printf("Player: %s | Attack: %d | Armor: %d | Life: %d\n", player.Nickname, player.Attack, player.Armor, player.Life)
	fmt.Printf("Enemy: %s | Attack: %d | Armor: %d | Life: %d\n", enemy.Nickname, enemy.Attack, enemy.Armor, enemy.Life)

	if player.Life <= 0 || enemy.Life <= 0 {
		return nil, "", errors.New("tanto o jogador quanto o inimigo devem ter vida > 0 para batalhar")
	}

	battle := entity.NewBattle(player.ID, enemy.ID, player.Nickname, enemy.Nickname)
	dice := battle.DiceThrown

	var result string
	var damage int

	if dice <= 3 {

		damage = bs.calculateDamage(enemy.Attack, player.Armor)
		if bs.isCriticalHit() {
			damage *= 2
			fmt.Println("Crítico!")
		}
		player.Life -= damage
		if player.Life < 0 {
			player.Life = 0
		}
		if err := bs.PlayerRepository.SavePlayer(player.ID, player); err != nil {
			return nil, "", errors.New("falha ao atualizar a vida do jogador")
		}
		result = "Inimigo atacou. Dano causado: " + strconv.Itoa(damage) + ". Vida restante do jogador: " + strconv.Itoa(player.Life)
	} else {

		damage = bs.calculateDamage(player.Attack, enemy.Armor)
		if bs.isCriticalHit() {
			damage *= 2
			fmt.Println("Crítico! Dano do jogador foi dobrado.")
		}
		enemy.Life -= damage
		if enemy.Life < 0 {
			enemy.Life = 0
		}
		if err := bs.EnemyRepository.SaveEnemy(enemy.ID, enemy); err != nil {
			return nil, "", errors.New("falha ao atualizar a vida do inimigo")
		}
		result = "Jogador atacou. Dano causado: " + strconv.Itoa(damage) + ". Vida restante do inimigo: " + strconv.Itoa(enemy.Life)
	}

	if player.Life == 0 {
		battle.Result = "Inimigo venceu"
		result = "Inimigo venceu a batalha"
	} else if enemy.Life == 0 {
		battle.Result = "Jogador venceu"
		result = "Jogador venceu a batalha"
	} else {
		battle.Result = "A batalha continua"
	}

	if _, err := bs.BattleRepository.AddBattle(battle); err != nil {
		return nil, "", err
	}

	return battle, result, nil
}

func (bs *BattleService) LoadBattles() ([]*entity.Battle, error) {
	return bs.BattleRepository.LoadBattles()
}

func (bs *BattleService) LoadBattle(id string) (*entity.Battle, error) {
	return bs.BattleRepository.LoadBattleById(id)
}
