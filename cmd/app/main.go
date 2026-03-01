package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"drone/internal/config"
	"drone/internal/services"
	"drone/internal/visualization"
)

func main() {
	// Парсинг флагов командной строки
	configPath := flag.String("config", "", "Путь к файлу конфигурации (JSON)")
	noVis := flag.Bool("no-vis", false, "Отключить визуализацию")
	debug := flag.Bool("debug", false, "Режим отладки")
	flag.Parse()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Применение флагов
	if *noVis {
		cfg.VisualizationEnabled = false
	}
	if *debug {
		cfg.Debug = true
	}

	// Валидация конфигурации
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка валидации конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Настройка логгера
	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	logger.Info("Запуск симуляции дронов",
		"drones", cfg.DroneCount,
		"visualization", cfg.VisualizationEnabled,
	)

	// Создание сервиса симуляции
	simService := services.NewSimulationService(cfg, logger)

	if err := simService.Init(); err != nil {
		logger.Error("Ошибка инициализации симуляции", "error", err)
		os.Exit(1)
	}

	if err := simService.Start(); err != nil {
		logger.Error("Ошибка запуска симуляции", "error", err)
		os.Exit(1)
	}

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов - в отдельной горутине
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал остановки...")
		cancel()
	}()

	fmt.Println("\nСимуляция запущена. Нажмите Ctrl+C для остановки.")

	// Запуск визуализации
	if cfg.VisualizationEnabled {
		// Lock OS thread для OpenGL - должно быть в главной горутине
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		visManager := visualization.NewManager(cfg, simService.GetMainDrone())
		if err := visManager.Init(); err != nil {
			logger.Error("Ошибка инициализации визуализации", "error", err)
			simService.Stop()
			os.Exit(1)
		}

		// Передаём симуляцию в визуализатор для управления дроном
		visManager.SetSimulation(simService)

		// Запуск рендеринга в главной горутине (блокирующий вызов)
		if err := visManager.Run(ctx); err != nil {
			logger.Error("Ошибка визуализации", "error", err)
		}

		if err := visManager.Close(); err != nil {
			logger.Warn("Ошибка при закрытии визуализации", "error", err)
		}
	} else {
		// Без визуализации - просто ждём сигнал
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
		}()
		wg.Wait()
	}

	// Graceful shutdown
	simService.Stop()
	logger.Info("Симуляция завершена")
}
