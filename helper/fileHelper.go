package helper

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
)

func CreateDirIfNotExists(dir string) error {
	exists, err := IsExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		return os.Mkdir(dir, 0755)
	}
	return nil
}

func IsExists(filename string) (exists bool, err error) {
	_, err = os.Stat(filename)
	if err == nil {
		exists = true
		return
	}
	if os.IsNotExist(err) {
		exists = false
		err = nil
		return
	}
	return false, err
}

type HashInfo struct {
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
	Sha512 string `json:"sha512"`
}

// Calculate hashes
// FROM: http://marcio.io/2015/07/calculating-multiple-file-hashes-in-a-single-pass/
func CalculateBasicHashes(filename string) (info HashInfo, err error) {
	rd, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer rd.Close()

	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	sha512 := sha512.New()

	// For optimum speed, Getpagesize returns the underlying system's memory page size.
	pagesize := os.Getpagesize()

	// wraps the Reader object into a new buffered reader to read the files in chunks
	// and buffering them for performance.
	reader := bufio.NewReaderSize(rd, pagesize)

	// creates a multiplexer Writer object that will duplicate all write
	// operations when copying data from source into all different hashing algorithms
	// at the same time
	multiWriter := io.MultiWriter(md5, sha1, sha256, sha512)

	// Using a buffered reader, this will write to the writer multiplexer
	// so we only traverse through the file once, and can calculate all hashes
	// in a single byte buffered scan pass.
	//
	_, err = io.Copy(multiWriter, reader)
	if err != nil {
		return
	}

	info.Md5 = hex.EncodeToString(md5.Sum(nil))
	info.Sha1 = hex.EncodeToString(sha1.Sum(nil))
	info.Sha256 = hex.EncodeToString(sha256.Sum(nil))
	info.Sha512 = hex.EncodeToString(sha512.Sum(nil))

	return info, nil
}
