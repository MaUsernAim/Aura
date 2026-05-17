package main

import (
	"math"
	"sort"
"time"
)

// TrackData - структура, описывающая все уровни данных одной песни (необходима для метода FindFivePoints)
type TrackData struct {
	Title            string        // Имя песни (из тегов или имя файла)
	Artist           string        // Автор / Исполнитель
	Duration         time.Duration // Длительность трека
	RawTimeline      []float64     // Полный, не сокращенный формат (энергия каждого фрейма ~46мс)
	FilteredTimeline []float64     // Сокращенный вариант (после локального Gated-фильтра)
}

// FiveVectors представляет собой 5 сглаженных макро-характеристик для 5 участков песни
type FiveVectors struct {
	Intro   float64 // Участок 1: Вступление
	Verse   float64 // Участок 2: Куплет
	BuildUp float64 // Участок 3: Разгон / Ожидание
	Chorus  float64 // Участок 4: Припев / Катарсис / Объяснение
	Outro   float64 // Участок 5: Финал / Затухание
}

// FindFivePoints анализирует сокращенный массив данных, находит 5 ключевых участков,
// фильтрует аномалии и возвращает сглаженные макро-векторы.
func (td *TrackData) FindFivePoints() FiveVectors {
	timeline := td.FilteredTimeline
	n := len(timeline)

	// Защита от пустых или слишком коротких треков
	if n < 10 {
		return FiveVectors{}
	}

	// =========================================================================
	// ШАГ 1: ПОИСК ТОЧЕК ПЕРЕХОДА (Через скорость изменения энергии)
	// =========================================================================
	type Point struct {
		Index float64
		Value float64
	}
	var gradients []Point

	// Вычисляем первую производную (разницу между соседними фреймами)
	for i := 1; i < n; i++ {
		grad := math.Abs(timeline[i] - timeline[i-1])
		gradients = append(gradients, Point{Index: float64(i), Value: grad})
	}

	// Сортируем изменения по убыванию силы всплеска
	sort.Slice(gradients, func(i, j int) bool {
		return gradients[i].Value > gradients[j].Value
	})

	// Нам нужно найти 4 точки раздела, чтобы получить 5 участков.
	// Чтобы точки не сгрудились в одном припеве, вводим минимальную дистанцию между ними (микро-окно)
	var boundaryIndices []int
	minDistance := n / 10 // Участки не должны быть короче 10% от всей песни

	for _, g := range gradients {
		idx := int(g.Index)
		tooClose := false
		for _, b := range boundaryIndices {
			if math.Abs(float64(idx-b)) < float64(minDistance) {
				tooClose = true
				break
			}
		}
		if !tooClose {
			boundaryIndices = append(boundaryIndices, idx)
		}
		if len(boundaryIndices) == 4 {
			break
		}
	}

	// Если песня слишком монотонна и алгоритм не нашел 4 ярких пика,
	// делим трек на 5 равных пропорциональных зон (резервный математический план)
	if len(boundaryIndices) < 4 {
		boundaryIndices = []int{n / 5, (n * 2) / 5, (n * 3) / 5, (n * 4) / 5}
	}

	// Сортируем границы по времени хронологически
	sort.Ints(boundaryIndices)

	// Массив индексов границ: [0, точка1, точка2, точка3, точка4, конец_песни]
	boundaries := []int{0, boundaryIndices[0], boundaryIndices[1], boundaryIndices[2], boundaryIndices[3], n}

	// =========================================================================
	// ШАГ 2: КУСОЧНАЯ СГЛАЖИВАЮЩАЯ ФИЛЬТРАЦИЯ (С ОТБРАСЫВАНИЕМ АНОМАЛИЙ)
	// =========================================================================
	sectionMeans := make([]float64, 5)

	for s := 0; s < 5; s++ {
		start := boundaries[s]
		end := boundaries[s+1]

		// 2.1 Находим предварительное среднее арифметическое участка (базовый ориентир)
		sum := 0.0
		count := 0
		for i := start; i < end; i++ {
			if timeline[i] > 0 { // Считаем только живые, не зануленные ранее фреймы
				sum += timeline[i]
				count++
			}
		}

		if count == 0 {
			sectionMeans[s] = 0.0
			continue
		}
		globalSectionMean := sum / float64(count)

		// 2.2 ВТОРОЙ ПРОХОД: Фильтр аномалий.
		// Если точка отличается от среднего значения участка более чем на 50% (высокое отличие),
		// мы её ОТБРАСЫВАЕМ (игнорируем). Это убирает случайные шумы и технические всплески.
		cleanSum := 0.0
		cleanCount := 0
		maxDeviation := 1.5 // Порог допустимого отличия (от 0.5 до 1.5 от среднего)
		minDeviation := 0.5

		for i := start; i < end; i++ {
			val := timeline[i]
			if val == 0 {
				continue
			}

			ratio := val / globalSectionMean
			if ratio >= minDeviation && ratio <= maxDeviation {
				cleanSum += val
				cleanCount++
			}
		}

		// Если после жесткой фильтрации что-то осталось — берем чистое среднее,
		// иначе откатываемся к первичному среднему значению участка
		if cleanCount > 0 {
			sectionMeans[s] = cleanSum / float64(cleanCount)
		} else {
			sectionMeans[s] = globalSectionMean
		}
	}
	// Упаковываем очищенную математику в 5 векторов
	return FiveVectors{
		Intro:   sectionMeans[0],
		Verse:   sectionMeans[1],
		BuildUp: sectionMeans[2],
		Chorus:  sectionMeans[3],
		Outro:   sectionMeans[4],
	}
}
