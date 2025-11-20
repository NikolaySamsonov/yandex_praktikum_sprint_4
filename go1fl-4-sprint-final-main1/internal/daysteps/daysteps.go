package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Неверный формат, дата должна иметь 2 части %d. получено", len(parts))
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка перевода шагов'%s': %w", stepsStr, err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("Шагов должно быть больше чем 0, получено %d", steps)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка перевода продолжительности'%s': %w", durationStr, err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Printf("Ошибка входных данных: %v\n", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, float64(weight), float64(height), duration)

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)

	return result
}
