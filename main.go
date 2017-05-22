// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Copyright 2017 Jay Kiline <jay@slushpupie.com>

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

func main() {

	getMainEngine().Run(":4000")
}

func getMainEngine() *gin.Engine {

	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")
	//r.LoadHTMLGlob("templates/*.json")

	r.GET("/install", installClient)
	r.GET("/uninstall", uninstallClient)
	r.GET("/clients", listClients)
	r.GET("/", indexPage)
	r.POST("/", generateInstall)

	return r
}

type oauthClient struct {
	OauthID         string `json:"oauthId"`
	CapabilitiesURL string `json:"capabilitiesUrl"`
	RoomID          int    `json:"roomId"`
	OauthSecret     string `json:"oauthSecret"`
}

var oauthClients = []oauthClient{}

func uninstallClient(c *gin.Context) {
	log.Print("Remove client request")
	redirectURL := c.Query("redirect_url")
	//installableURL := c.Query("installable_url")

	//meh..
	c.Redirect(http.StatusFound, redirectURL)
}

func installClient(c *gin.Context) {
	log.Print("New client request")

	redirectURL := c.Query("redirect_url")
	installableURL := c.Query("installable_url")

	req, err := http.NewRequest("GET", installableURL, nil)
	if err != nil {
		log.Print("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("Do: ", err)
		return
	}

	defer resp.Body.Close()

	var oclient oauthClient
	if err := json.NewDecoder(resp.Body).Decode(&oclient); err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	oauthClients = append(oauthClients, oclient)

	log.Print(oclient)
	c.Redirect(http.StatusFound, redirectURL)
}

func listClients(c *gin.Context) {
	c.JSON(200, gin.H{"clients": oauthClients})
}

func indexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func generateInstall(c *gin.Context) {

	data := struct {
		Host          string
		HipChatServer string
		Name          string
		Description   string
		Key           string
		AvatarURL     string
		Avatar2xURL   string
		FromName      string
		Scopes        []string
		AllowRoom     bool
		AllowGlobal   bool
		VendorName    string
		VendorURL     string
	}{
		Host:          c.Request.Host,
		HipChatServer: template.JSEscapeString(c.PostForm("server_name")),
		Name:          template.JSEscapeString(c.PostForm("name")),
		Description:   template.JSEscapeString(c.PostForm("description")),
		Key:           template.JSEscapeString(c.PostForm("key")),
		AvatarURL:     template.JSEscapeString(c.PostForm("avatar_url")),
		Avatar2xURL:   template.JSEscapeString(c.PostForm("avatar2x_url")),
		FromName:      template.JSEscapeString(c.PostForm("from_name")),
		Scopes:        c.PostFormArray("scopes"),
		AllowRoom:     c.PostForm("allow_room") == "true",
		AllowGlobal:   c.PostForm("allow_global") == "true",
		VendorName:    template.JSEscapeString(c.PostForm("vendor_name")),
		VendorURL:     template.JSEscapeString(c.PostForm("vendor_url")),
	}

	t := template.New("capabilities.json")
	t, err := t.ParseFiles("templates/capabilities.json")
	if err != nil {
		log.Print(err)
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Print(err)
	}

	jsonString := tpl.String()

	c.HTML(200, "index.html", gin.H{
		"Values":     data,
		"Host":       c.Request.Host,
		"InstallURL": fmt.Sprintf("https://%s/addons/install?url=data:application/json;base64,%s", template.URLQueryEscaper(c.PostForm("server_name")), base64.StdEncoding.EncodeToString([]byte(jsonString))),
		"Json":       jsonString,
		"Data":       base64.StdEncoding.EncodeToString([]byte(jsonString)),
	})
}
