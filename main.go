package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"regexp"
	"syscall"
	"time"
	"unsafe"

	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var config Config

type Config struct {
	ButtonText             string       `json:"button_text"`
	TelegramWindowName     string       `json:"telegram_window_name"`
	Colors                 []Color      `json:"colors"`
	ClickDelayMin          int          `json:"click_delay_min"`
	ClickDelayMax          int          `json:"click_delay_max"`
	ButtonRelativePosition ButtonStruct `json:"button_relative_position"`
}

type Color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type ButtonStruct struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	setWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	getMessage          = user32.NewProc("GetMessageW")
	getWindowRectProc   = user32.NewProc("GetWindowRect")
	paused              bool
	hookHandle          HHOOK
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type HHOOK uintptr
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X int32
	Y int32
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println(`
	BLUM Bot v1.03 - Auto Clicker
                                     by @Firsim

	Telegram group: https://t.me/projectby

	Controls: Press 'P' to pause/resume
	Make sure Telegram is running and BLUM bot is active

	Example run with 5 games: blum_bot.exe -play 5

	`)

	// Загрузка конфигурации
	if err := loadConfig(); err != nil {
		fmt.Println("Configuration loading error:", err)
		return
	}

	// Обработка аргументов командной строки
	playCount := flag.Int("play", 0, "Number of games to run (0 = 999 games)")
	flag.Parse()

	// Если количество игр не указано или равно 0, устанавливаем значение по умолчанию (999)
	maxGames := *playCount
	if maxGames == 0 {
		maxGames = 999
	}

	fmt.Printf("Games planned: %d\n", maxGames)

	// Инициализация счетчика игр
	gamesPlayed := 0

	// Проверка, запущена ли программа от имени администратора
	if !isRunAsAdmin() {
		fmt.Println("Error: The program must be run as an administrator!")
		fmt.Println("Run the program with administrator privileges.")

		// Запускаем бесконечный цикл, что бы окно не закрывалось сразу
		for {
			time.Sleep(1 * time.Second) // Программа "засыпает" на 1 секунду
		}
	}

	// Поиск окна "TelegramWindowName"
	blumHwnd := findTelegramWindow(config.TelegramWindowName)
	if blumHwnd == 0 {
		fmt.Println("Window '%s' not found. Please start Blum first.", config.TelegramWindowName)

		// Запускаем бесконечный цикл, что бы окно не закрывалось сразу
		for {
			time.Sleep(1 * time.Second) // Программа "засыпает" на 1 секунду
		}
	}

	fmt.Printf("Found window '%s': %v\n", config.TelegramWindowName, blumHwnd)

	fmt.Printf("Starting work (%s)\n", config.TelegramWindowName)

	// Запуск обработчика паузы
	go handlePause()

	// Паузу для того, что бы успела запустить горутина handlePause
	time.Sleep(time.Second)

	// Просим пользователя запустить первую игру вручную
	fmt.Println("")
	fmt.Println("!!! Please start the first game manually and unpause!")
	fmt.Println("")

	// Устанавливаем паузу, что бы пользователь успел запустить первую игру врусную
	paused = !paused

	// Основной цикл программы
	for {
		// Проверка паузы
		if paused {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Проверяем, достигнут ли лимит игр
		if checkGamesCompleted(gamesPlayed, maxGames, &paused) {
			time.Sleep(1 * time.Second)
			continue
		}

		// Запуск кликера для игры
		gameClicker(blumHwnd)

		// Увеличиваем счетчик игр
		gamesPlayed++
		fmt.Printf("Game completed. Total played: %d out of %d games, %d games remaining.\n", gamesPlayed, maxGames, int(maxGames-gamesPlayed))

		// Проверяем, достигнут ли лимит игр
		if checkGamesCompleted(gamesPlayed, maxGames, &paused) {
			time.Sleep(1 * time.Second)
			continue
		}

		// Поиск кнопки "Играть"
		for {
			// Вычисляем глобальные координаты кнопки
			buttonPos, err := calculateButtonPosition(blumHwnd, &config)
			if err != nil {
				fmt.Println(err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Нажатие на кнопку "Играть" через случайный интервал времени
			randomDelay := time.Duration(rand.Intn(config.ClickDelayMax-config.ClickDelayMin)+config.ClickDelayMin) * time.Second
			fmt.Printf("Button '%s' found at coordinates (%d, %d). Clicking in %v seconds.\n",
				config.ButtonText, buttonPos.X, buttonPos.Y, int(randomDelay/time.Second))
			time.Sleep(randomDelay)

			// Проверка паузы
			for {
				if paused {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}

			// Клик по кнопке
			click(buttonPos.X+rand.Intn(5), buttonPos.Y+rand.Intn(5)) // Добавляем случайное смещение
			// Еще один клик
			click(buttonPos.X+rand.Intn(5), buttonPos.Y+rand.Intn(5)) // Добавляем случайное смещение

			// Пауза перед следующей итерацией
			time.Sleep(500 * time.Millisecond)
			break
		}
	}
}

// Проверка запуска от имени администратора
func isRunAsAdmin() bool {
	// Загружаем библиотеку shell32.dll
	shell32 := syscall.NewLazyDLL("shell32.dll")
	procIsUserAnAdmin := shell32.NewProc("IsUserAnAdmin")

	// Вызываем функцию IsUserAnAdmin
	ret, _, _ := procIsUserAnAdmin.Call()
	return ret != 0
}

// Проверка на количество оставшихся игр
func checkGamesCompleted(gamesPlayed, maxGames int, paused *bool) bool {
	if gamesPlayed >= maxGames {
		fmt.Println("All games completed. The program is paused.")
		*paused = !*paused
		for {
			// Блокируем выполнение основного цикла, пока не будет снята пауза
			if !*paused {
				fmt.Println("The program is resumed. Continuing the game...")
				return true // Возобновляем выполнение программы
			}
			time.Sleep(1 * time.Second)
		}
	}
	return false
}

// Вычесляем позицию кнопки "Играть"
func calculateButtonPosition(hwnd uintptr, config *Config) (*image.Point, error) {
	// Получаем координаты окна
	rect := getClientRect(hwnd)
	if rect.Empty() {
		return nil, fmt.Errorf("не удалось получить координаты окна")
	}

	// Вычисляем ширину и высоту окна
	windowWidth := rect.Dx()
	windowHeight := rect.Dy()

	// Вычисляем процентные доли (от 0 до 1)
	percentX := float64(config.ButtonRelativePosition.X) / 100
	percentY := float64(config.ButtonRelativePosition.Y) / 100

	// Вычисляем относительные координаты внутри окна
	relativeX := int(float64(windowWidth) * percentX)
	relativeY := int(float64(windowHeight) * percentY)

	// Вычисляем глобальные координаты кнопки
	globalX := rect.Min.X + relativeX
	globalY := rect.Min.Y + relativeY

	// fmt.Printf("windowWidth: %d | windowHeight: %d\n", windowWidth, windowHeight)
	// fmt.Printf("percentX: %.2f | percentY: %.2f\n", percentX, percentY)
	// fmt.Printf("relativeX: %d | relativeY: %d\n", relativeX, relativeY)
	// fmt.Printf("globalX: %d | globalY: %d\n", globalX, globalY)

	// Проверяем, не выходят ли координаты за пределы окна
	if globalX < rect.Min.X || globalX > rect.Max.X || globalY < rect.Min.Y || globalY > rect.Max.Y {
		return nil, fmt.Errorf("координаты кнопки выходят за пределы окна")
	}

	return &image.Point{X: globalX, Y: globalY}, nil
}

// Загрузка конфигурации
func loadConfig() error {
	// Открываем файл конфигурации
	file, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Декодируем JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Валидация полей конфигурации
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// Валидация полей конфигурации
func validateConfig(cfg Config) error {
	// Проверка button_text
	if len(cfg.ButtonText) < 3 || len(cfg.ButtonText) > 100 {
		return fmt.Errorf("button_text must be between 3 and 100 characters")
	}
	if !isAlphanumeric(cfg.ButtonText) {
		return fmt.Errorf("button_text must contain only letters and numbers")
	}

	// Проверка telegram_window_name
	if len(cfg.TelegramWindowName) < 5 || len(cfg.TelegramWindowName) > 100 {
		return fmt.Errorf("telegram_window_name must be between 5 and 100 characters")
	}

	// Проверка colors
	if len(cfg.Colors) > 10 {
		return fmt.Errorf("colors array must contain no more than 10 colors")
	}

	for i, color := range cfg.Colors {
		if color.Red < 0 || color.Red > 255 ||
			color.Green < 0 || color.Green > 255 ||
			color.Blue < 0 || color.Blue > 255 {
			return fmt.Errorf("invalid color at index %d: RGB values must be between 0 and 255", i)
		}
	}

	// Проверка click_delay_min
	if cfg.ClickDelayMin < 0 || cfg.ClickDelayMin > 100 {
		return fmt.Errorf("click_delay_min must be between 0 and 100")
	}

	// Проверка click_delay_max
	if cfg.ClickDelayMax < 0 || cfg.ClickDelayMax > 100 {
		return fmt.Errorf("click_delay_max must be between 0 and 100")
	}

	// Проверка button_relative_position (x и y)
	if cfg.ButtonRelativePosition.X < 0 || cfg.ButtonRelativePosition.X > 100 ||
		cfg.ButtonRelativePosition.Y < 0 || cfg.ButtonRelativePosition.Y > 100 {
		return fmt.Errorf("button_relative_position x and y must be between 0 and 100")
	}

	return nil
}

// Регулярное выражение для проверки букв и цифр
func isAlphanumeric(s string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, s)
	return matched
}

// Поиск окна Telegram
func findTelegramWindow(windowName string) uintptr {
	title, _ := windows.UTF16PtrFromString(windowName)
	hwnd := win.FindWindow(nil, title)

	if hwnd == 0 {
		return 0
	}

	// Принудительно активируем окно
	// win.SetForegroundWindow(win.HWND(hwnd))
	// time.Sleep(500 * time.Millisecond) // Даем время окну активироваться

	return uintptr(hwnd)
}

// Установка кнопки для паузы
func hookCallback(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 && wParam == WM_KEYDOWN {
		kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		if kbdStruct.VkCode == 0x50 { // 0x50 - это код клавиши "P"
			paused = !paused
			if paused {
				fmt.Println("The program is paused. Press 'P' again to continue.")
			} else {
				fmt.Println("The program has resumed.")
			}
		}
	}
	ret, _, _ := callNextHookEx.Call(uintptr(0), uintptr(nCode), wParam, lParam)
	return ret
}

// Настройка захвата кнопки для паузы
func setupGlobalHook() error {
	hook, _, err := setWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		windows.NewCallback(hookCallback),
		uintptr(0),
		uintptr(0),
	)
	if err.(syscall.Errno) != 0 {
		return err
	}
	hookHandle = HHOOK(hook)
	return nil
}

// Очистить Hook
func cleanupHook() {
	if hookHandle != 0 {
		unhookWindowsHookEx.Call(uintptr(hookHandle))
	}
}

// Цикл обработки сообщений Windows
func messageLoop() {
	var msg MSG
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg)), uintptr(0), 0, 0)
		if ret == 0 {
			break
		}
	}
}

