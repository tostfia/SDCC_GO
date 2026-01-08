package service

import "fmt" 
type Service struct{
	Name string
	Host string 
	Port int
	Weight int

}

func (s *Service) DoWork (input string, output *string) error{
	*output= fmt.Sprintf("Risposta da %s: ho ricevuto '%s'",s.Name,input)
	return nil
}