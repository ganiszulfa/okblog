FROM node:20-alpine as build

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install
RUN npm ci

COPY . .

ARG VITE_ADMIN_API_BASE_URL
ARG VITE_UPLOADED_FILE_HOST

ENV VITE_ADMIN_API_BASE_URL=$VITE_ADMIN_API_BASE_URL
ENV VITE_UPLOADED_FILE_HOST=$VITE_UPLOADED_FILE_HOST

RUN echo "Building with VITE_ADMIN_API_BASE_URL=${VITE_ADMIN_API_BASE_URL} and VITE_UPLOADED_FILE_HOST=${VITE_UPLOADED_FILE_HOST}"

RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"] 