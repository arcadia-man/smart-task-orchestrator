#!/bin/bash

# Format Go files with 4-space indentation

echo "🔧 Formatting Go files with 4-space indentation..."

# Function to format a Go file with 4 spaces
format_go_file() {
    local file=$1
    echo "  📝 Formatting: $file"
    
    # First run gofmt to standardize
    gofmt -w "$file"
    
    # Then convert tabs to 4 spaces
    sed -i '' 's/\t/    /g' "$file"
}

# Find all Go files and format them
find . -name "*.go" -type f | while read -r file; do
    format_go_file "$file"
done

echo "✅ All Go files formatted with 4-space indentation"