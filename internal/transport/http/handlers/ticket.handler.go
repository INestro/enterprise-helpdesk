package handlers

import (
	"enterprise-helpdesk/internal/transport/http/dto"
	"enterprise-helpdesk/internal/transport/http/middleware"
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

func (h *TicketHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateTicketRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(c)

	id := generateID()
	ticketNumber := generateTicketNumber()

	_, err := h.deps.DB.DB.Exec(`
		INSERT INTO tickets (
			id, ticket_number, subject, description,
		    status, priority, created_by, created_at
		)
		VALUES ($1,$2,$3,$4,'open',$5,$6,$7)
	`,
		id,
		ticketNumber,
		req.Subject,
		req.Description,
		req.Priority,
		user.UserID,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":            id,
		"ticket_number": ticketNumber,
	})
}

func (h *TicketHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	_, err := h.deps.DB.DB.Exec(`
		UPDATE tickets
		SET status = COALESCE(NULLIF($1, ''), status),
			priority = COALESCE(NULLIF($2, ''), priority)
		WHERE id = $3
	`, req.Status, req.Priority, id)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id": id,
	})
}

func (h *TicketHandler) GetComments(c *fiber.Ctx) error {
	id := c.Params("id")

	rows, err := h.deps.DB.DB.Query(`
		SELECT id, body, author_id, created_at
		FROM ticket_comments
		WHERE ticket_id = $1
		ORDER BY created_at ASC 
	`, id)

	if err != nil {
		return err
	}
	defer rows.Close()

	var items []fiber.Map

	for rows.Next() {
		var (
			commendID string
			body      string
			authorID  string
			createdAt time.Time
		)

		if err := rows.Scan(&commendID, &body, &authorID, &createdAt); err != nil {
			return err
		}

		items = append(items, fiber.Map{
			"id":         commendID,
			"body":       body,
			"author_id":  authorID,
			"created_at": createdAt,
		})
	}

	return c.JSON(fiber.Map{
		"items": items,
	})
}

func (h *TicketHandler) AddComement(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.AddCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	user := middleware.GetUser(c)

	_, err := h.deps.DB.DB.Exec(`
		INSERT INTO ticket_comments (id, ticket_id, body, author_id, created_at)
		VALUES ($1,$2,$3,$4,$5)
	`,
		generateID(),
		id,
		req.Body,
		user.UserID,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"ok": true,
	})
}

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func generateTicketNumber() string {
	return "INC-" + strconv.FormatInt(time.Now().Unix(), 10)
}
