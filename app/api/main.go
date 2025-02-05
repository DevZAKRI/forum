package api

import (
	"database/sql"
	"encoding/json"
	"forum/app/api/auth"
	"forum/app/api/comments"
	"forum/app/api/posts"
	"forum/app/api/reactions"
	"forum/app/modules"
	"net/http"
	"strings"
)

func Router(resp http.ResponseWriter, req *http.Request, forumDB *sql.DB) {
	conn := &modules.Connection{
		Resp: resp,
		Req:  req,
	}

	path := strings.Split(req.URL.Path[5:], "/")
	switch path[0] {
	case "auth":
		auth.Entry(conn, forumDB)
	case "posts":
		if req.Method == http.MethodGet {
			data, err := posts.GetPosts(conn, forumDB)
			if err != nil {
				http.Error(resp, "500 - internal server error", 500)
				return
			}
			resp.Header().Set("Content-Type", "application/json")
			resp.Write(data)
		} else if req.Method == http.MethodPost {
			posts.AddPost(conn, forumDB)
		}
	case "coments":
		if req.Method == http.MethodPost {
			err := comments.AddComment(conn, forumDB)
			if err != nil {
				http.Error(resp, err.Error()+"500 - internal server error", 500)
				return
			}
		} else if req.Method == http.MethodGet {
			coments, err := comments.GetComents(req.URL, forumDB)
			if err != nil {
				http.Error(resp, err.Error()+"500 - internal server error", 500)
				return
			}
			err = json.NewEncoder(resp).Encode(coments)
			if err != nil {
				http.Error(resp, err.Error()+"500 - internal server error", 500)
				return
			}
		} else {
			http.Error(conn.Resp, "405 - method not allowed ", 405)
			return
		}
	case "reactions":
		reactions.HandleReactions(conn, path, conn.Req.Method, forumDB)
	default:

		return
	}
}
