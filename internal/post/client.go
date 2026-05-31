package post

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PostsClient struct {
	client  *http.Client
	baseURL string
}

func NewPostsClient(client *http.Client, baseURL string) *PostsClient {
	return &PostsClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (pc *PostsClient) GetPosts() ([]Post, error) {
	res, err := pc.client.Get(pc.baseURL + "/posts")
	if err != nil {
		return nil, err
	}
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

func (pc *PostsClient) GetPost(id int) (*Post, error) {
	res, err := pc.client.Get(fmt.Sprintf("%s/posts/%d", pc.baseURL, id))
	if err != nil {
		return nil, err
	}
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
