#!/bin/bash

# PMGO Installation Script
# This script installs the latest version of PMGO

set -e

# Configuration
REPO="struCoder/pmgo"
BINARY_NAME="pmgo"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os arch

    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "$arch" in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        armv7l) arch="arm" ;;
        i386|i686) arch="386" ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    case "$os" in
        linux|darwin) ;;
        mingw*|msys*|cygwin*) os="windows" ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
              grep '"tag_name":' | \
              sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        log_error "Failed to get latest version"
        exit 1
    fi

    echo "$version"
}

# Download and install binary
install_binary() {
    local version platform url temp_dir

    version=$(get_latest_version)
    platform=$(detect_platform)

    log_info "Installing PMGO $version for $platform"

    # Create temporary directory
    temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT

    # Construct download URL
    url="https://github.com/$REPO/releases/download/$version/pmgo-$platform"
    if [[ $platform == *"windows"* ]]; then
        url="${url}.exe"
    fi

    log_info "Downloading from $url"

    # Download binary
    if ! curl -L "$url" -o "$temp_dir/$BINARY_NAME"; then
        log_error "Failed to download binary"
        exit 1
    fi

    # Make executable
    chmod +x "$temp_dir/$BINARY_NAME"

    # Install to destination
    if [ ! -w "$INSTALL_DIR" ]; then
        log_warning "Need sudo permissions to install to $INSTALL_DIR"
        sudo mv "$temp_dir/$BINARY_NAME" "$INSTALL_DIR/"
    else
        mv "$temp_dir/$BINARY_NAME" "$INSTALL_DIR/"
    fi

    log_success "PMGO installed to $INSTALL_DIR/$BINARY_NAME"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        log_success "PMGO is installed and working: $version"
    else
        log_warning "PMGO binary not found in PATH. You may need to add $INSTALL_DIR to your PATH."
        log_info "Add this line to your shell profile (.bashrc, .zshrc, etc.):"
        log_info "export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

# Create config directory
setup_config() {
    local config_dir="$HOME/.pmgo"

    if [ ! -d "$config_dir" ]; then
        mkdir -p "$config_dir"
        log_info "Created config directory: $config_dir"
    fi

    # Create default config if not exists
    local config_file="$config_dir/config.yaml"
    if [ ! -f "$config_file" ]; then
        cat > "$config_file" << 'EOF'
# PMGO Configuration File
server:
  host: "localhost"
  port: 9876
  web_port: 8080

logging:
  level: "info"
  format: "text"

processes:
  default_restart_policy: "always"
  max_restart_attempts: 5
  restart_delay: "1s"
EOF
        log_info "Created default config: $config_file"
    fi
}

# Main installation function
main() {
    log_info "Starting PMGO installation..."

    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        log_error "curl is required but not installed"
        exit 1
    fi

    # Install binary
    install_binary

    # Setup configuration
    setup_config

    # Verify installation
    verify_installation

    log_success "PMGO installation completed!"
    log_info "Run 'pmgo --help' to get started"
    log_info "Documentation: https://github.com/$REPO"
}

# Show usage
show_usage() {
    cat << EOF
PMGO Installation Script

Usage: $0 [OPTIONS]

OPTIONS:
    -h, --help          Show this help message
    -d, --dir DIR       Installation directory (default: /usr/local/bin)

ENVIRONMENT VARIABLES:
    INSTALL_DIR         Installation directory

EXAMPLES:
    # Install to default location
    $0

    # Install to custom directory
    $0 --dir /opt/bin

    # Install to custom directory via environment
    INSTALL_DIR=/opt/bin $0

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Run main function
main