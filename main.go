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

	rate := vegeta.Rate{Freq: 440, Per: time.Minute}
	dur := 5 * time.Minute

	var wg sync.WaitGroup
	wg.Add(3)

	for i := range 3 {
		go func() {
			defer wg.Done()
			port := 18080 + i
			url := fmt.Sprintf("http://localhost:%d/?token=495386773&user=-2104222730&config=%d", port, i)
			outputFile := fmt.Sprintf("perfomance-result/bins/load_config%d.bin", i)

			log.Info(fmt.Sprintf("=== Конфигурация %d (порт %d) ===", i, port))
			log.Info(fmt.Sprintf("URL: %s", url))
			log.Info(fmt.Sprintf("Запуск Vegeta attack..."))

			targeter := vegeta.NewStaticTargeter(vegeta.Target{Method: "GET", URL: url})

			f, err := os.Create(outputFile)
			if err != nil {
				log.Error(fmt.Sprintf("Не удалось создать файл %s: %v", outputFile, err))
				return
			}
			defer f.Close()

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

			log.Info(fmt.Sprintf("--- Итоги для конфигурации %d ---", i))
			log.Info(fmt.Sprintf("Requests: %d", metrics.Requests))
			log.Info(fmt.Sprintf("Success rate: %.2f%%", metrics.Success*100))
			log.Info(fmt.Sprintf("P95 latency: %s", metrics.Latencies.P95))
			log.Info(fmt.Sprintf("Errors: %d", len(metrics.Errors)))
			log.Info(fmt.Sprintf("Результат сохранён в %s", outputFile))
		}()
	}
	wg.Wait()
}
