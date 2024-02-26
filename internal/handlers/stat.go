package handlers

import (
	"github.com/grafchitaru/summarize/internal/config"
	"net/http"
)

func Stat(ctx config.HandlerContext, res http.ResponseWriter) {
	/*
		Получение статистики по текущему пользователю
		Хендлер: GET /api/user/stat.
		Хендлер доступен только авторизованному пользователю. Данные должны группироваться по пользователю и статусу саммаризации.
		Доступные статусы саммаризации:
		Init — саммаризация текста инициирована;
		Proсess — текст в процессе саммаризации;
		Error — в результате саммаризации произошла ошибка;
		Success — саммаризация успешно завершена.
		Формат запроса:
		GET /api/user/stat HTTP/1.1
		Content-Length: 0
		Возможные коды ответа:
		200 — успешная обработка запроса.
		  Формат ответа:
		 200 OK HTTP/1.1
		  Content-Type: application/json
		  ...

		  [
		      {
		          "user_id": "9278923470",
		          "status": "Init",
		          "count": 1,
		          "tokens": "16000"
		      },
		      {
		          "user_id": "9278923470",
		          "status": "Error",
		          "count": 2,
		          "tokens": "16000"
		      }
		  ]

		204 — нет данных для ответа. Можно просто пустой массив вернуть?
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
}
