package models

import "github.com/docker/docker/api/types"

type CtrMgt struct {
	Id     string       `json:"id,omitempty"`
	Name   string       `json:"name,omitempty"`
	Image  string       `json:"image,omitempty"`
	State  string       `json:"state,omitempty"`
	Status string       `json:"status,omitempty"`
	Ports  []types.Port `json:"ports,omitempty"`
}
