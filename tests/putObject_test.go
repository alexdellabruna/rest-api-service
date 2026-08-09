package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"task-red-hat/internal/handlers"
	"task-red-hat/internal/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
)

func TestPutObject(t *testing.T) {

	testCases := []struct {
		name                  string
		bucket                string
		objectID              string
		content               string
		needBucketDataSeeding bool
		needObjectDataSeeding bool
		expectedResponseCode  int
		expectedResponseBody  handlers.ResponseTemplatePutDelete
	}{
		{
			name:                  "Put object in non existing bucket",
			bucket:                "test-bucket",
			objectID:              "verylongID",
			content:               "{\"content\":\"This is a test object content.\"}",
			needBucketDataSeeding: false,
			needObjectDataSeeding: false,
			expectedResponseCode:  http.StatusCreated,
			expectedResponseBody: handlers.ResponseTemplatePutDelete{
				Bucket:   "test-bucket",
				ObjectID: "verylongID",
				Message:  "Object stored successfully",
			},
		},
		{
			name:                  "Put object in existent bucket",
			bucket:                "test-bucket",
			objectID:              "verylongID",
			content:               "{\"content\":\"This is a test object content.\"}",
			needBucketDataSeeding: true,
			needObjectDataSeeding: false,
			expectedResponseCode:  http.StatusCreated,
			expectedResponseBody: handlers.ResponseTemplatePutDelete{
				Bucket:   "test-bucket",
				ObjectID: "verylongID",
				Message:  "Object stored successfully",
			},
		},
		{
			name:                  "Put existing object in existent bucket",
			bucket:                "test-bucket",
			objectID:              "verylongID",
			content:               "{\"content\":\"This is a test object content.\"}",
			needBucketDataSeeding: true,
			needObjectDataSeeding: true,
			expectedResponseCode:  http.StatusConflict,
			expectedResponseBody: handlers.ResponseTemplatePutDelete{
				Bucket:   "test-bucket",
				ObjectID: "verylongID",
				Message:  "An object with the same content already exists in the bucket",
			},
		},
	}

	db := dbConnectAndInit(t)
	defer db.Close()

	genericHTTPHandler := &handlers.GenericHTTPHandler{DBHandler: &storage.DBHandler{DB: db}}

	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tc.needBucketDataSeeding {
				seedTestBucketData(t, db, ctx, tc.bucket)
			}

			if tc.needObjectDataSeeding {
				seedTestObjectData(t, db, ctx, tc.bucket, tc.objectID)
			}

			t.Cleanup(func() {
				// this must be called every time to clean up the test data
				cleanupTestObjectData(t, db, ctx, tc.bucket, tc.objectID)
				cleanupTestBucketData(t, db, ctx, tc.bucket)
			})

			// Now test the PUT request
			ctx.Request = httptest.NewRequest("PUT", "/objects/"+tc.bucket+"/"+tc.objectID, strings.NewReader(tc.content))

			ctx.Params = gin.Params{
				{Key: "bucket", Value: tc.bucket},
				{Key: "objectID", Value: tc.objectID},
			}

			genericHTTPHandler.PutObject(ctx)

			if w.Code != tc.expectedResponseCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedResponseCode, w.Code)
			}

			var actualResponse handlers.ResponseTemplatePutDelete
			if err := json.Unmarshal(w.Body.Bytes(), &actualResponse); err != nil {
				t.Fatalf("Failed to parse JSON response: %v", err)
			}

			if diff := cmp.Diff(tc.expectedResponseBody, actualResponse); diff != "" {
				t.Errorf("PutObject() mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
