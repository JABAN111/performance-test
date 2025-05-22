package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	// немного высшей математики: 40 req/min × 11 пользователей = 440 req/min
	rate := vegeta.Rate{Freq: 440, Per: time.Minute}
	dur := 10 * time.Minute

	var wg sync.WaitGroup
	wg.Add(3)

	for i := range 3 {
		go func() {
			defer wg.Done()
			port := 18081 + i
			url := fmt.Sprintf("http://localhost:%d/?token=495386773&user=-2104222730&config=%d", port, i+1)
			outputFile := fmt.Sprintf("perfomance-result/bins/load_config%d.bin", i)

			log.Info(fmt.Sprintf("Запуск Vegeta attack"), "Конфигурация", fmt.Sprintf("%d (порт %d)", i, port), "URL", url)

			targeter := vegeta.NewStaticTargeter(vegeta.Target{Method: "GET", URL: url})

			f, err := os.Create(outputFile)
			if err != nil {
				log.Error(fmt.Sprintf("Не удалось создать файл %s: %v", outputFile, err))
				return
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					panic(err)
				}
			}(f)

			enc := vegeta.NewEncoder(f)

			attacker := vegeta.NewAttacker()

			results := attacker.Attack(targeter, rate, dur, fmt.Sprintf("config-%d", i))

			var metrics vegeta.Metrics
			for res := range results {
				metrics.Add(res)
				if err := enc.Encode(res); err != nil {
					log.Error(fmt.Sprintf("Ошибка записи результата: %v", err))
				}
			}
			metrics.Close()

			log.Info(fmt.Sprintf("Итоги для конфигурации %d", i),
				"Requests:", metrics.Requests,
				"Success rate", fmt.Sprintf("%.2f%%", metrics.Success*100),
				"P95 latency", fmt.Sprintf("%s", metrics.Latencies.P95),
				"Errors", len(metrics.Errors),
				"Результат сохранён", outputFile)
		}()
	}
	wg.Wait()
}
