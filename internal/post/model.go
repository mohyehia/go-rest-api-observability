package post

type Post struct {
	UserID       int64  `json:"userId"`
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Link         string `json:"link"`
	CommentCount int32  `json:"comment_count"`
}
