#!/bin/bash

echo "Pulizia vecchi container..."
docker compose down --remove-orphans

echo " Build dei servizi..."
docker compose build

echo "Avvio del sistema..."
docker compose up -d

echo "Attendo 3 secondi per permettere ai servizi di avviarsi..."
sleep 3

echo "Log del client:"
docker compose logs -f client