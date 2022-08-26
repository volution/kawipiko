

package archiver


import "bytes"
import "encoding/gob"
import "encoding/hex"
import "encoding/json"
import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "log"
import "net/http"
import "net/url"
import "os"
import "path/filepath"
import "runtime"
import "runtime/debug"
import "sort"
import "strings"
import "time"

import "github.com/colinmarc/cdb"
import "github.com/zeebo/blake3"
import "go.etcd.io/bbolt"

import . "github.com/volution/kawipiko/lib/common"
import . "github.com/volution/kawipiko/lib/archiver"

import _ "embed"




type context struct {
	cdbWriter *cdb.Writer
	cdbWriteCount int
	cdbWriteKeySize int
	cdbWriteDataSize int
	cdbWriteKeys map[string]bool
	storedFilePaths []string
	storedFolderPaths []string
	storedRedirectPaths []string
	storedDataMeta map[string]bool
	storedDataContent map[string]bool
	storedDataContentMeta map[string]map[string]string
	storedFiles map[string][2]string
	storedKeys map[string]uint64
	archivedReferences uint
	compress string
	compressLevel int
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
	includeFolderListing bool
	includePathsIndex bool
	progress bool
	progressStarted time.Time
	progressLast time.Time
	debug bool
}




func archiveFile (_context *context, _pathResolved string, _pathInArchive string, _name string, _statusPerhaps uint) (error) {
	
	var _fileDev uint64
	var _fileInode uint64
	var _fileSize uint64
	var _fileTimestamp uint64
	if _stat, _error := os.Stat (_pathResolved); _error == nil {
		_fileDev, _fileInode, _fileSize, _fileTimestamp, _error = SysStatExtract (_stat)
		if _error != nil {
			return _error
		}
	} else {
		return _error
	}
	
	_fileIdText := fmt.Sprintf ("%d.%d.%d.%d", _fileDev, _fileInode, _fileSize, _fileTimestamp)
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
					_fileDev_0, _fileInode_0, _fileSize_0, _fileTimestamp_0, _error := SysStatExtract (_stat)
					if _error != nil {
						return nil, _error
					}
					if (
							(_fileDev != _fileDev_0) ||
							(_fileInode != _fileInode_0) ||
							(_fileSize != _fileSize_0) ||
							(_fileTimestamp != _fileTimestamp_0)) {
						return nil, fmt.Errorf ("[3a07643b]  file changed while reading:  `%s`!", _pathResolved)
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
					_, _, _fileSize_0, _fileTimestamp_0, _error := SysStatExtract (_stat)
					if _error != nil {
						return nil, _error
					}
					if (
							(_fileSize != _fileSize_0) ||
							(_fileTimestamp != _fileTimestamp_0)) {
						return nil, fmt.Errorf ("[9689146e]  file changed while reading:  `%s`!", _pathResolved)
					}
				} else {
					return nil, _error
				}
				return _data, nil
			}
		
		if _fingerprintContent_0, _dataContent_0, _dataMeta_0, _error := prepareDataContent (_context, _pathResolved, _pathInArchive, _name, _fileId, _dataContentRead, "", _statusPerhaps); _error != nil {
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




func archiveFolder (_context *context, _pathResolved string, _pathInArchive string, _names []string, _stats map[string]os.FileInfo, _statusPerhaps uint) (error) {
	
	if ! strings.HasSuffix (_pathInArchive, "/") {
		_pathInArchive = _pathInArchive + "/"
	}
	
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
			_indexPathInArchive := _pathInArchive
			archiveFile (_context, _indexPathResolved, _indexPathInArchive, _indexNameFirst, _statusPerhaps)
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
	
	if _, _, _error := archiveReferenceAndData (_context, NamespaceFoldersContent, _pathResolved, _pathInArchive, "", _data, MimeTypeJson, _statusPerhaps); _error != nil {
		return _error
	}
	
	return nil
}




func archiveReferenceAndData (_context *context, _namespace string, _pathResolved string, _pathInArchive string, _name string, _dataContent []byte, _dataType string, _statusPerhaps uint) (string, string, error) {
	
	var _fingerprintContent string
	var _fingerprintMeta string
	var _dataMeta map[string]string
	var _dataMetaRaw []byte
	
	_dataContentRead := func () ([]byte, error) {
			return _dataContent, nil
		}
	
	if _fingerprintContent_0, _dataContent_0, _dataMeta_0, _error := prepareDataContent (_context, _pathResolved, _pathInArchive, _name, "", _dataContentRead, _dataType, _statusPerhaps); _error != nil {
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


func archiveReferenceAndDataWithMeta (_context *context, _namespace string, _pathInArchive string, _dataContent []byte, _dataMeta map[string]string) (string, string, error) {
	
	var _fingerprintContent string
	var _fingerprintMeta string
	var _dataMetaRaw []byte
	
	_fingerprintContentRaw := blake3.Sum256 (_dataContent)
	_fingerprintContent = hex.EncodeToString (_fingerprintContentRaw[:])
	
	if _wasStored, _ := _context.storedDataContent[_fingerprintContent]; _wasStored {
		_dataContent = nil
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
		if _key_0, _error := prepareKeyString (_context, NamespaceDataContent, _fingerprintContent); _error == nil {
			_key = _key_0
		} else {
			return _error
		}
		if _context.debug {
			log.Printf ("[dd] [085d83ec]  data-content ++ `%s` %d\n", _key, len (_dataContent))
		}
		if _found, _ := _context.cdbWriteKeys[_key]; _found {
			return fmt.Errorf ("[53aeea4b]  duplicate key encountered:  `%s`", _key)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataContent); _error != nil {
			return _error
		}
		_context.cdbWriteKeys[_key] = true
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
		if _key_0, _error := prepareKeyString (_context, NamespaceDataMetadata, _fingerprintMeta); _error == nil {
			_key = _key_0
		} else {
			return _error
		}
		if _context.debug {
			log.Printf ("[dd] [07737b98]  data-meta    ++ `%s` %d\n", _key, len (_dataMeta))
		}
		if _found, _ := _context.cdbWriteKeys[_key]; _found {
			return fmt.Errorf ("[8f2c6911]  duplicate key encountered:  `%s`", _key)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataMeta); _error != nil {
			return _error
		}
		_context.cdbWriteKeys[_key] = true
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
		case NamespaceRedirectsContent :
			_context.storedRedirectPaths = append (_context.storedRedirectPaths, _pathInArchive)
		default :
			return fmt.Errorf ("[051a102a]")
	}
	_context.archivedReferences += 1
	
	_namespacePrefix := KeyNamespacePrefix (_namespace)
	_key := fmt.Sprintf ("%c:%s", _namespacePrefix, _pathInArchive)
	
	var _keyMeta, _keyContent uint64
	if _key_0, _error := prepareKeyUint (_context, NamespaceDataMetadata, _fingerprintMeta); _error == nil {
		_keyMeta = _key_0
	} else {
		return _error
	}
	if _key_0, _error := prepareKeyUint (_context, NamespaceDataContent, _fingerprintContent); _error == nil {
		_keyContent = _key_0
	} else {
		return _error
	}
	var _references string
	if _references_0, _error := EncodeKeysPairToString (NamespaceDataMetadata, _keyMeta, NamespaceDataContent, _keyContent); _error == nil {
		_references = _references_0
	} else {
		return _error
	}
	
	if _context.debug {
		log.Printf ("[dd] [2b9c053a]  reference    ++ `%s` :: `%s` -> `%s`\n", _namespace, _pathInArchive, _references)
	}
	
	if _found, _ := _context.cdbWriteKeys[_key]; _found {
		return fmt.Errorf ("[3d856291]  duplicate key encountered:  `%s`", _key)
	}
	if _error := _context.cdbWriter.Put ([]byte (_key), []byte (_references)); _error != nil {
		return _error
	}
	_context.cdbWriteKeys[_key] = true
	_context.cdbWriteCount += 1
	_context.cdbWriteKeySize += len (_key)
	_context.cdbWriteDataSize += len (_references)
	
	if _context.progress {
		if _context.archivedReferences <= 1 {
			_context.progressLast = time.Now ()
		}
		if (
				(((_context.archivedReferences % 1000) == 0) && (time.Since (_context.progressLast) .Seconds () >= 2)) ||
				(((_context.archivedReferences % 10) == 0) && (time.Since (_context.progressLast) .Seconds () >= 6))) {
			log.Printf ("[ii] [5193276e]  pogress      -- %0.2f min -- %d fil, %d fol, %d red, %0.2f M (%0.2f M/s) -- %d comp, %0.2f%% -- %d rec, %0.2f M\n",
					time.Since (_context.progressStarted) .Minutes (),
					len (_context.storedFilePaths),
					len (_context.storedFolderPaths),
					len (_context.storedRedirectPaths),
					float32 (_context.dataUncompressedSize) / 1024 / 1024,
					float64 (_context.dataUncompressedSize) / 1024 / 1024 / (time.Since (_context.progressStarted) .Seconds () + 0.001),
					_context.dataCompressedCount,
					(float32 (_context.dataUncompressedSize - _context.dataCompressedSize) / float32 (_context.dataUncompressedSize) * 100),
					_context.cdbWriteCount,
					float32 (_context.cdbWriteKeySize + _context.cdbWriteDataSize) / 1024 / 1024,
				)
			_context.progressLast = time.Now ()
		}
	}
	
	return nil
}




func prepareDataContent (_context *context, _pathResolved string, _pathInArchive string, _name string, _dataContentId string, _dataContentRead func () ([]byte, error), _dataType string, _statusPerhaps uint) (string, []byte, map[string]string, error) {
	
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
				if _error := gobUnmarshal (_dataPreparedRaw, &_dataPrepared); _error != nil {
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
	
	_dataContent_1 := _dataContent
	_dataSize_1 := _dataSize
	_dataContentRead_1 := _dataContentRead
	_dataContentRead = func () ([]byte, error) {
			if _dataContent_1 != nil {
				return _dataContent_1, nil
			}
			if (_context.sourcesCache != nil) && (_dataContentId != "") && (_dataSize_1 <= 16 * 1024) {
				_cacheTxn, _error := _context.sourcesCache.Begin (false)
				if _error != nil {
					AbortError (_error, "[4e7853f4]  unexpected sources cache error!")
				}
				_cacheBucket := _cacheTxn.Bucket ([]byte ("content"))
				if _cacheBucket != nil {
					if _dataContent_0 := _cacheBucket.Get ([]byte (_fingerprintContent)); _dataContent_0 != nil {
						_dataContent_1 = _dataContent_0
					}
				}
				if _error := _cacheTxn.Rollback (); _error != nil {
					AbortError (_error, "[a0dd23c2]  unexpected sources cache error!")
				}
			}
			if _dataContent_1 != nil {
				return _dataContent_1, nil
			}
			if _dataContent_0, _error := _dataContentRead_1 (); _error == nil {
				_dataContent_1 = _dataContent_0
			} else {
				return nil, _error
			}
			if (_context.sourcesCache != nil) && (_dataContentId != "") && (_dataSize_1 <= 16 * 1024) {
				_cacheTxn, _error := _context.sourcesCache.Begin (true)
				if _error != nil {
					AbortError (_error, "[deecec9d]  unexpected sources cache error!")
				}
				_cacheBucket := _cacheTxn.Bucket ([]byte ("content"))
				if _cacheBucket == nil {
					if _bucket_0, _error := _cacheTxn.CreateBucket ([]byte ("content")); _error == nil {
						_cacheBucket = _bucket_0
					} else {
						AbortError (_error, "[40236265]  unexpected sources cache error!")
					}
				}
				_cacheBucket.FillPercent = 0.9
				if _error := _cacheBucket.Put ([]byte (_fingerprintContent), _dataContent_1); _error != nil {
					AbortError (_error, "[84d20b6d]  unexpected sources cache error!")
				}
				if _error := _cacheTxn.Commit (); _error != nil {
					AbortError (_error, "[5468cced]  unexpected sources cache error!")
				}
			}
			return _dataContent_1, nil
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
		if _dataType != "" {
			if _dataType_0, _ := MimeTypesAliases[_dataType]; _dataType_0 != "" {
				_dataType = _dataType_0
			}
		}
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
		if _data_0, _error := gobMarshal (_dataPrepared); _error == nil {
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
		var _dataCompressedZero bool
		var _dataCompressedCached bool
		
		_cacheBucketName := fmt.Sprintf ("%s:%d", _compressAlgorithm, _context.compressLevel)
		
		if _context.compressCache != nil {
			_cacheTxn, _error := _context.compressCache.Begin (false)
			if _error != nil {
				AbortError (_error, "[91a5b78a]  unexpected compression cache error!")
			}
			_cacheBucket := _cacheTxn.Bucket ([]byte (_cacheBucketName))
			if _cacheBucket != nil {
				_dataCompressed = _cacheBucket.Get ([]byte (_fingerprintContent))
				_dataCompressedCached = _dataCompressed != nil
				if _dataCompressedCached {
					_dataCompressedZero = len (_dataCompressed) == 0
					if _dataCompressedZero {
						_dataCompressed = nil
					}
				}
			}
			if _error := _cacheTxn.Rollback (); _error != nil {
				AbortError (_error, "[a06cfe46]  unexpected compression cache error!")
			}
		}
		
		if _dataCompressed == nil && !_dataCompressedZero {
			if _dataContent == nil {
				if _data_0, _error := _dataContentRead (); _error == nil {
					_dataContent = _data_0
				} else {
					return "", nil, nil, _error
				}
			}
			if _data_0, _error := Compress (_dataContent, _compressAlgorithm, _context.compressLevel); _error == nil {
				_dataCompressed = _data_0
			} else {
				return "", nil, nil, _error
			}
		}
		
		_dataCompressedSize := len (_dataCompressed)
		_dataCompressedDelta := _dataUncompressedSize - _dataCompressedSize
		_dataCompressedRatio := float32 (_dataCompressedDelta) * 100 / float32 (_dataUncompressedSize)
		_accepted := false
		if !_dataCompressedZero {
			_accepted = _accepted || ((_dataUncompressedSize > (1024 * 1024)) && (_dataCompressedRatio >= 5))
			_accepted = _accepted || ((_dataUncompressedSize > (64 * 1024)) && (_dataCompressedRatio >= 10))
			_accepted = _accepted || ((_dataUncompressedSize > (32 * 1024)) && (_dataCompressedRatio >= 15))
			_accepted = _accepted || ((_dataUncompressedSize > (16 * 1024)) && (_dataCompressedRatio >= 20))
			_accepted = _accepted || ((_dataUncompressedSize > (8 * 1024)) && (_dataCompressedRatio >= 25))
			_accepted = _accepted || ((_dataUncompressedSize > (4 * 1024)) && (_dataCompressedRatio >= 30))
			_accepted = _accepted || ((_dataUncompressedSize > (2 * 1024)) && (_dataCompressedRatio >= 35))
			_accepted = _accepted || ((_dataUncompressedSize > (1 * 1024)) && (_dataCompressedRatio >= 40))
			_accepted = _accepted || (_dataCompressedRatio >= 90)
		}
		if _accepted {
			_dataContent = _dataCompressed
			_dataEncoding = _compressEncoding
			_dataSize = _dataCompressedSize
		}
		
		if (_context.compressCache != nil) && !_dataCompressedCached {
			_cacheTxn, _error := _context.compressCache.Begin (true)
			if _error != nil {
				AbortError (_error, "[ddbe6a70]  unexpected compression cache error!")
			}
			_cacheBucket := _cacheTxn.Bucket ([]byte (_cacheBucketName))
			if _cacheBucket == nil {
				if _bucket_0, _error := _cacheTxn.CreateBucket ([]byte (_cacheBucketName)); _error == nil {
					_cacheBucket = _bucket_0
				} else {
					AbortError (_error, "[b7766792]  unexpected compression cache error!")
				}
			}
			_cacheBucket.FillPercent = 0.9
			_dataCompressed_1 := _dataCompressed
			if !_accepted {
				_dataCompressed_1 = []byte {}
			}
			if _error := _cacheBucket.Put ([]byte (_fingerprintContent), _dataCompressed_1); _error != nil {
				AbortError (_error, "[51d57220]  unexpected compression cache error!")
			}
			if _error := _cacheTxn.Commit (); _error != nil {
				AbortError (_error, "[a47c7c10]  unexpected compression cache error!")
			}
		}
		
		if !_dataCompressedZero {
			if _dataSize < _dataUncompressedSize {
				if _context.debug {
					log.Printf ("[dd] [271e48d6]  compress     -- %.1f%% %d (%d) `%s`\n", _dataCompressedRatio, _dataCompressedSize, _dataCompressedDelta, _pathInArchive)
				}
			} else {
				if _context.debug {
					log.Printf ("[dd] [2174c2d6]  compress-NOK -- %.1f%% %d (%d) `%s`\n", _dataCompressedRatio, _dataCompressedSize, _dataCompressedDelta, _pathInArchive)
				}
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
	
	_status := _statusPerhaps
	if _status == 0 {
		_status = 200
	}
	_dataMeta["!Status"] = fmt.Sprintf ("%d", _status)
	
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
	if _dataMetaRaw_0, _error := MetadataEncodeBinary (_dataMeta); _error == nil {
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




func prepareKeyUint (_context *context, _namespace string, _fingerprint string) (uint64, error) {
	_qualified := fmt.Sprintf ("%s:%s", _namespace, _fingerprint)
	if _keyValue, _found := _context.storedKeys[_qualified]; _found {
		return _keyValue, nil
	}
	_keyIndex := len (_context.storedKeys) + 1
	if _keyValue, _error := PrepareKey (_namespace, uint64 (_keyIndex)); _error == nil {
		_context.storedKeys[_qualified] = _keyValue
		return _keyValue, nil
	} else {
		return 0, _error
	}
}

func prepareKeyString (_context *context, _namespace string, _fingerprint string) (string, error) {
	if _keyValue, _error := prepareKeyUint (_context, _namespace, _fingerprint); _error == nil {
		if _key, _error := EncodeKeyToString (_namespace, _keyValue); _error == nil {
			return _key, nil
		} else {
			return "", _error
		}
	} else {
		return "", _error
	}
}




func walkPath (_context *context, _pathResolved string, _pathInArchive string, _name string, _recursed map[string]uint, _recurse bool, _statusPerhaps uint) (os.FileInfo, error) {
	
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
			log.Printf ("[dd] [da429eaa]  file         -- `%s` -> `%s`\n", _pathInArchive, _pathResolved)
		}
		if _error := archiveFile (_context, _pathResolved, _pathInArchive, _name, _statusPerhaps); _error != nil {
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
		var _wildcardStatus uint
		
		var _redirectsName string
		
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
							
							if _childStat_0, _error := walkPath (_context, _childPathResolved, _childPathInArchive, _childName, _recursed, false, _statusPerhaps); _error == nil {
								_childStat = _childStat_0
							} else {
								return nil, _error
							}
							
							_childsPathResolved[_childName] = _childPathResolved
							_childsPathInArchive[_childName] = _childPathInArchive
							_childsStat[_childName] = _childStat
							
							if strings.HasPrefix (_childName, "200.") || strings.HasPrefix (_childName, "_200.") || strings.HasPrefix (_childName, "_wildcard.") {
								if _wildcardName != "" {
									return nil, fmt.Errorf ("[b5afdb05]  duplicate wildcard files found:  `%s` and `%s`!", _childName, _wildcardName)
								}
								_wildcardName = _childName
								_wildcardStatus = 200
								continue
							}
							if strings.HasPrefix (_childName, "404.") || strings.HasPrefix (_childName, "_404.") {
								if _wildcardName != "" {
									return nil, fmt.Errorf ("[4c14c7bb]  duplicate wildcard files found:  `%s` and `%s`!", _childName, _wildcardName)
								}
								_wildcardName = _childName
								_wildcardStatus = 404
								continue
							}
							if (_childName == "_redirects") || (_childName == "_redirects.txt") {
								if _redirectsName != "" {
									return nil, fmt.Errorf ("[fbc0ee12]  duplicate redirects files found:  `%s` and `%s`!", _childName, _redirectsName)
								}
								_redirectsName = _childName
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
		if _error := archiveFolder (_context, _pathResolved, _pathInArchive, _childsName, _childsStat, _statusPerhaps); _error != nil {
			return nil, _error
		}
		
		if _wildcardName != "" {
			_childPathInArchive := filepath.Join (_pathInArchive, "*")
			_childPathResolved := _childsPathResolved[_wildcardName]
			if _context.debug {
				log.Printf ("[dd] [351804ff]  wildcard     -- `%s` -> `%s`\n", _childPathInArchive, _childPathResolved)
			}
			if _, _error := walkPath (_context, _childPathResolved, _childPathInArchive, _wildcardName, _recursed, true, _wildcardStatus); _error != nil {
				return nil, _error
			}
		}
		if _redirectsName != "" {
			_childPathInArchive := _pathInArchive
			_childPathResolved := _childsPathResolved[_redirectsName]
			if _context.debug {
				log.Printf ("[dd] [f49aa9f5]  redirects    -- `%s` -> `%s`\n", _childPathInArchive, _childPathResolved)
			}
			if _error := walkRedirects (_context, _childPathResolved, _childPathInArchive); _error != nil {
				return nil, _error
			}
		}
		
		if _context.debug {
			log.Printf ("[dd] [ce1fe181]  folder       |> `%s`\n", _pathInArchive)
		}
		for _, _childName := range _childsName {
			if _childsStat[_childName] .Mode () .IsRegular () {
				if _, _error := walkPath (_context, _childsPathResolved[_childName], _childsPathInArchive[_childName], _childName, _recursed, true, _statusPerhaps); _error != nil {
					return nil, _error
				}
			}
		}
		for _, _childName := range _childsName {
			if _childsStat[_childName] .Mode () .IsDir () {
				if _, _error := walkPath (_context, _childsPathResolved[_childName], _childsPathInArchive[_childName], _childName, _recursed, true, _statusPerhaps); _error != nil {
					return nil, _error
				}
			}
		}
		for _, _childName := range _childsName {
			if (! _childsStat[_childName] .Mode () .IsRegular ()) && (! _childsStat[_childName] .Mode () .IsDir ()) {
				if _, _error := walkPath (_context, _childsPathResolved[_childName], _childsPathInArchive[_childName], _childName, _recursed, true, _statusPerhaps); _error != nil {
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




func walkRedirects (_context *context, _pathResolved string, _pathInArchive string) (error) {
	
	_redirectsData, _error := os.ReadFile (_pathResolved)
	if _error != nil {
		return _error
	}
	
	_redirectLines := strings.Split (string (_redirectsData), "\n")
	
	for _, _redirectLine := range _redirectLines {
		
		_redirectLine = strings.ReplaceAll (_redirectLine, "\t", " ")
		_redirectLine = strings.Trim (_redirectLine, " ")
		for {
			_redirectLine_0 := strings.ReplaceAll (_redirectLine, "  ", " ")
			if _redirectLine == _redirectLine_0 {
				break
			}
			_redirectLine = _redirectLine_0
		}
		
		if _redirectLine == "" {
			continue
		}
		if _redirectLine[0] == '#' {
			continue
		}
		
		_redirectParts := strings.Split (_redirectLine, " ")
		if len (_redirectParts) != 3 {
			return fmt.Errorf ("[bf358c0a]  invalid redirect statement:  `%s` -- `%s`!", _pathResolved, _redirectLine)
		}
		
		_source := _redirectParts[0]
		_target := _redirectParts[1]
		_code := 0
		switch _redirectParts[2] {
			case "301" : _code = 301
			case "302" : _code = 302
			case "303" : _code = 303
			case "307" : _code = 307
			case "308" : _code = 308
			default :
				return fmt.Errorf ("[fb9ac5a5]  invalid redirect code:  `%s` -- `%s`!", _pathResolved, _redirectLine)
		}
		
		_sourceParsed := false
		if ! _sourceParsed {
			for _, _schemePrefix := range []string { "://", "http://", "https://" } {
				if strings.HasPrefix (_source, _schemePrefix) {
					if _pathInArchive != "/" {
						return fmt.Errorf ("[ceac537a]  invalid redirect absolute source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
					}
					_source_0 := "http://" + strings.TrimPrefix (_source, _schemePrefix)
					_sourceUrl, _error := url.ParseRequestURI (_source_0)
					if _error != nil {
						return fmt.Errorf ("[dd43230e]  invalid redirect source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
					}
					_sourceCanonical := _sourceUrl.String ()
					if _sourceCanonical != _source_0 {
						return fmt.Errorf ("[2c6dc950]  invalid redirect source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
					}
					_source = strings.TrimPrefix (_source, _schemePrefix)
					_source = "://" + _source
					_sourceParsed = true
					break
				}
			}
		}
		if ! _sourceParsed {
			if _source[0] == '.' {
				if _pathInArchive == "/" {
					_source = "/" + _source[2:]
				} else {
					_source = _pathInArchive + "/" + _source[2:]
				}
			} else if _source[0] == '/' {
				if _pathInArchive != "/" {
					return fmt.Errorf ("[99e025fe]  invalid redirect absolute source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
				}
			} else {
				return fmt.Errorf ("[2e716cd2]  invalid redirect source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
			}
			_sourceUrl, _error := url.ParseRequestURI (_source)
			if _error != nil {
				return fmt.Errorf ("[be9abd4d]  invalid redirect source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
			}
			_sourceCanonical := _sourceUrl.String ()
			if _sourceCanonical != _source {
				return fmt.Errorf ("[3c31deb4]  invalid redirect source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
			}
			_sourceParsed = true
		}
		if ! _sourceParsed {
			return fmt.Errorf ("[8305376f]  invalid redirect source:  `%s` -- `%s`!", _pathResolved, _redirectLine)
		}
		
		{
			_targetUrl, _error := url.Parse (_target)
			if _error != nil {
				return fmt.Errorf ("[968f7208]  invalid redirect target:  `%s` -- `%s`!", _pathResolved, _redirectLine)
			}
			_targetCanonical := _targetUrl.String ()
			if _targetCanonical != _target {
				return fmt.Errorf ("[c22eea24]  invalid redirect target:  `%s` -- `%s`!", _pathResolved, _redirectLine)
			}
		}
		
		if _context.debug {
			log.Printf ("[dd] [c8139953]  redirect     -- `%s` -> `%s` (%d)\n", _source, _target, _code)
		}
		
		_dataMeta := map[string]string {
				"!Status" : fmt.Sprintf ("%d", _code),
				"Location" : _target,
			}
		
		_fingerprintContentRaw := blake3.Sum256 ([]byte (_redirectLine))
		_fingerprintContent := hex.EncodeToString (_fingerprintContentRaw[:])
		
		if _context.includeCache {
			_dataMeta["Cache-Control"] = "public, immutable, max-age=3600"
		}
		if _context.includeEtag {
			_dataMeta["ETag"] = _fingerprintContent
		}
		
		if _, _, _error := archiveReferenceAndDataWithMeta (_context, NamespaceRedirectsContent, _source, []byte (""), _dataMeta); _error != nil {
			return _error
		}
	}
	
	return nil
}




func Main () () {
	
	if len (os.Args) == 2 {
		switch os.Args[1] {
			case "--help", "-h" :
				os.Stderr.WriteString (usageText)
				return
			case "--man" :
				os.Stderr.WriteString (manualText)
				return
		}
	}
	
	runtime.GOMAXPROCS (2)
	debug.SetGCPercent (75)
	debug.SetMaxThreads (8)
	
	Main_0 (main_0)
}


func main_0 () (error) {
	
	
	var _sourcesFolder string
	var _sourcesCache string
	var _archiveFile string
	var _compress string
	var _compressLevel int
	var _compressCache string
	var _includeIndex bool
	var _includeStripped bool
	var _includeCache bool
	var _includeEtag bool
	var _includeFolderListing bool
	var _includePathsIndex bool
	var _progress bool
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("kawipiko-archiver", flag.ContinueOnError)
		
		_flags.Usage = func () () {
			fmt.Fprintf (os.Stderr, "%s", usageText)
		}
		
		_sourcesFolder_0 := _flags.String ("sources", "", "")
		_sourcesCache_0 := _flags.String ("sources-cache", "", "")
		_archiveFile_0 := _flags.String ("archive", "", "")
		_compress_0 := _flags.String ("compress", "", "")
		_compressLevel_0 := _flags.Int ("compress-level", -1, "")
		_compressCache_0 := _flags.String ("compress-cache", "", "")
		_excludeIndex_0 := _flags.Bool ("exclude-index", false, "")
		_excludeStripped_0 := _flags.Bool ("exclude-strip", false, "")
		_excludeCache_0 := _flags.Bool ("exclude-cache", false, "")
		_includeEtag_0 := _flags.Bool ("include-etag", false, "")
		_includeFolderListing_0 := _flags.Bool ("include-folder-listing", false, "")
		_excludePathsIndex_0 := _flags.Bool ("exclude-paths-index", false, "")
		_progress_0 := _flags.Bool ("progress", false, "")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_sourcesFolder = *_sourcesFolder_0
		_sourcesCache = *_sourcesCache_0
		_archiveFile = *_archiveFile_0
		_compress = *_compress_0
		_compressLevel = *_compressLevel_0
		_compressCache = *_compressCache_0
		_includeIndex = ! *_excludeIndex_0
		_includeStripped = ! *_excludeStripped_0
		_includeCache = ! *_excludeCache_0
		_includeEtag = *_includeEtag_0
		_includeFolderListing = *_includeFolderListing_0
		_includePathsIndex = ! *_excludePathsIndex_0
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
			cdbWriteKeys : make (map[string]bool, 16 * 1024),
			storedFilePaths : make ([]string, 0, 16 * 1024),
			storedFolderPaths : make ([]string, 0, 16 * 1024),
			storedRedirectPaths : make ([]string, 0, 16 * 1024),
			storedDataMeta : make (map[string]bool, 16 * 1024),
			storedDataContent : make (map[string]bool, 16 * 1024),
			storedDataContentMeta : make (map[string]map[string]string, 16 * 1024),
			storedFiles : make (map[string][2]string, 16 * 1024),
			storedKeys : make (map[string]uint64, 16 * 1024),
			compress : _compress,
			compressLevel : _compressLevel,
			compressCache : _compressCacheDb,
			sourcesCache : _sourcesCacheDb,
			includeIndex : _includeIndex,
			includeStripped : _includeStripped,
			includeCache : _includeCache,
			includeEtag : _includeEtag,
			includeFolderListing : _includeFolderListing,
			includePathsIndex : _includePathsIndex,
			progress : _progress,
			debug : _debug,
		}
	
	_context.progressStarted = time.Now ()
	if _error := _context.cdbWriter.Put ([]byte (NamespaceSchemaVersion), []byte (CurrentSchemaVersion)); _error != nil {
		AbortError (_error, "[43228812]  failed writing archive!")
	}
	_context.cdbWriteKeys[NamespaceSchemaVersion] = true
	_context.cdbWriteCount += 1
	_context.cdbWriteKeySize += len (NamespaceSchemaVersion)
	_context.cdbWriteDataSize += len (CurrentSchemaVersion)
	
	if _, _error := walkPath (_context, _sourcesFolder, "/", filepath.Base (_sourcesFolder), nil, true, 0); _error != nil {
		AbortError (_error, "[b6a19ef4]  failed walking folder!")
	}
	
	if _includePathsIndex {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedFilePaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _key, _error := PrepareKeyToString (NamespaceFilesIndex, 1); _error == nil {
			if _found, _ := _context.cdbWriteKeys[_key]; _found {
				_error := fmt.Errorf ("[1d4dcde6]  duplicate key encountered:  `%s`", _key)
				AbortError (_error, "[a2d60ec1]  failed writing archive!")
			}
			if _error := _context.cdbWriter.Put ([]byte (_key), _buffer); _error != nil {
				AbortError (_error, "[1dbdde05]  failed writing archive!")
			}
			_context.cdbWriteKeys[_key] = true
			_context.cdbWriteCount += 1
			_context.cdbWriteKeySize += len (_key)
			_context.cdbWriteDataSize += len (_buffer)
		} else {
			AbortError (_error, "[b94abb82]  failed writing archive!")
		}
	}
	
	if _includePathsIndex {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedFolderPaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _key, _error := PrepareKeyToString (NamespaceFoldersIndex, 1); _error == nil {
			if _found, _ := _context.cdbWriteKeys[_key]; _found {
				_error := fmt.Errorf ("[c427b7f7]  duplicate key encountered:  `%s`", _key)
				AbortError (_error, "[651f521a]  failed writing archive!")
			}
			if _error := _context.cdbWriter.Put ([]byte (_key), _buffer); _error != nil {
				AbortError (_error, "[e2dd2de0]  failed writing archive!")
			}
			_context.cdbWriteKeys[_key] = true
			_context.cdbWriteCount += 1
			_context.cdbWriteKeySize += len (_key)
			_context.cdbWriteDataSize += len (_buffer)
		} else {
			AbortError (_error, "[0afc95f6]  failed writing archive!")
		}
	}
	
	if _includePathsIndex {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedRedirectPaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _key, _error := PrepareKeyToString (NamespaceRedirectsIndex, 1); _error == nil {
			if _found, _ := _context.cdbWriteKeys[_key]; _found {
				_error := fmt.Errorf ("[1d4dcde6]  duplicate key encountered:  `%s`", _key)
				AbortError (_error, "[d57123af]  failed writing archive!")
			}
			if _error := _context.cdbWriter.Put ([]byte (_key), _buffer); _error != nil {
				AbortError (_error, "[24833b01]  failed writing archive!")
			}
			_context.cdbWriteKeys[_key] = true
			_context.cdbWriteCount += 1
			_context.cdbWriteKeySize += len (_key)
			_context.cdbWriteDataSize += len (_buffer)
		} else {
			AbortError (_error, "[05b9d10b]  failed writing archive!")
		}
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
		log.Printf ("[ii] [56f63575]  completed    -- %0.2f min -- %d fil, %d fol, %d red, %0.2f M (%0.2f M/s) -- %d comp, %0.2f%% -- %d rec, %0.2f M\n",
				time.Since (_context.progressStarted) .Minutes (),
				len (_context.storedFilePaths),
				len (_context.storedFolderPaths),
				len (_context.storedRedirectPaths),
				float32 (_context.dataUncompressedSize) / 1024 / 1024,
				float64 (_context.dataUncompressedSize) / 1024 / 1024 / (time.Since (_context.progressStarted) .Seconds () + 0.001),
				_context.dataCompressedCount,
				(float32 (_context.dataUncompressedSize - _context.dataCompressedSize) / float32 (_context.dataUncompressedSize) * 100),
				_context.cdbWriteCount,
				float32 (_context.cdbWriteKeySize + _context.cdbWriteDataSize) / 1024 / 1024,
			)
	}
	
	return nil
}




func gobMarshal (_object interface{}) ([]byte, error) {
	_buffer := bytes.NewBuffer (nil)
	_encoder := gob.NewEncoder (_buffer)
	if _error := _encoder.Encode (_object); _error != nil {
		return nil, _error
	}
	return _buffer.Bytes (), nil
}

func gobUnmarshal (_data []byte, _object interface{}) (error) {
	_buffer := bytes.NewBuffer (_data)
	_encoder := gob.NewDecoder (_buffer)
	if _error := _encoder.Decode (_object); _error != nil {
		return _error
	}
	return nil
}




//go:embed usage.txt
var usageText string

//go:embed manual.txt
var manualText string

func init () {
	usageText = strings.ReplaceAll (usageText, "@{SCHEMA}", CurrentSchemaVersion)
}

