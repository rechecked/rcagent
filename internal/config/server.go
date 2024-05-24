package config

type Endpoint func(cv Values) interface{}

var Endpoints = make(map[string]Endpoint)
