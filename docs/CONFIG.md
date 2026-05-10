# 配置指南

这份文档解决 4 个最常见的问题：

1. 配置文件在哪里
2. 默认 API 域名在哪里改
3. Agent 应该先看哪里
4. 哪些功能分别需要哪些凭证

如果你现在卡在：

- 不知道 AppID / AppSecret 去哪拿
- 不知道微信 IP 白名单在哪配
- 明明配了凭证但还是 `ip not in whitelist`

先看：

- [微信凭证与 IP 白名单指南](WECHAT-CREDENTIALS.md)

如果你只想先跑通主路径，先看下面这 3 步。

## 3 步完成基础配置

### 1. 生成示例配置

```bash
md2wechat config init
```

默认会生成到：

```text
~/.config/md2wechat/config.yaml
```

你也可以显式指定输出位置：

```bash
md2wechat config init ./md2wechat.yaml
```

### 2. 打开配置文件，先填最小必需项

```yaml
wechat:
  appid: "你的微信公众号 AppID"
  secret: "你的微信公众号 Secret"

api:
  md2wechat_key: "你的 md2wechat API Key"
  md2wechat_base_url: "https://www.md2wechat.cn"
  convert_mode: "api"
  default_theme: "default"
```

### 3. 验证当前配置

```bash
md2wechat config validate
md2wechat config show --format json
```

---

## Agent 和用户应该先看哪里

如果你不知道去哪改配置，按这个顺序找：

1. `~/.config/md2wechat/config.yaml`
2. 环境变量
3. 当前目录下的 `md2wechat.yaml` / `md2wechat.yml` / `md2wechat.json`

对 Agent 来说，**默认应该优先检查 `~/.config/md2wechat/config.yaml`**。
如果用户说“把 API 域名改成备用域名”“切换图片服务”“检查当前配置”，先运行：

```bash
md2wechat config show --format json
```

这样可以直接看到当前生效的：

- `config_file`
- `md2wechat_base_url`
- `image_provider`
- `image_api_base`
- `default_convert_mode`

注意这里看到的是 **`config show --format json` 的扁平输出字段名**，不是配置文件里的嵌套 YAML 键名。

例如：

- 配置文件里写的是 `api.image_base_url`
- `config show --format json` 里看到的是 `image_api_base`

---

## 默认 API 域名在哪里改

项目当前默认值是：

```text
https://www.md2wechat.cn
```

它**不是写死不可改**。你有两种常用改法。

### 方式一：改配置文件

编辑 `~/.config/md2wechat/config.yaml`：

```yaml
api:
  md2wechat_base_url: "https://www.md2wechat.cn"
```

如果你要切到备用域名：

```yaml
api:
  md2wechat_base_url: "https://md2wechat.app"
```

### 方式二：用环境变量临时覆盖

```bash
export MD2WECHAT_BASE_URL="https://md2wechat.app"
```

环境变量优先级高于配置文件，适合：

- 临时切换备用域名
- CI / Agent 自动化
- 不想修改全局配置文件的场景

### 关于默认转换模式

当前 CLI 的默认行为是固定的：

- 不传 `--mode` 时，`md2wechat convert ...` 始终默认走 `api`
- 只有显式传入 `--mode ai` 时，才会走 AI 模式

也就是说，下面这个命令：

```bash
md2wechat convert article.md
```

当前一定等价于：

```bash
md2wechat convert article.md --mode api
```

所以如果用户没有填写配置，或者没有显式传 `--mode`，默认也是 `api`。

`api.convert_mode` / `CONVERT_MODE` 当前主要用于配置展示、校验和兼容字段；**不会覆盖 `convert` 命令在未传 `--mode` 时的默认行为**。

---

## 内置资产

当前仓库把官方默认 `themes` 和默认 `writer style` 随二进制一起提供。
这意味着即使 Agent 服务器上没有仓库目录，默认主题和默认写作风格也应该可用。

### 主题加载顺序

`themes` 的优先级从高到低如下：

1. `MD2WECHAT_THEMES_DIR`
2. 当前项目目录下的 `themes/`
3. `~/.config/md2wechat/themes/`
4. 二进制内置的官方默认 themes

同名主题以前面的来源覆盖后面的来源。

### 写作风格加载顺序

`writers` 的优先级从高到低如下：

