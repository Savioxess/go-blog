package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Savioxess/blog/internals/database"
	"github.com/Savioxess/blog/internals/models"
	"github.com/Savioxess/blog/internals/utils"
)

type GetUserBlogsHanlder struct{}

func (handler *GetUserBlogsHanlder) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userid")

	blogs, err := getUsersBlogs([]byte(userID))

	if err != nil {
		response := &utils.Error{
			Error: err.Error(),
		}

		responseJSON, err := utils.EncodeJSONResponse(response)

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		utils.ClientErrorResponse(409, responseJSON, r.Method, r.URL.Path, w)
		return
	}

	response := utils.Success{
		Success: true,
		Message: map[string]interface{}{
			"blogs": blogs,
		},
	}

	responseJSON, err := utils.EncodeJSONResponse(response)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	utils.SuccessResponse(200, responseJSON, r.Method, r.URL.Path, w)
}

func getUsersBlogs(authorID []byte) ([]models.Blog, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	rows, err := database.DB.QueryContext(ctx, "SELECT id, author_id, title, content, posted_on password FROM blog WHERE author_id=?", authorID)

	if err != nil {
		return []models.Blog{}, err
	}

	var blogs []models.Blog = []models.Blog{}

	for rows.Next() {
		var blog models.Blog

		if err := rows.Scan(&blog.ID, &blog.AuthorID, &blog.Title, &blog.Content, &blog.PostedOn); err != nil {
			return []models.Blog{}, nil
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
