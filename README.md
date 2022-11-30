# Sistemi Distribuiti e Cloud Computing - A.A. 2021/22
## Progetto B2: Algoritmi di mutua esclusione distribuita in Go

Lo scopo del progetto è realizzare nel linguaggio di programmazione Go un’applicazione distribuita che
implementi tre algoritmi di mutua esclusione, in particolare:
- algoritmo token-asking,
- algoritmo di Lamport distribuito,
- algoritmo di Ricart-Agrawala

# Configurazione 
All'interno del file configurations.co, contenuto nella cartella /src/utilities, si hanno diversi parametri configurabili manualmente:
- MAXPEERS: numero di peer che si vuole partecipino al gruppo di mutua esclusione. Attualmente tale parametro è posto pari a 3, nel caso venga modificato è importante anche cambiare il numero di repliche per il container "peer" nel file docker-compose.yml, posto nella cartella /src. 
- TEST: da porre pari a true se si vogliono eseguire i test
- VERBOSE: da porre pari a true se si vogliono scrivere sui file di log le informazioni riguardanti i messaggi scambiati e le entrate in sezione critica dei vari peer.
È importante, per l'esecuzione dei test, che tale parametro sia posto pari a true.

## Esecuzione
Per lanciare l'applicazione eseguire il comando
```sh
sudo docker-compose up --build
```
Per potersi connettere ad un peer specifico, identificato da un intero ID, eseguire
``` sh
sudo docker exec -it src_peer_ID /bin/sh
``` 
(il nome del container può essere letto tramite il comando sudo docker ps)

Una volta essersi connessi al container, l'utente inizialmente dovrà scegliere da riga di comando lo username da dare al peer, e successivamente si aprirà un menu tramite il quale l'utente può decidere quale algoritmo eseguire e eventualmente far mandare un messaggio di richiesta al peer, come mostrato nel video seguente: 

https://user-images.githubusercontent.com/54946553/204897320-91158479-8269-4430-a124-03af3779bfe3.MP4

## Test
Per lanciare i test, porre pari a 'true' le variabili TEST e VERBOSE all'interno del file configurations.go. Una volta essersi connessi al container dei vari peer, l'utente potrà scegliere quale test effettuare tramite un menu:

https://user-images.githubusercontent.com/54946553/204896366-e5bc9bfa-269b-4cf3-a994-a6889b8dbc11.MP4





