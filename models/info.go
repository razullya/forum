package models

type Info struct {
	User       User
	ThatUser   User
	Posts      []Post
	Post       Post
	Comments   []Comment
	Categories []string
}
