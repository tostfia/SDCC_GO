package main

import(
	"fmt"
	"log"
	"net/rpc"
	"SDCC_GO/registry"
)


func main(){
	//connessione al registry
	regClient,err:=rpc.Dial("tcp","localhost:9000")
	if err!=nil{
		log.Fatalf("Errore nella connessione al servizio:%v",err)
	}
	//Chiude la connessione quando  la funzione main termina
	defer regClient.Close()
	
	//Lookup dei servizi
	var services []registry.ServiceInfo
	err=regClient.Call("Registry.Lookup",struct{}{},&services)
	if err!=nil{
		log.Fatalf("Errore nella lookup:%v",err)
	}
	fmt.Println("Servizi disponibili:", services)

	//Scelta un servizio con LB
	chosen:= RandomLB(services)
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