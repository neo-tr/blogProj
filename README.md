# Markdown-блог 

Простое веб-приложение для публикации, редактирования постов и дополнительно удаление комментариев. Реализована авторизация, управление правами доступа, работа с базой данных и защита от XSS. 

## Возможности 

- Регистрация и вход пользователей (пароли хэшируются и не хранятся в чистом виде) 
- Создание, редактирование постов 
- Создание, редактирование и удаление комментариев 
- Защита от XSS 
- Хранение данных в PostgreSQL 
- Сессии пользователей через cookies - JWT 
- Простая архитектура на Go + HTML шаблоны

- Пример главнгой страницы 
![mainPage](img/mainPage.jpg) 
- Пример двух постов (изменение доступно только автору написанного поста или комментария) 
| ![post1](img/post1.jpg) | ![post2](img/post2.jpg) | 

- Хранение хэшей пароля 
![holdPass](img/holdPass.jpg)
