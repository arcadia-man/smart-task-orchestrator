FROM python:3.11-slim

# Install common test dependencies and tools for execution
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Install common python data libraries as per prompt
RUN pip install --no-cache-dir \
    pandas \
    numpy \
    requests

# Create working directory
WORKDIR /sandbox

# Define an entrypoint script that keeps the container alive
# so we can execute commands inside it when needed.
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'trap "exit" TERM INT' >> /entrypoint.sh && \
    echo 'while true; do sleep 1; done' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]
