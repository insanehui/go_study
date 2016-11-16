package main

import (
	"gopkg.in/yaml.v2"
	"fmt"
	// "github.com/ghodss/yaml" // 找到另一个yaml库，但貌似没有强大的功能
)

func main() {

	yfile := `
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: db2
  namespace: linksame-simplest
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
  selector:
    matchLabels:
      name: "db"
`
	var m map[string]interface{}

	yaml.Unmarshal([]byte(yfile), &m)

	spec := m["spec"]

	// 用类型断言，比reflect要简便一些
	if spec, ok := spec.(map[interface{}]interface{}); ok {
		sel := spec["selector"]
		if sel, ok := sel.(map[interface{}]interface{}); ok {
			mt := sel["matchLabels"]
			if mt, ok := mt.(map[interface{}]interface{}); ok {
				mt["xxx"] = 1234214
			}
		}
	}

	b, _ := yaml.Marshal(m)
	fmt.Printf("%s", b)
}
