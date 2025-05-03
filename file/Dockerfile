FROM python:3.13-slim

WORKDIR /app

# Install AWS CLI, netcat and other dependencies
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    unzip \
    netcat-traditional \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf awscliv2.zip aws \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

ENV PORT=8080

# Make the init script executable
RUN if [ -f /app/scripts/init-localstack.sh ]; then chmod +x /app/scripts/init-localstack.sh; fi

CMD exec gunicorn --bind :$PORT --workers 1 --threads 2 --timeout 0 --worker-class=gthread app:app 