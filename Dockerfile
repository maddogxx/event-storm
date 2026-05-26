# syntax=docker/dockerfile:1.7
FROM node:22-bookworm-slim AS base
ENV NODE_ENV=production \
    NUXT_TELEMETRY_DISABLED=1
WORKDIR /app

FROM base AS deps
ENV NODE_ENV=development
RUN apt-get update \
    && apt-get install -y --no-install-recommends python3 make g++ \
    && rm -rf /var/lib/apt/lists/*
COPY package*.json ./
RUN npm install

FROM deps AS build
COPY . .
RUN npm run build

FROM node:22-bookworm-slim AS runner
ENV NODE_ENV=production \
    NUXT_TELEMETRY_DISABLED=1 \
    HOST=0.0.0.0 \
    PORT=3000 \
    NUXT_DB_PATH=/app/data/event-storm.db \
    NUXT_DATA_DIR=/app/data
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends tini gosu \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 1001 nuxt \
    && useradd --system --create-home --uid 1001 --gid 1001 nuxt \
    && mkdir -p /app/data \
    && chown -R nuxt:nuxt /app

COPY --from=build /app/.output ./.output
COPY --from=build /app/node_modules/better-sqlite3 ./node_modules/better-sqlite3
COPY --from=build /app/node_modules/bindings ./node_modules/bindings
COPY --from=build /app/node_modules/file-uri-to-path ./node_modules/file-uri-to-path
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh \
    && chown -R nuxt:nuxt /app

VOLUME ["/app/data"]
EXPOSE 3000
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/docker-entrypoint.sh"]
CMD ["node", ".output/server/index.mjs"]
