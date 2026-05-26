# Event Storm

Aplicação web para criação de fluxos de **Event Storming** com post-its arrastáveis,
organização por projetos e dashboards, persistência em SQLite e exportação em SVG.

Construída com **Nuxt 3**, **TypeScript** e **Tailwind CSS**. Todo o estilo é
aplicado via classes utilitárias Tailwind diretamente nas tags — nenhum CSS
customizado é escrito (o único arquivo CSS é `assets/css/tailwind.css`, contendo
apenas as três diretivas `@tailwind`).

## Recursos

- **Projetos**: agrupam um conjunto de dashboards relacionados.
- **Dashboards**: o quadro onde o fluxo de Event Storm é desenhado.
- **Painel lateral de post-its** com os tipos clássicos do Event Storming
  (Evento de Domínio, Comando, Ator, Agregado, Política, Read Model, Sistema Externo, Hot Spot).
- **Drag & drop** do painel para o quadro e arraste livre dentro do quadro.
- **Edição** de título e descrição de cada post-it (duplo clique abre o editor).
- **Persistência** em arquivo SQLite (`./data/event-storm.db` por padrão).
- **Exportação em SVG** do dashboard com um clique.

## Requisitos

- Node.js 22+ (o repositório inclui um `.nvmrc`; rode `nvm use` para alinhar)
- npm 10+ (ou pnpm/yarn equivalentes)
- Compilador C++ disponível como fallback para `better-sqlite3` (na maioria dos casos
  o pacote baixa um prebuild compatível)

## Rodando localmente

```bash
npm install
npm run dev
```

A aplicação ficará disponível em `http://localhost:3000`.

O banco SQLite é criado automaticamente em `./data/event-storm.db`. Você pode
sobrescrever esse caminho com a variável de ambiente `NUXT_DB_PATH`:

```bash
NUXT_DB_PATH=/caminho/para/event-storm.db npm run dev
```

## Build de produção

```bash
npm run build
npm run start
```

## Docker

Para executar via Docker:

```bash
docker build -t event-storm .
docker run --rm -p 3000:3000 -v "$(pwd)/data:/app/data" event-storm
```

O volume `-v "$(pwd)/data:/app/data"` garante que o banco SQLite sobreviva
entre execuções do container.

O `Dockerfile` usa um entrypoint que ajusta automaticamente o owner do diretório
montado para o usuário não-root `nuxt` (UID 1001) antes de iniciar o servidor,
então o bind mount funciona mesmo quando o diretório no host pertence a outro
UID. Os arquivos `event-storm.db*` no host ficarão pertencendo ao UID 1001 — se
precisar removê-los do host, rode:

```bash
docker run --rm -v "$(pwd)/data:/data" alpine sh -c 'rm -rf /data/*'
```

Como alternativa, use um named volume e deixe o Docker gerenciar permissões:

```bash
docker run --rm -p 3000:3000 -v event-storm-data:/app/data event-storm
```

## Estrutura de pastas

```
.
├── app.vue                    # Layout raiz
├── assets/css/tailwind.css    # Diretivas @tailwind (único arquivo CSS)
├── components/                # PostItPanel, PostItNote, PostItEditor, DashboardCanvas
├── composables/               # usePostItTypes, useDashboardExport
├── pages/                     # index, projects/[id], dashboards/[id]
├── server/api/                # Endpoints REST (projects, dashboards, postits)
├── server/utils/db.ts         # Conexão e migrations SQLite
├── types/                     # Interfaces TypeScript do domínio
├── tailwind.config.cjs        # Tema e content paths do Tailwind
├── Dockerfile
└── nuxt.config.ts
```

## API REST

| Método | Rota                                   | Descrição                              |
| ------ | -------------------------------------- | -------------------------------------- |
| GET    | `/api/projects`                        | Lista projetos                         |
| POST   | `/api/projects`                        | Cria projeto                           |
| GET    | `/api/projects/:id`                    | Detalhe do projeto                     |
| PUT    | `/api/projects/:id`                    | Atualiza projeto                       |
| DELETE | `/api/projects/:id`                    | Remove projeto (cascade em dashboards) |
| GET    | `/api/projects/:id/dashboards`         | Lista dashboards do projeto            |
| POST   | `/api/dashboards`                      | Cria dashboard                         |
| GET    | `/api/dashboards/:id`                  | Detalhe do dashboard                   |
| PUT    | `/api/dashboards/:id`                  | Atualiza dashboard                     |
| DELETE | `/api/dashboards/:id`                  | Remove dashboard                       |
| GET    | `/api/dashboards/:id/postits`          | Lista post-its do dashboard            |
| POST   | `/api/postits`                         | Cria post-it                           |
| PUT    | `/api/postits/:id`                     | Atualiza post-it (texto e/ou posição)  |
| DELETE | `/api/postits/:id`                     | Remove post-it                         |

## Como usar

1. Crie um **projeto** na tela inicial.
2. Dentro do projeto, crie um **dashboard**.
3. No editor do dashboard, **arraste** um post-it do painel da esquerda para o quadro.
4. **Clique duas vezes** em qualquer post-it para abrir o editor lateral e
   adicionar título e descrição.
5. **Arraste** os post-its livremente para organizar o fluxo.
6. Use o botão **Exportar SVG** para baixar o dashboard como imagem vetorial.
