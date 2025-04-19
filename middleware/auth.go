package middleware

import (
	"blogProj/db"
	"blogProj/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		// Логирование данных
		log.Printf("Пользователь авторизован: userID=%d, username=%s", claims.UserID, claims.Username)

		// Получаем никнейм из базы данных
		var nickname string
		err = db.Pool.QueryRow(c, "SELECT nickname FROM users WHERE id = $1", claims.UserID).Scan(&nickname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить никнейм пользователя"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("nickname", nickname)

		c.Next()
	}
}

func OptionalJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Next() // Если токена нет, продолжаем, но без авторизации
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			// Если токен неверный, можно просто продолжить или вернуть ошибку
			c.Next()
			return
		}

		// Устанавливаем данные пользователя в контекст
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("nickname", claims.Username)

		c.Next()
	}
}
