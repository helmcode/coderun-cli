#!/bin/bash

# CodeRun CLI Installer
# This script installs the latest version of CodeRun CLI

set -e

# Variables
GITHUB_REPO="helmcode/coderun-cli"
BINARY_NAME="coderun"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""
    
    case "$OSTYPE" in
        linux*)   os="linux" ;;
        darwin*)  os="darwin" ;;
        msys*|cygwin*|mingw*) os="windows" ;;
        *)        
            print_error "Unsupported operating system: $OSTYPE"
            exit 1
            ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [[ -z "$version" ]]; then
        print_error "Failed to get latest version"
        exit 1
    fi
    
    echo "$version"
}

# Download and install binary
install_binary() {
    local platform="$1"
    local version="$2"
    local temp_dir
    local binary_url
    local binary_file
    
    # Create temporary directory
    temp_dir=$(mktemp -d)
    
    # Determine binary file name
    if [[ "$platform" == *"windows"* ]]; then
        binary_file="${BINARY_NAME}-${platform}.exe"
    else
        binary_file="${BINARY_NAME}-${platform}"
    fi
    
    # Construct download URL
    binary_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${binary_file}"
    
    print_info "Downloading CodeRun CLI ${version} for ${platform}..."
    print_info "URL: ${binary_url}"
    
    # Download binary
    if ! curl -L -o "${temp_dir}/${BINARY_NAME}" "$binary_url"; then
        print_error "Failed to download binary"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Make executable
    chmod +x "${temp_dir}/${BINARY_NAME}"
    
    # Install binary
    print_info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    
    if [[ "$EUID" -ne 0 ]] && [[ ! -w "$INSTALL_DIR" ]]; then
        print_warning "Need sudo privileges to install to ${INSTALL_DIR}"
        sudo mv "${temp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        mv "${temp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    
    print_info "CodeRun CLI installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        print_info "Installation verified!"
        print_info "Version: $($BINARY_NAME --version)"
        print_info ""
        print_info "Get started with:"
        echo "  coderun login"
        echo "  coderun deploy nginx:latest --name my-app --http-port 80"
    else
        print_error "Installation verification failed"
        print_warning "You may need to add ${INSTALL_DIR} to your PATH"
        exit 1
    fi
}

# Main installation process
main() {
    print_info "CodeRun CLI Installer"
    print_info "====================="
    
    # Check if curl is available
    if ! command -v curl &> /dev/null; then
        print_error "curl is required but not installed"
        exit 1
    fi
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    print_info "Detected platform: $platform"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    print_info "Latest version: $version"
    
    # Install binary
    install_binary "$platform" "$version"
    
    # Verify installation
    verify_installation
}

# Run main function
main "$@"
