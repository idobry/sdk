package sdk

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"strings"

	"github.com/gosimple/slug"
)

var (
	boardID uint
)

// Constants for templating
const (
	TemplatingHideNone = iota
	TemplatingHideLabel
	TemplatingHideVariable
)

type (
	// Board represents Grafana dashboard.
	Board struct {
		// Default fields as described by schema v17.
		// See https://grafana.com/docs/grafana/latest/reference/dashboard/
		ID           uint       `json:"id,omitempty"`
		UID          string     `json:"uid,omitempty"`
		Title        string     `json:"title"`
		Tags         []string   `json:"tags"`
		Style        string     `json:"style"`
		Timezone     string     `json:"timezone"`
		Editable     bool       `json:"editable"`
		HideControls bool       `json:"hideControls,omitempty" graf:"hide-controls"`
		GraphTooltip int        `json:"graphTooltip,omitempty"`
		Panels       []*Panel   `json:"panels"`
		Time         Time       `json:"time"`
		Timepicker   Timepicker `json:"timepicker"`
		Templating   Templating `json:"templating"`
		Annotations  struct {
			List []*Annotation `json:"list"`
		} `json:"annotations"`
		Refresh       *BoolString `json:"refresh"`
		SchemaVersion uint        `json:"schemaVersion"`
		Version       uint        `json:"version"`
		Links         []*Link     `json:"links"`

		// Optional fields.
		Description     *string `json:"description,omitempty"`
		Slug            *string `json:"slug,omitempty""`
		OriginalTitle   *string `json:"originalTitle,omitempty"`
		SharedCrosshair bool    `json:"sharedCrosshair,omitempty" graf:"shared-crosshair"`
		Rows            []*Row  `json:"rows,omitempty"`
	}
	Time struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	Timepicker struct {
		// Default fields.
		RefreshIntervals []string `json:"refresh_intervals"`
		TimeOptions      []string `json:"time_options"`
		// Optional fields.
		Now *bool `json:"now,omitempty"`
	}
	Templating struct {
		List []*TemplateVar `json:"list"`
	}
	TemplateVar struct {
		Name        string    `json:"name"`
		Type        string    `json:"type"`
		Auto        bool      `json:"auto,omitempty"`
		AutoCount   *int      `json:"auto_count,omitempty"`
		Datasource  *string   `json:"datasource"`
		Refresh     BoolInt   `json:"refresh"`
		Options     []*Option `json:"options"`
		IncludeAll  bool      `json:"includeAll"`
		AllFormat   string    `json:"allFormat"`
		AllValue    string    `json:"allValue"`
		Multi       bool      `json:"multi"`
		MultiFormat string    `json:"multiFormat"`
		Query       string    `json:"query"`
		Regex       string    `json:"regex"`
		Current     Current   `json:"current"`
		Label       string    `json:"label"`
		Hide        uint8     `json:"hide"`
		Sort        int       `json:"sort"`
	}
	// for templateVar
	Option struct {
		Text     string `json:"text"`
		Value    string `json:"value"`
		Selected bool   `json:"selected"`
	}
	// for templateVar
	Current struct {
		Tags  []string    `json:"tags,omitempty"`
		Text  string      `json:"text"`
		Value interface{} `json:"value"` // TODO select more precise type
	}
	Annotation struct {
		BuiltIn    int      `json:"builtIn,omitempty"`
		Datasource *string  `json:"datasource"`
		Enable     bool     `json:"enable"`
		Hide       bool     `json:"hide,omitempty"`
		IconColor  string   `json:"iconColor"`
		IconSize   uint     `json:"iconSize,omitempty"`
		Limit      int      `json:"limit,omitempty"`
		LineColor  string   `json:"lineColor,omitempty"`
		Name       string   `json:"name"`
		Query      string   `json:"query,omitempty"`
		ShowIn     int      `json:"showIn"`
		ShowLine   bool     `json:"showLine,omitempty"`
		Tags       []string `json:"tags,omitempty"`
		TagsField  string   `json:"tagsField,omitempty"`
		TextField  string   `json:"textField,omitempty"`
		Type       string   `json:"type"`
	}
)

// Link represents Link to another dashboard or external weblink
type Link struct {
	AsDropdown  *bool    `json:"asDropdown,omitempty"`
	DashURI     *string  `json:"dashUri,omitempty"`
	Dashboard   *string  `json:"dashboard,omitempty"`
	Icon        *string  `json:"icon,omitempty"`
	IncludeVars bool     `json:"includeVars"`
	KeepTime    *bool    `json:"keepTime,omitempty"`
	Params      *string  `json:"params,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	TargetBlank *bool    `json:"targetBlank,omitempty"`
	Title       string   `json:"title,omitempty"`
	Tooltip     *string  `json:"tooltip,omitempty"`
	Type        string   `json:"type"`
	URL         *string  `json:"url,omitempty"`
}

func NewBoard(title string) *Board {
	boardID++
	return &Board{
		ID:           boardID,
		Title:        title,
		Style:        "dark",
		Timezone:     "browser",
		Editable:     true,
		HideControls: false,
		Rows:         []*Row{},
	}
}

func (b *Board) RemoveTags(tags ...string) {
	// order might change after removing the tags
	for _, toRemoveTag := range tags {
		tagLen := len(b.Tags)
		for i, tag := range b.Tags {
			if tag == toRemoveTag {
				b.Tags[tagLen-1], b.Tags[i] = b.Tags[i], b.Tags[tagLen-1]
				b.Tags = b.Tags[:tagLen-1]
				break
			}
		}
	}
}

func (b *Board) AddTags(tags ...string) {
	tagFound := make(map[string]bool, len(b.Tags))
	for _, tag := range b.Tags {
		tagFound[tag] = true
	}
	for _, tag := range tags {
		if tagFound[tag] {
			continue
		}
		b.Tags = append(b.Tags, tag)
		tagFound[tag] = true
	}
}

func (b *Board) HasTag(tag string) bool {
	for _, t := range b.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (b *Board) AddRow(title string) *Row {
	if title == "" {
		title = "New row"
	}
	row := &Row{
		Title:    title,
		Collapse: false,
		Editable: true,
		Height:   NewFloatOrString(FloatOrStringString("250px")),
	}
	b.Rows = append(b.Rows, row)
	return row
}

func (b *Board) UpdateSlug() string {
	s := strings.ToLower(slug.Make(b.Title))
	b.Slug = &s
	return s
}
