# Projet-Forum

Un forum web développé en Go, dédié à la communauté du jeu **Teamfight Tactics (TFT)**, permettant aux utilisateurs de s’inscrire, se connecter, créer des fils de discussion (feeds), poster des messages, et gérer leur profil.

## Fonctionnalités principales

- **Inscription et connexion** des utilisateurs (avec gestion de session par JWT)
- **Création de fils de discussion** (feeds) avec tags, état (ouvert/fermé)
- **Affichage de tous les fils**, des fils par utilisateur, ou par tag
- **Ajout de posts** dans un fil de discussion
- **Consultation des posts** d’un fil
- **Mise à jour de l’état** d’un fil (ouvert/fermé)
- **Gestion du profil utilisateur** (bio, photo, etc.)
- **Déconnexion**

## Stack technique

- **Backend** : Go 1.22
- **Base de données** : MySQL 
- **Gestion des variables d’environnement** : godotenv
- **Authentification** : JWT 

## Installation

1. **Cloner le dépôt**  
   ```bash
   git clone https://github.com/LilianLPV/Projet-Forum.git
   cd Projet-Forum
   ```

2. **Configurer la base de données**  
   - Créer une base MySQL (gestion via **phpMyAdmin** recommandée).
   - Pour avoir un serveur MySQL/phpMyAdmin en local, installer **XAMPP** ou **WAMP**.
   - Adapter les variables d’environnement dans un fichier `.env` (voir `config/` pour les variables attendues).

4. **Lancer les migrations**  
   - Le dossier `migration/` contient les scripts SQL à exécuter dans phpMyAdmin pour créer les tables de la base de données.

5. **Démarrer le serveur**  
   ```bash
   go run main.go
   ```
   Le serveur sera accessible sur [http://localhost:8080](http://localhost:8080).

## Structure du projet

- `main.go` : point d’entrée, configuration du serveur
- `controllers/` : gestion des routes et logique web
- `services/` : logique métier
- `models/` : structures de données
- `repositories/` : accès à la base de données
- `template/` : templates HTML
- `static/` : fichiers statiques (CSS, JS, images)
- `auth/` : gestion de l’authentification
- `config/` : configuration et initialisation

## Exemple de routes

- `/signup`, `/login`, `/logout`
- `/feeds/create`, `/feeds/all`, `/feeds/{id}`
- `/posts/create`, `/posts/feed/{id}`
- `/profile`

## Auteurs

- Lilian LE PIVER
- Bastien JOUFFRE