// Caminho: cmd/api/main.go

package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"Ultimo_trabalho_Go/internal/handler"
	"Ultimo_trabalho_Go/internal/repository"
	"Ultimo_trabalho_Go/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgresql://postgres:5000@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}

	// Repositórios
	playerRepository := repository.NewPlayerRepository(db)
	enemyRepository := repository.NewEnemyRepository(db)
	battleRepository := repository.NewBattleRepository(db)

	// Serviços
	playerService := service.NewPlayerService(*playerRepository)
	enemyService := service.NewEnemyService(*enemyRepository)
	battleService := service.NewBattleService(*playerRepository, *enemyRepository, *battleRepository)

	// Handlers
	playerHandler := handler.NewPlayerHandler(playerService)
	enemyHandler := handler.NewEnemyHandler(enemyService)
	battleHandler := handler.NewBattleHandler(battleService)

	r := mux.NewRouter()

	// Rotas para jogadores
	r.HandleFunc("/player", playerHandler.AddPlayer).Methods("POST")
	r.HandleFunc("/player", playerHandler.LoadPlayers).Methods("GET")
	r.HandleFunc("/player/{id}", playerHandler.LoadPlayer).Methods("GET")
	r.HandleFunc("/player/{id}", playerHandler.SavePlayer).Methods("PUT")
	r.HandleFunc("/player/{id}", playerHandler.DeletePlayer).Methods("DELETE")

	// Rotas para inimigos
	r.HandleFunc("/enemy", enemyHandler.AddEnemy).Methods("POST")
	r.HandleFunc("/enemy", enemyHandler.LoadEnemies).Methods("GET")
	r.HandleFunc("/enemy/{id}", enemyHandler.LoadEnemy).Methods("GET")
	r.HandleFunc("/enemy/{id}", enemyHandler.SaveEnemy).Methods("PUT")
	r.HandleFunc("/enemy/{id}", enemyHandler.DeleteEnemy).Methods("DELETE")

	// Rotas para batalhas
	r.HandleFunc("/battle", battleHandler.CreateBattle).Methods("POST")
	r.HandleFunc("/battle", battleHandler.LoadBattles).Methods("GET")
	r.HandleFunc("/battle/{id}", battleHandler.LoadBattle).Methods("GET")

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println(err)
	}
}
