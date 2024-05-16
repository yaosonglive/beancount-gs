package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/beancount-gs/script"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	CORPID            = "" // "APPID"                                              // 填上自己的参数
	AGENTID           = "" // "AGENTID"                                            // 填上自己的参数
	AgentSecret       = "" // "APPSECRET"                                          // 填上自己的参数
	oauth2RedirectURI = "" // "http://192.168.1.129:8080/api/work_weixin/callback" // 填上自己的参数
	oauth2Scope       = "" // "snsapi_base"                                        // 填上自己的参数
	LedgerName        = "" // "snsapi_base"                                        // 填上自己的参数
	LedgerSecret      = "" // "snsapi_base"                                        // 填上自己的参数
)

func InitWorkWeixin() {
	config := script.GetServerConfig()
	CORPID = config.CORPID
	AGENTID = config.AGENTID
	AgentSecret = config.AgentSecret
	oauth2RedirectURI = config.DOMAIN + "/api/oauth2/work_weixin/callback"
	LedgerName = config.LedgerName
	LedgerSecret = config.LedgerSecret
}

type WorkWeixinAccessToken struct {
	Token  string
	Expire int64
}

var workWeixinAccessToken = &WorkWeixinAccessToken{}

var (

// sessionStorage                 = session.New(20*60, 60*60)
// oauth2Endpoint oauth2.Endpoint = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
)

func GetAccessToken() string {
	if workWeixinAccessToken.Expire > time.Now().Unix() {
		return workWeixinAccessToken.Token
	}
	url := fmt.Sprintf("https://wiki.yaosong.live/cgi-bin/gettoken?corpid=%s&corpsecret=%s", CORPID, AgentSecret) //qyapi.weixin.qq.com
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	res := struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	respData, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respData, &res)
	if res.Errcode == 0 && res.AccessToken != "" {
		workWeixinAccessToken.Token = res.AccessToken
		workWeixinAccessToken.Expire = time.Now().Unix() + int64(res.ExpiresIn) - 10
	}
	return workWeixinAccessToken.Token
}

func Oauth2WorkWeixinCallback(c *gin.Context) {
	code, _ := c.GetQuery("code")
	token := GetAccessToken()

	url := fmt.Sprintf("https://wiki.yaosong.live/cgi-bin/auth/getuserinfo?access_token=%s&code=%s", token, code)
	resp, err := http.Get(url)
	if err != nil {
		InternalError(c, "微信接口异常")
		return
	}
	res := struct {
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
		Userid     string `json:"userid"`
		UserTicket string `json:"user_ticket"`
	}{}

	respData, _ := ioutil.ReadAll(resp.Body)
	script.LogSystemInfo("Oauth2WorkWeixinCallback " + string(respData))
	err = json.Unmarshal(respData, &res)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if res.UserTicket == "" {
		c.String(200, "没有访问权限")
		return
	}
	t := sha1.New()
	_, err = io.WriteString(t, LedgerName+LedgerSecret)
	if err != nil {
		LedgerIsNotAllowAccess(c)
		return
	}

	ledgerId := hex.EncodeToString(t.Sum(nil))
	c.Data(200, "text/html; charset=utf-8", []byte("<script>window.localStorage.setItem('ledgerId', \""+ledgerId+"\");window.location.href=\"/\";</script>"))
	//c.String(200, "<script>window.localStorage.setItem('ledgerId', \""+ledgerId+"\");window.location.href=\"/\";</script>")
}

func Oauth2WorkWeixinCheck(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", []byte("<script>if(window.localStorage.getItem(\"ledgerId\")){window.location.href=\"/\";}else{window.location.href=\"/api/oauth2/work_weixin\";}</script>"))
}

func Oauth2WorkWeixin(c *gin.Context) {
	url := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_privateinfo&state=STATE&agentid=%s#wechat_redirect", CORPID, oauth2RedirectURI, AGENTID)

	script.LogSystemInfo("Oauth2WorkWeixin " + url)
	c.Redirect(302, url)
}
