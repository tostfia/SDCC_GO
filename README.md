# **Sistema Distribuito con Service Registry e Load Balancing lato Client**

**Corso:** Sistemi Distribuiti e Cloud Computing  
**Anno Accademico:** 2025/2026  

Questo progetto implementa un sistema distribuito che integra  
service discovery, load balancing lato client, caching dinamica  
e orchestrazione tramite Docker Compose.

---

## **Architettura del Sistema**

Il sistema è composto dai seguenti componenti:

- *Service Registry*  
- *Servizi RPC (`service1`, `service2`, …)*  
- *Client con load balancing lato client*  
- *Caching dinamica dei servizi*  
- *Registrazione e deregistrazione automatica*  
- *Orchestrazione tramite Docker Compose*  

I client scoprono dinamicamente i servizi disponibili tramite  
il registry, selezionano un servizio utilizzando diversi algoritmi  
di load balancing e inviano richieste RPC.  
I servizi possono mantenere stato interno.

---

## **Service Registry**

Il service registry mantiene la lista dei servizi attivi e  
supporta le seguenti operazioni:

- *Register: i servizi si registrano all’avvio fornendo nome, host, porta e peso*  
- *Deregister: i servizi si deregistrano automaticamente alla chiusura*  
- *Lookup: i client richiedono la lista aggiornata dei servizi disponibili*  

Il registry viene utilizzato esclusivamente per la scoperta  
dei servizi. Il bilanciamento del carico è effettuato  
interamente lato client.

---

## **Servizi RPC**

Ogni servizio espone un metodo RPC che elabora le richieste  
dei client e restituisce una risposta.

*Sottotitolo – Caratteristiche principali:*

- *Registrazione automatica all’avvio*  
- *Deregistrazione automatica alla terminazione*  
- *Possibilità di mantenere stato interno (ad esempio, conteggio delle richieste per client)*  
- *Comunicazione RPC su TCP*  

---

## **Client con Load Balancing lato Client**

Ogni client esegue le seguenti operazioni:

- *Effettua una lookup iniziale presso il service registry*  
- *Memorizza la lista dei servizi in una cache locale*  
- *Seleziona dinamicamente un servizio tramite un algoritmo di load balancing*  
- *Invia una richiesta RPC al servizio selezionato*  

In caso di errore RPC:

- *La cache viene invalidata*  
- *Viene effettuata una nuova lookup*  
- *Viene selezionato un nuovo servizio*  

---

## **Algoritmi di Load Balancing**

Il progetto implementa sia algoritmi stateless sia stateful.

*Sottotitolo – Stateless:*

- *Random: selezione casuale del servizio*  
- *Round Robin: selezione ciclica dei servizi*  

*Sottotitolo – Stateful:*

- *Weighted: selezione basata su un peso assegnato dal registry,  
  utile per distribuire il carico in base alle risorse disponibili*  

---

## **Caching Dinamica**

I client utilizzano una cache locale per ridurre il numero  
di lookup verso il registry.

- *Il TTL della cache è calcolato dinamicamente in base al numero di servizi attivi*  
- *La cache viene invalidata automaticamente in caso di errore RPC*  

---

## **Supporto Multi-Client**

Il sistema supporta più client in esecuzione parallela,  
ciascuno con un proprio algoritmo di load balancing.

Questo consente di osservare:

- *Concorrenza reale*  
- *Distribuzione del carico tra i servizi*  
- *Comportamento stateful dei servizi*  
- *Caching indipendente per ciascun client*  

---

## **Esecuzione con Docker Compose**

L’intero sistema è orchestrato tramite Docker Compose.

*Sottotitolo – Avvio del sistema:*

Lo script esegue automaticamente:

- *Ricostruzione delle immagini Docker*
- *Avvio di registry, servizi e client*
- *Visualizzazione dei log dei client*

```bash
./run.sh
