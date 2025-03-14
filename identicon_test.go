package identicon

import (
	"bytes"
	"image/png"
	"testing"
)

func TestNewInitializesCorrectly(t *testing.T) {
	size := 300
	rows := 5
	cols := 5

	icon := New(size, rows, cols)

	if icon.maxX != size {
		t.Errorf("Expected maxX to be %d, got %d", size, icon.maxX)
	}
	if icon.maxY != size {
		t.Errorf("Expected maxY to be %d, got %d", size, icon.maxY)
	}
	if icon.rows != rows {
		t.Errorf("Expected rows to be %d, got %d", rows, icon.rows)
	}
	if icon.cols != cols {
		t.Errorf("Expected cols to be %d, got %d", cols, icon.cols)
	}
}

func TestRenderProducesImage(t *testing.T) {
	icon := New(300, 5, 5)
	data := []byte("test data")

	img := icon.Make(data)

	if img == nil {
		t.Fatal("Expected non-nil image, got nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 300 || bounds.Dy() != 300 {
		t.Errorf("Expected image size 300x300, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderWithEmptyData(t *testing.T) {
	icon := New(300, 5, 5)
	data := []byte{}

	img := icon.Make(data)

	if img == nil {
		t.Fatal("Expected non-nil image even with empty data, got nil")
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	icon := New(300, 5, 5)
	data := []byte("test data")

	img1 := icon.Make(data)
	img2 := icon.Make(data)

	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)

	png.Encode(buf1, img1)
	png.Encode(buf2, img2)

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Error("Expected identical images for the same input data")
	}
}

func TestDifferentInputsProduceDifferentImages(t *testing.T) {
	icon := New(300, 5, 5)
	data1 := []byte("test data 1")
	data2 := []byte("test data 2")

	img1 := icon.Make(data1)
	img2 := icon.Make(data2)

	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)

	png.Encode(buf1, img1)
	png.Encode(buf2, img2)

	if bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Error("Expected different images for different input data")
	}
}

func TestRenderWithDifferentDimensions(t *testing.T) {
	sizes := []int{100, 200, 300}
	for _, size := range sizes {
		t.Run("Size "+string(rune(size)), func(t *testing.T) {
			icon := New(size, 3, 3)
			img := icon.Make([]byte("test"))

			bounds := img.Bounds()
			if bounds.Dx() != size || bounds.Dy() != size {
				t.Errorf("Expected image size %dx%d, got %dx%d", size, size, bounds.Dx(), bounds.Dy())
			}
		})
	}
}
