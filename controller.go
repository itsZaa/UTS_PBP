package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "strconv"
	m "uts/model"
	// "github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id, room_name FROM rooms"

	name := r.URL.Query()["room_name"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE room_name='" + name[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		sendErrorResponse(w, "Query Error")
		return
	}

	var room m.RoomDetail
	var rooms []m.RoomDetail
	for rows.Next() {
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Failed to scan data")
			return
		} else {
			rooms = append(rooms, room)
		}
	}

	var response m.AllRoomResponse
	if len(rooms) > 0 {
		response.Status = 200
		response.Data = rooms
	} else {
		response.Status = 400
	}

	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}

func GetDetailRoom(w http.ResponseWriter, r *http.Request) {
	db, err := connectGorm()

	err = r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed to parse form data")
		return
	}

	roomID := r.URL.Query().Get("id")

	var room m.Room
	err = db.Where("id = ?", roomID).First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			sendErrorResponse(w, "Room not found")
		} else {
			sendErrorResponse(w, "Failed to retrieve room")
		}
		return
	}

	var participants []m.Participant
	err = db.Where("id_room = ?", roomID).Find(&participants).Error
	if err != nil {
		sendErrorResponse(w, "Failed to retrieve participants")
		return
	}

	var participantsResponse []m.RoomDetailParticipant
	for _, participant := range participants {
		var account m.Account
		err := db.Where("id = ?", participant.IDAccount).First(&account).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendErrorResponse(w, "Account not found")
			} else {
				sendErrorResponse(w, "Failed to retrieve account")
			}
			return
		}

		participantResponse := m.RoomDetailParticipant{
			ID:        participant.ID,
			AccountID: participant.IDAccount,
			Username:  account.Username,
		}

		participantsResponse = append(participantsResponse, participantResponse)
	}
	response := m.RoomDetailResponse{
		Status: 200,
		Data: m.RoomDetailResponse2{
			Room: m.RoomDetails{
				ID:           room.ID,
				RoomName:     room.RoomName,
				Participants: participantsResponse,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
}

func EnterRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed to parse form data")
		return
	}

	roomID := r.Form.Get("roomid")
	accountID := r.Form.Get("accountid")

	if roomID == "" || accountID == "" {
		sendErrorResponse(w, "Both roomid and accountid must be provided")
		return
	}

	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		sendErrorResponse(w, "Invalid roomid format")
		return
	}

	accountIDInt, err := strconv.Atoi(accountID)
	if err != nil {
		sendErrorResponse(w, "Invalid accountid format")
		return
	}

	db, err := connectGorm()

	var room m.Room
	err = db.First(&room, roomIDInt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			sendErrorResponse(w, "Room not found")
		} else {
			sendErrorResponse(w, "Failed to retrieve room")
		}
		return
	}

	// Check if the room has reached MaxPlayers
	var participantsCount int
	err = db.Model(&m.Participant{}).Where("room_id = ?", roomIDInt).Count(&participantsCount).Error
	if err != nil {
		sendErrorResponse(w, "Failed to count participants")
		return
	}

	// Insert participant into the room
	participant := m.Participant{
		IDRoom:    roomIDInt,
		IDAccount: accountIDInt,
	}

	err = db.Create(&participant).Error
	if err != nil {
		sendErrorResponse(w, "Failed to enter the room")
		return
	}

	// Return success response
	response := m.ErrorResponse{
		Status:  http.StatusOK,
		Message: "Entered the room successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed to leave room")
		return
	}

	// Extract roomID and accountID from the form data
	roomID := r.Form.Get("roomid")
	accountID := r.Form.Get("accountid")

	// Validate if roomID and accountID are provided
	if roomID == "" || accountID == "" {
		sendErrorResponse(w, "Failed to leave room")
		return
	}

	// Convert roomID and accountID to integers
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		sendErrorResponse(w, "Failed to leave room")
		return
	}

	accountIDInt, err := strconv.Atoi(accountID)
	if err != nil {
		sendErrorResponse(w, "Failed to leave room")
		return
	}

	// Connect to the database
	db, err := connectGorm()

	// Check if the participant exists in the room
	var participant m.Participant
	err = db.Where("room_id = ? AND account_id = ?", roomIDInt, accountIDInt).First(&participant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			sendErrorResponse(w, "Failed to leave room")
		} else {
			sendErrorResponse(w, "Failed to leave room")
		}
		return
	}

	// Delete the participant from the room
	err = db.Delete(&participant).Error
	if err != nil {
		sendErrorResponse(w, "Failed to leave room")
		return
	}

	// Return success response
	response := m.ErrorResponse{
		Status:  http.StatusOK,
		Message: "Left the room successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
