package post

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PostsClient struct {
	client *http.Client
}

func NewPostsClient(client *http.Client) *PostsClient {
	return &PostsClient{
		client: client,
	}
}

func (pc *PostsClient) GetPosts() ([]Post, error) {
	res, err := pc.client.Get("https://json-placeholder.mock.beeceptor.com/posts")
	if err != nil {
		return nil, err
	}
	log.Printf("Response status: %s\n", res.Status)
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	var posts []Post
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (pc *PostsClient) GetPost(id string) (*Post, error) {
	res, err := pc.client.Get(fmt.Sprintf("https://json-placeholder.mock.beeceptor.com/posts/%s", id))
	if err != nil {
		return nil, err
	}
	log.Printf("Response status: %s\n", res.Status)
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}
	var post Post
	err = json.NewDecoder(res.Body).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
