/*
Copyright CETCS. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	SPDX-License-Identifier: Apache-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package gm

import "hash"

type sm3sig struct {
	msg []byte
}

func NewSM3Sig() hash.Hash {
	return &sm3sig{}
}

func (d *sm3sig) Write(p []byte) (n int, err error) {
	d.msg = append(d.msg, p...)
	return len(d.msg), nil
}

func (d *sm3sig) Sum(b []byte) []byte {
	if b != nil {
		panic("sm3sig fail: b must be nil")
	}

	return d.msg
}

func (d *sm3sig) Reset() {
	d.msg = d.msg[:0]
}

func (d *sm3sig) Size() int {
	return 0
}

func (d *sm3sig) BlockSize() int {
	return 0
}