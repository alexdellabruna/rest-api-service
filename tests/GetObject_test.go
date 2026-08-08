package tests

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"task-red-hat/handlers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
)

func seedTestData(t *testing.T, db *sql.DB, ctx *gin.Context, bucket string, objectID string) {
	_, err := db.ExecContext(ctx, "INSERT INTO buckets (id) VALUES (?)", bucket)
	if err != nil {
		t.Errorf("Failed to seed test data: %v", err)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO objects (id, bucket_id, content, sha256_hash) VALUES (?, ?, ?, ?)", objectID, bucket, "This is a test object content.", "abcd1234efgh5678ijkl9012mnop3456qrst7890uvwx1234yzab5678cdef9012")

	if err != nil {
		t.Errorf("Failed to seed test data: %v", err)
	}
}

func cleanupTestData(t *testing.T, db *sql.DB, ctx *gin.Context, bucket string, objectID string) {
	_, err := db.ExecContext(ctx, "DELETE FROM objects WHERE id = ? AND bucket_id = ?", objectID, bucket)
	if err != nil {
		t.Errorf("Failed to clean up test data: %v", err)
	}

	_, err = db.ExecContext(ctx, "DELETE FROM buckets WHERE id = ?", bucket)
	if err != nil {
		t.Errorf("Failed to clean up test data: %v", err)
	}
}

func TestGetObject(t *testing.T) {

	testCases := []struct {
		name                 string
		bucket               string
		objectID             string
		content              string
		needDataSeeding      bool
		expectedResponseCode int
		expectedResponseBody handlers.ResponseTemplateGet
	}{
		{
			name:                 "Get existing object",
			bucket:               "test-bucket",
			objectID:             "verylongID",
			content:              "This is a test object content.",
			needDataSeeding:      true,
			expectedResponseCode: 200,
			expectedResponseBody: handlers.ResponseTemplateGet{
				ResponseTemplatePutDelete: handlers.ResponseTemplatePutDelete{
					Bucket:   "test-bucket",
					ObjectID: "verylongID",
					Message:  "Object retrieved successfully",
					Error:    "",
				},
				Content: "This is a test object content.",
			},
		},
		{
			name:                 "Get non-existing object",
			bucket:               "nonexistent-bucket",
			objectID:             "nonexistent-object",
			content:              "",
			needDataSeeding:      false,
			expectedResponseCode: 404,
			expectedResponseBody: handlers.ResponseTemplateGet{
				ResponseTemplatePutDelete: handlers.ResponseTemplatePutDelete{
					Bucket:   "nonexistent-bucket",
					ObjectID: "nonexistent-object",
					Message:  "Object not found",
					Error:    "Not found",
				},
				Content: "",
			},
		},
	}

	db := dbConnectAndInit(t)
	defer db.Close()

	dbHandler := &handlers.DBHandler{DB: db}

	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tc.needDataSeeding {
				seedTestData(t, db, ctx, tc.bucket, tc.objectID)
				t.Cleanup(func() {
					cleanupTestData(t, db, ctx, tc.bucket, tc.objectID)
				})
			}

			// Now test the GET request success case
			ctx.Request = httptest.NewRequest("GET", "/objects/"+tc.bucket+"/"+tc.objectID, nil)

			ctx.Params = gin.Params{
				{Key: "bucket", Value: tc.bucket},
				{Key: "objectID", Value: tc.objectID},
			}

			dbHandler.GetObject(ctx)

			if w.Code != tc.expectedResponseCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedResponseCode, w.Code)
			}

			var actualResponse handlers.ResponseTemplateGet
			if err := json.Unmarshal(w.Body.Bytes(), &actualResponse); err != nil {
				t.Fatalf("Failed to parse JSON response: %v", err)
			}

			if diff := cmp.Diff(tc.expectedResponseBody, actualResponse); diff != "" {
				t.Errorf("GetObject() mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
