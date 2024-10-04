// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package svg

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	gitea_html "code.gitea.io/gitea/modules/html"
)

var svgIcons map[string]string
const cloudfrontBaseURL = "https://d21gfi7kzrpyzn.cloudfront.net/"

const defaultSize = 16

func fetchSVGContent(url string) (string, error) {
	// Send HTTP GET request to fetch the SVG content
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch SVG from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch SVG, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read SVG content: %v", err)
	}

	return string(body), nil
}

// Init discovers SVG icons and populates the `svgIcons` variable
func Init() error {
	iconMap := generateSVGIconMap()
	svgIcons = make(map[string]string)

	for name, url := range iconMap {
		// Fetch the content from the CloudFront URL
		content, err := fetchSVGContent(cloudfrontBaseURL+url)
		if err != nil {
			fmt.Printf("Error fetching %s from %s: %v\n", name, url, err)
			continue
		}

		// Store the fetched content in the svgIcons map
		svgIcons[name] = content
	}

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
/*func RenderHTML(icon string, others ...any) template.HTML {
	size, class := gitea_html.ParseSizeAndClass(defaultSize, "", others...)
	if svgURL, ok := svgIcons[icon]; ok {
		return template.HTML(fmt.Sprintf(`<img src="%s" class="svg %s" width="%d" height="%d" />`, fmt.Sprintf(`https://d21gfi7kzrpyzn.cloudfront.net/%s`, svgURL), class, size, size))
	}
	return template.HTML(fmt.Sprintf("<span>%s(%d/%s)</span>", template.HTMLEscapeString(icon), size, template.HTMLEscapeString(class)))
}*/
func RenderHTML(icon string, others ...any) template.HTML {
	size, class := gitea_html.ParseSizeAndClass(defaultSize, "", others...)
	if svgStr, ok := svgIcons[icon]; ok {
		// the code is somewhat hacky, but it just works, because the SVG contents are all normalized
		if size != defaultSize {
			svgStr = strings.Replace(svgStr, fmt.Sprintf(`width="%d"`, defaultSize), fmt.Sprintf(`width="%d"`, size), 1)
			svgStr = strings.Replace(svgStr, fmt.Sprintf(`height="%d"`, defaultSize), fmt.Sprintf(`height="%d"`, size), 1)
		}
		if class != "" {
			svgStr = strings.Replace(svgStr, `class="`, fmt.Sprintf(`class="%s `, class), 1)
		}
		return template.HTML(svgStr)
	}
	// during test (or something wrong happens), there is no SVG loaded, so use a dummy span to tell that the icon is missing
	return template.HTML(fmt.Sprintf("<span>%s(%d/%s)</span>", template.HTMLEscapeString(icon), size, template.HTMLEscapeString(class)))
}