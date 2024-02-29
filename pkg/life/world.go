package life

import (
	"errors"
	"math/rand"
	"os"
	"time"
)

type World struct {
	Height int // высота сетки
	Width  int // ширина сетки
	Cells  [][]bool
}

// NewWorld выделяет память под сетку
func NewWorld(height, width int) (*World, error) {
	if height <= 0 || width <= 0 {
		return nil, errors.New("Размеры не могут быть отрицательны")
	}
	// создаём тип World с количеством слайсов hight (количество строк)
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width) // создаём новый слайс в каждой строке
	}
	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}, nil
}

// Преобразование bool в int для вычисления количества соседей
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// KnowNeighbors ведёт подсчёт соседей на торе(с учетом рамок за выхождения).
func (w *World) KnowNeighbors(x, y int) int {
	// Если вышли за границы - сосед тор
	x = (x + w.Width) % w.Width
	y = (y + w.Height) % w.Height
	// Если со
	return boolToInt(w.Cells[y][x])
}

// Neighbors вычисляет количество соседей на торе(с учетом рамок за выхождения).
func (w *World) Neighbors(x int, y int) int {
	return w.KnowNeighbors(x-1, y-1) + w.KnowNeighbors(x, y-1) + w.KnowNeighbors(x+1, y-1) + w.KnowNeighbors(x+1, y) + w.KnowNeighbors(x+1, y+1) + w.KnowNeighbors(x, y+1) + w.KnowNeighbors(x-1, y+1) + w.KnowNeighbors(x-1, y)
}

func (w *World) Next(x, y int) bool {
	n := w.Neighbors(x, y)       // получим количество живых соседей
	alive := w.Cells[y][x]       // текущее состояние клетки
	if n < 4 && n > 1 && alive { // если соседей двое или трое, а клетка жива
		return true // то следующее состояние — жива
	}
	if n == 3 && !alive { // если клетка мертва, но у неё трое соседей
		return true // клетка оживает
	}

	return false // в любых других случаях — клетка мертва
}

func NextState(oldWorld, newWorld *World) {
	// переберём все клетки, чтобы понять, в каком они состоянии
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			// для каждой клетки получим новое состояние
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}
}

// Seed заполнит сетку живыми клетками в случайном порядке
func (w *World) Seed() {
	// снова переберём все клетки
	for _, row := range w.Cells {
		for i := range row {
			//rand.Intn(10) возвращает случайное число из диапазона	от 0 до 9
			if rand.Intn(10) == 1 {
				row[i] = true
			}
		}
	}
}

// SaveState сохранит текущего состояния сетки в файл
func (w *World) SaveState(filename string) error {
	// Создаём новый файл
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()
	// переберём все клетки
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			// _, err := f.WriteString(string(boolToInt(w.Cells[i][j])))
			_, err := f.WriteString(func() string {
				if w.Cells[i][j] {
					return "1"
				}
				return "0"
			}())
			if err != nil {
				return err
			}
		}
		_, err := f.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// RandInit заполняет поля на указанное число процентов
func (w *World) RandInit(percentage int) {
	// Количество живых клеток
	numAlive := percentage * w.Height * w.Width / 100
	// Заполним живыми первые клетки
	w.fillAlive(numAlive)
	// Получаем рандомные числа
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// Рандомно меняем местами
	for i := 0; i < w.Height*w.Width; i++ {
		randRowLeft := r.Intn(w.Width)
		randColLeft := r.Intn(w.Height)
		randRowRight := r.Intn(w.Width)
		randColRight := r.Intn(w.Height)

		w.Cells[randRowLeft][randColLeft] = w.Cells[randRowRight][randColRight]
	}
}

func (w *World) fillAlive(num int) {
	aliveCount := 0
	for j, row := range w.Cells {
		for k := range row {
			w.Cells[j][k] = true
			aliveCount++
			if aliveCount == num {

				return
			}
		}
	}
}

// Вывод на экран в виде цветных клеток
func (w *World) String() string {
	s := ""
	deadSquare := "\xF0\x9F\x9F\xAB"
	liveSquare := "\xF0\x9F\x9F\xA9"
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			if w.Cells[i][j] {
				s += liveSquare
			} else {
				s += deadSquare
			}
		}
		if i != w.Height-1 {
			s += "\n"
		}
	}
	return s
}
