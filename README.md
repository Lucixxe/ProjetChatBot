### Création de la base de données ###

Simplement tapez dans le terminal
```sh
	sqlite chatbot.db < schema.sql
```

### Lancement du serveur ###
Dans un terminal :
```
mkdir ollama
cd ollama
curl -L https://ollama.com/download/ollama-linux-amd64 -o ollama
chmod u+x ollama
./ollama serve
```
Dans un autre terminal, dans le répertoire ollama créé précédemment :
```
./ollama run llama3
```

### Compilation ###
Pour compiler le projet tapez :
```sh
	go get
	go build
```

Pour lancer l'application il suffit de lancer l'exécutable chatbot généré par la compilation avec la commande suivante :
```
go run chatbot
```

### Lancement du site ###
Le site démarre sur le port 3333. Pour le lancer, il vous faut ouvrir un navigateur puis taper dans la barre de recherche :
```
localhost:3333
```
Vous serez alors redirigé vers la page de connexion.
