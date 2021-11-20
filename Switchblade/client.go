package main

import (
        "bytes"
        "crypto/tls"
        "crypto/x509"
        "io/ioutil"
        "log"
        "net/http"
)

func main() {
        // load client cert
        cert, err := tls.LoadX509KeyPair("/etc/nginx/certs/server.crt", "/etc/nginx/certs/server.key")
        if err != nil {
                log.Fatal(err)
        }

        // load CA cert
        caCert, err := ioutil.ReadFile("/etc/nginx/certs/ca.crt")
        if err != nil {
                log.Fatal(err)
        }
        caCertPool := x509.NewCertPool()
        caCertPool.AppendCertsFromPEM(caCert)

        // https client tls config
        // InsecureSkipVerify true means not validate server certificate (so no need to set RootCAs)
        tlsConfig := &tls.Config{
                Certificates: []tls.Certificate{cert},
                //RootCAs:            caCertPool,
                InsecureSkipVerify: true,
        }
        tlsConfig.BuildNameToCertificate()
        transport := &http.Transport{TLSClientConfig: tlsConfig}

        // https client request
        url := "https://test.temp:443"
        j := []byte(`{"id": "3232323", "name": "lambda"}`)
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{Transport: transport}

        // read response
        resp, err := client.Do(req)
        if err != nil {
                log.Fatal(err)
        }
        contents, err := ioutil.ReadAll(resp.Body)
        log.Println(string(contents))
}
