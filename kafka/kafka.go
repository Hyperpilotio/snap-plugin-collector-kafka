package kafka

import (
	"reflect"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
)

// const defines constant varaibles
const (
	// Name of plugin
	Name = "kafka"
	// Version of plugin
	Version = 1
	// Type of plugin
	PluginType = plugin.CollectorPluginType

	// Timeout duration
	DefaultTimeout = 5 * time.Second

	Mx4jURL    = "mx4j_url"
	Mx4jPORT   = "mx4j_port"
	InvalidURL = "Invalid URL in Global configuration"
	NoHostname = "No hostname define in Global configuration"
)

// Meta returns the snap plug.PluginMeta type
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, PluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// NewKafkaCollector returns a new instance of Kafka struct
func NewKafkaCollector() *Kafka {
	return &Kafka{}
}

// Kafka struct
type Kafka struct {
	client *Mx4jClient
}

// CollectMetrics collects metrics from Kafka through JMX
func (p *Kafka) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}

	err := p.loadMetricAPI(mts[0].Config())
	if err != nil {
		return nil, err
	}

	for _, m := range mts {
		results := []nodeData{}
		search := strings.Split(replaceUnderscoreToDot(strings.TrimLeft(m.Namespace().String(), "/")), "/")
		if len(search) > 3 {
			p.client.Root.Get(p.client.client.GetUrl(), search[4:], 0, &results)
		}

		for _, result := range results {
			ns := append([]string{"hyperpilot", "kafka", "node", p.client.host}, strings.Split(result.Path, Slash)...)
			metrics = append(metrics, plugin.MetricType{
				Namespace_: core.NewNamespace(ns...),
				Timestamp_: time.Now(),
				Data_:      result.Data,
				Unit_:      reflect.TypeOf(result.Data).String(),
			})
		}

	}

	return metrics, nil
}

// GetMetricTypes returns the metric types exposed by Kafka
func (p *Kafka) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	return NewEmptyMx4jClient().getMetricType(cfg)
}

// GetConfigPolicy returns a ConfigPolicy
func (p *Kafka) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// loadMetricAPI returns the root node
func (p *Kafka) loadMetricAPI(config *cdata.ConfigDataNode) error {
	var err error
	// inits KafkaClient
	p.client, err = initClient(plugin.ConfigType{ConfigDataNode: config})
	if err != nil {
		return err
	}

	// reads the root metric node from the memory
	// nod, err := readMetricAPI()
	// if err != nil {
	err = p.client.buidMetricAPI()
	if err != nil {
		return err
	}
	return nil
}
