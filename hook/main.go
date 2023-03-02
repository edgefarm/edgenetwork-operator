package main

import (
	"io/ioutil"
	api "leaf-nats-controller/api/edgenetwork/v1alpha1"
	"log"
	"net/http"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

type SyncRequest struct {
	Parent   api.EdgeNetwork     `json:"parent"`
	Children SyncRequestChildren `json:"children"`
}

type SyncRequestChildren struct {
	Pods map[string]*v1.Pod `json:"Pod.v1"`
}

type SyncResponse struct {
	Status   api.EdgeNetworkStatus `json:"status"`
	Children []runtime.Object      `json:"children"`
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received sync request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(string(body))
	request := &SyncRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	specGenerator := &NatsSpecGenerator{}

	response, err := specGenerator.GenerateManifestResponse(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	log.Println("Starting hook")

	http.HandleFunc("/sync", syncHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
