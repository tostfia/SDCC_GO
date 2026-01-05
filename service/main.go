package main

import(
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"SDCC_GO/registry"
)


//Service è il servizio applicativo 
type Service struct {
	Name string
}

//DoWork costruisce una risposta usando il nome del servizio e l'input ricevuto
func (s *Service) DoWork(input string, output *string) error {
	*output= fmt.Sprintf("Risposta da %s: ho ricevuto '%s'", s.Name, input)
	return nil
}

func main(){
	//Nome del servizo da argomento oppure default S1
	name:= "S1"
	if len(os.Args)>1{
		name=os.Args[1]
	}
	//crea istanza del servizio
	srv:=&Service{Name:name}
	

	//Espone il servizio RPC
	err:= rpc.Register(srv)
	if err!=nil{
		log.Fatalf("Errore nella registrazione del servizio RPC:%v",err)
	}

	//Ascolta su una porta (es.8001)
	port:= ":8001"
	//Crea un listener TCP sulla porta specificata
	listener, err:= net.Listen("tcp",port)
	if err!=nil{
		log.Fatalf("Errore nel listen sulla porta %s:%v", port, err)
	}
	fmt.Printf("Servizo %s in ascolto sulla porta %s ... \n", name, port)

	//Connessione al registry 
	regClient, err:= rpc.Dial("tcp","localhost:9000")
	if err!=nil{
		log.Fatalf("Errore nella connessione al registry:%v",err)

	}
	defer regClient.Close()

	info:=registry.ServiceInfo{
		Name: name,
		Host: "localhost",
		Port: 8001,
		Weight: 1,
	}

	var ok bool
	err=regClient.Call("Registry.Register",info,&ok)
	if err!=nil{
		log.Fatalf("Errore nella registrazione del servizio: %v",err)
	}
	fmt.Println("Servizio registrato correttamente presso il registry")

	//Ciclo infinito per le connessioni
	for{
		conn,err:= listener.Accept()
		if err!=nil{
			log.Printf("Errore nell'accettare connessione:%v",err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}