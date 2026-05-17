package main

import (
	"fmt"
	"os"
//	"sort"
	"time"
)

func main() {
	folderPath := "/storage/emulated/0/Download/" // Заполните путь к папке с MP3 файлами

	if folderPath == "" {
		fmt.Println("Ошибка: переменная folderPath не инициализирована")
		fmt.Println("Пожалуйста, укажите путь к папке с музыкой в переменной folderPath")
		return
	}

	fmt.Printf("Начинаю анализ папки: %s\n\n", folderPath)
	tracks, err := AnalyzeMusicFolder(folderPath)
	if err != nil {
		fmt.Printf("Ошибка при анализе папки: %v\n", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("В указанной папке не найдено MP3 файлов")
		return
	}

	fmt.Printf("Найдено треков: %d\n", len(tracks))
	fmt.Println("Создаю HTML-таблицу...")

	// Создаём HTML файл
	htmlFile, err := os.Create("music_analysis.html")
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer htmlFile.Close()

	// Заголовок HTML
	htmlFile.WriteString(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Анализ музыкальных треков</title>
    <style>
        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            margin: 20px;
            background: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .info {
            text-align: center;
            color: #666;
            margin-bottom: 20px;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            background: white;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: center;
            vertical-align: top;
        }
        th {
            background-color: #4CAF50;
            color: white;
            font-weight: bold;
            position: sticky;
            top: 0;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .good {
            background-color: #c8e6c9;
            font-weight: bold;
        }
        .bad {
            background-color: #ffcdd2;
        }
        .section-table {
            font-size: 12px;
            margin-top: 5px;
        }
        .section-table td {
            padding: 4px;
        }
        .track-title {
            font-weight: bold;
            color: #2196F3;
        }
        .note {
            margin-top: 20px;
            padding: 10px;
            background: #fff3e0;
            border-left: 4px solid #ff9800;
        }
    </style>
</head>
<body>
    <h1>🎵 Анализ музыкальных треков</h1>
    <div class="info">
        Дата: ` + time.Now().Format("02.01.2006 15:04:05") + `<br>
        Обработано треков: ` + fmt.Sprintf("%d", len(tracks)) + `
    </div>
`)

	// Таблица с треками
	htmlFile.WriteString(`    <table id="track-table">
        <thead>
            <tr>
                <th>#</th>
                <th>Трек</th>
                <th>Длит.</th>
                <th>Mode</th>
                <th>Variance</th>
                <th>Asymmetry</th>
                <th>КЭР</th>
                <th>Пики, %</th>
                <th>5 участков</th>
                <th>Мезо-профиль</th>
            </tr>
        </thead>
        <tbody>
`)

var arr []FullTrackDopamine


	for idx, track := range tracks {
		// Получаем профили
		fiveVectors := track.FindFivePoints()
		dopamineProfile := track.AnalyzeDopamineMeters()
		arr=append(arr,dopamineProfile)
		// Считаем КЭР
		var variances []float64
		sections := []SectionDopamine{
			dopamineProfile.IntroDopamine,
			dopamineProfile.VerseDopamine,
			dopamineProfile.BuildUpDopamine,
			dopamineProfile.ChorusDopamine,
			dopamineProfile.OutroDopamine,
		}
		for _, s := range sections {
			if s.Variance > 0 {
				variances = append(variances, s.Variance)
			}
		}
		var KER float64
		if len(variances) > 0 {
			maxVar, minVar := variances[0], variances[0]
			sumVar := 0.0
			for _, v := range variances {
				if v > maxVar {
					maxVar = v
				}
				if v < minVar {
					minVar = v
				}
				sumVar += v
			}
			meanVar := sumVar / float64(len(variances))
			if meanVar > 0 {
				KER = (maxVar - minVar) / meanVar
			}
		}

		// Считаем процент сохранённых пиков
		survived := 0
		for _, v := range track.FilteredTimeline {
			if v > 0 {
				survived++
			}
		}
		peakPercent := float64(survived) / float64(len(track.FilteredTimeline)) * 100

		// Определяем тип трека по асимметрии
		asymType := ""
		if dopamineProfile.TotalAsymmetry > 0.5 {
			asymType = "🔥 Энергичный"
		} else if dopamineProfile.TotalAsymmetry < -0.5 {
			asymType = "🌙 Спокойный"
		} else {
			asymType = "⚖️ Сбалансированный"
		}

		// Строка таблицы
		htmlFile.WriteString(fmt.Sprintf(`             <tr>
                <td>%d</td>
                <td class="track-title">%s<br><small>%s</small></td>
                <td>%s</td>
                <td>%.4f</td>
                <td>%.4f</td>
                <td>%.4f<br><small>%s</small></td>
                <td>%.2f</td>
                <td>%.1f%%</td>
                <td style="font-size:12px">
                    I:%.3f<br>V:%.3f<br>B:%.3f<br>C:%.3f<br>O:%.3f
                </td>
                <td style="font-size:11px">
                    <table class="section-table">
                        <tr><td>Intro</td><td>M:%.4f</td><td>D:%.4f</td><td>P:%.3f</td></tr>
                        <tr><td>Verse</td><td>M:%.4f</td><td>D:%.4f</td><td>P:%.3f</td></tr>
                        <tr><td>BuildUp</td><td>M:%.4f</td><td>D:%.4f</td><td>P:%.3f</td></tr>
                        <tr><td>Chorus</td><td>M:%.4f</td><td>D:%.4f</td><td>P:%.3f</td></tr>
                        <tr><td>Outro</td><td>M:%.4f</td><td>D:%.4f</td><td>P:%.3f</td></tr>
                    </table>
                </td>
             </tr>
`,
			idx+1,
			track.Title,
			track.Artist,
			track.Duration.String(),
			dopamineProfile.TotalMode,
			dopamineProfile.TotalVariance,
			dopamineProfile.TotalAsymmetry,
			asymType,
			KER,
			peakPercent,
			fiveVectors.Intro,
			fiveVectors.Verse,
			fiveVectors.BuildUp,
			fiveVectors.Chorus,
			fiveVectors.Outro,
			sections[0].Mode, sections[0].Variance, sections[0].Predictability,
			sections[1].Mode, sections[1].Variance, sections[1].Predictability,
			sections[2].Mode, sections[2].Variance, sections[2].Predictability,
			sections[3].Mode, sections[3].Variance, sections[3].Predictability,
			sections[4].Mode, sections[4].Variance, sections[4].Predictability,
		))
	}

	htmlFile.WriteString(`        </tbody>
    </table>
    <div class="note">
        <strong>📖 Легенда:</strong><br>
        • <strong>Mode</strong> — базовая энергия трека (мода)<br>
        • <strong>Variance</strong> — дисперсия (контрастность, неожиданность)<br>
        • <strong>Asymmetry</strong> — асимметрия (>0.5 = громкие участки доминируют, <-0.5 = тихие участки доминируют)<br>
        • <strong>КЭР</strong> — коэффициент эмоционального рельефа (насколько секции отличаются друг от друга)<br>
        • <strong>Пики, %</strong> — процент фреймов, выживших после фильтрации<br>
        • <strong>I/V/B/C/O</strong> — Intro, Verse, BuildUp, Chorus, Outro (энергия по секциям)<br>
        • <strong>M/D/P</strong> — Mode/Dispersion/Predictability для каждой секции
    </div>
    <script>
        // Сортировка таблицы при клике на заголовок
        document.querySelectorAll('th').forEach(header => {
            header.addEventListener('click', () => {
                const table = document.getElementById('track-table');
                const tbody = table.querySelector('tbody');
                const rows = Array.from(tbody.querySelectorAll('tr'));
                const index = header.cellIndex;
                const type = index === 2 ? 'string' : 'number';
                
                rows.sort((a, b) => {
                    let aVal = a.cells[index].innerText;
                    let bVal = b.cells[index].innerText;
                    if (type === 'number') {
                        aVal = parseFloat(aVal) || 0;
                        bVal = parseFloat(bVal) || 0;
                        return aVal - bVal;
                    }
                    return aVal.localeCompare(bVal);
                });
                
                rows.forEach(row => tbody.appendChild(row));
            });
        });
    </script>
</body>
</html>
`)

	fmt.Printf("\n✅ Готово! Файл сохранён: music_analysis.html\n")
	fmt.Println("Откройте его в браузере для просмотра таблицы.")
txt,html :=ProcessAndFormatTracks(arr)
htmlYaFile, err := os.Create("htmlTable.html")
        if err != nil {
                fmt.Printf("Ошибка создания файла: %v\n", err)
                return
        }
        defer htmlYaFile.Close()
htmlYaFile.WriteString(html)
textFile, err := os.Create("txtTable.html")
        if err != nil {
                fmt.Printf("Ошибка создания файла: %v\n", err)
                return
        }
        defer textFile.Close()
textFile.WriteString(txt)
}
