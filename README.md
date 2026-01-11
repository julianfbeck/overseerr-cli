# overseerr-cli

A command-line interface for [Overseerr](https://overseerr.dev/), the media request management tool.

## Installation

```bash
# Homebrew (macOS/Linux)
brew install julianfbeck/tap/overseerr

# From source
go install github.com/julianfbeck/overseerr-cli@latest
```

## Configuration

Set your Overseerr server URL and API key:

```bash
# Using environment variables (recommended)
export OVERSEERR_URL="https://overseerr.example.com"
export OVERSEERR_API_KEY="your-api-key"

# Or save to config file
overseerr config set-url https://overseerr.example.com
overseerr config set-key your-api-key
```

Get your API key from Overseerr Settings → General → API Key.

## Usage

### Status

```bash
# Check server status
overseerr status
```

### Requests

```bash
# List requests
overseerr requests list
overseerr requests list --limit 10 --filter pending

# Get request details
overseerr requests get 123

# Request a movie (by TMDB ID)
overseerr requests movie 550

# Request a TV show (by TMDB ID)
overseerr requests tv 1396
overseerr requests tv 1396 --seasons 1,2,3

# Approve/decline requests
overseerr requests approve 123
overseerr requests decline 123

# Delete a request
overseerr requests delete 123 --force
```

### Discover

```bash
# Discover popular movies
overseerr discover movies
overseerr discover movies --page 2

# Discover popular TV shows
overseerr discover tv

# Trending content
overseerr discover trending --json
```

### Search

```bash
# Search for movies and TV shows
overseerr search "Breaking Bad" --json
```

### Users

```bash
# List users
overseerr users list

# Get current user
overseerr users me
```

### Media Details

```bash
# Get movie details (by TMDB ID)
overseerr media movie 550

# Get TV show details (by TMDB ID)
overseerr media tv 1396
```

## Options

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |
| `-q, --quiet` | Suppress non-essential output |
| `--no-color` | Disable color output |
| `-u, --url` | Override server URL |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `OVERSEERR_URL` | Server URL |
| `OVERSEERR_API_KEY` | API key |

## License

MIT
