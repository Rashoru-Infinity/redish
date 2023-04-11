package redish

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/go-redis/redis/v8"
	"github.com/rashoru-infinity/redish/internal"
)

func HandleStrings(w http.ResponseWriter, r *http.Request) {
	var (
		rc *redis.Client
		ctx = context.Background()
	)
	if rc = internal.GetRedisClient(r); rc == nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer rc.Close()
	switch r.Method {
	case "GET":
		key, _ := internal.GetStringsKV(r)
		if key == nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		value, err := rc.Get(ctx, *key).Result()
		if err == redis.Nil {
			http.Error(w, http.StatusText(404), 404)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		j, err := json.Marshal(map[string]string {
			*key : value,
		})
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", string(j))
		break
	case "POST":
		key, value := internal.GetStringsKV(r)
		if key == nil || value == nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if err := rc.Set(ctx, *key, *value, 0).Err(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		http.Error(w, http.StatusText(201), 201)
		break
	case "PUT":
		key, value := internal.GetStringsKV(r)
		if key == nil || value == nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		keyCount, err := rc.Exists(ctx, *key).Result()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if err := rc.Set(ctx, *key, *value, 0).Err(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if keyCount == 0 {
			http.Error(w, http.StatusText(201), 201)
			return
		}
		http.Error(w, http.StatusText(204), 204)
		break
	case "DELETE":
		key, _ := internal.GetStringsKV(r)
		if key == nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if err := rc.Del(ctx, *key).Err(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		http.Error(w, http.StatusText(204), 204)
		break
	default:
		http.Error(w, http.StatusText(405), 405)
		return
	}
}
