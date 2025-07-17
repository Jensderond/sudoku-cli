#!/bin/bash

# Sudoku CLI Install Script
# This script automatically detects your system and installs the appropriate binary

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# GitHub repository information
GITHUB_REPO="jensderond/sudoku-cli"
BINARY_NAME="sudoku"

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

# Function to detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "macos"
            ;;
        Linux*)
            echo "linux"
            ;;
        CYGWIN*|MINGW*|MSYS*)
            echo "windows"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Function to detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        i386|i686)
            echo "386"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to download and install
install_binary() {
    local os=$1
    local arch=$2
    local version=${3:-"latest"}
    
    print_info "Detected system: $os-$arch"
    
    # Determine file extension
    local file_ext="tar.gz"
    if [ "$os" = "windows" ]; then
        file_ext="zip"
    fi
    
    # Construct download URL
    local filename="sudoku-${os}-${arch}.${file_ext}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/${version}/download/${filename}"
    
    print_info "Downloading from: $download_url"
    
    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Download the file
    if command_exists curl; then
        curl -fsSL "$download_url" -o "$filename"
    elif command_exists wget; then
        wget -q "$download_url" -O "$filename"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    # Extract the file
    if [ "$file_ext" = "tar.gz" ]; then
        tar -xzf "$filename"
    elif [ "$file_ext" = "zip" ]; then
        if command_exists unzip; then
            unzip -q "$filename"
        else
            print_error "unzip is not available. Please install it."
            exit 1
        fi
    fi
    
    # Determine binary name
    local binary_file="$BINARY_NAME"
    if [ "$os" = "windows" ]; then
        binary_file="${BINARY_NAME}.exe"
    fi
    
    # Make binary executable
    chmod +x "$binary_file"
    
    # Determine installation directory
    local install_dir="/usr/local/bin"
    if [ "$os" = "windows" ]; then
        install_dir="$HOME/bin"
        mkdir -p "$install_dir"
    fi
    
    # Install the binary
    if [ -w "$install_dir" ]; then
        mv "$binary_file" "$install_dir/$binary_file"
        print_info "Successfully installed $BINARY_NAME to $install_dir"
    else
        print_warning "No write permission to $install_dir. Trying with sudo..."
        sudo mv "$binary_file" "$install_dir/$binary_file"
        print_info "Successfully installed $BINARY_NAME to $install_dir (with sudo)"
    fi
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$temp_dir"
    
    # Verify installation
    if command_exists "$BINARY_NAME"; then
        print_info "Installation successful! You can now run '$BINARY_NAME'"
        print_info "Run '$BINARY_NAME --help' to get started."
    else
        print_warning "Installation completed, but '$BINARY_NAME' is not in your PATH."
        if [ "$os" = "windows" ]; then
            print_info "Please add $install_dir to your PATH environment variable."
        else
            print_info "Please add $install_dir to your PATH or restart your terminal."
        fi
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -v, --version VERSION    Install specific version (default: latest)"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                      Install latest version"
    echo "  $0 -v v1.0.0           Install version v1.0.0"
}

# Main function
main() {
    local version="latest"
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                version="$2"
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
    
    print_info "Starting Sudoku CLI installation..."
    
    # Detect system
    local os
    local arch
    os=$(detect_os)
    arch=$(detect_arch)
    
    # Install binary
    install_binary "$os" "$arch" "$version"
}

# Run main function
main "$@"