1. `MD2WECHAT_WRITERS_DIR`
2. 当前项目目录下的 `writers/`
3. `~/.config/md2wechat/writers/`
4. `~/.md2wechat-writers/`
5. 二进制内置的默认 writer style

同名写作风格同样以前面的来源覆盖后面的来源。

### 什么时候改哪里

如果你想：

- 仅当前项目生效，放到项目目录
- 所有项目都生效，放到 `~/.config/md2wechat/...`
- Agent 服务器显式指定，设置 `MD2WECHAT_THEMES_DIR` 或 `MD2WECHAT_WRITERS_DIR`
- 保持官方默认不变，直接用内置资产

---

## 配置文件搜索顺序

程序会按以下顺序查找配置文件：

1. `~/.config/md2wechat/config.yaml`
2. `~/.md2wechat.yaml`
3. `~/.md2wechat.yml`
4. `./md2wechat.yaml`
5. `./md2wechat.yml`
6. `./md2wechat.json`
7. `./.md2wechat.yaml`
8. `./.md2wechat.yml`
9. `./.md2wechat.json`

实践上建议：

- 全局默认配置放 `~/.config/md2wechat/config.yaml`
- 项目特殊配置再放当前目录

---

## 完整示例配置

仓库里提供了一份可直接参考的示例：

- [config.yaml.example](examples/config.yaml.example)

完整示例：

```yaml
wechat:
  appid: "your_wechat_appid"
  secret: "your_wechat_secret"

api:
  md2wechat_key: "your_md2wechat_api_key"
  md2wechat_base_url: "https://www.md2wechat.cn"
  image_key: "your_image_api_key"
  image_base_url: "https://ark.cn-beijing.volces.com/api/v3"
  image_provider: "volcengine"
  image_model: "doubao-seedream-5-0-260128"
  image_size: "2K"
  convert_mode: "api"
  default_theme: "default"
  background_type: "none"
  http_timeout: 30

image:
  compress: true
  max_width: 1920
  max_size_mb: 5
```

## 三套命名要分清

当前最容易混淆的是：同一个配置项会同时出现在 3 个地方，但名字不完全一样。

### 1. 配置文件字段名

这是你在 `config.yaml` 里实际填写的名字，例如：

- `wechat.appid`
- `api.md2wechat_key`
- `api.image_base_url`
- `api.background_type`

### 2. 环境变量名

这是终端或 CI 里覆盖配置时使用的名字，例如：

- `WECHAT_APPID`
- `MD2WECHAT_API_KEY`
- `IMAGE_API_BASE`
- `DEFAULT_BACKGROUND_TYPE`

### 3. `config show --format json` 输出字段名

这是 CLI 为了更稳定的 machine-readable 输出而提供的扁平字段，例如：

- `wechat_appid`
- `md2wechat_api_key`
- `image_api_base`
- `default_background_type`

所以如果你是在：

- 改配置文件：用 `api.image_base_url`
- 查环境变量：看 `IMAGE_API_BASE`
- 解析 `config show --format json`：看 `image_api_base`

不要把这三套名字混成一个层次。

---

## 配置项说明

### 微信配置

| 配置项 | 必需 | 说明 |
|--------|------|------|
| `wechat.appid` | 创建草稿、上传图片时需要 | 微信公众号 AppID |
| `wechat.secret` | 创建草稿、上传图片时需要 | 微信公众号 Secret |

### API 转换配置

| 配置项 | 必需 | 说明 | 默认值 |
|--------|------|------|--------|
| `api.md2wechat_key` | API 模式需要 | md2wechat API Key | - |
| `api.md2wechat_base_url` | 否 | 排版 API 域名 | `https://www.md2wechat.cn` |
| `api.convert_mode` | 否 | 默认转换模式 | `api` |
| `api.default_theme` | 否 | 默认主题 | `default` |
| `api.background_type` | 否 | 背景类型 | `none` |
| `api.http_timeout` | 否 | HTTP 超时秒数 | `30` |

### 图片生成配置

| 配置项 | 必需 | 说明 | 默认值 |
|--------|------|------|--------|
| `api.image_key` | AI 图片时需要 | 图片生成 API Key | - |
| `api.image_provider` | 否 | 图片服务提供方 | `openai` |
| `api.image_base_url` | 否 | 图片服务地址 | `https://api.openai.com/v1` |
| `api.image_model` | 否 | 图片模型 | `gpt-image-1.5` |
| `api.image_size` | 否 | 默认图片执行尺寸/宽高比 | 跟随当前 provider，例如 `openai=1024x1024`、`volcengine=2K` |

