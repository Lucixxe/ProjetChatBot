#!/bin/bash

mkdir -p ollama
cd ollama

# Download ollama
# curl -L https://ollama.com/download/ollama-linux-amd64 -o ollama

chmod u+x ollama

./ollama serve &

sleep 5

./ollama run llama3 &
