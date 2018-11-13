

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

// import "github.com/colinmarc/cdb"
import cdb "github.com/cipriancraciun/go-cdb-lib"

import . "github.com/cipriancraciun/go-cdb-http/lib/common"
import . "github.com/cipriancraciun/go-cdb-http/lib/archiver"




type context struct {
	cdbWriter *cdb.Writer
	storedData map[string]bool
	compress string
	debug bool
}




func archiveFile (_context *context, _pathResolved string, _pathInArchive string, _name string, _stat os.FileInfo) (error) {
	
	var _data []byte
	if _data_0, _error := ioutil.ReadFile (_pathResolved); _error == nil {
		_data = _data_0
	} else {
		return _error
	}
	
	if _, _error := archiveData (_context, NamespaceFilesContent, _pathInArchive, _name, _data, ""); _error != nil {
		return _error
	}
	
	return nil
}




func archiveFolder (_context *context, _pathResolved string, _pathInArchive string, _names []string, _stats map[string]os.FileInfo) (error) {
	
	type Entry struct {
		Name string `json:"name",omitempty`
		Type string `json:"type",omitempty`
		Size uint64 `json:"size",omitempty"`
	}
	
	type Folder struct {
		Entries []Entry `json:"entries"`
	}
	
	_entries := make ([]Entry, 0, len (_names))
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
	
	_folder := Folder {
			Entries : _entries,
		}
	
	if _data, _error := json.Marshal (&_folder); _error == nil {
		if _, _error := archiveData (_context, NamespaceFoldersEntries, _pathInArchive, "", _data, MimeTypeJson); _error != nil {
			return _error
		}
	} else {
		return _error
	}
	
	return nil
}




