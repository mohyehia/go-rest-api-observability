package post

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohyehia/rest-api-observability/internal/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type mockPostClient struct {
	client  *http.Client
	baseURL string
}

func (mock *mockPostClient) GetPosts() ([]Post, error) {
	return []Post{
		{
			UserID:       1,
			Id:           1,
			Title:        "Test Post",
			Body:         "Test post body",
			Link:         "dummy link",
			CommentCount: 10,
		},
	}, nil
}

func (mock *mockPostClient) GetPost(id int) (*Post, error) {
	return &Post{
		UserID:       1,
		Id:           1,
		Title:        "test post",
		Body:         "test post body",
		Link:         "dummy link",
		CommentCount: 0,
	}, nil
}

func TestGetPostsHandler_IncrementMetrics(t *testing.T) {
	// Arrange
	registry := prometheus.NewRegistry()
	appMetrics := core.NewApplicationMetrics(registry)

	mockClient := &mockPostClient{}

	mux := http.NewServeMux()

	// Track baseline value of the counter before hitting the endpoint
	baseGetPostsTotal := testutil.ToFloat64(appMetrics.GetPostsTotal)

	// Act
	RegisterHandlers(mux, mockClient, appMetrics)
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("incorrect status code: got %v want %v", rec.Code, http.StatusOK)
	}

	var actualPost []Post
	if err := json.NewDecoder(rec.Body).Decode(&actualPost); err != nil {
		t.Fatalf("failed to parse handler json response: %v", err)
	}

	// Assert that Prometheus Registry isolation increments correctly
	finalGetPostsTotal := testutil.ToFloat64(appMetrics.GetPostsTotal)
	if finalGetPostsTotal-baseGetPostsTotal != 1 {
		t.Errorf("expected metric counter to increase by 1, got baseline: %f, final: %f", baseGetPostsTotal, finalGetPostsTotal)
	}
}

func TestGetPostByIDHandler_IncrementMetrics(t *testing.T) {
	registry := prometheus.NewRegistry()
	appMetrics := core.NewApplicationMetrics(registry)
	mockClient := &mockPostClient{}
	mux := http.NewServeMux()

	baseGetPostByIDTotal := testutil.ToFloat64(appMetrics.GetPostByIDTotal)

	RegisterHandlers(mux, mockClient, appMetrics)
	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("incorrect status code: got %v want %v", rec.Code, http.StatusOK)
	}

	var actualPost Post
	if err := json.NewDecoder(rec.Body).Decode(&actualPost); err != nil {
		t.Fatalf("failed to parse handler json response: %v", err)
	}

	finalGetPostByIDTotal := testutil.ToFloat64(appMetrics.GetPostByIDTotal)
	if finalGetPostByIDTotal-baseGetPostByIDTotal != 1 {
		t.Errorf("expected metric counter to increase by 1, got baseline: %f, final: %f", baseGetPostByIDTotal, finalGetPostByIDTotal)
	}
}
