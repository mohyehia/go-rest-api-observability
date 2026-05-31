package post

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mohyehia/rest-api-observability/internal/core"
)

type PostsService interface {
	GetPosts() ([]Post, error)

	GetPost(postId int) (*Post, error)
}

type PostsHandler struct {
	PostsService PostsService
	appMetrics   *core.ApplicationMetrics
}

func RegisterHandlers(mux *http.ServeMux, ps PostsService, appMetrics *core.ApplicationMetrics) {
	handler := &PostsHandler{
		PostsService: ps,
		appMetrics:   appMetrics,
	}
	mux.HandleFunc("GET /posts", handler.getPostsHandler)
	mux.HandleFunc("GET /posts/{post_id}", handler.getPostByIDHandler)
}

func (ph *PostsHandler) getPostsHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("GetPostsHandler :: start")

	posts, err := ph.PostsService.GetPosts()
	if err != nil {
		// Increment counter vector with failure label
		ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts", strconv.Itoa(http.StatusInternalServerError)).Inc()
		core.NewErrorResponse(w, err.Error(), http.StatusInternalServerError, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts", strconv.Itoa(http.StatusInternalServerError)).Inc()
		core.NewErrorResponse(w, err.Error(), http.StatusInternalServerError, "")
		return
	}
	ph.appMetrics.GetPostsTotal.Inc()
	ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts", strconv.Itoa(http.StatusOK)).Inc()
	log.Println("GetPostsHandler :: end")
}

func (ph *PostsHandler) getPostByIDHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("GetPostByIDHandler :: start")
	postIdStr := req.PathValue("post_id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts/{post_id}", strconv.Itoa(http.StatusBadRequest)).Inc()
		core.NewErrorResponse(w, err.Error(), http.StatusBadRequest, "PostID should be a number")
		return
	}
	post, err := ph.PostsService.GetPost(postId)
	if err != nil {
		ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts/{post_id}", strconv.Itoa(http.StatusInternalServerError)).Inc()
		core.NewErrorResponse(w, err.Error(), http.StatusInternalServerError, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts/{post_id}", strconv.Itoa(http.StatusInternalServerError)).Inc()
		core.NewErrorResponse(w, err.Error(), http.StatusInternalServerError, "")
		return
	}
	ph.appMetrics.GetPostByIDTotal.Inc()
	ph.appMetrics.HttpRequestsTotal.WithLabelValues("/posts/{post_id}", strconv.Itoa(http.StatusOK)).Inc()
	log.Println("GetPostByIDHandler :: end")
}
