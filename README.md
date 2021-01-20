# loschinitos_backend_go

Este es el Backend creado en go para LosChinitos

<p align="center">
  <br>
  <a href="https://github.com/luismgluis/loschinitos" alt="LosChinitos Frontend in Vue"></a>
  <br>
</p>

App que tiene como objetivo conocer Vue.js, Golang, Dgraph database y crear una app
funcional con estas tecnologias

Este backend es necesario para el funcionamiento del Frontend de <a href="https://github.com/luismgluis/loschinitos" alt="LosChinitos Vue"></a>

Para poder usar este backend debe crear un Docker con la siguiente imagen :

```
docker run -it -p 6080:6080 -p 8080:8080 -p 9080:9080 -p 8000:8000 -v /mnt/dgraph:/dgraph dgraph/standalone:v20.03.0

```

Puede ver otras alternativas a docker en <a href="https://dgraph.io/downloads" alt="Dgraph/Downloads"></a>, aunque personalmente no tuve suerte con las demas soluciones, solo con docker.

Puede obtener mas informacion de como iniciar una maquina Dgraph en <a href="https://dgraph.io/docs/tutorial-1/" alt="tutorial en la pagina oficial"></a>.

Para ejecutar el codigo se esta usando Golang en:
```
go version go1.15.6 linux/amd64
```

Finalmente para ejecutar el servidor usamos:
```
go run .
```

Tener en cuenta que el archivo main.go se ejecuta el servidor en el puerto 3000

```
http.ListenAndServe(":3000", Logger())
```

Para fines especificos fue creado el endpoint /importx
el cual importa los datos desde urls especificas con diferente complejidad cada una, clientes en JSON, productos separado por comillas y saltos de linea y transacciones en un archivo desconocido con separadores desconocidos

```
http://localhost:3000/importx
```

adicionalmente las urls especificas admiten la posibilidad de pasar como parametro date la fecha en formato UNIX y lo podemos usar nosotros tambien, por ejemplo para extraer los datos de 20/01/2020 usamos 1611118800 de la siguiente forma

```
http://localhost:3000/importx/1611118800
```

la importacion de productos en la funcion analisisTransaccionesX puede establecer si quiere que se sobre escriban o no las transacciones por default es true

```
func dataFromInternet(w http.ResponseWriter, r *http.Request) {
  .
  ..
  ...
    analisisTransaccionesX(transacciones, fechaInt, true)
```

Es importante realizar la importacion con en endpoint /importx antes de usar el frontend en vue.