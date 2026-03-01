package config

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the simulation configuration
type Config struct {
	// DroneCount количество дочерних дронов в формировании
	DroneCount int `json:"drone_count"`

	// MinDistance минимальная дистанция от главного дрона
	MinDistance float64 `json:"min_distance"`

	// MaxDistance максимальная дистанция от главного дрона
	MaxDistance float64 `json:"max_distance"`

	// UpdateInterval интервал обновления позиции дронов
	UpdateInterval time.Duration `json:"update_interval"`

	// FormationRadius радиус формирования (базовый)
	FormationRadius float64 `json:"formation_radius"`

	// MovementVariation вариация движения для симуляции реалистичности
	MovementVariation float64 `json:"movement_variation"`

	// SmoothingFactor коэффициент сглаживания движения (0-1)
	SmoothingFactor float64 `json:"smoothing_factor"`

	// Debug режим отладки
	Debug bool `json:"debug"`

	// VisualizationEnabled включить визуализацию
	VisualizationEnabled bool `json:"visualization_enabled"`

	// WindowWidth ширина окна визуализации
	WindowWidth int `json:"window_width"`

	// WindowHeight высота окна визуализации
	WindowHeight int `json:"window_height"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		DroneCount:           100,
		MinDistance:          5.0,
		MaxDistance:          10.0,
		UpdateInterval:       2 * time.Second,
		FormationRadius:      7.5,
		MovementVariation:    0.5,
		SmoothingFactor:      0.2,
		Debug:                false,
		VisualizationEnabled: true,
		WindowWidth:          800,
		WindowHeight:         600,
	}
}

// LoadConfig загружает конфигурацию из файла или возвращает значения по умолчанию
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	if path == "" {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// LoadEnvConfig загружает конфигурацию из .env файла и применяет её к конфигурации
// Возвращает Config с применёнными значениями из .env или error если файл не найден
// Поддерживаемые переменные окружения:
//   - DRONE_COUNT (int)
//   - MIN_DISTANCE (float)
//   - MAX_DISTANCE (float)
//   - UPDATE_INTERVAL (int, наносекунды)
//   - FORMATION_RADIUS (float)
//   - MOVEMENT_VARIATION (float)
//   - SMOOTHING_FACTOR (float)
//   - DEBUG (bool)
//   - VISUALIZATION_ENABLED (bool)
//   - WINDOW_WIDTH (int)
//   - WINDOW_HEIGHT (int)
func LoadEnvConfig(path string) (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	config := DefaultConfig()

	// Читаем значения из переменных окружения
	if v := os.Getenv("DRONE_COUNT"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			config.DroneCount = val
		}
	}
	if v := os.Getenv("MIN_DISTANCE"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			config.MinDistance = val
		}
	}
	if v := os.Getenv("MAX_DISTANCE"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			config.MaxDistance = val
		}
	}
	if v := os.Getenv("UPDATE_INTERVAL"); v != "" {
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			config.UpdateInterval = time.Duration(val)
		}
	}
	if v := os.Getenv("FORMATION_RADIUS"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			config.FormationRadius = val
		}
	}
	if v := os.Getenv("MOVEMENT_VARIATION"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			config.MovementVariation = val
		}
	}
	if v := os.Getenv("SMOOTHING_FACTOR"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			config.SmoothingFactor = val
		}
	}
	if v := os.Getenv("DEBUG"); v != "" {
		config.Debug = v == "true"
	}
	if v := os.Getenv("VISUALIZATION_ENABLED"); v != "" {
		config.VisualizationEnabled = v == "true"
	}
	if v := os.Getenv("WINDOW_WIDTH"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			config.WindowWidth = val
		}
	}
	if v := os.Getenv("WINDOW_HEIGHT"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			config.WindowHeight = val
		}
	}

	return config, nil
}

// Save сохраняет конфигурацию в файл
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Validate проверяет валидность конфигурации
func (c *Config) Validate() error {
	if c.DroneCount <= 0 {
		c.DroneCount = 100
	}
	if c.MinDistance < 0 {
		c.MinDistance = 5.0
	}
	if c.MaxDistance < c.MinDistance {
		c.MaxDistance = c.MinDistance + 5.0
	}
	if c.UpdateInterval <= 0 {
		c.UpdateInterval = 2 * time.Second
	}
	if c.SmoothingFactor <= 0 || c.SmoothingFactor > 1 {
		c.SmoothingFactor = 0.2
	}
	// Валидация размеров окна
	if c.WindowWidth <= 0 {
		c.WindowWidth = 800
	}
	if c.WindowHeight <= 0 {
		c.WindowHeight = 600
	}
	// Ограничение минимальных размеров окна
	if c.WindowWidth < 320 {
		c.WindowWidth = 320
	}
	if c.WindowHeight < 240 {
		c.WindowHeight = 240
	}
	return nil
}
