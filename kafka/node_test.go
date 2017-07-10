package kafka

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"strings"
	"testing"
)

func TestReadXMLAttribute(t *testing.T) {
	Convey("ReadXMLAttributes shuould get all attributes within a Mbean.", t, func() {
		body := `<?xml version="1.0" encoding="UTF-8"?>
		<MBean classname="com.yammer.metrics.reporting.JmxReporter$Meter" description="Information on the management interface of the MBean" objectname="kafka.server:type=BrokerTopicMetrics,name=BytesInPerSec">
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Count" strinit="true" type="long" value="313358992"/>
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="EventType" strinit="true" type="java.lang.String" value="bytes"/>
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="FifteenMinuteRate" strinit="true" type="double" value="9892.210108862237"/>
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="FiveMinuteRate" strinit="true" type="double" value="176.8137594030284"/>
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="MeanRate" strinit="true" type="double" value="37602.763294272896"/>
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="OneMinuteRate" strinit="true" type="double" value="9.599225732359006E-11"/>
		    <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="RateUnit" strinit="false" type="java.util.concurrent.TimeUnit" value="SECONDS"/>
		    <Operation description="Operation exposed for management" impact="unknown" name="objectName" return="javax.management.ObjectName"/>
		</MBean>`
		content, err := ioutil.ReadAll(strings.NewReader(body))
		So(err, ShouldBeNil)
		attrs, err := readXMLAttributes(content)
		So(err, ShouldBeNil)
		So(len(attrs), ShouldEqual, 7)
	})

}
