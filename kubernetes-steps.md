# kubernetes steps - v1
## Como crear un deployment?

### Herramientas
1. *kind*

	Guia de instalacion de kind:
	[https://github.com/kubernetes-sigs/kind](https://github.com/kubernetes-sigs/kind)


2. *kubectl*

	Guia de instalacion de kubectl:
	[https://kubernetes.io/docs/tasks/tools/install-kubectl/#before-you-begin](https://kubernetes.io/docs/tasks/tools/install-kubectl/#before-you-begin)

### Procedimiento
*Observacion: los comandos van a ser ejecutados en el directorio raiz del proyecto*

Primero, debemos crear el cluster:

```bash
kind create cluster --name stack-test-cluster --wait 1m --config kind/kind-cluster.yaml
```
Con este comando, kind creara un cluster con el nombre *stack-test-cluster* definido en el flag --name .
 El flag --wait bloquea la consola y espera el tiempo configurado hasta que el control-plane este ejecutandose correctamente, pasado ese tiempo timeoutea la operacion. 
 El flag --config pasa como parametro el archivo de configuracion .yaml que especifiquemos. 
Dentro del archivo *kind-cluster.yaml* podremos observar los siguientes campos:
1. *kind* : tipo de objeto que estas declarando con la configuracion
2. *apiVersion* : version de api que estas utilizando
3. *nodes*: configuracion del cluster

Cada archivo de configuracion en kubernetes declara el estado de un objeto. En este caso, estamos declarando el estado del cluster.

Al ejecutar el comando nos encontraremos con algo similar a :

```
kind create cluster --name stack-test-cluster --wait 1m --config kind/kind-cluster.yaml

Creating cluster "stack-test-cluster" ...
 ✓ Ensuring node image (kindest/node:v1.15.3) 🖼
 ✓ Preparing nodes 📦📦📦📦 
 ✓ Creating kubeadm config 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✓ Joining worker nodes 🚜 
 ✓ Waiting ≤ 1m0s for control-plane = Ready ⏳ 
 • Ready after 21s 💚
Cluster creation complete. You can now use the cluster with:

export KUBECONFIG="$(kind get kubeconfig-path --name="stack-test-cluster")"
kubectl cluster-info

```


Ejecutando el export que nos retorna la ejecucion, apuntaremos al cluster creado anteriormente. De otra, cada vez que querramos usar kubectl debemos especificar en cada comando a que cluster nos estamos refiriendo con el flag --kubeconfig "$(kind get kubeconfig-path --name="stack-test-cluster")"

Ahora, podremos obtener informacion sobre nuestros clusters con el comando: 

```
kubectl get nodes -o wide

NAME                               STATUS   ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                                  KERNEL-VERSION      CONTAINER-RUNTIME
stack-test-cluster-control-plane   Ready    master   13m   v1.15.3   172.17.0.4    <none>        Ubuntu Disco Dingo (development branch)   4.15.0-58-generic   containerd://1.2.6-0ubuntu1
stack-test-cluster-worker          Ready    <none>   12m   v1.15.3   172.17.0.3    <none>        Ubuntu Disco Dingo (development branch)   4.15.0-58-generic   containerd://1.2.6-0ubuntu1
stack-test-cluster-worker2         Ready    <none>   12m   v1.15.3   172.17.0.2    <none>        Ubuntu Disco Dingo (development branch)   4.15.0-58-generic   containerd://1.2.6-0ubuntu1
stack-test-cluster-worker3         Ready    <none>   12m   v1.15.3   172.17.0.5    <none>        Ubuntu Disco Dingo (development branch)   4.15.0-58-generic   containerd://1.2.6-0ubuntu1
```

Lo siguiente, es definir que imagen vamos a introducir en nuestro cluster. Recordarmos que esta imagen ya debio ser buildeada con el comando *docker build ...*

```
kind load docker-image "docker-kubernetes-sandbox:v1" --name stack-test-cluster
```

Los archivos de configuracion(objecto de kubernetes) que le pasamos como parametro a kubectl poseen 4 parametros requeridos para su correcto funcionamiento:
1. *kind* : tipo de objeto que estas declarando con la configuracion
2. *apiVersion* : version de api que estas utilizando
3. *metadata* : informacion sobre el objeto
4. *spec* : dependiendo del kind de objeto que estemos creando, el tipo de informacion que va aqui dentro cambia. Este campo refiere a la configuracion y funcionamiento del objeto en si. Vamos a verlo mas adelante en algun ejemplo.


Nuestro siguiente paso es crear el deployment y el servicio de ese deployment.

Para esto, vamos a utilizar el siguiente comando:
```
kubectl apply -R -f ./kubernetes-config/event-service-demo --kubeconfig "$(kind get kubeconfig-path --name="stack-test-cluster")"
```

Este comando apply, tomara nuestras configuraciones y las creara o actualizara(en caso de existir un cambio) en nuestro cluster, actualizando el estado de los objetos de kubernetes definidos en esos archivos. 
Veamos los casos:
1.  *event-service-demo-srv.yaml *: Declara un objeto del tipo Service de kubernetes que expone el puerto 9290 y lo rutea al puerto 9290 de los servicios con el selector *event-service-demo*
2. *event-service-demo-deployment.yaml*: Declara un objeto del tipo Deployment de kubernetes que declara su matchLabel(con esto matchea el servicio para redirigirle trafico) *event-service-demo* y dentro de template, declara como es la configuracion de los los pods, quiza lo mas interesante, es que declara el estado de los containers de ese deployment( imagen que utilizara el container, que puerto expone, los recursos de memoria y cpu, health-check).

```

kubectl apply -R -f ./kubernetes-config/event-service-demo --kubeconfig "$(kind get kubeconfig-path --name="stack-test-cluster")"

deployment.apps/event-service-demo-deployment created
service/event-service-demo created
```

Vamos a hacer el deploy de nuestros objetos en kubernetes:
```
kubectl wait --for=condition=available --timeout=10s --kubeconfig "$(kind get kubeconfig-path --name="stack-test-cluster")" deploy/event-service-demo-deployment

deployment.extensions/event-service-demo-deployment condition met

```

Podremos ver nuestro deployment con:
```
kubectl get deployments

NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
event-service-demo-deployment   1/1     1            1           2m53s
```

```
kubectl get pods

NAME                                            READY   STATUS    RESTARTS   AGE
event-service-demo-deployment-9dc9ddcf6-hxxjq   1/1     Running   0          3m30s
```

```
curl 127.0.0.1:9290
curl: (7) Failed to connect to localhost port 9290: Conexión rehusada
```

Actualmente, no podemos probar el servicio, ya que se encuentra el la red privada de kubernetes. Podemos utilizar el siguiente comando para fowardear nuestro puerto 9290 al host de kubernetes.
```
kubectl port-forward  event-service-demo-deployment-9dc9ddcf6-hxxjq 9290:9290
```
```
curl 127.0.0.1:9290
Hello, World!
```

Ya terminada la prueba, solo resta limpiar nuestro entorno de kubernetes.

```
kind delete cluster --name="stack-test-cluster"

Deleting cluster "stack-test-cluster" ...
```


### Bibliografia

[https://kubernetes.io/docs/concepts/](https://kubernetes.io/docs/concepts/)
[https://kubernetes.io/docs/tasks/tools/install-kubectl/#before-you-begin](https://kubernetes.io/docs/tasks/tools/install-kubectl/#before-you-begin)
[https://kind.sigs.k8s.io/docs/user/quick-start/](https://kind.sigs.k8s.io/docs/user/quick-start/)
[https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands)