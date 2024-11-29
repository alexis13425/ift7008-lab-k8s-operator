# lab8
Code pour l'opérateur du dernier laboratoire de ift7008. Il sert à l'introduction au développement et déploiement d'operateur.

## Description
Repository pour le laboratoire sur les operators Kubernetes du cours ift-7008. L'opérateur va créer un deployment et un service. Cette opérateur est très simple et ne représente pas une bonne utilisation du Operator pattern de Kubernetes. Créer un opérateur est utile quand il est nécessaire d'ajouter des fonctionnalités à Kubernetives ou pour le déploiement d'application complexe.
### prérequis
- installer go (https://go.dev/doc/install)
- operator-sdk (https://master.sdk.operatorframework.io/docs/installation/)
### étape pour faire le laboratoire
1. Créer un cluster kind
```sh
kind create cluster
```
2. Initialiser le projet grâce à l'outil operator-sdk
```sh
mkdir lab8
cd lab8
// initialiser le framework de notre operator
operator-sdk init --domain=example.com --repo=github.com/example-inc/lab8-operator
// installer les dépendances nécessaire pour le laboratoire
go mod tidy
go mod vendor
```
3. Initialisation du controller pour l'opérateur. La CustomResource aura le nom de Traveller
```sh
operator-sdk create api --group traveller --version v1 --kind Traveller --resource --controller
```
L'ajout de --group traveller changera la version de l'api pour `group.domain/version`. Dans le laboratoire, l'api version sera `traveller.example.com/v1`
Après avoir créer votre controller, il y aura un `fichier internal/controller/traveller_controller.go`. Il faudra ajouter les fichiers `traveller_controller.go`, `service.go` et `deployment.go` ce trouvant dans le répertoire `internal/controller` de ce repo github.

4. Générer les fichier manifests de l'opérateur.
```sh
make generate
make manifests
```
make manifests va créer les fichiers manifests présents dans le répertoire `config/`. Les annotations kubebuilder dans le fichier `traveller_controller.go` spécifie les configurations nécessaires pour créer ces fichiers. Nous avons ajouté deux annotations.
```
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
```
Puisque l'opérateur doit créer une Deployment et un Service, il faut ajouter les droits d'accès et de creation pour ces ressources.
5. Build l'image de l'opérateur
```sh
make bundle IMG="<some-registry>/lab8-operator:v0.0.1"
make docker-build docker-push IMG="<some-registry>/lab8-operator:v0.0.1"

// modifier l'image du manifest du controller pour l'image qui vient d'être créé.
nano config/manager/manager.yaml
```
6. Déployer l'opérateur ainsi que la custom resource.
```sh
make deploy
// vérifier que le déploiement a fonctionné.
kubectl get all -n <Nom du namespace qui a été créé>
// ajouter la custom resource
kubectl apply -f config/samples/traveller_v1_traveller.yaml
// vérifier que l'opérateur à eu le comportement désiré.
kubectl get svc
kubectl get deployment
```

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**
```sh
make docker-build docker-push IMG=<some-registry>/lab8:tag
```


**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/lab8:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/lab8:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/lab8/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

