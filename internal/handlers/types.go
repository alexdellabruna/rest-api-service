package handlers

import "database/sql"

type DBHandler struct {
	DB *sql.DB
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
	Error    string `json:"error"`
}

type ResponseTemplateGet struct {
	ResponseTemplatePutDelete
	Content string `json:"content"`
}
