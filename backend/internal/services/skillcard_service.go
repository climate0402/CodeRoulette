package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

type SkillCardService struct {
	redis *redis.Client
}

type SkillCard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Cost        int    `json:"cost"`
	Rarity      string `json:"rarity"`
	Effect      string `json:"effect"`
}

type SkillCardUsage struct {
	CardID   string    `json:"card_id"`
	PlayerID string    `json:"player_id"`
	MatchID  string    `json:"match_id"`
	UsedAt   time.Time `json:"used_at"`
	Effect   string    `json:"effect"`
	Target   string    `json:"target,omitempty"`
}

// Available skill cards
var availableSkillCards = []SkillCard{
	{
		ID:          "peek_code",
		Name:        "Code Peek",
		Description: "View one line of your opponent's code",
		Type:        "peek",
		Cost:        1,
		Rarity:      "common",
		Effect:      "peek_line",
	},
	{
		ID:          "swap_test",
		Name:        "Test Swap",
		Description: "Swap one test case with your opponent",
		Type:        "swap",
		Cost:        2,
		Rarity:      "rare",
		Effect:      "swap_test_case",
	},
	{
		ID:          "hint",
		Name:        "Hint",
		Description: "Get a hint for the current problem",
		Type:        "hint",
		Cost:        1,
		Rarity:      "common",
		Effect:      "show_hint",
	},
	{
		ID:          "time_boost",
		Name:        "Time Boost",
		Description: "Get 30 seconds extra time",
		Type:        "boost",
		Cost:        2,
		Rarity:      "rare",
		Effect:      "add_time",
	},
	{
		ID:          "code_lock",
		Name:        "Code Lock",
		Description: "Lock opponent's code for 10 seconds",
		Type:        "lock",
		Cost:        3,
		Rarity:      "epic",
		Effect:      "lock_opponent",
	},
	{
		ID:          "perfect_score",
		Name:        "Perfect Score",
		Description: "Guarantee 100% score on next submission",
		Type:        "boost",
		Cost:        5,
		Rarity:      "legendary",
		Effect:      "perfect_submission",
	},
}

func NewSkillCardService(redis *redis.Client) *SkillCardService {
	return &SkillCardService{redis: redis}
}

// GetAvailableCards returns all available skill cards
func (s *SkillCardService) GetAvailableCards() []SkillCard {
	return availableSkillCards
}

// GetRandomCard returns a random skill card based on rarity
func (s *SkillCardService) GetRandomCard() SkillCard {
	rand.Seed(time.Now().UnixNano())

	// Weighted random selection based on rarity
	weights := map[string]int{
		"common":    50,
		"rare":      30,
		"epic":      15,
		"legendary": 5,
	}

	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	random := rand.Intn(totalWeight)
	currentWeight := 0

	var selectedRarity string
	for rarity, weight := range weights {
		currentWeight += weight
		if random < currentWeight {
			selectedRarity = rarity
			break
		}
	}

	// Get cards of selected rarity
	var cardsOfRarity []SkillCard
	for _, card := range availableSkillCards {
		if card.Rarity == selectedRarity {
			cardsOfRarity = append(cardsOfRarity, card)
		}
	}

	if len(cardsOfRarity) == 0 {
		// Fallback to common cards
		for _, card := range availableSkillCards {
			if card.Rarity == "common" {
				cardsOfRarity = append(cardsOfRarity, card)
			}
		}
	}

	return cardsOfRarity[rand.Intn(len(cardsOfRarity))]
}

