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
    if len(os.Args) <4 {
        log.Fatalf("Uso:go run ./service <Nome> <Porta> <Peso>")
    }


    name:= os.Args[1]
    // Porta come string convertito a int
    portStr := os.Args[2]
    port, err := strconv.Atoi(portStr)
    if err != nil {
        log.Fatalf("Porta non valida: %v", err)
    }

    // Peso come string convertito a int
    weightStr := os.Args[3]
    weight, err := strconv.Atoi(weightStr)
    if err != nil {
        log.Fatalf("Peso non valido: %v", err)
    }


    srv := &impl.Service{Name: name, Host: "localhost", Port: port, Weight: weight,}

    // Registra il servizio RPC
    rpc.Register(srv)

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