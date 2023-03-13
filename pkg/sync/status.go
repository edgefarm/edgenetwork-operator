package sync

import (
	"strconv"

	api "github.com/edgefarm/edgenetwork-operator/apis/edgenetwork/v1alpha1"
)

func CalculateStatus(request *SyncRequest) *api.EdgeNetworkStatus {
	for _, ds := range request.Children.DaemonSet {
		// return after the first one because we only have one daemonset
		return &api.EdgeNetworkStatus{
			Current: strconv.Itoa(int(ds.Status.CurrentNumberScheduled)),
			Desired: strconv.Itoa(int(ds.Status.DesiredNumberScheduled)),
			Ready:   strconv.Itoa(int(ds.Status.NumberReady)),
		}
	}
	return &api.EdgeNetworkStatus{
		Current: "0",
		Desired: "0",
		Ready:   "0",
	}
}
