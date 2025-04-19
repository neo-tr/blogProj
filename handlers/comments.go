package handlers

import (
	"net/http"
	"strconv"

	"blogProj/models"
	"github.com/gin-gonic/gin"
)

// AddCommentHandler обрабатывает POST-запрос на добавление комментария к посту.
func AddCommentHandler(c *gin.Context) {
	// Получаем ID поста из URL-параметров
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID поста"})
		return
	}

	// Извлекаем userID из контекста (установленный JWTAuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	userID := userIDInterface.(int)

	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Комментарий не может быть пустым"})
		return
	}

	err = models.CreateComment(postID, userID, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания комментария: " + err.Error()})
		return
	}

	// Можно редиректнуть обратно на страницу поста.
	c.Redirect(http.StatusSeeOther, "/post/"+postIDStr)
}

func EditCommentForm(c *gin.Context) {
	commentIDStr := c.Param("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID комментария")
		return
	}

	comment, err := models.GetCommentByID(commentID)
	if err != nil {
		c.String(http.StatusNotFound, "Комментарий не найден")
		return
	}

	c.HTML(http.StatusOK, "edit_comment.html", gin.H{"comment": comment})
}

func UpdateCommentHandler(c *gin.Context) {
	// Получаем commentID из URL
	commentIDStr := c.Param("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID комментария")
		return
	}

	content := c.PostForm("content")
	if content == "" {
		c.String(http.StatusBadRequest, "Комментарий не может быть пустым")
		return
	}

	// Получаем комментарий из базы, чтобы узнать PostID и автора
	comment, err := models.GetCommentByID(commentID)
	if err != nil {
		c.String(http.StatusNotFound, "Комментарий не найден")
		return
	}

	// Обновляем комментарий
	err = models.UpdateComment(commentID, content)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка обновления комментария: %v", err)
		return
	}

	// Редиректим на страницу поста, используя comment.PostID
	redirectURL := "/post/" + strconv.Itoa(comment.PostID)
	c.Redirect(http.StatusSeeOther, redirectURL)
}

// DeleteCommentHandler удаляет комментарий, если текущий пользователь является его автором.
func DeleteCommentHandler(c *gin.Context) {
	// Получаем commentID из URL-параметров
	commentIDStr := c.Param("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID комментария")
		return
	}

	// Получаем комментарий по ID, чтобы проверить, кому он принадлежит
	comment, err := models.GetCommentByID(commentID)
	if err != nil {
		c.String(http.StatusNotFound, "Комментарий не найден")
		return
	}

	// Удаляем комментарий
	if err := models.DeleteComment(commentID); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка удаления комментария: %v", err)
		return
	}

	// Перенаправляем пользователя на страницу поста, к которому принадлежит комментарий
	redirectURL := "/post/" + strconv.Itoa(comment.PostID)
	c.Redirect(http.StatusSeeOther, redirectURL)
}
