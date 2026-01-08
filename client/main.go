package main

import(
	"fmt"
	"log"
	"net/rpc"
	"SDCC_GO/registry"
	"os"
	
	
)



func main(){
	//Scelta dell'algoritmo da terminale 
	algorithm:="random"//default
	if len(os.Args)>1{
		algorithm=os.Args[1]
	}
	fmt.Println("Algoritmo scelto:",algorithm)

	//connessione al registry
	regClient,err:=rpc.Dial("tcp","localhost:9000")
	if err!=nil{
		log.Fatalf("Errore nella connessione al servizio:%v",err)
	}
	//Chiude la connessione quando  la funzione main termina
	defer regClient.Close()
	
	//Caching, in caso di miss facciamo lookup
	services:=getServicesWithCache()

	//Scelgo un servizio con LB
	var chosen registry.ServiceInfo

	switch algorithm {
	case "random":
		chosen = RandomLB(services)

	case "rr":
		chosen = RoundRobinLB(services)

	case "weighted":
		chosen = WeightedLB(services)

	default:
		fmt.Println("Algoritmo non riconosciuto, uso random")
		chosen = RandomLB(services)
	}
	fmt.Println("Servizio scelto:", chosen.Name)

	//Connessione al servizio scelto
	addr := fmt.Sprintf("%s:%d", chosen.Host, chosen.Port)
	srvClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Errore nella connessione al servizio scelto: %v", err)
	}
	defer srvClient.Close()

	var reply string
	err = srvClient.Call("Service.DoWork", "ciao dal client", &reply)
	if err != nil {
		log.Fatalf("Errore nella chiamata RPC: %v", err)
	}

	fmt.Println("Risposta:", reply)


}