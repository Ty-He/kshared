package utils // for controller used

import (
    "os"
    "io"
    "mime/multipart"
    "fmt"
    "sync"
)

// tempFile should allow multiple go visit.
type TempFile struct {
    file *os.File
}

var filename = uint64(0)
var mut sync.Mutex

// if err == nil, should call defer t.file.Close()
func NewTemp() (*TempFile, error) {
    mut.Lock()
    curf := filename 
    filename++
    mut.Unlock()

    // '-' index is 17
    f, err := os.Create(fmt.Sprintf("resource/article/-%d.md", curf))
    if err != nil {
        return nil, err
    }
    return &TempFile{ f }, nil 
}


// copy from src
func (t *TempFile) Copy(src multipart.File) error {
    _, err := io.Copy(t.file, src)
    return err
}

// save temp file 
func (t *TempFile) Save(id uint32) {
    os.Rename(t.file.Name(), fmt.Sprintf("resource/article/%d.md", id))
}


// when close, if file still temp, will be removed
func (t *TempFile) Close() {
    name := t.file.Name()
    if name[17] == '-' {
        os.Remove(name)
    }
    t.file.Close()
}
