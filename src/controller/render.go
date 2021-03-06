package controller

import (
	"container/list"
	"framework/base/config"
	"info"
	"sort"
	"strings"
	"sync"
	"time"
)

type tagRender struct {
	Tag   string
	Count int
}

type timeRender struct {
	Date  string
	Year  int
	Month int
}

type sideRender struct {
	BlogTagList  []*tagRender
	BlogTimeList []*timeRender
}

type hostRender struct {
	Host string
}

var staticHostRender *hostRender = nil
var hostRenderOnce sync.Once

func buildHostRender() *hostRender {
	hostRenderOnce.Do(func() {
		staticHostRender = &hostRender{}
		staticHostRender.Host = config.GetDefaultConfigJsonReader().Get("net.host").(string)
		if !strings.HasPrefix(staticHostRender.Host, "http://") {
			staticHostRender.Host = "http://" + staticHostRender.Host
		}
	})
	return staticHostRender
}

func buildSideRender(blogList *list.List) *sideRender {
	var topRender sideRender
	var tagMap map[string]int = make(map[string]int)
	var timeMap map[string]int64 = make(map[string]int64)
	for iter := blogList.Front(); iter != nil; iter = iter.Next() {
		info := iter.Value.(info.BlogInfo)
		for tag := range info.BlogTagList {
			tagMap[info.BlogTagList[tag]]++
		}
		timeMap[time.Unix(info.BlogTime, 0).Format("2006年01月")] = info.BlogTime
	}
	var tagList []*tagRender = nil
	for k, v := range tagMap {
		tagList = append(tagList, &tagRender{k, v})
	}
	var blogTimeList BlogTimeList = nil
	for k, v := range timeMap {
		blogTimeList = append(blogTimeList, &BlogTime{k, v})
	}
	sort.Sort(blogTimeList)
	var blogTimeStringList []*timeRender = nil
	for i := range blogTimeList {
		renderTime := time.Unix(blogTimeList[i].Time, 0)
		blogTimeStringList = append(blogTimeStringList, &timeRender{blogTimeList[i].Tag,
			renderTime.Year(), int(renderTime.Month())})
	}
	topRender.BlogTagList = tagList
	topRender.BlogTimeList = blogTimeStringList
	return &topRender
}
