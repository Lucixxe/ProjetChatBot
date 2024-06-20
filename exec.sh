#!/bin/bash

handler_sigint() {
    echo "Arrêt du programme..."
    kill $child_pid
    wait $child_pid
    echo "Programme arrêté."
}

# Associe la fonction handler_sigint à l'interruption du script
trap handler_sigint SIGINT

# Vérifiez si la base de données existe
if [ -f chatbot.db ]
then
    echo "La base de données existe déjà. Passage à l'étape suivante."
else
    echo "Initialisation de la base de données SQLite..."
    sqlite chatbot.db < schema.sql
fi

# On compile
echo "Compilation du programme..."
go get
go build
go run chatbot &

child_pid=$!

# On attend au cas où... À modifier si nécessaire
sleep 1

xdg-open 'http://localhost:3333'

wait $child_pid
