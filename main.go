package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gorilla/mux"
)

var ossClient *oss.Client

var (
	ENDPOINT = ""
	listTemp *template.Template
	infoTemp *template.Template
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("need args: endpoint key secret")
		os.Exit(1)
		return
	}
	ENDPOINT = os.Args[1]

	// 创建OSSClient实例。
	client, err := oss.New(ENDPOINT, os.Args[2], os.Args[3])
	if err != nil {
		panic(err)
	}
	ossClient = client

	initTemplate()

	m := mux.NewRouter()
	m.HandleFunc("/", home)
	m.HandleFunc("/list/{bucket}", info)
	m.HandleFunc("/list/{bucket}/{name}", info)
	m.HandleFunc("/list/{bucket}/{name}/{action}", eprocess)
	http.ListenAndServe(":9090", m)

}
func initTemplate() {
	var err error
	listTemp, err = template.New("list").Parse(_lt)
	if err != nil {
		panic(err)
	}
	infoTemp, err = template.New("info").Parse(_it)
	if err != nil {
		panic(err)
	}

}

type ListModel struct {
	Path string
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.Write([]byte(fmt.Sprint(err)))
			w.Write(debug.Stack())
		}
	}()
	fmt.Println("Req home")
	// 列举存储空间。
	ls := []ListModel{}
	marker := ""
	for {
		lsRes, err := ossClient.ListBuckets(oss.Marker(marker))
		if err != nil {
			panic(err)
		}

		// 默认情况下一次返回100条记录。
		for _, bucket := range lsRes.Buckets {
			ls = append(ls, ListModel{
				Path: "/list/" + bucket.Name,
				Name: bucket.Name,
			})
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	if err := listTemp.Execute(w, ls); err != nil {
		panic(err)
	}
}

type InfoModel struct {
	Bucket string
	Path   string
	Rame   string
	Name   string
}

func info(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.Write([]byte(fmt.Sprint(err)))
			w.Write(debug.Stack())
		}
	}()
	vars := mux.Vars(r)
	fmt.Printf("Req info:%+v\n", vars)
	bucketName, ok := vars["bucket"]
	if !ok {
		w.Write([]byte("参数错误"))
		return
	}
	// 获取存储空间。
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	fname, ok := vars["name"]
	i := 0
	if ok {
		fname = strings.ReplaceAll(fname, "@", "/")
		i++
	}
	fmt.Println("ossReq:", fname)

	lsRes, err := bucket.ListObjects(oss.Marker(fname), oss.MaxKeys(2))
	if err != nil {
		panic(err)
	}

	if len(lsRes.Objects) == 0 {
		w.Write([]byte("无"))
		return
	}

	fmt.Printf("ossResp:%+v\n", lsRes.Objects)

	if err := infoTemp.Execute(w, InfoModel{
		Bucket: bucketName,
		Rame:   strings.ReplaceAll(lsRes.Objects[0].Key, "/", "@"),
		Name:   strings.ReplaceAll(lsRes.Objects[i].Key, "/", "@"),
		Path:   fmt.Sprintf("https://%s.%s/%s", bucketName, ENDPOINT, lsRes.Objects[0].Key),
	}); err != nil {
		panic(err)
	}
}

func eprocess(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.Write([]byte(fmt.Sprint(err)))
			w.Write(debug.Stack())
		}
	}()
	vars := mux.Vars(r)
	fmt.Printf("Req eprocess:%+v\n", vars)
	bucketName, ok := vars["bucket"]
	if !ok {
		w.Write([]byte("bucket 参数错误"))
		return
	}
	// 获取存储空间。
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	fname, ok := vars["name"]
	if !ok {
		w.Write([]byte("name 参数错误"))
		return
	}
	fname = strings.ReplaceAll(fname, "@", "/")

	action, ok := vars["action"]
	switch action {
	case "90", "180", "270":
		style := "image/rotate," + action
		process := fmt.Sprintf("%s|sys/saveas,o_%v", style, base64.URLEncoding.EncodeToString([]byte(fname)))
		result, err := bucket.ProcessObject(fname, process)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("处理完成:" + result.Status))
	case "l90", "l180", "l270":
		// 下载图片
		// 下载文件到本地文件。
		err = bucket.GetObjectToFile(fname, "tmp/"+vars["name"])
		if err != nil {
			panic(err)
		}
		w.Write([]byte("处理完成"))
	default:
		w.Write([]byte("name 参数错误"))
		return
	}

}
