package meta

// Class 是图片的分类结果。
type Class struct {
	Aspect  string // landscape / portrait / square
	ResTier string // low / hd / fullhd / 4k
}

// Classify 按宽高比与分辨率（百万像素）对图片归类。
func Classify(m Meta) Class {
	var c Class
	switch {
	case m.Width > m.Height:
		c.Aspect = "landscape"
	case m.Width < m.Height:
		c.Aspect = "portrait"
	default:
		c.Aspect = "square"
	}
	mp := m.Width * m.Height / 1_000_000
	switch {
	case mp < 1:
		c.ResTier = "low"
	case mp < 2:
		c.ResTier = "hd"
	case mp < 8:
		c.ResTier = "fullhd"
	default:
		c.ResTier = "4k"
	}
	return c
}
