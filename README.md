# NWS Conference – Installation & Démarrage via Docker

Ce projet est une web app de gestion de conférences (backend Go + frontend Vue 3), entièrement conteneurisée avec Docker.

--------------------------------
## Prérequis

- Docker
- Docker Compose

--------------------------------
## 1. Cloner le dépôt

```bash
git clone https://github.com/LCRpro/partiel-master1.git
cd partiel-master1
```

--------------------------------
## 2. Créer un fichier `.env`

Crée un fichier `.env` à la racine du projet avec les variables suivantes :

```dotenv
voir classroom
```



--------------------------------
## 3. Démarrer les services

Depuis la racine du projet, lancez :

```bash
docker compose --env-file .env up --build
```

Cela va :
- Lancer MySQL sur le port 3306
- Lancer le backend Go sur le port 8080
- Lancer le frontend Vue 3 (Vite) sur le port 5173

--------------------------------
## 4. Accès

- Frontend : http://localhost:5173  
- Backend API : http://localhost:8080

--------------------------------
## 5. Lancer les tests backend

```bash
docker exec -it <nom_du_conteneur_backend> go test -v ./controllers
```

Vous pouvez obtenir le nom exact avec `docker ps`.

--------------------------------
✨ Auteur

Liam Cariou
