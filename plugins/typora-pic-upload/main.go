package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"typoraPicUpload/pkg"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please input url or picture path")
		return
	}
	err := pkg.Login()
	if err != nil {
		fmt.Printf("登陆失败; INFO:%+v", err)
		return
	}
	for i := 1; i < len(os.Args); i++ {
		filePath := strings.Replace(os.Args[i], "\n", "", -1)
		if strings.HasPrefix(filePath, "http") {
			var mp map[string]interface{}
			func() {
				resp, err := pkg.UpLoadPic("url", filePath)
				defer resp.Body.Close()
				data, _ := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				mp = make(map[string]interface{})
				err = json.Unmarshal(data, &mp)
				if err != nil {
					log.Fatal(err)
				}
			}()
			if !mp["flag"].(bool) {
				log.Fatal(errors.New("转存外链图片失败"))
			} else {
				fmt.Println(mp["data"])
			}

		} else {
			var mp map[string]interface{}
			func() {
				resp, err := pkg.UpLoadPic("file", filePath)
				defer resp.Body.Close()
				data, _ := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				mp = make(map[string]interface{})
				err = json.Unmarshal(data, &mp)
				if err != nil {
					log.Fatal(err)
				}
			}()
			if !mp["flag"].(bool) {
				log.Fatal(errors.New("上传图片失败"))
			} else {
				fmt.Println(mp["data"])
			}
		}

	}
}
