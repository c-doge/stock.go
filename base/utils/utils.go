package utils

import (
	"io"
	"os"
	"fmt"
	"sort"
	"bytes"
	"time"
	"math"
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
    "archive/zip"
    "path/filepath"
	"encoding/binary"
)

var DefaultTimeZone = time.FixedZone("UTC-8", int((8 * time.Hour).Seconds()))
var EarlyDate       = time.Date(1970,  1,  1, 0, 0, 0, 0, DefaultTimeZone);
var FutureDate      = time.Date(2100, 12, 31, 0, 0, 0, 0, DefaultTimeZone);


func DateToTime(date uint32) time.Time {
	y := int(date / 10000);
	m := int((date % 10000) / 100);
	d := int(date % 100);

	ts := time.Date(y, time.Month(m), d, 0, 0, 0, 0, DefaultTimeZone);

	return ts;
}

func TimeToDate(t time.Time) uint32 {
	y := t.Year();
	m := int(t.Month());
	d := t.Day();
	return uint32(y * 10000 + m * 100 + d);
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}


func Mkdir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0777);
		} else {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("Not Direcotory")
	}
	return err;
}


func ListFile(dirPath string) ([]string, error) {
    files := make([]string, 0, 10)
    dir, err := ioutil.ReadDir(dirPath)
    if err != nil {
        return nil, err
    }
    //PathSep := string(os.PathSeparator)
    for _, file := range dir {
        if !file.IsDir() { // 忽略目录
            files = append(files, file.Name())
        }
    }
    sort.SliceStable(files, func(i, j int) bool {
        return strings.Compare(files[i], files[j]) <=0
    })
    return files, nil
}

func UnzipFile(zipFile, dstPath string) error {
	f, err := os.Open(zipFile)
    defer f.Close()
    if err != nil {
        return err;
    }
    fi, err := f.Stat()
    if err != nil {
        return err; 
    }
    
    err = Mkdir(dstPath)
    if err != nil {
    	return err;
    }
   	err = Unzip(f, fi.Size(), dstPath)
    return err
}

func Unzip(r io.ReaderAt, s int64, dstPath string) error {
    reader, err := zip.NewReader(r, s);
    if err != nil {
        return err
    }
    //defer reader.Close()
    for _, f := range reader.File {
        path := filepath.Join(dstPath, f.Name)
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()
        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            dstFile, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, f.Mode())
            if err != nil {
                return err;
            }
            _, err = io.Copy(dstFile, rc);
            if err != nil {
                return err;
            }
        }
    }
    return nil;
}

func Download(url, path string,  unzip bool) error {
    resp, err := http.Get(url);
    if err != nil {
        return err;
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return fmt.Errorf("Http Response %d %s", resp.StatusCode, resp.Status);
    }
    

    if unzip {
    	var fi os.FileInfo;
	    fi, err = os.Stat(path);
	    if err != nil {
	    	return err;
	    } else if !fi.IsDir() {
	    	return fmt.Errorf("path(%s) is not a dir", path);
	    }
    	var buf []byte;
    	buf, err = io.ReadAll(resp.Body)
	    if err == nil {
	        r := bytes.NewReader(buf)
	    	err = Unzip(r, int64(len(buf)), path);
	    }
    } else {
    	var file *os.File = nil;
    	file, err = os.Create(path);
    	if err == nil {
    		_, err = io.Copy(file, resp.Body);
    		file.Close();
    	}
    }
	return err
}

func TrimString(s string) string {
	return strings.Trim(s, " \t\r\n\f\x00")
}


func Float32ToBytes(f float32) []byte{
    bits := math.Float32bits(f)
    b := make([]byte, 4)
    binary.LittleEndian.PutUint32(b, bits)
    return b
}

func Int32ToBytes(i int32) []byte {
    b := bytes.NewBuffer([]byte{});
    binary.Write(b, binary.LittleEndian, i);
    return b.Bytes();
}
func Int16ToBytes(i int16) []byte {
    b := bytes.NewBuffer([]byte{});
    binary.Write(b, binary.LittleEndian, i);
    return b.Bytes();
}

func BytesToUint16(b []byte) uint16 {
    var x uint16;
    byteBuf := bytes.NewBuffer(b)
    binary.Read(byteBuf, binary.LittleEndian, &x);
    return x;
}

func BytesToUint32(b []byte) uint32 {
   var x uint32;
   byteBuf := bytes.NewBuffer(b)
   binary.Read(byteBuf, binary.LittleEndian, &x);
   return x
}

// func BytesToUint32(buf []byte) uint32 {
//     var x uint32;
//     x = (uint32(buf[1]) << 8) | uint32(buf[0]);
//     x = x | (uint32(buf[2]) << 16) | (uint32(buf[3]) << 24);
//     return x;
// }

func BytesToFloat32(b []byte) float32 {
   var x float32;
   byteBuf := bytes.NewBuffer(b)
   binary.Read(byteBuf, binary.LittleEndian, &x);
   return x;
}
