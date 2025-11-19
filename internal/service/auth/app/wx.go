package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/jwt"
)

const (
	appID     = "wx123456789" // 微信后台获取
	appSecret = "abcdef"      // 微信后台获取
)

type WxResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func Code2Session(ctx context.Context, code string) (*WxResp, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var wx WxResp
	if err := json.Unmarshal(body, &wx); err != nil {
		return nil, err
	}
	if wx.ErrCode != 0 {
		return nil, fmt.Errorf("wx err:%d %s", wx.ErrCode, wx.ErrMsg)
	}
	return &wx, nil
}

func GenTokens(openID string) (acc, ref string, err error) {
	now := time.Now()
	acc, err = jwt.Sign(jwt.Claims{
		Sub: openID,
		Exp: now.Add(30 * time.Minute).Unix(),
	})
	if err != nil {
		return "", "", err
	}
	ref, err = jwt.Sign(jwt.Claims{
		Sub: openID,
		Exp: now.Add(7 * 24 * time.Hour).Unix(),
	})
	return acc, ref, err
}
