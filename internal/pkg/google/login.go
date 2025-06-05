package google

import (
	"context"

	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

func VerifyToken(token string) (map[string]interface{}, error) {
	clientId := config.Conf.Server.Google.Id
	token, err := exchangeCode(token)
	if err != nil {
		return nil, err
	}
	log.Logger.Infof("google token: %s\n", token)
	payload, err := idtoken.Validate(context.Background(), token, clientId)
	if err != nil {
		return nil, err
	}
	claims := payload.Claims
	log.Logger.Infof("google返回信息：%+v\n", claims)
	claims["id"] = payload.Subject
	return claims, nil
}

func exchangeCode(code string) (string, error) {
	config := &oauth2.Config{
		ClientID:     config.Conf.Server.Google.Id,
		ClientSecret: config.Conf.Server.Google.Secret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://bmdev.gss.run",
		Scopes:       []string{"openid", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}
	return t.Extra("id_token").(string), nil
}
