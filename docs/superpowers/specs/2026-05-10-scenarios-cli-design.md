# Scenarios CLI 命令接入设计规格

> **版本**：v1.0  
> **日期**：2026-05-10  
> **状态**：待实现  
> **目标版本**：md2wechat-skill v2.2.0

---

## 1. 问题与目标

### 背景

wechat-markdown-editor API 服务已于 v1（2026-05）上线 Scenarios API，提供三个端点：

- `GET /api/scenarios` — 列出 5 个商业内容场景
- `GET /api/scenarios/:id` — 场景详情（字段定义、samplePayload、Top 3 推荐主题）
- `POST /api/scenarios/apply` — 传入结构化字段，返回微信兼容 HTML

**当前 md2wechat-skill CLI v2.1.0 完全未接入 Scenarios API。**

### 解决的问题

Scenarios API 把「怎么排版」的决策权收到 API 侧，调用方只管填内容字段：

```
/api/convert 流程：   Markdown（调用方写布局）→ HTML
Scenarios 流程：      结构化字段（调用方填内容）→ 布局 Markdown → HTML
```

CLI 接入后，Agent 可以通过结构化接口完成公众号文章生成，不需要了解 `:::block` 语法、主题命名或模块排列规则。

### 成功标准

1. `md2wechat scenarios list/show/apply` 在本地 API（`http://localhost:3000`）和线上（`https://www.md2wechat.cn`）均可用
2. 所有命令输出符合现有 JSON envelope 协议（schema_version: v1）
3. Agent 可以完成以下工作流，无需额外文档：
   - `scenarios list` → `scenarios show <id>` → 用 samplePayload 作内容基础 → `scenarios apply <id>` → 得到 HTML
4. check.sh 静态检查通过（5 个场景 ID 在所有 SKILL.md 中存在）
5. `make quality-gates` 通过

---

## 2. 范围边界

### 本次实现（In Scope）

- `scenarios list` — 列出所有场景
- `scenarios show <id>` — 场景详情
- `scenarios apply <id> [content.json]` — 生成 HTML/Markdown
- `capabilities --json` 加入 scenarios 条目
- `scripts/check.sh` 静态一致性检查
- 文档同步（README / DISCOVERY.md / SCENARIOS.md / FAQ.md / SKILL.md×2 / CHANGELOG.md）
- Contract tests + Integration tests（对 localhost:3000）

### 本次不做（Out of Scope）

- `scenarios apply --draft`（直接发布草稿，下一版本）
- 交互式填字段（`--interactive` 模式）
- `scenarios generate`（AI 从 brief 自动填字段）
- 新增 scenario 类型（数量保持与 API 一致：5 个）

---

## 3. 架构设计

### 包结构

```
md2wechat-skill/
├── cmd/md2wechat/
│   ├── scenarios.go           # NEW：Cobra 子命令定义（list/show/apply）
│   ├── main.go                # MODIFY：注册 scenariosCmd，加 SCENARIOS_* code 常量
│   └── discovery.go           # MODIFY：capabilities 输出加 scenarios 条目
├── internal/
│   └── scenarios/
│       ├── client.go          # NEW：HTTP 客户端，调 /api/scenarios/* 端点
│       └── types.go           # NEW：Go 类型对应 API 响应
└── scripts/
    └── check.sh               # NEW：静态 + 动态一致性检查
```

### 包职责边界

| 包 | 职责 | 不做什么 |
|---|---|---|
| `internal/scenarios` | HTTP 请求/响应、类型定义、错误翻译 | 不知道 Cobra 存在，不做 UI，不做文件 I/O |
| `cmd/md2wechat/scenarios.go` | 参数解析、文件读取、调用 scenarios client、格式化 JSON envelope | 不直接构造 HTTP 请求 |
| `internal/converter` | Markdown→HTML 渲染（现有） | 不处理 scenarios 结构字段 |

此设计与现有 `internal/wechat`（WeChat SDK）、`internal/converter`（转换）、`cmd/md2wechat/convert.go`（命令层）三层分离模式一致。

