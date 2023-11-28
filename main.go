package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	updateAccessTokenURL = "https://auth.aliyundrive.com/v2/account/token"
	signinURL            = "https://member.aliyundrive.com/v1/activity/sign_in_list"
)

// 阿里云盘使用AccessToken鉴权，RefreshToken用来获取AccessToken，这两个Token都有期限
// AccessToken好像是几个小时，RefreshToken不知道，好像好几天都不变
// 脚本原理就是使用用户输入的的RefreshToken获取AccessToken，然后访问对应的签到链接
var refreshTokenArray = []string{
	// 手动填写自己的refresh_token
	"5565f62eaf5040bf9961xxxxxxxxxxxx",
	"",
}

// 用来发起http请求，返回响应包
func makeRequest(url string, queryBody map[string]string) ([]byte, error) {
	body, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// 更新AccessToken
func updateAccessToken(refreshToken string) (string, error) {
	queryBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	respBody, err := makeRequest(updateAccessTokenURL, queryBody)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("获取accessToken错误")
	}

	return accessToken, nil
}

func signIn(accessToken string) error {
	queryBody := map[string]string{
		"grant_type":   "refresh_token",
		"access_token": accessToken,
	}

	respBody, err := makeRequest(signinURL, queryBody)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))
	return nil
}

func main() {
	for _, refreshToken := range refreshTokenArray {
		accessToken, err := updateAccessToken(refreshToken)
		if err != nil {
			fmt.Printf("更新accessToken错误: %v\n", err)
			continue
		}
		fmt.Println(accessToken)

		err = signIn(accessToken)
		if err != nil {
			fmt.Printf("登录错误: %v\n", err)
		}
	}
}
