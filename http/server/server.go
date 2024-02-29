package server

import (
	"context"
	"fmt"
	"github.com/katenester/game/http/server/handler"
	"github.com/katenester/game/internal/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// маршрутизация
func new(ctx context.Context, logger *zap.Logger, lifeService service.LifeService) (http.Handler, error) {
	// Регистрируем маршрут с обработкой. muxHandler стал объектом Handler с обработчиком
	muxHandler, err := handler.New(ctx, lifeService)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}
	// middleware для обработчиков
	muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

// Run c родительским контекстом, с специально настроенным логгиром, высотой и шириной по умолчанию(версия 1.0)
func Run(ctx context.Context, logger *zap.Logger, height, width int) (func(context.Context) error, error) {
	// сервис с игрой. Получение текущего и следующего состояния игры
	lifeService, err := service.New(height, width)
	if err != nil {
		return nil, err
	}
	// Регистрация маршрутизатора
	muxHandler, err := new(ctx, logger, *lifeService)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: ":8081", Handler: muxHandler}

	go func() {
		// Запускаем сервер
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe",
				zap.String("err", err.Error()))
		}
	}()
	// вернем функцию для завершения работы сервера
	return srv.Shutdown, nil
}

// middleware для логированя запросов
func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропуск запроса к следующему обработчику
			next.ServeHTTP(w, r)

			// Завершение логирования после выполнения запроса
			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
