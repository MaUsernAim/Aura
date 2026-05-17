package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
type FullTrackDopamine struct {
	TotalMode      float64
	TotalVariance  float64
	TotalAsymmetry float64

	IntroDopamine   SectionDopamine
	VerseDopamine   SectionDopamine
	BuildUpDopamine SectionDopamine
	ChorusDopamine  SectionDopamine
	OutroDopamine  SectionDopamine
}

// Метод для сортировки и форматирования треков
func ProcessAndFormatTracks(tracks []FullTrackDopamine) (string, string) {
	// Сортируем треки по TotalAsymmetry (по возрастанию)
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].TotalAsymmetry < tracks[j].TotalAsymmetry
	})

	// Текстовый формат
	var textBuilder strings.Builder
	textBuilder.WriteString("ТРЕКИ, ОТСОРТИРОВАННЫЕ ПО TOTAL ASYMMETRY\n")
	textBuilder.WriteString(strings.Repeat("=", 80) + "\n")

	for i, track := range tracks {
		textBuilder.WriteString(fmt.Sprintf("Трек #%d (Asymmetry: %.3f)\n", i+1, track.TotalAsymmetry))
		textBuilder.WriteString(fmt.Sprintf("  Total Mode: %.3f, Total Variance: %.3f\n", track.TotalMode, track.TotalVariance))
		textBuilder.WriteString("  Сегменты:\n")
		textBuilder.WriteString(fmt.Sprintf("    Intro: Mode=%.3f, Variance=%.3f\n", track.IntroDopamine.Mode, track.IntroDopamine.Variance))
		textBuilder.WriteString(fmt.Sprintf("    Verse: Mode=%.3f, Variance=%.3f\n", track.VerseDopamine.Mode, track.VerseDopamine.Variance))
		textBuilder.WriteString(fmt.Sprintf("    BuildUp: Mode=%.3f, Variance=%.3f\n", track.BuildUpDopamine.Mode, track.BuildUpDopamine.Variance))
		textBuilder.WriteString(fmt.Sprintf("    Chorus: Mode=%.3f, Variance=%.3f\n", track.ChorusDopamine.Mode, track.ChorusDopamine.Variance))
		textBuilder.WriteString(fmt.Sprintf("    Outro: Mode=%.3f, Variance=%.3f\n", track.OutroDopamine.Mode, track.OutroDopamine.Variance))
		textBuilder.WriteString("\n")
	}

	// HTML формат
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<h1>Треки, отсортированные по Total Asymmetry</h1>\n")
	htmlBuilder.WriteString("<table border='1'>\n")
	htmlBuilder.WriteString("<thead><tr>\n")
	htmlBuilder.WriteString("<th>№</th><th>Total Asymmetry</th><th>Total Mode</th><th>Total Variance</th>\n")
	htmlBuilder.WriteString("<th>Intro</th><th>Verse</th><th>BuildUp</th><th>Chorus</th><th>Outro</th>\n")
	htmlBuilder.WriteString("</tr></thead>\n<tbody>\n")

	for i, track := range tracks {
		htmlBuilder.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%.3f</td><td>%.3f</td><td>%.3f</td>\n",
			i+1, track.TotalAsymmetry, track.TotalMode, track.TotalVariance))

		// Форматируем сегменты как "Mode/Variance"
		formatSection := func(s SectionDopamine) string {
			return fmt.Sprintf("%.2f/%.2f", s.Mode, s.Variance)
	}

		htmlBuilder.WriteString(fmt.Sprintf(
    "<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
    formatSection(track.IntroDopamine),
    formatSection(track.VerseDopamine),
    formatSection(track.BuildUpDopamine),
    formatSection(track.ChorusDopamine),
    formatSection(track.OutroDopamine),
))


	htmlBuilder.WriteString("</tbody></table>")
}
	return textBuilder.String(), htmlBuilder.String()


}
*//*
func ProcessAndFormatTracks(tracks []FullTrackDopamine) (string, string) {
    // Сортируем треки по TotalAsymmetry (по возрастанию)
    sort.Slice(tracks, func(i, j int) bool {
        return tracks[i].TotalAsymmetry < tracks[j].TotalAsymmetry
    })

    // Текстовый формат
    var textBuilder strings.Builder
    textBuilder.WriteString("ТРЕКИ, ОТСОРТИРОВАННЫЕ ПО TOTAL ASYMMETRY\n")
    textBuilder.WriteString(strings.Repeat("=", 80) + "\n")

    for i, track := range tracks {
        textBuilder.WriteString(fmt.Sprintf("Трек #%d (Asymmetry: %.3f)\n", i+1, track.TotalAsymmetry))
        textBuilder.WriteString(fmt.Sprintf("  Total Mode: %.3f, Total Variance: %.3f\n", track.TotalMode, track.TotalVariance))
        textBuilder.WriteString("  Сегменты:\n")
        textBuilder.WriteString(fmt.Sprintf("    Intro: Mode=%.3f, Variance=%.3f\n", track.IntroDopamine.Mode, track.IntroDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    Verse: Mode=%.3f, Variance=%.3f\n", track.VerseDopamine.Mode, track.VerseDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    BuildUp: Mode=%.3f, Variance=%.3f\n", track.BuildUpDopamine.Mode, track.BuildUpDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    Chorus: Mode=%.3f, Variance=%.3f\n", track.ChorusDopamine.Mode, track.ChorusDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    Outro: Mode=%.3f, Variance=%.3f\n", track.OutroDopamine.Mode, track.OutroDopamine.Variance))
        textBuilder.WriteString("\n")
    }

    // HTML формат
    var htmlBuilder strings.Builder
    htmlBuilder.WriteString("<h1>Треки, отсортированные по Total Asymmetry</h1>\n")
    htmlBuilder.WriteString("<table border='1'>\n")
    htmlBuilder.WriteString("<thead><tr>\n")
    htmlBuilder.WriteString("<th>№</th><th>Total Asymmetry</th><th>Total Mode</th><th>Total Variance</th>\n")
    htmlBuilder.WriteString("<th>Intro</th><th>Verse</th><th>BuildUp</th><th>Chorus</th><th>Outro</th>\n")
    htmlBuilder.WriteString("</tr></thead>\n<tbody>\n")

    for i, track := range tracks {
        htmlBuilder.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%.3f</td><td>%.3f</td><td>%.3f</td>\n",
            i+1, track.TotalAsymmetry, track.TotalMode, track.TotalVariance))


        // Форматируем сегменты как "Mode/Variance"
        formatSection := func(s SectionDopamine) string {
            return fmt.Sprintf("%.2f/%.2f", s.Mode, s.Variance)
        }

        htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
    formatSection(track.IntroDopamine),
    formatSection(track.VerseDopamine),
    formatSection(track.BuildUpDopamine),
    formatSection(track.ChorusDopamine),
    formatSection(track.OutroDopamine))
)

    htmlBuilder.WriteString("</tbody></table>")

    return textBuilder.String(), htmlBuilder.String()
}

