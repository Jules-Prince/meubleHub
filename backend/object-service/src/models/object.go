package models

type Object struct {
	ID          string `json:"id"`          // Unique identifier
	Name        string `json:"name"`        // Name of the object
	Type        string `json:"type"`        // Type of the object
	IsReserved  bool   `json:"isReserved"`  // Indicates if the object is reserved
	ReservedBy  string `json:"reservedBy"`  // User who reserved the object
    RoomID      string `json:"room_id"`     // ID of the room this object belongs to
}
