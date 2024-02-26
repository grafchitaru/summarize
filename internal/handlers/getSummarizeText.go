package handlers

import (
	"github.com/grafchitaru/summarize/internal/config"
	"net/http"
)

func GetSummarizeText(ctx config.HandlerContext, res http.ResponseWriter) {
	/*
		Получение саммаризированного текста
		Хендлер: GET /api/user/summarize/{id}.
		Хендлер доступен только аутентифицированным пользователям.


		Формат запроса:
		GET /api/user/summarize/{id} HTTP/1.1
		Content-Type: application/json
		Возможные коды ответа:
		200 — возвращается саммаризированный текст;
		401 — пользователь не аутентифицирован;
		500 — внутренняя ошибка сервера.
		Здесь наверное еще бы 404 добавить? И статус запроса в ответ
	*/
}
