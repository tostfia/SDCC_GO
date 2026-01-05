package registry

import (
	"sync"					//pacchetto per la mutua esclusione sui dati condivisi
)


//Struttura ServiceInfo descrive un servizio registrato
type ServiceInfo struct{
	Name string //Nome logico del servizio (es. S1 S2 ...)
	Host string //Host su cui gira il servizio 
	Port int  // Porta su cui il servizio ascolta
	Weight int  //peso per il load balancing stateful
}

//Registry mantiene la lista dei servizi registrati
type Registry struct{
	mu sync.Mutex 			//Mutex per proteggere l'accesso concorrente alla mappa
	services map[string] ServiceInfo   //Mappa: nome servizio -> informazioni sul servizio
}

//NewRegistry crea un nuovo registry vuoto
func NewRegistry() *Registry{
	return &Registry{ 
		services: make(map[string]ServiceInfo), //inizializza la mappa nuova
	}
}


//metodi esportati per registrazione, deregistrazione e lookup

//Register aggiunge o aggiorna un servizio nel registry (argomento,*risposta)
func (r *Registry) Register(info ServiceInfo,ok *bool) error{
	r.mu.Lock()  //Blocca il mutex per evitare accessi concorrenti
	defer r.mu.Unlock() //Sblocco alla fine delle funzione


	r.services[info.Name]=info //Inserisco e aggiorno 
	*ok=true //se va bene 
	return nil
}


//Deregister rimuove un servizio alla mappa
func (r *Registry) Deregistry(info ServiceInfo,ok *bool) error{
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.services,info.Name)//Rimuovo
	*ok=true
	return nil
}

//Lookup restituisce la lista
//L'argomento è una struct vuota perchè non serve input
func (r *Registry) Lookup(_ struct{}, list *[]ServiceInfo) error{
	r.mu.Lock() 
	defer r.mu.Unlock()

	//Pulisce la lista prima di riempirla evitare confusione
	*list=(*list)[:0]

	//Aggiunge tutti i servizi nella slice di input
	for _,s:= range r.services{
		*list=append(*list,s)
	}
	return nil
}



