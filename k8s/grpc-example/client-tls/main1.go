package main

import (
	"crypto/tls"
	//"crypto/x509"
	//"io/ioutil"
	"log"
)

func main() {
	//cert, err := ioutil.ReadFile("/run/certs/root.crt")
	//if err != nil {
	//	log.Fatalf("Couldn't load file", err)
	//}
	//certPool := x509.NewCertPool()
	//certPool.AppendCertsFromPEM(cert)

	conf := &tls.Config{
		//RootCAs: certPool,
        InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "localhost:5443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))
}
