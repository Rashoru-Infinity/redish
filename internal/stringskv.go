package internal

import (
        "encoding/json"
        "io"
        "net/http"
        "strings"
)

func GetStringsKV(r *http.Request) (*string, *string) {
        var kv map[string]string
        uri := r.RequestURI
        head := strings.LastIndex(uri, "/")
        if head < 0 || head == len(uri) - 1 {
                return nil, nil
        }
        key := uri[head + 1:]
        b, err := io.ReadAll(r.Body)
        if err != nil {
                return &key, nil
        }
        if err = json.Unmarshal(b, &kv); err != nil {
                return &key, nil
        }
        value, _ := kv[key]
        return &key, &value
}
