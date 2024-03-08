package model

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Game struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MaxPlayers int    `json:"max_players"`
}

type Room struct {
	ID       int    `json:"id"`
	RoomName string `json:"room_name"`
	IDGame   int    `json:"game_id"`
}

type Participant struct {
	ID        int `json:"id"`
	IDRoom    int `json:"room_id"`
	IDAccount int `json:"account_id"`
}

type ParticipantDetail struct {
	ID      int     `json:"id"`
	Account Account `json:""`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AllRoomResponse struct {
	Status int          `json:"status"`
	Data   []RoomDetail `json:"room"`
}

type RoomDetail struct {
	ID       int    `json:"id"`
	RoomName string `json:"room_name"`
}

type RoomDetailParticipant struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Username  string `json:"username"`
}

type RoomDetails struct {
	ID           int                     `json:"id"`
	RoomName     string                  `json:"room_name"`
	Participants []RoomDetailParticipant `json:"participants"`
}

type RoomDetailResponse2 struct {
	Room RoomDetails `json:"Room`
}
type RoomDetailResponse struct {
	Status int                 `json:"Status"`
	Data   RoomDetailResponse2 `json:"Data`
}
