/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package blockfilemgr

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
)

// ErrUnexpectedEndOfBlockfile error used to indicate an unexpected end of a file segment
// this can happen mainly if a crash occurs during appening a block and partial block contents
// get written towards the end of the file
var ErrUnexpectedEndOfBlockfile = errors.New("unexpected end of blockfile")

// blockfileStream reads blocks sequentially from a single file.
// It starts from the given offset and can traverse till the end of the file
type blockfileStream struct {
	fileNum       int
	reader        *bufio.Reader
	currentOffset int64
	fileSize int64
}

const (
	blockfilePrefix = "blockfile_"
)

var (
	blkMgrInfoKey = []byte("blkMgrInfo")
)

// blockStream reads blocks sequentially from multiple files.
// it starts from a given file offset and continues with the next
// file segment until the end of the last segment (`endFileNum`)
type BlockStream struct {
	rootDir           string
	currentFileNum    int
	endFileNum        int
	currentFileStream *blockfileStream
}

// blockPlacementInfo captures the information related
// to block's placement in the file.
type blockPlacementInfo struct {
	fileNum          int
	blockStartOffset int64
	blockBytesOffset int64
}

func deriveBlockfilePath(rootDir string, suffixNum int) string {
	return rootDir + "/" + blockfilePrefix + fmt.Sprintf("%06d", suffixNum)
}

func (s *blockfileStream) nextBlockBytes() ([]byte, error) {
	blockBytes, _, err := s.nextBlockBytesAndPlacementInfo()
	return blockBytes, err
}

// nextBlockBytesAndPlacementInfo returns bytes for the next block
// along with the offset information in the block file.
// An error `ErrUnexpectedEndOfBlockfile` is returned if a partial written data is detected
// which is possible towards the tail of the file if a crash had taken place during appending of a block
func (s *blockfileStream) nextBlockBytesAndPlacementInfo() ([]byte, *blockPlacementInfo, error) {
	var lenBytes []byte
	var err error
	moreContentAvailable := true

	if s.currentOffset == s.fileSize {
		fmt.Printf("Finished reading file number [%d]", s.fileNum)
		return nil, nil, nil
	}
	remainingBytes := s.fileSize - s.currentOffset
	// Peek 8 or smaller number of bytes (if remaining bytes are less than 8)
	// Assumption is that a block size would be small enough to be represented in 8 bytes varint
	peekBytes := 8
	if remainingBytes < int64(peekBytes) {
		peekBytes = int(remainingBytes)
		moreContentAvailable = false
	}
	//fmt.Printf("Remaining bytes=[%d], Going to peek [%d] bytes", remainingBytes, peekBytes)
	if lenBytes, err = s.reader.Peek(peekBytes); err != nil {
		return nil, nil, err
	}
	length, n := proto.DecodeVarint(lenBytes)
	if n == 0 {
		// proto.DecodeVarint did not consume any byte at all which means that the bytes
		// representing the size of the block are partial bytes
		if !moreContentAvailable {
			return nil, nil, ErrUnexpectedEndOfBlockfile
		}
		panic(fmt.Errorf("Error in decoding varint bytes [%#v]", lenBytes))
	}
	bytesExpected := int64(n) + int64(length)
	if bytesExpected > remainingBytes {
		fmt.Printf("At least [%d] bytes expected. Remaining bytes = [%d]. Returning with error [%s]",
			bytesExpected, remainingBytes, ErrUnexpectedEndOfBlockfile)
		return nil, nil, ErrUnexpectedEndOfBlockfile
	}
	// skip the bytes representing the block size
	if _, err = s.reader.Discard(n); err != nil {
		return nil, nil, err
	}
	blockBytes := make([]byte, length)
	if _, err = io.ReadAtLeast(s.reader, blockBytes, int(length)); err != nil {
		fmt.Printf("Error while trying to read [%d] bytes from fileNum [%d]: %s", length, s.fileNum, err)
		return nil, nil, err
	}
	blockPlacementInfo := &blockPlacementInfo{
		fileNum:          s.fileNum,
		blockStartOffset: s.currentOffset,
		blockBytesOffset: s.currentOffset + int64(n)}
	s.currentOffset += int64(n) + int64(length)
	//fmt.Printf("Returning blockbytes - length=[%d], placementInfo={%s}", len(blockBytes), blockPlacementInfo)
	return blockBytes, blockPlacementInfo, nil
}


///////////////////////////////////
// blockStream functions
////////////////////////////////////
func NewBlockStream(ledgerBytes []byte) (*BlockStream, error) {
	bufReader := bytes.NewReader(ledgerBytes)
	startFileStream := &blockfileStream{0, bufio.NewReader(bufReader), 0, int64(len(ledgerBytes))}
	return &BlockStream{"", 0, 0, startFileStream}, nil
}

func (s *BlockStream) NextBlockBytes() ([]byte, error) {
	blockBytes, _, err := s.nextBlockBytesAndPlacementInfo()
	return blockBytes, err
}

func (s *BlockStream) nextBlockBytesAndPlacementInfo() ([]byte, *blockPlacementInfo, error) {
	var blockBytes []byte
	var blockPlacementInfo *blockPlacementInfo
	var err error
	if blockBytes, blockPlacementInfo, err = s.currentFileStream.nextBlockBytesAndPlacementInfo(); err != nil {
		fmt.Printf("current file [%d] length of blockbytes [%d]. Err:%s", s.currentFileNum, len(blockBytes), err)
		return nil, nil, err
	}
	return blockBytes, blockPlacementInfo, nil
}

func (i *blockPlacementInfo) String() string {
	return fmt.Sprintf("fileNum=[%d], startOffset=[%d], bytesOffset=[%d]",
		i.fileNum, i.blockStartOffset, i.blockBytesOffset)
}
