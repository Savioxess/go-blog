package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Savioxess/blog/internals/database"
	"github.com/Savioxess/blog/internals/models"
	"github.com/Savioxess/blog/internals/utils"
	"github.com/google/uuid"
)

type CreatBlogHandler struct{}

func (handler *CreatBlogHandler) Handle(w http.ResponseWriter, r *http.Request) {
	requestBlogBody := &models.Blog{}
	err := utils.GetRequestBodyJSON(r.Body, requestBlogBody)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	user, err := getUserByID([]byte(r.Header.Get("UserID")))

	if err != nil {
		response := &utils.Error{
			Error: "User Does Not Exist",
		}

		responseJSON, err := utils.EncodeJSONResponse(response)

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		utils.ClientErrorResponse(409, responseJSON, r.Method, r.URL.Path, w)
		return
	}

	err = createBlog(user.ID, requestBlogBody.Title, requestBlogBody.Content)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	response := &utils.Success{
		Success: true,
		Message: "Blog Has Been Created",
	}

	responseJSON, err := utils.EncodeJSONResponse(response)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	utils.SuccessResponse(200, responseJSON, r.Method, r.URL.Path, w)
}

func createBlog(authorId []byte, title, content string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	blogId, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	postedOnDate := time.Now().UTC()
	postedOnDateString := postedOnDate.Format("2006-01-02")

	_, err = database.DB.ExecContext(ctx, "INSERT INTO blog VALUES(?, ?, ?, ?, ?)", blogId, authorId, title, content, postedOnDateString)

	if err != nil {
		return err
	}

	return nil
}

func getUserByID(userID []byte) (*models.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	row := database.DB.QueryRowContext(ctx, "SELECT id, email, username, joined_on password FROM user WHERE id=?", userID)

	userResultFromDatabase := &models.User{}

	if err := row.Scan(&userResultFromDatabase.ID, &userResultFromDatabase.Email, &userResultFromDatabase.Username, &userResultFromDatabase.JoinedOn); err != nil {
		return &models.User{}, err
	}

	return userResultFromDatabase, nil
}
