// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package svg

import (
	"fmt"
	"html/template"

	gitea_html "code.gitea.io/gitea/modules/html"
)

var svgIcons map[string]string

const defaultSize = 16

// Init discovers SVG icons and populates the `svgIcons` variable
func Init() error {
	svgIcons = generateSVGIconMap()
	return nil
}

// MockIcon replaces the current icon temporarily for testing or mocking purposes
func MockIcon(icon string) func() {
	if svgIcons == nil {
		svgIcons = generateSVGIconMap()
	}
	orig, exist := svgIcons[icon]
	svgIcons[icon] = fmt.Sprintf(`https://d21gfi7kzrpyzn.cloudfront.net/public/assets/img/%s.svg`, icon)
	return func() {
		if exist {
			svgIcons[icon] = orig
		} else {
			delete(svgIcons, icon)
		}
	}
}

// RenderHTML renders icons - arguments icon name (string), size (int), class (string)
func RenderHTML(icon string, others ...any) template.HTML {
	size, class := gitea_html.ParseSizeAndClass(defaultSize, "", others...)
	if svgURL, ok := svgIcons[icon]; ok {
		return template.HTML(fmt.Sprintf(`<img src="%s" class="svg %s" width="%d" height="%d" />`, svgURL, class, size, size))
	}
	return template.HTML(fmt.Sprintf("<span>%s(%d/%s)</span>", template.HTMLEscapeString(icon), size, template.HTMLEscapeString(class)))
}