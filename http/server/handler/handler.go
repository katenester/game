package handler

import (
	"context"
	"encoding/json"
	"github.com/katenester/game/internal/service"
	"net/http"
)

// Decorator создаёт новый тип для добавления middleware к обработчикам
type Decorator func(http.Handler) http.Handler

// LifeStates является объектом для хранения состояния игры
type LifeStates struct {
	service.LifeService
}

func New(ctx context.Context, lifeService service.LifeService) (http.Handler, error) {
	serveMux := http.NewServeMux()
	lifeState := LifeStates{
		LifeService: lifeService,
	}
	// Регистрируем функцию обработчика для пути "/nextstate". Спойлер - вывод слайс байт в http.ResponseWriter
	serveMux.HandleFunc("/nextstate", lifeState.nextState)
	return serveMux, nil
}

// функция добавления middleware
func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

// Обработчик пути "/nextstate" для handler -LifeStates,
// Получение очередного состояния игры
func (ls *LifeStates) nextState(w http.ResponseWriter, r *http.Request) {
	// Выполнение хода игры, получение нового состояния
	worldState := ls.LifeService.NewState()
	// Записываем ответ(булевое поле-Cells) в w в json(в слайс байт) формат
	// NewEncoder возвращает новый кодировщик, который записывает данные в w.
	// Encode записывает JSON-кодировку v в поток, за которым следует символ новой строки.
	err := json.NewEncoder(w).Encode(worldState.Cells)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
