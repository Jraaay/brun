# brun

一个极简的命令包装器：运行你的命令，结束时通过 [Bark](https://github.com/Finb/Bark) 推送通知到 iPhone。

适合跑耗时任务（编译、训练、备份、批处理……）时，离开电脑也能第一时间知道成功还是失败。

## 特性

- 透传 stdin / stdout / stderr，和直接运行命令一样
- 捕获退出码，成功/失败用不同标题推送
- 显示命令内容、耗时，失败时附带退出码
- 转发 `SIGINT` / `SIGTERM`, `Ctrl+C` 行为正常
- 零依赖（仅 Go 标准库）

## 安装

```bash
go install github.com/Jraaay/brun@latest
```

或从源码构建：

```bash
git clone https://github.com/Jraaay/brun.git
cd brun
go build -o brun
```

然后把生成的 `brun` 放到 `$PATH` 中。

## 配置

设置环境变量 `BARK_URL` 为你的 Bark 推送地址：

```bash
export BARK_URL="https://api.day.app/你的key"
```

建议写入 `~/.zshrc` 或 `~/.bashrc`。

## 使用

在原命令前加上 `brun` 即可：

```bash
brun make build
brun python train.py --epochs 100
brun rsync -av /data/ backup:/data/
```

命令结束后，手机会收到类似下面的通知：

- ✅ 成功：**命令运行成功** — `make build` / 耗时: 1m23s
- ❌ 失败：**命令运行失败** — `python train.py` / 耗时: 5.2s / 退出码: 1

## 退出码

`brun` 的退出码与被包装命令一致：

- 命令正常结束 → 返回命令自身的退出码
- 命令启动失败 → 返回 `127`
- 其它等待错误 → 返回 `1`

## License

MIT
