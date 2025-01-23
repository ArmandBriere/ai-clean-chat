# AI powered clean chat

This is an application developed for ConFoo 2025. It is a video chat application that uses AI to transcript live conversations and identify inappropriate content.

The application is built with Svelte and Python.

## Requirements

- Docker
- Docker Compose
- mkcert

First, you need to install the local CA into the system trust store. This need to be ran only once on your system.

```bash
cd traefik/certs/
mkcert ai-clean-chat.home.site
mkcert -install
```

The commands above will generate a `ai-clean-chat.home.site.pem` and `ai-clean-chat.home.site-key.pem` files that are tied to your machine and user.

Then update your `/etc/hosts` file to add the following entry:

```txt
127.0.0.1 ai-clean-chat.home.site
```

You will now be able to access the application at [https://ai-clean-chat.home.site](https://ai-clean-chat.home.site) with a proper ssl certificate.

## Installation

1. Initialize the project

```bash
make init
```

2. Install dependencies

```bash
make install
```

3. Run the vite dev server

```bash
make dev
```

4. Visit the application at [https://ai-clean-chat.home.site](https://ai-clean-chat.home.site)

## Backend

The backend is a Python application that uses the AI model to analyze the conversation. It uses the `faster-whisper` model to analyze the conversation and a custom model to identify inappropriate content.

The backend received a WebRTC stream from the frontend and have the following roles:

- Create a WebRTC connection between the two users
- Transcribe the conversation of the two users
  - Stream back the full transcript to each user
  - Invoke the AI model to analyze the conversation and censor inappropriate content
  - Stream back the censored words to each owner with an explanation of why it was censored

Both of those model are defined, trained (if needed) and tested in the `ai` folder.

### Setup

To setup the backend, you need to install [uv](https://github.com/astral-sh/uv) which handles python dependencies and virtual environment.

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

Then you can install the dependencies with:

```bash
uv install
```

### Run

To run any command, prefix it with `uv run` as usual, for example:

```bash
# Run a script
uv run main.py

# Run a command
uv run isort .
```

## Frontend

The frontend is a Svelte application that uses WebRTC to handle the video chat and the backend to analyze the conversation.

The frontend dependencies are managed with [bun](https://bun.sh/), you can install it with:

```bash
curl -fsSL https://bun.sh/install | bash
```

Then you can install the dependencies with:

```bash
bun install
```

### Run

To run any command, prefix it with `bun run` as usual, for example:

```bash
# Run a script
bun run dev

# Run a command
bun run prettier . --check
```

## AI

The AI folder contains the AI models used by the backend.
