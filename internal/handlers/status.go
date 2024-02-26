package handlers

import (
	"github.com/grafchitaru/summarize/internal/config"
	"net/http"
)

func Status(ctx config.HandlerContext, res http.ResponseWriter) {
	/*
		Получение текущего статуса пользователя - сколько еще есть запросов/токенов
		Хендлер: GET /api/user/status.
		Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о доступном пользователю количества запросов и токенов.
		Формат запроса:
		Скопировать код
		GET /api/user/status HTTP/1.1
		Content-Length: 0
		Возможные коды ответа:
		200 — успешная обработка запроса.
		  Формат ответа:
		Скопировать код
		 200 OK HTTP/1.1
		  Content-Type: application/json
		  ...

		  {
		      "count": 1000,
		      "tokens": 100000
		  }

		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.

	*/
}
