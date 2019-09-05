#https://docs.docker.com/develop/develop-images/dockerfile_best-practices/

#Imagen base sobre la cual vamos a trabajar
FROM golang:1.12-alpine

#Copiamos los archivos que se encuentran en nuestro directorio donde se encuentra el Dockerfile que vamos a buildear
COPY . .

#Cambiamos de directorio
WORKDIR main/

#Ejecutamos la siguiente sentencia para generar el ejecutable que levanta la aplicacion
RUN go build main.go

#Indicamos en que puerto nuestro container va a escuchar el trafico
EXPOSE 9290

#Comando para levantar nuestra app
CMD ["./main"]