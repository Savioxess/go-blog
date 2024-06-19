package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Savioxess/blog/internals/database"
	"github.com/Savioxess/blog/internals/models"
	"github.com/Savioxess/blog/internals/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct{}

func (handler *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	loginUserDetails := &models.User{}
	err := utils.GetRequestBodyJSON(r.Body, loginUserDetails)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	userFromDatabase, err := getUserAndPasswordFromDatabase(loginUserDetails.Email)

	if err != nil {
		response := &utils.Error{
			Error: "User With Email Does Not Exist",
		}

		responseJSON, err := utils.EncodeJSONResponse(response)

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		utils.ClientErrorResponse(409, responseJSON, r.Method, r.URL.Path, w)
		return
	}

	isPasswordMatching := comparePassword(userFromDatabase.Password, loginUserDetails.Password)

	if !isPasswordMatching {
		response := &utils.Error{
			Error: "Invalid Credentials",
		}

		responseJSON, err := utils.EncodeJSONResponse(response)

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		utils.ClientErrorResponse(401, responseJSON, r.Method, r.URL.Path, w)
		return
	}

	token, err := generateJWTToken(userFromDatabase.ID)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	response := &utils.Success{
		Success: true,
		Message: map[string]string{
			"token": token},
	}

	responseJSON, err := utils.EncodeJSONResponse(response)

	if err != nil {
		utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
		return
	}

	utils.SuccessResponse(201, responseJSON, r.Method, r.URL.Path, w)
}

func getUserAndPasswordFromDatabase(email string) (models.User, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	row := database.DB.QueryRowContext(ctx, "SELECT id, email, password FROM user WHERE email=?", email)

	userResultFromDatabase := models.User{}

	if err := row.Scan(&userResultFromDatabase.ID, &userResultFromDatabase.Email, &userResultFromDatabase.Password); err != nil {
		return models.User{}, err
	}

	return userResultFromDatabase, nil
}

func comparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func generateJWTToken(payload []byte) (string, error) {
	godotenv.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Blog",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		Subject:   string(payload),
	})

	tokenAsString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenAsString, nil
}
