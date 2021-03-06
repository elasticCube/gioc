// Copyright 2017 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package types

import "reflect"

type Error struct {
	Type    reflect.Type
	Name    string
	Code    ErrorCode
	Message string
}

type DependencyDescription struct {
	Type    reflect.Type
	Name    string
	Index   int
	Default reflect.Value
	Depend  Dependency
	Flags   DependencyFlag
}