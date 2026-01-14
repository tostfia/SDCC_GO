package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"

	"SDCC_GO/registry"
	"SDCC_GO/service/impl"
)

func main() {
	algorithm := "random"
	if len(os.Args) > 1 {
		algorithm = os.Args[1]
	}

	services := getServicesWithCache()
	if len(services) == 0 {
		log.Fatalf("Nessun servizio disponibile")
	}

	var chosen registry.ServiceInfo
	switch algorithm {
	case "random":
		chosen = RandomLB(services)
	case "rr":
		chosen = RoundRobinLB(services)
	case "weighted":
		chosen = WeightedLB(services)
	default:
		chosen = RandomLB(services)
	}

	addr := fmt.Sprintf("%s:%d", chosen.Host, chosen.Port)
	srvClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Printf("Errore connessione → fallback")
		invalidateCache()
		services = getServicesWithCache()
		chosen = RandomLB(services)
		addr = fmt.Sprintf("%s:%d", chosen.Host, chosen.Port)
		srvClient, err = rpc.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("Errore anche nel fallback: %v", err)
		}
	}
	defer srvClient.Close()

	for i := 1; i <= 5; i++ {
		req := impl.WorkRequest{
			ClientID: "C1",
			Payload:  fmt.Sprintf("richiesta #%d", i),
		}

		var reply string
		err = srvClient.Call("Service.DoWork", req, &reply)
		if err != nil {
			log.Printf("Errore RPC → cambio servizio")
			invalidateCache()
			services = getServicesWithCache()
			chosen = RandomLB(services)
			addr = fmt.Sprintf("%s:%d", chosen.Host, chosen.Port)
			srvClient, _ = rpc.Dial("tcp", addr)
			continue
		}

		fmt.Println("Risposta:", reply)
		time.Sleep(1 * time.Second)
	}
}