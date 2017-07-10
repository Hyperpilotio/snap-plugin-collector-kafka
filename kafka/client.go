package kafka

import (
	"encoding/xml"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap/control/plugin"
	"io"
	"io/ioutil"
)

// const defines constant varaibles
const (
	MetricQuery    = "/serverbydomain?querynames=*:*&template=identity"
	MbeanQuery     = "/mbean?objectname="
	QuerySuffix    = "&template=identity"
	JavaStringType = "java.lang.String"
)

// XMLServer represents Server element
type XMLServer struct {
	XMLName xml.Name  `xml:"Server"`
	Domain  XMLDomain `xml:"Domain"`
}

// XMLDomain represents Domain element
type XMLDomain struct {
	XMLName xml.Name   `xml:"Domain"`
	MBeans  []XMLMBean `xml:"MBean"`
}

// XMLMBean represents MBean element
type XMLMBean struct {
	XMLName    xml.Name `xml:"MBean"`
	Objectname string   `xml:"objectname,attr"`
}

//XMLAttributes represents list of Attribute elements
type XMLAttributes struct {
	XMLName    xml.Name       `xml:"MBean"`
	Attributes []XMLAttribute `xml:"Attribute"`
}

// XMLAttribute represents Attribute element
type XMLAttribute struct {
	XMLName xml.Name `xml:"Attribute"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Value   float64  `xml:"value,attr"`
}

// Mx4jClient defines the URL of Kafka
type Mx4jClient struct {
	client *HTTPClient
	host   string
	Root   *node
}

// NewMx4jClient returns a new instance of Mx4jClient
func NewMx4jClient(url, host string) *Mx4jClient {
	return &Mx4jClient{
		client: NewHTTPClient(url, "", DefaultTimeout),
		host:   host,
		Root:   &node{Name: Root, Children: map[string]*node{}},
	}
}

// NewEmptyMx4jClient returns an empty instance of Mx4jClient
func NewEmptyMx4jClient() *Mx4jClient {
	return &Mx4jClient{}
}

// getMetricType returns all available metric types.
func (cc *Mx4jClient) getMetricType(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	return cc.buildMetricType(cfg)
}

// buildMetricType builds all metric types and write them into
func (cc *Mx4jClient) buildMetricType(cfg plugin.ConfigType) ([]plugin.MetricType, error) {

	cc, err := initClient(cfg)
	if err != nil {
		return nil, err
	}
	resp, err := cc.client.httpClient.Get(cc.client.GetUrl() + MetricQuery)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	mbeans, err := readObjectname(resp.Body)
	if err != nil {
		return nil, err
	}

	nspace := map[string]plugin.MetricType{}
	for _, mbean := range mbeans {
		// mbean.Objectname represents each callable measurement
		ns, _ := cc.getElementTypes(mbean.Objectname)
		for _, n := range ns {
			nspace[n.Namespace().String()] = n
		}
	}

	mtsType := []plugin.MetricType{}
	for _, v := range nspace {
		mtsType = append(mtsType, v)
	}

	// writeMetricTypes(mtsType)
	return mtsType, nil
}

// buildMetricAPI builds the base searchable tree and write it
// into CassandraMetricAPI.json file.
func (cc *Mx4jClient) buidMetricAPI() error {
	resp, err := cc.client.httpClient.Get(cc.client.GetUrl() + MetricQuery)
	if err != nil {
		return err
	}

	mbeans, err := readObjectname(resp.Body)
	if err != nil {
		return err
	}

	for _, mbean := range mbeans {
		nodes := makeLitteralNamespace(mbean.Objectname, "")
		cc.Root.Add(nodes, 0, mbean.Objectname)
	}
	// writeMetricAPIs(cc.Root)
	return nil
}

// getElementTypes returns specific XML element namespace along with its unit
func (cc *Mx4jClient) getElementTypes(url string) ([]plugin.MetricType, error) {
	resp, err := cc.client.httpClient.Get(cc.client.GetUrl() + MbeanQuery + url + QuerySuffix)
	if err != nil {
		kafkaLog.WithFields(log.Fields{
			"_block": "getTypes",
			"error":  err,
		}).Error(ReadDocErr)
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil || string(contents) == EmptyRespErr {
		kafkaLog.WithFields(log.Fields{
			"_block": "getTypes",
			"error":  err,
		}).Error(QueryDocErr)
		return nil, errors.New(QueryDocErr)
	}

	attrs, _ := readXMLAttributes(contents)
	ns := []plugin.MetricType{}
	for _, attr := range attrs {
		if attr.Type != JavaStringType {
			ns = append(ns, plugin.MetricType{
				Namespace_: makeDynamicNamespace(cc.host, url, attr.Name),
				Unit_:      attr.Type,
			})
		}
	}
	return ns, nil
}

func readObjectname(reader io.Reader) ([]XMLMBean, error) {
	var xmlServer XMLServer
	err := xml.NewDecoder(reader).Decode(&xmlServer)
	if err != nil {
		return nil, err
	}
	return xmlServer.Domain.MBeans, nil
}

func readXMLAttributes(content []byte) ([]XMLAttribute, error) {
	var xmlAttributes XMLAttributes
	xml.Unmarshal(content, &xmlAttributes)
	return xmlAttributes.Attributes, nil
}
