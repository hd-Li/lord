package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	//managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	//"k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/norman/types/convert"
)

type employ struct {
	student
	name string
}

type student struct {
	level string
	age string
}

func main(){
	printSchema()
}

func reflTest(){
	t := reflect.TypeOf(v3.Cluster{})
	fmt.Println(t.Name())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Println(".......")
		fmt.Println("json name is : ", strings.SplitN(field.Tag.Get("json"), ",", 2)[0])
		fmt.Println("field name is : ", field.Name)
		fmt.Println("filed name converted is :", convert.LowerTitle(field.Name))
		fmt.Println("is 匿名成员 or not: ", field.Anonymous)
		fmt.Println("pkgPath is :", field.PkgPath)
	}
}

func printSchema(){
	schemas := Schemas.Schemas()
	for k, v := range schemas {
		fmt.Println("the num is : ", k)
		if v.CanList(nil) == nil {
		}
		b, err := json.Marshal(*v)
	    if err != nil {
	        fmt.Sprintf("%+v", *v)
	    }
	    var out bytes.Buffer
	    err = json.Indent(&out, b, "", "    ")
	    if err != nil {
	        fmt.Sprintf("%+v", *v)
	    }
		fmt.Println(out.String())
		fmt.Println("...................................")
	}
	
}