package main 

import(
	"fmt"
	"log"
	"net/rpc"
	"time"
	"SDCC_GO/registry"
)


//Cache dei servizi
var cachedServices []registry.ServiceInfo
var cacheTimestamp time.Time

//TTL della cache dinamica in base al numero di servizi
func computeTTL(serviceCount int) time.Duration{
	switch  {
	case serviceCount <=1:
		return 2*time.Second
	case serviceCount <=5:
		return 5*time.Second
	default:
		return 10*time.Second
		
	}
}

func invalidateCache() {
	cachedServices = nil
	cacheTimestamp = time.Time{}
}


//Funzione che restituisce i servizi usando la cache
func getServicesWithCache() []registry.ServiceInfo{
	//Se la cache è ancora valida, si usa
	if len(cachedServices)>0{
		ttl:=computeTTL(len(cachedServices))
		if time.Since(cacheTimestamp)<ttl{
			fmt.Printf("Uso della cache (TTL dinamico: %v)\n", ttl)
			return cachedServices
		}
	}
	//Altrimenti faccio lookup dal registry
	regClient,err:=rpc.Dial("tcp","registry:9000")
	if err!=nil{
		log.Println("Registry non raggiungibile, uso cache se disponibile")
        return cachedServices
    }
    defer regClient.Close()

	var services []registry.ServiceInfo
	err=regClient.Call("Registry.Lookup", struct{}{}, &services)
    if err != nil {
        log.Println("Errore nella lookup, uso cache se disponibile")
        return cachedServices
    }

    // Aggiorna la cache
    cachedServices = services
    cacheTimestamp = time.Now()
	ttl:=computeTTL(len(services))
    fmt.Printf("Cache aggiornata (TTL dinamico: %v): %v\n", ttl, services)
    return services
}

