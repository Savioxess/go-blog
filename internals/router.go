package internals

import (
	"net/http"
	"os"

	"github.com/Savioxess/blog/internals/api"
	"github.com/joho/godotenv"
)

var Router *http.ServeMux

type Handler interface {
	Handle(http.ResponseWriter, *http.Request)
}

func init() {
	godotenv.Load()

	Router = http.NewServeMux()
	cfg := APIConfig{JWT_SECRET: os.Getenv("JWT_SECRET")}
	SignupHandler := api.SignupHandler{}
	LoginHandler := api.LoginHandler{}
	CreateBlogHandler := api.CreatBlogHandler{}
	GetUserBlogsHandler := api.GetUserBlogsHanlder{}
	GetAllBlogsHandler := api.GetAllBlogsHandler{}

	Router.Handle("POST /api/signup", cfg.Handler(&SignupHandler))
	Router.Handle("POST /api/login", cfg.Handler(&LoginHandler))
	Router.Handle("POST /api/blog", cfg.GetUserIDFromToken(&CreateBlogHandler))
	Router.Handle("GET /api/blog/author/{userid}", cfg.Handler(&GetUserBlogsHandler))
	Router.Handle("GET /api/blog", cfg.Handler(&GetAllBlogsHandler))
}
