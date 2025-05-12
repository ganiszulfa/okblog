FROM nginx:stable-alpine

# Copy nginx configuration files
COPY ./api-gateway.conf /etc/nginx/nginx.conf
COPY ./cors-preflight.conf /etc/nginx/cors-preflight.conf

# Create cache directory
RUN mkdir -p /data/nginx/cache

# Expose port 80
EXPOSE 80

# Run nginx in foreground
CMD ["nginx", "-g", "daemon off;"] 