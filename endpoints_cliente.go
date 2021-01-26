package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/mholt/binding"
)

//-------------- API ENDPOINT ------------------//
func GetClienteIDfromOldID(olduid string, fn FunctionBackCliente) {
	q := fmt.Sprintf(`
	{
		clientes(func: eq(UIDOLD,"%s") ) @filter(has(name) AND has(UIDOLD) AND eq(dgraph.type,["Cliente"]) )   {
		  uid
		  UIDOLD
		  name
		  age
		  avatar
		  date
		  dgraph.type
		}
	}
	`, olduid)

	ConsultaDataBase(q, func(data []byte) {
		ccc := Clientes{}
		micliente := Cliente{}
		err3312 := json.Unmarshal(data, &ccc)
		if err3312 == nil {
			if ccc.Clientes != nil {
				clientes := ccc.Clientes
				micliente = clientes[0]
				if micliente.UID != "" {
					fn(micliente)
				} else {
					fn(micliente)
				}
			} else {
				fn(micliente)
			}
		} else {
			fn(micliente)
		}

	})
}

// AllClientes returns todos los clientes en la DB
func AllClientes(w http.ResponseWriter, r *http.Request) {

	q := `
	{
		clientes(func: eq(dgraph.type,["Cliente"])) @filter(has(name))  @cascade {
		  uid
		  name
		  age
		  UIDOLD
		  dgraph.type
		}
	} 
	 `

	ConsultaDataBase(q, func(data []byte) {
		ccc := Clientes{}
		err3312 := json.Unmarshal(data, &ccc)
		if err3312 == nil {
			fmt.Println("Clientes = " + parseString(len(ccc.Clientes)))
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else {
			respondwithJSON(w, http.StatusOK, `{"result":"error"}`)
		}
	})
}

func GetClienteByIDData(id string, fn FunctionBackCliente) {
	/*query := fmt.Sprintf(`
	{
		clientes(func: uid(%s)){
			uid
			name
			age
			UIDOLD
		}
	}
	`, id)*/
	q := fmt.Sprintf(`
	{
		clientes(func: uid(%s) ) @filter(has(name) AND has(UIDOLD) AND eq(dgraph.type,["Cliente"]) )   {
		  uid
		  UIDOLD
		  name
		  age
		  avatar
		  date
		  dgraph.type
		}
	}
	`, id)
	ConsultaDataBase(q, func(data []byte) {
		clis := Clientes{}
		err33 := json.Unmarshal(data, &clis)
		if err33 == nil {
			cli := clis.Clientes[0]
			fn(cli)
		}
	})
}
func GetClienteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	GetClienteByIDData(id, func(data Cliente) {
		clis := Clientes{}
		clis.Clientes = append(clis.Clientes, data)
		jsonbytes, err := json.Marshal(clis)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonbytes)
		} else {
			respondwithJSON(w, http.StatusOK, `{"result":"error"}`)
		}
	})

}

func GetClienteDetailsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	type ClienteDetallado struct {
		Cliente               Cliente       `json:"cliente,omitempty"`
		Transacciones         []Transaccion `json:"transacciones,omitempty"`
		Productos             []Producto    `json:"productos,omitempty"`
		Transaccionesporip    []Transaccion `json:"transaccionesporip,omitempty"`
		ProductosRecomendados []Producto    `json:"productosrecomendados,omitempty"`
	}
	supercli := ClienteDetallado{}
	query := fmt.Sprintf(`
	{
		clientes(func: uid(%s) ) @filter(has(name) AND has(UIDOLD) AND eq(dgraph.type,["Cliente"]) )   {
			uid
			name
			age
			UIDOLD
			dgraph.type
			transacciones @filter(has(TRANSID) AND has(ip)) {
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
			unaip := ""
			productosporrevisar := []string{}
			for i := range tss {
				transa := tss[i]
				if transa.IP != "" {
					unaip = transa.IP
				}
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
						GetTransaccionesAsociadasIP(unaip, func(data Transacciones) {
							supercli.Transaccionesporip = data.Transacciones
							realizados = 0
							asegurador := []string{}
							terminar := func(supercli ClienteDetallado) {
								jsonbytes, err := json.Marshal(supercli)
								if err == nil {
									w.Header().Set("Content-Type", "application/json")
									w.Write(jsonbytes)
								} else {
									retornarError()
								}
							}
							//reseteamos el contador para reusarlo
							for ii := range supercli.Transaccionesporip { //analizamos las transacciones para obtener un producto de los comprados
								trr := supercli.Transaccionesporip[ii]
								if len(trr.ProductIDS) == 0 {
									realizados++
									if len(supercli.Transaccionesporip) == realizados {
										terminar(supercli)
									}
								} else {
									//aqui podriamos consultar todos los productos y sugerir el mas costoso o el mas barato
									//tambien podriamos comparar con las transacciones del cliente para sugerir productos que no
									//hubiera comprado antes o productos en el mismo rango de precio
									GetProductoByIDData(trr.ProductIDS[0], func(productoresult Producto) {
										if !contains(asegurador, productoresult.UID) {
											supercli.ProductosRecomendados = append(supercli.ProductosRecomendados, productoresult)
											asegurador = append(asegurador, productoresult.UID)
										}
										realizados++
										if len(supercli.Transaccionesporip) == realizados {
											terminar(supercli)
										}
									})
								}

							}
						})

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
