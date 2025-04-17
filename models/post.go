package models

import (
	"blogProj/db"
	"context"
	"time"
)

type Post struct {
	ID             int
	UserID         *int // Автор (ссылка на пользователя)
	Title          string
	Content        string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	Edited         bool
	AuthorNickname string // Новый параметр для хранения ника автора
}

// CreatePost с сохранением автора
func CreatePost(userID int, title, content string) error {
	sql := INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3)
	_, err := db.Pool.Exec(context.Background(), sql, userID, title, content)
	return err
}

// GetAllPosts возвращает все посты, отсортированные по дате создания.
func GetAllPosts() ([]Post, error) {
	sql := SELECT id, user_id, title, content, created_at FROM posts ORDER BY created_at DESC
	rows, err := db.Pool.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// GetPostByID возвращает пост по id (с информацией об авторе может быть объединен через JOIN)
func GetPostByID(id int) (Post, error) {
	// Обратите внимание: здесь выбираем дополнительные данные из таблицы users (u.nickname)
	sql := `SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at, p.edited, u.nickname 
          FROM posts p 
          LEFT JOIN users u ON p.user_id = u.id 
          WHERE p.id=$1`
	var p Post
	err := db.Pool.QueryRow(context.Background(), sql, id).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.Edited, &p.AuthorNickname,
	)
	return p, err
}

// UpdatePost позволяет редактировать пост
func UpdatePost(postID int, title, content string) error {
	sql := UPDATE posts SET title=$1, content=$2, updated_at=NOW(), edited=true WHERE id=$3
	_, err := db.Pool.Exec(context.Background(), sql, title, content, postID)
	return err
}

func SearchPosts(query string) ([]Post, error) {
	sql := `SELECT id, user_id, title, content, created_at, updated_at, edited 
          FROM posts WHERE title ILIKE '%'  $1  '%' OR content ILIKE '%'  $1  '%'
          ORDER BY created_at DESC`
	rows, err := db.Pool.Query(context.Background(), sql, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.Edited)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}