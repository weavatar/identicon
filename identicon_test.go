// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package identicon

import (
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	dir, _ := os.Getwd()
	dir = filepath.Join(dir, "testdata")
	assert.NoError(t, os.MkdirAll(dir, os.ModePerm))
	defer os.RemoveAll(dir)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		t.Errorf("can not save generated images to %s", dir)
	}

	backColor := color.White
	imgMaker, err := New(64, backColor, DarkColors...)
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		s := strconv.Itoa(i)
		img := imgMaker.Make([]byte(s))

		f, err := os.Create(filepath.Join(dir, s+".png"))
		if !assert.NoError(t, err) {
			continue
		}
		defer f.Close()
		err = png.Encode(f, img)
		assert.NoError(t, err)
	}
}