*/
func ProcessAndFormatTracks(tracks []FullTrackDopamine) (string, string) {
    // Сортируем треки по TotalAsymmetry (по возрастанию)
    sort.Slice(tracks, func(i, j int) bool {
        return tracks[i].TotalAsymmetry < tracks[j].TotalAsymmetry
    })

    // Текстовый формат
    var textBuilder strings.Builder
    textBuilder.WriteString("ТРЕКИ, ОТСОРТИРОВАННЫЕ ПО TOTAL ASYMMETRY\n")
    textBuilder.WriteString(strings.Repeat("=", 80) + "\n")

    for i, track := range tracks {
        textBuilder.WriteString(fmt.Sprintf("Трек #%d (Asymmetry: %.3f)\n", i+1, track.TotalAsymmetry))
        textBuilder.WriteString(fmt.Sprintf("  Total Mode: %.3f, Total Variance: %.3f\n", track.TotalMode, track.TotalVariance))
        textBuilder.WriteString("  Сегменты:\n")
        textBuilder.WriteString(fmt.Sprintf("    Intro: Mode=%.3f, Variance=%.3f\n", track.IntroDopamine.Mode, track.IntroDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    Verse: Mode=%.3f, Variance=%.3f\n", track.VerseDopamine.Mode, track.VerseDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    BuildUp: Mode=%.3f, Variance=%.3f\n", track.BuildUpDopamine.Mode, track.BuildUpDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    Chorus: Mode=%.3f, Variance=%.3f\n", track.ChorusDopamine.Mode, track.ChorusDopamine.Variance))
        textBuilder.WriteString(fmt.Sprintf("    Outro: Mode=%.3f, Variance=%.3f\n", track.OutroDopamine.Mode, track.OutroDopamine.Variance))
        textBuilder.WriteString("\n")
    }

    // HTML формат
    var htmlBuilder strings.Builder
    htmlBuilder.WriteString("<h1>Треки, отсортированные по Total Asymmetry</h1>\n")
    htmlBuilder.WriteString("<table border='1'>\n")
    htmlBuilder.WriteString("<thead><tr>\n")
    htmlBuilder.WriteString("<th>№</th><th>Total Asymmetry</th><th>Total Mode</th><th>Total Variance</th>\n")
    htmlBuilder.WriteString("<th>Intro</th><th>Verse</th><th>BuildUp</th><th>Chorus</th><th>Outro</th>\n")
    htmlBuilder.WriteString("</tr></thead>\n<tbody>\n")

    for i, track := range tracks {
        htmlBuilder.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%.3f</td><td>%.3f</td><td>%.3f</td>",
            i+1, track.TotalAsymmetry, track.TotalMode, track.TotalVariance))

        // Форматируем сегменты как "Mode/Variance"
        formatSection := func(s SectionDopamine) string {
            return fmt.Sprintf("%.2f/%.2f", s.Mode, s.Variance)
        }

        htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
            formatSection(track.IntroDopamine),
            formatSection(track.VerseDopamine),
            formatSection(track.BuildUpDopamine),
            formatSection(track.ChorusDopamine),
            formatSection(track.OutroDopamine)))
    }

    htmlBuilder.WriteString("</tbody></table>")

    return textBuilder.String(), htmlBuilder.String()
}

