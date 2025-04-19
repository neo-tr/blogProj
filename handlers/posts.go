package handlers

import (
	"bytes"
	"github.com/yuin/goldmark"
	"log"
	"net/http"
	"strconv"

	"blogProj/models"
	"github.com/gin-gonic/gin"
)

func EditPostForm(c *gin.Context) {
	idParam := c.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID")
		return
	}

	post, err := models.GetPostByID(postID)
	if err != nil {
		c.String(http.StatusNotFound, "Пост не найден")
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	// Приводим к int
	currentUserID := userIDInterface.(int)
	if post.UserID == nil || *post.UserID != currentUserID {
		c.String(http.StatusForbidden, "Нет прав на редактирование")
		return
	}

	c.HTML(http.StatusOK, "edit_post.html", gin.H{"post": post})
}

func UpdatePostHandler(c *gin.Context) {
	idParam := c.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID")
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	if title == "" || content == "" {
		c.String(http.StatusBadRequest, "Поля не могут быть пустыми")
		return
	}

	post, err := models.GetPostByID(postID)
	if err != nil {
		c.String(http.StatusNotFound, "Пост не найден")
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}
	currentUserID := userIDInterface.(int)
	if post.UserID == nil || *post.UserID != currentUserID {
		c.String(http.StatusForbidden, "Нет прав на редактирование")
		return
	}

	err = models.UpdatePost(postID, title, content)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка обновления поста")
		return
	}

	c.Redirect(http.StatusFound, "/post/"+idParam)
}

func ShowIndex(c *gin.Context) {
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Printf("Ошибка получения постов: %v", err)
		c.String(500, "Ошибка сервера")
		return
	}

	username, _ := c.Get("username") // будет "" если не задан
	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts":    posts,
		"username": username,
	})
}

// ShowPost отображает один пост по ID, преобразуя Markdown в HTML.
func ShowPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID поста")
		return
	}

	post, err := models.GetPostByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "Пост не найден")
		return
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(post.Content), &buf); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка рендеринга Markdown: %v", err)
		return
	}

	// Определяем, является ли текущий пользователь автором поста
	userIDVal, exists := c.Get("userID")
	isAuthor := false
	if exists && post.UserID != nil {
		if uid, ok := userIDVal.(int); ok && uid == *post.UserID {
			isAuthor = true
		}
	}

	// Получаем комментарии для поста
	comments, err := models.GetCommentsByPostID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка получения комментариев: %v", err)
		return
	}

	userIDValTwo, existsTwo := c.Get("userID")
	var currentUserID int
	if existsTwo {
		currentUserID = userIDValTwo.(int)
	}

	// Формируем для каждого комментария флаг, является ли текущий пользователь автором
	for i := range comments {
		comments[i].Editable = (comments[i].UserID == currentUserID)
	}

	// Извлекаем username, если он есть в контексте
	var currentUsername string
	if uname, ok := c.Get("username"); ok {
		if unameStr, ok := uname.(string); ok {
			currentUsername = unameStr
		}
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"post":        post,
		"contentHTML": buf.String(),
		"isAuthor":    isAuthor,
		"username":    currentUsername, // передаем username в шаблон
		"comments":    comments,
	})
}

// NewPostForm отображает HTML-форму для создания нового поста.
func NewPostForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", nil)
}

// CreatePost обрабатывает POST-запрос для создания нового поста.
func CreatePost(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.String(http.StatusBadRequest, "Заголовок и содержание не могут быть пустыми")
		return
	}

	// Получаем userID из контекста (middleware JWT должен его установить)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}
	userID := userIDInterface.(int)

	if err := models.CreatePost(userID, title, content); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка создания поста: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}
