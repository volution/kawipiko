

package archiver


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
import "strings"
import "syscall"
import "time"

import "github.com/colinmarc/cdb"
import "github.com/zeebo/blake3"
import "go.etcd.io/bbolt"

import . "github.com/volution/kawipiko/lib/common"
import . "github.com/volution/kawipiko/lib/archiver"




type context struct {
	cdbWriter *cdb.Writer
	cdbWriteCount int
	cdbWriteKeySize int
	cdbWriteDataSize int
	storedFilePaths []string
	storedFolderPaths []string
	storedDataMeta map[string]bool
	storedDataContent map[string]bool
	storedDataContentMeta map[string]map[string]string
	storedFiles map[string][2]string
	storedKeys map[string]string
	archivedReferences uint
	compress string
	compressCache *bbolt.DB
	sourcesCache *bbolt.DB
	dataUncompressedCount int
	dataUncompressedSize int
	dataCompressedCount int
	dataCompressedSize int
	includeIndex bool
	includeStripped bool
	includeCache bool
	includeEtag bool
	includeFileListing bool
	includeFolderListing bool
	progress bool
	progressStarted time.Time
	progressLast time.Time
	debug bool
}




func archiveFile (_context *context, _pathResolved string, _pathInArchive string, _name string) (error) {
	
	var _fileDev uint64
	var _fileInode uint64
	var _fileSize uint64
	var _fileTimestamp [2]uint64
	if _stat, _error := os.Stat (_pathResolved); _error == nil {
		_stat := _stat.Sys()
		if _stat, _ok := _stat.(*syscall.Stat_t); _ok {
			_fileDev = uint64 (_stat.Dev)
			_fileInode = uint64 (_stat.Ino)
			_fileSize = uint64 (_stat.Size)
			_fileTimestamp = [2]uint64 { uint64 (_stat.Mtim.Sec), uint64 (_stat.Mtim.Nsec) }
		} else {
			return fmt.Errorf ("[6578d2d7]  failed `stat`-ing:  `%s`!", _pathResolved)
		}
	} else {
		return _error
	}
	
	_fileIdText := fmt.Sprintf ("%d.%d-%d-%d.%d", _fileDev, _fileInode, _fileSize, _fileTimestamp[0], _fileTimestamp[1])
	_fileIdRaw := blake3.Sum256 ([]byte (_fileIdText))
	_fileId := hex.EncodeToString (_fileIdRaw[:])
	
	var _wasStored bool
	var _fingerprintContent string
	var _fingerprintMeta string
	if _fingerprints, _wasStored_0 := _context.storedFiles[_fileId]; _wasStored_0 {
		_fingerprintContent = _fingerprints[0]
		_fingerprintMeta = _fingerprints[1]
		_wasStored = true
	}
	
	var _dataContent []byte
	var _dataMeta map[string]string
	var _dataMetaRaw []byte
	
	if ! _wasStored {
		
		_dataContentRead := func () ([]byte, error) {
				if _context.debug {
					log.Printf ("[dd] [30ef6c2f]  file         <= `%s`\n", _pathInArchive)
				}
				var _file *os.File
				if _file_0, _error := os.Open (_pathResolved); _error == nil {
					_file = _file_0
				} else {
					return nil, _error
				}
				defer _file.Close ()
				if _stat, _error := _file.Stat (); _error == nil {
					_stat := _stat.Sys()
					if _stat, _ok := _stat.(*syscall.Stat_t); _ok {
						if (
								(_fileDev != uint64 (_stat.Dev)) ||
								(_fileInode != uint64 (_stat.Ino)) ||
								(_fileSize != uint64 (_stat.Size)) ||
								(_fileTimestamp[0] != uint64 (_stat.Mtim.Sec)) ||
								(_fileTimestamp[1] != uint64 (_stat.Mtim.Nsec))) {
							return nil, fmt.Errorf ("[3a07643b]  file changed while reading:  `%s`!", _pathResolved)
						}
					} else {
						return nil, fmt.Errorf ("[4daf593a]  failed `stat`-ing:  `%s`!", _pathResolved)
					}
				} else {
					return nil, _error
				}
				var _data []byte
				if _data_0, _error := ioutil.ReadAll (_file); _error == nil {
					_data = _data_0
				} else {
					return nil, _error
				}
				if _stat, _error := _file.Stat (); _error == nil {
					_stat := _stat.Sys()
					if _stat, _ok := _stat.(*syscall.Stat_t); _ok {
						if (
								(_fileSize != uint64 (_stat.Size)) ||
								(_fileTimestamp[0] != uint64 (_stat.Mtim.Sec)) ||
								(_fileTimestamp[1] != uint64 (_stat.Mtim.Nsec))) {
							return nil, fmt.Errorf ("[9689146e]  file changed while reading:  `%s`!", _pathResolved)
						}
					} else {
						return nil, fmt.Errorf ("[523fa3d1]  failed `stat`-ing:  `%s`!", _pathResolved)
					}
				} else {
					return nil, _error
				}
				return _data, nil
			}
		
		if _fingerprintContent_0, _dataContent_0, _dataMeta_0, _error := prepareDataContent (_context, _pathResolved, _pathInArchive, _name, _fileId, _dataContentRead, ""); _error != nil {
			return _error
		} else {
			_fingerprintContent = _fingerprintContent_0
			_dataContent = _dataContent_0
			_dataMeta = _dataMeta_0
		}
		if _fingerprintMeta_0, _dataMetaRaw_0, _error := prepareDataMeta (_context, _dataMeta); _error != nil {
			return _error
		} else {
			_fingerprintMeta = _fingerprintMeta_0
			_dataMetaRaw = _dataMetaRaw_0
		}
	}
	
	if _context.includeStripped {
		for _, _suffix := range StripSuffixes {
			if strings.HasSuffix (_pathInArchive, _suffix) {
				_pathInArchive := _pathInArchive [: len (_pathInArchive) - len (_suffix)]
				if _error := archiveReference (_context, NamespaceFilesContent, _pathInArchive, _fingerprintContent, _fingerprintMeta); _error != nil {
					return _error
				}
				break
			}
		}
	}
	
	if _error := archiveReference (_context, NamespaceFilesContent, _pathInArchive, _fingerprintContent, _fingerprintMeta); _error != nil {
		return _error
	}
	if _dataMetaRaw != nil {
		if _error := archiveDataMeta (_context, _fingerprintMeta, _dataMetaRaw); _error != nil {
			return _error
		}
	}
	if _dataContent != nil {
		if _error := archiveDataContent (_context, _fingerprintContent, _dataContent); _error != nil {
			return _error
		}
	}
	
	if ! _wasStored {
		_context.storedFiles[_fileId] = [2]string { _fingerprintContent, _fingerprintMeta }
	}
	
	return nil
}




