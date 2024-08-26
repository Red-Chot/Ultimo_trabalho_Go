package handler

import (
	"encoding/json"
	"net/http"

	"Ultimo_trabalho_Go/internal/entity"
	"Ultimo_trabalho_Go/internal/service"

	"github.com/gorilla/mux"
)

type BattleHandler struct {
	BattleService *service.BattleService
}

func NewBattleHandler(battleService *service.BattleService) *BattleHandler {
	return &BattleHandler{BattleService: battleService}
}

func (bh *BattleHandler) CreateBattle(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Player string `json:"Player"`
		Enemy  string `json:"Enemy"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	battle, result, err := bh.BattleService.CreateBattle(request.Player, request.Enemy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Battle *entity.Battle `json:"battle"`
		Result string         `json:"result"`
	}{
		Battle: battle,
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Método para carregar todas as batalhas
func (bh *BattleHandler) LoadBattles(w http.ResponseWriter, r *http.Request) {
	battles, err := bh.BattleService.LoadBattles()
	if err != nil {
		http.Error(w, "Error loading battles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(battles); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// Método para carregar uma batalha específica por ID
func (bh *BattleHandler) LoadBattle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	battle, err := bh.BattleService.LoadBattle(id)
	if err != nil {
		http.Error(w, "Error loading battle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(battle); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
