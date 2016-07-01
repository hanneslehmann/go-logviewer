package main

import (
    "log"
    "fmt"
    "encoding/json"
    "path/filepath"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "github.com/ActiveState/tail"
)

var FileURLtesting = map[string]FileList {
   "AccessRatesService":   FileList {Directory: "AccessRate", Path:"/tmp/*.*"},
   "AWBStockService": FileList {Directory: "AWBStockService", Path:"/home/hannes/*.*"}}

var basePath="/data/logs/applications/"
var FileURLBuild  = map[string]FileList {
      "AccessRatesService":   FileList {Directory: "AccessRate", Path: basePath+"AccessRatesService/*.*"},
      "AWBStockService": FileList {Directory: "AWBStockService", Path: basePath+"AWBStockService/*.*"}}

var FileURL = FileURLtesting
var _ = FileURLtesting
var _ = FileURLBuild

type FileList struct {
  Directory string
  Path string
  URL string
  Files FileListStruct
}

type FileListStruct struct {
    Files []FileStruct  `json:"files"`
}

type FileStruct struct {
  Name string `json:"name"`
  Size string `json:"size"`
  Mod string `json:"modified"`
}

func createFileList(fl FileList) FileListStruct {
  var returnList FileListStruct
  files, _ := filepath.Glob(fl.Path)
  for _,file := range files {
     returnList.Files=append(returnList.Files, FileStruct{filepath.Base(file),"",""})
  }
  return returnList
}

func getFileIndex(i FileList) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        b, _ := json.Marshal(createFileList(i))
        w.Write(b)
    }
}

func getFile(u string) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        http.ServeFile(w, r, filepath.Dir(FileURL[u].Path)+"/"+ps.ByName("file"))
    }
}

func getFileList(fl map[string]FileList) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        var list []string
        for key,_ := range fl {
           list=append(list,key)
        }
        b, _ := json.Marshal(struct{F[]string `json:"folders"`}{list})
        w.Write(b)
    }
}


func getFileStream(u string) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        t, _ := tail.TailFile(filepath.Dir(FileURL[u].Path)+"/"+ps.ByName("file"), tail.Config{Follow: true})
        for line := range t.Lines {
          fmt.Fprintln(w, line.Text)
          if f, ok := w.(http.Flusher); ok {
            f.Flush()
          } else {
            log.Println("Damn, no flush");
          }
        }
    }
}

func setupRoutes(router *httprouter.Router) {
  router.GET("/api/logfiles", getFileList(FileURL))
  for k, j:= range FileURL {
    router.GET("/api/logfiles/"+k, getFileIndex(j))
    router.GET("/api/logfiles/"+k+"/:file", getFile(k))
    router.GET("/api/stream"+k+"/:file", getFileStream(k))
  }
}

func main() {
    router := httprouter.New()
    setupRoutes(router)
    log.Fatal(http.ListenAndServe(":9891", router))
}
