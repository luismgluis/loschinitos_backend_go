package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/mholt/binding"
)

//-------------- API ENDPOINT ------------------//

// AllClientes returns todos los clientes en la DB
func AllTransaccion(w http.ResponseWriter, r *http.Request) {
	dg, cancel := getDgraphClient()
	defer cancel()

	ctx := context.Background()

	txn := dg.NewTxn()
	defer txn.Discard(ctx)
	q := `
	  {
		clientes(func: has(name)) {
		  name
		  age
		  uid
		  followers {
			uid
			name
		  }
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

func AllTransaccionesRangeData(inicio int, final int, fn FunctionBackTransacciones) {

	query := `
	{
		transacciones(func: eq(dgraph.type,["Transaccion"]), first: 1000 )  {
			uid
			TRANSID
			buyerid
			buyer {
				uid
				name
			}
			date
			ip
			device
			produtids
		}
	}
	`
	/* //Deberiamos realizar una consulta de este tipo peeero aun no logro que funcione
		query := fmt.Sprintf(`
		`, parseString(inicio), parseString(final)) .... -> query ->
		{
			transacciones(func: eq(dgraph.type,["Transaccion"])) @filter( ge(date,16) AND le(date,1611118800) )  {
				uid
				TRANSID
				buyerid
				buyer {
					uid
					name
				}
				date
				ip
				device
				produtids
			}
	}
	*/
	fmt.Println(query)
	ConsultaDataBase(query, func(data []byte) {
		trans := Transacciones{}
		err33 := json.Unmarshal(data, &trans)
		if err33 == nil {
			transaccionesO := trans.Transacciones
			transacciones := []Transaccion{}
			tra := []Transaccion{}
			for i := range transaccionesO {
				t := transaccionesO[i]
				if t.Date <= final && t.Date >= inicio {
					transacciones = append(transacciones, t)
				}
			}
			for i := range transacciones {
				t := transacciones[i]
				GetClienteIDfromOldID(t.BuyerID, func(data Cliente) {
					t.Buyer = data
					tra = append(tra, t)
					if len(transacciones) == len(tra) {
						fn(Transacciones{Transacciones: tra})
					}
				})
			}
		}
	})
}

func AllTransaccionesRange(w http.ResponseWriter, r *http.Request) {
	rango := chi.URLParam(r, "range") //conseguimos el ID pasado por URL
	fmt.Printf("%s", rango)

	arrgo := strings.Split(rango, "-")
	inicio := parseInt(arrgo[0])
	final := parseInt(arrgo[1])
	AllTransaccionesRangeData(inicio, final, func(data Transacciones) {
		jsonbytes, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonbytes)
		}
	})
}

func GetTransaccionByIDData(idtrans string, fn FunctionBackTransaccion) {

	query := fmt.Sprintf(`
	{
		transacciones(func: uid(%s)) {
			uid
			ip
			device
			produtids
			date
			buyer {
				uid
			}
		}
	}
	`, idtrans) //0x6e5
	ConsultaDataBase(query, func(data []byte) {
		trans := Transacciones{}
		err33 := json.Unmarshal(data, &trans)
		if err33 == nil {
			transaccion := trans.Transacciones[0]
			conteo1 := len(transaccion.ProductIDS)
			GetClienteByIDData(transaccion.Buyer.UID, func(data Cliente) {
				transaccion.Buyer = data
				prods := transaccion.ProductIDS
				for i := range prods {
					prod := prods[i]
					GetProductoByIDData(prod, func(data2 Producto) {
						transaccion.Products.Productos = append(transaccion.Products.Productos, data2)
						if conteo1 == len(transaccion.Products.Productos) {
							fn(transaccion)
						}
					})
				}
			})
		}
	})
}

func GetTransaccionesAsociadasIP(ipBusqueda string, fn FunctionBackTransacciones) {
	query := fmt.Sprintf(`
	{
		transacciones(func: eq(ip,"%s") ) @filter(has(TRANSID) AND eq(dgraph.type,["Transaccion"]) )   {
			uid
			TRANSID
			buyerid
			ip
			device
			produtids
			date
		}
	}
	`, ipBusqueda)
	ConsultaDataBase(query, func(data []byte) {
		trans := Transacciones{}
		err33 := json.Unmarshal(data, &trans)
		if err33 == nil {
			transacciones := trans.Transacciones
			tra := []Transaccion{}
			for i := range transacciones {
				t := transacciones[i]
				GetClienteIDfromOldID(t.BuyerID, func(data Cliente) {
					t.Buyer = data
					tra = append(tra, t)
					if len(transacciones) == len(tra) {
						fn(Transacciones{Transacciones: tra})
					}
				})
			}
		}
	})
}

func GetTransaccionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //conseguimos el ID pasado por URL
	fmt.Printf("%s", id)
	query := fmt.Sprintf(`
	{
		cliente(func: uid(%s)){
			uid
			name
			age
		}
	}
	`, id)
	ConsultaDataBaseJson(query, func(data string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})
}

// PostCliente crea un nuevo cliente
func PostTransaccion(w http.ResponseWriter, r *http.Request) {
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
					//fmt.Printf("%s", string(data))
					w.Write([]byte(`{"result":"ok"}`))
				} else {
					w.Write([]byte(`{"result":"error"}`))
				}
			})
		}
	})
}

// PutCliente actualiza un cliente
func PutTransaccion(w http.ResponseWriter, r *http.Request) {

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
func DeleteTransaccion(w http.ResponseWriter, r *http.Request) {
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
