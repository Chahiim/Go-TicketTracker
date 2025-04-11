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
	Quantity  string
}

func ValidateTicket(v *validator.Validator, ticket *Ticket) {
	v.Check(validator.NotBlank(ticket.CName), "customername", "must be provided")
	v.Check(validator.MaxLength(ticket.CName, 50), "Customer Name", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(ticket.IName), "Item Name", "must be provided")
	v.Check(validator.MaxLength(ticket.IName, 50), "Item Name", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(ticket.Quantity), "Quantity", "must be provided")
	v.Check(validator.MaxLength(ticket.Quantity, 25), "Quantity", "must not be more than 25 bytes")
}

type TicketModel struct {
	DB *sql.DB
}

func (m *TicketModel) Insert(ticket *Ticket) error {
	query := `
		INSERT INTO ticket (cname, iname, quantity)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

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

func (m *TicketModel) ReadAll(ticket *Ticket) error {
	query := `
		SELECT id, created_at, cname, iname, quantity
		FROM ticket
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (m *TicketModel) Update(ticket *Ticket) error {
	query := `
		UPDATE ticket
		SET cname = $1, iname = $2, quantity = $3
		WHERE id = $4
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(
		ctx,
		query,
		ticket.CName,
		ticket.IName,
		ticket.Quantity,
		ticket.ID,
	).Scan(&ticket.ID, &ticket.CreatedAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m *TicketModel) Delete(ticket *Ticket) error {
	query := `
		DELETE FROM ticket
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, ticket.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
