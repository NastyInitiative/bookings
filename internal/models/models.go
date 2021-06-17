package models

import "time"

// Users - user data model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms - rooms data model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions - restrictions data model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservations - reservations data model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	RoomID    int
	Room      Room
}

// RoomRestrictions - single room restriction data model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	Reservation   Reservation
	Restriction   Restriction
	Room          Room
}