当前内置 provider：`openai`、`tuzi`、`modelscope` (`ms`)、`openrouter` (`or`)、`gemini` (`google`)、`volcengine` (`volc`)。

### 图片处理配置

| 配置项 | 必需 | 说明 | 默认值 |
|--------|------|------|--------|
| `image.compress` | 否 | 是否自动压缩 | `true` |
| `image.max_width` | 否 | 最大宽度 | `1920` |
| `image.max_size_mb` | 否 | 最大大小（MB） | `5` |

---

## 环境变量对照表

| 环境变量 | 对应配置项 |
|----------|------------|
| `WECHAT_APPID` | `wechat.appid` |
| `WECHAT_SECRET` | `wechat.secret` |
| `MD2WECHAT_API_KEY` | `api.md2wechat_key` |
| `MD2WECHAT_BASE_URL` | `api.md2wechat_base_url` |
| `IMAGE_API_KEY` | `api.image_key` |
| `IMAGE_API_BASE` | `api.image_base_url` |
| `IMAGE_PROVIDER` | `api.image_provider` |
| `IMAGE_MODEL` | `api.image_model` |
| `IMAGE_SIZE` | `api.image_size` |
| `CONVERT_MODE` | `api.convert_mode` |
| `DEFAULT_THEME` | `api.default_theme` |
| `DEFAULT_BACKGROUND_TYPE` | `api.background_type` |
| `HTTP_TIMEOUT` | `api.http_timeout` |
| `COMPRESS_IMAGES` | `image.compress` |
| `MAX_IMAGE_WIDTH` | `image.max_width` |
| `MAX_IMAGE_SIZE` | `image.max_size_mb` |
| `MD2WECHAT_THEMES_DIR` | `themes` 覆盖目录 |
| `MD2WECHAT_WRITERS_DIR` | `writers` 覆盖目录 |

图片生成相关命令还支持 `--model`，用于单次覆盖当前调用的图片模型。优先级顺序为：

1. `--model`
2. `IMAGE_MODEL`
3. `api.image_model`
4. provider 默认模型

## `config show --format json` 常见字段对照

如果你是在排查 Agent / 脚本实际读到的配置，最常见的不是 YAML 字段，而是下面这些扁平 key：

| `config show --format json` 字段 | 对应配置文件字段 |
|---|---|
| `wechat_appid` | `wechat.appid` |
| `wechat_secret` | `wechat.secret` |
| `md2wechat_api_key` | `api.md2wechat_key` |
| `md2wechat_base_url` | `api.md2wechat_base_url` |
| `image_api_key` | `api.image_key` |
| `image_api_base` | `api.image_base_url` |
| `image_provider` | `api.image_provider` |
| `image_model` | `api.image_model` |
| `image_size` | `api.image_size` |
| `default_convert_mode` | `api.convert_mode` |
| `default_theme` | `api.default_theme` |
| `default_background_type` | `api.background_type` |
| `compress_images` | `image.compress` |
| `max_image_width` | `image.max_width` |
| `max_image_size_mb` | `image.max_size_mb` |
| `http_timeout` | `api.http_timeout` |
| `config_file` | 当前实际命中的配置文件路径 |

---

## 常见场景怎么配

### 只预览，不创建草稿

最小需要：

```yaml
api:
  md2wechat_key: "your_md2wechat_api_key"
  md2wechat_base_url: "https://www.md2wechat.cn"
  convert_mode: "api"
```

### 需要上传图片和创建草稿

最小需要：

```yaml
wechat:
  appid: "your_wechat_appid"
  secret: "your_wechat_secret"

api:
  md2wechat_key: "your_md2wechat_api_key"
```

### 需要 AI 图片生成

最小需要：

```yaml
wechat:
  appid: "your_wechat_appid"
  secret: "your_wechat_secret"

api:
  image_key: "your-ark-api-key"
  image_provider: "volcengine"
  image_model: "seedream-3-0"
  image_size: "2K"
```

补充说明：

