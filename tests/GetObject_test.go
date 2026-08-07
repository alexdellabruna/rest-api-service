package tests

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http/httptest"
	"strings"
	"task-red-hat/handlers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
)

func TestGetObjectSuccess(t *testing.T) {
	db, err := sql.Open("sqlite3", "../db_data/local.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbHandler := &handlers.DBHandler{DB: db}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	bucket := "test-bucket"
	objectID := "verylongID"

	ctx.Request = httptest.NewRequest("PUT", "/objects/"+bucket+"/"+objectID, strings.NewReader("{\"content\": \"This is a test object content.\"}"))

	ctx.Params = gin.Params{
		{Key: "bucket", Value: bucket},
		{Key: "objectID", Value: objectID},
	}

	dbHandler.PutObject(ctx)

	if w.Code != 201 {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	// Now test the GET request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest("GET", "/objects/"+bucket+"/"+objectID, nil)

	ctx.Params = gin.Params{
		{Key: "bucket", Value: bucket},
		{Key: "objectID", Value: objectID},
	}

	dbHandler.GetObject(ctx)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	expectedResponse := handlers.ResponseTemplateGet{
		// first two fields are maintened for reference, the content field is the actual object content
		ResponseTemplatePutDelete: handlers.ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Object retrieved successfully",
			Error:    "",
		},
		Content: "This is a test object content.",
	}

	var actualResponse handlers.ResponseTemplateGet
	if err := json.Unmarshal(w.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if diff := cmp.Diff(expectedResponse, actualResponse); diff != "" {
		t.Errorf("GetObject() mismatch (-expected +actual):\n%s", diff)
	}
}
