package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-red-hat/internal/handlers"
	"task-red-hat/internal/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
)

func TestDeleteObject(t *testing.T) {

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
			name:                  "Delete existing object in existing bucket",
			bucket:                "test-bucket",
			objectID:              "verylongID",
			content:               "",
			needBucketDataSeeding: true,
			needObjectDataSeeding: true,
			expectedResponseCode:  http.StatusOK,
			expectedResponseBody: handlers.ResponseTemplatePutDelete{
				Bucket:   "test-bucket",
				ObjectID: "verylongID",
				Message:  "Object deleted successfully",
			},
		},
		{
			name:                  "Delete non-existing object in existing bucket",
			bucket:                "test-bucket",
			objectID:              "nonexistent-object-id",
			content:               "",
			needBucketDataSeeding: true,
			needObjectDataSeeding: false,
			expectedResponseCode:  http.StatusNotFound,
			expectedResponseBody: handlers.ResponseTemplatePutDelete{
				Bucket:   "test-bucket",
				ObjectID: "nonexistent-object-id",
				Message:  "Object not found",
			},
		},
		{
			name:                  "Delete non-existing object in non-existing bucket",
			bucket:                "non-existing-bucket",
			objectID:              "verylongID",
			content:               "",
			needBucketDataSeeding: false,
			needObjectDataSeeding: false,
			expectedResponseCode:  http.StatusNotFound,
			expectedResponseBody: handlers.ResponseTemplatePutDelete{
				Bucket:   "non-existing-bucket",
				ObjectID: "verylongID",
				Message:  "Object not found",
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
				if tc.needObjectDataSeeding {
					cleanupTestObjectData(t, db, ctx, tc.bucket, tc.objectID)
				}

				if tc.needBucketDataSeeding {
					cleanupTestBucketData(t, db, ctx, tc.bucket)
				}
			})

			// Now test the DELETE request
			ctx.Request = httptest.NewRequest("DELETE", "/objects/"+tc.bucket+"/"+tc.objectID, nil)

			ctx.Params = gin.Params{
				{Key: "bucket", Value: tc.bucket},
				{Key: "objectID", Value: tc.objectID},
			}

			genericHTTPHandler.DeleteObject(ctx)

			if w.Code != tc.expectedResponseCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedResponseCode, w.Code)
			}

			var actualResponse handlers.ResponseTemplatePutDelete
			if err := json.Unmarshal(w.Body.Bytes(), &actualResponse); err != nil {
				t.Fatalf("Failed to parse JSON response: %v", err)
			}

			if diff := cmp.Diff(tc.expectedResponseBody, actualResponse); diff != "" {
				t.Errorf("DeleteObject() mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
