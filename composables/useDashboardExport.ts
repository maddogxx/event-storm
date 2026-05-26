import type { Dashboard, PostIt } from '~/types'

const escapeXml = (input: string): string =>
  input.replace(/[<>&"']/g, (char) => {
    switch (char) {
      case '<': return '&lt;'
      case '>': return '&gt;'
      case '&': return '&amp;'
      case '"': return '&quot;'
      case "'": return '&apos;'
      default: return char
    }
  })

const wrapText = (text: string, maxChars: number): string[] => {
  if (!text) return []
  const words = text.split(/\s+/)
  const lines: string[] = []
  let current = ''
  for (const word of words) {
    const candidate = current ? `${current} ${word}` : word
    if (candidate.length > maxChars && current) {
      lines.push(current)
      current = word
    } else {
      current = candidate
    }
  }
  if (current) lines.push(current)
  return lines
}

interface BuildSvgOptions {
  dashboard: Pick<Dashboard, 'name'>
  postits: PostIt[]
}

const buildSvg = ({ dashboard, postits }: BuildSvgOptions): string => {
  const { byKind } = usePostItTypes()
  const padding = 40
  const headerHeight = 80

  const maxX = postits.reduce((max, p) => Math.max(max, p.x + p.width), 1200)
  const maxY = postits.reduce((max, p) => Math.max(max, p.y + p.height), 600)
  const width = maxX + padding * 2
  const height = maxY + padding * 2 + headerHeight

  const items = postits
    .map((p) => {
      const type = byKind(p.kind)
      const x = p.x + padding
      const y = p.y + padding + headerHeight
      const titleLines = wrapText(p.title || type.label, Math.max(10, Math.floor(p.width / 9)))
      const descLines = p.description
        ? wrapText(p.description, Math.max(12, Math.floor(p.width / 7)))
        : []

      const titleSvg = titleLines
        .map((line, i) =>
          `<text x="${x + 12}" y="${y + 26 + i * 18}" font-family="Inter, sans-serif" font-size="14" font-weight="700" fill="${type.textColor}">${escapeXml(line)}</text>`
        )
        .join('')

      const descSvg = descLines
        .slice(0, Math.max(0, Math.floor((p.height - 50 - titleLines.length * 18) / 14)))
        .map((line, i) =>
          `<text x="${x + 12}" y="${y + 26 + titleLines.length * 18 + 14 + i * 14}" font-family="Inter, sans-serif" font-size="11" fill="${type.textColor}" opacity="0.85">${escapeXml(line)}</text>`
        )
        .join('')

      return `
        <g>
          <rect x="${x}" y="${y}" width="${p.width}" height="${p.height}" rx="6" ry="6"
                fill="${type.color}" stroke="#1f2937" stroke-opacity="0.15" stroke-width="1" />
          <rect x="${x}" y="${y}" width="${p.width}" height="18" rx="6" ry="6"
                fill="#000" fill-opacity="0.08" />
          <text x="${x + 8}" y="${y + 13}" font-family="Inter, sans-serif" font-size="10"
                font-weight="600" fill="${type.textColor}" opacity="0.85">
            ${escapeXml(type.label.toUpperCase())}
          </text>
          ${titleSvg}
          ${descSvg}
        </g>
      `
    })
    .join('')

  return `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
  <rect width="100%" height="100%" fill="#F9FAFB" />
  <text x="${padding}" y="${padding + 30}" font-family="Inter, sans-serif" font-size="24" font-weight="700" fill="#111827">
    ${escapeXml(dashboard.name)}
  </text>
  <text x="${padding}" y="${padding + 54}" font-family="Inter, sans-serif" font-size="12" fill="#6B7280">
    Event Storm · exportado em ${new Date().toLocaleString('pt-BR')}
  </text>
  ${items}
</svg>`
}

export const useDashboardExport = () => {
  const exportToSvg = (dashboard: Dashboard, postits: PostIt[]): void => {
    const svg = buildSvg({ dashboard, postits })
    const blob = new Blob([svg], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const safeName = dashboard.name.replace(/[^a-z0-9-_]+/gi, '_').toLowerCase()
    link.href = url
    link.download = `${safeName || 'dashboard'}.svg`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  return { exportToSvg }
}
