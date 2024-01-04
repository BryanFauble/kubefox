/*
Copyright © 2023 XigXog

This Source Code Form is subject to the terms of the Mozilla Public License,
v2.0. If a copy of the MPL was not distributed with this file, You can obtain
one at https://mozilla.org/MPL/2.0/.
*/

package v1alpha1

import (
	"github.com/xigxog/kubefox/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterVirtualEnvSpec struct {
	ReleasePolicies *ReleasePolicy `json:"releasePolicies,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clustervirtualenvs,scope=Cluster
// +kubebuilder:printcolumn:name="Title",type=string,JSONPath=`.details.title`,priority=1
// +kubebuilder:printcolumn:name="Description",type=string,JSONPath=`.details.description`,priority=1

type ClusterVirtualEnv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec    ClusterVirtualEnvSpec `json:"spec,omitempty"`
	Data    api.VirtualEnvData    `json:"data,omitempty"`
	Details VirtualEnvDetails     `json:"details,omitempty"`
}

// +kubebuilder:object:root=true
type ClusterVirtualEnvList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ClusterVirtualEnv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterVirtualEnv{}, &ClusterVirtualEnvList{})
}
