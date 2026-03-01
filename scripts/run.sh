#!/bin/bash

# Скрипт запуска симуляции дронов
# Использование: ./scripts/run.sh [OPTIONS]
#
# Опции:
#   -d, --debug       Режим отладки
#   -n, --no-vis      Без визуализации
#   -c, --config FILE Путь к конфигурации
#   -h, --help        Показать справку

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Параметры по умолчанию
DEBUG=""
NO_VIS=""
CONFIG=""

# Парсинг аргументов
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--debug)
            DEBUG="-debug"
            shift
            ;;
        -n|--no-vis)
            NO_VIS="-no-vis"
            shift
            ;;
        -c|--config)
            CONFIG="-config $2"
            shift 2
            ;;
        -h|--help)
            echo "Использование: $0 [OPTIONS]"
            echo ""
            echo "Опции:"
            echo "  -d, --debug       Режим отладки"
            echo "  -n, --no-vis      Без визуализации"
            echo "  -c, --config FILE Путь к конфигурации"
            echo "  -h, --help        Показать справку"
            exit 0
            ;;
        *)
            echo -e "${RED}Неизвестная опция: $1${NC}"
            echo "Используйте -h для справки"
            exit 1
            ;;
    esac
done

# Проверка наличия Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Ошибка: Go не найден. Установите Go 1.21+${NC}"
    exit 1
fi

# Вывод информации
echo -e "${GREEN}=== Drone Simulation ===${NC}"
echo -e "${YELLOW}Запуск симуляции...${NC}"
if [ -n "$DEBUG" ]; then
    echo -e "  Режим: ${YELLOW}Отладка${NC}"
fi
if [ -n "$NO_VIS" ]; then
    echo -e "  Визуализация: ${YELLOW}Отключена${NC}"
fi
if [ -n "$CONFIG" ]; then
    echo -e "  Конфигурация: ${YELLOW}$CONFIG${NC}"
fi
echo ""

# Запуск
go run cmd/app/main.go $DEBUG $NO_VIS $CONFIG
