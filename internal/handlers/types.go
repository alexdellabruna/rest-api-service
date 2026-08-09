package handlers

import "task-red-hat/internal/storage"

type GenericHTTPHandler struct {
	DBHandler *storage.DBHandler
}

type ApiRequestGetPutDelete struct {
	Bucket   string `uri:"bucket" binding:"required"`
	ObjectID string `uri:"objectID" binding:"required"`
}

type ApiRequestPutBody struct {
	Content string `json:"content" binding:"required"`
}

type ResponseTemplatePutDelete struct {
	Bucket   string `json:"bucket"`
	ObjectID string `json:"objectID"`
	Message  string `json:"message"`
}

type ResponseTemplateGet struct {
	ResponseTemplatePutDelete
	Content string `json:"content"`
}
