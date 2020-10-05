module github.com/stakater/slack-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/operator-framework/operator-sdk v0.18.1
	github.com/slack-go/slack v0.6.5
	github.com/stakater/operator-utils v0.1.2
	github.com/stretchr/testify v1.5.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.2
)

replace k8s.io/client-go => k8s.io/client-go v0.18.8
