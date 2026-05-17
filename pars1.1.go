package main

import (
        "fmt"
        "io/ioutil"
        "log"
        "regexp"
        "sort"
        "strconv"
"math"
)

// Структура для хранения данных трека перед сортировкой
type TrackResult struct {
        Number    string
        Asymmetry float64
        Mode      float64
        Variance  float64
        Result    float64 // Базовый индекс ритма
        Energy    float64 // Твой новый крупный индекс (Result / Variance)
}

func main() {
        // Укажите имя вашего файла здесь
        fileName := "./txtTable.html"

        // Читаем содержимое файла
        content, err := ioutil.ReadFile(fileName)
        if err != nil {
                log.Fatalf("Не удалось открыть файл %s: %v", fileName, err)
        }
        text := string(content)

        // Регулярные выражения для поиска данных
        trackRegex := regexp.MustCompile(`Трек\s+#(\d+)\s+\(Asymmetry:\s+([-\d.]+)\)`)
        modeVarianceRegex := regexp.MustCompile(`Total\s+Mode:\s+([-\d.]+),\s+Total\s+Variance:\s+([-\d.]+)`)

        // Ищем все совпадения
        trackMatches := trackRegex.FindAllStringSubmatch(text, -1)
        mvMatches := modeVarianceRegex.FindAllStringSubmatch(text, -1)

        if len(trackMatches) != len(mvMatches) {
                fmt.Printf("[Внимание] Количество треков (%d) и Total Mode/Variance (%d) не совпадает!\n\n",
                        len(trackMatches), len(mvMatches))
        }

        // Слайс для хранения результатов
        var results []TrackResult

        // Обрабатываем данные и сохраняем в слайс
        for i := 0; i < len(trackMatches) && i < len(mvMatches); i++ {
                trackNum := trackMatches[i][1]
                asymmetryStr := trackMatches[i][2]

                modeStr := mvMatches[i][1]
                varianceStr := mvMatches[i][2]

                asymmetry, err1 := strconv.ParseFloat(asymmetryStr, 64)
                mode, err2 := strconv.ParseFloat(modeStr, 64)
                variance, err3 := strconv.ParseFloat(varianceStr, 64)

                if err1 != nil || err2 != nil || err3 != nil {
                        fmt.Printf("Ошибка конвертации чисел на треке #%s\n", trackNum)
                        continue
                }

                // Защита от деления на ноль, если вдруг Variance равен 0.000
                if variance == 0 {
                        variance = 0.0001
                }

                // Вычисление базового индекса: (2 + Asymmetry) * Total Mode
                // Перед использованием math.Exp убедись, что пакет "math" импортирован вверху файла

// 1. Твой базовый ритмический индекс
finalResult := (2.0 + asymmetry) * mode

// 2. Настройка чувствительности экспоненты (K)
// Чем больше K, тем сильнее дисперсия будет занижать итоговый индекс
k := 50.0 

// 3. Вычисление экспоненциального коэффициента динамики
// math.Exp(-k * variance) вернет число от 0.0 до 1.0
dynamicFactor := math.Exp(-k * variance)

// 4. Финальный индекс энергии
energyResult := finalResult * dynamicFactor

                // Добавляем запись в слайс
                results = append(results, TrackResult{
                        Number:    trackNum,
                        Asymmetry: asymmetry,
                        Mode:      mode,
                        Variance:  variance,
                        Result:    finalResult,
                        Energy:    energyResult,
                })
        }

        // Сортировка по твоей новой крупной переменной (Energy) от большего к меньшему
        sort.Slice(results, func(i, j int) bool {
                return results[i].Energy > results[j].Energy
        })

        // Заголовок таблицы в консоли
        fmt.Printf("%-10s | %-10s | %-10s | %-10s | %-12s | %s\n", "№ Трека", "Asymmetry", "Total Mode", "Total Var", "Индекс Ритма", "Новый Индекс (~100)")
        fmt.Println("---------------------------------------------------------------------------------------------------------")

        // Вывод отсортированных по новой метрике данных
        for _, track := range results {
                fmt.Printf("Трек #%-3s | %-10.3f | %-10.3f | %-10.4f | %-12.4f | %.4f\n", 
                        track.Number, track.Asymmetry, track.Mode, track.Variance, track.Result, track.Energy)
        }
}

