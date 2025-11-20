package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "0", 0, fmt.Errorf("Неверный формат, данные должны иметь 3 части %d. получено", len(parts))
	}

	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "0", 0, fmt.Errorf("Ошибка перевода шагов '%s': %w", stepsStr, err)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "0", 0, fmt.Errorf("Ошибка перевода продолжительности '%s': %w", durationStr, err)
	}

	if duration <= 0 {
		return 0, "0", 0, fmt.Errorf("Продолжительность должна быть больше 0, имеем  %v", duration)
	}

	if activity == "" {
		return 0, "0", 0, fmt.Errorf("Отсутсвует активность")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	if steps <= 0 || height <= 0 {
		return 0.0
	}

	stepLength := height * stepLengthCoefficient

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	return distanceKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0.0
	}

	if steps <= 0 || height <= 0 {
		return 0.0
	}

	distKm := distance(steps, height)

	if distKm <= 0 {
		return 0.0
	}

	hours := duration.Hours()

	speedKmh := distKm / hours

	return speedKmh
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Шагов должно быть больше 0, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0, получено: %.2f кг", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0, получено: %.2f м", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0, получено: %v", duration)
	}

	// 2. Рассчитываем среднюю скорость в км/ч
	speedKmh := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speedKmh * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0, получено: %.2f кг", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0, получено: %.2f м", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0, получено: %v", duration)
	}

	speedKmh := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	baseCalories := (weight * speedKmh * durationInMinutes) / minInH

	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка данных тренировки:", err)
		return "", err
	}

	if weight <= 0 || height <= 0 {
		err := fmt.Errorf("некорректные данные: вес=%.2f кг, рост=%.2f м", weight, height)
		log.Println(err)
		return "", err
	}

	distanceKm := distance(steps, height)
	speedKmh := meanSpeed(steps, height, duration)

	var calories float64
	var normalizedActivity string
	activityNorm := strings.ToLower(strings.TrimSpace(activity))

	switch activityNorm {
	case "ходьба":
		normalizedActivity = "Ходьба"
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println("Ошибка расчёта калорий для ходьбы:", err)
			return "", err
		}

	case "бег":
		normalizedActivity = "Бег"
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println("Ошибка расчёта калорий для бега:", err)
			return "", err
		}

	default:
		err := fmt.Errorf("неизвестный тип тренировки: %q", activity)
		log.Println(err)
		return "", err
	}

	hours := duration.Hours()

	// 4. Формируем итоговую строку
	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		normalizedActivity,
		hours,
		distanceKm,
		speedKmh,
		calories,
	)

	return result, nil
}
