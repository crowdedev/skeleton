package routes

import (
	"net/http"
	"regexp"

	"google.golang.org/grpc"
)

const API_DOC_PATH = "/api/docs/"

type ApiDoc struct {
}

func (a *ApiDoc) Path() string {
	return API_DOC_PATH
}

func (a *ApiDoc) SetClient(client *grpc.ClientConn) {
}

func (a *ApiDoc) Handle(w http.ResponseWriter, r *http.Request) {
	regex := regexp.MustCompile(API_DOC_PATH)
	http.ServeFile(w, r, regex.ReplaceAllString(r.URL.Path, "swaggers/"))
}
