#!/bin/bash

# Variables
DOCKER_IMAGE="cchalop1/justdeploy"
GIT_TAG=$(git describe --tags --abbrev=0 || echo "latest")

# Affiche le message de bienvenue
echo "================================"
echo "Construction de l'image Docker JustDeploy"
echo "Version: $GIT_TAG"
echo "================================"

# Vérifier si Docker est installé
if ! command -v docker &> /dev/null; then
    echo "Erreur: Docker n'est pas installé. Veuillez l'installer avant de continuer."
    exit 1
fi

# Construction de l'image Docker
echo "🔨 Construction de l'image Docker $DOCKER_IMAGE:$GIT_TAG..."
docker build --build-arg VERSION=$GIT_TAG -t $DOCKER_IMAGE:$GIT_TAG -t $DOCKER_IMAGE:latest -f Dockerfile  ../

# Vérifier si la construction a réussi
if [ $? -ne 0 ]; then
    echo "❌ Erreur lors de la construction de l'image Docker."
    exit 1
fi

echo "✅ Construction réussie!"

# Demander à l'utilisateur s'il souhaite pousser l'image
read -p "Voulez-vous pousser l'image sur Docker Hub? (o/n): " PUSH_IMAGE

if [ "$PUSH_IMAGE" = "o" ] || [ "$PUSH_IMAGE" = "O" ]; then
    echo "🚀 Envoi de l'image vers Docker Hub..."
    
    # Vérifier si l'utilisateur est connecté à Docker Hub
    docker login
    
    # Pousser les images
    docker push $DOCKER_IMAGE:$GIT_TAG
    docker push $DOCKER_IMAGE:latest
    
    echo "✅ Images $DOCKER_IMAGE:$GIT_TAG et $DOCKER_IMAGE:latest poussées avec succès!"
else
    echo "ℹ️ L'image n'a pas été poussée sur Docker Hub."
fi

echo "================================"
echo "✅ Terminé!"
echo "Image Docker: $DOCKER_IMAGE:$GIT_TAG"
echo "================================" 