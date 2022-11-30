# Sistemi Distribuiti e Cloud Computing - A.A. 2021/22
## Progetto B2: Algoritmi di mutua esclusione distribuita in Go

Lo scopo del progetto è realizzare nel linguaggio di programmazione Go un’applicazione distribuita che
implementi tre algoritmi di mutua esclusione, in particolare:
- algoritmo token-asking,
- algoritmo di Lamport distribuito,
- algoritmo di Ricart-Agrawala

# Configurazione iniziale
All'interno del file configurations.co, contenuto nella cartella /src/utilities, si hanno diversi parametri configurabili manualmente:
- MAXPEERS: numero di peer che si vuole partecipino al gruppo di mutua esclusione
- TEST: da porre pari a true se si vogliono eseguire i test
- VERBOSE: da porre pari a true se si vogliono scrivere sui file di log le informazioni riguardanti i messaggi scambiati e le entrate in sezione critica dei vari peer.
È importante, per l'esecuzione dei test, che anche il parametro VERBOSE sia posto pari a true.



fare in un terminale
``` sh
sudo docker-compose up --build
```

per collegarsi al nodo register, in un altro terminale fare
``` sh
sudo docker exec -it src_register_node_1 /bin/sh
```
dove il nome del nodo lo prendo leggendo da
``` sh
sudo docker ps
```
in questo modo il register si mette in attesa.


poi fare
``` sh
sudo docker exec -it src_peer_1 /bin/sh
``` 

per eliminare i file nei volumi fare da 
``` 
sudo docker system prune
``` 

per rimuovere permessi
``` 
sudo chmod -rwx directoryname
``` 
con il "-" rimuovo i permessi, con il "+" li aggiungo (fare prima - poi +)

il peer è il nodo su cui viene eseguito il processo! processo sono i task, il codice
