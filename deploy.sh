#!/bin/bash
set -e

# Deployment script for EC2 instances
# This should be run on your EC2 instance

APP_NAME="jared-wallace-blog"
IMAGE_NAME="jw-blog"
CONTAINER_NAME="jw-blog-container"
DATA_PATH="/var/www/html/data"  # Using the EBS mount

echo "🚀 Starting deployment of $APP_NAME..."

# Ensure data directory exists on EBS volume
sudo mkdir -p $DATA_PATH
sudo chown $USER:$USER $DATA_PATH

# Stop and remove existing container if it exists
if [ $(docker ps -q -f name=$CONTAINER_NAME) ]; then
    echo "📦 Stopping existing container..."
    docker stop $CONTAINER_NAME
fi

if [ $(docker ps -aq -f name=$CONTAINER_NAME) ]; then
    echo "🗑️  Removing existing container..."
    docker rm $CONTAINER_NAME
fi

# Remove old image
if [ $(docker images -q $IMAGE_NAME) ]; then
    echo "🗑️  Removing old image..."
    docker rmi $IMAGE_NAME
fi

# Build new image
echo "🔨 Building new Docker image..."
docker build -t $IMAGE_NAME .

# Run new container
echo "🚀 Starting new container..."
docker run -d \
    --name $CONTAINER_NAME \
    --restart unless-stopped \
    -p 80:8080 \
    -v $DATA_PATH:/data \
    -e DB_PATH=/data/blog.db \
    -e PORT=8080 \
    --health-cmd="wget --quiet --tries=1 --spider http://localhost:8080/ || exit 1" \
    --health-interval=30s \
    --health-timeout=10s \
    --health-retries=3 \
    $IMAGE_NAME

echo "⏳ Waiting for container to be healthy..."
sleep 10

# Check if container is running
if [ $(docker ps -q -f name=$CONTAINER_NAME) ]; then
    echo "✅ Deployment successful! Container is running."
    docker logs --tail=20 $CONTAINER_NAME
else
    echo "❌ Deployment failed! Container is not running."
    docker logs $CONTAINER_NAME
    exit 1
fi

# Clean up unused images
echo "🧹 Cleaning up unused Docker images..."
docker image prune -f

echo "🎉 Deployment complete!"
