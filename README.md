Pour coder la partie serveur de notre site, nous utilisons la version 1.22.0 de Golang.


Pour le lancer sur un ordinateur, vous pouvez suivre les indications suivantes:

### Création de la base de données ###

Tapez dans un terminal la commande suivante:
```sh
sqlite chatbot.db < schema.sql
```

### Lancement du serveur ###
Entrez ensuite les suivantes:
```
mkdir ollama
cd ollama
curl -L https://ollama.com/download/ollama-linux-amd64 -o ollama
chmod u+x ollama
```
et enfin:
```
./ollama serve
```
Puis, dans un autre terminal, toujours dans le répertoire ollama créé précédemment:
```
./ollama run llama3
```

Attendez que tout se charge avant de passer aux étapes suivantes.

### Compilation ###
Pour compiler le projet tapez:
```sh
go get
go build
```

Pour lancer l'application il suffit de lancer l'exécutable chatbot généré par la compilation avec la commande suivante:
```
go run chatbot
```

### Lancement du site ###
Le site démarre sur le port 3333. Pour le lancer, il vous faut ouvrir un navigateur puis taper dans la barre de recherche:
```
localhost:3333
```
Vous serez alors redirigé vers la page de connexion et pourrez enfin utiliser l'application.
