package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "os"
    "strconv"

    impl "SDCC_GO/service/impl"
)

func main() {

    //Parametri da terminale
    if len(os.Args) <3 {
        log.Fatalf("Uso:go run ./service  <Porta> <Peso>")
    }


    // Nome del servizio preso dalla variabile d'ambiente
    name := os.Getenv("SERVICE_NAME")
    if name == "" {
        log.Fatalf("Variabile SERVICE_NAME non impostata")
    }

    // Porta come string convertito a int
    portStr := os.Args[1]
    port, err := strconv.Atoi(portStr)
    if err != nil {
        log.Fatalf("Porta non valida: %v", err)
    }

    // Peso come string convertito a int
    weightStr := os.Args[2]
    weight, err := strconv.Atoi(weightStr)
    if err != nil {
        log.Fatalf("Peso non valido: %v", err)
    }

    srv := &impl.Service{Name:name,Host:os.Getenv("SERVICE_HOST"), Port:port, Weight: weight, }

    // Registra il servizio RPC
    err = rpc.Register(srv)
	if err != nil {
		log.Fatalf("Errore registrazione RPC servizio: %v", err)
	}


    // Avvia listener
    address := fmt.Sprintf(":%d", port)
    listener, err := net.Listen("tcp", address)
    if err != nil {
        log.Fatalf("Errore listen: %v", err)
    }

    fmt.Printf("Servizio %s in ascolto sulla porta %d...\n", name, port)

    // Self-registration
    var ok bool
    err = srv.Register(struct{}{}, &ok)
    if err != nil {
        log.Println("Errore registrazione:", err)
    }

    // Deregistrazione automatica alla chiusura
    defer srv.Deregister(struct{}{}, &ok)

    // Loop di accettazione
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Errore accept:", err)
            continue
        }
        go rpc.ServeConn(conn)
    }
}