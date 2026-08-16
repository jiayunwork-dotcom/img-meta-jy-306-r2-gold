// Package rename 依据元数据生成批量重命名方案并按需应用。
package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"img-meta/internal/meta"
)

// Rename 是一条重命名记录。
type Rename struct {
	From string
	To   string
}

// Plan 依据模板与分类生成重命名方案；同名冲突以序号消歧。
// 模板支持占位符：{aspect} {tier} {w} {h} {i}。
func Plan(metas []meta.Meta, template string) []Rename {
	if template == "" {
		template = "{aspect}-{w}x{h}-{i}"
	}
	seen := map[string]int{}
	var plans []Rename
	for i, m := range metas {
		class := meta.Classify(m)
		base := render(template, m, class, i)
		ext := strings.TrimPrefix(filepath.Ext(m.Path), ".")
		if ext == "" {
			ext = m.Format
		}
		count := seen[base]
		seen[base] = count + 1
		to := base
		if count > 0 {
			to = fmt.Sprintf("%s_%d", base, count)
		}
		to = to + "." + ext
		plans = append(plans, Rename{From: m.Path, To: to})
	}
	return plans
}

// Apply 执行重命名（原地移动文件）。
func Apply(plans []Rename) error {
	for _, p := range plans {
		if err := os.Rename(p.From, p.To); err != nil {
			return err
		}
	}
	return nil
}

func render(tpl string, m meta.Meta, c meta.Class, i int) string {
	repl := strings.NewReplacer(
		"{aspect}", c.Aspect,
		"{tier}", c.ResTier,
		"{w}", fmt.Sprint(m.Width),
		"{h}", fmt.Sprint(m.Height),
		"{i}", fmt.Sprint(i),
	)
	return repl.Replace(tpl)
}
