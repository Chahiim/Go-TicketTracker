package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/chahiim/ticket_tracker/internal/validator"
)

type Ticket struct {
	ID        int64
	CreatedAt time.Time
	CName     string
	IName     string
	Quantity  int
}

type TicketModel struct {
	DB *sql.DB
}

func ValidateTicket(v *validator.Validator, ticket *Ticket) {
	v.Check(validator.NotBlank(ticket.CName), "Customer Name", "must be provided")
	v.Check(validator.MaxLength(ticket.CName, 50), "Customer Name", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(ticket.IName), "Item Name", "must be provided")
	v.Check(validator.MaxLength(ticket.IName, 50), "Item Name", "must not be more than 50 bytes long")
	v.Check(ticket.Quantity > 0, "Quantity", "must be greater than zero")
	v.Check(ticket.Quantity <= 100, "Quantity", "must not exceed 100")
}

// Insert a ticket into the database
func (m *TicketModel) Insert(ticket *Ticket) error {
	query := `
		INSERT INTO tickets (cname, iname, quantity)
		VALUES ($1, $2, $3)
		RETURNING ticket_id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		ticket.CName,
		ticket.IName,
		ticket.Quantity,
	).Scan(&ticket.ID, &ticket.CreatedAt)
}

// Get all tickets from the database
func (m *TicketModel) GetAll() ([]*Ticket, error) {
	query := `
		SELECT ticket_id, created_at, cname, iname, quantity
		FROM tickets`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*Ticket
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(
			&ticket.ID,
			&ticket.CreatedAt,
			&ticket.CName,
			&ticket.IName,
			&ticket.Quantity,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)
	}
	return tickets, nil
}

// Delete a ticket
func (m *TicketModel) Delete(id int64) error {
	query := `DELETE FROM tickets WHERE ticket_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Update a ticket's information
func (m *TicketModel) Update(ticket *Ticket) error {
	query := `UPDATE tickets
			  SET cname = $1, iname = $2, quantity = $3
			  WHERE ticket_id = $4`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(
		ctx,
		query,
		ticket.CName,
		ticket.IName,
		ticket.Quantity,
	)
	return err
}

func (m *TicketModel) GetByID(id int64) (*Ticket, error) {
	query := `SELECT ticket_id, created_at, cname, iname, quantity FROM tickets WHERE ticket_id = $1`
	var ticket Ticket
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID, &ticket.CreatedAt, &ticket.CName, &ticket.IName, &ticket.Quantity)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
