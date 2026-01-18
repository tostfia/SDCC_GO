package impl

import (
    "fmt"
    "os"
    "strconv"
)


// Struttura della richiesta RPC
type WorkRequest struct {
	ClientID string
	Payload  string
}

type Service struct {
	Name   string
	Host   string
	Port   int
	Weight int

	RequestCount int
	ByClient     map[string]int
}



//Stateless
func (s *Service) Echo(payload string, reply *string) error {
    *reply = "Echo: " + payload
    return nil
}


//Statefull
func (s *Service) DoWork(req WorkRequest, output *string) error {
    // Percorso del file condiviso tra repliche
    path := "/app/state/counter.txt"

    // Se il file non esiste, crealo con valore 0
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.WriteFile(path, []byte("0"), 0644)
    }

    // Leggi il contatore globale
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("errore lettura stato: %v", err)
    }

    count, err := strconv.Atoi(string(data))
    if err != nil {
        count = 0
    }

    // Aggiorna il contatore
    count++

    // Scrivi il nuovo valore
    err = os.WriteFile(path, []byte(strconv.Itoa(count)), 0644)
    if err != nil {
        return fmt.Errorf("errore scrittura stato: %v", err)
    }

    // Risposta
    *output = fmt.Sprintf(
        "Servizio %s → richiesta #%d (client %s) → payload='%s'",
        s.Name,
        count,
        req.ClientID,
        req.Payload,
    )

    return nil
}

// Usato dal registry per verificare se il servizio è vivo
func (s *Service) HealthCheck(_ struct{}, ok *bool) error {
	*ok = true
	return nil
}