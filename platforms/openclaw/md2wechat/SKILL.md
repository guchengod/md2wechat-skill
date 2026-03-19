---
name: md2wechat
description: Convert Markdown to WeChat Official Account HTML, process images, and create drafts for OpenClaw.
metadata: {"openclaw":{"emoji":"📝","homepage":"https://github.com/geekjourneyx/md2wechat-skill","primaryEnv":"WECHAT_APPID","requires":{"env":["WECHAT_APPID","WECHAT_SECRET"]},"install":[{"id":"openclaw-installer-shell","kind":"download","label":"Download fixed-version OpenClaw installer (shell)","url":"https://github.com/geekjourneyx/md2wechat-skill/releases/download/v1.11.1/install-openclaw.sh","os":["darwin","linux"]},{"id":"openclaw-installer-powershell","kind":"download","label":"Download fixed-version installer (PowerShell)","url":"https://github.com/geekjourneyx/md2wechat-skill/releases/download/v1.11.1/install.ps1","os":["win32"]},{"id":"openclaw-skill-bundle","kind":"download","label":"Download OpenClaw skill bundle","url":"https://github.com/geekjourneyx/md2wechat-skill/releases/download/v1.11.1/md2wechat-openclaw-skill.tar.gz","archive":"tar.gz","targetDir":"~/.openclaw/skills","os":["darwin","linux","win32"]},{"id":"openclaw-runtime-linux","kind":"download","label":"Download md2wechat runtime (Linux amd64)","url":"https://github.com/geekjourneyx/md2wechat-skill/releases/download/v1.11.1/md2wechat-linux-amd64","targetDir":"~/.openclaw/tools/md2wechat","os":["linux"]},{"id":"openclaw-runtime-darwin","kind":"download","label":"Download md2wechat runtime (macOS amd64)","url":"https://github.com/geekjourneyx/md2wechat-skill/releases/download/v1.11.1/md2wechat-darwin-amd64","targetDir":"~/.openclaw/tools/md2wechat","os":["darwin"]}]}}
---

# md2wechat for OpenClaw

透明披露：

- 会读取本地 Markdown 文件和本地图片。
- 可能把处理后的图片和 HTML 上传到微信草稿箱或素材接口。
- 可能调用外部图像服务来生成图片或补图。
- 草稿上传需要 `WECHAT_APPID` 和 `WECHAT_SECRET`。
- 图片生成通常还需要 `IMAGE_API_KEY`，以及可选的 `IMAGE_PROVIDER` / `IMAGE_API_BASE`。

## 运行边界

- 本 skill 只执行已经安装好的 `md2wechat` runtime。
- 优先查找 `~/.openclaw/tools/md2wechat` 下的 runtime，然后查找 `PATH`。
- 运行时不会下载二进制，不会写入缓存下载器，也不会静默回退到远程拉取。
- `metadata.openclaw.install` 当前提供的是固定版本安装资源与 installer 入口，真正完整的自动安装主线仍以 `install-openclaw.sh` 为准。

## 推荐流程

1. 先通过固定版本 OpenClaw installer 完成 skill 和 runtime 安装。
2. 确认 `md2wechat` 可以直接运行。
3. 再执行以下任务：
   - `convert <file.md> --preview`
   - `convert <file.md> --draft --cover <cover.jpg>`
   - `create-image-post --from-markdown <file.md>`
4. 如果要使用 AI 转换或 AI 图片，再补齐图像服务配置。

See [references/runtime.md](references/runtime.md) for the runtime lookup contract.
