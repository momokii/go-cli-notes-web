# Knowledge Garden CLI - Web Dashboard

A standalone web dashboard showcasing the Knowledge Garden CLI application. Built with Go Fiber and Tailwind CSS.

## What is This?

This is a marketing/informational website for the [Knowledge Garden CLI](https://github.com/momokii/go-cli-notes) - a terminal-based Personal Knowledge Management system for developers.

**Design**: The dashboard uses a playful, Notion-inspired design with emoji and warm colors. Alternative design variants (`index1.html`, `index2.html`) are available in the templates folder but not actively served.

## Features

- **Go Fiber v2** - Fast, lightweight web framework
- **Tailwind CSS** - Utility-first CSS via CDN (no build step)
- **Live Demos** - 4 GIF demos showcasing CLI features
- **Rate Limiting** - 120 requests/minute per IP
- **Security Headers** - Proper HTTP security headers
- **Custom 404 Page** - Friendly error page
- **Health Check** - `/health` endpoint for monitoring
- **Responsive Design** - Works on all device sizes

## Quick Start

### Prerequisites

- Go 1.24 or later
- Docker (optional, for containerized deployment)

### Run Locally

```bash
# From the web/ directory
cd web

# Download dependencies
go mod download

# Run the server
go run main.go
```

The dashboard will be available at `http://localhost:3000`

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `ENV` | Environment (development/production) | `development` |

### Example .env File

```bash
PORT=3000
ENV=development
```

## Building for Production

### Build Binary

```bash
# Build the binary
go build -o dashboard main.go

# Run the binary
./dashboard
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o dashboard-linux main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o dashboard-macos main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dashboard-macos-arm main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o dashboard.exe main.go
```

## Deployment

### Docker Deployment (Recommended)

#### Using Docker Compose

```bash
# Build and start the container
docker compose up -d

# View logs
docker compose logs -f

# Stop the container
docker compose down
```

#### Using Docker Directly

```bash
# Build the image
docker build -t kg-dashboard .

# Run the container
docker run -d -p 3000:3000 --name kg-dashboard kg-dashboard

# Run on different port
docker run -d -p 8080:3000 --name kg-dashboard kg-dashboard

# View logs
docker logs -f kg-dashboard

# Stop the container
docker stop kg-dashboard
docker rm kg-dashboard
```

#### Docker Configuration

The Dockerfile uses a multi-stage build for efficiency:
- **Builder stage**: Compiles Go binary with build flags for smaller size
- **Runtime stage**: Minimal Alpine Linux image with ca-certificates
- **Security**: Runs as non-root user (dashboard:1000)
- **Health check**: Built-in health check on `/health` endpoint

### Using systemd (Linux)

Create `/etc/systemd/system/kg-dashboard.service`:

```ini
[Unit]
Description=Knowledge Garden CLI Dashboard
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/web
Environment="PORT=3000"
Environment="ENV=production"
ExecStart=/path/to/web/dashboard
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Then:

```bash
sudo systemctl daemon-reload
sudo systemctl enable kg-dashboard
sudo systemctl start kg-dashboard
```

## Project Structure

```
web/
├── main.go                 # Fiber server entry point
├── go.mod                 # Go module definition
├── go.sum                 # Dependencies checksum
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Docker Compose configuration
├── .dockerignore          # Docker build exclusions
├── templates/             # HTML templates
│   ├── index1.html       # Minimalist design (archived)
│   ├── index2.html       # Technical design (archived)
│   ├── index3.html       # Playful design (active)
│   └── 404.html          # Custom 404 page
├── static/                # Static assets
│   ├── css/              # Shared scoped styles
│   └── demo/             # Demo GIFs
│       ├── cli-notes-1.gif
│       ├── cli-notes-2.gif
│       ├── cli-notes-3.gif
│       └── cli-notes-4.gif
└── README.md             # This file
```

## Middleware

The server includes the following middleware:

1. **Logger** - Request logging with timestamps and latency
2. **Recovery** - Panic recovery
3. **Compress** - Gzip compression
4. **CORS** - Cross-origin resource sharing
5. **Security Headers** - X-Frame-Options, X-XSS-Protection, etc.
6. **Rate Limiting** - 120 req/min per IP using Fiber's built-in limiter

## Routes

| Route | Description |
|-------|-------------|
| `GET /` | Main page (index3.html) |
| `GET /health` | Health check endpoint |
| `GET /static/*` | Static files (CSS, JS, images) |

## Customization

### Changing Content

Edit the content directly in `templates/index3.html`. The template is self-contained with:
- Header with navigation
- Hero section with CTAs
- Features section (3 cards)
- Demo section (4 GIFs)
- How It Works section
- Installation section with code block
- Footer with links

### Changing Colors

The template uses Tailwind CSS utility classes. To customize colors, search for gradient classes:
- Purple gradients: `from-purple-* to-pink-*`
- Orange gradients: `from-orange-* to-pink-*`

### Using Alternative Templates

To switch to a different template (e.g., index1.html or index2.html):

1. Edit `main.go` line ~125
2. Change `"index3.html"` to your preferred template
3. Restart the server

### Adding New Pages

1. Create new HTML file in `templates/`
2. Add route in `main.go`:
   ```go
   app.Get("/your-page", func(c *fiber.Ctx) error {
       c.Set("Content-Type", "text/html; charset=utf-8")
       content, err := os.ReadFile(filepath.Join(templatesPath, "your-page.html"))
       if err != nil {
           return c.Status(fiber.StatusNotFound).SendString("Template not found")
       }
       return c.SendString(string(content))
   })
   ```

## Troubleshooting

### Port Already in Use

```bash
# Find what's using port 3000
lsof -i :3000

# Or use a different port
PORT=8080 go run main.go

# Or for Docker
docker run -d -p 8080:3000 kg-dashboard
```

### Template Not Found

Ensure you're running the command from the `web/` directory, or use absolute paths in `main.go`.

### Dependencies Issues

```bash
# Clean and re-download dependencies
go mod tidy
go mod download
```

### Docker Build Issues

```bash
# Clean Docker build cache
docker builder prune -a

# Build with no cache
docker build --no-cache -t kg-dashboard .
```

## Performance

- **Cold start**: ~10ms
- **Request handling**: <1ms for static files
- **Memory usage**: ~15MB
- **Binary size**: ~8MB
- **Docker image size**: ~45MB (Alpine-based)

## Security

- Rate limiting prevents abuse (120 req/min per IP)
- Security headers protect against common attacks
- No user input processing = no XSS risk
- Read-only templates = no injection risk
- Docker container runs as non-root user
- CORS configured for specific origins only

## License

MIT License - see main project repository for details.

## Links

- **Main Project**: https://github.com/momokii/go-cli-notes
- **Cloud API**: https://cli-notes-api.kelanach.xyz/
- **Installation Docs**: https://github.com/momokii/go-cli-notes/blob/main/docs/INSTALL.md

## Support

For issues or questions:
1. Check the main project documentation
2. Open an issue on GitHub
3. Contact the maintainers
