// Copyright (c) 2019 aimerforreimu. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  GNU GENERAL PUBLIC LICENSE
//                        Version 3, 29 June 2007
//
//  Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
//  Everyone is permitted to copy and distribute verbatim copies
// of this license document, but changing it is not allowed.
//
// repo: https://github.com/aimerforreimu/common

package picbed

import (
	"github.com/jinfagang/picbed/common"
	"github.com/jinfagang/picbed/tools"
)

type JueJin struct {
	FileLimit []string
	MaxSize   int
}

func (s *JueJin) Upload(image *ImageParam) (ImageReturn, error) {
	url := "https://cdn-ms.juejin.im/v1/upload?bucket=gold-user-assets"

	file := &common.FormFile{
		Name:  image.Name,
		Key:   "file",
		Value: *image.Content,
		Type:  image.Type,
	}
	var header map[string]string
	data := tools.FormPost(file, url, header)
	j := common.JueJinResp{}
	err := j.UnmarshalJSON([]byte(data))
	if err != nil {
		return ImageReturn{}, err
	}
	return ImageReturn{
		Url: j.D.URL.HTTPS,
		ID:  10,
	}, nil
}

func (s *JueJin) UploadToJueJin(img []byte, imgInfo string, imgType string) string {
	url := "https://cdn-ms.juejin.im/v1/upload?bucket=gold-user-assets"
	name := tools.GetFileNameByMimeType(imgInfo)

	file := &common.FormFile{
		Name:  name,
		Key:   "file",
		Value: img,
		Type:  imgType,
	}
	var header map[string]string
	data := tools.FormPost(file, url, header)
	j := common.JueJinResp{}
	j.UnmarshalJSON([]byte(data))
	return j.D.URL.HTTPS
}
