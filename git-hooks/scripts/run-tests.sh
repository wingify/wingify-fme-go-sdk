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

echo "🧪 Running Wingify FME Go SDK Test Suite"
echo "===================================="

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "   Please install Go to run tests."
    exit 1
fi

# Get Go version
GO_VERSION=$(go version)
echo "✓ Go is available: $GO_VERSION"
echo ""

# Run tests with different options
echo "Running comprehensive test suite..."
echo ""

# Run all tests with verbose output
echo "📋 Running all tests with verbose output:"
go test ./test/... -v

# Check exit code
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All tests passed successfully!"
    echo ""
    echo "Test Summary:"
    echo "  • E2E Tests: Feature flag functionality"
    echo "  • Unit Tests: Segmentation, operators, and validation"
    echo "  • Integration Tests: Storage and network components"
    echo ""
    echo "You can also run:"
    echo "  • go test ./test/e2e/... -v    # E2E tests only"
    echo "  • go test ./test/unit/... -v    # Unit tests only"
    echo "  • go test ./test/... -cover     # With coverage"
    echo "  • go test ./test/... -race      # With race detection"
else
    echo ""
    echo "❌ Some tests failed!"
    echo "   Please check the output above for details."
    exit 1
fi
