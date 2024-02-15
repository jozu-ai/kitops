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

# hide the evidence
clear


pe "./jmm version" 

# Let's check if there are any models locally
pe "./jmm models"

# clean the local models and check again
pe "rm -rf ~/.jozu"
pe "./jmm models"

pe "./jmm build --help"

# Let's build the onnx model
pe "./jmm build ../examples/onnx -t localhost:5050/test-repo:test-tag"

# Let's check if the model is built
pe "./jmm models"

# run a local registry
pe "docker run --name registry --rm -d -p 5050:5050 -e REGISTRY_HTTP_ADDR=:5050 registry" 

# Let's push the model to the local registry
pe "./jmm push localhost:5050/test-repo:test-tag --http"
# Let's check if the model is pushed
pe "./jmm models localhost:5050/test-repo --http"
# clean the local models and check again
pe "rm -rf ~/.jozu"
pe "./jmm models"

# Let's pull the model to the local registry
pe "./jmm pull localhost:5050/test-repo:test-tag --http"

pe "./jmm models"

