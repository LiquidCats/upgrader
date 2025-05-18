package configs

import "github.com/LiquidCats/upgrader/internal/app/domain/entities"

type Chains map[entities.ChainName]Chain

type Chain struct {
	Topic string `json:"topic" envconfig:"TOPIC"`
	Path  string `json:"path" envconfig:"PATH"`
}
