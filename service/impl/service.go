package impl

import "fmt"

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

func (s *Service) DoWork(req WorkRequest, output *string) error {
	s.RequestCount++

	if s.ByClient == nil {
		s.ByClient = make(map[string]int)
	}
	s.ByClient[req.ClientID]++

	*output = fmt.Sprintf(
		"Servizio %s → richiesta #%d (client %s: %d richieste) → payload='%s'",
		s.Name,
		s.RequestCount,
		req.ClientID,
		s.ByClient[req.ClientID],
		req.Payload,
	)

	return nil
}

// Usato dal registry per verificare se il servizio è vivo
func (s *Service) HealthCheck(_ struct{}, ok *bool) error {
	*ok = true
	return nil
}