# docker steps - v1
## Como buildear una imagen y garantizarnos que esta corriendo?

### Build de la imagen

Nos situamos en la carpeta donde se encuentra nuestro Dockerfile
y corremos el siguiente comando:
```bash
docker image build -t docker-kubernetes-sandbox:v1 . 
```

Podremos observar si la imagen se genero correctamente con el comando 

```bash
docker image ls
```
### Test de la imagen

Una vez generada la imagen, vamos a testear que funcione correctamente.
Primero vamos a generar y ejecutar un container con la imagen generada:

```bash
 docker container run -d --name docker-kubernetes-image-test  -p 9290:9290 docker-kubernetes-sandbox:v1 
```

De esta forma, corremos un container con la imagen *docker-kubernetes-sandbox:v1*, bindeamos el puerto 9290 de nuestro host a la red privada de kubernetes, con lo cual todo el trafico que pase por nuestro 127.0.0.1:9290 se redirigira a nuestro container.

Podemos ver los containers que estan ejecutandose con el comando: 
```bash
 docker container ls
```
y se mostrara de la siguiente forma:
```
CONTAINER ID        IMAGE                          COMMAND             CREATED             STATUS              PORTS                    NAMES
17f0c88c545e        docker-kubernetes-sandbox:v1   "./main"            3 minutes ago       Up 3 minutes        0.0.0.0:9290->9290/tcp   docker-kubernetes-image-test
```
Tambien podemos acceder a los logs del container con el comando 

```bash
docker container logs docker-kubernetes-image-test
```
En caso de querer obtener informacion mas detallada del container, podremos hacerlo con el siguiente comando:

```bash
docker container inspect docker-kubernetes-image-test
```


Lo unico que nos resta hacer es testear que, cuando le peguemos a 127.0.0.1:9290 nos devuelva un **Hello, World!**

```bash
curl 127.0.0.1:9290
```

Una vez, que testeamos que la imagen funciona, podemos interrumpir y eliminar el container.

```bash
docker container stop docker-kubernetes-image-test
```

Ahora podremos observar, que si hacemos un 
```bash
 docker container ls
```
este comando no nos mostrara nuestro container

```bash
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES

```

Podremos ver nuestro container que no esta ejecutandose con el flag -a
```bash
 docker container ls -a
```
```

CONTAINER ID        IMAGE                          COMMAND             CREATED             STATUS                     PORTS               NAMES
17f0c88c545e        docker-kubernetes-sandbox:v1   "./main"            17 minutes ago      Exited (2) 2 minutes ago                       docker-kubernetes-image-test
```

Nos resta unicamente eliminar el container con el comando:



``` bash
docker container rm docker-kubernetes-image-test
```


```
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```