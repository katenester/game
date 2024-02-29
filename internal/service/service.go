package service

// сервис, который инициализирует и хранит состояния игры
import (
	"github.com/katenester/game/pkg/life"
	"math/rand"
	"time"
)

// LifeService хранит состояния
type LifeService struct {
	currentWorld *life.World
	nextWorld    *life.World
}

func New(height, width int) (*LifeService, error) {
	rand.NewSource(time.Now().UTC().UnixNano())
	// Получение объявленной структуры World с заданными размерами
	currentWorld, err := life.NewWorld(height, width)
	if err != nil {
		return nil, err
	}
	// Инициализация поля игры(заполнение живых клеток) предыдущего состояния ( для упрощения примера используется 40%)
	currentWorld.RandInit(40)
	// Получение объявленной структуры World с заданными размерами нового состояния
	newWorld, err := life.NewWorld(height, width)
	if err != nil {
		return nil, err
	}
	ls := LifeService{
		currentWorld: currentWorld,
		nextWorld:    newWorld,
	}

	return &ls, nil
}

// NewState получает очередное состояние игры
func (ls *LifeService) NewState() *life.World {
	life.NextState(ls.currentWorld, ls.nextWorld)

	ls.currentWorld = ls.nextWorld

	return ls.currentWorld
}
