package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ze674/EZLine/templates"
	"log"
	"net/http"
)

func main() {
	// Создаем роутер
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Главная страница
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Page(templates.Home()).Render(r.Context(), w)
	})

	// Статические файлы (если понадобятся)
	// fileServer := http.FileServer(http.Dir("./static"))
	// r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Запускаем сервер
	log.Println("Запуск сервера на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
