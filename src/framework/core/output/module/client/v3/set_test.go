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

package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/v3"
	//"configcenter/src/framework/core/types"
	"fmt"
	"testing"
)

func TestSearchSet(t *testing.T) {

	cli := cli := client.NewForConfig(config.Config{"supplierAccount": "0", "user": "build_user", "http://test.apiserver:8080": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_set_name").Like("作业").Field("bk_biz_id").Eq(2)
	dataMap, err := cli.SearchSets(cond)

	if nil != err {
		t.Errorf("failed to search, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		fmt.Println("item:", item.String("bk_set_name"))
	}
}
