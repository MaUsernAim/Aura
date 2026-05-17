package main

import (
//	"fmt"
	"math"
	"sort"
//"time"
	"gonum.org/v1/gonum/stat"
)
/*
// TrackData - структура, описывающая все уровни данных одной песни (необходима для метода AnalyzeDopamineMeters)
type TrackData struct {
	Title            string        // Имя песни (из тегов или имя файла)
	Artist           string        // Автор / Исполнитель
	Duration         time.Duration // Длительность трека
	RawTimeline      []float64     // Полный, не сокращенный формат (энергия каждого фрейма ~46мс)
	FilteredTimeline []float64     // Сокращенный вариант (после локального Gated-фильтра)
}
*/
// SectionDopamine описывает математические триггеры внутри одного конкретного участка
type SectionDopamine struct {
	Mode           float64 // Мода (базовый предсказуемый уровень энергии)
	Variance       float64 // Дисперсия (мера случайности и интересности звука)
	Predictability float64 // Индекс угадываемости (коэффициент вариации)
}

// FullTrackDopamine Профиль всей песни на макро- и мезо-уровнях
type FullTrackDopamine struct {
	// Макро-уровень (Вся песня)
	TotalMode      float64
	TotalVariance  float64
	TotalAsymmetry float64 // Асимметрия распределения (куда сдвинут звук: тихо/громко)

	// Мезо-уровень (5 отрезков по отдельности)
	IntroDopamine   SectionDopamine
	VerseDopamine   SectionDopamine
	BuildUpDopamine SectionDopamine
	ChorusDopamine  SectionDopamine
	OutroDopamine   SectionDopamine
}

// AnalyzeDopamineMeters считает ТВиМС-параметры на макро-уровне и внутри каждого из 5 отрезков
func (td *TrackData) AnalyzeDopamineMeters() FullTrackDopamine {
	timeline := td.FilteredTimeline
	n := len(timeline)

	if n < 10 {
		return FullTrackDopamine{}
	}

	var profile FullTrackDopamine

	// =========================================================================
	// 1. МАКРО-УРОВЕНЬ: Расчет по всей песне целиком
	// =========================================================================
	// Очищаем массив от нулей для честной статистики ТВиМС
	var activeSamples []float64
	for _, v := range timeline {
		if v > 0 {
			activeSamples = append(activeSamples, v)
		}
	}
	// Если весь трек занулен фильтром, берем сырые данные
	if len(activeSamples) < 5 {
		activeSamples = td.RawTimeline
	}

	mode, _ := stat.Mode(activeSamples, nil)
profile.TotalMode = mode
	profile.TotalVariance = stat.Variance(activeSamples, nil)

	// Считаем асимметрию (куда сдвинут звук: в тихо или в громко)
	mean := stat.Mean(activeSamples, nil)
	stdDev := stat.StdDev(activeSamples, nil)
	if stdDev > 0 {
		// Формула асимметрии Пирсона: (Среднее - Мода) / СКО
		profile.TotalAsymmetry = (mean - profile.TotalMode) / stdDev
	}

	// =========================================================================
	// 2. МЕЗО-УРОВЕНЬ: Находим границы 5 участков (повторяем логику градиентов)
	// =========================================================================
	type Point struct {
		Index int
		Value float64
	}
	var gradients []Point
	for i := 1; i < n; i++ {
		gradients = append(gradients, Point{Index: i, Value: math.Abs(timeline[i] - timeline[i-1])})
	}
	sort.Slice(gradients, func(i, j int) bool { return gradients[i].Value > gradients[j].Value })

	var boundaryIndices []int
 	minDistance := n / 10
	for _, g := range gradients {
		tooClose := false
		for _, b := range boundaryIndices {
			if math.Abs(float64(g.Index-b)) < float64(minDistance) {
				tooClose = true
				break
			}
		}
		if !tooClose {
			boundaryIndices = append(boundaryIndices, g.Index)
		}
		if len(boundaryIndices) == 4 {
			break
		}
	}
	if len(boundaryIndices) < 4 {
		boundaryIndices = []int{n / 5, (n * 2) / 5, (n * 3) / 5, (n * 4) / 5}
	}
	sort.Ints(boundaryIndices)
	boundaries := []int{0, boundaryIndices[0], boundaryIndices[1], boundaryIndices[2], boundaryIndices[3], n}

	// Вспомогательная функция для расчета ТВиМС внутри конкретного отрезка индекса
	calcSection := func(start, end int) SectionDopamine {
		var sectionSamples []float64
		for i := start; i < end; i++ {
			if timeline[i] > 0 {
				sectionSamples = append(sectionSamples, timeline[i])
			}
		}
		// Защита от пустых отрезков
		if len(sectionSamples) < 2 {
			return SectionDopamine{}
		}

		var sd SectionDopamine
	sMode, _ := stat.Mode(sectionSamples, nil)
sd.Mode = sMode

		sd.Variance = stat.Variance(sectionSamples, nil)
		// Считаем угадываемость через Коэффициент вариации (СКО / Среднее) [1]
		sMean := stat.Mean(sectionSamples, nil)
		sStd := stat.StdDev(sectionSamples, nil)
		if sMean > 0 {
			// Чем ближе к 0 — тем монотоннее и предсказуемее ритмика
			// Чем выше — тем больше неожиданных контрастов и "интересности" [1]
			sd.Predictability = sStd / sMean
		}
		return sd
	}

	// Обсчитываем каждый из 5 участков
	profile.IntroDopamine = calcSection(boundaries[0], boundaries[1])
	profile.VerseDopamine = calcSection(boundaries[1], boundaries[2])
	profile.BuildUpDopamine = calcSection(boundaries[2], boundaries[3])
	profile.ChorusDopamine = calcSection(boundaries[3], boundaries[4])
	profile.OutroDopamine = calcSection(boundaries[4], boundaries[5])

	return profile
}
