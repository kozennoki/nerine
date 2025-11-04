#!/bin/bash

set -e

# Install mise
curl -sSfL https://mise.run | sh

# Add mise activation to bashrc
echo 'eval "$(~/.local/bin/mise activate bash)"' >> ~/.bashrc

# Install mise tools
~/.local/bin/mise trust
~/.local/bin/mise install

# Install project tools via task
~/.local/bin/mise exec -- task install-tools
