package session

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"net/http"

	session "github.com/bhojpur/session/pkg/engine"
)

// FileSessionStore File session store
type FileSessionStore session.FileSessionStore

// Set value to file session
func (fs *FileSessionStore) Set(key, value interface{}) error {
	return (*session.FileSessionStore)(fs).Set(context.Background(), key, value)
}

// Get value from file session
func (fs *FileSessionStore) Get(key interface{}) interface{} {
	return (*session.FileSessionStore)(fs).Get(context.Background(), key)
}

// Delete value in file session by given key
func (fs *FileSessionStore) Delete(key interface{}) error {
	return (*session.FileSessionStore)(fs).Delete(context.Background(), key)
}

// Flush Clean all values in file session
func (fs *FileSessionStore) Flush() error {
	return (*session.FileSessionStore)(fs).Flush(context.Background())
}

// SessionID Get file session store id
func (fs *FileSessionStore) SessionID() string {
	return (*session.FileSessionStore)(fs).SessionID(context.Background())
}

// SessionRelease Write file session to local file with Gob string
func (fs *FileSessionStore) SessionRelease(w http.ResponseWriter) {
	(*session.FileSessionStore)(fs).SessionRelease(context.Background(), w)
}

// FileProvider File session provider
type FileProvider session.FileProvider

// SessionInit Init file session provider.
// savePath sets the session files path.
func (fp *FileProvider) SessionInit(maxlifetime int64, savePath string) error {
	return (*session.FileProvider)(fp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead Read file session by sid.
// if file is not exist, create it.
// the file path is generated from sid string.
func (fp *FileProvider) SessionRead(sid string) (Store, error) {
	s, err := (*session.FileProvider)(fp).SessionRead(context.Background(), sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionExist Check file session exist.
// it checks the file named from sid exist or not.
func (fp *FileProvider) SessionExist(sid string) bool {
	res, _ := (*session.FileProvider)(fp).SessionExist(context.Background(), sid)
	return res
}

// SessionDestroy Remove all files in this save path
func (fp *FileProvider) SessionDestroy(sid string) error {
	return (*session.FileProvider)(fp).SessionDestroy(context.Background(), sid)
}

// SessionGC Recycle files in save path
func (fp *FileProvider) SessionGC() {
	(*session.FileProvider)(fp).SessionGC(context.Background())
}

// SessionAll Get active file session number.
// it walks save path to count files.
func (fp *FileProvider) SessionAll() int {
	return (*session.FileProvider)(fp).SessionAll(context.Background())
}

// SessionRegenerate Generate new sid for file session.
// it delete old file and create new file named from new sid.
func (fp *FileProvider) SessionRegenerate(oldsid, sid string) (Store, error) {
	s, err := (*session.FileProvider)(fp).SessionRegenerate(context.Background(), oldsid, sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}
