# Multi-stage build for optimal production image
FROM node:20-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Copy package files
COPY package*.json ./
# Use npm install if no package-lock.json exists, otherwise use npm ci
RUN if [ -f package-lock.json ]; then npm ci --only=production; else npm install --only=production; fi && npm cache clean --force

# Build the application
FROM base AS builder
WORKDIR /app
COPY package*.json ./
# Use npm install if no package-lock.json exists, otherwise use npm ci
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

COPY . .
RUN npm run build

# Production image
FROM base AS runner
WORKDIR /app

# Create a non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Copy built application
COPY --from=builder --chown=nuxtjs:nodejs /app/.output /app/.output

USER nuxtjs

# Expose port
EXPOSE 3000

# Set environment variables
ENV NODE_ENV=production
ENV NUXT_HOST=0.0.0.0
ENV NUXT_PORT=3000
# Enable logging to stdout/stderr for Docker
ENV NITRO_LOG_LEVEL=info
ENV DEBUG=nuxt:*

# Start the application
CMD ["node", ".output/server/index.mjs"] 