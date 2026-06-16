# FOR MENTOR — Forum web en Go

FOR MENTOR est un forum web développé en **Go**, avec une interface en **HTML/CSS pur** et une base de données **SQLite**.  
Le projet permet à des utilisateurs de s'inscrire, se connecter, publier des contenus, commenter, liker/disliker et filtrer les publications.

Le projet est également **containerisé avec Docker**, afin de pouvoir être lancé facilement sans installation manuelle de Go ou des dépendances sur la machine d'évaluation.

---

## Sommaire

- [Objectif du projet](#objectif-du-projet)
- [Fonctionnalités](#fonctionnalités)
- [Technologies utilisées](#technologies-utilisées)
- [Architecture du projet](#architecture-du-projet)
- [Base de données](#base-de-données)
- [Lancement avec Docker](#lancement-avec-docker)
- [Lancement local sans Docker](#lancement-local-sans-docker)
- [Routes principales](#routes-principales)
- [Sessions et cookies](#sessions-et-cookies)
- [Upload d'images](#upload-dimages)
- [Commandes utiles](#commandes-utiles)
- [Dépannage](#dépannage)
- [Pistes d'amélioration](#pistes-damélioration)

---

## Objectif du projet

L'objectif est de réaliser un forum complet côté serveur en respectant les contraintes suivantes :

- backend en **Go uniquement** ;
- aucune bibliothèque frontend type React, Vue, Angular ou Bootstrap ;
- base de données **SQLite** ;
- authentification avec mots de passe hashés ;
- gestion des sessions avec cookies ;
- posts, commentaires, catégories, likes et dislikes ;
- upload d'images ;
- lancement complet via Docker.

---

## Fonctionnalités

### Authentification

- Inscription avec pseudo, email et mot de passe.
- Vérification des champs côté serveur.
- Hashage des mots de passe avec `bcrypt`.
- Connexion via email et mot de passe.
- Création d'une session utilisateur après connexion.
- Déconnexion avec suppression du cookie de session.

### Publications

- Création de posts par les utilisateurs connectés.
- Titre obligatoire.
- Contenu obligatoire.
- Association d'une ou plusieurs catégories.
- Upload d'image optionnel.
- Affichage des posts sur la page d'accueil.
- Affichage de l'auteur, du titre, du contenu, de l'image et de la date.

### Commentaires

- Ajout de commentaires par les utilisateurs connectés.
- Affichage des commentaires sous les posts.
- Les visiteurs peuvent lire les commentaires.

### Likes / dislikes

- Like et dislike sur les posts.
- Un utilisateur ne peut pas liker et disliker le même post en même temps.
- Compteurs de likes et dislikes affichés sur les publications.

### Filtres

- Affichage de tous les posts.
- Filtrage par catégorie.
- Filtrage des posts créés par l'utilisateur connecté.
- Filtrage des posts likés par l'utilisateur connecté.

### Profil

- Page de profil accessible aux utilisateurs connectés.
- Affichage des informations du compte.
- Accès rapide à la création de post.
- Bouton de déconnexion.

---

## Technologies utilisées

| Technologie | Utilisation |
|---|---|
| Go | Serveur HTTP, routes, handlers, logique backend |
| SQLite | Base de données locale du forum |
| HTML | Templates des pages |
| CSS | Mise en page et design |
| JavaScript | Prévisualisation d'image sur la page de création |
| Docker | Containerisation et lancement simplifié |
| bcrypt | Hashage sécurisé des mots de passe |
| go-sqlite3 | Driver SQLite pour Go |

---

## Architecture du projet

```text
Projet-forum/
├── Dockerfile
├── compose.yaml
├── .dockerignore
├── go.mod
├── go.sum
│
├── back-end/
│   ├── cmd/
│   │   └── main.go
│   │
│   ├── database/
│   │   ├── database.go
│   │   └── _sql/
│   │       ├── DataBase.sql
│   │       └── forum.db
│   │
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── comment.go
│   │   │   ├── creation.go
│   │   │   ├── index.go
│   │   │   ├── liked_posts.go
│   │   │   ├── likes.go
│   │   │   ├── logout.go
│   │   │   ├── profil.go
│   │   │   ├── register.go
│   │   │   └── singIn.go
│   │   │
│   │   ├── models/
│   │   │   ├── category.go
│   │   │   ├── comment.go
│   │   │   ├── cookies.go
│   │   │   ├── home.go
│   │   │   ├── liked_posts.go
│   │   │   ├── likes.go
│   │   │   ├── post.go
│   │   │   └── user.go
│   │   │
│   │   └── render/
│   │       └── render.go
│   │
│   ├── pkg/
│   │   └── utils/
│   │       └── utils.go
│   │
│   └── server/
│       └── server.go
│
└── front-end/
    ├── static/
    │   ├── css/
    │   ├── icons/
    │   ├── js/
    │   └── uploads/
    │
    └── template/
        ├── HomePage.html
        ├── LikedPosts.html
        ├── ProfilPage.html
        ├── SettingsPage.html
        ├── SignIn.html
        ├── SignUp.html
        └── creationPage.html
```

---

## Base de données

Le projet utilise une base SQLite située ici :

```text
back-end/database/_sql/forum.db
```

Le schéma SQL de référence est situé ici :

```text
back-end/database/_sql/DataBase.sql
```

### Tables principales

| Table | Rôle |
|---|---|
| `users` | Stocke les comptes utilisateurs |
| `posts` | Stocke les publications |
| `comments` | Stocke les commentaires |
| `category` | Stocke les catégories disponibles |
| `post_categories` | Associe les posts à leurs catégories |
| `posts_like` | Stocke les likes des posts |
| `posts_dislike` | Stocke les dislikes des posts |
| `comments_like` | Stocke les likes des commentaires |
| `comments_dislike` | Stocke les dislikes des commentaires |
| `user_avatars` | Stocke les avatars utilisateurs sous forme de BLOB |

### Relations principales

- Un utilisateur peut créer plusieurs posts.
- Un utilisateur peut écrire plusieurs commentaires.
- Un post peut recevoir plusieurs commentaires.
- Un post peut être associé à plusieurs catégories.
- Un utilisateur peut liker ou disliker un post.
- Les suppressions utilisent des clés étrangères avec `ON DELETE CASCADE` lorsque nécessaire.

---

## Lancement avec Docker

### Prérequis

- Docker Desktop installé.
- Docker Desktop lancé.
- Le code du projet récupéré en local.

Aucune installation manuelle de Go, SQLite ou GCC n'est nécessaire pour lancer le projet avec Docker.

### Commande de lancement

Depuis la racine du projet :

```bash
docker compose up --build
```

Une fois le serveur lancé, ouvrir :

```text
http://localhost:8080
```

### Arrêter le projet Docker

```bash
docker compose down
```

### Fonctionnement Docker

Le `Dockerfile` compile l'application Go dans une première image, puis copie le binaire compilé dans une image finale plus légère.

Le fichier `compose.yaml` :

- construit l'image Docker ;
- lance le conteneur `projet-forum` ;
- expose le port `8080` ;
- monte la base SQLite locale dans le conteneur ;
- conserve les images uploadées dans un volume Docker.

---

## Lancement local sans Docker

### Prérequis

- Go `1.25.0` ou supérieur.
- Un compilateur compatible CGO pour `go-sqlite3`.
- SQLite CLI optionnel pour inspecter ou modifier la base.

### Installation des dépendances

Depuis la racine du projet :

```bash
go mod tidy
```

### Lancement du serveur

```bash
go run ./back-end/cmd
```

Puis ouvrir :

```text
http://localhost:8080
```

---

## Routes principales

| Route | Méthode | Description | Accès |
|---|---:|---|---|
| `/` | GET | Page d'accueil avec les posts et filtres | Public |
| `/register` | GET | Affichage du formulaire d'inscription | Public |
| `/register` | POST | Création d'un compte utilisateur | Public |
| `/login` | GET | Affichage du formulaire de connexion | Public |
| `/login` | POST | Connexion utilisateur | Public |
| `/logout` | POST | Déconnexion utilisateur | Connecté |
| `/post/create` | GET | Affichage du formulaire de création de post | Connecté |
| `/post/create` | POST | Création d'un nouveau post | Connecté |
| `/comment/create` | POST | Ajout d'un commentaire à un post | Connecté |
| `/post/like` | POST | Like d'un post | Connecté |
| `/post/dislike` | POST | Dislike d'un post | Connecté |
| `/liked-posts` | GET | Page des posts likés | Connecté |
| `/profile` | GET | Page de profil utilisateur | Connecté |

---

## Sessions et cookies

Le forum utilise un système de sessions côté serveur basé sur un cookie nommé :

```text
session_id
```

Lorsqu'un utilisateur se connecte :

1. le serveur vérifie l'email et le mot de passe ;
2. un identifiant de session aléatoire est généré ;
3. la session est stockée en mémoire côté serveur ;
4. un cookie `session_id` est envoyé au navigateur ;
5. les pages protégées peuvent retrouver l'utilisateur connecté grâce à ce cookie.

Les cookies sont configurés avec :

- `HttpOnly` pour empêcher leur lecture par JavaScript ;
- `SameSite=Lax` ;
- une expiration de 24 heures.

> Les sessions sont actuellement stockées en mémoire. Si le serveur est redémarré, les utilisateurs doivent se reconnecter.

---

## Upload d'images

La création de post permet d'envoyer une image optionnelle.

Formats acceptés :

- JPEG ;
- PNG ;
- GIF.

Taille maximale :

```text
20 Mo
```

Les images sont enregistrées dans :

```text
front-end/static/uploads/
```

Le lien de l'image est ensuite stocké dans la colonne `image_link` de la table `posts`.

---

## Commandes utiles

### Lancer avec Docker

```bash
docker compose up --build
```

### Arrêter Docker

```bash
docker compose down
```

### Voir les conteneurs actifs

```bash
docker ps
```

### Voir les logs Docker

```bash
docker compose logs -f
```

### Lancer en local

```bash
go run ./back-end/cmd
```

### Vider les données de la base SQLite locale

Depuis PowerShell, à la racine du projet :

```powershell
@"
PRAGMA foreign_keys = OFF;

DELETE FROM comments_like;
DELETE FROM comments_dislike;
DELETE FROM posts_like;
DELETE FROM posts_dislike;
DELETE FROM post_categories;
DELETE FROM comments;
DELETE FROM posts;
DELETE FROM category;
DELETE FROM users;

PRAGMA foreign_keys = ON;
VACUUM;
"@ | sqlite3 .\back-end\database\_sql\forum.db
```

### Vérifier les tables SQLite

```bash
sqlite3 back-end/database/_sql/forum.db ".tables"
```

---

## Dépannage

### Le port 8080 est déjà utilisé

Erreur possible :

```text
listen tcp :8080: bind: Only one usage of each socket address is normally permitted
```

Solution : arrêter l'ancien serveur ou le conteneur Docker.

```bash
docker compose down
```

Ou chercher le processus sur Windows :

```powershell
netstat -ano | findstr :8080
```

Puis arrêter le processus avec son PID :

```powershell
taskkill /PID <PID> /F
```

### Docker ne répond pas

Erreur possible :

```text
failed to connect to the docker API
```

Solution : ouvrir Docker Desktop et attendre que le moteur soit lancé, puis relancer :

```bash
docker info
```

### Les données visibles en local ne sont pas visibles dans Docker

Le projet utilise un montage dans `compose.yaml` :

```yaml
volumes:
  - ./back-end/database/_sql:/app/back-end/database/_sql
```

Ce montage permet à Docker d'utiliser la même base SQLite que le projet local.

Si les données ne correspondent pas, vérifier :

- que le serveur local est bien arrêté ;
- que Docker utilise le bon dossier ;
- que la base consultée est bien `back-end/database/_sql/forum.db`.

### Erreur `no such table`

Cela signifie généralement que la base SQLite existe mais que le schéma n'a pas été créé.

Vérifier que le fichier suivant existe :

```text
back-end/database/_sql/forum.db
```

Et que le schéma de référence existe :

```text
back-end/database/_sql/DataBase.sql
```

---

## Pistes d'amélioration

- Remplacer les sessions en mémoire par une table `sessions` persistante.
- Finaliser la page paramètres.
- Finaliser la modification du profil.
- Ajouter la suppression et modification des posts.
- Ajouter la suppression et modification des commentaires.
- Ajouter les likes/dislikes sur les commentaires.
- Ajouter une initialisation automatique de la base depuis `DataBase.sql` si `forum.db` n'existe pas.
- Ajouter des tests unitaires sur les modèles et handlers.
- Nettoyer les routes non utilisées et harmoniser les noms de templates.

---

## Auteurs

Projet réalisé dans le cadre d'un projet de forum web en B1.

Membres du groupe :

- Decung Tessa
- Jeannot Louis
- Royer Nicolas
- Lamarche Nathan

