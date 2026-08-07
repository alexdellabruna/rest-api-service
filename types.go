package main

import "database/sql"

type DBHandler struct {
	DB *sql.DB
}

type ApiRequestGetDelete struct {
	BucketID int64 `form:"bucketID" binding:"required"`
	ObjectID int64 `form:"objectID" binding:"required"`
}

type ApiRequestPut struct {
	ApiRequestGetDelete
	Content string `form:"content" binding:"required"`
}

type responseTemplatePutDelete struct {
	BucketID int64  `json:"bucketID"`
	ObjectID int64  `json:"objectID"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

type responseTemplateGet struct {
	responseTemplatePutDelete
	Content string `json:"content"`
}
