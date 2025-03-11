package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// Server оборачивает стандартный http.Server и добавляет дополнительный метод AddRouter
type Server struct {
	*http.Server
}

// NewServer создает новый HTTP-сервер с заданным адресом
func NewServer(addr string) *Server {
	mux := chi.NewRouter()
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// AddRouter добавляет маршруты из другого роутера в текущий сервер
func (s *Server) AddRouter(router chi.Router, prefix string) {
	if mux, ok := s.Handler.(*chi.Mux); ok {
		mux.Mount(prefix, router)
	}
}

// Run запускает сервер с контекстом для правильной обработки завершения
func (s *Server) Run(ctx context.Context) error {
	// Канал для ловли ошибок
	errs := make(chan error)

	// Запуск сервера в горутине
	go func() {
		// Запуск ListenAndServe с контекстом
		errs <- s.ListenAndServe()
	}()

	// Ожидание завершения работы или тайм-аута
	select {
	case <-ctx.Done(): // Если контекст отменен, то пытаемся корректно завершить сервер
		// Применяем тайм-аут для graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %v", err)
		}
		return nil
	case err := <-errs: // Если сервер вернул ошибку
		return err
	}
}
