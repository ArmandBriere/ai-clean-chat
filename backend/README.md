# Backend

## Requirements

Linux (as you should)

```bash
sudo apt install libopus-dev libopusfile-dev
```

MacOS

```bash
brew install opus opusfile
```

## Backend

### Golang

The project requires Go `1.23.3` and is using the `Sherpa-onnx` ai model.

Download the model files:

```bash
cd backend
./download-model.sh
```

### Python

Install [uv](https://docs.astral.sh/uv/getting-started/installation/):

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

Create the `.venv`:

```bash
uv venv
```

Install the dependencies:

```bash
uv sync
```

#### Formatting and Linting

We are using [ruff](https://docs.astral.sh/ruff/) to format and lint the code. The configuration of the linter and formatter is in the `pyproject.toml` file.

Lint the code:

```bash
ruff check --fix
```

Format the code:

```bash
ruff format
```

## WebRTC offers

- `minptime` - Minimal time in milliseconds that a packet is allowed to stay in the jitter buffer before being played out.
- `maxptime` - Maximum time in milliseconds that a packet is allowed to stay in the jitter buffer before being played out.
