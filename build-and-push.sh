#!/bin/bash

# Docker Build and Push Script for Go Application
# Exit immediately if any command fails
set -e

# Default configuration
VERSION="latest"
IMAGE_NAME="k8s-demo"
DOCKER_USERNAME="iamnilotpal"
DOCKERFILE_PATH="infra/docker/Dockerfile"

# Colors for better output formatting
NC='\033[0m'
RED='\033[0;31m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'

# Helper functions for colored output
print_info() {
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

# Display usage information
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Build and push Docker images for the k8s-demo application"
    echo ""
    echo "Options:"
    echo "  -u, --username USERNAME    Docker Hub username (default: $DOCKER_USERNAME)"
    echo "  -n, --name IMAGE_NAME      Docker image name (default: $IMAGE_NAME)"
    echo "  -v, --version VERSION      Image version tag (default: $VERSION)"
    echo "  -f, --dockerfile PATH      Path to Dockerfile (default: $DOCKERFILE_PATH)"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                         # Build with defaults"
    echo "  $0 -v v1.2.3              # Build with specific version"
    echo "  $0 -u myuser -n myapp      # Build with custom user and image name"
    echo ""
}

# Verify Docker is running and accessible
check_docker() {
    print_info "Checking Docker daemon status..."

    if ! docker info >/dev/null 2>&1; then
        print_error "Docker daemon is not running or not accessible"
        print_error "Please start Docker Desktop or Docker service and try again"
        exit 1
    fi

    print_success "Docker daemon is running and accessible"
}

# Verify build context and Dockerfile exist
verify_build_context() {
    print_info "Verifying build context and Dockerfile..."

    # Check if we're in the correct directory (should contain go.mod)
    if [[ ! -f "go.mod" ]]; then
        print_error "go.mod not found. Please run this script from the project root directory"
        exit 1
    fi

    # Check if Dockerfile exists
    if [[ ! -f "$DOCKERFILE_PATH" ]]; then
        print_error "Dockerfile not found at: $DOCKERFILE_PATH"
        print_error "Please verify the Dockerfile path or use -f flag to specify correct path"
        exit 1
    fi

    print_success "Build context verified - found go.mod and Dockerfile"
}

# Login to Docker Hub with error handling
docker_login() {
    print_info "Authenticating with Docker Hub..."

    # Check if already logged in by testing a simple command
    if docker info | grep -q "Username"; then
        print_info "Already authenticated with Docker Hub"
        return 0
    fi

    # Attempt to login
    if ! docker login; then
        print_error "Failed to authenticate with Docker Hub"
        print_error "Please check your credentials and network connection"
        exit 1
    fi

    print_success "Successfully authenticated with Docker Hub"
}

# Build the Docker image with multi-stage optimization
build_image() {
    local full_image_name="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

    print_info "Building Docker image: ${full_image_name}"
    print_info "Using multi-stage Dockerfile: ${DOCKERFILE_PATH}"
    print_info "Build context: current directory (project root)"
    print_info ""
    print_info "This process will:"
    print_info "  1. Download Go dependencies and verify them"
    print_info "  2. Build the Go application with static linking"
    print_info "  3. Create minimal final image from scratch"
    print_info ""

    # Build with detailed progress and no cache for consistency
    if docker build \
        --file "$DOCKERFILE_PATH" \
        --tag "$full_image_name" \
        --progress=plain \
        --no-cache \
        .; then

        print_success "Successfully built Docker image: ${full_image_name}"

        # Display image information
        local image_size=$(docker images "$full_image_name" --format "table {{.Size}}" | tail -n +2)
        local image_id=$(docker images "$full_image_name" --format "table {{.ID}}" | tail -n +2)

        print_info "Image Details:"
        print_info "  Image ID: ${image_id}"
        print_info "  Size: ${image_size}"
        print_info "  Architecture: Multi-stage (builder + scratch runtime)"

    else
        print_error "Failed to build Docker image"
        print_error "Common issues:"
        print_error "  - Ensure you're in the project root directory"
        print_error "  - Check that go.mod and source files are present"
        print_error "  - Verify Docker has sufficient resources"
        exit 1
    fi
}

# Push the image to Docker Hub with retry logic
push_image() {
    local full_image_name="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

    print_info "Pushing image to Docker Hub: ${full_image_name}"

    # Push with retry on failure
    local max_retries=3
    local retry_count=0

    while [[ $retry_count -lt $max_retries ]]; do
        if docker push "$full_image_name"; then
            print_success "Successfully pushed image to Docker Hub"
            print_success "Image available at: https://hub.docker.com/r/${DOCKER_USERNAME}/${IMAGE_NAME}"
            print_info "Pull command: docker pull ${full_image_name}"
            return 0
        else
            retry_count=$((retry_count + 1))
            if [[ $retry_count -lt $max_retries ]]; then
                print_warning "Push failed, retrying... (attempt $retry_count/$max_retries)"
                sleep 5
            fi
        fi
    done

    print_error "Failed to push image after $max_retries attempts"
    print_error "Please check your network connection and Docker Hub permissions"
    exit 1
}

# Optionally clean up local images to save space
cleanup_local_image() {
    local full_image_name="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

    echo ""
    read -p "Remove local image to save disk space? (y/N): " -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Removing local image: ${full_image_name}"
        if docker rmi "$full_image_name" 2>/dev/null; then
            print_success "Local image removed successfully"
        else
            print_warning "Could not remove local image (may be in use)"
        fi
    else
        print_info "Keeping local image for future use"
    fi
}

# Parse command line arguments
parse_arguments() {
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
            -f|--dockerfile)
                DOCKERFILE_PATH="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo ""
                show_usage
                exit 1
                ;;
        esac
    done
}

# Main execution flow
main() {
    # Parse command line arguments first
    parse_arguments "$@"

    print_info "Starting Docker build and push process for k8s-demo application"
    print_info ""
    print_info "Configuration:"
    print_info "  Docker Username: ${DOCKER_USERNAME}"
    print_info "  Image Name: ${IMAGE_NAME}"
    print_info "  Version Tag: ${VERSION}"
    print_info "  Dockerfile: ${DOCKERFILE_PATH}"
    print_info "  Full Image: ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"
    echo ""

    # Execute all build steps
    verify_build_context
    check_docker
    docker_login
    build_image
    push_image
    cleanup_local_image

    print_success "Build and push process completed successfully!"
    print_info ""
    print_info "Next steps:"
    print_info "  1. Update your Helm values.yaml with the new image tag"
    print_info "  2. Deploy to Kubernetes: helm upgrade k8s-demo ./infra/helm"
    print_info "  3. Verify deployment: kubectl get pods"
}

# Run the main function with all arguments
main "$@"