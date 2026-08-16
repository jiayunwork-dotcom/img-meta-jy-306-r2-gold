# img-meta

提取图片的几何与格式元数据（基于标准库 `image` 包，支持 PNG/JPEG/GIF，无需外部依赖），按宽高比与分辨率归类，并批量生成重命名方案。

## 用法

```bash
# 预览重命名方案（不实际改名）
img-meta scan --dir example

# 自定义模板
img-meta scan --dir example --template "{tier}-{w}x{h}-{i}"

# 实际执行重命名
img-meta scan --dir example --apply
```

非图片文件会被跳过并提示；目录缺失时返回受控错误。示例图片由 `example/gen.go`（`go run example/gen.go`）生成。

## 构建

```bash
go build ./...
go test ./...
```
