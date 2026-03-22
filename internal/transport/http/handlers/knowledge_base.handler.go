package handlers

import "enterprise-helpdesk/internal/transport/http/routes"

type KnowledgeBaseHandler struct {
	deps routes.Dependencies
}

func NewKnowledgeBaseHandler(deps routes.Dependencies) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{deps: deps}
}
