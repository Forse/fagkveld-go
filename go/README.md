## Go-kode

Her finner du instruksjoner for oppsett og gjennomgang av kode.

### Oppsett

<ul>
	<li>Bruk Linux/Mac (jeg bruker WSL2 Ubuntu) og VSCode</li>
	<li><a href="https://go.dev/doc/install" target="_blank">Installer Go v1.20</a></li>
	<li>Installer `@recommended` extensions i VSCode</li>
	<li>
		Go extension i VSCode vil be deg installere noen tools (bl. a. language server), 
		<a href="https://github.com/golang/vscode-go/blob/master/docs/tools.md" target="_blank">gjør det evt manuelt</a>
	</li>
	<li>Github Copilot for assistert læring</li>
</ul>

```bash
wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
go version
source ~/.profile
```

Erstatt eventuelt `linux-amd64` med din arkitektur. Tilgjengelige archives finnes på [go.dev/dl/](https://go.dev/dl/)

### Hello, world!

```bash
go run cmd/hello/main.go
```

### Eksempler fra slides

Se [/go/internal/slides](/go/internal/slides)-mappa for en del av eksempel-koden som gås gjennom i slides.

### Øvelser - Exercism

Exercism er en plattform for å lære seg språk. 
Man kan løse små oppgaver gjennom VSCode ved hjelp av Exercism CLI,
eller bare gjøre det i editor på web.

* Logg inn på Exercism - jeg brukte Github-konto
* [Join Go track](https://exercism.org/tracks/go) 
* [Se instruksjoner for installasjon av CLI her](https://exercism.org/cli-walkthrough), eller følg nedenfor

```bash
wget https://github.com/exercism/cli/releases/download/v3.1.0/exercism-3.1.0-linux-x86_64.tar.gz
tar -xf exercism-3.1.0-linux-x86_64.tar.gz
mkdir -p ~/bin
mv exercism ~/bin
~/bin/exercism

# Check if ~/bin is in path
[[ ":$PATH:" == *":$HOME/bin:"* || ":$PATH:" == *":~/bin:"* ]] && echo "~/bin is in PATH" || echo "~/bin is not in PATH"

// Add to path
echo 'export PATH=~/bin:$PATH' >> ~/.bash_profile
source ~/.bash_profile
exercism
# Token can be found on the settings page: https://exercism.org/settings/api_cli
exercism configure --token=<insert-your-token>
exercism download --exercise=hello-world --track=go
cd ~/exercism/go/hello-world
code .
```

Nå kan du løse oppgaven. Se terminal-vindu for submission av løsning.
![Exercism og VSCode](/slides/public/exercism-helloworld.png)

### Øvelse - docker-ctx-analyze

Et CLI som gjør analyse for å finne ut hvor stor docker context blir i bygging av et container image.
Se [cmd/docker-ctx-analyze/main.go](/go/cmd/docker-ctx-analyze/main.go)-filen for instruksjoner.

```bash
go run cmd/docker-ctx-analyze/main.go
```

Hvis jeg får tid, skal jeg pushe en 'ferdig' implementasjon til branchen `completed/docker-ctx-analyze`

### Øvelse - todo-api

Et HTTP API for en todo-applikasjon.
Se [cmd/todo-api/main.go](/go/cmd/todo-api/main.go)-filen for instruksjoner.

```bash
go run cmd/todo-api/main.go
# Hvis du har Docker
docker build -t todo-api:latest .
```

### Kjør alle tester

```bash
go test ./...
```
