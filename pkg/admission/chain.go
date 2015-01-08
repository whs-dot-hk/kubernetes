/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package admission

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
)

// chainAdmissionHandler is an instance of admission.Interface that performs admission control using a chain of admission handlers
type chainAdmissionHandler []Interface

// New returns an admission.Interface that will enforce admission control decisions
func NewFromPlugins(client client.Interface, pluginNames []string, configFilePath string) Interface {
	plugins := []Interface{}
	for _, pluginName := range pluginNames {
		plugin := InitPlugin(pluginName, client, configFilePath)
		if plugin != nil {
			plugins = append(plugins, plugin)
		}
	}
	return chainAdmissionHandler(plugins)
}

// Admit performs an admission control check using a chain of handlers, and returns immediately on first error
func (admissionHandler chainAdmissionHandler) Admit(a Attributes) (err error) {
	for _, handler := range admissionHandler {
		err := handler.Admit(a)
		if err != nil {
			return err
		}
	}
	return nil
}
