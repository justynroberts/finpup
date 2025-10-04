#!/bin/bash

echo "🔍 Validating Finton Build..."
echo ""

# Check binary exists
if [ -f "finton" ]; then
    echo "✅ Binary exists: $(ls -lh finton | awk '{print $5}')"
else
    echo "❌ Binary not found"
    exit 1
fi

# Check binary is executable
if [ -x "finton" ]; then
    echo "✅ Binary is executable"
else
    echo "❌ Binary is not executable"
    exit 1
fi

# Run tests
echo ""
echo "🧪 Running tests..."
if go test ./... > /dev/null 2>&1; then
    echo "✅ All tests passing"
else
    echo "❌ Some tests failing"
    go test ./...
    exit 1
fi

# Check key files
echo ""
echo "📄 Checking documentation..."
for file in README.md QUICKSTART.md CLAUDE.md PROJECT_SUMMARY.md; do
    if [ -f "$file" ]; then
        echo "✅ $file exists"
    else
        echo "❌ $file missing"
    fi
done

# Check structure
echo ""
echo "🏗️  Checking project structure..."
for dir in cmd/finton internal/buffer internal/editor internal/ui internal/ai internal/config internal/highlight pkg/themes; do
    if [ -d "$dir" ]; then
        echo "✅ $dir/ exists"
    else
        echo "❌ $dir/ missing"
    fi
done

echo ""
echo "📊 Project Stats:"
echo "   Lines of Go: $(find . -name '*.go' -type f -exec wc -l {} + | tail -1 | awk '{print $1}')"
echo "   Go files: $(find . -name '*.go' -type f | wc -l | tr -d ' ')"
echo "   Test files: $(find . -name '*_test.go' -type f | wc -l | tr -d ' ')"

echo ""
echo "✨ Finton is ready!"
echo ""
echo "📦 Install: ./install.sh"
echo "🚀 Run: finton <filename>"
echo ""