---

## 4. 命令 API

### 4.1 `scenarios list`

```bash
md2wechat scenarios list [--json]
```

**输出（--json）：**

```json
{
  "success": true,
  "code": "SCENARIOS_LISTED",
  "message": "5 scenarios available",
  "schema_version": "v1",
  "status": "completed",
  "retryable": false,
  "data": {
    "scenarios": [
      {
        "id": "tutorial-guide",
        "name": "教程指南 / 工具教学",
        "description": "把教程、工具教学、AI 工作流指南组织成清晰步骤、注意事项和工具箱。",
        "primaryBuyer": "SaaS 内容团队、AI 课程团队、工具型产品运营",
        "readerJob": "让读者按步骤完成一个具体任务。",
        "moduleSequence": ["hero", "toc", "steps", "image-steps", "notice", "checklist", "toolbox", "cta"],
        "recommendedTheme": "bytedance"
      }
    ]
  },
  "error": ""
}
```

**Human-readable 输出（无 --json）：**

```
5 个可用场景：

  tutorial-guide      教程指南 / 工具教学        推荐主题: bytedance
  product-launch      产品发布 / 功能上线        推荐主题: bytedance
  course-conversion   课程 / 训练营转化          推荐主题: apple
  case-study          案例复盘 / 用户故事        推荐主题: bytedance
  industry-report     行业报告 / 白皮书导读      推荐主题: default

用 `md2wechat scenarios show <id>` 查看字段定义。
```

---

### 4.2 `scenarios show <id>`

```bash
md2wechat scenarios show tutorial-guide [--json]
```

**输出（--json）：**

```json
{
  "success": true,
  "code": "SCENARIO_SHOWN",
  "schema_version": "v1",
  "status": "completed",
  "retryable": false,
  "data": {
    "scenario": {
      "id": "tutorial-guide",
      "name": "教程指南 / 工具教学",
      "description": "...",
      "primaryBuyer": "...",
      "readerJob": "...",
      "moduleSequence": ["hero", "toc", "steps", ...],
      "contentMeta": {
        "bestFor": ["知识", "产品"],
        "density": "balanced",
        "contrast": "medium"
      },
      "recommendedThemes": [
        { "id": "bytedance", "score": 24, "mood": ["科技", "现代"] },
        { "id": "apple",     "score": 14, "mood": ["简洁", "高端"] },
        { "id": "default",   "score": 10, "mood": ["经典", "通用"] }
      ],
      "fields": [
        { "key": "title",    "label": "主标题",  "type": "text",      "required": true,  "maxLength": 80 },
        { "key": "outcome",  "label": "完成后效果", "type": "textarea", "required": true,  "maxLength": 200 },
        { "key": "steps",    "label": "教程步骤",  "type": "step-list", "required": true, "maxItems": 8 },
        { "key": "cta",      "label": "行动引导",  "type": "cta",      "required": false, "maxLength": 240 }
      ],
      "samplePayload": { ... }
    }
  },
  "error": ""
}
```

> **Agent 提示**：`samplePayload` 可直接作为 `scenarios apply` 的 content.json 基础。

---

### 4.3 `scenarios apply <id> [content.json]`

```bash
# 从文件读取（主要路径）
md2wechat scenarios apply tutorial-guide content.json [flags]

# 从 stdin 读取（Agent 管道场景）
echo '{"title":"...","steps":[...]}' | md2wechat scenarios apply tutorial-guide

# stdin 与文件同时存在：以文件为准
```

