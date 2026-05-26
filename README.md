# Event Storm

Aplicação desktop em Go, construída com [Fyne](https://fyne.io), para facilitar
sessões de **Event Storming**: organização de eventos, comandos, agregados,
políticas, atores e demais elementos em um quadro visual.

## Funcionalidades

- **Projetos** agrupam múltiplos dashboards (um por bounded context, fluxo, etc.).
- **Dashboards** são quadros editáveis onde você desenha o fluxo Event Storming.
- **Painel lateral** com a paleta clássica de post-its (Evento, Comando, Ator,
  Agregado, Sistema Externo, Política, Read Model, Hotspot, Informação) — clique
  para adicionar ao quadro.
- **Edição visual**: arraste post-its livremente, clique para editar a descrição
  ou trocar o tipo.
- **Persistência local** em SQLite (`event-storm.db`).
- **Exportação SVG** do dashboard com cores, textos e layout preservados.

## Pré-requisitos

- **Go 1.25+**
- **Toolchain CGO** (o Fyne usa GLFW/OpenGL via CGO).
  - Linux (Ubuntu/Debian/WSL):

    ```bash
    sudo apt update
    sudo apt install -y gcc pkg-config libgl1-mesa-dev xorg-dev \
        libxkbcommon-dev libxcursor-dev libxrandr-dev libxinerama-dev \
        libxi-dev libwayland-dev
    ```

  - macOS: instale o Xcode Command Line Tools (`xcode-select --install`).
  - Windows: instale MinGW-w64 ou MSYS2.

## Como executar

```bash
git clone <repo>
cd event-storm
go mod tidy
go run .
```

O banco de dados é criado automaticamente em
`$XDG_CONFIG_HOME/event-storm/event-storm.db` (Linux),
`~/Library/Application Support/event-storm/event-storm.db` (macOS) ou
`%AppData%\event-storm\event-storm.db` (Windows).

Você pode sobrescrever o caminho via variável de ambiente:

```bash
EVENT_STORM_DB=/tmp/storm.db go run .
```

## Build de produção

```bash
go build -o event-storm .
./event-storm
```

Ou empacotamento nativo via Fyne CLI (opcional):

```bash
go install fyne.io/tools/cmd/fyne@latest
fyne package -os linux  # ou darwin / windows
```

## Estrutura do projeto

```
event-storm/
├── main.go                 # Bootstrap (configura DB e dispara a UI)
├── internal/
│   ├── domain/             # Entidades Project, Dashboard, PostIt + paleta de tipos
│   ├── storage/            # Repositório SQLite (modernc.org/sqlite, sem CGO)
│   ├── export/             # Renderização SVG do dashboard
│   └── ui/                 # Telas Fyne (projetos, dashboards, editor) e widget de post-it
└── README.md
```

A camada `domain` é pura — sem dependências de Fyne ou SQLite — e descreve a
paleta de cores Event Storming canônica.

## Atalhos no editor

| Ação                       | Como fazer                                 |
| -------------------------- | ------------------------------------------ |
| Adicionar post-it          | Clique no item desejado da paleta lateral. |
| Mover post-it              | Clique e arraste o post-it.                |
| Editar descrição / tipo    | Clique uma vez no post-it.                 |
| Remover post-it            | Clique no post-it → "Remover post-it".     |
| Navegar pelo quadro        | Use as barras de rolagem.                  |
| Exportar dashboard em SVG  | Botão "Exportar SVG" no topo do editor.    |

## Testes

```bash
go test ./internal/domain/... ./internal/storage/... ./internal/export/...
```

Os pacotes não-UI rodam sem CGO. Testes de UI exigem o ambiente gráfico do Fyne.

## Paleta Event Storming utilizada

| Tipo            | Cor         | Significado                                            |
| --------------- | ----------- | ------------------------------------------------------ |
| Evento          | Laranja     | Fato relevante que aconteceu (passado).                |
| Comando         | Azul        | Intenção que dispara um evento.                        |
| Ator            | Amarelo     | Pessoa ou papel que emite comandos.                    |
| Agregado        | Amarelo-pal | Entidade do domínio que mantém invariantes.            |
| Sistema Externo | Rosa        | Sistema de terceiros que participa do fluxo.           |
| Política        | Roxo        | Regra reativa ("quando X acontecer, faça Y").          |
| Read Model      | Verde       | Visão/projeção usada para tomar decisões.              |
| Hotspot         | Vermelho    | Risco, dúvida ou ponto de atenção.                     |
| Informação      | Branco      | Nota, contexto ou referência.                          |
