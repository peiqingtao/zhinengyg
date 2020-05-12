package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// 脚手架
var jsonName, mName, mTitle, rName, cName string
var tConfig map[string]interface{}
var tFields []map[string]string

func main() {
	// # 1 获取需要生成脚手架的模型参数
	jsonName = os.Args[1]
	// 解析 JSON
	parseJSON()

	// 模型名
	mName = tConfig["modelName"].(string)
	// 标题
	mTitle = tConfig["title"].(string)

	// # 2 利用参数获取需要的数据
	//路由资源名称
	// Brand -> brand, UserRoom => user-room
	// 思路，将不处于开头的大写字母前增加-，转换为小写
	p := `([A-Z])`
	re, _ := regexp.Compile(p)
	// 正则替换 | 小写转换 | 去掉首字母的-
	rName = strings.ToLower(re.ReplaceAllString(mName, "-${1}"))[1:]
	// 控制器名字前缀
	cName = mName

	// 模型构建
	genModel()

	// # 3 生成路由代码
	genRouter()

	// # 4 生成控制器代码
	genController()

	// # 5 生成前端路由
	genManagerRouter()
	// # 6 生成前端组件
	genManagerCP()

}

func genManagerCP() {
	// 从模板读取内容
	tplFile := "./scaffoldTemplate/manager/listComponent"
	content, _ := ioutil.ReadFile(tplFile) // []byte
	code := string(content)
	// 占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%r-name%", rName)
	cpName := mName + "List"
	code = strings.ReplaceAll(code, "%cp-name%", cpName)
	// 展示列列表，筛选字段
	tcList := ""
	filterList := ""
	setFieldList := ""
	for _, field := range tFields { //map[name:Name type:string, isList:yes, label:]
		if field["isList"] == "yes" {
			// 需要列表展示的情况下
			sortAble := ""
			if field["isSort"] == "yes" {
				sortAble = `sortable="custom"`
			}
			tcList += fmt.Sprintf("\t\t\t\t\t\t\t\t<el-table-column prop=\"%s\" label=\"%s\" %s></el-table-column>\n",
				field["name"], field["label"], sortAble,
			)
		}

		if field["isFilter"] == "yes" {
			//需要筛选
			filterList += fmt.Sprintf(`
						<el-form-item label="%s">
							<el-input v-model="filterForm.filter%s"></el-input>
						</el-form-item>`,
						field["label"], field["name"],
			)
		}

		if field["isSet"] == "yes" {
			// 需要出现在表单中
			setFieldList += fmt.Sprintf(`
						<el-form-item label="%s" prop="%s">
							<el-input v-model="itemSetForm.%s"></el-input>
						</el-form-item>`,
				field["label"], field["name"], field["name"],
			)
		}

	}
	code = strings.ReplaceAll(code, "%tc-list%",tcList)
	code = strings.ReplaceAll(code, "%filter-list%",filterList)
	code = strings.ReplaceAll(code, "%set-field-list%",setFieldList)

	// 将 code 写入到指定文件
	cPath := strings.ToLower(mName[0:1]) + mName[1:]
	codeFile := "../manager/src/components/" + cPath + "/" + cpName + ".vue"
	// 子目录
	os.Mkdir("../manager/src/components/" + cPath, 0755)
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY | os.O_TRUNC | os.O_CREATE , 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("manager 生成的 component 代码位于 components/ 中")
}

// manager 端路由
func genManagerRouter() {
	code := `{ path: '%r-name%', component: ()=>import('../components/%c-path%/%c-name%List.vue'), },`
	// 占位符替换
	code = strings.ReplaceAll(code, "%r-name%", rName)
	cPath := strings.ToLower(mName[0:1]) + mName[1:]
	code = strings.ReplaceAll(code, "%c-path%", cPath)
	code = strings.ReplaceAll(code, "%c-name%", cName)

	codeFile := "./ManagerRouterCode"
	handle, _ := os.OpenFile(codeFile, os.O_APPEND | os.O_CREATE, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code + "\n")
	writer.Flush()

	log.Println("manager 生成的 router 代码位于 scaffold/ManagerRouterCode 中，请将代码拷贝到 router/router.js 中。")
}


func genController() {
	// 从模板读取内容
	tplFile := "./scaffoldTemplate/backend/controller"
	content, _ := ioutil.ReadFile(tplFile) // []byte
	code := string(content)
	// 占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%m-name%", mName)
	code = strings.ReplaceAll(code, "%c-name%", cName)
	// 将 code 写入到指定文件
	// controller file name 。ClassRoom => classRoom, Brand => brand
	cfName := strings.ToLower(mName[0:1]) + mName[1:]
	codeFile := "../backend/src/controller/" + cfName + ".go"
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY | os.O_TRUNC | os.O_CREATE , 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("backend 生成的 controller 代码位于 controller/ 中")
}

func genRouter() {
	// 从模板读取内容
	tplFile := "./scaffoldTemplate/backend/router"
	content, _ := ioutil.ReadFile(tplFile) // []byte
	code := string(content)
	// 占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%r-name%", rName)
	code = strings.ReplaceAll(code, "%c-name%", cName)
	// 将 code 写入到指定文件
	codeFile := "./routerCode"
	handle, _ := os.OpenFile(codeFile, os.O_APPEND | os.O_CREATE, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("backend 生成的 router 代码位于 scaffold/routerCode 中，请将代码拷贝到 router/router.go 中。")
}

func parseJSON() {
	// 读取json 文件
	file := "./config/" + jsonName +".json"
	content, _ := ioutil.ReadFile(file)
	json.Unmarshal(content, &tConfig)
	// 将tConfig处理成可以直接操作的结构
	fields := tConfig["fields"].([]interface{})
	for _, f := range fields {
		field :=  map[string]string{}
		for k, v := range f.(map[string]interface{}) {
			field[k] = v.(string)
		}
		tFields = append(tFields, field)
	}
}

func genModel() {
	// 从模板读取内容
	tplFile := "./scaffoldTemplate/backend/model"
	content, _ := ioutil.ReadFile(tplFile) // []byte
	code := string(content)
	// 占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%m-name%", mName)
	// 替换 field-list
	fieldList := ""
	for _, field := range tFields { //map[name:Name type:string]
		fieldList += fmt.Sprintf("\t%s %s\n", field["name"], field["type"])
	}
	code = strings.ReplaceAll(code, "%field-list%", fieldList)

	// 将 code 写入到指定文件
	codeFile := "../backend/src/model/" + mName + ".go"
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY | os.O_TRUNC | os.O_CREATE , 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("backend 生成的 model 代码位于 model/ 中")

}