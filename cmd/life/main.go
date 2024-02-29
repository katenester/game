package main

import (
	"context"
	"github.com/katenester/game/internal/application"
	"os"
)

func main() {
	// Создаём родительский контекст.
	ctx := context.Background()
	// Exit приводит к завершению программы с заданным кодом(0 - успешно, иначе - ошибка).
	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) int {
	cfg := application.Config{
		Width:  100,
		Height: 100,
	}
	app := application.New(cfg)

	return app.Run(ctx)
}
