package utils

import (
	"image"
	"image/jpeg" // 导入jpeg包用于JPEG编码
	"mime/multipart"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

const maxWidth = uint(800)

// 可以调整jpegQuality值以控制JPEG图片的质量，范围是 1 到 100。
const jpegQuality = 80 // 示例中设置为80

// CompressAndConvertToJPEG 将上传的图片压缩并转换为 JPEG 格式
func CompressAndConvertToJPEG(fileHeader *multipart.FileHeader, filePath string) error {
	// 1. 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// 2. 自动解码图像（支持 JPEG, PNG, GIF, BMP, TIFF 等）
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 3. 获取原始尺寸
	bounds := img.Bounds()
	width := uint(bounds.Dx())

	// 4. 如果宽度超过限制，则缩放
	if width > maxWidth {
		height := uint(bounds.Dy())
		newHeight := (maxWidth * height) / width
		img = resize.Resize(maxWidth, newHeight, img, resize.Lanczos3)
	}

	// 5. 创建保存目录
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. 创建输出文件
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 7. 编码为 JPEG（有损压缩，质量由jpegQuality控制）
	options := &jpeg.Options{Quality: jpegQuality} // 设置JPEG质量

	err = jpeg.Encode(out, img, options)
	if err != nil {
		return err
	}

	return nil
}