func archiveData (_context *context, _namespace string, _pathInArchive string, _name string, _data []byte, _dataType string) (string, error) {
	
	_fingerprintRaw := sha256.Sum256 (_data)
	_fingerprint := hex.EncodeToString (_fingerprintRaw[:])
	
	_wasStored, _ := _context.storedData[_fingerprint]
	
	if (_dataType == "") && (_name != "") {
		_extension := filepath.Ext (_pathInArchive)
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
	
	if ! _wasStored {
		
		var _dataEncoding string
		if _data_0, _dataEncoding_0, _error := Compress (_data, _context.compress); _error == nil {
			_data = _data_0
			_dataEncoding = _dataEncoding_0
		}
		
		_metadata := make (map[string]string, 16)
		_metadata["content-type"] = _dataType
		_metadata["content-encoding"] = _dataEncoding
		_metadata["etag"] = _fingerprint
		
		var _metadataRaw []byte
		if _metadataRaw_0, _error := MetadataEncode (_metadata); _error == nil {
			_metadataRaw = _metadataRaw_0
		} else {
			return "", _error
		}
		
		{
			_key := fmt.Sprintf ("%s:%s", NamespaceDataContent, _fingerprint)
			if _context.debug {
				log.Printf ("[  ] ++ %s", _key)
			}
			if _error := _context.cdbWriter.Put ([]byte (_key), _data); _error != nil {
				return "", _error
			}
		}
		
		{
			_key := fmt.Sprintf ("%s:%s", NamespaceDataMetadata, _fingerprint)
			if _context.debug {
				log.Printf ("[  ] ++ %s", _key)
			}
			if _error := _context.cdbWriter.Put ([]byte (_key), _metadataRaw); _error != nil {
				return "", _error
			}
		}
	}
	
	if _namespace != "" {
		_key := fmt.Sprintf ("%s:%s", _namespace, _pathInArchive)
		if _context.debug {
			log.Printf ("[  ] ++ %s %s", _key, _fingerprint)
		}
		if _error := _context.cdbWriter.Put ([]byte (_key), []byte (_fingerprint)); _error != nil {
			return "", _error
		}
	}
	
	return _dataType, nil
}




func walkPath (_context *context, _path string, _prefix string, _name string, _recursed map[string]uint) (error) {
	
	if _recursed == nil {
		_recursed = make (map[string]uint, 128)
	}
	
	_pathInArchive := filepath.Join (_prefix, _name)
	
	var _stat os.FileInfo
	if _stat_0, _error := os.Lstat (_path); _error == nil {
		_stat = _stat_0
	} else {
		return _error
	}
	
	_isSymlink := false
	if (_stat.Mode () & os.ModeSymlink) != 0 {
		_isSymlink = true
		if _stat_0, _error := os.Stat (_path); _error == nil {
			_stat = _stat_0
		} else {
			return _error
		}
	}
	
	var _pathResolved string
	if _isSymlink {
		if _resolved, _error := filepath.EvalSymlinks (_path); _error == nil {
			_pathResolved = _resolved
		} else {
			return _error
		}
	} else {
		_pathResolved = _path
	}
	
	if _isSymlink && _context.debug {
		log.Printf ("[  ] ~~ %s -> %s\n", _pathInArchive, _pathResolved)
	}
	
	if _stat.Mode () .IsRegular () {
		
		if _context.debug {
			log.Printf ("[  ] ## %s\n", _pathInArchive)
		}
		return archiveFile (_context, _pathResolved, _pathInArchive, _name, _stat)
		
	} else if _stat.Mode () .IsDir () {
		
		_wasRecursed, _ := _recursed[_pathResolved]
		if _wasRecursed > 0 {
			log.Printf ("[ww] [2e1744c9]  detected directory loop for `%s` resolving to `%s`;  ignoring!\n", _path, _pathResolved)
			return nil
		}
		_recursed[_pathResolved] = _wasRecursed + 1
		
		if _context.debug {
			log.Printf ("[  ] >> %s\n", _pathInArchive)
		}
		
		_names := make ([]string, 0, 16)
		_stats := make (map[string]os.FileInfo, 16)
		
		if _stream, _error := os.Open (_path); _error == nil {
			defer _stream.Close ()
			_prefix = filepath.Join (_prefix, _name)
			_loop : for {
				switch _buffer, _error := _stream.Readdir (128); _error {
					case nil :
						for _, _stat := range _buffer {
							_name := _stat.Name ()
							_names = append (_names, _name)
							_stats[_name] = _stat
							if _error := walkPath (_context, filepath.Join (_path, _name), _prefix, _name, _recursed); _error != nil {
								return _error
							}
						}
					case io.EOF :
						break _loop
					default :
						return _error
				}
			}
		}
		
		sort.Strings (_names)
		
		if _context.debug {
			log.Printf ("[  ] << %s\n", _pathInArchive)
		}
		
		if _context.debug {
			log.Printf ("[  ] <> %s\n", _pathInArchive)
		}
		if _error := archiveFolder (_context, _pathResolved, _pathInArchive, _names, _stats); _error != nil {
			return _error
		}
		
		_recursed[_pathResolved] = _wasRecursed
		return nil
		
	} else {
		return fmt.Errorf ("[d9b836d7]  unexpected file type for `%s`:  `%s`!", _path, _stat.Mode ())
	}
}




func main () () {
	Main (main_0)
}


func main_0 () (error) {
	
	
	var _sourcesFolder string
	var _archiveFile string
	var _compress string
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("cdb-http-archiver", flag.ContinueOnError)
		
		_sourcesFolder_0 := _flags.String ("sources", "", "<path>")
		_archiveFile_0 := _flags.String ("archive", "", "<path>")
		_compress_0 := _flags.String ("compress", "", "gzip | brotli")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_sourcesFolder = *_sourcesFolder_0
		_archiveFile = *_archiveFile_0
		_compress = *_compress_0
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
			compress : _compress,
			debug : _debug,
		}
	
	if _error := walkPath (_context, _sourcesFolder, "/", "", nil); _error != nil {
		AbortError (_error, "[b6a19ef4]  failed walking folder!")
	}
	
	if _error := _cdbWriter.Close (); _error != nil {
		AbortError (_error, "[bbfb8478]  failed creating archive (while closing)!")
	}
	
	
	return nil
}

