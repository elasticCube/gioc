// Copyright 2017 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package depend

import (
	"github.com/vlorc/gioc/types"
	"github.com/vlorc/gioc/utils"
	"reflect"
)

func NewDependencyInject(scan types.DependencyScan, ref types.Reflect) types.DependencyInject {
	return &CoreDependencyInject{
		scan,
		ref,
	}
}
func (di *CoreDependencyInject) SetValue(v reflect.Value){
	di.Set(di.DependencyScan, v)
}

func (di *CoreDependencyInject) SetInterface(v interface{}){
	di.SetValue(reflect.ValueOf(v))
}

func (di *CoreDependencyInject) SubInject(provider types.Provider) types.DependencyInject {
	src := di.Get(di.DependencyScan)
	dst := utils.NewOf(src)

	return di.Depend().AsInject(dst)
}