func archiveFolder (_context *context, _pathResolved string, _pathInArchive string, _names []string, _stats map[string]os.FileInfo) (error) {
	
	type Entry struct {
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
		Size uint64 `json:"size,omitempty"`
	}
	
	type Folder struct {
		Entries []Entry `json:"entries,omitempty"`
		Indices []string `json:"indices,omitempty"`
	}
	
	_entries := make ([]Entry, 0, len (_names))
	if _context.includeFolderListing {
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
	
	if ! _context.includeFolderListing {
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
	
	if _, _, _error := archiveReferenceAndData (_context, NamespaceFoldersContent, _pathResolved, _pathInArchive, "", _data, MimeTypeJson); _error != nil {
		return _error
	}
	
	return nil
}




func archiveReferenceAndData (_context *context, _namespace string, _pathResolved string, _pathInArchive string, _name string, _dataContent []byte, _dataType string) (string, string, error) {
	
	var _fingerprintContent string
	var _fingerprintMeta string
	var _dataMeta map[string]string
	var _dataMetaRaw []byte
	
	_dataContentRead := func () ([]byte, error) {
			return _dataContent, nil
		}
	
	if _fingerprintContent_0, _dataContent_0, _dataMeta_0, _error := prepareDataContent (_context, _pathResolved, _pathInArchive, _name, "", _dataContentRead, _dataType); _error != nil {
		return "", "", _error
	} else {
		_fingerprintContent = _fingerprintContent_0
		_dataContent = _dataContent_0
		_dataMeta = _dataMeta_0
	}
	if _fingerprintMeta_0, _dataMetaRaw_0, _error := prepareDataMeta (_context, _dataMeta); _error != nil {
		return "", "", _error
	} else {
		_fingerprintMeta = _fingerprintMeta_0
		_dataMetaRaw = _dataMetaRaw_0
	}
	
	if _error := archiveReference (_context, _namespace, _pathInArchive, _fingerprintContent, _fingerprintMeta); _error != nil {
		return "", "", _error
	}
	if _dataMetaRaw != nil {
		if _error := archiveDataMeta (_context, _fingerprintMeta, _dataMetaRaw); _error != nil {
			return "", "", _error
		}
	}
	if _dataContent != nil {
		if _error := archiveDataContent (_context, _fingerprintContent, _dataContent); _error != nil {
			return "", "", _error
		}
	}
	
	return _fingerprintContent, _fingerprintMeta, nil
}




func archiveDataContent (_context *context, _fingerprintContent string, _dataContent []byte) (error) {
	
	if _wasStored, _ := _context.storedDataContent[_fingerprintContent]; _wasStored {
		return fmt.Errorf ("[256cde78]  data content already stored:  `%s`!", _fingerprintContent)
	}
	
	{
		var _key string
		if _key_0, _error := prepareKey (_context, NamespaceDataContent, _fingerprintContent); _error == nil {
			_key = fmt.Sprintf ("%s:%s", NamespaceDataContent, _key_0)
		} else {
			return _error
		}
		if _context.debug {
			log.Printf ("[dd] [085d83ec]  data-content ++ `%s` %d\n", _key, len (_dataContent))
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataContent); _error != nil {
			return _error
		}
		_context.cdbWriteCount += 1
		_context.cdbWriteKeySize += len (_key)
		_context.cdbWriteDataSize += len (_dataContent)
	}
	
	_context.storedDataContent[_fingerprintContent] = true
	
	return nil
}


func archiveDataMeta (_context *context, _fingerprintMeta string, _dataMeta []byte) (error) {
	
	if _wasStored, _ := _context.storedDataMeta[_fingerprintMeta]; _wasStored {
		return fmt.Errorf ("[2918c4e2]  data meta already stored:  `%s`!", _fingerprintMeta)
	}
	
	{
		var _key string
		if _key_0, _error := prepareKey (_context, NamespaceDataMetadata, _fingerprintMeta); _error == nil {
			_key = fmt.Sprintf ("%s:%s", NamespaceDataMetadata, _key_0)
		} else {
			return _error
		}
		if _context.debug {
			log.Printf ("[dd] [07737b98]  data-meta    ++ `%s` %d\n", _key, len (_dataMeta))
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataMeta); _error != nil {
			return _error
		}
		_context.cdbWriteCount += 1
		_context.cdbWriteKeySize += len (_key)
		_context.cdbWriteDataSize += len (_dataMeta)
	}
	
	_context.storedDataMeta[_fingerprintMeta] = true
	
	return nil
}




func archiveReference (_context *context, _namespace string, _pathInArchive string, _fingerprintContent string, _fingerprintMeta string) (error) {
	
	switch _namespace {
		case NamespaceFilesContent :
			_context.storedFilePaths = append (_context.storedFilePaths, _pathInArchive)
		case NamespaceFoldersContent :
			_context.storedFolderPaths = append (_context.storedFolderPaths, _pathInArchive)
		default :
			return fmt.Errorf ("[051a102a]")
	}
	_context.archivedReferences += 1
	
	_key := fmt.Sprintf ("%s:%s", _namespace, _pathInArchive)
	
	var _keyMeta, _keyContent string
	if _key_0, _error := prepareKey (_context, NamespaceDataMetadata, _fingerprintMeta); _error == nil {
		_keyMeta = _key_0
	} else {
		return _error
	}
	if _key_0, _error := prepareKey (_context, NamespaceDataContent, _fingerprintContent); _error == nil {
		_keyContent = _key_0
	} else {
		return _error
	}
	_references := fmt.Sprintf ("%s:%s", _keyMeta, _keyContent)
	
	if _context.debug {
		log.Printf ("[dd] [2b9c053a]  reference    ++ `%s` :: `%s` -> `%s` ~ `%s`\n", _namespace, _pathInArchive, _keyMeta, _keyContent)
	}
	
	if _error := _context.cdbWriter.Put ([]byte (_key), []byte (_references)); _error != nil {
		return _error
	}
	_context.cdbWriteCount += 1
	_context.cdbWriteKeySize += len (_key)
	_context.cdbWriteDataSize += len (_references)
	
	if _context.progress {
		if _context.archivedReferences <= 1 {
			_context.progressLast = time.Now ()
		}
		if ((_context.archivedReferences % 1000) == 0) || (((_context.archivedReferences % 10) == 0) && (time.Since (_context.progressLast) .Seconds () >= 6)) {
			log.Printf ("[ii] [5193276e]  pogress      -- %0.2f minutes -- %d files, %d folders, %0.2f MiB (%0.2f MiB/s) -- %d compressed (%0.2f MiB, %0.2f%%) -- %d records (%0.2f MiB)\n",
					time.Since (_context.progressStarted) .Minutes (),
					len (_context.storedFilePaths),
					len (_context.storedFolderPaths),
					float32 (_context.dataUncompressedSize) / 1024 / 1024,
					float64 (_context.dataUncompressedSize) / 1024 / 1024 / (time.Since (_context.progressStarted) .Seconds () + 0.001),
					_context.dataCompressedCount,
					float32 (_context.dataUncompressedSize - _context.dataCompressedSize) / 1024 / 1024,
					(float32 (_context.dataUncompressedSize - _context.dataCompressedSize) / float32 (_context.dataUncompressedSize) * 100),
					_context.cdbWriteCount,
					float32 (_context.cdbWriteKeySize + _context.cdbWriteDataSize) / 1024 / 1024,
				)
			_context.progressLast = time.Now ()
		}
	}
	
	return nil
}




func prepareDataContent (_context *context, _pathResolved string, _pathInArchive string, _name string, _dataContentId string, _dataContentRead func () ([]byte, error), _dataType string) (string, []byte, map[string]string, error) {
	
	type DataPrepared struct {
		DataFingerprint string
		DataSize int
		DataType string
	}
	
	var _dataPreparedCached bool
	var _dataPrepared *DataPrepared
	if (_context.sourcesCache != nil) && (_dataContentId != "") {
		_cacheTxn, _error := _context.sourcesCache.Begin (false)
		if _error != nil {
			AbortError (_error, "[5fe9ada0]  unexpected sources cache error!")
		}
		_cacheBucket := _cacheTxn.Bucket ([]byte ("prepare"))
		if _cacheBucket != nil {
			if _dataPreparedRaw := _cacheBucket.Get ([]byte (_dataContentId)); _dataPreparedRaw != nil {
				if _error := json.Unmarshal (_dataPreparedRaw, &_dataPrepared); _error != nil {
					AbortError (_error, "[6865d963]  unexpected sources cache error!")
				}
			}
			_dataPreparedCached = _dataPrepared != nil
		}
		if _error := _cacheTxn.Rollback (); _error != nil {
			AbortError (_error, "[5137d84a]  unexpected sources cache error!")
		}
	}
	
	var _fingerprintContent string
	var _dataContent []byte
	var _dataSize int
	if _dataPrepared != nil {
		_fingerprintContent = _dataPrepared.DataFingerprint
		_dataSize = _dataPrepared.DataSize
	} else {
		if _data_0, _error := _dataContentRead (); _error == nil {
			_dataContent = _data_0
		} else {
			return "", nil, nil, _error
		}
		_fingerprintContentRaw := blake3.Sum256 (_dataContent)
		_fingerprintContent = hex.EncodeToString (_fingerprintContentRaw[:])
		_dataSize = len (_dataContent)
	}
	
	if _wasStored, _ := _context.storedDataContent[_fingerprintContent]; _wasStored {
		_dataMeta := _context.storedDataContentMeta[_fingerprintContent]
		return _fingerprintContent, nil, _dataMeta, nil
	}
	
	if (_dataType == "") && (_dataPrepared != nil) {
		_dataType = _dataPrepared.DataType
	}
	if (_dataType == "") && (_name != "") {
		_extension := filepath.Ext (_name)
		if _extension != "" {
			_extension = _extension[1:]
		}
		_dataType, _ = MimeTypesByExtension[_extension]
	}
	if (_dataType == "") {
		_dataType = http.DetectContentType (_dataContent)
	}
	if _dataType == "" {
		_dataType = MimeTypeRaw
	}
	
	if (_context.sourcesCache != nil) && (_dataContentId != "") && !_dataPreparedCached {
		_dataPrepared = & DataPrepared {
				DataFingerprint : _fingerprintContent,
				DataSize : _dataSize,
				DataType : _dataType,
			}
		_cacheTxn, _error := _context.sourcesCache.Begin (true)
		if _error != nil {
			AbortError (_error, "[acf09d20]  unexpected sources cache error!")
		}
		_cacheBucket := _cacheTxn.Bucket ([]byte ("prepare"))
		if _cacheBucket == nil {
			if _bucket_0, _error := _cacheTxn.CreateBucket ([]byte ("prepare")); _error == nil {
				_cacheBucket = _bucket_0
			} else {
				AbortError (_error, "[c21b0972]  unexpected sources cache error!")
			}
		}
		_cacheBucket.FillPercent = 0.9
		var _dataPreparedRaw []byte
		if _data_0, _error := json.Marshal (_dataPrepared); _error == nil {
			_dataPreparedRaw = _data_0
		} else {
			AbortError (_error, "[5538658b]  unexpected sources cache error!")
		}
		if _error := _cacheBucket.Put ([]byte (_dataContentId), _dataPreparedRaw); _error != nil {
			AbortError (_error, "[b4a6b6f9]  unexpected sources cache error!")
		}
		if _error := _cacheTxn.Commit (); _error != nil {
			AbortError (_error, "[5581f8ec]  unexpected sources cache error!")
		}
	}
	
	_dataEncoding := "identity"
	_dataUncompressedSize := _dataSize
	
	var _compressAlgorithm string
	var _compressEncoding string
	if _algorithm_0, _encoding_0, _error := CompressEncoding (_context.compress); _error == nil {
		_compressAlgorithm = _algorithm_0
		_compressEncoding = _encoding_0
	} else {
		return "", nil, nil, _error
	}
	
	if _compressAlgorithm != "identity" {
		
		var _dataCompressed []byte
		var _dataCompressedCached bool
		
		if _context.compressCache != nil {
			_cacheTxn, _error := _context.compressCache.Begin (false)
			if _error != nil {
				AbortError (_error, "[91a5b78a]  unexpected compression cache error!")
			}
			_cacheBucket := _cacheTxn.Bucket ([]byte (_compressAlgorithm))
			if _cacheBucket != nil {
				_dataCompressed = _cacheBucket.Get ([]byte (_fingerprintContent))
				_dataCompressedCached = _dataCompressed != nil
			}
			if _error := _cacheTxn.Rollback (); _error != nil {
				AbortError (_error, "[a06cfe46]  unexpected compression cache error!")
			}
		}
		
		if _dataCompressed == nil {
			if _dataContent == nil {
				if _data_0, _error := _dataContentRead (); _error == nil {
					_dataContent = _data_0
				} else {
					return "", nil, nil, _error
				}
			}
			if _data_0, _error := Compress (_dataContent, _compressAlgorithm); _error == nil {
				_dataCompressed = _data_0
			} else {
				return "", nil, nil, _error
			}
		}
		
		_dataCompressedSize := len (_dataCompressed)
		_dataCompressedDelta := _dataUncompressedSize - _dataCompressedSize
		_dataCompressedRatio := float32 (_dataCompressedDelta) * 100 / float32 (_dataUncompressedSize)
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
			_dataContent = _dataCompressed
			_dataEncoding = _compressEncoding
			_dataSize = _dataCompressedSize
		}
		
		if (_context.compressCache != nil) && !_dataCompressedCached && _accepted {
			_cacheTxn, _error := _context.compressCache.Begin (true)
			if _error != nil {
				AbortError (_error, "[ddbe6a70]  unexpected compression cache error!")
			}
			_cacheBucket := _cacheTxn.Bucket ([]byte (_compressAlgorithm))
			if _cacheBucket == nil {
				if _bucket_0, _error := _cacheTxn.CreateBucket ([]byte (_compressAlgorithm)); _error == nil {
					_cacheBucket = _bucket_0
				} else {
					AbortError (_error, "[b7766792]  unexpected compression cache error!")
				}
			}
			_cacheBucket.FillPercent = 0.9
			if _error := _cacheBucket.Put ([]byte (_fingerprintContent), _dataCompressed); _error != nil {
				AbortError (_error, "[51d57220]  unexpected compression cache error!")
			}
			if _error := _cacheTxn.Commit (); _error != nil {
				AbortError (_error, "[a47c7c10]  unexpected compression cache error!")
			}
		}
		
		if _dataSize < _dataUncompressedSize {
			if _context.debug {
				log.Printf ("[dd] [271e48d6]  compress     -- %.1f%% %d (%d) `%s`\n", _dataCompressedRatio, _dataCompressedSize, _dataCompressedDelta, _pathInArchive)
			}
		} else {
			if _context.debug || _context.progress {
				log.Printf ("[dd] [2174c2d6]  compress-NOK -- %.1f%% %d (%d) `%s`\n", _dataCompressedRatio, _dataCompressedSize, _dataCompressedDelta, _pathInArchive)
			}
		}
		
	} else {
		if _context.debug && (_context.compress != "identity") {
			log.Printf ("[dd] [a9d7a281]  compress-NOK -- %d `%s`\n", _dataUncompressedSize, _pathInArchive)
		}
	}
	
	if _dataContent == nil {
		if _data_0, _error := _dataContentRead (); _error == nil {
			_dataContent = _data_0
		} else {
			return "", nil, nil, _error
		}
	}
	
	_context.dataUncompressedSize += _dataUncompressedSize
	_context.dataUncompressedCount += 1
	_context.dataCompressedSize += _dataSize
	if _dataSize != _dataUncompressedSize {
		_context.dataCompressedCount += 1
	}
	
	_dataMeta := make (map[string]string, 16)
	
	// _dataMeta["Content-Length"] = fmt.Sprintf ("%d", _dataSize)
	_dataMeta["Content-Type"] = _dataType
	_dataMeta["Content-Encoding"] = _dataEncoding
	
	if _context.includeCache {
		_dataMeta["Cache-Control"] = "public, immutable, max-age=3600"
	}
	if _context.includeEtag {
		_dataMeta["ETag"] = _fingerprintContent
	}
	
	_context.storedDataContentMeta[_fingerprintContent] = _dataMeta
	
	return _fingerprintContent, _dataContent, _dataMeta, nil
}


func prepareDataMeta (_context *context, _dataMeta map[string]string) (string, []byte, error) {
	
	var _dataMetaRaw []byte
	if _dataMetaRaw_0, _error := MetadataEncode (_dataMeta); _error == nil {
		_dataMetaRaw = _dataMetaRaw_0
	} else {
		return "", nil, _error
	}
	
	_fingerprintMetaRaw := blake3.Sum256 (_dataMetaRaw)
	_fingerprintMeta := hex.EncodeToString (_fingerprintMetaRaw[:])
	
	if _wasStored, _ := _context.storedDataMeta[_fingerprintMeta]; _wasStored {
		return _fingerprintMeta, nil, nil
	}
	
	return _fingerprintMeta, _dataMetaRaw, nil
}




func prepareKey (_context *context, _namespace string, _fingerprint string) (string, error) {
	_qualified := fmt.Sprintf ("%s:%s", _namespace, _fingerprint)
	if _key, _found := _context.storedKeys[_qualified]; _found {
		return _key, nil
	}
	_keyIndex := len (_context.storedKeys) + 1
	if _keyIndex >= (1 << 32) {
		return "", fmt.Errorf ("[aba09b4d]  maximum stored keys reached!")
	}
	_key := fmt.Sprintf ("%x", _keyIndex)
	_context.storedKeys[_qualified] = _key
	return _key, nil
}




func walkPath (_context *context, _pathResolved string, _pathInArchive string, _name string, _recursed map[string]uint, _recurse bool) (os.FileInfo, error) {
	
	if _recursed == nil {
		_recursed = make (map[string]uint, 128)
	}
	
	_pathOriginal := _pathResolved
	
	var _stat os.FileInfo
	if _stat_0, _error := os.Lstat (_pathResolved); _error == nil {
		_stat = _stat_0
	} else {
		return nil, _error
	}
	_statMode := _stat.Mode ()
	
	_isSymlink := false
	if (_stat.Mode () & os.ModeSymlink) != 0 {
		_isSymlink = true
		if _stat_0, _error := os.Stat (_pathResolved); _error == nil {
			_stat = _stat_0
		} else {
			return nil, _error
		}
	}
	_statMode = _stat.Mode ()
	
	if ! _recurse {
		return _stat, nil
	}
	
	if _isSymlink {
		if _pathResolved_0, _error := filepath.EvalSymlinks (_pathResolved); _error == nil {
			_pathResolved = _pathResolved_0
		} else {
			return nil, _error
		}
	}
	
	if _isSymlink && _context.debug {
		log.Printf ("[dd] [7eed523e]  symlink      -- `%s` -> `%s`\n", _pathInArchive, _pathResolved)
	}
	
	
	if _statMode.IsRegular () {
		
		if _context.debug {
			log.Printf ("[dd] [da429eaa]  file         -- `%s`\n", _pathInArchive)
		}
		if _error := archiveFile (_context, _pathResolved, _pathInArchive, _name); _error != nil {
			return nil, _error
		}
		return _stat, nil
		
	} else if _statMode.IsDir () {
		
		_wasRecursed, _ := _recursed[_pathResolved]
		if _wasRecursed > 0 {
			log.Printf ("[ww] [2e1744c9]  detected directory loop for `%s` resolving to `%s`;  ignoring!\n", _pathOriginal, _pathResolved)
			return _stat, nil
		}
		_recursed[_pathResolved] = _wasRecursed + 1
		
		_childsName := make ([]string, 0, 16)
		_childsPathResolved := make (map[string]string, 16)
		_childsPathInArchive := make (map[string]string, 16)
		_childsStat := make (map[string]os.FileInfo, 16)
		
		var _wildcardName string
		
		if _context.debug {
			log.Printf ("[dd] [2d22d910]  folder       |> `%s`\n", _pathInArchive)
		}
		if _stream, _error := os.Open (_pathResolved); _error == nil {
			defer _stream.Close ()
			_loop : for {
				switch _buffer, _error := _stream.Readdir (1024); _error {
					case nil :
						if _context.debug && (len (_buffer) > 0) {
							log.Printf ("[dd] [d4c30c66]  folder       |~ `%s` %d\n", _pathInArchive, len (_buffer))
						}
						for _, _childStat := range _buffer {
							
							_childName := _childStat.Name ()
							_childPathResolved := filepath.Join (_pathResolved, _childName)
							_childPathInArchive := filepath.Join (_pathInArchive, _childName)
							
							if _childStat_0, _error := walkPath (_context, _childPathResolved, _childPathInArchive, _childName, _recursed, false); _error == nil {
								_childStat = _childStat_0
							} else {
								return nil, _error
							}
							
							_childsPathResolved[_childName] = _childPathResolved
							_childsPathInArchive[_childName] = _childPathInArchive
							_childsStat[_childName] = _childStat
							
							if strings.HasPrefix (_childName, "_wildcard.") {
								_wildcardName = _childName
								continue
							}
							if ShouldSkipName (_childName) {
								if _context.debug {
									log.Printf ("[dd] [f11b5ba1]  skip         !! `%s`\n", _childPathInArchive)
								}
								continue
							}
							
							_childsName = append (_childsName, _childName)
						}
					case io.EOF :
						break _loop
					default :
						return nil, _error
				}
			}
		}
		if _context.debug {
			log.Printf ("[dd] [2d61b965]  folder       |< `%s`\n", _pathInArchive)
		}
		
		sort.Strings (_childsName)
		
		if _context.debug {
			log.Printf ("[dd] [a4475a48]  folder       -- `%s`\n", _pathInArchive)
		}
		if _error := archiveFolder (_context, _pathResolved, _pathInArchive, _childsName, _childsStat); _error != nil {
			return nil, _error
		}
		
		if _wildcardName != "" {
			_childPathInArchive := filepath.Join (_pathInArchive, "*")
			if _, _error := walkPath (_context, _childsPathResolved[_wildcardName], _childPathInArchive, _wildcardName, _recursed, true); _error != nil {
				return nil, _error
			}
		}
		
		if _context.debug {
			log.Printf ("[dd] [ce1fe181]  folder       |> `%s`\n", _pathInArchive)
		}
		for _, _childName := range _childsName {
			if _childsStat[_childName] .Mode () .IsRegular () {
				if _, _error := walkPath (_context, _childsPathResolved[_childName], _childsPathInArchive[_childName], _childName, _recursed, true); _error != nil {
					return nil, _error
				}
			}
		}
		for _, _childName := range _childsName {
			if _childsStat[_childName] .Mode () .IsDir () {
				if _, _error := walkPath (_context, _childsPathResolved[_childName], _childsPathInArchive[_childName], _childName, _recursed, true); _error != nil {
					return nil, _error
				}
			}
		}
		for _, _childName := range _childsName {
			if (! _childsStat[_childName] .Mode () .IsRegular ()) && (! _childsStat[_childName] .Mode () .IsDir ()) {
				if _, _error := walkPath (_context, _childsPathResolved[_childName], _childsPathInArchive[_childName], _childName, _recursed, true); _error != nil {
					return nil, _error
				}
			}
		}
		if _context.debug {
			log.Printf ("[dd] [6b7de03b]  folder       |< `%s`\n", _pathInArchive)
		}
		
		_recursed[_pathResolved] = _wasRecursed
		return _stat, nil
		
	} else {
		return nil, fmt.Errorf ("[d9b836d7]  unexpected file type for `%s`:  `%s`!", _pathResolved, _statMode)
	}
}




func Main () () {
	Main_0 (main_0)
}


func main_0 () (error) {
	
	
	var _sourcesFolder string
	var _sourcesCache string
	var _archiveFile string
	var _compress string
	var _compressCache string
	var _includeIndex bool
	var _includeStripped bool
	var _includeCache bool
	var _includeEtag bool
	var _includeFileListing bool
	var _includeFolderListing bool
	var _progress bool
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("kawipiko-archiver", flag.ContinueOnError)
		
		_flags.Usage = func () () {
			fmt.Fprintf (os.Stderr, "%s",
`
  ====  kawipiko -- blazingly fast static HTTP server  ====

  |  Documentation, issues and sources:
  |      * https://github.com/volution/kawipiko
  |  Authors:
  |      * Ciprian Dorin Craciun
  |          ciprian@volution.ro
  |          ciprian.craciun@gmail.com
  |          https://volution.ro/ciprian
  -----------------------------------------------------------

  kawipiko-archiver

    --sources <path>

    --archive <path>

    --compress <gzip | zopfli | brotli | identity>
    --compress-cache <path>

    --exclude-index
    --exclude-strip
    --exclude-cache
    --include-etag

    --exclude-file-listing
    --include-folder-listing

    --progress
    --debug

  ** for details see:
     https://github.com/volution/kawipiko#kawipiko-archiver

`)
		}
		
		_sourcesFolder_0 := _flags.String ("sources", "", "")
		_sourcesCache_0 := _flags.String ("sources-cache", "", "")
		_archiveFile_0 := _flags.String ("archive", "", "")
		_compress_0 := _flags.String ("compress", "", "")
		_compressCache_0 := _flags.String ("compress-cache", "", "")
		_excludeIndex_0 := _flags.Bool ("exclude-index", false, "")
		_excludeStripped_0 := _flags.Bool ("exclude-strip", false, "")
		_excludeCache_0 := _flags.Bool ("exclude-cache", false, "")
		_includeEtag_0 := _flags.Bool ("include-etag", false, "")
		_excludeFileListing_0 := _flags.Bool ("exclude-file-listing", false, "")
		_includeFolderListing_0 := _flags.Bool ("include-folder-listing", false, "")
		_progress_0 := _flags.Bool ("progress", false, "")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_sourcesFolder = *_sourcesFolder_0
		_sourcesCache = *_sourcesCache_0
		_archiveFile = *_archiveFile_0
		_compress = *_compress_0
		_compressCache = *_compressCache_0
		_includeIndex = ! *_excludeIndex_0
		_includeStripped = ! *_excludeStripped_0
		_includeCache = ! *_excludeCache_0
		_includeEtag = *_includeEtag_0
		_includeFileListing = ! *_excludeFileListing_0
		_includeFolderListing = *_includeFolderListing_0
		_progress = *_progress_0
		_debug = *_debug_0
		
		if _sourcesFolder == "" {
			AbortError (nil, "[515ee462]  expected sources folder argument!")
		}
		if _archiveFile == "" {
			AbortError (nil, "[5e8da985]  expected archive file argument!")
		}
	}
	
	
	var _cdbWriter *cdb.Writer
	if _db_0, _error := cdb.Create (_archiveFile); _error == nil {
		_cdbWriter = _db_0
	} else {
		AbortError (_error, "[85234ba0]  failed creating archive (while opening)!")
	}
	
	var _compressCacheDb *bbolt.DB
	if _compressCache != "" {
		_options := bbolt.Options {
				PageSize : 128 * 1024,
				InitialMmapSize : 128 * 1024 * 1024,
				NoFreelistSync : true,
				NoSync : true,
			}
		if _db_0, _error := bbolt.Open (_compressCache, 0600, &_options); _error == nil {
			_compressCacheDb = _db_0
		} else {
			AbortError (_error, "[eaff07f6]  failed opening compression cache!")
		}
	}
	
	var _sourcesCacheDb *bbolt.DB
	if _sourcesCache != "" {
		_options := bbolt.Options {
				PageSize : 4 * 1024,
				InitialMmapSize : 4 * 1024 * 1024,
				NoFreelistSync : true,
				NoSync : true,
			}
		if _db_0, _error := bbolt.Open (_sourcesCache, 0600, &_options); _error == nil {
			_sourcesCacheDb = _db_0
		} else {
			AbortError (_error, "[17a308dc]  failed opening sources cache!")
		}
	}
	
	_context := & context {
			cdbWriter : _cdbWriter,
			storedFilePaths : make ([]string, 0, 16 * 1024),
			storedFolderPaths : make ([]string, 0, 16 * 1024),
			storedDataMeta : make (map[string]bool, 16 * 1024),
			storedDataContent : make (map[string]bool, 16 * 1024),
			storedDataContentMeta : make (map[string]map[string]string, 16 * 1024),
			storedFiles : make (map[string][2]string, 16 * 1024),
			storedKeys : make (map[string]string, 16 * 1024),
			compress : _compress,
			compressCache : _compressCacheDb,
			sourcesCache : _sourcesCacheDb,
			includeIndex : _includeIndex,
			includeStripped : _includeStripped,
			includeCache : _includeCache,
			includeEtag : _includeEtag,
			includeFileListing : _includeFileListing,
			includeFolderListing : _includeFolderListing,
			progress : _progress,
			debug : _debug,
		}
	
	_context.progressStarted = time.Now ()
	
	if _error := _context.cdbWriter.Put ([]byte (NamespaceSchemaVersion), []byte (CurrentSchemaVersion)); _error != nil {
		AbortError (_error, "[43228812]  failed writing archive!")
	}
	_context.cdbWriteCount += 1
	_context.cdbWriteKeySize += len (NamespaceSchemaVersion)
	_context.cdbWriteDataSize += len (CurrentSchemaVersion)
	
	if _, _error := walkPath (_context, _sourcesFolder, "/", filepath.Base (_sourcesFolder), nil, true); _error != nil {
		AbortError (_error, "[b6a19ef4]  failed walking folder!")
	}
	
	if _includeFileListing {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedFilePaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _error := _context.cdbWriter.Put ([]byte (NamespaceFilesIndex), _buffer); _error != nil {
			AbortError (_error, "[1dbdde05]  failed writing archive!")
		}
		_context.cdbWriteCount += 1
		_context.cdbWriteKeySize += len (NamespaceFilesIndex)
		_context.cdbWriteDataSize += len (_buffer)
	}
	
	if _includeFolderListing {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedFolderPaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _error := _context.cdbWriter.Put ([]byte (NamespaceFoldersIndex), _buffer); _error != nil {
			AbortError (_error, "[e2dd2de0]  failed writing archive!")
		}
		_context.cdbWriteCount += 1
		_context.cdbWriteKeySize += len (NamespaceFilesIndex)
		_context.cdbWriteDataSize += len (_buffer)
	}
	
	if _error := _context.cdbWriter.Close (); _error != nil {
		AbortError (_error, "[bbfb8478]  failed creating archive (while closing)!")
	}
	_context.cdbWriter = nil
	
	if _context.compressCache != nil {
		if _error := _context.compressCache.Close (); _error != nil {
			AbortError (_error, "[53cbe28d]  failed closing compression cache!")
		}
	}
	
	if _context.sourcesCache != nil {
		if _error := _context.sourcesCache.Close (); _error != nil {
			AbortError (_error, "[7fe3692c]  failed closing compression cache!")
		}
	}
	
	if true {
		log.Printf ("[ii] [56f63575]  completed    -- %0.2f minutes -- %d files, %d folders, %0.2f MiB (%0.2f MiB/s) -- %d compressed (%0.2f MiB, %0.2f%%) -- %d records (%0.2f MiB)\n",
				time.Since (_context.progressStarted) .Minutes (),
				len (_context.storedFilePaths),
				len (_context.storedFolderPaths),
				float32 (_context.dataUncompressedSize) / 1024 / 1024,
				float64 (_context.dataUncompressedSize) / 1024 / 1024 / (time.Since (_context.progressStarted) .Seconds () + 0.001),
				_context.dataCompressedCount,
				float32 (_context.dataUncompressedSize - _context.dataCompressedSize) / 1024 / 1024,
				(float32 (_context.dataUncompressedSize - _context.dataCompressedSize) / float32 (_context.dataUncompressedSize) * 100),
				_context.cdbWriteCount,
				float32 (_context.cdbWriteKeySize + _context.cdbWriteDataSize) / 1024 / 1024,
			)
	}
	
	return nil
}

