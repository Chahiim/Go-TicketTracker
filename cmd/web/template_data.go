package main

import "github.com/chahiim/ticket_tracker/internal/data"

type TemplateData struct {
	Title      string
	HeaderText string
	CSRFToken  string
	Flash      string
	TicketList []*data.Ticket
	FormErrors map[string]string
	FormData   map[string]string
}

//factory function

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
		TicketList: []*data.Ticket{},
	}
}
