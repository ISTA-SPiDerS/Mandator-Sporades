package benchmark

import (
	"async-consensus/proto"
	"context"
	"github.com/go-redis/redis/v8"
)

/*
	struct defining the benchmark
*/

type Benchmark struct {
	mode        int           // 0 for resident k/v store, 1 for redis
	RedisClient *redis.Client // redis client
	RedisCtx    context.Context
	KVStore     map[string]string
	name        int32 // name of the server
	keyLen      int
	valueLen    int
}

/*
	Initialize a new benchmark
*/

func Init(mode int, name int32, keyLen int, valueLen int) *Benchmark {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdsContext := context.Background()
	client.FlushAll(rdsContext) // delete the data base

	b := Benchmark{
		mode:        mode,
		RedisClient: client,
		RedisCtx:    rdsContext,
		KVStore:     make(map[string]string),
		name:        name,
		keyLen:      keyLen,
		valueLen:    valueLen,
	}

	return &b
}

/*
	external API to call
*/

func (b *Benchmark) Execute(miniBlock *proto.MemPoolMini) *proto.MemPoolMini {
	var commands []*proto.MemPoolMini_ClientBatch
	if b.mode == 0 {
		commands = b.residentExecute(miniBlock.Commands)
	} else {
		commands = b.redisExecute(miniBlock.Commands)
	}
	return &proto.MemPoolMini{
		Sender:   b.name,
		Receiver: miniBlock.Creator,
		UniqueId: miniBlock.UniqueId,
		Type:     8,
		Note:     miniBlock.Note,
		Commands: commands,
		Creator:  miniBlock.Creator,
	}
}

/*
	resident key value store operation: for each client request invoke the resident k/v store
*/

func (b *Benchmark) residentExecute(commands []*proto.MemPoolMini_ClientBatch) []*proto.MemPoolMini_ClientBatch {
	returnCommands := make([]*proto.MemPoolMini_ClientBatch, len(commands))

	for clientBatchIndex := 0; clientBatchIndex < len(commands); clientBatchIndex++ {

		returnCommands[clientBatchIndex] = &proto.MemPoolMini_ClientBatch{
			Commands: make([]*proto.MemPoolMini_SingleOperation, len(commands[clientBatchIndex].Commands)),
			Id:       commands[clientBatchIndex].Id,
			Creator:  commands[clientBatchIndex].Creator,
		}

		for clientRequestIndex := 0; clientRequestIndex < len(commands[clientBatchIndex].Commands); clientRequestIndex++ {
			returnCommands[clientBatchIndex].Commands[clientRequestIndex] = &proto.MemPoolMini_SingleOperation{
				UniqueId: commands[clientBatchIndex].Commands[clientRequestIndex].UniqueId,
				Command:  "",
			}

			cmd := commands[clientBatchIndex].Commands[clientRequestIndex].Command
			typ := cmd[0:1]
			key := cmd[1 : 1+b.keyLen]
			val := cmd[1+b.keyLen:]
			if typ == "0" { // write
				b.KVStore[key] = val
				returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = "0" + key + "ok"
			} else { // read
				v, ok := b.KVStore[key]
				if ok {
					returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = "1" + key + v
				} else {
					returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = "1" + key + "nil"
				}
			}
		}
	}
	return returnCommands
}

/*
	redis commands execution: batch the requests and execute
*/

func (b *Benchmark) redisExecute(commands []*proto.MemPoolMini_ClientBatch) []*proto.MemPoolMini_ClientBatch {
	returnCommands := make([]*proto.MemPoolMini_ClientBatch, len(commands))

	mset := make([]string, 0) // pending MSET requests
	mget := make([]string, 0) // pending MGET requests

	for clientBatchIndex := 0; clientBatchIndex < len(commands); clientBatchIndex++ {

		returnCommands[clientBatchIndex] = &proto.MemPoolMini_ClientBatch{
			Commands: make([]*proto.MemPoolMini_SingleOperation, len(commands[clientBatchIndex].Commands)),
			Id:       commands[clientBatchIndex].Id,
			Creator:  commands[clientBatchIndex].Creator,
		}

		for clientRequestIndex := 0; clientRequestIndex < len(commands[clientBatchIndex].Commands); clientRequestIndex++ {
			returnCommands[clientBatchIndex].Commands[clientRequestIndex] = &proto.MemPoolMini_SingleOperation{
				UniqueId: commands[clientBatchIndex].Commands[clientRequestIndex].UniqueId,
				Command:  "",
			}

			cmd := commands[clientBatchIndex].Commands[clientRequestIndex].Command
			typ := cmd[0:1]
			key := cmd[1 : 1+b.keyLen]
			val := cmd[1+b.keyLen:]
			if typ == "0" { // write
				mset = append(mset, key)
				mset = append(mset, val)
				returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = "0" + key + "ok" // writes always succeed
			} else { // read
				mget = append(mget, key)
				returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = ""
			}
		}
	}

	if len(mset) > 0 {
		// execute writes in a batch
		if err := b.RedisClient.MSet(b.RedisCtx, mset).Err(); err != nil {
			panic(err)
		}
	}

	if len(mget) > 0 {
		// execute reads in a batch
		vs, err := b.RedisClient.MGet(b.RedisCtx, mget...).Result()
		if err != nil {
			panic(err)
		}

		vsCount := 0

		for clientBatchIndex := 0; clientBatchIndex < len(commands); clientBatchIndex++ {
			for clientRequestIndex := 0; clientRequestIndex < len(commands[clientBatchIndex].Commands); clientRequestIndex++ {
				cmd := commands[clientBatchIndex].Commands[clientRequestIndex].Command
				typ := cmd[0:1]
				key := cmd[1 : 1+b.keyLen]
				if typ == "0" {
					// we already set the response for writes
				} else { // read
					if vs[vsCount] == nil {
						// key not found
						returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = "1" + key + "nil"
					} else {
						if rep, ok := vs[vsCount].(string); !ok {
							panic(vs[vsCount])
						} else {
							returnCommands[clientBatchIndex].Commands[clientRequestIndex].Command = "1" + key + rep
						}
					}
					vsCount++
				}
			}
		}
	}

	return returnCommands
}
