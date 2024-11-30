// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"application-design-test/cmd/hotel-booking/setup"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	router := chi.NewRouter()

	repos := setup.InitRepositories()
	services := setup.InitServices(repos)
	//r.Use(middleware.Recoverer)
	setup.InitHandlers(router, services)

	slog.Info("Server listening localhost:80")
	if err := http.ListenAndServe(":80", router); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