- `api.image_size` / `IMAGE_SIZE` 控制的是实际发给图片 provider 的默认执行尺寸
- `generate_image --size ...` 会覆盖配置文件里的 `api.image_size`
- 图片 prompt 里的 `default_aspect_ratio` 是 preset 的语义默认画幅，用于渲染 prompt 与默认视觉比例
- 对于 Gemini / OpenRouter 这类支持比例格式的 provider，`api.image_size` 可以直接写成 `16:9`、`3:4`、`21:9`
- 对于 Volcengine Ark 当前接入，`api.image_size` 使用尺寸等级，例如 `2K`、`3K`；如果省略，当前默认值是 `2K`
- `api.image_base_url` 对 OpenAI、TuZi、ModelScope、OpenRouter、Volcengine 生效；Gemini 直连模式当前固定走官方 Go SDK backend，不读取该配置

---

## 配置优先级

优先级从高到低：

```text
命令行参数 > 环境变量 > 配置文件 > 默认值
```

举例：

1. 配置文件里写了：

```yaml
api:
  md2wechat_base_url: "https://www.md2wechat.cn"
```

2. 当前终端又执行了：

```bash
export MD2WECHAT_BASE_URL="https://md2wechat.app"
```

最终生效的是：

```text
https://md2wechat.app
```

---

## 自检命令

```bash
md2wechat config init
md2wechat config show --format json
md2wechat config validate
```

推荐排查顺序：

1. 先看 `config_file` 指向哪个文件
2. 再看 `md2wechat_base_url` 是否真是你想要的域名
3. 再看 `image_provider` / `image_api_base` 是否匹配
   这里的 `image_api_base` 是 `config show --format json` 的输出字段；配置文件里对应的是 `api.image_base_url`
4. 最后检查环境变量是否把文件里的值覆盖掉了

---

---

## Brand Profile

Brand Profile 是 Agent 读取的品牌与风格配置文件，**CLI 不解析此文件**。

与 CLI 运行时配置（`~/.md2wechat.yaml`）不同，Brand Profile 专门为 Agent 设计，用于记录内容生成的风格约束和偏好。

### 快速开始

```bash
# 初始化 Brand Profile（幂等操作，文件存在时不覆盖）
md2wechat brand init

# 查看当前 Brand Profile
md2wechat brand show
md2wechat brand show --json
```

Brand Profile 位置：

```text
~/.config/md2wechat/brand.yaml
```

### Schema 完整说明

以下是一份注释示例，展示所有支持的字段：

```yaml
# md2wechat Brand Profile
schema_version: 1

# 品牌名称或你的名字
name: "极客杰尼"

# 声音与语气定义
voice:
  # 你的语气风格（自由文本或关键字）
  # 例如：犀利实用，第一人称 | 深度分析，学术风格 | 温暖鼓励，故事导向
  tone: "犀利实用，第一人称"
  
  # 可选：风格参考文档或目录
  # 指向本地文件或目录，Agent 可参考此规范进行内容创作
  # 例如：~/Documents/brand/voice-guide.md 或 ~/brand-assets/
  # style_ref: ~/Documents/brand/voice-guide.md
  
  # 要避免的表达方式（列表）
  avoid:
    - 过度营销
    - 空泛鸡汤
    - 过多 emoji

# 排版与布局偏好
layout:
  # 文章开头风格
  # 可选值：verdict_first（结论先行）| story_first（故事开场）
  #        question_first（提问开场）| data_first（数据开场）
  opening: "verdict_first"

# 内容约束与限制
limits:
  # 单篇最多使用的高级排版模块数（上限 43）
  # 设为 0 表示使用默认值 6
  max_modules: 6
  
  # 最多行动按钮数（上限 2）
  max_cta: 1
  
  # 最多引用块数（上限 10）
  max_quotes: 2
  
  # 最多 Hero 块数（固定 1）
  max_hero: 1

# 默认行动按钮 (CTA)
cta:
  # 默认按钮标题
  default_title: "如果这篇对你有启发"
  
  # 默认按钮正文
  default_body: "欢迎关注极客杰尼的专栏，获取更多 AI 工具与独立开发实践"
  
  # 默认按钮行动（例如：关注 / 咨询 / 订阅）
  default_action: "关注"

# 作者卡片信息
author_card:
  # 显示名称
  name: "极客杰尼"
  
  # 职位或身份标签
  title: "AI 应用开发者 / 独立开发者"
  
  # 自我介绍
  bio: "记录 AI 工具、内容系统和独立产品实践。"
```

