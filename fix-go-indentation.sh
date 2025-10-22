#!/bin/bash

# Fix Go files indentation to 4 spaces

echo "🔧 Fixing Go files indentation to 4 spaces..."

# Function to convert tabs to 4 spaces in Go files
fix_go_file() {
    local file=$1
    echo "  📝 Fixing: $file"
    
    # Convert tabs to 4 spaces
    sed -i '' 's/\t/    /g' "$file"
}

# Find all Go files and fix them
find backend -name "*.go" -type f | while read -r file; do
    fix_go_file "$file"
done

echo "✅ All Go files fixed to use 4-space indentation"