**Flags：**

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--output` | string | `html` | `html` / `markdown` / `both` |
| `--theme` | string | （空）| 不传则 API 自动语义推荐 |
| `--font-size` | string | `medium` | `small` / `medium` / `large` |
| `--background` | string | `default` | `default` / `grid` / `none` |
| `--write-html` | string | （空）| HTML 另存为文件路径（可选） |
| `--json` | bool | false | 输出 JSON envelope |

**输出（--json）：**

```json
{
  "success": true,
  "code": "SCENARIO_APPLIED",
  "message": "scenario applied: tutorial-guide → bytedance (auto)",
  "schema_version": "v1",
  "status": "completed",
  "retryable": false,
  "data": {
    "scenario": "tutorial-guide",
    "html": "<section style=\"...\">...</section>",
    "markdown": ":::hero\n...",
    "theme": "bytedance",
    "themeAutoSelected": true,
    "fontSize": "medium",
    "backgroundType": "default",
    "modules": ["hero", "toc", "steps", "checklist", "cta"],
    "wordCount": 320,
    "estimatedReadTime": 2
  },
  "error": ""
}
```

**`--output html`（默认）时，`data.markdown` 不出现。`--output markdown` 时，`data.html` 不出现。**

---

## 5. 类型定义（`internal/scenarios/types.go`）

```go
type ScenarioListItem struct {
    ID               string   `json:"id"`
    Name             string   `json:"name"`
    Description      string   `json:"description"`
    PrimaryBuyer     string   `json:"primaryBuyer"`
    ReaderJob        string   `json:"readerJob"`
    ModuleSequence   []string `json:"moduleSequence"`
    RecommendedTheme string   `json:"recommendedTheme"`
}

type RecommendedThemeEntry struct {
    ID    string   `json:"id"`
    Score int      `json:"score"`
    Mood  []string `json:"mood"`
}

type ScenarioField struct {
    Key       string `json:"key"`
    Label     string `json:"label"`
    Type      string `json:"type"`
    Required  bool   `json:"required"`
    MaxLength int    `json:"maxLength,omitempty"`
    MaxItems  int    `json:"maxItems,omitempty"`
}

type ScenarioDetail struct {
    ID                string                  `json:"id"`
    Name              string                  `json:"name"`
    Description       string                  `json:"description"`
    PrimaryBuyer      string                  `json:"primaryBuyer"`
    ReaderJob         string                  `json:"readerJob"`
    ModuleSequence    []string                `json:"moduleSequence"`
    ContentMeta       map[string]interface{}  `json:"contentMeta"`
    RecommendedThemes []RecommendedThemeEntry `json:"recommendedThemes"`
    Fields            []ScenarioField         `json:"fields"`
    SamplePayload     map[string]interface{}  `json:"samplePayload"`
}

type ApplyRequest struct {
    Scenario       string                 `json:"scenario"`
    Content        map[string]interface{} `json:"content"`
    Output         string                 `json:"output,omitempty"`
    Theme          string                 `json:"theme,omitempty"`
    FontSize       string                 `json:"fontSize,omitempty"`
    BackgroundType string                 `json:"backgroundType,omitempty"`
}

