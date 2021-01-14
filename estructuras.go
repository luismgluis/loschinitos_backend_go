package main

import (
	"net/http"

	"github.com/mholt/binding"
)

type Cliente struct {
	UID    string   `json:"uid,omitempty"`
	Name   string   `json:"name,omitempty"`
	Age    int      `json:"age,omitempty"`
	Avatar string   `json:"avatar,omitempty"`
	DType  []string `json:"dgraph.type,omitempty"`
}

// Then provide a field mapping (pointer receiver is vital)
func (cf *Cliente) FieldMap(req *http.Request) binding.FieldMap {
	cf.DType = []string{"Cliente"}
	return binding.FieldMap{
		&cf.UID: binding.Field{
			Form:     "id",
			Required: true,
		},
		&cf.Name:   "name",
		&cf.Age:    "age",
		&cf.Avatar: "avatar",
	}
}

//DType  []string `json:"dgraph.type,omitempty"`

type Clientes struct {
	Clientes []Cliente `json:"clientes"`
}

type Producto struct {
	UID   string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Price int      `json:"age,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Transaccion struct {
	UID        string   `json:"id,omitempty"`
	BuyerID    string   `json:"buyerid,omitempty"`
	IP         string   `json:"ip,omitempty"`
	Device     string   `json:"device,omitempty"`
	ProductIDS []string `json:"produts,omitempty"`
	DType      []string `json:"dgraph.type,omitempty"`
}
