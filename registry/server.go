package registry

import ( 
	"fmt"
	"log"
	"net"		//per creare listener di rete (TCP)
	"net/rpc"  //pacchetto RPC 
)

//Avvia il server RPC del registry
func StartRegistry(){
	reg:=NewRegistry()

	err:=rpc.Register(reg)
	if err!=nil{
		log.Fatalf("Errore nella registrazione RPC:%v",err)
	}

	listener,err := net.Listen("tcp",":9000")
	if err!=nil{
		log.Fatalf("Errore nel listen sulla porta 9000:%v",err)
	}

	fmt.Println("Service Registry in ascolto sulla porta 9000...")


	//Ciclo infinito: accetta connessioni in arrivo.
	for{
		conn,err:=listener.Accept()//Attende una nuova connessione
		if err!=nil{
			log.Printf("Errore nell'accettare connessione:%v",err)
			continue
		}
		//Connessione in una goroutine separata
		go rpc.ServeConn(conn)
	}
}