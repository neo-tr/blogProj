package models

import (
	"context"
	"log"
	"time"

	"blogProj/db"
)

// Comment представляет комментарий к посту.
type Comment struct {
	ID             int
	PostID         int       // ID поста, к которому относится комментарий
	UserID         int       // ID пользователя, который оставил комментарий
	Content        string    // Текст комментария
	CreatedAt      time.Time // Дата и время создания комментария
	AuthorNickname string    // Ник автора комментария
	Editable       bool      // Поле для шаблона: может ли текущий пользователь редактировать этот комментарий
}

// CreateComment добавляет новый комментарий к посту.
func CreateComment(postID, userID int, content string) error {
	sql := INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)
	_, err := db.Pool.Exec(context.Background(), sql, postID, userID, content)
	log.Printf("Сохраняем комментарий: postID=%d, userID=%d, content=%s", postID, userID, content)
	return err
}

// GetCommentsByPostID возвращает все комментарии для заданного поста с информацией об авторе.
func GetCommentsByPostID(postID int) ([]Comment, error) {
	sql := `
    SELECT 
      c.id, 
      c.post_id, 
      c.user_id, 
      c.content, 
      c.created_at, 
      u.nickname 
    FROM comments c 
    LEFT JOIN users u ON c.user_id = u.id 
    WHERE c.post_id = $1 
    ORDER BY c.created_at ASC`
	rows, err := db.Pool.Query(context.Background(), sql, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var cmt Comment
		if err := rows.Scan(&cmt.ID, &cmt.PostID, &cmt.UserID, &cmt.Content, &cmt.CreatedAt, &cmt.AuthorNickname); err != nil {
			return nil, err
		}
		comments = append(comments, cmt)
	}
	return comments, nil
}

// GetCommentByID возвращает комментарий по ID.
func GetCommentByID(commentID int) (Comment, error) {
	sql := `
    SELECT 
      c.id, 
      c.post_id, 
      c.user_id, 
      c.content, 
      c.created_at, 
      u.nickname 
    FROM comments c 
    LEFT JOIN users u ON c.user_id = u.id 
    WHERE c.id = $1`
	var cmt Comment
	err := db.Pool.QueryRow(context.Background(), sql, commentID).Scan(
		&cmt.ID, &cmt.PostID, &cmt.UserID, &cmt.Content, &cmt.CreatedAt, &cmt.AuthorNickname,
	)
	return cmt, err
}

// UpdateComment обновляет содержание комментария.
func UpdateComment(commentID int, content string) error {
	sql := UPDATE comments SET content=$1 WHERE id=$2
	_, err := db.Pool.Exec(context.Background(), sql, content, commentID)
	return err
}

// DeleteComment удаляет комментарий по его ID.
func DeleteComment(commentID int) error {
	sql := DELETE FROM comments WHERE id = $1
	_, err := db.Pool.Exec(context.Background(), sql, commentID)
	return err
}