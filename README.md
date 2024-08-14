-----------------------------------
https://github.com/kubernetes/client-go?tab=readme-ov-file#client-go
-----------------------------------
-----------------------------------
GO K8S Client
-----------------------------------
-----------------------------------
Cluster Durumu:
Cluster'da 9 pod var
Cluster'da 5 namespace var
Cluster'da 1 node var
Son 1 saatte 0 event var
Cluster'da 0 PersistentVolumeClaim var

Detaylı pod kontrolü:
Pod alpine-deployment-548dbddc9b-dnq9r namespace default içinde bulundu
Pod durumu: Running
Pod IP: 10.42.0.45
Node: enesce-toshiba
-----------------------------------
-----------------------------------
go mod init go-k8s-client
go get k8s.io/client-go@v0.27.0
go get k8s.io/apimachinery@v0.27.0
go mod tidy
go run main.go --kubeconfig=/home/enesce/kubeconfig
-----------------------------------

