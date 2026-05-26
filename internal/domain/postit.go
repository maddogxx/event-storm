package domain

import "image/color"

// Type identifies the semantic category of a post-it in an Event Storming session.
type Type string

const (
	TypeDomainEvent    Type = "domain_event"
	TypeCommand        Type = "command"
	TypeActor          Type = "actor"
	TypeAggregate      Type = "aggregate"
	TypeExternalSystem Type = "external_system"
	TypePolicy         Type = "policy"
	TypeReadModel      Type = "read_model"
	TypeHotspot        Type = "hotspot"
	TypeInformation    Type = "information"
)

// TypeInfo describes a post-it category for UI palettes and SVG exports.
type TypeInfo struct {
	Type        Type
	Label       string
	Description string
	Fill        color.NRGBA
	Stroke      color.NRGBA
	TextColor   color.NRGBA
}

// HexFill returns the fill color in CSS-compatible hex notation.
func (t TypeInfo) HexFill() string { return hex(t.Fill) }

// HexStroke returns the stroke color in CSS-compatible hex notation.
func (t TypeInfo) HexStroke() string { return hex(t.Stroke) }

// HexText returns the text color in CSS-compatible hex notation.
func (t TypeInfo) HexText() string { return hex(t.TextColor) }

func hex(c color.NRGBA) string {
	const digits = "0123456789ABCDEF"
	b := []byte{'#',
		digits[c.R>>4], digits[c.R&0x0F],
		digits[c.G>>4], digits[c.G&0x0F],
		digits[c.B>>4], digits[c.B&0x0F],
	}
	return string(b)
}

// typeCatalog is the canonical palette used across UI and exports.
var typeCatalog = []TypeInfo{
	{
		Type:        TypeDomainEvent,
		Label:       "Evento de Domínio",
		Description: "Algo relevante que aconteceu (passado).",
		Fill:        color.NRGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF},
		Stroke:      color.NRGBA{R: 0xCC, G: 0x84, B: 0x00, A: 0xFF},
		TextColor:   color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
	},
	{
		Type:        TypeCommand,
		Label:       "Comando",
		Description: "Intenção do usuário/sistema que dispara um evento.",
		Fill:        color.NRGBA{R: 0x5B, G: 0x9B, B: 0xD5, A: 0xFF},
		Stroke:      color.NRGBA{R: 0x2E, G: 0x75, B: 0xB6, A: 0xFF},
		TextColor:   color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	},
	{
		Type:        TypeActor,
		Label:       "Ator",
		Description: "Pessoa ou papel que executa um comando.",
		Fill:        color.NRGBA{R: 0xFF, G: 0xE6, B: 0x99, A: 0xFF},
		Stroke:      color.NRGBA{R: 0xBF, G: 0x90, B: 0x00, A: 0xFF},
		TextColor:   color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
	},
	{
		Type:        TypeAggregate,
		Label:       "Agregado",
		Description: "Conceito do domínio que mantém invariantes.",
		Fill:        color.NRGBA{R: 0xFF, G: 0xF2, B: 0xCC, A: 0xFF},
		Stroke:      color.NRGBA{R: 0xBF, G: 0x8F, B: 0x00, A: 0xFF},
		TextColor:   color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
	},
	{
		Type:        TypeExternalSystem,
		Label:       "Sistema Externo",
		Description: "Sistema de terceiro que participa do fluxo.",
		Fill:        color.NRGBA{R: 0xFF, G: 0x99, B: 0xCC, A: 0xFF},
		Stroke:      color.NRGBA{R: 0xC9, G: 0x57, B: 0x92, A: 0xFF},
		TextColor:   color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
	},
	{
		Type:        TypePolicy,
		Label:       "Política",
		Description: "Regra reativa: quando X acontece, faça Y.",
		Fill:        color.NRGBA{R: 0xB5, G: 0x7E, B: 0xDC, A: 0xFF},
		Stroke:      color.NRGBA{R: 0x6F, G: 0x42, B: 0xC1, A: 0xFF},
		TextColor:   color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	},
	{
		Type:        TypeReadModel,
		Label:       "Read Model",
		Description: "Visão/projeção usada para tomar decisões.",
		Fill:        color.NRGBA{R: 0x93, G: 0xC4, B: 0x7D, A: 0xFF},
		Stroke:      color.NRGBA{R: 0x55, G: 0x8C, B: 0x3F, A: 0xFF},
		TextColor:   color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
	},
	{
		Type:        TypeHotspot,
		Label:       "Hotspot",
		Description: "Ponto de atenção, dúvida ou risco.",
		Fill:        color.NRGBA{R: 0xE0, G: 0x2B, B: 0x20, A: 0xFF},
		Stroke:      color.NRGBA{R: 0x96, G: 0x1B, B: 0x14, A: 0xFF},
		TextColor:   color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	},
	{
		Type:        TypeInformation,
		Label:       "Informação",
		Description: "Nota informativa, contexto ou referência.",
		Fill:        color.NRGBA{R: 0xFA, G: 0xFA, B: 0xFA, A: 0xFF},
		Stroke:      color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF},
		TextColor:   color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
	},
}

// TypeCatalog returns a copy of the canonical post-it palette.
func TypeCatalog() []TypeInfo {
	out := make([]TypeInfo, len(typeCatalog))
	copy(out, typeCatalog)
	return out
}

// InfoFor returns the type metadata for the given category, falling back to Information.
func InfoFor(t Type) TypeInfo {
	for _, info := range typeCatalog {
		if info.Type == t {
			return info
		}
	}
	return typeCatalog[len(typeCatalog)-1]
}

// PostIt is a single sticky note placed on a dashboard canvas.
type PostIt struct {
	ID          int64
	DashboardID int64
	Type        Type
	Description string
	X           float32
	Y           float32
	Width       float32
	Height      float32
}

// DefaultSize returns the default post-it dimensions in canvas units.
func DefaultSize() (width, height float32) { return 160, 100 }
