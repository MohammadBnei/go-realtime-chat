# Projet pédagogique : création d'une API de chat en temps réel avec gestion des utilisateurs et JWT

## Objectif du projet
L'objectif de ce projet est d'apprendre à créer une API de chat en temps réel en utilisant la bibliothèque broadcast et en apprenant les concepts de goroutine, de channels et d'interfaces en Go. En plus de cela, il faut intégrer les fonctionnalités de gestion des utilisateurs et d'authentification JWT.

## Prérequis 
- Connaissance de base de Go
- Connaissance de base des concepts de programmation orientée objet
- Connaissance de base des concepts de concurrency en Go
- Connaissance de base de JWT

## Tâches du projet 
1. Création de la structure de base pour l'API de chat en temps réel.
2. Implémentation de la gestion des utilisateurs :
    1. Ajouter un utilisateur
    2. Supprimer un utilisateur
    3. Authentification des utilisateurs avec JWT
3. Implémentation de la gestion des salons :
    1. Créer un salon
    2. Supprimer un salon
4. Implémentation de la gestion des messages :
    1. Envoyer un message à tous les utilisateurs d'un salon
5. Implémentation de l'API de chat en temps réel :
    1. Gérer les utilisateurs (CRUD + Login)
    2. Gérer les salons en utilisant les routes HTTP : /rooms/create, /rooms/delete
    3. Permettre aux utilisateurs de joindre un salon en utilisant la route HTTP : /rooms/join
    4. Permettre aux utilisateurs d'envoyer un message dans un salon en utilisant la route HTTP : /rooms/message
    5. Mettre en place du SSE pour la communication en temps réel entre les utilisateurs et les salons.

## Livrables :
- Le code source de l'API de chat en temps réel en Go avec les packages pour la gestion des salons, la gestion des utilisateurs et l'API de chat en temps réel.
- Un document expliquant les fonctionnalités implémentées par les différents membres du groupe.

## Critères d'évaluation :
- Qualité du code (bonne gestion des erreurs, organisation du code, lisibilité, etc.)
- Fonctionnalités de l'API de chat en temps réel implémentées correctement.
- Utilisation correcte des fonctionnalités de gestion des utilisateurs et d'authentification JWT.
- Tout ajout de fonctionnalités, utilisation de technologies autres que le REST (gRPC, websocket, AMQP...)

## Références
- [Tutoriel création d'une API de gestion des utilisateurs avec utilisation des JWT](https://github.com/MohammadBnei/gorm-user-auth)
- [Tutoriel de création d'un broadcaster et d'un room manager en golang](https://github.com/MohammadBnei/go-realtime-chat)