module github.com/evalsocket/policyreport-octant-plugin

go 1.16

require (
	github.com/fatih/color v1.12.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/vmware-tanzu/octant v0.20.0
	go.uber.org/zap v1.16.1-0.20210329175301-c23abee72d19 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210421230115-4e50805a0758 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver v0.21.0 // indirect
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0 // indirect
	k8s.io/klog/v2 v2.8.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	sigs.k8s.io/controller-runtime v0.8.1 // indirect
	sigs.k8s.io/wg-policy-prototypes v0.0.0-20210614220510-051e95e52a33
)

replace (
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
