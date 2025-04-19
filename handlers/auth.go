package handlers

import (
	"blogProj/models"
	"blogProj/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")

	if username == "" || password == "" || nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля обязательны"})
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хэширования пароля"})
		return
	}

	err = models.CreateUser(username, string(hashedPassword), nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации: " + err.Error()})
		return
	}

	// После успешной регистрации редиректим на страницу логина с параметром success
	c.Redirect(http.StatusSeeOther, "/login?success=1")
}

func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Получаем параметр success из запроса
	success := c.Query("success")
	var successMessage string
	if success == "1" {
		successMessage = "Регистрация прошла успешно! Теперь войдите."
	}

	// Если GET-запрос (значит просто зашли на страницу, не отправляли форму)
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"SuccessMessage": successMessage,
			"ErrorMessage":   "",
		})
		return
	}

	// POST-логика — проверка пользователя
	user, err := models.GetUserByUsername(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"ErrorMessage":   "Неверный логин или пароль",
			"SuccessMessage": "",
		})
		return
	}

	// Генерация токена и установка cookie
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"ErrorMessage": "Ошибка генерации токена",
		})
		return
	}
	c.SetCookie("token", token, 3600*24, "/", "localhost", false, true)

	// Успешный вход → редирект на главную
	c.Redirect(http.StatusSeeOther, "/")
}
