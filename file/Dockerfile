FROM python:3.13-alpine AS builder

WORKDIR /app

# Install build dependencies and AWS CLI more efficiently
RUN apk add --no-cache \
    curl \
    gnupg \
    unzip \
    build-base \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf awscliv2.zip aws

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Final stage
FROM python:3.13-alpine

WORKDIR /app

# Install runtime dependencies only
RUN apk add --no-cache \
    netcat-openbsd \
    libstdc++

# Copy AWS CLI from builder
COPY --from=builder /usr/local/aws-cli /usr/local/aws-cli
COPY --from=builder /usr/local/bin/aws* /usr/local/bin/

# Copy installed Python packages
COPY --from=builder /usr/local/lib/python3.13/site-packages /usr/local/lib/python3.13/site-packages

# Copy application code
COPY . .

ENV PORT=8080

# Make the init script executable
RUN if [ -f /app/scripts/init-minio.sh ]; then chmod +x /app/scripts/init-minio.sh; fi

CMD exec gunicorn --bind :$PORT --workers 1 --threads 2 --timeout 0 --worker-class=gthread app:app 