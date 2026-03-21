package handlers

import (
	"enterprise-helpdesk/internal/transport/http/dto"
	"enterprise-helpdesk/internal/transport/http/routes"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Ticket struct {
	ID          int       `json:"id"`
	TicketNum   string    `json:"ticket_number"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

type TicketHandler struct {
	deps routes.Dependencies
}

func NewTicketHandler(deps routes.Dependencies) *TicketHandler {
	return &TicketHandler{deps: deps}
}

func (h *TicketHandler) GetList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	search := c.Query("search")
	status := c.Query("status")
	priority := c.Query("priority")

	offset := (page - 1) * limit

	// здесь будет repository в следующих частях
	rows, err := h.deps.DB.DB.Query(`
		SELECT id, ticket_number, subject, status, priority, created_at
		FROM tickets
		WHERE ($1 = '' OR subject ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR status = $2)
		  AND ($3 = '' OR priority = $3)
		ORDER BY created_at DESC 
		LIMIT $4 OFFSET $5
	`, search, status, priority, limit, offset)

	if err != nil {
		return err
	}
	defer rows.Close()

	var items []fiber.Map

	for rows.Next() {
		var (
			id           string
			ticketNumber string
			subject      string
			statusVal    string
			priorityVal  string
			createdAt    time.Time
		)

		if err := rows.Scan(&id, &ticketNumber, &subject, &statusVal, &priorityVal, &createdAt); err != nil {
			return err
		}

		items = append(items, fiber.Map{
			"id":            id,
			"ticket_number": ticketNumber,
			"subject":       subject,
			"status":        statusVal,
			"priority":      priorityVal,
			"created_at":    createdAt,
		})
	}
	interfaceItems := make([]interface{}, len(items))
	for i, v := range items {
		interfaceItems[i] = v
	}

	var total int
	_ = h.deps.DB.DB.QueryRow(`SELECT COUNT(*) FROM tickets`).Scan(&total)

	return c.JSON(dto.TicketListResponse{
		Items: interfaceItems,
		Total: total,
	})
}

func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	//var ticket fiber.Map
	var t Ticket

	err := h.deps.DB.DB.QueryRow(`
		SELECT id, ticket_number, subject, description, status, priority, created_at
		FROM tickets
		WHERE id = $1
	`, id).Scan(
		&t.ID,
		&t.TicketNum,
		&t.Subject,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedAt,
	)

	if err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(t)
}
