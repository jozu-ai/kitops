#!/usr/bin/env bash

source ./.demo-magic/demo-magic.sh


########################
# Configure the options
########################

#
# speed at which to simulate typing. bigger num = faster
#
 TYPE_SPEED=25

#
# custom prompt
#
# see http://www.tldp.org/HOWTO/Bash-Prompt-HOWTO/bash-prompt-escape-sequences.html for escape sequences
#
DEMO_PROMPT="${GREEN}➜ ${CYAN}\W ${COLOR_RESET}"

# text color
# DEMO_CMD_COLOR=$BLACK

if command -v docker; then
  DOCKER=docker
elif command -v podman; then
  DOCKER=podman
else
  echo "You need docker or podman installed for this demo"
  exit 1
fi

# hide the evidence
clear


pe "./kit version"

# Let's check if there are any model kits locally
pe "./kit list"

# clean the local model kits and check again
pe "rm -rf ~/.kitops"
pe "./kit list"

pe "./kit build --help"

# Let's build the onnx model
pe "./kit build ../examples/onnx -t localhost:5050/test-repo:test-tag"

# Let's check if the model kit is built
pe "./kit list"

# run a local registry
pe "$DOCKER run --name registry --rm -d -p 5050:5050 -e REGISTRY_HTTP_ADDR=:5050 registry"

# Let's push the model to the local registry
pe "./kit push localhost:5050/test-repo:test-tag --http"
# Let's check if the model is pushed
pe "./kit list localhost:5050/test-repo --http"
# clean the local models and check again
pe "rm -rf ~/.kitops"
pe "./kit list"

# Let's pull the model kit to the local registry
pe "./kit pull localhost:5050/test-repo:test-tag --http"

pe "./kit list"

