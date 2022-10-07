package service

import (
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"strings"
)

type Post interface {
	//	CRUD
	CreatePost(post models.Post, user models.User) error
	UpdatePost(postId int, post models.Post, user models.User) error
	DeletePost(postId int, user models.User) error
	//	GET
	GetAllPosts(filter string) ([]models.Post, error)
	GetPostById(id int) (models.Post, error)
	GetPostsByUsername(username string) ([]models.Post, error)
	UpdateCountsReactionsPost(likes int, dislikes int, postId int) error
}

type PostService struct {
	storage storage.Post
}

func newPostService(storage storage.Post) *PostService {
	return &PostService{
		storage: storage,
	}
}

func (p *PostService) CreatePost(post models.Post, user models.User) error {
	post.Category = strings.Fields(post.Category[0])

	if err := p.storage.CreatePost(user.Username, post); err != nil {
		return err
	}
	return nil
}

func (p *PostService) GetPostById(id int) (models.Post, error) {

	post, err := p.storage.GetPostById(id)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil

}
func (p *PostService) DeletePost(postId int, user models.User) error {
	post, err := p.storage.GetPostById(postId)
	if err != nil {
		return err
	}
	if user.Username != post.Creator {
		return fmt.Errorf("service: u cant delete dat post")
	}
	if err := p.storage.DeletePost(post); err != nil {
		return err
	}
	return nil
}
func (p *PostService) UpdatePost(postId int, post models.Post, user models.User) error {

	if user.Username != post.Creator {
		return fmt.Errorf("service: u cant change dat post")
	}
	if err := p.storage.UpdatePost(postId, post); err != nil {
		return err
	}
	return nil
}
func (p *PostService) GetAllPosts(filter string) ([]models.Post, error) {
	return p.storage.GetAllPosts(filter)
}

func (p *PostService) GetPostsByUsername(username string) ([]models.Post, error) {
	posts, err := p.storage.GetPostsByUsername(username)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
func (p *PostService) UpdateCountsReactionsPost(likes int, dislikes int, postId int) error {
	if err := p.storage.UpdateCountsReactionsPost(likes, dislikes, postId); err != nil {
		return err
	}
	return nil
}
