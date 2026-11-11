package options

import (
	_ "embed"
	"encoding/json"
)

//go:embed zh_hans.json
var zhHansJSON []byte

type zhEntry struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

var zhByKey map[string]zhEntry

func init() {
	_ = json.Unmarshal(zhHansJSON, &zhByKey)
	if zhByKey == nil {
		zhByKey = map[string]zhEntry{}
	}
	// Keys missing from AriaNg zh_Hans.
	extra := map[string]zhEntry{
		"index-out":              {Name: "索引输出路径", Desc: "按 INDEX=PATH 指定种子内文件的输出路径，可重复。"},
		"media-request-contexts": {Name: "媒体请求上下文", Desc: "按来源设置的 JSON 请求头。"},
		"media-input":            {Name: "媒体输入清单", Desc: "浏览器捕获的 JSON 清单。"},
		"rpc-secret":             {Name: "RPC 密钥", Desc: "JSON-RPC 鉴权令牌。"},
		"pause":                  {Name: "添加时暂停", Desc: "通过 RPC 添加的任务以暂停状态加入。"},
	}
	for k, v := range extra {
		if _, ok := zhByKey[k]; !ok {
			zhByKey[k] = v
		}
	}
}

// Display returns localized name and description for an option key.
// lang: "zh" or "en". In zh mode, never fall back to English prose
// (empty Chinese description stays empty to avoid mixed-language rows).
func Display(key, fallbackDesc, lang string) (name, desc string) {
	if lang == "en" {
		return key, fallbackDesc
	}
	if e, ok := zhByKey[key]; ok {
		name = e.Name
		if name == "" {
			name = key
		}
		return name, e.Desc
	}
	return key, ""
}
