

package main


import "crypto/sha256"
import "encoding/hex"
import "encoding/json"
import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "log"
import "net/http"
import "path/filepath"
import "os"
import "sort"
import "syscall"

// import "github.com/colinmarc/cdb"
import cdb "github.com/cipriancraciun/go-cdb-lib"

import . "github.com/cipriancraciun/go-cdb-http/lib/common"
import . "github.com/cipriancraciun/go-cdb-http/lib/archiver"




type context struct {
	cdbWriter *cdb.Writer
	storedData map[string]bool
	storedFiles map[[2]uint64]string
	compress string
	includeIndex bool
	includeMetadata bool
	debug bool
}




func archiveFile (_context *context, _pathResolved string, _pathInArchive string, _name string) (error) {
	
	var _file *os.File
	if _file_0, _error := os.Open (_pathResolved); _error == nil {
		_file = _file_0
	} else {
		return _error
	}
	
	defer _file.Close ()
	
	var _fileId [2]uint64
	if _stat, _error := _file.Stat (); _error == nil {
		_stat := _stat.Sys()
		if _stat, _ok := _stat.(*syscall.Stat_t); _ok {
			_fileId = [2]uint64 { _stat.Dev, _stat.Ino }
		} else {
			return fmt.Errorf ("[6578d2d7]  failed `stat`-ing:  `%s`!", _pathResolved)
		}
	}
	
	_fingerprint, _wasStored := _context.storedFiles[_fileId]
	
	if _wasStored {
		if _error := archiveReference (_context, NamespaceFilesContent, _pathInArchive, _fingerprint); _error != nil {
			return _error
		}
		return nil
	}
	
	var _data []byte
	if _data_0, _error := ioutil.ReadAll (_file); _error == nil {
		_data = _data_0
	} else {
		return _error
	}
	
	if _fingerprint_0, _error := archiveReferenceAndData (_context, NamespaceFilesContent, _pathResolved, _pathInArchive, _name, _data, ""); _error != nil {
		return _error
	} else {
		_fingerprint = _fingerprint_0
	}
	
	_context.storedFiles[_fileId] = _fingerprint
	
	return nil
}




func archiveFolder (_context *context, _pathResolved string, _pathInArchive string, _names []string, _stats map[string]os.FileInfo) (error) {
	
	type Entry struct {
		Name string `json:"name",omitempty`
		Type string `json:"type",omitempty`
		Size uint64 `json:"size",omitempty"`
	}
	
	type Folder struct {
		Entries []Entry `json:"entries",omitempty`
		Indices []string `json:"indices",omitempty`
	}
	
	_entries := make ([]Entry, 0, len (_names))
	if _context.includeMetadata {
		for _, _name := range _names {
			_entry := Entry {
					Name : _name,
					Type : "unknown",
				}
			_stat := _stats[_name]
			_statMode := _stat.Mode ()
			if _statMode.IsRegular () {
				_entry.Type = "file"
				_entry.Size = uint64 (_stat.Size ())
			} else if _statMode.IsDir () {
				_entry.Type = "folder"
			}
			_entries = append (_entries, _entry)
		}
	}
	
	_indexNames := make ([]string, 0, 4)
	if _context.includeIndex {
		var _indexNameFirst string
		for _, _indexName := range IndexNames {
			_indexNameFound := sort.SearchStrings (_names, _indexName)
			if _indexNameFound == len (_names) {
				continue
			}
			if _names[_indexNameFound] != _indexName {
				continue
			}
			_stat := _stats[_indexName]
			_statMode := _stat.Mode ()
			if ! _statMode.IsRegular () {
				continue
			}
			if _indexNameFirst == "" {
				_indexNameFirst = _indexName
			}
			_indexNames = append (_indexNames, _indexName)
		}
		if _indexNameFirst != "" {
			_indexPathResolved := filepath.Join (_pathResolved, _indexNameFirst)
			_indexPathInArchive := _pathInArchive + "/"
			if _pathInArchive == "/" {
				_indexPathInArchive = "/"
			}
			archiveFile (_context, _indexPathResolved, _indexPathInArchive, _indexNameFirst)
		}
	}
	
	if ! _context.includeMetadata {
		return nil
	}
	
	_folder := Folder {
			Entries : _entries,
			Indices : _indexNames,
		}
	
	var _data []byte
	if _data_0, _error := json.Marshal (&_folder); _error == nil {
		_data = _data_0
	} else {
		return _error
	}
	
	if _, _error := archiveReferenceAndData (_context, NamespaceFoldersContent, _pathResolved, _pathInArchive, "", _data, MimeTypeJson); _error != nil {
		return _error
	}
	
	return nil
}




