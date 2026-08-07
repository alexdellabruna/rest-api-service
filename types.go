package main

import "database/sql"

type DBHandler struct {
	DB *sql.DB
}

type ApiRequestGetDelete struct {
	BucketID string `uri:"bucketID" binding:"required"`
	ObjectID string `uri:"objectID" binding:"required"`
}

type ApiRequestPut struct {
	ApiRequestGetDelete
	Content string `form:"content" binding:"required"`
}

type responseTemplatePutDelete struct {
	BucketID string `json:"bucketID"`
	ObjectID string `json:"objectID"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

type responseTemplateGet struct {
	responseTemplatePutDelete
	Content string `json:"content"`
}
