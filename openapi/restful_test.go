// Copyright 1999-2019 Alibaba Group Holding Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package openapi

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/aliyun-cli/cli"
	"github.com/stretchr/testify/assert"

	"bufio"
	"testing"
)

func TestRestfulInvoker_Prepare(t *testing.T) {
	a := &RestfulInvoker{
		BasicInvoker: &BasicInvoker{
			request: requests.NewCommonRequest(),
		},
	}

	a.BasicInvoker.request.RegionId = "cn-hangzhou"
	a.BasicInvoker.request.Content = []byte("{")
	w := new(bufio.Writer)
	ctx := cli.NewCommandContext(w)

	bodyflag := NewBodyFlag()
	bodyflag.SetAssigned(true)
	ctx.Flags().Add(bodyflag)

	secureflag := NewSecureFlag()
	secureflag.SetAssigned(true)
	ctx.Flags().Add(secureflag)

	bodyfile := NewBodyFileFlag()
	bodyfile.SetAssigned(true)
	ctx.Flags().Add(bodyfile)

	ctx.SetUnknownFlags(cli.NewFlagSet())
	ctx.UnknownFlags().Add(NewBodyFlag())
	err := a.Prepare(ctx)
	assert.Nil(t, err)

	BodyFlag(ctx.Flags()).SetAssigned(false)
	BodyFileFlag(ctx.Flags()).SetAssigned(false)
	a.BasicInvoker.request.Content = []byte("{")
	err = a.Prepare(ctx)
	assert.Nil(t, err)

	a.BasicInvoker.request.Headers = map[string]string{}
	a.BasicInvoker.request.Content = []byte("<")
	err = a.Prepare(ctx)
	assert.Nil(t, err)
}

func TestRestfulInvoker_Call(t *testing.T) {
	client, err := sdk.NewClientWithAccessKey("regionid", "accesskeyid", "accesskeysecret")
	assert.Nil(t, err)

	a := &RestfulInvoker{
		BasicInvoker: &BasicInvoker{
			client:  client,
			request: requests.NewCommonRequest(),
		},
	}
	_, err = a.Call()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "[SDK.CanNotResolveEndpoint] Can not resolve endpoint")
}

func Test_checkRestfulMethod(t *testing.T) {
	w := new(bufio.Writer)
	ctx := cli.NewCommandContext(w)
	methodOrPath := "get"
	pathPattern := "/user"
	ok, method, path, err := checkRestfulMethod(ctx, methodOrPath, "")
	assert.False(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "", method)
	assert.Equal(t, "", path)

	ok, method, path, err = checkRestfulMethod(ctx, methodOrPath, pathPattern)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "GET", method)
	assert.Equal(t, "/user", path)

	pathPattern = "user"
	ok, method, path, err = checkRestfulMethod(ctx, methodOrPath, pathPattern)
	assert.True(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, "bad restful path user", err.Error())
	assert.Equal(t, "GET", method)
	assert.Equal(t, "", path)

	ctx.Flags().Add(NewRoaFlag())
	methodOrPath = "update"
	ok, method, path, err = checkRestfulMethod(ctx, methodOrPath, pathPattern)
	assert.False(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "", method)
	assert.Equal(t, "", path)

	RoaFlag(ctx.Flags()).SetAssigned(true)
	RoaFlag(ctx.Flags()).SetValue("get")
	ok, method, path, err = checkRestfulMethod(ctx, methodOrPath, pathPattern)
	assert.True(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, "bad restful path update", err.Error())
	assert.Equal(t, "get", method)
	assert.Equal(t, "", path)
}
