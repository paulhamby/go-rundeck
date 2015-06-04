package rundeck

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
	//"log"
)

func TestExecutionOutput(t *testing.T) {
	xmlfile, err := os.Open("assets/test/execution.xml")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer xmlfile.Close()
	xmlData, _ := ioutil.ReadAll(xmlfile)
	var s ExecutionOutput
	xml.Unmarshal(xmlData, &s)
	/*
	<output>
	<id>23</id>
	<offset>424</offset>
	<completed>true</completed>
	<execCompleted>true</execCompleted>
	<hasFailedNodes>false</hasFailedNodes>
	<execState>succeeded</execState>
	<lastModified>1433212362000</lastModified>
	<execDuration>409</execDuration>
	<percentLoaded>98.6046511627907</percentLoaded>
	<totalSize>430</totalSize>
	<entries>
	<entry time='02:32:42' absolute_time='2015-06-02T02:32:42Z' log='hello' level='NORMAL' user='rundeck' command='' stepctx='1' node='localhost' />
	</entries>
	</output>
	*/
	//log.Printf("%+v\n", s.Entries.Entry)
	intexpects(s.ID, 23, t)
	intexpects(s.Offset, 424, t)
	strexpects(s.ExecState, "succeeded", t)
	intexpects(s.ExecDuration, 409, t)
	//expects(s.PercentLoaded, 98.6046511627907, t)
	intexpects(s.TotalSize, 430, t)
	strexpects(s.Entries.Entry[0].Log, "hello", t)
	strexpects(s.Entries.Entry[1].Log, "bye", t)
}
