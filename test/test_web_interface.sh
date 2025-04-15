#!/bin/bash

# Test script for the Snap web interface
# This script creates a test repository, adds some files, commits them,
# creates some issues, and then starts the web interface.

# Exit on error
set -e

# Create a test directory
TEST_DIR=$(mktemp -d -t snap-web-test-XXXXXX)
echo "Created test directory: $TEST_DIR"
cd "$TEST_DIR"

# Build snap
echo "Building snap..."
go build -o snap github.com/stanlocht/snap

# Initialize repository
echo "Initializing repository..."
./snap init

# Create some files
echo "Creating files..."
echo "# Test Repository" > README.md
echo "This is a test file." > test.txt
echo "Another test file." > another.txt

# Add files
echo "Adding files..."
./snap add README.md test.txt another.txt

# Commit files
echo "Committing files..."
./snap commit -m "✨ Initial commit" -a "testuser" -e "test@example.com"

# Create some issues
echo "Creating issues..."
./snap issue new -t "Test issue 1" -d "This is a test issue."
./snap issue new -t "Test issue 2" -d "This is another test issue."
./snap issue assign 1 testuser
./snap issue close 2

# Start web interface
echo "Starting web interface..."
echo "Press Ctrl+C to stop the server"
./snap web --open

# Clean up
# Note: This will only run if the user manually stops the server with Ctrl+C
echo "Cleaning up..."
cd -
rm -rf "$TEST_DIR"
echo "Done!"
