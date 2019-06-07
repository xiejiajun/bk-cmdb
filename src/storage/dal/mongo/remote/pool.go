/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package remote

import (
	"sync"

	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/rpc"
)

type pool struct {
	cache map[string]rpc.Client
	conn  rpc.Client
	sync.RWMutex
}

type client struct {
	p   *pool
	opt *dal.JoinOption
}

func NewPool(client rpc.Client) *pool {
	return &pool{
		cache: make(map[string]rpc.Client, 0),
		conn:  client,
	}
}

func (p *pool) JoinOption(opt *dal.JoinOption) *client {
	return &client{
		p:   p,
		opt: opt,
	}
}

func (c *client) Call(cmd string, input interface{}, result interface{}) error {
	rpcClient := c.p.conn
	if c.opt != nil && c.opt.TMAddr != "" {
		var err error
		rpcClient, err = c.GetRPCByAddr(c.opt.TMAddr)
		if err != nil {
			blog.ErrorJSON("client call addr(%s) err:%s, rid:%s", c.opt.TMAddr, err.Error(), c.opt.RequestID)
			return err
		}
	}

	return rpcClient.Call(cmd, input, result)
}

// GetRPCByAddr get rpc client by cache
func (c *client) GetRPCByAddr(addr string) (rpc.Client, error) {
	rpc, ok := c.getRPCByAddr(addr)
	if ok {
		return rpc, nil
	}
	return c.addRPCByAddr(addr)
}

func (c *client) getRPCByAddr(addr string) (rpc.Client, bool) {
	c.p.RLock()
	defer c.p.RUnlock()
	rpc, ok := c.p.cache[addr]
	return rpc, ok
}

func (c *client) addRPCByAddr(addr string) (rpc.Client, error) {
	c.p.Lock()
	defer c.p.Unlock()
	getSrvFunc := func() ([]string, error) {
		return []string{addr}, nil
	}
	return rpc.NewClientPool("tcp", getSrvFunc, "/txn/v3/rpc")
}
