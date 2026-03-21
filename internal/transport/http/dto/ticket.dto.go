package dto

type TicketListResponse struct {
	Items []interface{} `json:"items"`
	Total int           `json:"total"`
}
