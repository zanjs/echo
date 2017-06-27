/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package param

import (
	"strconv"
)

type StringSlice []string

func (p StringSlice) String() []string {
	return []string(p)
}

func (p StringSlice) Int(filters ...func(int) bool) []int {
	var filter func(int) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var ids []int
	for _, id := range p {
		i, _ := strconv.Atoi(id)
		if filter == nil || filter(i) {
			ids = append(ids, i)
		}
	}
	return ids
}

func (p StringSlice) Int64(filters ...func(int64) bool) []int64 {
	var filter func(int64) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var ids []int64
	for _, id := range p {
		i, _ := strconv.ParseInt(id, 10, 64)
		if filter == nil || filter(i) {
			ids = append(ids, i)
		}
	}
	return ids
}

func (p StringSlice) Int32(filters ...func(int32) bool) []int32 {
	var filter func(int32) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var ids []int32
	for _, id := range p {
		i, _ := strconv.ParseInt(id, 10, 32)
		iv := int32(i)
		if filter == nil || filter(iv) {
			ids = append(ids, iv)
		}
	}
	return ids
}

func (p StringSlice) Uint(filters ...func(uint) bool) []uint {
	var filter func(uint) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var ids []uint
	for _, id := range p {
		i, _ := strconv.ParseUint(id, 10, 64)
		iv := uint(i)
		if filter == nil || filter(iv) {
			ids = append(ids, iv)
		}
	}
	return ids
}

func (p StringSlice) Uint64(filters ...func(uint64) bool) []uint64 {
	var filter func(uint64) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var ids []uint64
	for _, id := range p {
		i, _ := strconv.ParseUint(id, 10, 64)
		if filter == nil || filter(i) {
			ids = append(ids, i)
		}
	}
	return ids
}

func (p StringSlice) Uint32(filters ...func(uint32) bool) []uint32 {
	var filter func(uint32) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var ids []uint32
	for _, id := range p {
		i, _ := strconv.ParseUint(id, 10, 32)
		iv := uint32(i)
		if filter == nil || filter(iv) {
			ids = append(ids, iv)
		}
	}
	return ids
}

func (p StringSlice) Float32(filters ...func(float32) bool) []float32 {
	var filter func(float32) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var values []float32
	for _, v := range p {
		i, _ := strconv.ParseFloat(v, 32)
		iv := float32(i)
		if filter == nil || filter(iv) {
			values = append(values, iv)
		}
	}
	return values
}

func (p StringSlice) Float64(filters ...func(float64) bool) []float64 {
	var filter func(float64) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var values []float64
	for _, v := range p {
		i, _ := strconv.ParseFloat(v, 64)
		if filter == nil || filter(i) {
			values = append(values, i)
		}
	}
	return values
}

func (p StringSlice) Bool(filters ...func(bool) bool) []bool {
	var filter func(bool) bool
	if len(filters) > 0 {
		filter = filters[0]
	}
	var values []bool
	for _, v := range p {
		i, e := strconv.ParseBool(v)
		if e != nil {
			continue
		}
		if filter == nil || filter(i) {
			values = append(values, i)
		}
	}
	return values
}
