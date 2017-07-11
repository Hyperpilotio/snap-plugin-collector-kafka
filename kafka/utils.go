package kafka

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/core"
	"net"
	"strings"
)

// const defines constant varaibles
const (
	EmptyRespErr        = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	ReadDocErr          = "Read document error"
	QueryDocErr         = "Queried document not found"
	EmptyNamespaceErr   = "To be collected metric namespace is empty"
	InvalidNamespaceErr = "To be collected metric namespace is invalid"

	Dot        = "."
	Underscore = "_"
	Root       = "Root"
)

var (
	kafkaLog = log.WithField("_module", "kafka-collector-client")
)

func initClient(cfg interface{}) (*Mx4jClient, error) {
	items, err := config.GetConfigItems(cfg, Mx4jURL, Mx4jPORT)
	if err != nil {
		return nil, err
	}

	url := items[Mx4jURL].(string)
	port := items[Mx4jPORT].(int)
	hostname, err := net.LookupAddr(url)
	if err != nil {
		hostname = []string{url}
	}
	server := fmt.Sprintf("%s:%d", url, port)
	fmt.Println("PrintingURL")
	fmt.Print(url, port, hostname)
	return NewMx4jClient(server, hostname[0]), nil
}

// TODO: Cache results
// func readMetricType() ([]plugin.MetricType, error) {
// 	data, err := Asset("data/CassandraMetricType.json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(data) == 0 {
// 		return nil, errors.New(ReadDocErr)
// 	}
// 	var metricTypes []plugin.MetricType
// 	err = json.Unmarshal(data, &metricTypes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return metricTypes, nil
// }

// TODO: Cache results
// func writeMetricTypes(types []plugin.MetricType) error {
// 	tys, err := json.Marshal(types)
// 	if err != nil {
// 		return err
// 	}

// 	jsonFile, err := os.Create("data/CassandraMetricType.json")
// 	if err != nil {
// 		return err
// 	}
// 	defer jsonFile.Close()

// 	jsonFile.Write(tys)
// 	jsonFile.Close()
// 	return nil
// }

// TODO: Cache results
// func writeMetricAPIs(node *node) error {
// 	tree, err := json.MarshalIndent(node, "", " ")
// 	if err != nil {
// 		return err
// 	}

// 	jsonFile, err := os.Create("data/CassandraMetricAPI.json")
// 	if err != nil {
// 		return err
// 	}
// 	defer jsonFile.Close()

// 	jsonFile.Write(tree)
// 	jsonFile.Close()
// 	return nil
// }

// TODO: Cache results
// func readMetricAPI() (*node, error) {
// 	content, err := Asset("data/CassandraMetricAPI.json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(content) == 0 {
// 		return nil, errors.New(ReadDocErr)
// 	}
// 	var jtree *node
// 	err = json.Unmarshal(content, &jtree)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jtree, nil
// }

func replaceDotToUnderscore(s string) string {
	if strings.Contains(s, Dot) {
		return strings.Replace(s, Dot, Underscore, -1)
	}
	return s
}

func replaceUnderscoreToDot(s string) string {
	if strings.Contains(s, Underscore) {
		return strings.Replace(s, Underscore, Dot, -1)
	}
	return s
}

// makeLitteralNamespace returns a string array without any modification
// to input string
func makeLitteralNamespace(url, name string) []string {
	ns := []string{}

	sp := strings.Split(url, ":")
	ns = append(ns, sp[0])

	sp1 := strings.Split(sp[1], ",")
	for _, s := range sp1 {
		v := strings.Split(s, "=")
		ns = append(ns, v...)
	}

	if len(name) > 0 {
		ns = append(ns, name)
	}
	return ns
}

// makeDynamicNamespace returns a dynamic namespace
func makeDynamicNamespace(host, url, name string) core.Namespace {
	ns := core.NewNamespace("hyperpilot", "kafka", "node").AddDynamicElement("nodeName", "The name of a Kafka node")

	sp := strings.Split(replaceDotToUnderscore(url), ":")
	ns = ns.AddStaticElement(sp[0])

	sp1 := strings.Split(sp[1], ",")
	for _, s := range sp1 {
		v := strings.Split(s, "=")
		ns = ns.AddStaticElement(v[0])
		ns = ns.AddDynamicElement(v[0]+" value", "The value of "+v[0])
	}

	if name != "" {
		ns = ns.AddStaticElement(name)
	}
	return ns
}