### 字段说明表

| 字段 | 类型 | 必需 | 说明 | 默认值 |
|------|------|------|------|--------|
| `schema_version` | 整数 | 是 | 格式版本 | `1` |
| `name` | 字符串 | 否 | 品牌名称或作者名 | （无） |
| `voice.tone` | 字符串 | 否 | 语气风格描述 | （无） |
| `voice.style_ref` | 字符串 | 否 | 风格参考文件/目录路径 | （无） |
| `voice.avoid` | 字符串数组 | 否 | 要避免的表达方式 | `[]` |
| `layout.opening` | 字符串 | 否 | 开头风格选择 | `verdict_first` |
| `limits.max_modules` | 整数 | 否 | 最多模块数（0=默认6） | `6` |
| `limits.max_cta` | 整数 | 否 | 最多 CTA 数 | `1` |
| `limits.max_quotes` | 整数 | 否 | 最多引用块数 | `2` |
| `limits.max_hero` | 整数 | 否 | 最多 Hero 块数 | `1` |
| `cta.default_title` | 字符串 | 否 | 默认按钮标题 | （无） |
| `cta.default_body` | 字符串 | 否 | 默认按钮正文 | （无） |
| `cta.default_action` | 字符串 | 否 | 默认按钮行动 | （无） |
| `author_card.name` | 字符串 | 否 | 作者显示名称 | （无） |
| `author_card.title` | 字符串 | 否 | 作者职位/身份 | （无） |
| `author_card.bio` | 字符串 | 否 | 作者自我介绍 | （无） |

### 降级行为与容错

1. **文件不存在**：Agent 继续工作，不报错；使用内置默认值。

2. **字段缺失**：未指定的字段使用表中的默认值。

3. **超出上限自动截断（Sanity Caps）**：
   - 若 `max_modules: 100`（超过上限 43），Agent 自动降至 `43`
   - 若 `max_cta: 10`（超过上限 2），Agent 自动降至 `2`
   - 若 `max_quotes: 50`（超过上限 10），Agent 自动降至 `10`
   - 若 `max_hero: 2`（超过上限 1），Agent 自动降至 `1`

4. **YAML 语法错误**：
   - `md2wechat brand show` 会检测并返回错误信息
   - Agent 应在读取前先执行 `brand show --json` 验证

### 与 CLI 运行时配置的区别

| 配置文件 | 位置 | 用途 | 解析方 | 必需 |
|---------|------|------|--------|------|
| CLI 运行时配置 | `~/.md2wechat.yaml` | API Keys、Provider、主题 | CLI | 创建草稿时必需 |
| Brand Profile | `~/.config/md2wechat/brand.yaml` | 内容风格、排版偏好、约束 | Agent | 可选（无则使用默认） |

**CLI 运行时配置** 典型场景：切换图片 Provider、配置 WeChat AppID、选择主题。

**Brand Profile** 典型场景：Agent 生成内容时遵守品牌约束、追踪作者信息、统一语气风格。

### 常见场景

#### 只初始化，保持最小配置

```bash
md2wechat brand init
# 编辑 ~/.config/md2wechat/brand.yaml，填入 name 和 voice.tone 即可
```

结果：Agent 会尊重你的品牌名和语气，但使用所有其他默认值。

#### Agent 读取并应用 Brand Profile

Agent 应该：

```bash
# 1. 检查 Brand Profile 是否存在
md2wechat brand show --json

# 2. 如果存在，解析并应用规则
# 3. 生成内容时遵守 limits 约束、应用 voice.tone、使用 author_card 信息
```

#### Sanity Caps 自动截断示例

配置文件：

```yaml
limits:
  max_modules: 50    # 超过上限 43
  max_cta: 5         # 超过上限 2
```

Agent 行为：

```
实际使用的上限：
  max_modules → 43（自动降至上限）
  max_cta → 2（自动降至上限）
```

---

## 相关文档

- [新手快速开始](QUICKSTART.md)
- [安装指南](INSTALL.md)
- [图片服务配置](IMAGE_PROVISIONERS.md)
- [真实烟雾测试记录](SMOKE.md)
- [内置资产](#内置资产)