type ApplyResponse struct {
    Scenario          string `json:"scenario"`
    HTML              string `json:"html,omitempty"`
    Markdown          string `json:"markdown,omitempty"`
    Theme             string `json:"theme"`
    ThemeAutoSelected bool   `json:"themeAutoSelected"`
    FontSize          string `json:"fontSize"`
    BackgroundType    string `json:"backgroundType"`
    Modules           []string `json:"modules"`
    WordCount         int    `json:"wordCount"`
    EstimatedReadTime int    `json:"estimatedReadTime"`
}
```

---

## 6. HTTP 客户端设计（`internal/scenarios/client.go`）

```go
type Client struct {
    BaseURL    string
    APIKey     string
    HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client
func (c *Client) ListScenarios(ctx context.Context) ([]ScenarioListItem, error)
func (c *Client) GetScenario(ctx context.Context, id string) (*ScenarioDetail, error)
func (c *Client) ApplyScenario(ctx context.Context, req ApplyRequest) (*ApplyResponse, error)
```

**超时**：30 秒（与 converter API client 保持一致）。

**错误类型：**
```go
type ScenarioError struct {
    Code       string // 内部错误码
    Message    string // 用户可读的错误信息
    HTTPStatus int    // 原始 HTTP 状态码
    Retryable  bool
}
```

---

## 7. 错误处理

| 场景 | Code | Retryable | 消息 |
|------|------|-----------|------|
| content.json 不存在 | `SCENARIO_CONTENT_FILE_NOT_FOUND` | false | "内容文件不存在: %s，请检查路径" |
| content.json 格式错误 | `SCENARIO_CONTENT_INVALID_JSON` | false | "内容文件 JSON 格式错误: %s" |
| content.json > 256KB | `SCENARIO_CONTENT_TOO_LARGE` | false | "内容文件超过 256KB 限制" |
| API Key 未配置 | `API_KEY_MISSING` | false | "API Key 未配置，请运行 md2wechat config set api_key <key>" |
| 场景 ID 不存在 | `SCENARIO_NOT_FOUND` | false | "场景 '%s' 不存在，用 `md2wechat scenarios list` 查看有效 ID" |
| API 返回 400 | `SCENARIO_VALIDATION_FAILED` | false | 透传 API msg 字段 |
| API 返回 401 | `SCENARIO_AUTH_FAILED` | false | "API Key 无效，请重新获取" |
| 请求超时 | `SCENARIO_REQUEST_TIMEOUT` | true | "请求超时，请重试" |
| API 不可达 | `SCENARIO_API_UNREACHABLE` | true | "API 服务不可达，请检查网络连接" |
| --write-html 目录不存在 | `SCENARIO_OUTPUT_DIR_NOT_FOUND` | false | "输出目录不存在: %s" |

---

## 8. `capabilities --json` 变更

在 `data` 下新增：

```json
"scenarios": {
  "scenario_ids": [
    "tutorial-guide",
    "product-launch",
    "course-conversion",
    "case-study",
    "industry-report"
  ],
  "subcommands": ["list", "show", "apply"],
  "output_modes": ["html", "markdown", "both"],
  "theme_auto_select": true
}
```

---

## 9. `scripts/check.sh`

**检查项：**

```bash
#!/usr/bin/env bash
set -euo pipefail

SCENARIO_IDS=("tutorial-guide" "product-launch" "course-conversion" "case-study" "industry-report")
SKILL_FILES=(
  "skills/md2wechat/SKILL.md"
  "platforms/openclaw/md2wechat/SKILL.md"
  "docs/SCENARIOS.md"
)

# 1. 每个 scenario ID 必须在所有 SKILL 文件中出现
for id in "${SCENARIO_IDS[@]}"; do
  for file in "${SKILL_FILES[@]}"; do
    grep -q "$id" "$file" || { echo "FAIL: $id missing in $file"; exit 1; }
  done
done

# 2. capabilities 文档中提到 scenarios
grep -q "scenarios" docs/DISCOVERY.md || { echo "FAIL: scenarios missing in DISCOVERY.md"; exit 1; }

# 3. main.go 包含 scenarios 代码常量
grep -q "SCENARIOS_LISTED" cmd/md2wechat/main.go || { echo "FAIL: SCENARIOS_LISTED constant missing"; exit 1; }

# 4. 可选：本地 API 存活检查（有 localhost:3000 才跑）
if curl -s --connect-timeout 2 http://localhost:3000/api/scenarios > /dev/null 2>&1; then
  count=$(curl -s http://localhost:3000/api/scenarios \
    -H "X-API-Key: ${MD2WECHAT_API_KEY:-test}" \
    | python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d.get('data',{}).get('scenarios',[])))" 2>/dev/null || echo 0)
  [ "$count" -ge 5 ] || { echo "FAIL: API returned fewer than 5 scenarios (got $count)"; exit 1; }
  echo "OK: local API returned $count scenarios"
fi

echo "check.sh: all checks passed"
```

**集成到 Makefile：**
```makefile
check:
	@bash scripts/check.sh

