package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//douyu 斗鱼直播
type douyu struct{}

//Site 实现接口
func (i *douyu) Site() string { return "斗鱼直播" }

//SiteURL 实现接口
func (i *douyu) SiteURL() string {
	return "http://www.douyu.com"
}

//SiteIcon 实现接口
func (i *douyu) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *douyu) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *douyu) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *douyu) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	html, err := httpGet(url)
	if strings.Contains(html, "data-onlineid=") {
		reg, _ := regexp.Compile("data-onlineid=(\\d+)")
		id = reg.FindStringSubmatch(html)[1]
		url = fmt.Sprintf("http://www.douyu.com/ztCache/WebM/room/%s", id)
		html, err = httpGet(url)
		if err == nil && html != "[]" {
			live = strings.Contains(html, "\\\"show_status\\\":1")
		} else {
			id = ""
		}
	} else {
		reg, _ := regexp.Compile("\\\\?\"room_id\\\\?\":(\\d+)")
		tmp := reg.FindString(html)
		live = !strings.Contains(html, "上次直播")
		id = reg.FindStringSubmatch(tmp)[1]
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *douyu) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://www.douyutv.com/api/v1/"
	args := fmt.Sprintf("room/%s?aid=wp&client_sys=wp&time=%d", id, getUnixTimesTamp())
	url = fmt.Sprintf("%s%s&auth=%s", url, args, getMD5String(args+"zNzMV1y4EMxOHS6I5WKm"))
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).JToken("data"))
	video := fmt.Sprintf("%s/%s", json["rtmp_url"], json["rtmp_live"])
	img := json["room_src"].(string)
	title := json["room_name"].(string)
	details := json["show_details"].(string)
	nick := json["nickname"].(string)
	live.LiveNick = nick
	live.LivingIMG = img
	live.RoomDetails = details
	live.RoomTitle = title
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
