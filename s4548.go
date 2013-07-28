package s4548

import (
    "image"
    //"image/draw"
    "image/color"
    "os"
)

var bitmasks = []byte{ 0x80, 0x40, 0x20, 0x10, 0x8, 0x4, 0x2, 0x1 }

type s4548 struct {
        *image.Paletted
        
        fb []uint8
        fp *os.File
        path string
}

func (s *s4548) Scanout() {
    r := image.Rect(0, 0, 101, 40)
    
    for y := r.Min.Y; y < r.Max.Y; y++ {
        row_offset := ((y / 8) * s.Bounds().Max.X)
        for x := r.Min.X; x < r.Max.X; x++ {
            if s.ColorIndexAt(x, y) > 0 {
                s.fb[row_offset + x] |= bitmasks[y % 8]
            } else {
                s.fb[row_offset + x] &= ^bitmasks[y % 8]
            }
        }
    }
    
    n, err := s.fp.WriteAt(s.fb, 0)
    if (err != nil) || (n != len(s.fb)) {
        panic(err)
    }
}

func (s *s4548) getPath() string {
    return s.path
}

func NewS4548(path string) *s4548 {
    
    screen := new(s4548)
    screen.path = path
    
    var err error
    screen.fp = new(os.File)
    screen.fp, err = os.OpenFile(path, os.O_RDWR, os.ModeDevice | os.ModeCharDevice)
    if err != nil { 
        panic(err) 
    }
    
    screen.fb = make([]uint8, 101*5)
    
    p := color.Palette{color.White, color.Black}
    r := image.Rect(0, 0, 101, 40)
    screen.Paletted = image.NewPaletted(r, p)
    
    return screen
}

func (s *s4548) Close() error {
    return s.fp.Close()
}
