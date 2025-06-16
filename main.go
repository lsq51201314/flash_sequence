package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

type sliceInfo struct {
	CharacterID int `json:"id"`
	Width       int `json:"width"`
	Height      int `json:"height"`
	PositionX   int `json:"positionX"`
	PositionY   int `json:"positionY"`
	OffsetX     int `json:"offsetX"`
	OffsetY     int `json:"offsetY"`
}

func main() {
	//获取参数
	if len(os.Args) != 3 {
		fmt.Println("参数错误：请提供合并目录的完整路径以及序列名称。")
		return
	}
	dirs := os.Args[1]
	files := strings.ReplaceAll(os.Args[2], "[", "")
	files = strings.ReplaceAll(files, "]", "")
	arr := strings.Split(files, ",")
	//解析文件
	var W, H int
	items := make(map[string][]sliceInfo)
	imgs := make(map[string]image.Image)
	for _, v := range arr {
		jsonData, err := os.ReadFile(dirs + "\\" + v + ".json")
		if err != nil {
			fmt.Println("配置错误：" + err.Error())
			return
		}
		pngData, err := os.ReadFile(dirs + "\\" + v + ".png")
		if err != nil {
			fmt.Println("序列错误：" + err.Error())
			return
		}
		var slice []sliceInfo
		if err := json.Unmarshal(jsonData, &slice); err != nil {
			fmt.Println("解析错误：" + err.Error())
			return
		}
		buf := bytes.NewBuffer(pngData)
		img, err := png.Decode(buf)
		if err != nil {
			fmt.Println("图片错误：" + err.Error())
			return
		}
		//更改纵向位置
		for i, im := range slice {
			im.PositionY += H
			slice[i] = im
		}
		//计算最大宽高
		if img.Bounds().Max.X > W {
			W = img.Bounds().Max.X
		}
		H += img.Bounds().Max.Y
		items[v] = slice
		imgs[v] = img
	}
	//绘制图像
	dc := gg.NewContext(W, H)
	for k, v := range items {
		dc.DrawImage(imgs[k], 0, v[0].PositionY)
	}
	//保存图像
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		fmt.Println("转换错误：" + err.Error())
		return
	}
	info, err := json.Marshal(&items)
	if err != nil {
		fmt.Println("索引错误：" + err.Error())
		return
	}
	if err := os.WriteFile(dirs+".json", info, 0777); err != nil {
		fmt.Println("配置错误：" + err.Error())
		return
	}
	if err := os.WriteFile(dirs+".png", buf.Bytes(), 0777); err != nil {
		fmt.Println("保存错误：" + err.Error())
		return
	}
	fmt.Println("转换成功：" + dirs + ".png")
}
