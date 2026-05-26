// Package export renders an Event Storming dashboard into a static SVG file.
package export

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/madd0gxx/event-storm/internal/domain"
)

const (
	margin     = 24.0
	titleArea  = 56.0
	defaultBg  = "#F5F5F0"
	fontFamily = "Helvetica, Arial, sans-serif"
)

// Dashboard writes the dashboard and its post-its as an SVG file at path.
func Dashboard(path string, board *domain.Dashboard, postits []*domain.PostIt) error {
	if board == nil {
		return fmt.Errorf("export: dashboard is required")
	}

	svg := render(board, postits)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create svg: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(svg); err != nil {
		return fmt.Errorf("write svg: %w", err)
	}
	return nil
}

func render(board *domain.Dashboard, postits []*domain.PostIt) string {
	width, height := bounds(postits)

	var b strings.Builder
	fmt.Fprintf(&b, `<?xml version="1.0" encoding="UTF-8"?>`+"\n")
	fmt.Fprintf(&b,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n",
		width, height, width, height,
	)

	fmt.Fprintf(&b, `<rect x="0" y="0" width="%.0f" height="%.0f" fill="%s"/>`+"\n",
		width, height, defaultBg)

	fmt.Fprintf(&b,
		`<text x="%.0f" y="%.0f" font-family="%s" font-size="22" font-weight="700" fill="#222">%s</text>`+"\n",
		margin, margin+18, fontFamily, escape(board.Name),
	)
	fmt.Fprintf(&b,
		`<text x="%.0f" y="%.0f" font-family="%s" font-size="11" fill="#666">Exportado em %s</text>`+"\n",
		margin, margin+36, fontFamily, time.Now().Format("02/01/2006 15:04"),
	)

	for _, p := range postits {
		info := domain.InfoFor(p.Type)
		x := p.X - minX(postits) + margin
		y := p.Y - minY(postits) + margin + titleArea
		drawPostIt(&b, x, y, p.Width, p.Height, info, p.Description)
	}

	b.WriteString(`</svg>` + "\n")
	return b.String()
}

func drawPostIt(b *strings.Builder, x, y, w, h float32, info domain.TypeInfo, description string) {
	fmt.Fprintf(b,
		`<g transform="translate(%.2f %.2f)">`+"\n",
		x, y,
	)
	fmt.Fprintf(b,
		`  <rect width="%.2f" height="%.2f" rx="4" ry="4" fill="%s" stroke="%s" stroke-width="1.5"/>`+"\n",
		w, h, info.HexFill(), info.HexStroke(),
	)
	fmt.Fprintf(b,
		`  <text x="8" y="16" font-family="%s" font-size="10" font-weight="700" fill="%s">%s</text>`+"\n",
		fontFamily, info.HexText(), escape(strings.ToUpper(info.Label)),
	)
	writeWrappedText(b, description, w-16, h-26, info.HexText())
	b.WriteString(`</g>` + "\n")
}

func writeWrappedText(b *strings.Builder, text string, maxWidth, maxHeight float32, fill string) {
	const (
		fontSize   = 12.0
		lineHeight = 14.0
		approxChar = 6.5 // heuristic for average glyph advance at fontSize
	)
	if text == "" {
		return
	}
	maxChars := int(maxWidth / approxChar)
	if maxChars < 8 {
		maxChars = 8
	}
	maxLines := int(maxHeight / lineHeight)
	if maxLines < 1 {
		maxLines = 1
	}

	lines := wrap(text, maxChars, maxLines)
	startY := 32.0
	fmt.Fprintf(b, `  <text x="8" y="%.0f" font-family="%s" font-size="%.0f" fill="%s">`+"\n",
		startY, fontFamily, fontSize, fill)
	for i, line := range lines {
		dy := 0.0
		if i > 0 {
			dy = lineHeight
		}
		fmt.Fprintf(b, `    <tspan x="8" dy="%.0f">%s</tspan>`+"\n", dy, escape(line))
	}
	b.WriteString(`  </text>` + "\n")
}

func wrap(text string, maxChars, maxLines int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var (
		lines []string
		cur   strings.Builder
	)
	flush := func() {
		if cur.Len() > 0 {
			lines = append(lines, cur.String())
			cur.Reset()
		}
	}
	for _, word := range words {
		if cur.Len() == 0 {
			cur.WriteString(word)
			continue
		}
		if cur.Len()+1+len(word) > maxChars {
			flush()
			cur.WriteString(word)
			continue
		}
		cur.WriteByte(' ')
		cur.WriteString(word)
	}
	flush()
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		if len(lines[maxLines-1]) > 3 {
			lines[maxLines-1] = strings.TrimRight(lines[maxLines-1][:len(lines[maxLines-1])-1], " ") + "…"
		}
	}
	return lines
}

func bounds(postits []*domain.PostIt) (width, height float32) {
	if len(postits) == 0 {
		return 800, 480
	}
	var maxX, maxY float32 = 0, 0
	mnX, mnY := minX(postits), minY(postits)
	for _, p := range postits {
		x := p.X - mnX + p.Width
		y := p.Y - mnY + p.Height
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}
	width = maxX + margin*2
	height = maxY + margin*2 + titleArea
	if width < 480 {
		width = 480
	}
	if height < 320 {
		height = 320
	}
	return width, height
}

func minX(postits []*domain.PostIt) float32 {
	if len(postits) == 0 {
		return 0
	}
	m := postits[0].X
	for _, p := range postits {
		if p.X < m {
			m = p.X
		}
	}
	return m
}

func minY(postits []*domain.PostIt) float32 {
	if len(postits) == 0 {
		return 0
	}
	m := postits[0].Y
	for _, p := range postits {
		if p.Y < m {
			m = p.Y
		}
	}
	return m
}

func escape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}
