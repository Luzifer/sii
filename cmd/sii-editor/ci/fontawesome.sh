#!/bin/bash
set -euxo pipefail

# Ensure deletion of older version
rm -rf frontend/fontawesome

# Download and unpack fontawesome-free
curl -sSfLo frontend/fa.zip "https://use.fontawesome.com/releases/v${FA_VERSION}/fontawesome-free-${FA_VERSION}-web.zip"
unzip frontend/fa.zip -d frontend
rm frontend/fa.zip

# Move to generic path
mv frontend/fontawesome-free-${FA_VERSION}-web frontend/fontawesome
