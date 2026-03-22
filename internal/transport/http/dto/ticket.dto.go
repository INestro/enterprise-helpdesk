package dto

type TicketListResponse struct {
	Items []interface{} `json:"items"`
	Total int           `json:"total"`
}

type CreateTicketRequest struct {
	Type         string                 `json:"type"`
	Priority     string                 `json:"priority"`
	Subject      string                 `json:"subject"`
	Description  string                 `json:"description"`
	CategoryId   string                 `json:"category_id"`
	TeamID       string                 `json:"team_id"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type UpdateTicketRequest struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

type AddCommentRequest struct {
	Body string `json:"body"`
}