// UseSkillCard uses a skill card in a match
func (s *SkillCardService) UseSkillCard(ctx context.Context, matchID, playerID, cardID string) (*SkillCardUsage, error) {
	// Find the card
	var card SkillCard
	for _, c := range availableSkillCards {
		if c.ID == cardID {
			card = c
			break
		}
	}

	if card.ID == "" {
		return nil, fmt.Errorf("skill card not found: %s", cardID)
	}

	// Create usage record
	usage := &SkillCardUsage{
		CardID:   cardID,
		PlayerID: playerID,
		MatchID:  matchID,
		UsedAt:   time.Now(),
		Effect:   card.Effect,
	}

	// Store usage in Redis
	usageKey := fmt.Sprintf("skill_usage:%s:%s", matchID, playerID)
	usageData, err := json.Marshal(usage)
	if err != nil {
		return nil, err
	}

	// Store with expiration (match duration)
	if err := s.redis.Set(ctx, usageKey, usageData, 30*time.Minute).Err(); err != nil {
		return nil, err
	}

	// Publish skill card usage event
	event := map[string]interface{}{
		"type":      "skill_card_used",
		"match_id":  matchID,
		"player_id": playerID,
		"card":      card,
		"usage":     usage,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	// Publish to match room
	roomKey := fmt.Sprintf("room:%s", matchID)
	s.redis.Publish(ctx, roomKey, eventData)

	return usage, nil
}

// GetPlayerCards returns skill cards available to a player
func (s *SkillCardService) GetPlayerCards(ctx context.Context, playerID string) ([]SkillCard, error) {
	cardsKey := fmt.Sprintf("player_cards:%s", playerID)

	// Get player's cards from Redis
	cardsData, err := s.redis.Get(ctx, cardsKey).Result()
	if err == redis.Nil {
		// Initialize with random cards
		cards := s.generateInitialCards()
		if err := s.setPlayerCards(ctx, playerID, cards); err != nil {
			return nil, err
		}
		return cards, nil
	} else if err != nil {
		return nil, err
	}

	var cards []SkillCard
	if err := json.Unmarshal([]byte(cardsData), &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

// generateInitialCards generates initial skill cards for a new player
func (s *SkillCardService) generateInitialCards() []SkillCard {
	cards := make([]SkillCard, 0, 5)

	// Give 3 common cards, 1 rare card, 1 epic card
	for i := 0; i < 3; i++ {
		cards = append(cards, s.getRandomCardByRarity("common"))
	}
	cards = append(cards, s.getRandomCardByRarity("rare"))
	cards = append(cards, s.getRandomCardByRarity("epic"))

	return cards
}

// getRandomCardByRarity returns a random card of specific rarity
func (s *SkillCardService) getRandomCardByRarity(rarity string) SkillCard {
	var cardsOfRarity []SkillCard
	for _, card := range availableSkillCards {
		if card.Rarity == rarity {
			cardsOfRarity = append(cardsOfRarity, card)
		}
	}

	if len(cardsOfRarity) == 0 {
		// Fallback to common
		for _, card := range availableSkillCards {
			if card.Rarity == "common" {
				cardsOfRarity = append(cardsOfRarity, card)
			}
		}
	}

	return cardsOfRarity[rand.Intn(len(cardsOfRarity))]
}

// setPlayerCards stores player's skill cards
func (s *SkillCardService) setPlayerCards(ctx context.Context, playerID string, cards []SkillCard) error {
	cardsKey := fmt.Sprintf("player_cards:%s", playerID)
	cardsData, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, cardsKey, cardsData, 0).Err() // No expiration
}

// AddCardToPlayer adds a skill card to player's collection
func (s *SkillCardService) AddCardToPlayer(ctx context.Context, playerID string, card SkillCard) error {
	cards, err := s.GetPlayerCards(ctx, playerID)
	if err != nil {
		return err
	}

	cards = append(cards, card)
	return s.setPlayerCards(ctx, playerID, cards)
}

// RemoveCardFromPlayer removes a skill card from player's collection
func (s *SkillCardService) RemoveCardFromPlayer(ctx context.Context, playerID, cardID string) error {
	cards, err := s.GetPlayerCards(ctx, playerID)
	if err != nil {
		return err
	}

	// Remove the card
	for i, card := range cards {
		if card.ID == cardID {
			cards = append(cards[:i], cards[i+1:]...)
			break
		}
	}

	return s.setPlayerCards(ctx, playerID, cards)
}

// GetMatchSkillUsage returns skill card usage for a match
func (s *SkillCardService) GetMatchSkillUsage(ctx context.Context, matchID string) ([]SkillCardUsage, error) {
	// Get all skill usage keys for this match
	pattern := fmt.Sprintf("skill_usage:%s:*", matchID)
	keys, err := s.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var usages []SkillCardUsage
	for _, key := range keys {
		usageData, err := s.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var usage SkillCardUsage
		if err := json.Unmarshal([]byte(usageData), &usage); err != nil {
			continue
		}

		usages = append(usages, usage)
	}

	return usages, nil
}
