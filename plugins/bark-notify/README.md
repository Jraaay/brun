# bark-notify

一个 Claude Code 插件,在 Claude 每次回复结束时通过 [Bark](https://github.com/Finb/Bark) 推送通知到你的 iPhone。

和同仓库的 [`brun`](../../) CLI 是一对姊妹工具:`brun` 管命令行,`bark-notify` 管对话。

## 工作原理

注册一个 `Stop` hook,在 Claude 回复结束时执行 `hooks/notify.sh`:

1. 从 stdin 读取 hook 的 JSON(`session_id` / `cwd` / `stop_hook_active` 等)
2. 如果 `stop_hook_active=true` (Claude 处于强制继续状态),跳过,避免重复推送
3. 组装标题和正文(项目目录名 + session 前 8 位),POST 到 `$BARK_URL`

## 依赖

- `jq` —— 解析 hook 输入
- `curl` —— 发送 HTTP 请求

macOS:`brew install jq curl`

如果缺任一依赖或 `BARK_URL` 未设置,脚本静默退出,不影响 Claude Code 使用。

## 配置

设置环境变量(建议写入 `~/.zshrc`):

```bash
export BARK_URL="https://api.day.app/你的key"
```

## 安装

### 本地测试

```bash
claude --plugin-dir /path/to/brun/plugins/bark-notify
```

### 通过 plugin marketplace

如果你把 `brun` 仓库配置成 Claude Code 的 marketplace,可以直接:

```
/plugin install bark-notify
```

## 自定义

想改通知标题、正文格式或添加更多字段(例如从 transcript 里抓最后一条用户消息),直接编辑 `hooks/notify.sh`。

## License

MIT