func archiveReferenceAndData (_context *context, _namespace string, _pathResolved string, _pathInArchive string, _name string, _data []byte, _dataType string) (string, error) {
	
	var _fingerprint string
	var _dataMeta map[string]string
	if _fingerprint_0, _data_0, _dataMeta_0, _error := prepareData (_context, _pathResolved, _pathInArchive, _name, _data, _dataType); _error != nil {
		return "", _error
	} else {
		_fingerprint = _fingerprint_0
		_data = _data_0
		_dataMeta = _dataMeta_0
	}
	
	if _error := archiveReference (_context, _namespace, _pathInArchive, _fingerprint); _error != nil {
		return "", _error
	}
	
	if (_data != nil) && (_dataMeta != nil) {
		if _error := archiveData (_context, _fingerprint, _data, _dataMeta); _error != nil {
			return "", _error
		}
	}
	
	return _fingerprint, nil
}



func archiveData (_context *context, _fingerprint string, _data []byte, _dataMeta map[string]string) (error) {
	
	if _wasStored, _ := _context.storedData[_fingerprint]; _wasStored {
		return fmt.Errorf ("[256cde78]  data already stored:  `%s`!", _fingerprint)
	}
	
	var _dataMetaRaw []byte
	if _dataMetaRaw_0, _error := MetadataEncode (_dataMeta); _error == nil {
		_dataMetaRaw = _dataMetaRaw_0
	} else {
		return _error
	}
	
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataMetadata, _fingerprint)
		if _context.debug {
			log.Printf ("[  ] blob-meta    ++ `%s`\n", _key)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataMetaRaw); _error != nil {
			return _error
		}
	}
	
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataContent, _fingerprint)
		if _context.debug {
			log.Printf ("[  ] blob-data    ++ `%s`\n", _key)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _data); _error != nil {
			return _error
		}
	}
	
	_context.storedData[_fingerprint] = true
	
	return nil
}




func archiveReference (_context *context, _namespace string, _pathInArchive string, _fingerprint string) (error) {
	
	_key := fmt.Sprintf ("%s:%s", _namespace, _pathInArchive)
	if _context.debug {
		log.Printf ("[  ] reference    ++ `%s` :: `%s` -> `%s`\n", _namespace, _pathInArchive, _fingerprint)
	}
	if _error := _context.cdbWriter.Put ([]byte (_key), []byte (_fingerprint)); _error != nil {
		return _error
	}
	
	return nil
}




func prepareData (_context *context, _pathResolved string, _pathInArchive string, _name string, _data []byte, _dataType string) (string, []byte, map[string]string, error) {
	
	_fingerprintRaw := sha256.Sum256 (_data)
	_fingerprint := hex.EncodeToString (_fingerprintRaw[:])
	
	if _wasStored, _ := _context.storedData[_fingerprint]; _wasStored {
		return _fingerprint, nil, nil, nil
	}
	
	if (_dataType == "") && (_name != "") {
		_extension := filepath.Ext (_name)
		if _extension != "" {
			_extension = _extension[1:]
		}
		_dataType, _ = MimeTypesByExtension[_extension]
	}
	if _dataType == "" {
		_dataType = http.DetectContentType (_data)
	}
	if _dataType == "" {
		_dataType = MimeTypeRaw
	}
	
	_dataEncoding := "identity"
	_dataUncompressedSize := len (_data)
	_dataSize := _dataUncompressedSize
	if _dataSize > 512 {
		if _data_0, _dataEncoding_0, _error := Compress (_data, _context.compress); _error == nil {
			if _dataEncoding_0 != "identity" {
				_dataCompressedSize := len (_data_0)
				_dataCompressedDelta := _dataUncompressedSize - _dataCompressedSize
				_dataCompressedRatio := (_dataCompressedDelta * 100) / _dataUncompressedSize
				_accepted := false
				_accepted = _accepted || ((_dataUncompressedSize > (1024 * 1024)) && (_dataCompressedRatio >= 5))
				_accepted = _accepted || ((_dataUncompressedSize > (64 * 1024)) && (_dataCompressedRatio >= 10))
				_accepted = _accepted || ((_dataUncompressedSize > (32 * 1024)) && (_dataCompressedRatio >= 15))
				_accepted = _accepted || ((_dataUncompressedSize > (16 * 1024)) && (_dataCompressedRatio >= 20))
				_accepted = _accepted || ((_dataUncompressedSize > (8 * 1024)) && (_dataCompressedRatio >= 25))
				_accepted = _accepted || ((_dataUncompressedSize > (4 * 1024)) && (_dataCompressedRatio >= 30))
				_accepted = _accepted || ((_dataUncompressedSize > (2 * 1024)) && (_dataCompressedRatio >= 35))
				_accepted = _accepted || ((_dataUncompressedSize > (1 * 1024)) && (_dataCompressedRatio >= 40))
				_accepted = _accepted || (_dataCompressedRatio >= 90)
				if _accepted {
					_data = _data_0
					_dataEncoding = _dataEncoding_0
					_dataSize = _dataCompressedSize
				}
				if _dataSize < _dataUncompressedSize {
					if _context.debug {
						log.Printf ("[  ] compress          %02d %8d %8d `%s`\n", _dataCompressedRatio, _dataUncompressedSize, _dataCompressedDelta, _pathInArchive)
					}
				} else {
					if _context.debug {
						log.Printf ("[  ] compress-NOK      %02d %8d %8d `%s`\n", _dataCompressedRatio, _dataUncompressedSize, _dataCompressedDelta, _pathInArchive)
					}
				}
			}
		} else {
			return "", nil, nil, _error
		}
	} else {
		if _context.debug && (_context.compress != "identity") {
			log.Printf ("[  ] compress-NOK         %8d          `%s`\n", _dataUncompressedSize, _pathInArchive)
		}
	}
	
	_dataMeta := make (map[string]string, 16)
	_dataMeta["Content-Length"] = fmt.Sprintf ("%d", _dataSize)
	_dataMeta["Content-Type"] = _dataType
	_dataMeta["Content-Encoding"] = _dataEncoding
	_dataMeta["ETag"] = _fingerprint
	
	return _fingerprint, _data, _dataMeta, nil
}




