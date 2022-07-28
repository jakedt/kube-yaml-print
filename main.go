package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/xlab/treeprint"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type Resource struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}

func main() {
	input := os.Stdin
	if len(os.Args) > 1 {
		var err error
		input, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	d := yaml.NewDecoder(input)

	namespaced := make(map[string]map[string][]*Resource, 0)

	// Parse the resources
	for {
		// create new spec here
		spec := new(Resource)
		// pass a reference to spec reference
		err := d.Decode(&spec)
		// check it was parsed
		if spec == nil {
			continue
		}
		// break the loop in case of EOF
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}

		var namespace string
		if spec.ApiVersion == "v1" && spec.Kind == "Namespace" {
			namespace = spec.Metadata.Name
		} else if spec.Metadata.Namespace != "" {
			namespace = spec.Metadata.Namespace
		}

		gvk := fmt.Sprintf("%s/%s", spec.ApiVersion, spec.Kind)

		namespacedGVK, ok := namespaced[namespace]
		if !ok {
			namespacedGVK = make(map[string][]*Resource, 0)
			namespaced[namespace] = namespacedGVK
		}

		gvkResources := namespacedGVK[gvk]
		gvkResources = append(gvkResources, spec)
		namespacedGVK[gvk] = gvkResources
	}

	// Build the tree
	tree := treeprint.New()

	nsNames := make([]string, 0, len(namespaced))
	for ns := range namespaced {
		nsNames = append(nsNames, ns)
	}

	sort.Slice(nsNames, func(i, j int) bool {
		return len(nsNames[i]) < len(nsNames[j])
	})

	for _, ns := range nsNames {
		gvks := namespaced[ns]

		nsNode := tree
		if ns != "" {
			nsNode = tree.AddMetaBranch("namespace", ns)
		}

		nsGVKNames := make([]string, 0, len(gvks))
		for gvk := range gvks {
			nsGVKNames = append(nsGVKNames, gvk)
		}

		sort.Slice(nsGVKNames, func(i, j int) bool {
			return len(nsGVKNames[i]) < len(nsGVKNames[j])
		})

		for _, gvkName := range nsGVKNames {
			resources := gvks[gvkName]
			gvkNode := nsNode.AddBranch(gvkName)
			for _, resource := range resources {
				gvkNode.AddNode(resource.Metadata.Name)
			}
		}
	}

	// Print the tree
	fmt.Println(tree.String())
}
