package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DroneCount != 100 {
		t.Errorf("Expected DroneCount=100, got %d", cfg.DroneCount)
	}
	if cfg.MinDistance != 5.0 {
		t.Errorf("Expected MinDistance=5.0, got %f", cfg.MinDistance)
	}
	if cfg.MaxDistance != 10.0 {
		t.Errorf("Expected MaxDistance=10.0, got %f", cfg.MaxDistance)
	}
	if cfg.UpdateInterval != 2*time.Second {
		t.Errorf("Expected UpdateInterval=2s, got %v", cfg.UpdateInterval)
	}
	if cfg.SmoothingFactor != 0.2 {
		t.Errorf("Expected SmoothingFactor=0.2, got %f", cfg.SmoothingFactor)
	}
	if cfg.Debug != false {
		t.Error("Expected Debug=false")
	}
	if cfg.VisualizationEnabled != true {
		t.Error("Expected VisualizationEnabled=true")
	}
	if cfg.WindowWidth != 800 {
		t.Errorf("Expected WindowWidth=800, got %d", cfg.WindowWidth)
	}
	if cfg.WindowHeight != 600 {
		t.Errorf("Expected WindowHeight=600, got %d", cfg.WindowHeight)
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	cfg, err := LoadConfig("nonexistent.json")
	if err != nil {
		t.Fatalf("Expected no error for nonexistent file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("Expected config to be returned")
	}
	if cfg.DroneCount != 100 {
		t.Errorf("Expected default DroneCount, got %d", cfg.DroneCount)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(tmpFile, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err := LoadConfig(tmpFile)
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestLoadConfigValid(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.json")

	jsonContent := `{
		"drone_count": 50,
		"min_distance": 3.0,
		"max_distance": 8.0,
		"debug": true
	}`

	if err := os.WriteFile(tmpFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.DroneCount != 50 {
		t.Errorf("Expected DroneCount=50, got %d", cfg.DroneCount)
	}
	if cfg.MinDistance != 3.0 {
		t.Errorf("Expected MinDistance=3.0, got %f", cfg.MinDistance)
	}
	if cfg.Debug != true {
		t.Error("Expected Debug=true")
	}
}

func TestConfigSave(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.json")

	cfg := DefaultConfig()
	cfg.DroneCount = 75

	err := cfg.Save(tmpFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Проверяем, что файл существует
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Загружаем и проверяем
	loadedCfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedCfg.DroneCount != 75 {
		t.Errorf("Expected DroneCount=75, got %d", loadedCfg.DroneCount)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		wantErr  bool
		checkFunc func(*Config) bool
	}{
		{
			name:   "ValidConfig",
			config: DefaultConfig(),
			wantErr: false,
			checkFunc: func(c *Config) bool { return true },
		},
		{
			name: "ZeroDroneCount",
			config: &Config{DroneCount: 0},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.DroneCount == 100 },
		},
		{
			name: "NegativeMinDistance",
			config: &Config{MinDistance: -1.0},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.MinDistance == 5.0 },
		},
		{
			name: "MaxDistanceLessThanMin",
			config: &Config{MinDistance: 10.0, MaxDistance: 5.0},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.MaxDistance == 15.0 },
		},
		{
			name: "InvalidSmoothingFactor",
			config: &Config{SmoothingFactor: 1.5},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.SmoothingFactor == 0.2 },
		},
		{
			name: "ZeroWindowWidth",
			config: &Config{WindowWidth: 0},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.WindowWidth == 800 },
		},
		{
			name: "SmallWindowWidth",
			config: &Config{WindowWidth: 100},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.WindowWidth == 320 },
		},
		{
			name: "SmallWindowHeight",
			config: &Config{WindowHeight: 100},
			wantErr: false,
			checkFunc: func(c *Config) bool { return c.WindowHeight == 240 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.checkFunc(tt.config) {
				t.Error("Config validation did not correct values as expected")
			}
		})
	}
}

func TestLoadEnvConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, ".env")

	envContent := `DRONE_COUNT=50
MIN_DISTANCE=3.0
MAX_DISTANCE=8.0
DEBUG=true
VISUALIZATION_ENABLED=false
WINDOW_WIDTH=1024
WINDOW_HEIGHT=768
`

	if err := os.WriteFile(tmpFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Очищаем окружение перед тестом
	_ = godotenv.Overload(tmpFile)
	defer func() {
		// Очищаем переменные после теста
		for _, key := range []string{"DRONE_COUNT", "MIN_DISTANCE", "MAX_DISTANCE", "DEBUG", 
			"VISUALIZATION_ENABLED", "WINDOW_WIDTH", "WINDOW_HEIGHT", "UPDATE_INTERVAL",
			"FORMATION_RADIUS", "MOVEMENT_VARIATION", "SMOOTHING_FACTOR"} {
			os.Unsetenv(key)
		}
	}()

	cfg, err := LoadEnvConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadEnvConfig failed: %v", err)
	}

	if cfg.DroneCount != 50 {
		t.Errorf("Expected DroneCount=50, got %d", cfg.DroneCount)
	}
	if cfg.MinDistance != 3.0 {
		t.Errorf("Expected MinDistance=3.0, got %f", cfg.MinDistance)
	}
	if cfg.MaxDistance != 8.0 {
		t.Errorf("Expected MaxDistance=8.0, got %f", cfg.MaxDistance)
	}
	if cfg.Debug != true {
		t.Error("Expected Debug=true")
	}
	if cfg.VisualizationEnabled != false {
		t.Error("Expected VisualizationEnabled=false")
	}
	if cfg.WindowWidth != 1024 {
		t.Errorf("Expected WindowWidth=1024, got %d", cfg.WindowWidth)
	}
	if cfg.WindowHeight != 768 {
		t.Errorf("Expected WindowHeight=768, got %d", cfg.WindowHeight)
	}
}

func TestLoadEnvConfigNotFound(t *testing.T) {
	_, err := LoadEnvConfig("nonexistent.env")
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestLoadEnvConfigPartial(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, ".env")

	// Только некоторые параметры
	envContent := `DRONE_COUNT=25
DEBUG=true
`

	if err := os.WriteFile(tmpFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Очищаем окружение перед тестом
	_ = godotenv.Overload(tmpFile)
	defer func() {
		for _, key := range []string{"DRONE_COUNT", "MIN_DISTANCE", "MAX_DISTANCE", "DEBUG", 
			"VISUALIZATION_ENABLED", "WINDOW_WIDTH", "WINDOW_HEIGHT", "UPDATE_INTERVAL",
			"FORMATION_RADIUS", "MOVEMENT_VARIATION", "SMOOTHING_FACTOR"} {
			os.Unsetenv(key)
		}
	}()

	cfg, err := LoadEnvConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadEnvConfig failed: %v", err)
	}

	if cfg.DroneCount != 25 {
		t.Errorf("Expected DroneCount=25, got %d", cfg.DroneCount)
	}
	if cfg.Debug != true {
		t.Error("Expected Debug=true")
	}
	// Остальные параметры должны быть по умолчанию
	if cfg.MinDistance != 5.0 {
		t.Errorf("Expected default MinDistance=5.0, got %f", cfg.MinDistance)
	}
}

func TestLoadEnvConfigWithComments(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, ".env")

	envContent := `# Это комментарий
DRONE_COUNT=30
# Ещё комментарий
MIN_DISTANCE=4.0

DEBUG=true
`

	if err := os.WriteFile(tmpFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Очищаем окружение перед тестом
	_ = godotenv.Overload(tmpFile)
	defer func() {
		for _, key := range []string{"DRONE_COUNT", "MIN_DISTANCE", "MAX_DISTANCE", "DEBUG", 
			"VISUALIZATION_ENABLED", "WINDOW_WIDTH", "WINDOW_HEIGHT", "UPDATE_INTERVAL",
			"FORMATION_RADIUS", "MOVEMENT_VARIATION", "SMOOTHING_FACTOR"} {
			os.Unsetenv(key)
		}
	}()

	cfg, err := LoadEnvConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadEnvConfig failed: %v", err)
	}

	if cfg.DroneCount != 30 {
		t.Errorf("Expected DroneCount=30, got %d", cfg.DroneCount)
	}
	if cfg.MinDistance != 4.0 {
		t.Errorf("Expected MinDistance=4.0, got %f", cfg.MinDistance)
	}
}

func TestLoadEnvConfigUpdateInterval(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, ".env")

	envContent := `UPDATE_INTERVAL=5000000000
`

	if err := os.WriteFile(tmpFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Очищаем окружение перед тестом
	_ = godotenv.Overload(tmpFile)
	defer func() {
		for _, key := range []string{"DRONE_COUNT", "MIN_DISTANCE", "MAX_DISTANCE", "DEBUG", 
			"VISUALIZATION_ENABLED", "WINDOW_WIDTH", "WINDOW_HEIGHT", "UPDATE_INTERVAL",
			"FORMATION_RADIUS", "MOVEMENT_VARIATION", "SMOOTHING_FACTOR"} {
			os.Unsetenv(key)
		}
	}()

	cfg, err := LoadEnvConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadEnvConfig failed: %v", err)
	}

	if cfg.UpdateInterval != 5*time.Second {
		t.Errorf("Expected UpdateInterval=5s, got %v", cfg.UpdateInterval)
	}
}
