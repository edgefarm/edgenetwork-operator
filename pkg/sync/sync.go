package sync

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	api "github.com/edgefarm/edgenetwork-operator/apis/edgenetwork/v1alpha1"
	"github.com/edgefarm/edgenetwork-operator/pkg/generate"
)

type SyncRequest struct {
	Parent   api.EdgeNetwork     `json:"parent"`
	Children SyncRequestChildren `json:"children"`
}

type SyncRequestChildren struct {
	DaemonSet  map[string]*appsv1.DaemonSet `json:"DaemonSet.apps/v1"`
	Configmaps map[string]*v1.ConfigMap     `json:"Configmap.v1"`
}

type SyncResponse struct {
	Status   *api.EdgeNetworkStatus `json:"status"`
	Children []runtime.Object       `json:"children"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug().Msgf("Received sync request [raw]: ", string(body))

	request := &SyncRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		log.Error().Msgf("Error unmarshalling request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info().Msgf("Request is for %s in namespace %s", request.Parent.Name, request.Parent.Namespace)

	manifests, err := generate.Manifests(&request.Parent)
	if err != nil {
		log.Error().Msgf("Error generating manifests: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := SyncResponse{
		Status:   CalculateStatus(request),
		Children: manifests,
	}

	body, err = json.Marshal(&response)
	if err != nil {
		log.Error().Msgf("Error marshalling response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug().Msg(string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
