package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Logger return log message
func Logger() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now(), r.Method, r.URL)
		router.ServeHTTP(w, r) // dispatch the request
	})
}

func concat(text1 string, text2 string) string {
	str_slices := []string{text1, text2}
	str_concat := strings.Join(str_slices, "-")
	return str_concat
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//OJO QUE ESTO DEVUELVE CERO SI FALLA
func parseInt(ss string) int {
	aa, err := strconv.Atoi(ss)
	if err != nil {
		aa = 0
	}
	return aa
}
func parseString(ss int) string {
	aa := strconv.Itoa(ss)
	return aa
}

func httpReques(url string, fn FunctionBackBytes) {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fn(body)
}

func analisisClientesJson(jsontext string, fecha int) {
	//jsontext = `[{"id":"1b4dc721","name":"Zsa","age":67},{"id":"59ee2a09","name":"Fiske","age":79},{"id":"b8b70557","name":"Brest","age":85},{"id":"7decab7f","name":"Pages","age":45},{"id":"994bec35","name":"Bohlin","age":63},{"id":"243693e6","name":"Goldsmith","age":68},{"id":"b9606046","name":"Mattland","age":28},{"id":"caebe766","name":"Hersch","age":76},{"id":"d5763a87","name":"Peppel","age":55},{"id":"9b6b8e4f","name":"Kass","age":59},{"id":"56963771","name":"Patrick","age":35},{"id":"50f09cf3","name":"Persian","age":85},{"id":"d03d560c","name":"Binette","age":58},{"id":"5fe59670","name":"Bergin","age":42},{"id":"b9867908","name":"Borchert","age":44},{"id":"12359714","name":"Yonah","age":78},{"id":"bcbe455d","name":"Phenica","age":71},{"id":"9e1a2e49","name":"Zenobia","age":49},{"id":"cc5ab785","name":"Isiah","age":61},{"id":"7818d2d2","name":"Graehl","age":30},{"id":"594bca8a","name":"Zildjian","age":60},{"id":"2e89fe10","name":"Deena","age":24},{"id":"87267f51","name":"Pfeffer","age":67},{"id":"8b2fc525","name":"Hirza","age":31},{"id":"6baaba10","name":"Hoopen","age":40},{"id":"632354b2","name":"LaRue","age":62},{"id":"ece1bb7e","name":"Marteena","age":50},{"id":"20e32ed7","name":"Syl","age":61},{"id":"88cee7ca","name":"Jannery","age":66},{"id":"53db951c","name":"Richer","age":20}]`
	rr := `{"id":"1b4dc721","name":"Zsa","age":67}`
	type Clientex struct {
		UID  string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	clientesolo := Clientex{}
	ee1 := json.Unmarshal([]byte(rr), &clientesolo)
	if ee1 != nil {
		fmt.Printf("%s", "paila")
	}

	clientes := []Clientex{}
	Data := []byte(jsontext)
	err := json.Unmarshal(Data, &clientes)
	if err != nil {
		fmt.Printf("%s", "paila")
	}
	ArrClientes := []Cliente{}
	total := 0
	for i := range clientes {
		cli := clientes[i]
		ncliente := Cliente{
			UID:    "_:elid",
			UIDOLD: cli.UID,
			Name:   cli.Name,
			Age:    cli.Age,
			Date:   fecha,
		}
		ArrClientes = append(ArrClientes, ncliente)
		total++
	}

	ArrClientesViejos := []Cliente{}
	ArrClientesNuevos := []Cliente{}
	totalTTT := 0
	totalNuevos := 0
	for i := range ArrClientes {
		cli := ArrClientes[i]
		query := fmt.Sprintf(`
		{
			people(func: eq(UIDOLD,"%s")) {
			  uid
			}
		}
		`, cli.UIDOLD)
		ConsultaDataBase(query, func(data []byte) {
			algo := string(data)
			if data != nil {
				if algo != `{"people":[]}` {
					ArrClientesViejos = append(ArrClientesViejos, cli)
					fmt.Println(algo + "-Ya existe")
				} else {
					ArrClientesNuevos = append(ArrClientesNuevos, cli)
					totalNuevos++
					fmt.Println(algo + "-Ingresado!")
				}
			} else {
				ArrClientesNuevos = append(ArrClientesNuevos, cli)
				fmt.Println(algo + "-Ingresado!")
			}
			totalTTT++
			if total == totalTTT {
				for i := range ArrClientesNuevos {
					clien := ArrClientesNuevos[i]
					jsonbytes, err := json.Marshal(clien)
					if err != nil {
						fmt.Println(err)
					} else {
						MutacionDataBase(jsonbytes, func(data []byte) {})
					}
				}
				fmt.Printf("total nuevo %s , viejos %s", totalNuevos, totalTTT-totalNuevos)
			}
		})
	}

}

func analisisProductosComillas(texto string, fecha int, sobreescribir bool) {
	fmt.Println(texto)
	filas := strings.Split(texto, "\n")
	productos := []Producto{}

	total := 0
	totalTTT := 0
	analizar := func(fila string) {
		datos := strings.Split(fila, "|")
		if len(datos) >= 2 {
			p, _ := strconv.Atoi(datos[2])
			prod := Producto{
				UID:    "_:elid",
				PRODID: datos[0],
				Name:   datos[1],
				Price:  p,
				Date:   fecha,
			}
			productos = append(productos, prod)
			total++
		} else {
			fmt.Println("algo con este item")
		}
	}
	for filanum := range filas {
		//fmt.Println()
		fila := filas[filanum]
		if !strings.Contains(fila, `"`) {
			analizar(strings.ReplaceAll(fila, "'", "|"))
		} else { //fila con mas comillas
			guardados := []string{}
			swiche := 0 // 0= no lee ,1= lee, 2=cerrar
			palabra := ""
			for i, r := range fila {
				char := string(r)
				if char == `"` {
					palabra += char
					if swiche == 0 {
						swiche = 1
					} else if swiche == 1 {
						swiche = 0
						guardados = append(guardados, palabra)
						palabra = ""
					}
				} else {
					if swiche == 1 {
						palabra += char
					}
				}
				fmt.Println(i, r)
			}
			for item := range guardados { //cambiamos el "hola's" por ${1}
				pl := guardados[item]
				r := "${" + strconv.Itoa(item) + "}"
				fila = strings.ReplaceAll(fila, pl, r)
			}
			fila = strings.ReplaceAll(fila, "'", "|")
			for item := range guardados { //cambiamos el ${1} por "hola's"
				pl := strings.ReplaceAll(guardados[item], `"`, "") //quitamos las comillas
				r := "${" + strconv.Itoa(item) + "}"
				fila = strings.ReplaceAll(fila, r, pl)
			}
			analizar(fila)
		}

	}
	productosResult := []Producto{}
	for i := range productos {
		productolisto := productos[i]
		query := fmt.Sprintf(`
		{
			productos(func: eq(PRODID,"%s")) {
			  uid
			  PRODID
			}
		  }
		`, productolisto.PRODID)
		ConsultaDataBase(query, func(data []byte) {
			algo := string(data)
			totalTTT++
			fmt.Println(algo)
			yaesta := false
			if data != nil {
				if algo != `{"productos":[]}` {
					yaesta = true
					if sobreescribir { //comparamos los datos para actualizar los ya existentes
						ppp := Productos{}
						err33 := json.Unmarshal(data, &ppp)
						if err33 == nil {
							for item := range ppp.Productos {
								p := ppp.Productos[item]
								if p.UID != "" {
									yaesta = false
									productolisto.UID = p.UID
								}
							}
						}
					}

				}
			}
			if !yaesta {
				productosResult = append(productosResult, productolisto) //esto lo usariamos para retornar el los clientes depronto
				//lo ingresamos a la db
				jsonbytes, err := json.Marshal(productolisto)
				if err != nil {
					fmt.Println(err)
				} else {
					MutacionDataBase(jsonbytes, func(data []byte) {
						fmt.Println(string(data))
					})
				}
			}
			if total == totalTTT {
				fmt.Println("finish" + strconv.Itoa(len(productosResult)))
			}
		})
	}

}

func replaceAtIndex1(str string, replacement rune, index int) string {
	out := []rune(str)
	out[index] = replacement
	return string(out)
}

func analisisTransaccionesX(texto string, fecha int, sobreescribir bool) {
	transacciones := []Transaccion{}
	filas := []string{}
	counter := 0
	normalizarProductosIDS := func(arrproductos []string, fn FunctionBackArrStrings) {
		newarr := []string{}
		for i := range arrproductos {
			item := arrproductos[i]
			GetProductoIDfromOldID(item, func(data string) {
				newarr = append(newarr, data)
				if len(newarr) == len(arrproductos) {
					fn(newarr)
				}
			})
		}

	}
	actualizarCliente := func(t Transaccion) {
		type ClienteRicacho struct {
			UID           string        `json:"uid,omitempty"`
			Transacciones []Transaccion `json:"transacciones,omitempty"`
		}
		arrT := []Transaccion{}
		arrT = append(arrT, t) //le metemos la transaccion actual a ese arreglo
		GetClienteIDfromOldID(t.BuyerID, func(uidcliente string) {
			if uidcliente != "" {
				t.BuyerID = uidcliente
				t.Buyer = Cliente{UID: t.BuyerID}
				cli := ClienteRicacho{UID: uidcliente,
					Transacciones: arrT,
				}
				jsonbytes, err := json.Marshal(cli)
				if err != nil {
					fmt.Println(err)
				} else {
					MutacionDataBase(jsonbytes, func(data []byte) {
						//fmt.Println(string(data))
						fmt.Println(t.BuyerID + " ingresado, van " + strconv.Itoa(counter))
					})
				}
			}
		})

	}

	uploadTransaccion := func(t Transaccion) {
		query := fmt.Sprintf(`
		{
			transacciones(func: eq(TRANSID,"%s")) {
			  uid
			  TRANSID
			}
		}
		`, t.TRANSID)
		ConsultaDataBase(query, func(data []byte) {
			algo := string(data)
			fmt.Println(algo)
			yaesta := false
			if data != nil {
				if algo != `{"transacciones":[]}` {
					if sobreescribir {
						ppp := Transacciones{}
						err33 := json.Unmarshal(data, &ppp)
						if err33 == nil {
							if t.UID != "" {
								t.UID = ppp.Transacciones[0].UID
								fmt.Println("Actualizando transaccion: " + t.UID)
							}
						}
					} else {
						yaesta = true
					}

				}
			}
			if !yaesta {
				transacciones = append(transacciones, t) //POS SI QUEREMOS RETORNARLAS
				normalizarProductosIDS(t.ProductIDS, func(data []string) {
					t.ProductIDS = data
					ps := Productos{}
					for i := range data { //metemos con su clase para indexacion
						ps.Productos = append(ps.Productos, Producto{UID: data[i]})
					}
					t.Products = ps
					actualizarCliente(t)
				})
			}
			counter++
			if counter == len(filas) {
				fmt.Println("listo" + strconv.Itoa(counter))
			}
		})
	}
	//recorreomos carater a carater

	conteochats := 0
	type remo struct {
		indice int
		runa   rune
	}
	losremplazados := []remo{}
	for i, r := range texto { //recorremos todo eso letra por letra
		if r == 0 { //si el caracter es extraño
			conteochats++
			losremplazados = append(losremplazados, remo{indice: i, runa: r})
		}
	}
	out := []rune(texto)
	separadorRune := []rune("|")
	for i := range losremplazados { //remplazamos los caracteres raros por barritas
		nremo := losremplazados[i]
		out[nremo.indice] = separadorRune[0]
	}
	texto = string(out)
	fmt.Println("Normalizamos " + strconv.Itoa(conteochats) + " caracteres")
	filas = strings.Split(texto, "||")
	fmt.Println("Vamos a procesar " + strconv.Itoa(len(filas)) + " transacciones")
	for i := range filas {
		fila := filas[i]
		Rfila := strings.Split(fila, "|")
		if len(Rfila) >= 4 {
			//#00006004cf8c|1d646993|143.125.42.1|android|(a1122fc4,1a56a1bf
			ppp := Rfila[4]
			ppp = strings.ReplaceAll(ppp, "(", "")
			ppp = strings.ReplaceAll(ppp, ")", "")
			prodsids := strings.Split(ppp, ",")
			ntrans := Transaccion{
				UID:        "_:elid",
				TRANSID:    Rfila[0],
				BuyerID:    Rfila[1],
				IP:         Rfila[2],
				Device:     Rfila[3],
				ProductIDS: prodsids,
				Date:       fecha,
			}
			uploadTransaccion(ntrans)
			fmt.Println(i)
		}
	}

	//---forma anterior pero muy demorada
	/*for i, r := range texto { //recorremos todo eso letra por letra
		//fmt.Println(i)
		conteochats = i
		if r == 0 { //esto indica el caracter raro
			if saltar { //aqui entra cuando sean dos caracteres raros seguidos osea como un salto de linea
				saltar = false
				Rfila := strings.Split(fila, "|")
				if len(Rfila) >= 4 {
					//#00006004cf8c|1d646993|143.125.42.1|android|(a1122fc4,1a56a1bf
					ppp := Rfila[4]
					ppp = strings.ReplaceAll(ppp, "(", "")
					ppp = strings.ReplaceAll(ppp, ")", "")
					prodsids := strings.Split(ppp, ",")
					ntrans := Transaccion{
						UID:        "_:elid",
						TRANSID:    Rfila[0],
						BuyerID:    Rfila[1],
						IP:         Rfila[2],
						Device:     Rfila[3],
						ProductIDS: prodsids,
						Date:       fecha,
					}
					uploadTransaccion(ntrans)
					filas = append(filas, fila)
					fila = ""
				}

			} else {
				saltar = true
				fila += "|"
				//fmt.Println(i, r, "|")
			}
		} else {
			saltar = false
			//fmt.Println(i, r, string(r))
			fila += string(r)
		}

	}
	fmt.Println(conteochats)*/
}
