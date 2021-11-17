

package archiver


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
import "strings"
import "syscall"

import "github.com/colinmarc/cdb"

import . "github.com/volution/kawipiko/lib/common"
import . "github.com/volution/kawipiko/lib/archiver"




type context struct {
	cdbWriter *cdb.Writer
	storedFilePaths []string
	storedFolderPaths []string
	storedDataMeta map[string]bool
	storedDataContent map[string]bool
	storedDataContentMeta map[string]map[string]string
	storedFiles map[[2]uint64][2]string
	compress string
	includeIndex bool
	includeStripped bool
	includeCache bool
	includeEtag bool
	includeFileListing bool
	includeFolderListing bool
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
			_fileId = [2]uint64 { uint64 (_stat.Dev), uint64 (_stat.Ino) }
		} else {
			return fmt.Errorf ("[6578d2d7]  failed `stat`-ing:  `%s`!", _pathResolved)
		}
	}
	
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
		
		if _dataContent_0, _error := ioutil.ReadAll (_file); _error == nil {
			_dataContent = _dataContent_0
		} else {
			return _error
		}
		
		if _fingerprintContent_0, _dataContent_0, _dataMeta_0, _error := prepareDataContent (_context, _pathResolved, _pathInArchive, _name, _dataContent, ""); _error != nil {
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
	
	if _fingerprintContent_0, _dataContent_0, _dataMeta_0, _error := prepareDataContent (_context, _pathResolved, _pathInArchive, _name, _dataContent, _dataType); _error != nil {
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
		_key := fmt.Sprintf ("%s:%s", NamespaceDataContent, _fingerprintContent)
		if _context.debug {
			log.Printf ("[  ] data-content ++ `%s`\n", _key)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataContent); _error != nil {
			return _error
		}
	}
	
	_context.storedDataContent[_fingerprintContent] = true
	
	return nil
}


func archiveDataMeta (_context *context, _fingerprintMeta string, _dataMeta []byte) (error) {
	
	if _wasStored, _ := _context.storedDataMeta[_fingerprintMeta]; _wasStored {
		return fmt.Errorf ("[2918c4e2]  data meta already stored:  `%s`!", _fingerprintMeta)
	}
	
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataMetadata, _fingerprintMeta)
		if _context.debug {
			log.Printf ("[  ] data-meta    ++ `%s`\n", _key)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), _dataMeta); _error != nil {
			return _error
		}
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
	
	_key := fmt.Sprintf ("%s:%s", _namespace, _pathInArchive)
	if _context.debug {
		log.Printf ("[  ] reference    ++ `%s` :: `%s` -> `%s` ~ `%s`\n", _namespace, _pathInArchive, _fingerprintContent[:16], _fingerprintMeta[:16])
	}
	
	_fingerprints := fmt.Sprintf ("%s:%s", _fingerprintContent, _fingerprintMeta)
	if _error := _context.cdbWriter.Put ([]byte (_key), []byte (_fingerprints)); _error != nil {
		return _error
	}
	
	return nil
}




