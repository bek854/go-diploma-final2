package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "02.01.2006"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %w", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила повторения")
	}

	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("не указан интервал для ежедневного повторения")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный интервал для ежедневного повторения: %w", err)
		}

		if days < 1 || days > 400 {
			return "", fmt.Errorf("интервал должен составлять от 1 до 400 дней")
		}

		for {
			t = t.AddDate(0, 0, days)
			if t.After(now) {
				break
			}
		}

		return t.Format(dateFormat), nil

	case "y":
		for {
			t = t.AddDate(1, 0, 0)
			if t.After(now) {
				break
			}
		}
		return t.Format(dateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат даты 'now'", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}
