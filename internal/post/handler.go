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
	PostsClient PostsClient
	appMetrics  *core.ApplicationMetrics
}

func RegisterHandlers(mux *http.ServeMux, pc *PostsClient, appMetrics *core.ApplicationMetrics) {
	handler := &PostsHandler{
		PostsClient: *pc,
		appMetrics:  appMetrics,
	}
	mux.HandleFunc("GET /posts", handler.getPostsHandler)
	mux.HandleFunc("GET /posts/{post_id}", handler.getPostByIDHandler)
}

func (ph *PostsHandler) getPostsHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("GetPostsHandler :: start")

	posts, err := ph.PostsClient.GetPosts()
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
	postId := req.PathValue("post_id")
	post, err := ph.PostsClient.GetPost(postId)
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