func walkPath (_context *context, _path string, _prefix string, _name string, _recursed map[string]uint, _recurse bool) (os.FileInfo, error) {
	
	if _recursed == nil {
		_recursed = make (map[string]uint, 128)
	}
	
	_pathInArchive := filepath.Join (_prefix, _name)
	
	var _stat os.FileInfo
	if _stat_0, _error := os.Lstat (_path); _error == nil {
		_stat = _stat_0
	} else {
		return nil, _error
	}
	_statMode := _stat.Mode ()
	
	_isSymlink := false
	if (_stat.Mode () & os.ModeSymlink) != 0 {
		_isSymlink = true
		if _stat_0, _error := os.Stat (_path); _error == nil {
			_stat = _stat_0
		} else {
			return nil, _error
		}
	}
	_statMode = _stat.Mode ()
	
	if ! _recurse {
		return _stat, nil
	}
	
	var _pathResolved string
	if _isSymlink {
		if _resolved, _error := filepath.EvalSymlinks (_path); _error == nil {
			_pathResolved = _resolved
		} else {
			return nil, _error
		}
	} else {
		_pathResolved = _path
	}
	
	if _isSymlink && _context.debug {
		log.Printf ("[  ] symlink      :: `%s` -> `%s`\n", _pathInArchive, _pathResolved)
	}
	
	
	if _statMode.IsRegular () {
		
		if _context.debug {
			log.Printf ("[  ] file         :: `%s`\n", _pathInArchive)
		}
		if _error := archiveFile (_context, _pathResolved, _pathInArchive, _name); _error != nil {
			return _stat, _error
		}
		return _stat, nil
		
	} else if _statMode.IsDir () {
		
		_wasRecursed, _ := _recursed[_pathResolved]
		if _wasRecursed > 0 {
			log.Printf ("[ww] [2e1744c9]  detected directory loop for `%s` resolving to `%s`;  ignoring!\n", _path, _pathResolved)
			return _stat, nil
		}
		_recursed[_pathResolved] = _wasRecursed + 1
		
		_names := make ([]string, 0, 16)
		_stats := make (map[string]os.FileInfo, 16)
		
		if _context.debug {
			log.Printf ("[  ] folder       >> `%s`\n", _pathInArchive)
		}
		if _stream, _error := os.Open (_path); _error == nil {
			defer _stream.Close ()
			_prefix = filepath.Join (_prefix, _name)
			_loop : for {
				switch _buffer, _error := _stream.Readdir (128); _error {
					case nil :
						for _, _stat := range _buffer {
							_name := _stat.Name ()
							_names = append (_names, _name)
							if _stat, _error := walkPath (_context, filepath.Join (_path, _name), _prefix, _name, _recursed, false); _error == nil {
								_stats[_name] = _stat
							} else {
								return nil, _error
							}
						}
					case io.EOF :
						break _loop
					default :
						return nil, _error
				}
			}
		}
		if _context.debug {
			log.Printf ("[  ] folder       << `%s`\n", _pathInArchive)
		}
		
		sort.Strings (_names)
		
		if _context.debug {
			log.Printf ("[  ] folder       :: `%s`\n", _pathInArchive)
		}
		if _error := archiveFolder (_context, _pathResolved, _pathInArchive, _names, _stats); _error != nil {
			return nil, _error
		}
		
		if _context.debug {
			log.Printf ("[  ] folder       >> `%s`\n", _pathInArchive)
		}
		
		for _, _name := range _names {
			_stat := _stats[_name]
			_statMode := _stat.Mode ()
			if _statMode.IsRegular () {
				if _, _error := walkPath (_context, filepath.Join (_path, _name), _prefix, _name, _recursed, true); _error != nil {
					return _stat, _error
				}
			}
		}
		for _, _name := range _names {
			_stat := _stats[_name]
			_statMode := _stat.Mode ()
			if _statMode.IsDir () {
				if _, _error := walkPath (_context, filepath.Join (_path, _name), _prefix, _name, _recursed, true); _error != nil {
					return _stat, _error
				}
			}
		}
		for _, _name := range _names {
			_stat := _stats[_name]
			_statMode := _stat.Mode ()
			if ! _statMode.IsRegular () && ! _statMode.IsDir () {
				if _, _error := walkPath (_context, filepath.Join (_path, _name), _prefix, _name, _recursed, true); _error != nil {
					return _stat, _error
				}
			}
		}
		if _context.debug {
			log.Printf ("[  ] folder       << `%s`\n", _pathInArchive)
		}
		
		_recursed[_pathResolved] = _wasRecursed
		return _stat, nil
		
	} else {
		return nil, fmt.Errorf ("[d9b836d7]  unexpected file type for `%s`:  `%s`!", _path, _statMode)
	}
}