// Функция паузы
func handlePause() {
	fmt.Println("Press 'P' to pause/resume the program.")

	err := setupGlobalHook()
	if err != nil {
		log.Fatalf("Error setting global hook: %v", err)
	}
	defer cleanupHook()

	messageLoop()
}

// Захват экрана
func captureScreen(hwnd uintptr) image.Image {
	rect := getClientRect(hwnd)
	// fmt.Printf("Capturing screen with rect: %+v\n", rect) // Отладочная информация

	// Проверка размера области захвата
	// fmt.Printf("Width: %d, Height: %d\n", rect.Dx(), rect.Dy())

	img, err := screenshot.Capture(int(rect.Min.X), int(rect.Min.Y), rect.Dx(), rect.Dy())
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return nil
	}

	return img
}

// Получение координат окна
func getClientRect(hwnd uintptr) (rect image.Rectangle) {
	var rectWin RECT
	ret, _, _ := getWindowRectProc.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rectWin)))
	if ret == 0 { // Проверяем результат на равенство 0
		fmt.Println("Failed to get window coordinates.")
		return image.Rectangle{}
	}
	return image.Rect(int(rectWin.Left), int(rectWin.Top), int(rectWin.Right), int(rectWin.Bottom))
}

// Клик!
func click(x, y int) {
	win.SetCursorPos(int32(x), int32(y))

	// Определение функции mouse_event через syscall
	const MOUSEEVENTF_LEFTDOWN = 0x0002
	const MOUSEEVENTF_LEFTUP = 0x0004
	mouseEvent := windows.NewLazySystemDLL("user32.dll").NewProc("mouse_event")
	mouseEvent.Call(uintptr(MOUSEEVENTF_LEFTDOWN), 0, 0, 0, 0)
	mouseEvent.Call(uintptr(MOUSEEVENTF_LEFTUP), 0, 0, 0, 0)
}

