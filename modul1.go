package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/dhowden/tag"
//	"gonum.org/v1/gonum/stat"
)
/*
// TrackData — структура, описывающая все уровни данных одной песни
type TrackData struct {
	Title            string        // Имя песни (из тегов или имя файла)
	Artist           string        // Автор / Исполнитель
	Duration         time.Duration // Длительность трека
	RawTimeline      []float64     // Полный, не сокращенный формат (энергия каждого фрейма ~46мс)
	FilteredTimeline []float64     // Сокращенный вариант (после локального Gated-фильтра)
}
*/
// AnalyzeMusicFolder — главный метод, который сканирует папку и собирает массивы данных
func AnalyzeMusicFolder(folderPath string) ([]TrackData, error) {
	var result []TrackData

	// 1. Читаем все файлы в папке
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать папку: %v", err)
	}

	for _, file := range files {
		// Проверяем, что это MP3 файл
		if file.IsDir() || !strings.HasSuffix(strings.ToLower(file.Name()), ".mp3") {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		fmt.Printf("Анализирую: %s...\n", file.Name())

		// 2. Читаем метаданные (Имя и Автор) через dhowden/tag
		title, artist := extractMetadata(filePath, file.Name())

		// 3. Декодируем MP3 через hajimehoshi/go-mp3
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Ошибка открытия файла %s: %v\n", file.Name(), err)
			continue
		}
		defer f.Close()

		decoder, err := mp3.NewDecoder(f)
		if err != nil {
			fmt.Printf("Ошибка создания декодера MP3 для %s: %v\n", file.Name(), err)
			continue
		}

		// Читаем все PCM данные
		var pcmData []byte
		buf := make([]byte, 4096)
		for {
			n, err := decoder.Read(buf)
			if n > 0 {
				pcmData = append(pcmData, buf[:n]...)
			}
			if err != nil {
				break
			}
		}

		// Вычисляем длительность
		// Читаем все PCM данные
//var pcmData []byte
buf = make([]byte, 4096)
for {
    n, err := decoder.Read(buf)
    if n > 0 {
        pcmData = append(pcmData, buf[:n]...)
    }
    if err != nil {
        break
    }
}

// Получаем параметры аудио и преобразуем в int
// Получаем частоту дискретизации
sampleRate := int(decoder.SampleRate())

// MP3 обычно имеет 2 канала (стерео)
channels := 2

// Вычисляем длительность
duration := time.Duration(float64(len(pcmData)) / float64(sampleRate*channels*2) * float64(time.Second))

// Преобразуем в моно-сэмплы
samples := convertToMonoFloat64(pcmData, channels)

// Вычисляем длительность
duration = time.Duration(float64(len(pcmData)) / float64(sampleRate*channels*2) * float64(time.Second))

// 4. Превращаем сырые байты в моно-сэмплы float64
samples = convertToMonoFloat64(pcmData, channels)
		// 5. СОЗДАЕМ СЫРОЙ ПОЛНЫЙ ФОРМАТ (RawTimeline)
		// Бьем на окна по 2048 сэмплов (~46 мс при 44100 Гц) и считаем RMS энергию
		windowSize := 2048
		var rawTimeline []float64

		for i := 0; i < len(samples); i += windowSize {
			end := i + windowSize
			if end > len(samples) {
				end = len(samples)
			}

			sum := 0.0
			for j := i; j < end; j++ {
				sum += samples[j] * samples[j]
			}
			rms := math.Sqrt(sum / float64(end-i))
			rawTimeline = append(rawTimeline, rms)
		}

		// 6. СОЗДАЕМ СОКРАЩЕННЫЙ ВАРИАНТ (FilteredTimeline)
		// Применяем локальный фильтр: усредняем по окну в 11 фреймов
		filteredTimeline := applyGatedFilter(rawTimeline, 11, 1.2)

		// Упаковываем все данные в объект трека
		track := TrackData{
			Title:            title,
			Artist:           artist,
			Duration:         duration,
			RawTimeline:      rawTimeline,
			FilteredTimeline: filteredTimeline,
		}

		result = append(result, track)
	}

	return result, nil
}

// Вспомогательная функция: перевод int16-байтов в моно float64
func convertToMonoFloat64(data []byte, channels int) []float64 {
	// 2 байта на один int16 сэмпл
	totalSamples := len(data) / 2
	var monoSamples []float64

	for i := 0; i < totalSamples; i += channels {
		// Читаем левый канал (низкий и высокий байт)
		if i*2+1 >= len(data) {
			break
		}
		leftBits := int16(data[i*2]) | int16(data[i*2+1])<<8
		left := float64(leftBits) / 32768.0

		right := left
		if channels > 1 && (i+1)*2+1 < len(data) {
			// Читаем правый канал
			rightBits := int16(data[(i+1)*2]) | int16(data[(i+1)*2+1])<<8
			right = float64(rightBits) / 32768.0
		}

		// Микшируем в моно (L+R)/2
		monoSamples = append(monoSamples, (left+right)/2.0)
	}
	return monoSamples
}

// Вспомогательная функция: локальный Gated-фильтр для сокращения данных
func applyGatedFilter(raw []float64, windowSize int, thresholdFactor float64) []float64 {
	filtered := make([]float64, len(raw))
	half := windowSize / 2
	for i := 0; i < len(raw); i++ {
		start := i - half
		if start < 0 {
			start = 0
		}
		end := i + half
		if end >= len(raw) {
			end = len(raw) - 1
		}

		sum := 0.0
		count := 0
		for j := start; j <= end; j++ {
			sum += raw[j]
			count++
		}
		localMean := sum / float64(count)

		// Если это всплеск (припев/пик) — сохраняем локальное среднее, иначе — обнуляем
		if raw[i] > localMean*thresholdFactor {
			filtered[i] = localMean
		} else {
			filtered[i] = 0.0 // Сжимаем информацию, убирая куплетный монотонный шум
		}
	}
	return filtered
}

// Вспомогательная функция для безопасного извлечения тегов ID3
func extractMetadata(filePath string, fileName string) (string, string) {
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	artist := "Неизвестен"

	f, err := os.Open(filePath)
	if err == nil {
		defer f.Close()
		metadata, err := tag.ReadFrom(f)
		if err == nil {
			if metadata.Title() != "" {
				title = metadata.Title()
			}
			if metadata.Artist() != "" {
				artist = metadata.Artist()
			}
		}
	}
	return title, artist
}
