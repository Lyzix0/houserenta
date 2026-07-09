FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat
# Installed once here (not via corepack) so every derived stage below already
# has the binary — corepack would otherwise re-download it per stage, since the
# deps stage's download isn't part of the base image the other stages branch from.
RUN npm install -g pnpm@10.33.2
WORKDIR /app

FROM base AS deps
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .
# Baked into the client bundle at build time: must be an address the user's
# browser can reach (not the internal docker-compose service name).
ARG NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}
# next.config.ts's rewrites() is evaluated once during `next build` and baked
# into .next/routes-manifest.json — `next start` just replays it, so setting
# this only at container runtime (via docker-compose `environment:`) has no
# effect. It must also be present here, at build time.
ARG API_INTERNAL_BASE_URL
ENV API_INTERNAL_BASE_URL=${API_INTERNAL_BASE_URL}
RUN pnpm build

FROM base AS runner
ENV NODE_ENV=production
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/next.config.ts ./next.config.ts
EXPOSE 3000
CMD ["pnpm", "start"]
