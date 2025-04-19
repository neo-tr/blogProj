package main

import (
	"blogProj/db"
	"blogProj/handlers"
	"blogProj/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	// Инициализация бд
	db.Init()

	// Отложим закрытие при нормальном завершении
	defer db.Close()

	// Создаем шаблоны с функцией safeHTML и форматированием даты
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"formatDate": func(t interface{}) string {
			// Функция принимает t как значение типа time.Time или *time.Time
			switch tt := t.(type) {
			case time.Time:
				return tt.Format("2006-01-02 15:04")
			case *time.Time:
				if tt != nil {
					return tt.Format("2006-01-02 15:04")
				}
				return ""
			default:
				return ""
			}
		},
	}).ParseGlob("templates/*.html"))

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "./static")

	// Открытые маршруты
	r.GET("/", middleware.OptionalJWTAuthMiddleware(), handlers.ShowIndex)
	r.GET("/post/:id", middleware.OptionalJWTAuthMiddleware(), handlers.ShowPost)
	r.GET("/new", handlers.NewPostForm)
	r.GET("/register", func(c *gin.Context) { c.HTML(200, "register.html", nil) })
	r.POST("/register", handlers.RegisterHandler)
	r.GET("/login", func(c *gin.Context) { c.HTML(200, "login.html", nil) })
	r.POST("/login", handlers.LoginHandler)
	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	})

	// Защищенные маршруты (требуют JWT)
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/post/:id/edit", handlers.EditPostForm)
		auth.POST("/post/:id/edit", handlers.UpdatePostHandler)
		auth.POST("/post/:id/comment", handlers.AddCommentHandler)
		auth.GET("/comment/:commentID/edit", handlers.EditCommentForm)
		auth.POST("/comment/:commentID/edit", handlers.UpdateCommentHandler)
		auth.POST("/new", handlers.CreatePost)
		auth.POST("/comment/:commentID/delete", handlers.DeleteCommentHandler)
	}

	// Ловим сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Получен сигнал завершения, закрываем базу данных...")
		db.Close()
		os.Exit(0)
	}()

	// Запускаем сервер
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