// Функция обработки игры
func gameClicker(hwnd uintptr) {
	startTime := time.Now()
	for time.Since(startTime) < 33*time.Second {

		// Проверяем установлена ли пауза
		if paused {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Переключаемся на окно перед каждым действием
		win.SetForegroundWindow(win.HWND(hwnd))
		time.Sleep(5 * time.Millisecond)

		screen := captureScreen(hwnd)
		if screen == nil {
			fmt.Println("Failed to capture screenshot.")
			continue
		}

		width, height := screen.Bounds().Dx(), screen.Bounds().Dy()

		// Логирование: размеры захваченного экрана
		// fmt.Printf("Captured screen size: width=%d, height=%d\n", width, height)

		for x := 0; x < width; x += 20 {
			for y := 0; y < height; y += 20 {
				r, g, b, _ := screen.At(x, y).RGBA()
				r, g, b = r>>8, g>>8, b>>8 // Преобразуем значения RGBA к 0-255

				// Логирование: текущий пиксель и его цвет
				// fmt.Printf("Checking pixel at (%d, %d): R=%d, G=%d, B=%d\n", x, y, r, g, b)

				// Проверяем все цветовые диапазоны из конфигурации
				for _, color := range config.Colors {
					// Генерируем диапазоны ±10 для каждого цвета
					redRange := [2]int{color.Red - 10, color.Red + 10}
					greenRange := [2]int{color.Green - 10, color.Green + 10}
					blueRange := [2]int{color.Blue - 10, color.Blue + 10}

					// Логирование: генерируемые диапазоны
					// fmt.Printf("Generated ranges for color (R=%d, G=%d, B=%d): Red=[%d, %d], Green=[%d, %d], Blue=[%d, %d]\n",
					// 	color.Red, color.Green, color.Blue,
					// 	redRange[0], redRange[1],
					// 	greenRange[0], greenRange[1],
					// 	blueRange[0], blueRange[1])

					// Проверяем, попадает ли текущий пиксель в диапазон
					if inRange(int(r), redRange) && inRange(int(g), greenRange) && inRange(int(b), blueRange) {
						// fmt.Printf("Pixel (%d, %d) matches color range: R=%d, G=%d, B=%d\n", x, y, r, g, b)

						rect := getClientRect(hwnd)
						if rect.Empty() {
							fmt.Println("Failed to get window coordinates.")
							time.Sleep(time.Second)
							continue
						}

						// Вычисляем координаты клика на экране
						screenX := int(rect.Min.X) + x
						screenY := int(rect.Min.Y) + y

						// Логирование: координаты клика
						// fmt.Printf("Clicking at screen coordinates: X=%d, Y=%d\n", screenX, screenY)

						// Клик по координатам
						click(screenX, screenY)
						time.Sleep(10 * time.Millisecond)

						// Прерываем цикл, так как объект найден
						break
					}
				}
			}
		}
	}
}

// Проверка пикселя и цветового диапазона
func inRange(value int, rangeValues [2]int) bool {
	return value >= rangeValues[0] && value <= rangeValues[1]
}
