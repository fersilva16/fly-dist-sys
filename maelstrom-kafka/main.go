package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type SendRequest struct {
	maelstrom.MessageBody
	Key string `json:"key"`;
	Msg int `json:"msg"`;
}

type PollRequest struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`;
}

type CommitOffsetsRequest struct {
	maelstrom.MessageBody
	Offsets map[string]int `json:"offsets"`;
}

type ListCommittedOffsetsRequest struct {
	maelstrom.MessageBody
	Keys []string `json:"keys"`;
}

type Offsets struct {
	mu sync.RWMutex;
	offsets map[string]int;
}

type Messages struct {
	mu sync.RWMutex;
	messages map[string][][]int
}

var node = maelstrom.NewNode();
var offsets_offset = 0;
var o = Offsets{ offsets: make(map[string]int) };
var m = Messages{ messages: make(map[string][][]int) }

func main() {
	node.Handle("send", func(msg maelstrom.Message) error {
		var body SendRequest;
		
    if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
    }

		o.mu.Lock();

		if o.offsets[body.Key] == 0 {
			o.offsets[body.Key] = offsets_offset * 1000;
			offsets_offset++;
		}

		o.offsets[body.Key]++;

		o.mu.Unlock();

		m.mu.Lock();

		if m.messages[body.Key] == nil {
			m.messages[body.Key] = [][]int{};
		}

		m.messages[body.Key] = append(m.messages[body.Key], []int{ o.offsets[body.Key], body.Msg });

		m.mu.Unlock();
		
    res_body := map[string]any{
      "type": "send_ok",
			"offset": o.offsets[body.Key],
    };
    
    return node.Reply(msg, res_body);
  });

	node.Handle("poll", func(msg maelstrom.Message) error {
		var body PollRequest;
		
    if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
    }

		m.mu.RLock();

		res_messages := make(map[string][][]int);

		for key, msgs := range res_messages {
			if body.Offsets[key] == 0 {
				continue;
			}

			res_messages[key] = [][]int{};

			for i := 0; i < len(msgs); i++ {
				if msgs[i][0] < body.Offsets[key] {
					continue;
				}

				res_messages[key][i] = msgs[i];
			}
		}

		m.mu.RUnlock();

    res_body := map[string]any{
      "type": "poll_ok", 
			"msgs": res_messages,
    };
    
    return node.Reply(msg, res_body);
  });

	node.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body CommitOffsetsRequest;
		
    if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
    }

		o.mu.Lock();

		for key, offset := range body.Offsets {
			o.offsets[key] = offset;
		}
		
		o.mu.Unlock();

    res_body := map[string]any{
      "type": "commit_offsets_ok", 
    };
    
    return node.Reply(msg, res_body);
  });

	node.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body ListCommittedOffsetsRequest;
		
    if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err;
    }

		o.mu.RLock();

    res_body := map[string]any{
      "type": "list_committed_offsets_ok",
			"offsets": o.offsets,
    };

		o.mu.RUnlock();
    
    return node.Reply(msg, res_body);
  });

  
  if err := node.Run(); err != nil {
    log.Fatal(err);
  }
}