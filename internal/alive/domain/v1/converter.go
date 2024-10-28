package v1

import "github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/response"

func (a *Alive) FromDomainToApi() *response.AliveResponse {
	if a == nil {
		return nil
	}

	return &response.AliveResponse{
		ServerName: a.ServerName,
		Hostname:   a.Hostname,
		Port:       a.Port,
		Message:    a.Message,
	}
}

func (pa *PublicAlive) FromDomainToApi() *response.PublicAliveResponse {
	if pa == nil {
		return nil
	}

	return &response.PublicAliveResponse{
		Name:    pa.Name,
		Message: pa.Message,
	}
}
