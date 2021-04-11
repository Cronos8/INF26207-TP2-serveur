# INF26207-TP2-serveur

## Pour commencer

Simulation d'un serveur UDP avec les sockets dans le cadre d'un travail pratique du cours de téléinformatique à l'UQAR.

### Pré-requis

- Go v1.16.3

Testé sous :
- macOS 11.2.3
- windows 10 

### Installation

Dans la racine du répertoire, exécuter la commande : ``go build .``.

## Démarrage

À la racine du répertoire :

Exécuter la commande :``./INF26207-TP2-serveur IPServeur:PortServeur cheminDuFichierAEnvoyer``

Exemple : ``./INF26207-TP2-serveur 127.0.0.1:22222 "testfiles/alpaga.jpeg"``

Le dossier ``testfiles/`` contient des fichiers de tests.

## Suppression de l'exécutable 

À la racine du répertoire :

Exécuter la commande :``go clean``.

## Auteurs

Alexandre Nguyen