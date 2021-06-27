## Policy Report octant plugin 
[Under development]

Resource Policy Report Tab 
![alt text](./img/octant.png)

Namespace Policy Report Tab
![alt text](./img/octant-namespace.png)

Policy Report Navigation 
![alt text](./img/octant-cluster.png)
## Installation

Install policyreport-octant-plugin
```bash
$ curl -s https://raw.githubusercontent.com/evalsocket/policyreport-octant-plugin/master/install.sh | bash
$ mv bin/policyreport ~/.config/octant/plugins/
$ If you don't have any policy engine installed in your cluster then i will suggest you to install one who uses policy report. 
$ octant
```

Build the plugin manually:
`go build -o policyreport-octant-plugin  cmd/policyrepor-octant-plugin/main.go `
 
Then move the binary:

`mv policyreport-octant-plugin ~/.config/octant/plugins/`

You may need to create this directory if it does not exist.

TODO:
- Add Kyverno policy validation in octant editor
- Falco ecosystem 
