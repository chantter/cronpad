package rest

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"log"
	"net/http"
	"net/url"
)

const (
	ErrInternal = 0
)

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string, errCode int) {
	log.Printf("[WARN] %s", errMsg(r, httpStatusCode, err, details, errCode))
	render.Status(r, httpStatusCode)
	render.JSON(w, r, rest.JSON{"error": err.Error(), "details": details, "code": errCode})
}

func errMsg(r *http.Request, httpStatusCode int, err error, details string, errCode int) string {
	userInfo := ""
	if user, e := GetUserInfo(r); e == nil {
		userInfo = user.ID + "/" + user.ID + " - "
	}

	query := r.URL.String()
	if qun, e := url.QueryUnescape(query); e == nil {
		query = qun
	}

	return fmt.Sprintf("(user: %s) %v {%s} - %s - %d %d", userInfo, query, details, err, httpStatusCode, errCode)
}
