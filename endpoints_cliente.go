package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/mholt/binding"
)

//-------------- API ENDPOINT ------------------//
func GetClienteIDfromOldID(olduid string, fn FunctionBackString) {
	q := fmt.Sprintf(`
	{
		clientes(func: eq(UIDOLD,"%s")) {
		  uid
		  UIDOLD
		}
	}
	`, olduid)

	ConsultaDataBase(q, func(data []byte) {
		ccc := Clientes{}
		err3312 := json.Unmarshal(data, &ccc)
		if err3312 == nil {
			if ccc.Clientes != nil {
				clientes := ccc.Clientes
				micliente := clientes[0]
				if micliente.UID != "" {
					fn(micliente.UID)
				} else {
					fn("")
				}
			} else {
				fn("")
			}
		} else {
			fn("")
		}

	})
}

// AllClientes returns todos los clientes en la DB
func AllClientes(w http.ResponseWriter, r *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	txn := dg.NewTxn()
	defer txn.Discard(ctx)
	q := `
	  {
		clientes(func: has(UIDOLD)) {
		  name
		  age
		  uid
		  UIDOLD
		}
	  } 
	  `
	//clientes := []Cliente{}
	//res, err := txn.QueryWithVars(ctx, q, map[string]string{"$a": "Alice"})
	res, err := txn.Query(ctx, q)

	//s := string(`{"operation": "get", "key": "example"}`)

	if err == nil {
		fmt.Printf("%s\n", res.Json)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res.Json)
		//respondwithJSON(w, http.StatusOK, res.Json)
	} else {
		respondwithJSON(w, http.StatusOK, `{"result":"error"}`)
	}
}

func GetClienteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	query := fmt.Sprintf(`
	{
		cliente(func: uid(%s)){
			uid
			name
			age
			UIDOLD
		}
	}
	`, id)
	ConsultaDataBaseJson(query, func(data string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})
}

func GetClienteDetailsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	type ClienteDetallado struct {
		Cliente               Cliente       `json:"cliente,omitempty"`
		Transacciones         []Transaccion `json:"transacciones,omitempty"`
		Productos             []Producto    `json:"productos,omitempty"`
		Clientesporip         []Cliente     `json:"clientesporip,omitempty"`
		ProductosRecomendados []Producto    `json:"productosrecomendados,omitempty"`
	}
	supercli := ClienteDetallado{}
	query := fmt.Sprintf(`
	{
		clientes(func: uid(%s)){
			uid
			name
			age
			UIDOLD
			transacciones {
			  uid
			  TRANSID
			  buyerid
			  ip
			  device
			  produtids
			  products
			  date
			}
		}
	}
	`, id)
	retornarError := func() {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":"error"}`))
	}
	ConsultaDataBase(query, func(data []byte) {
		clis := Clientes{}
		err33 := json.Unmarshal(data, &clis)
		if err33 != nil {
			retornarError()
		} else {
			tss := clis.Clientes[0].Transacciones
			supercli.Cliente = clis.Clientes[0]
			supercli.Transacciones = tss

			productosporrevisar := []string{}
			for i := range tss {
				transa := tss[i]
				for ii := range transa.ProductIDS {
					idprod := transa.ProductIDS[ii]
					if !contains(productosporrevisar, idprod) {
						productosporrevisar = append(productosporrevisar, idprod)
					}
				}
			}

			prodsrevisados := []Producto{}
			realizados := 0
			for i := range productosporrevisar {
				id := productosporrevisar[i]

				GetProductoByIDData(id, func(data Producto) {
					prodsrevisados = append(prodsrevisados, data)
					realizados++
					if len(productosporrevisar) == realizados {
						supercli.Productos = prodsrevisados
						jsonbytes, err := json.Marshal(supercli)
						if err == nil {
							w.Header().Set("Content-Type", "application/json")
							w.Write(jsonbytes)
						} else {
							retornarError()
						}

					}
				})
			}

		}

	})
}

// PostCliente crea un nuevo cliente
func PostCliente(w http.ResponseWriter, r *http.Request) {
	cliente := new(Cliente)
	//esto nos ayuda a asignar los datos enviados desde frontend a el struct
	if errs := binding.Bind(r, cliente); errs != nil {
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf(`
	{
		cliente(func: uid(%s)){
			uid
			name
			age
		}
	}
	`, cliente.UID)

	ConsultaDataBase(query, func(data []byte) {
		if data != nil { //el cliente existe entonces decimos rechazamos la creacion
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result":"error"}`))
		} else { //no existe entonces si lo creamos
			cliente.UID = "_:elid" //asignamos esto para que se ponga el ID automaticamente
			jsonbytes, err := json.Marshal(cliente)
			if err != nil {
				log.Fatal(err)
			}
			MutacionDataBase(jsonbytes, func(data []byte) {
				w.Header().Set("Content-Type", "application/json")
				if data == nil {
					w.Write([]byte(`{"result":"ok"}`))
				} else {
					w.Write([]byte(`{"result":"error"}`))
				}
			})
		}
	})
}

// PutCliente actualiza un cliente
func PutCliente(w http.ResponseWriter, r *http.Request) {

	cliente := new(Cliente)
	//esto nos ayuda a asignar los datos enviados desde frontend a el struct
	if errs := binding.Bind(r, cliente); errs != nil {
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf(`
	{
		cliente(func: uid(%s)){
			uid
			name
			age
		}
	}
	`, cliente.UID)

	ConsultaDataBase(query, func(data []byte) {
		if data != nil { //el cliente existe entonces lo actualizamos
			jsonbytes, err := json.Marshal(cliente)
			if err != nil {
				log.Fatal(err)
			}
			MutacionDataBase(jsonbytes, func(data []byte) {
				w.Header().Set("Content-Type", "application/json")
				if data == nil {
					//fmt.Printf("%s", string(data))
					w.Write([]byte(`{"result":"ok"}`))
				} else {
					w.Write([]byte(`{"result":"error"}`))
				}
			})
		} else { //no existe entonces si lo rechazamos
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result":"error"}`))
		}
	})
}

// DeleteCliente remove a spesific post
func DeleteCliente(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	query := fmt.Sprintf(`
	{
		"delete": [
			{
				"uid": "%s"
			}
		]
	}
	`, id)
	ConsultaDataBase(query, func(data []byte) {
		w.Header().Set("Content-Type", "application/json")
		if data != nil {
			//ok
			w.Write([]byte(`{"result":"ok"}`))
		} else {
			//bad
			w.Write([]byte(`{"result":"error"}`))
		}
	})
}
