package main

import "./adapter/rest"
import _ "./services/broker"

func main() {
	rest.ListenAndServe()
}
