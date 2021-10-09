// Copyright 2017 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"reflect"
	"testing"

	"github.com/poleljy/kube-boost/api"
	metricapi "github.com/poleljy/kube-boost/integration/metric/api"
	"github.com/poleljy/kube-boost/resource/dataselect"
	"gopkg.in/square/go-jose.v2/json"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGetNodeList(t *testing.T) {
	kubeConfig := "D:/kube-47.config"
	//kubeConfig := "/root/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err)
	}
	clientSet, err := kube.NewForConfig(config)
	if err != nil {
		t.Fatal(err.Error())
	}

	// filter
	//dsQuery := dataselect.NewLabelSelectQuery(map[string]string{"ai.piesat.cn/jupyter": ""})
	dsQuery := dataselect.NewLabelSelectQuery(map[string]string{"ai.piesat.cn/model": "evaluation"})
	actual, err := GetNodeList(clientSet, dsQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range actual.Nodes {
		content, _ := json.MarshalIndent(node.AllocatedResources, "", " ")
		t.Log("Node:", string(content))

		//detail, err := GetNodeDetail(clientSet, nil, node.ObjectMeta.Name, dataselect.NoDataSelect)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//detailContent, _ := json.MarshalIndent(detail, "", " ")
		//t.Log("Node detail:", string(detailContent))
	}
	return

	cases := []struct {
		node     *v1.Node
		expected *NodeList
	}{
		{
			&v1.Node{
				ObjectMeta: metaV1.ObjectMeta{Name: "test-node"},
				Spec: v1.NodeSpec{
					Unschedulable: true,
				},
			},
			&NodeList{
				ListMeta: api.ListMeta{
					TotalItems: 1,
				},
				Errors:            []error{},
				CumulativeMetrics: make([]metricapi.Metric, 0),
				Nodes: []Node{{
					ObjectMeta: api.ObjectMeta{Name: "test-node"},
					TypeMeta:   api.TypeMeta{Kind: api.ResourceKindNode},
					Ready:      "Unknown",
					AllocatedResources: NodeAllocatedResources{
						CPURequests:            0,
						CPURequestsFraction:    0,
						CPULimits:              0,
						CPULimitsFraction:      0,
						CPUCapacity:            0,
						MemoryRequests:         0,
						MemoryRequestsFraction: 0,
						MemoryLimits:           0,
						MemoryLimitsFraction:   0,
						MemoryCapacity:         0,
						AllocatedPods:          0,
						PodCapacity:            0,
						PodFraction:            0,
					},
				},
				},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.node)
		actual, _ := GetNodeList(fakeClient, dataselect.NoDataSelect, nil)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetNodeList() == \ngot: %#v, \nexpected %#v", actual, c.expected)
		}
	}
}
