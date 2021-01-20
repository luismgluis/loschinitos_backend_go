package main

import (
	"net/http"

	"github.com/mholt/binding"
)

type Cliente struct {
	UID           string        `json:"uid,omitempty"`
	UIDOLD        string        `json:"UIDOLD,omitempty`
	Name          string        `json:"name,omitempty"`
	Age           int           `json:"age,omitempty"`
	Avatar        string        `json:"avatar,omitempty"`
	Date          int           `json:"date,omitempty"`
	DType         []string      `json:"dgraph.type,omitempty"`
	Transacciones []Transaccion `json:"transacciones,omitempty"`
}

// Then provide a field mapping (pointer receiver is vital)
func (cf *Cliente) FieldMap(req *http.Request) binding.FieldMap {
	cf.DType = []string{"Cliente"}
	return binding.FieldMap{
		&cf.UID: binding.Field{
			Form:     "id",
			Required: true,
		},
		&cf.UIDOLD: "uidold",
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
	UID    string   `json:"uid,omitempty"`
	PRODID string   `json:"PRODID,omitempty"`
	Name   string   `json:"name,omitempty"`
	Price  int      `json:"price,omitempty"`
	Date   int      `json:"date,omitempty"`
	DType  []string `json:"dgraph.type,omitempty"`
}
type Productos struct {
	Productos []Producto `json:"productos,omitempty"`
}

type Transaccion struct {
	UID        string    `json:"uid,omitempty"`
	TRANSID    string    `json:"TRANSID,omitempty"`
	BuyerID    string    `json:"buyerid,omitempty"`
	Buyer      Cliente   `json:"buyer,omitempty"`
	IP         string    `json:"ip,omitempty"`
	Device     string    `json:"device,omitempty"`
	ProductIDS []string  `json:"produtids,omitempty"`
	Products   Productos `json:"products,omitempty"`
	Date       int       `json:"date,omitempty"`
	DType      []string  `json:"dgraph.type,omitempty"`
}

type Transacciones struct {
	Transacciones []Transaccion `json:"transacciones,omitempty"`
}
