package comments

import (
	"database/sql"
	"errors"
	"net/url"
	"strconv"
)

func GetComents(URL *url.URL, forumDB *sql.DB) ([]Comment, error) {
	post_id := URL.Query().Get("p_id")
	Offset := URL.Query().Get("offset")
	if Offset == "" {
		return nil, errors.New("invalid url")
	}
	offset, err := strconv.Atoi(Offset)
	if err != nil {
		return nil, errors.New("invalid post id")
	}
	if post_id == "" {
		return nil, errors.New("invalid url")
	}
	p_id, err := strconv.Atoi(post_id)
	if err != nil {
		return nil, errors.New("invalid post id")
	}
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ? ORDER BY updated_at DESC LIMIT 10 OFFSET ?`
	rows, err := forumDB.Query(query, p_id, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.Item_id, &comment.Post_id, &comment.User_id, &comment.Content, &comment.Created_at); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
