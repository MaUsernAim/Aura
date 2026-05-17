package main

import (
        "fmt"
        "io/ioutil"
        "log"
        "regexp"
        "sort"
        "strconv"
)

// Структура для хранения данных трека перед сортировкой (поля дополнены)
type TrackResult struct {
        Number    string
        Asymmetry float64
        Mode      float64
        Variance  float64 // Добавлено поле для дисперсии
        Result    float64
        Energy    float64 // Добавлено поле для энергичности
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
        // Измененный паттерн: захватывает Mode во 2-ю группу, а Variance — в 3-ю группу (индексы 1 и 2 в матчах)
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

                // Достаем значения из объединенного регулярного выражения
                modeStr := mvMatches[i][1]
                varianceStr := mvMatches[i][2]

                asymmetry, err1 := strconv.ParseFloat(asymmetryStr, 64)
                mode, err2 := strconv.ParseFloat(modeStr, 64)
                variance, err3 := strconv.ParseFloat(varianceStr, 64)

                if err1 != nil || err2 != nil || err3 != nil {
                        fmt.Printf("Ошибка конвертации чисел на треке #%s\n", trackNum)
                        continue
                }

                // Вычисление по формуле: (2 + Asymmetry) * Total Mode
                finalResult := (2.0 + asymmetry) * mode

                // Вычисление энергичности: твой индекс * дисперсия * 1000
                energyResult := finalResult /variance 

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

        // Сортировка от большего к меньшему (по убыванию Result)
        sort.Slice(results, func(i, j int) bool {
                return results[i].Result > results[j].Result
        })

        // Заголовок таблицы в консоли (добавлены новые колонки)
        fmt.Printf("%-10s | %-10s | %-10s | %-10s | %-24s | %s\n", "№ Трека", "Asymmetry", "Total Mode", "Total Var", "Результат: (2+Asym)*Mode", "Энергичность")
        fmt.Println("---------------------------------------------------------------------------------------------------------")

        // Вывод отсортированных данных
        for _, track := range results {
                fmt.Printf("Трек #%-3s | %-10.3f | %-10.3f | %-10.4f | %-24.4f | %.4f\n", 
                        track.Number, track.Asymmetry, track.Mode, track.Variance, track.Result, track.Energy)
        }
}

