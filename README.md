# NWS Conference – Installation & Démarrage Local

Ce projet est une web app de gestion de conférences (backend Go + frontend Vue 3).

--------------------------------
## Prérequis
```bash
- Node.js >= 16.x  
- npm >= 8.x ou yarn  
- Go >= 1.20  
- MAMP/XAMPP ou autre MySQL local
```

--------------------------------
## 1. Cloner le dépôt

```bash
git clone https://github.com/LCRpro/partiel-master1.git
cd partiel-master1
```

--------------------------------
## 2. Backend (Go)

a) Aller dans le dossier backend

```bash
cd backend
```

b) Installer les dépendances Go

```bash
go mod tidy
```
c) Configuration des variables d'environnement

Crée un fichier .env à la racine de backend/ avec le contenu fourni sur Classroom ou les votres directement :
```bash
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=
JWT_SECRET=changemeplease
DB_CONN=user:mdp@tcp(host:port)/nombdd
```


d) Lancer le serveur backend
```bash
go run main.go
```
Le backend écoute sur http://localhost:8080

--------------------------------
## 3. Frontend (Vue 3)

a) Aller dans le dossier frontend
```bash
cd ../frontend
```
b) Installer les dépendances
```bash
npm install

```
c) Lancer le serveur de développement
```bash
npm run dev

```
Le frontend écoute sur http://localhost:5173

--------------------------------
## 4. Accès
```bash
- Frontend : http://localhost:5173  
- Backend API : http://localhost:8080
```
--------------------------------
## 5. Tester

Lances tes tests backend :
```bash
go test -v ./controllers
```
--------------------------------
✨ Auteur

Liam Cariou

--------------------------------
