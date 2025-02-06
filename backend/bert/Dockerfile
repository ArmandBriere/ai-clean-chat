FROM python:3.12-slim AS builder

# UV setup
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

WORKDIR /app

# Copy UV requirements
COPY pyproject.toml uv.lock /app/

# Install dependencies
RUN uv sync

COPY ./ /app

# Run Flask app
ENTRYPOINT [ "uv" ]
CMD [ "run", "python3", "-m" , "flask", "run", "--host=0.0.0.0", "--port=8080" ]
