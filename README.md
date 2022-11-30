# Sistemi Distribuiti e Cloud Computing - A.A. 2021/22
## Progetto B2: Algoritmi di mutua esclusione distribuita in Go

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
