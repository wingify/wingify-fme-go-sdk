#!/bin/bash

# Copyright 2025-2026 Wingify Software Pvt. Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "🔧 Setting up Git Hooks for Wingify FME Go SDK"
echo "=========================================="

# Get the project root directory (go up two levels from scripts directory)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GIT_HOOKS_DIR="$PROJECT_ROOT/git-hooks"
GIT_HOOKS_TARGET="$PROJECT_ROOT/.git/hooks"

echo "Project root: $PROJECT_ROOT"
echo "Git hooks source: $GIT_HOOKS_DIR"
echo "Git hooks target: $GIT_HOOKS_TARGET"

# Check if we're in a git repository
if [ ! -d "$PROJECT_ROOT/.git" ]; then
    echo "❌ Error: Not in a git repository!"
    exit 1
fi

# Create .git/hooks directory if it doesn't exist
mkdir -p "$GIT_HOOKS_TARGET"

# Copy git hooks
echo ""
echo "📋 Installing git hooks..."

# Copy pre-push hook
if [ -f "$GIT_HOOKS_DIR/pre-push" ]; then
    cp "$GIT_HOOKS_DIR/pre-push" "$GIT_HOOKS_TARGET/pre-push"
    chmod +x "$GIT_HOOKS_TARGET/pre-push"
    echo "✓ Installed pre-push hook"
else
    echo "❌ pre-push hook not found!"
fi

# Copy commit-msg hook
if [ -f "$GIT_HOOKS_DIR/commit-msg" ]; then
    cp "$GIT_HOOKS_DIR/commit-msg" "$GIT_HOOKS_TARGET/commit-msg"
    chmod +x "$GIT_HOOKS_TARGET/commit-msg"
    echo "✓ Installed commit-msg hook"
else
    echo "❌ commit-msg hook not found!"
fi

# Make sure all scripts are executable
echo ""
echo "🔐 Setting executable permissions..."

# Make Node.js scripts executable
find "$GIT_HOOKS_DIR/scripts" -name "*.js" -exec chmod +x {} \;
echo "✓ Made Node.js scripts executable"

# Make shell scripts executable
# Note: add-copyright.sh has been removed, use Node.js version instead

# Test Node.js availability
echo ""
echo "🧪 Testing dependencies..."

if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo "✓ Node.js is available: $NODE_VERSION"
else
    echo "❌ Node.js is not installed!"
    echo "   Please install Node.js to use the git hooks."
    exit 1
fi

if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "✓ Go is available: $GO_VERSION"
else
    echo "❌ Go is not installed!"
    echo "   Please install Go to use the git hooks."
    exit 1
fi

# Test the license check script
echo ""
echo "🧪 Testing license check script..."
if node "$GIT_HOOKS_DIR/scripts/check-license.js" &> /dev/null; then
    echo "✓ License check script is working"
else
    echo "❌ License check script test failed"
fi

echo ""
echo "🎉 Git hooks setup completed successfully!"
echo ""
echo "The following hooks are now active:"
echo "  • pre-push: Runs copyright check, go fmt, go vet, and go build"
echo "  • commit-msg: Validates commit message format"
echo ""
echo "You can now:"
echo "  • Run 'node git-hooks/scripts/add-copyright.js' to add copyright headers to all Go files"
echo "  • Run './git-hooks/scripts/run-tests.sh' to run the test suite"
echo "  • Git hooks will automatically run on commit and push"
echo ""
echo "To disable hooks temporarily:"
echo "  git config core.hooksPath /dev/null"
echo ""
echo "To re-enable hooks:"
echo "  git config --unset core.hooksPath"
