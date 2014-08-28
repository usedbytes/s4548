package s4548

import (
    "image"
    "image/color"
    "os"
    "syscall"
    "unsafe"
    "fmt"
)

var pBITMASKS = []byte{ 0x80, 0x40, 0x20, 0x10, 0x8, 0x4, 0x2, 0x1 }

const default_screen = "/dev/s4548-0"
const WIDTH = 101
const HEIGHT = 40

type S4548 struct {
        *image.Paletted
        damage image.Rectangle 
        fb []uint8
        fp *os.File
        path string
}

func (s *S4548) Scanout() {
    s.damage = image.Rect(0, 0, WIDTH, HEIGHT)
    s.Repair()
}

func (s *S4548) Repair() {
    r := s.Bounds().Intersect(s.damage)
    //image.Rect(0, 0, WIDTH, HEIGHT)
    fmt.Println("Screen Damage:", r)
    for y := r.Min.Y; y < r.Max.Y; y++ {
        row_offset := ((y / 8) * s.Bounds().Max.X)
        for x := r.Min.X; x < r.Max.X; x++ {
            if s.ColorIndexAt(x, y) > 0 {
                s.fb[row_offset + x] |= pBITMASKS[y % 8]
            } else {
                s.fb[row_offset + x] &= ^pBITMASKS[y % 8]
            }
        }
    }
    s.damage = image.ZR
    msync(s.fb, syscall.MS_ASYNC)
   /* 
    n, err := s.fp.WriteAt(s.fb, 0)
    if (err != nil) || (n != len(s.fb)) {
        panic(err)
    }
    */
}

func (s *S4548) GetPath() string {
    return s.path
}

func NewS4548(path string) *S4548 {
    
    screen := new(S4548)
    screen.path = path
    
    var err error
    screen.fp = new(os.File)
    screen.fp, err = os.OpenFile(path, os.O_RDWR, os.ModeDevice | os.ModeCharDevice)
    if err != nil { 
        panic(err) 
    }
    
    //screen.fb = make([]uint8, WIDTH*(HEIGHT / 8))
    screen.fb, err = syscall.Mmap(int(screen.fp.Fd()), 0, WIDTH*(HEIGHT / 8), 
        syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED);
    if (err != nil) {
        panic(err)
    }

    p := color.Palette{color.White, color.Black}
    r := image.Rect(0, 0, WIDTH, HEIGHT)
    screen.Paletted = image.NewPaletted(r, p)
    
    return screen
}

func GetS4548EnvPath() string {
    env := os.Getenv("S4548");
    if (env != "") {
        return env;
    } else {
        return default_screen;
    }
}

func (s *S4548) Close() error {
    return s.fp.Close()
}

func (s *S4548) Width() int {
    return WIDTH
}

func (s *S4548) Height() int {
    return HEIGHT
}

func (s *S4548) Damage(r image.Rectangle) {
    if (r != image.ZR) {
        if (s.damage == image.ZR) {
            s.damage = r
        } else {
            s.damage = s.damage.Union(r)
        }
    }
}
/*
func (s *S4548) Draw(dst draw.Image, r image.Rectangle, src image.Image, 
    sp image.Point, op draw.Op) {
    s.damage = s.damage.Union(r)
    draw.Draw(s.Paletted, r, src, sp, op)
}
*/

func msync(b []byte, flag int) (err error) {
    var _p0 unsafe.Pointer
    if len(b) > 0 {
        _p0 = unsafe.Pointer(&b[0])
    } else {
        _p0 = unsafe.Pointer(&b[0])
    }
    _, _, e1 := syscall.Syscall(syscall.SYS_MSYNC,
                    uintptr(_p0), uintptr(len(b)), uintptr(flag))
    if e1 != 0 {
        err = e1
    }
    return
}
