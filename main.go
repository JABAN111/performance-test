package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func attackAllConfig(log *slog.Logger, rate vegeta.Rate, dur time.Duration) {
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
				"Mean", metrics.Latencies.Mean,
				"Success rate", fmt.Sprintf("%.2f%%", metrics.Success*100),
				"P95 latency", fmt.Sprintf("%s", metrics.Latencies.P95),
				"Errors", len(metrics.Errors),
				"Результат сохранён", outputFile)
		}()
	}
	wg.Wait()
}

func attackSecondConfig(log *slog.Logger, rate vegeta.Rate, dur time.Duration) (int, time.Duration) {
	port := 18081 + 2
	url := fmt.Sprintf("http://localhost:%d/?token=495386773&user=-2104222730&config=%d", port, 3)

	log.Info(fmt.Sprintf("Запуск Vegeta attack"), "Конфигурация", fmt.Sprintf("%d (порт %d)", 2, port), "URL", url)

	targeter := vegeta.NewStaticTargeter(vegeta.Target{Method: "GET", URL: url})
	attacker := vegeta.NewAttacker()
	results := attacker.Attack(targeter, rate, dur, fmt.Sprintf("config-%d", 2))
	var metrics vegeta.Metrics
	for res := range results {
		metrics.Add(res)
	}
	metrics.Close()

	log.Info(fmt.Sprintf("Итоги для конфигурации %d", 3),
		"Frequency", rate.Freq,
		"Mean", metrics.Latencies.Mean,
		"Requests:", metrics.Requests,
		"Success rate", fmt.Sprintf("%.2f%%", metrics.Success*100),
		"P95 latency", fmt.Sprintf("%s", metrics.Latencies.P95),
		"Errors", len(metrics.Errors))
	return rate.Freq, metrics.Latencies.Mean
}

func AttackSecondConfig(log *slog.Logger, dur time.Duration) {
	users := 11
	req := 40

	rates := []vegeta.Rate{
		{Freq: users * req, Per: time.Minute},
		{Freq: users * req * 2, Per: time.Minute},
		{Freq: users * req * 4, Per: time.Minute},
		{Freq: users * req * 8, Per: time.Minute},
		{Freq: users * req * 16, Per: time.Minute},
		{Freq: users * req * 32, Per: time.Minute},
	}
	type freqMean struct {
		freq int
		mean time.Duration
	}

	var results []freqMean

	for _, r := range rates {
		freq, mean := attackSecondConfig(log, r, dur)
		results = append(results, freqMean{freq, mean})
	}
	log.Info("start reading results")
	for _, res := range results {
		log.Info("result got", "freq", res.freq, "mean", res.mean)
	}
}

func AttackAllConfigs(log *slog.Logger, rate vegeta.Rate, dur time.Duration) {
	attackAllConfig(log, rate, dur)
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	users := 11
	req := 40
	dur := time.Hour
	AttackSecondConfig(log, dur)

	rate := vegeta.Rate{Freq: users * req, Per: time.Minute}
	AttackAllConfigs(log, rate, dur)
}