quality-gates: fmt vet lint test check npm-pack release-check
```

---

## 10. 文档同步范围

| 文件 | 更新内容 |
|------|---------|
| `README.md` | scenarios 命令一行描述 + 快速示例 |
| `docs/DISCOVERY.md` | scenarios 发现命令说明（list/show 用法） |
| `docs/SCENARIOS.md` | 完整 CLI 侧使用文档（5 个场景速查 + Agent 工作流） |
| `docs/FAQ.md` | scenarios 常见问题 3 条 |
| `skills/md2wechat/SKILL.md` | scenarios 工作流 + 5 个 ID 速查 + 字段类型说明 |
| `platforms/openclaw/md2wechat/SKILL.md` | 同上 |
| `CHANGELOG.md` | v2.2.0 新增 scenarios 命令 |

---

## 11. 测试策略

### 层级 1：CLI 契约测试（无网络依赖）

位置：`cmd/md2wechat/main_contract_test.go`（追加到现有文件）

| 测试名 | 验证点 |
|-------|-------|
| `TestRunScenariosListJSON_EnvelopeShape` | JSON envelope 结构合法，含 schema_version |
| `TestRunScenariosListJSON_APIKeyMissing` | 无 key → 非 0 退出 + error 字段 |
| `TestRunScenariosShowJSON_NotFound` | id=invalid → SCENARIO_NOT_FOUND |
| `TestRunScenariosApplyJSON_ContentFileMissing` | 文件不存在 → SCENARIO_CONTENT_FILE_NOT_FOUND |
| `TestRunScenariosApplyJSON_InvalidJSON` | 内容 JSON 错误 → SCENARIO_CONTENT_INVALID_JSON |
| `TestRunScenariosApplyJSON_ContentTooLarge` | > 256KB → SCENARIO_CONTENT_TOO_LARGE |

### 层级 2：集成测试（需 localhost:3000）

位置：`cmd/md2wechat/scenarios_integration_test.go`（build tag: integration）

| 测试名 | 验证点 |
|-------|-------|
| `TestScenariosIntegration_ListAll` | 返回 5 个场景，每个有 recommendedTheme |
| `TestScenariosIntegration_ShowDetail` | tutorial-guide 有 fields + samplePayload + recommendedThemes |
| `TestScenariosIntegration_ApplyWithSample` | show 的 samplePayload 直接用于 apply → 成功 |
| `TestScenariosIntegration_ApplyAutoTheme` | 不传 theme → themeAutoSelected: true |
| `TestScenariosIntegration_ApplyExplicitTheme` | 传 theme="default" → themeAutoSelected: false |

### 层级 3：Unit 测试（`internal/scenarios/`）

| 测试名 | 验证点 |
|-------|-------|
| `TestClient_ParseListResponse` | JSON 解析正确性，含 recommendedTheme |
| `TestClient_ParseApplyResponse` | themeAutoSelected 字段正确解析 |
| `TestClient_4xxErrorPropagation` | API 400 msg 透传到 ScenarioError.Message |
| `TestClient_TimeoutRetryable` | 超时 → Retryable: true |

---

## 12. 版本与发布

- 版本号：`VERSION` 文件 + marketplace.json 更新为 `2.2.0`
- `CHANGELOG.md` 新增 v2.2.0 章节
- `make quality-gates` 必须通过（含新增的 `make check`）
- 发布前 E2E smoke test：`scenarios apply` 对 localhost:3000 验证真实渲染

---

## 13. Agent 推荐工作流（写入 SKILL.md）

```bash
# Step 1: 找到匹配的场景
md2wechat scenarios list --json

# Step 2: 获取字段定义和示例
md2wechat scenarios show tutorial-guide --json

# Step 3: 用 samplePayload 作为基础，填入真实内容
# → 将 .data.scenario.samplePayload 写入 content.json，替换为用户内容

# Step 4: 生成 HTML
md2wechat scenarios apply tutorial-guide content.json --json

# Step 5: （可选）写入 HTML 文件供后续命令消费
md2wechat scenarios apply tutorial-guide content.json --write-html output.html

# Step 6: 创建微信草稿
md2wechat create_draft output.html --title "文章标题" --cover cover.jpg
```
