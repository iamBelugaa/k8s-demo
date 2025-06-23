#!/bin/bash

# Docker Build and Push Script for Go Application
# This script automates the process of building and pushing your Docker image

set -e  # Exit on any error

# Configuration - Modify these variables for your specific setup
VERSION="latest"
IMAGE_NAME="k8s-demo"
DOCKER_USERNAME="iamnilotpal"

# Colors for output formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    print_status "Checking Docker status..."
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Login to Docker Hub
docker_login() {
    print_status "Logging into Docker Hub..."
    if ! docker login; then
        print_error "Failed to login to Docker Hub"
        exit 1
    fi
    print_success "Successfully logged into Docker Hub"
}

# Build the Docker image
build_image() {
    local full_image_name="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

    print_status "Building Docker image: ${full_image_name}"
    print_status "Using Dockerfile at: infra/docker/Dockerfile"
    print_status "Build context: current directory (project root)"
    print_status "This may take a few minutes..."

    # Build with progress output, specifying Dockerfile location
    # -f flag specifies the Dockerfile path relative to build context
    # Build context is current directory (.) which should be project root
    if docker build \
        --file infra/docker/Dockerfile \
        --tag "${full_image_name}" \
        --progress=plain \
        --no-cache \
        .; then
        print_success "Successfully built image: ${full_image_name}"
    else
        print_error "Failed to build Docker image"
        print_error "Make sure you're running this script from the project root directory"
        exit 1
    fi

    # Show image size
    local image_size=$(docker images "${full_image_name}" --format "table {{.Size}}" | tail -n +2)
    print_status "Final image size: ${image_size}"
}

# Push the image to Docker Hub
push_image() {
    local full_image_name="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

    print_status "Pushing image to Docker Hub: ${full_image_name}"

    if docker push "${full_image_name}"; then
        print_success "Successfully pushed image to Docker Hub"
        print_success "Image available at: https://hub.docker.com/r/${DOCKER_USERNAME}/${IMAGE_NAME}"
    else
        print_error "Failed to push image to Docker Hub"
        exit 1
    fi
}

# Clean up local images (optional)
cleanup() {
    local full_image_name="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

    read -p "Do you want to remove the local image to save space? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Removing local image..."
        docker rmi "${full_image_name}" || print_warning "Could not remove local image"
    fi
}

# Show usage instructions
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -u, --username USERNAME    Docker Hub username"
    echo "  -n, --name IMAGE_NAME      Image name"
    echo "  -v, --version VERSION      Image version/tag (default: latest)"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 -u myusername -n myapp -v v1.0.0"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--username)
            DOCKER_USERNAME="$2"
            shift 2
            ;;
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main execution flow
main() {
    print_status "Starting Docker build and push process..."
    print_status "Configuration:"
    print_status "  Docker Username: ${DOCKER_USERNAME}"
    print_status "  Image Name: ${IMAGE_NAME}"
    print_status "  Version: ${VERSION}"
    echo ""

    # Run all the steps
    check_docker
    docker_login
    build_image
    push_image
    cleanup

    print_success "Process completed successfully!"
    print_status "Your image is now available on Docker Hub"
}

# Run the main function
main "$@"