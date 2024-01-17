package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	ch := client.Watch(context.Background(), key)
	// listen to channel of key
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Type:%v key:%v value:%v\n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))

			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete { // only put and delete
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed, err:%v\n", err)
					continue
				}
			}
			// if del then pass a null strcut
			fmt.Printf(" get new conf:%v\n", newConf)
			newConfCh <- newConf
		}
	}
}