func main () () {
	Main (main_0)
}


func main_0 () (error) {
	
	
	var _sourcesFolder string
	var _archiveFile string
	var _compress string
	var _includeIndex bool
	var _includeMetadata bool
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("cdb-http-archiver", flag.ContinueOnError)
		
		_flags.Usage = func () () {
			fmt.Fprintf (os.Stderr, "%s",
`
cdb-http-archiver
	--sources <path>
	--archive <path>
	--compress <gzip | brotli | identity>
	--exclude-index
	--include-metadata
	--debug
`)
		}
		
		_sourcesFolder_0 := _flags.String ("sources", "", "")
		_archiveFile_0 := _flags.String ("archive", "", "")
		_compress_0 := _flags.String ("compress", "", "")
		_excludeIndex_0 := _flags.Bool ("exclude-index", false, "")
		_includeMetadata_0 := _flags.Bool ("include-metadata", false, "")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_sourcesFolder = *_sourcesFolder_0
		_archiveFile = *_archiveFile_0
		_compress = *_compress_0
		_includeIndex = ! *_excludeIndex_0
		_includeMetadata = *_includeMetadata_0
		_debug = *_debug_0
		
		if _sourcesFolder == "" {
			AbortError (nil, "[515ee462]  expected sources folder argument!")
		}
		if _archiveFile == "" {
			AbortError (nil, "[5e8da985]  expected archive file argument!")
		}
	}
	
	
	var _cdbWriter *cdb.Writer
	if _cdbWriter_0, _error := cdb.Create (_archiveFile); _error == nil {
		_cdbWriter = _cdbWriter_0
	} else {
		AbortError (_error, "[85234ba0]  failed creating archive (while opening)!")
	}
	
	_context := & context {
			cdbWriter : _cdbWriter,
			storedData : make (map[string]bool, 16 * 1024),
			storedFiles : make (map[[2]uint64]string, 16 * 1024),
			compress : _compress,
			includeIndex : _includeIndex,
			includeMetadata : _includeMetadata,
			debug : _debug,
		}
	
	if _, _error := walkPath (_context, _sourcesFolder, "/", "", nil, true); _error != nil {
		AbortError (_error, "[b6a19ef4]  failed walking folder!")
	}
	
	if _error := _cdbWriter.Close (); _error != nil {
		AbortError (_error, "[bbfb8478]  failed creating archive (while closing)!")
	}
	
	
	return nil
}

