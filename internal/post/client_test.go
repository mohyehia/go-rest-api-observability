package post

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostsClient_GetPosts_Success(t *testing.T) {
	// Arrange
	expectedResponse := []Post{
		{
			UserID:       1,
			Id:           1,
			Title:        "Test Post",
			Body:         "Test post body",
			Link:         "dummy link",
			CommentCount: 10,
		},
	}

	// Start a local, in-memory mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/posts" {
			t.Errorf("Expected path /posts, got %s", req.URL.Path)
		}
		// Return a 200 OK and serialize the response payload back
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResponse)
	}))
	// Always shut down the local server when the test finishes
	defer server.Close()

	pc := NewPostsClient(http.DefaultClient, server.URL)

	posts, err := pc.GetPosts()

	if err != nil {
		t.Fatalf("Expected no errors from GetPosts, got: %v", err)
	}
	if len(posts) != len(expectedResponse) {
		t.Fatalf("Expected %d posts, got %d", len(expectedResponse), len(posts))
	}

	got := posts[0]
	want := expectedResponse[0]

	assertRetrievedPost(t, got, want)
}

func TestPostsClient_GetPosts_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/posts" {
			t.Errorf("Expected path /posts, got %s", req.URL.Path)
		}
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer server.Close()

	pc := NewPostsClient(http.DefaultClient, server.URL)

	_, err := pc.GetPosts()

	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	got := err.Error()
	want := "status code 400"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestPostsClient_GetPost_Success(t *testing.T) {
	expectedResponse := &Post{
		UserID:       1,
		Id:           1,
		Title:        "Test Post",
		Body:         "Test post body",
		Link:         "dummy link",
		CommentCount: 10,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/posts/1" {
			t.Errorf("Expected path /posts/1, got %s", req.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResponse)
	}))

	defer server.Close()

	pc := NewPostsClient(http.DefaultClient, server.URL)

	post, err := pc.GetPost(1)

	if err != nil {
		t.Fatalf("Expected no errors from GetPost, got: %v", err)
	}

	assertRetrievedPost(t, *post, *expectedResponse)
}

func TestPostsClient_GetPost_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/posts/1" {
			t.Errorf("Expected path /posts/1, got %s", req.URL.Path)
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	pc := NewPostsClient(http.DefaultClient, server.URL)

	_, err := pc.GetPost(1)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	got := err.Error()
	want := "status code 400"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertRetrievedPost(t *testing.T, got Post, want Post) {
	if got.Id != want.Id {
		t.Errorf("Expected id %d, got %d", want.Id, got.Id)
	}
	if got.UserID != want.UserID {
		t.Errorf("Expected userID %d, got %d", want.UserID, got.UserID)
	}
	if got.Title != want.Title {
		t.Errorf("Expected title %s, got %s", want.Title, got.Title)
	}
	if got.Body != want.Body {
		t.Errorf("Expected body %s, got %s", want.Body, got.Body)
	}
	if got.Link != want.Link {
		t.Errorf("Expected link %s, got %s", want.Link, got.Link)
	}
	if got.CommentCount != want.CommentCount {
		t.Errorf("Expected comment count %d, got %d", want.CommentCount, got.CommentCount)
	}
}
