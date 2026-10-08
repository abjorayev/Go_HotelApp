package entity

import "time"

type Hotel struct { 
	ID   int   `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	Stars   int    `json:"stars"`
	ManagerID int   `json:"manager_id"`
	IsActive  bool   `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}
