// Author: Liyanjing
// Created: 2023-12-28
// Tool Description: report to html.

package HTML

import (
	pub "DepthInspection/api/PublicClass"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
)

// 定义HTML模板用于生成HTML报告
const temp = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>mysql-check-report</title>
    <style>
        body {
            font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
            font-size: 14px;
            line-height: 1.5;
            color: #333;
            background-color: #f5f5f5;
        }
		
        h2 {
            font-weight: bold;
            font-size: 28px;
            margin: 20px auto;
            text-align: center;
        }
        table {
            font-size: 14px;
			width:210mm;
            margin: auto auto 20px;
            border-collapse: collapse;
            border-spacing: 0;
            background-color: transparent;
        }
        th,
        td {
            padding: 8px;
            vertical-align: top;
            border-top: 1px solid #ddd;
        }
        th {
            font-weight: bold;
            text-align: left;
            background-color: #f9f9f9;
            border: 1px solid #ddd;
            color: #0074a3;
            background-color:#e5eefd;
            white-space: nowrap;
        }
        td:hover {
            background-color: #ddd;
        }
        .table-bordered {
            border: 1px solid #ddd;
            border-collapse: separate;
            border-left: 0;
            border-radius: 4px;
            overflow: hidden;
        }
        .table-bordered th,
        .table-bordered td {
            border-left: 1px solid #ddd;
        }
        .table-hover > tbody > tr:hover {
            background-color: #f5f5f5;
        }
        .table-striped > tbody > tr:nth-of-type(odd) {
            background-color: #f9f9f9;
        }
        .table-hover .table-striped > tbody > tr:hover {
            background-color: #e8e8e8;
        }
        .text-center {
            text-align: center;
        }
        .text-right {
            text-align: right;
        }
        .text-left {
            text-align: left;
        }
        .bold {
            font-weight: bold;
        }
        .float-right {
            float: right;
        }
        .generated-time {
            font-size: 12px;
            text-align: right;
			width:210mm;
            margin: auto auto 5px;
        }
		.generated-time>span {
            font-size: 12px;
            font-weight: bold;
        }
		.title-lever1 {
            font-size: 14px;
			font-weight: bold;
            text-align: left;			
			width:210mm;
			height:32px;
			display: flex;align-items: center; 
            margin: auto auto 5px;
			background-color: #e5eefd;
        }
		.title-lever2{
            font-size: 12px;
			font-weight: bold;
            color:#0074a3;
            text-align: left;
			width:210mm;
			display: flex;align-items: center; 
            margin: auto auto 5px;
            margin-top: 10px;
        }
		.panel-responsive {
            font-size: 12px;
            text-align: left;			
			width:208mm;
            margin: auto auto 0px;
			padding-left:10px;padding-top:10px;
            background-color: #fff;
            word-wrap: break-word;
        }
		.panel-responsive>span {
            font-size: 12px;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="row">
            <div class="col-md-12">
                <h2>MySQL数据库巡检报告</h2>
                    <div class="generated-time"><span>巡检级别：</span>{{.ReportHead.InspectionLevel}}  ，<span>巡检人员：</span>{{.ReportHead.InspectionPersonnel}} ，<span>巡检时间：</span>{{.ReportHead.InspectionTime}}</div>
                <div class="title-lever1">一、巡检结果概览 </div>
                <div class="table-responsive">
                    <table class="table table-bordered table-hover table-striped">
                        <thead>
                            <tr>
                                <th class="text-center">序号</th>
                                <th class="text-center">检测项</th>
                                <th class="text-center">检测数量</th>
                                <th class="text-center">正常</th>
                                <th class="text-center">异常</th>

                            </tr>
                        </thead>
                        <tbody>
                                {{range .resultSummary}}
                                <tr>
                                    <td class="text-center">{{ .Id}}</td>        
                                    <td class="text-center">{{ .Name}}</td>        
                                    <td class="text-center">{{ .Counts}}</td>
                                    <td class="text-center">{{ .NormalCounts}}</td>
                                    <td class="text-center">{{ .AbnormalCounts}}</td>
                                </tr>
                                {{end}}
                        </tbody>
                    </table>
                </div>
                <div class="title-lever1">二、巡检结果详情 </div>
                {{range .checkResult}}
                    <div class="title-lever2">2. {{.Id}} {{.Title}} </div>
                    {{range .Results}}
                        <div class="panel-responsive">
                            <span style="width:70px;"> {{.Id}}. 巡检项: </span> {{.Name}} ，<span>阈值: </span> {{.Threshold}} ，<span>错误码：</span>{{.ErrorCode}}，
                            <span>异常相关信息：</span>{{.AbnormalInformation}}
                        </div>
                     {{end}}
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>
`

type OutputWayStruct struct{}
type OutPutWayInter interface {
	OutHtml()
	ResultSummaryStringSlice() [][]string
}

// 定义巡检介绍
type ReportHead struct {
	InspectionLevel     string
	InspectionPersonnel string
	InspectionTime      string
}

// 定义巡检结果概要
type ReportSummary struct {
	Id             string
	Name           string
	Counts         string
	NormalCounts   string
	AbnormalCounts string
}

// 定义巡检结果详情
type ReportResults struct {
	Id                  string
	Name                string
	Threshold           string
	ErrorCode           string
	AbnormalInformation string
}

// 定义对象，检测类别
type CheckType struct {
	Id         int
	TypeName   string
	TypeSliect []string
}

// 定义对象，检查结果
type CheckResult struct {
	Id      int
	Title   string
	Results []ReportResults
}

func (out *OutputWayStruct) OutHtml() {

	cst := time.FixedZone("CST", 8*60*60)
	currentTime := time.Now().In(cst)
	resultFileName := fmt.Sprintf(pub.ResultOutput.OutputPath+"mysql_check_report-%s.html", currentTime.Format("20060102-150405"))

	// 构建实例（设置数据）-巡检介绍
	reportHead := ReportHead{
		InspectionLevel:     pub.ResultOutput.InspectionLevel,
		InspectionPersonnel: pub.ResultOutput.InspectionPersonnel,
		InspectionTime:      pub.CheckBeginTime,
	}

	// 1.设置数据-巡检结果概要（标题1）
	reportSummaryArray := []ReportSummary{}
	for _, value := range out.ResultSummaryStringSlice() {
		//fmt.Println(fmt.Sprintf("name is %s", value))
		record := ReportSummary{value[0], value[1], value[2], value[3], value[4]}
		reportSummaryArray = append(reportSummaryArray, record)
	}

	// 2.设置数据-巡检结果
	arrayCheckType := []CheckType{
		{1, "数据库配置", []string{"configParameter"}},
		{2, "数据库性能", []string{"binlogDiskUsageRate", "historyConnectionMaxUsageRate", "tmpDiskTableUsageRate",
			"tmpDiskfileUsageRate", "innodbBufferPoolUsageRate", "innodbBufferPoolDirtyPagesRate", "innodbBufferPoolHitRate",
			"openFileUsageRate", "openTableCacheUsageRate", "openTableCacheOverflowsUsageRate", "selectScanUsageRate", "selectfullJoinScanUsageRate",
			"tableAutoPrimaryKeyUsageRate", "tableRows", "diskFragmentationRate", "bigTable", "coldTable"}},
		{3, "数据库基线", []string{"tableCharset", "tableEngine", "tableForeign", "tableNoPrimaryKey", "tableAutoIncrement",
			"tableBigColumns", "indexColumnIsNull", "indexColumnType", "tableIncludeRepeatIndex", "tableProcedureFunc", "tableTrigger"}},
		{4, "数据库安全", []string{"anonymousUsers", "emptyPasswordUser", "rootUserRemoteLogin", "normalUserConnectionUnlimited",
			"userPasswordSame", "normalUserDatabaseAllPrivilages", "normalUserSuperPrivilages", "databasePort"}},
	}

	checkResultArray := []CheckResult{}
	// 遍历检查类型，获取检查结果，封装实例
	for _, value := range arrayCheckType {
		//fmt.Println(value.TypeName)
		tmpTypeResult := out.tmpResultSummary(value.TypeSliect)
		tmpTypeResultArray := []ReportResults{}
		for _, vr := range tmpTypeResult {
			tempRecord := ReportResults{vr[0], vr[1], vr[2], vr[3], vr[4]}
			tmpTypeResultArray = append(tmpTypeResultArray, tempRecord)
		}
		tmpResult := CheckResult{value.Id, value.TypeName, tmpTypeResultArray}
		checkResultArray = append(checkResultArray, tmpResult)
	}

	// 3.传给html模板的数据
	content := map[string]interface{}{"ReportHead": reportHead, "resultSummary": reportSummaryArray, "checkResult": checkResultArray}

	//传给html模板的数据
	//content := map[string]interface{}{"ReportHead": reportHead, "resultSummary": reportSummaryArray, "resultConfig": resultConfigArray, "resultPerformance": resultPerformanceArray}

	// 创建并写入HTML报告
	var report = template.Must(template.New(resultFileName).Parse(temp))
	file, err := os.Create(resultFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	report.Execute(file, content)
	fmt.Println(fmt.Sprintf("Output written to file %s", resultFileName))

}
