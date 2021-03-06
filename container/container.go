// Copyright 2017 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package container

import (
	"github.com/vlorc/gioc/provider"
	"github.com/vlorc/gioc/types"
	"github.com/vlorc/gioc/utils"
	"reflect"
)

func NewWithContainer(provider types.Provider,parent types.Container, deep int) types.Container {
	var binderFactory types.BinderFactory
	var selectorFactory types.SelectorFactory
	var registerFactory types.RegisterFactory

	provider.Assign(&binderFactory)
	provider.Assign(&selectorFactory)
	provider.Assign(&registerFactory)

	if nil == selectorFactory || nil == registerFactory{
		return nil
	}

	sel,err := selectorFactory.Instance(binderFactory)
	if nil != err {
		panic(err)
	}

	reg, err := registerFactory.Instance(sel)
	if nil != err {
		panic(err)
	}

	return NewContainer(reg, parent, deep)
}

func NewChildContainer(provider types.Provider,parent types.Container, deep int) types.Container {
	var reg types.Register

	err := provider.Assign(&reg)
	if nil != err {
		return NewWithContainer(provider,parent,deep)
	}

	return NewContainer(reg, parent, deep)
}

func NewContainer(register types.Register, parent types.Container, deep int) types.Container {
	return &CoreContainer{
		register: register,
		parent:   parent,
		deep:     deep,
	}
}

func (c *CoreContainer) AsRegister() types.Register {
	return c.register
}

func (c *CoreContainer) AsProvider() types.Provider {
	return provider.NewProxyProvider(c)
}

func (c *CoreContainer) Seal() types.Container {
	return c
}

func (c *CoreContainer) Readonly() types.Container {
	return c
}

func (c *CoreContainer) Parent() types.Container {
	return c.parent
}

func (c *CoreContainer) Child() types.Container {
	return NewChildContainer(c,c,c.deep)
}

func (c *CoreContainer) ResolveNamed(impType interface{}, name string, deep int) (interface{}, error) {
	return c.ResolveType(utils.TypeOf(impType), name, deep)
}

func (c *CoreContainer) ResolveType(typ reflect.Type, name string, deep int) (instance interface{}, err error) {
	if factory := c.register.AsSelector().FactoryOf(typ,name); nil != factory {
		instance, err = factory.Instance(c)
		return
	}

	if deep < 0 {
		deep = c.deep
	}
	if deep--; nil == c.Parent() || deep < 0 {
		err = types.NewError(types.ErrFactoryNotFound, typ, name)
	} else {
		instance, err = c.Parent().ResolveType(typ, name, deep)
	}
	return
}

func (c *CoreContainer) Resolve(impType interface{}, args ...string) (interface{}, error) {
	var name string = types.DEFAULT_NAME
	if len(args) > 0 {
		name = args[0]
	}
	return c.ResolveNamed(impType, name, -1)
}

func (c *CoreContainer) Assign(dst interface{}, args ...string) (err error) {
	var name string = types.DEFAULT_NAME
	if len(args) > 0 {
		name = args[0]
	}
	return c.AssignNamed(dst, nil, name, -1)
}

func (c *CoreContainer) AssignNamed(dst interface{}, impType interface{}, name string, deep int) (err error) {
	dstValue := reflect.ValueOf(dst)
	if !dstValue.CanSet() {
		if reflect.Ptr != dstValue.Kind() {
			err = types.NewError(types.ErrTypeNotSet, dst)
			return
		}
		dstValue = dstValue.Elem()
	}

	var srcType reflect.Type
	if nil == impType {
		srcType = dstValue.Type()
	} else {
		srcType = utils.TypeOf(impType)
	}

	instance, err := c.ResolveType(srcType, name, deep)
	if nil != err {
		return
	}

	srcValue := reflect.ValueOf(instance)
	dstValue.Set(srcValue)

	return
}
