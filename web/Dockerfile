# Build stage
FROM node:20-alpine as build-stage

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy all project files
COPY . .

# Build the project
RUN npm run build

# Production stage
FROM node:20-alpine as production-stage

WORKDIR /app

# Copy package files from build stage
COPY --from=build-stage /app/package*.json ./

# Install production dependencies only
RUN npm install --only=production

# Copy build output from build stage
COPY --from=build-stage /app/.nuxt ./.nuxt
COPY --from=build-stage /app/static ./static
COPY --from=build-stage /app/nuxt.config.js ./

# Expose the port Nuxt.js runs on
EXPOSE 3000

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=3000
ENV NODE_ENV=production

# Start the server
CMD ["npm", "start"] 