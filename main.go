package main

import (
	"DepthInspection/InspectionItem"
	"DepthInspection/ResultOutput/HTML"
	"DepthInspection/ResultOutput/PDF"
	"DepthInspection/api/PublicClass"
	"DepthInspection/flag"
	"strings"
	"time"
)

var OutputPdf PDF.OutPutWayInter = &PDF.OutputWayStruct{}
var OutputHtml HTML.OutPutWayInter = &HTML.OutputWayStruct{}

func main() {
	//配置文件初始化
	PublicClass.CheckBeginTime = time.Now().Format("2006-01-02 15:04:05")
	flag.ParameterCheck()
	PublicClass.ConfigInit()
	PublicClass.YamlconfigInit()
	//查询数据初始化，将用到的数据初始化到内存中
	PublicClass.QueryDbDateInit()
	InspectionItem.DatabaseConfigCheck(PublicClass.ConfigurationCanCheck)
	var c = &InspectionItem.DatabaseBaselineCheckStruct{}
	c.BaselineCheckTablesDesign()
	c.BaselineCheckColumnsDesign()
	c.BaselineCheckIndexColumnDesign()
	c.BaselineCheckProcedureTriggerDesign()
	c.BaselineCheckUserPriDesign()
	c.BaselineCheckPortDesign()
	c.DatabasePerformanceStatusCheck()
	c.DatabasePerformanceTableIndexCheck()
	PublicClass.CheckEndTime = time.Now().Format("2006-01-02 15:04:05")
	PublicClass.CheckTimeConsuming, _ = PublicClass.Strea.GetTimeSecondsArr(PublicClass.CheckBeginTime, PublicClass.CheckEndTime)

	outPutFileType := strings.TrimSpace(PublicClass.ResultOutput.OutputWay)
	if outPutFileType == "" {
		OutputPdf.OutPdf()
		OutputHtml.OutHtml()
	} else if outPutFileType == "pdf" {
		OutputPdf.OutPdf()
	} else {
		OutputHtml.OutHtml()
	}

}
