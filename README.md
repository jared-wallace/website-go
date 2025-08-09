# Jared Wallace's Personal Website

A minimalist blog and portfolio website built with Go, emphasizing performance, simplicity, and clean design inspired by mid-century modern aesthetics.

## Features

- **Minimalist Design**: Clean, focused interface with mid-century modern styling
- **High Performance**: Optimized for speed with minimal JavaScript and efficient caching
- **Responsive**: Works beautifully on all device sizes
- **Docker Ready**: Containerized for easy deployment and scaling
- **SQLite Database**: Lightweight, file-based database perfect for low-traffic sites
- **RSS News Ticker**: Live BBC World News feed integration
- **SEO Optimized**: Clean URLs, semantic HTML, and proper meta tags

## Tech Stack

- **Backend**: Go 1.21+ with Gorilla Mux router
- **Database**: SQLite3 with automatic migrations
- **Frontend**: Vanilla HTML/CSS/JavaScript (no frameworks)
- **Deployment**: Docker with multi-stage builds
- **Infrastructure**: AWS (EC2, ALB, Route53, EBS)

## Quick Start

### Prerequisites

- Go 1.21+
- Docker (optional, for containerized deployment)
- Make (optional, for using Makefile commands)

### Development Setup

1. **Clone and setup:**
   ```bash
   git clone <your-repo>
   cd jared-wallace-blog
   make setup  # Creates directory structure
   make tidy   # Download dependencies
   ```

2. **Initialize with sample data:**
   ```bash
   make init-db
   ```

3. **Run locally:**
   ```bash
   make run
   ```
   
   Visit `http://localhost:8080` to see your site.

### Docker Development

```bash
make docker-dev
```

This will start the application with Docker Compose, mounting your local files for easy development.

## Project Structure

```
├── cmd/server/          # Application entry point
├── templates/           # HTML templates
│   ├── layout.html      # Base layout
│   ├── home.html        # Homepage
│   ├── blog.html        # Blog listing
│   ├── post.html        # Individual post
│   ├── projects.html    # Projects page
│   └── about.html       # About page
├── static/              # Static assets
│   ├── css/             # Stylesheets
│   ├── js/              # JavaScript files
│   └── images/          # Images and media
├── data/                # SQLite database and persistent data
├── scripts/             # Database migrations and utilities
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Development environment
├── Makefile            # Build and deployment commands
└── deploy.sh           # Production deployment script
```

## Database Schema

The application uses SQLite with two main tables:

- **posts**: Blog posts with title, content, excerpt, slug, and publishing status
- **projects**: Portfolio projects with descriptions, URLs, and technology tags

Migrations run automatically on startup.

## Deployment

### Production Deployment on AWS

The included Terraform configuration sets up:

- VPC with public subnets
- Application Load Balancer with SSL termination
- Auto Scaling Group with health checks
- Route53 DNS management
- EBS volume for persistent data storage

### Deploy to your EC2 instance:

```bash
make deploy
```

This will:
1. Build a production Docker image
2. Transfer it to your EC2 instance
3. Run the deployment script
4. Start the new container with zero downtime

### Manual Docker Deployment

```bash
# Build the image
docker build -t jw-blog .

# Run with persistent data
mkdir -p ./data
docker run -d \
  --name jw-blog \
  --restart unless-stopped \
  -p 80:8080 \
  -v $(pwd)/data:/data \
  jw-blog
```

## Configuration

Configuration is handled through environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database file path (default: /data/blog.db)

## Content Management

### Adding Blog Posts

Blog posts are stored in the SQLite database. You can add them by:

1. **Direct database insertion** (for development):
   ```sql
   INSERT INTO posts (title, content, excerpt, published, slug) 
   VALUES ('Your Title', '<p>Your content...</p>', 'Brief excerpt', 1, 'your-slug');
   ```

2. **API endpoint** (add authentication for production use)

### Adding Projects

Similar to blog posts, projects are managed through the database:

```sql
INSERT INTO projects (title, description, url, image_url, technologies) 
VALUES ('Project Name', 'Description', 'https://example.com', '/static/images/project.jpg', 'Go, Docker, AWS');
```

## Performance Features

- **Critical CSS inlined** for faster first paint
- **Image optimization** with WebP support and fallbacks
- **Efficient caching** headers for static assets
- **Minimal JavaScript** only where necessary
- **Compressed responses** with gzip
- **CDN-ready** static asset serving

## Design Philosophy

This website embodies a minimalist approach inspired by mid-century modern design:

- **Clean typography** with careful hierarchy
- **Generous whitespace** for improved readability
- **Subtle animations** that enhance without distracting
- **Purposeful color palette** with accessible contrast ratios
- **Mobile-first responsive design**

## Security Considerations

- Input sanitization for all user-facing content
- HTTPS enforcement (handled by ALB)
- Minimal attack surface with simple architecture
- Regular dependency updates
- Proper CORS and security headers

## Monitoring and Logs

- Docker health checks for container monitoring
- Application logs to stdout for easy collection
- Performance monitoring through browser DevTools
- AWS CloudWatch integration for production metrics

## Contributing

This is a personal website, but if you find bugs or have suggestions:

1. Open an issue describing the problem
2. Fork the repository
3. Create a feature branch
4. Submit a pull request

## License

This project is personal and not open for general use, but feel free to reference the code for your own projects.

## Support

For questions or issues, please open a GitHub issue or contact me directly.

---

Built with ❤️ and minimal dependencies.
