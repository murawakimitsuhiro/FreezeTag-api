package models

type UserLocation struct {
	UserName string `json:"user_name"`
	PlayerKind string `json:"player_kind"`
	LocationLat	float64 `json:"location_lat"`
	LocationLng	float64 `json:"location_lng"`
	PlayerStatus string `json:"player_status"`
}
