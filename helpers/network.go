package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

func ArePortsReachable(ip string, ports ...string) error {
	var conn net.Conn
	var err error

external:
	for _, port := range ports {
		for i := 0; i < 3; i++ {
			conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), 5*time.Second)
			if err == nil {
				conn.Close()
				continue external
			}
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			return fmt.Errorf("cannot connect to %s on port %s: %v", ip, port, err)
		}
	}

	return nil
}

func DoReq[response any](url string, data []byte, method string, headers map[string]string, skipTlsVerification bool) (response, int, error) {
	var result response

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return result, http.StatusInternalServerError, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	if !skipTlsVerification {
		tlsConfig, err := LoadHTTPTLSCredentials(filepath.Join(GetMyPath(), "certs", "utm.crt"))
		if err != nil {
			return result, http.StatusInternalServerError, fmt.Errorf("failed to load TLS credentials: %v", err)
		}
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return result, resp.StatusCode, fmt.Errorf("while sending request to %s received status code: %d and response body: %s", url, resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}

	return result, resp.StatusCode, nil
}
