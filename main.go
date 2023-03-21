package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

func main() {

	ecb()
	cbc()
}

func ecb() {
	key := []byte("your-secret-key-") // AES-128: 16バイト

	// ディレクトリパス
	dirPath := "./"

	// 入力ファイル名
	inputFilename := filepath.Join(dirPath, "input.png")

	// 出力ファイル名
	outputFilename := filepath.Join(dirPath, "ecb_encrypted.png")

	// 画像ファイルを読み込む
	imgFile, err := os.Open(inputFilename)
	if err != nil {
		fmt.Printf("入力ファイルの読み取りに失敗: %v\n", err)
		return
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		fmt.Printf("画像のデコードに失敗: %v\n", err)
		return
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 暗号化後の画像を作成
	encryptedImg := image.NewRGBA(bounds)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("ブロック暗号の作成に失敗: %v\n", err)
		return
	}

	// 画像データを暗号化する
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// ゼロパディング
			rgb := []byte{byte(r), byte(g), byte(b), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			encryptECB(block, rgb)

			encryptedImg.Set(x, y, color.RGBA{rgb[0], rgb[1], rgb[2], byte(a)})
		}
	}

	// 暗号化されたデータをファイルに書き込む
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("出力ファイルの作成に失敗: %v\n", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, encryptedImg)
	if err != nil {
		fmt.Printf("暗号化された画像のエンコードに失敗: %v\n", err)
		return
	}

	fmt.Println("ECBの暗号化完了")
}

func encryptECB(block cipher.Block, buf []byte) {
	block.Encrypt(buf, buf)
}

func cbc() {
	key := []byte("your-secret-key-") // AES-128: 16バイト
	iv := []byte("your-iv-16-bytes")  // 16バイトの初期化ベクター

	// ディレクトリパス
	dirPath := "./"

	// 入力ファイル名
	inputFilename := filepath.Join(dirPath, "input.png")

	// 出力ファイル名
	outputFilename := filepath.Join(dirPath, "cbc_encrypted.png")

	// 画像ファイルを読み込む
	imgFile, err := os.Open(inputFilename)
	if err != nil {
		fmt.Printf("入力ファイルの読み取りに失敗: %v\n", err)
		return
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		fmt.Printf("画像のデコードに失敗: %v\n", err)
		return
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 暗号化後の画像を作成
	encryptedImg := image.NewRGBA(bounds)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("ブロック暗号の作成に失敗: %v\n", err)
		return
	}

	// CBCモードで暗号化する
	mode := cipher.NewCBCEncrypter(block, iv)

	// 画像データを暗号化する
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Zero-padding
			rgb := []byte{byte(r), byte(g), byte(b), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			mode.CryptBlocks(rgb, rgb)

			encryptedImg.Set(x, y, color.RGBA{rgb[0], rgb[1], rgb[2], byte(a)})
		}
	}

	// 暗号化されたデータをファイルに書き込む
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("出力ファイルの作成に失敗: %v\n", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, encryptedImg)
	if err != nil {
		fmt.Printf("暗号化された画像のエンコードに失敗: %v\n", err)
		return
	}

	fmt.Println("CBCの暗号化完了")
}
