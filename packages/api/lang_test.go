// MIT License
//
// Copyright (c) 2016 GenesisCommunity
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package api

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLang(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName("lng")
	utfName := randName("lngutf")

	_, id, err := postTxResult("NewLang", &url.Values{
		"Name":          {name},
		"Trans":         {`{"en": "My test", "fr": "French string", "en-US": "US locale"}`},
		"ApplicationId": {"1"},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	cases := []struct {
		url    string
		form   url.Values
		expect string
	}{
		{
			"NewLang",
			url.Values{
				"Name":          {utfName},
				"Trans":         {`{"en": "ัะตัั"}`},
				"ApplicationId": {"1"},
			},
			"",
		},
		{
			"NewPage",
			url.Values{
				"Name":          {name},
				"Value":         {fmt.Sprintf("Span($%s$)", name)},
				"Menu":          {"default_menu"},
				"Conditions":    {"ContractConditions(`MainCondition`)"},
				"ApplicationId": {"1"},
			},
			"",
		},
		{
			"content/page/" + name,
			url.Values{"lang": {"fr"}},
			`[{"tag":"span","children":[{"tag":"text","text":"French string"}]}]`,
		},
		{
			"content/page/" + name,
			url.Values{"lang": {"en-GB"}},
			`[{"tag":"span","children":[{"tag":"text","text":"My test"}]}]`,
		},
		{
			"content/page/" + name,
			url.Values{"lang": {"en-US"}},
			`[{"tag":"span","children":[{"tag":"text","text":"US locale"}]}]`,
		},
		{
			"content",
			url.Values{
				"template": {
					fmt.Sprintf(`Div(){
						Button(Body: $%[1]s$ $,  Page:test).Alert(Text: $%[1]s$, ConfirmButton: $confirm$, CancelButton: $cancel$)
						Button(Body: LangRes(%[1]s) LangRes, PageParams: "test", ).Alert(Text: $%[1]s$, CancelButton: $cancel$)
					}`, utfName),
				},
				"app_id": {"1"},
			},
			`[{"tag":"div","children":[{"tag":"button","attr":{"alert":{"cancelbutton":"$cancel$","confirmbutton":"$confirm$","text":"ัะตัั"},"page":"test"},"children":[{"tag":"text","text":"ัะตัั $"}]},{"tag":"button","attr":{"alert":{"cancelbutton":"$cancel$","text":"ัะตัั"},"pageparams":{"test":{"text":"test","type":"text"}}},"children":[{"tag":"text","text":"ัะตัั"},{"tag":"text","text":" LangRes"}]}]}]`,
		},
		{
			"content",
			url.Values{
				`template`: {fmt.Sprintf(`Span(Text LangRes(%s)+LangRes(%[1]s,fr))`, name)},
				`app_id`:   {`1`},
			},
			`[{"tag":"span","children":[{"tag":"text","text":"Text My test"},{"tag":"text","text":"+French string"}]}]`,
		},
		{
			"content",
			url.Values{
				"template": {fmt.Sprintf(`Span(Text LangRes(%s)+LangRes(%[1]s,fr))`, name)},
				"lang":     {"fr"},
				"app_id":   {"1"},
			},
			`[{"tag":"span","children":[{"tag":"text","text":"Text French string"},{"tag":"text","text":"+French string"}]}]`,
		},
		{
			"EditLang",
			url.Values{
				"Id":    {id},
				"Trans": {`{"en": "My test", "fr": "French string", "es": "Spanish text"}`},
			},
			"",
		},
		{
			"content",
			url.Values{
				"template": {fmt.Sprintf(`Table(mysrc,"$%[1]s$=name")Span(Text LangRes(%[1]s,es) $%[1]s$) Input(Class: form-control, Placeholder: $%[1]s$, Type: text, Name: Name)`, name)},
				"app_id":   {"1"},
			},
			`[{"tag":"table","attr":{"columns":[{"Name":"name","Title":"My test"}],"source":"mysrc"}},{"tag":"span","children":[{"tag":"text","text":"Text Spanish text"},{"tag":"text","text":" My test"}]},{"tag":"input","attr":{"class":"form-control","name":"Name","placeholder":"My test","type":"text"}}]`,
		},
		{
			"content",
			url.Values{
				"template": {fmt.Sprintf(`MenuGroup($%s$){MenuItem(Ooops, ooops)}MenuGroup(nolang){MenuItem(no, no)}`, name)},
				"app_id":   {"1"},
			},
			fmt.Sprintf(`[{"tag":"menugroup","attr":{"name":"$%s$","title":"My test"},"children":[{"tag":"menuitem","attr":{"page":"ooops","title":"Ooops"}}]},{"tag":"menugroup","attr":{"name":"nolang","title":"nolang"},"children":[{"tag":"menuitem","attr":{"page":"no","title":"no"}}]}]`, name),
		},
	}

	for _, v := range cases {
		var ret contentResult

		if len(v.expect) == 0 {
			assert.NoError(t, postTx(v.url, &v.form))
			continue
		}

		assert.NoError(t, sendPost(v.url, &v.form, &ret))
		assert.Equal(t, v.expect, RawToString(ret.Tree))
	}
}
