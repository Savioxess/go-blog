package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Savioxess/blog/internals/database"
	"github.com/Savioxess/blog/internals/models"
	"github.com/Savioxess/blog/internals/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct{}

func (handler *SignupHandler) Handle(w http.ResponseWriter, r *http.Request) {
	signupUserDetails := &models.User{}
	err := utils.GetRequestBodyJSON(r.Body, signupUserDetails)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	err = validateUser(signupUserDetails.Email, signupUserDetails.Username, signupUserDetails.Password)

	if err != nil {
		response := &utils.Error{
			Error: err.Error(),
		}

		responseJSON, err := utils.EncodeJSONResponse(response)

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		utils.ClientErrorResponse(400, responseJSON, r.Method, r.URL.Path, w)
		return
	}

	doesUserExist := checkIfUserWithEmailExists(signupUserDetails.Email)

	if doesUserExist {
		response := &utils.Error{
			Error: "User With Email Already Exists",
		}

		responseJSON, err := utils.EncodeJSONResponse(response)

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		utils.ClientErrorResponse(409, responseJSON, r.Method, r.URL.Path, w)
		return
	}

	err = createUser(signupUserDetails.Email, signupUserDetails.Username, signupUserDetails.Password)

	if err != nil {
		log.Println(err)
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	response := &utils.Success{
		Success: true,
		Message: "User Has Been Created",
	}

	responseJSON, err := utils.EncodeJSONResponse(response)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	utils.SuccessResponse(201, responseJSON, r.Method, r.URL.Path, w)
}

func validateUser(email, username, password string) error {
	matched, err := regexp.MatchString("@", email)

	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid Email")
	}

	if len(username) < 5 {
		return errors.New("username should be atleast 5 characters")
	}

	if len(password) < 4 {
		return errors.New("password should be more than 3 characters")
	}

	return nil
}

func createUser(email, username, password string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	userId, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	hashedUserPassword, err := hashUserPassword(password)

	if err != nil {
		return err
	}

	joinedOnDate := time.Now().UTC()
	joinedOnDateString := joinedOnDate.Format("2006-01-02")

	_, err = database.DB.ExecContext(ctx, "INSERT INTO user VALUES(?, ?, ?, ?, ?)", userId, email, username, hashedUserPassword, joinedOnDateString)

	if err != nil {
		return err
	}

	return nil
}

func hashUserPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func checkIfUserWithEmailExists(email string) bool {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	row := database.DB.QueryRowContext(ctx, "SELECT email FROM user WHERE email=?", email)

	if err := row.Scan(&email); err != nil {
		return false
	}

	return true
}
