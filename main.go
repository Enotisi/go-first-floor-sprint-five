package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // Вид тренировки
	Action       int           // Кол-во шагов/гребков
	LenStep      float64       // Длинна одного шага/гребка (в м)
	Duration     time.Duration // Продолжительность тренировки
	Weight       float64       // Вес пользователя
}

// distance возвращает дистанцию, которую преодолел пользователь
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	if t.Duration.Hours() == 0 {
		return 0
	}

	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке.
// Пока возвращает 0, так как этот метод будет переопределяться для каждого типа тренировки.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string        // Вид тренировки
	Duration     time.Duration // Длительность тренировки
	Distance     float64       // Рсстояние
	Speed        float64       // Средняя скорость
	Calories     float64       // Кол-во калорий
}

// TrainingInfo возвращает труктуру InfoMessage, в которой хранится вся информация о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потраченных килокалорий при беге.
const (
	CaloriesMeanSpeedMultiplier = 18   // множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // коэффициент изменения средней скорости
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возввращает количество потраченных килокалория при беге.
// Это переопределенный метод Calories() из Training.
func (r Running) Calories() float64 {
	speed := r.meanSpeed()
	duration := r.Duration.Hours()
	return (CaloriesMeanSpeedMultiplier * speed * CaloriesMeanSpeedShift) * r.Weight / MInKm * duration * MinInHours
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

// Константы для расчета потраченных килокалорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // коэффициент для роста
	KmHInMsec                     = 0.278 // коэффициент для перевода км/ч в м/с
)

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64 // Рост пользователя в метрах
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
// Это переопределенный метод Calories() из Training.
func (w Walking) Calories() float64 {
	speedMinSec := w.meanSpeed() * KmHInMsec
	duration := w.Duration.Hours()
	return (CaloriesWeightMultiplier*w.Weight + (math.Pow(speedMinSec, 2)/w.Height)*CaloriesSpeedHeightMultiplier*w.Weight) * duration * MinInHours
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

// Константы для расчета потраченных килокалорий при плавании.
const (
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // множитель веса пользователя
)

// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool int // длинна бассейна в метрах
	CountPool  int // Кол-во пересечений бассейна
}

// meanSpeed возвращает среднюю скорость при плавании.
// Это переопределенный метод meanSpeed() из Training.
func (s Swimming) meanSpeed() float64 {
	duration := s.Duration.Hours()
	return float64(s.LengthPool*s.CountPool) / MInKm / duration
}

// Calories возвращает количество калорий, потраченных при плавании.
// Это переопределенный метод Calories() из Training.
func (s Swimming) Calories() float64 {
	speed := s.meanSpeed()
	duration := s.Duration.Hours()

	return (speed + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * duration
}

// TrainingInfo returns info about swimming training.
// Это переопределенный метод TrainingInfo() из Training.
func (s Swimming) TrainingInfo() InfoMessage {
	return s.Training.TrainingInfo()
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	// получите количество затраченных калорий
	calories := training.Calories()

	// получите информацию о тренировке
	info := training.TrainingInfo()
	// добавьте полученные калории в структуру с информацией о тренировке
	info.Calories = calories

	return fmt.Sprint(info)
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
