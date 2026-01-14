package impl

import(
	"fmt"
	"log"
	"net/rpc"
	"SDCC_GO/registry"
)


//Metodo RPC: registra il servizio nel registry
func (s *Service) Register(_ struct{}, ok *bool) error{
	regClient, err :=rpc.Dial("tcp", "registry:9000")
	if err!=nil{
		return fmt.Errorf("errore connessione registry:%v",err)
	}
	defer regClient.Close()

	info := registry.ServiceInfo{
		Name: s.Name,
		Host: s.Host,
		Port: s.Port,
		Weight: s.Weight,
	}
	err = regClient.Call("Registry.Register", info, ok)
    if err != nil {
        return fmt.Errorf("errore registrazione: %v", err)
    }

    log.Println("Servizio registrato correttamente")
    return nil
}

// Metodo RPC: deregistra il servizio dal registry
func (s *Service) Deregister(_ struct{}, ok *bool) error {
    regClient, err := rpc.Dial("tcp", "registry:9000")
    if err != nil {
        return fmt.Errorf("errore connessione registry: %v", err)
    }
    defer regClient.Close()

    info := registry.ServiceInfo{
        Name: s.Name,
    }

    err = regClient.Call("Registry.Deregister", info, ok)
    if err != nil {
        return fmt.Errorf("errore deregistrazione: %v", err)
    }

    log.Println("Servizio deregistrato correttamente")
    return nil
}

