package models

import (
	"time"

	"example.com/app/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (evt *Event) Save() error {
	query := `
		INSERT INTO Events(
			name, description, location, dateTime, user_id
		) VALUES (
			?, ?, ?, ?, ?
		)
	`

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()
	result, err := statement.Exec(evt.Name, evt.Description, evt.Location, evt.DateTime, evt.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	evt.ID = id

	return err
}

func GetEventByID(id int64) (*Event, error) {
	query := `
		SELECT * FROM Events WHERE id=?
	`

	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func GetAllEvents() ([]Event, error) {
	query := `
		SELECT * FROM Events
	`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events = []Event{}

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (event Event) Update() error {
	query := `
		UPDATE Events
		SET name=?, description=?, location=?, dateTime=?
		WHERE id=?
	`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := `
		DELETE FROM Events
		WHERE id=?
	`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID)
	return err
}

func (event *Event) Register(userID int64) error {
	query := `
		INSERT INTO Registrations (event_id, user_id)
		VALUES (?, ?)
	`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID, userID)
	return err
}

func (event *Event) CancelRegistration(userID int64) error {
	query := `
		DELETE FROM Registrations
		WHERE event_id=? AND user_id=?
	`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID, userID)
	return err
}