func prepareDataContent (_context *context, _pathResolved string, _pathInArchive string, _name string, _dataContent []byte, _dataType string) (string, []byte, map[string]string, error) {
	
	_fingerprintContentRaw := sha256.Sum256 (_dataContent)
	_fingerprintContent := hex.EncodeToString (_fingerprintContentRaw[:])
	
	if _wasStored, _ := _context.storedDataContent[_fingerprintContent]; _wasStored {
		_dataMeta := _context.storedDataContentMeta[_fingerprintContent]
		return _fingerprintContent, nil, _dataMeta, nil
	}
	
	if (_dataType == "") && (_name != "") {
		_extension := filepath.Ext (_name)
		if _extension != "" {
			_extension = _extension[1:]
		}
		_dataType, _ = MimeTypesByExtension[_extension]
	}
	if _dataType == "" {
		_dataType = http.DetectContentType (_dataContent)
	}
	if _dataType == "" {
		_dataType = MimeTypeRaw
	}
	
	_dataEncoding := "identity"
	_dataUncompressedSize := len (_dataContent)
	_dataSize := _dataUncompressedSize
	if _dataSize > 512 {
		if _dataContent_0, _dataEncoding_0, _error := Compress (_dataContent, _context.compress); _error == nil {
			if _dataEncoding_0 != "identity" {
				_dataCompressedSize := len (_dataContent_0)
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
					_dataContent = _dataContent_0
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
	
	_fingerprintMetaRaw := sha256.Sum256 (_dataMetaRaw)
	_fingerprintMeta := hex.EncodeToString (_fingerprintMetaRaw[:])
	
	if _wasStored, _ := _context.storedDataMeta[_fingerprintMeta]; _wasStored {
		return _fingerprintMeta, nil, nil
	}
	
	return _fingerprintMeta, _dataMetaRaw, nil
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
		log.Printf ("[  ] symlink      :: `%s` -> `%s`\n", _pathInArchive, _pathResolved)
	}
	
	
	if _statMode.IsRegular () {
		
		if _context.debug {
			log.Printf ("[  ] file         :: `%s`\n", _pathInArchive)
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
			log.Printf ("[  ] folder       >> `%s`\n", _pathInArchive)
		}
		if _stream, _error := os.Open (_pathResolved); _error == nil {
			defer _stream.Close ()
			_loop : for {
				switch _buffer, _error := _stream.Readdir (128); _error {
					case nil :
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
									log.Printf ("[  ] skip         !! `%s`\n", _childPathInArchive)
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
			log.Printf ("[  ] folder       << `%s`\n", _pathInArchive)
		}
		
		sort.Strings (_childsName)
		
		if _context.debug {
			log.Printf ("[  ] folder       :: `%s`\n", _pathInArchive)
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
			log.Printf ("[  ] folder       >> `%s`\n", _pathInArchive)
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
			log.Printf ("[  ] folder       << `%s`\n", _pathInArchive)
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
	var _archiveFile string
	var _compress string
	var _includeIndex bool
	var _includeStripped bool
	var _includeCache bool
	var _includeEtag bool
	var _includeFileListing bool
	var _includeFolderListing bool
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

    --exclude-index
    --exclude-strip
    --exclude-cache
    --include-etag

    --exclude-file-listing
    --include-folder-listing

    --debug

  ** for details see:
     https://github.com/volution/kawipiko#kawipiko-archiver

`)
		}
		
		_sourcesFolder_0 := _flags.String ("sources", "", "")
		_archiveFile_0 := _flags.String ("archive", "", "")
		_compress_0 := _flags.String ("compress", "", "")
		_excludeIndex_0 := _flags.Bool ("exclude-index", false, "")
		_excludeStripped_0 := _flags.Bool ("exclude-strip", false, "")
		_excludeCache_0 := _flags.Bool ("exclude-cache", false, "")
		_includeEtag_0 := _flags.Bool ("include-etag", false, "")
		_excludeFileListing_0 := _flags.Bool ("exclude-file-listing", false, "")
		_includeFolderListing_0 := _flags.Bool ("include-folder-listing", false, "")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_sourcesFolder = *_sourcesFolder_0
		_archiveFile = *_archiveFile_0
		_compress = *_compress_0
		_includeIndex = ! *_excludeIndex_0
		_includeStripped = ! *_excludeStripped_0
		_includeCache = ! *_excludeCache_0
		_includeEtag = *_includeEtag_0
		_includeFileListing = ! *_excludeFileListing_0
		_includeFolderListing = *_includeFolderListing_0
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
	
	if _error := _cdbWriter.Put ([]byte (NamespaceSchemaVersion), []byte (CurrentSchemaVersion)); _error != nil {
		AbortError (_error, "[43228812]  failed writing archive!")
	}
	
	_context := & context {
			cdbWriter : _cdbWriter,
			storedFilePaths : make ([]string, 0, 16 * 1024),
			storedFolderPaths : make ([]string, 0, 16 * 1024),
			storedDataMeta : make (map[string]bool, 16 * 1024),
			storedDataContent : make (map[string]bool, 16 * 1024),
			storedDataContentMeta : make (map[string]map[string]string, 16 * 1024),
			storedFiles : make (map[[2]uint64][2]string, 16 * 1024),
			compress : _compress,
			includeIndex : _includeIndex,
			includeStripped : _includeStripped,
			includeCache : _includeCache,
			includeEtag : _includeEtag,
			includeFileListing : _includeFileListing,
			includeFolderListing : _includeFolderListing,
			debug : _debug,
		}
	
	if _, _error := walkPath (_context, _sourcesFolder, "/", filepath.Base (_sourcesFolder), nil, true); _error != nil {
		AbortError (_error, "[b6a19ef4]  failed walking folder!")
	}
	
	if _includeFileListing {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedFilePaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _error := _cdbWriter.Put ([]byte (NamespaceFilesIndex), _buffer); _error != nil {
			AbortError (_error, "[1dbdde05]  failed writing archive!")
		}
	}
	
	if _includeFolderListing {
		_buffer := make ([]byte, 0, 1024 * 1024)
		for _, _path := range _context.storedFolderPaths {
			_buffer = append (_buffer, _path ...)
			_buffer = append (_buffer, '\n')
		}
		if _error := _cdbWriter.Put ([]byte (NamespaceFoldersIndex), _buffer); _error != nil {
			AbortError (_error, "[e2dd2de0]  failed writing archive!")
		}
	}
	
	if _error := _cdbWriter.Close (); _error != nil {
		AbortError (_error, "[bbfb8478]  failed creating archive (while closing)!")
	}
	
	return nil
}

