package restconnector

import (
	"path/filepath"
	"regexp"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/mogo/os/osutil"
)

type ConnectionSet struct {
	Map ConnectionFile
}

func NewConnectionSet() ConnectionSet {
	set := ConnectionSet{}
	set.init()
	return set
}

// ReadDir reads Connection JSON files.
func (set *ConnectionSet) ReadDir(dir string, recursive bool) error {
	set.init()
	connJSONEntries, err := osutil.ReadDirMore(dir, regexp.MustCompile(`^Connection_.+\.json$`), false, true, false)
	if err != nil {
		return err
	}

	for _, entry := range connJSONEntries {
		connMap, err := ReadConnectionFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
		for k, v := range *connMap {
			set.Map[k] = v
		}
	}

	if recursive {
		sdirs, err := osutil.ReadDirMore(dir, nil, true, false, false)
		if err != nil {
			return err
		}
		for _, sdir := range sdirs {
			sdirName := sdir.Name()
			if sdirName == "." || sdirName == ".." {
				continue
			} else if err := set.ReadDir(filepath.Join(dir, sdirName), recursive); err != nil {
				return err
			}
		}
	}
	return nil
}

func (set *ConnectionSet) ExtendedAttributes() *histogram.HistogramSet {
	hs := histogram.NewHistogramSet("extended attributes")
	for ck, c := range set.Map {
		eaNames := c.ExternalAttrs.Names(ExternalAttrNamesOpts{
			ToUpper:      true,
			Dedupe:       true,
			Sort:         true,
			RequireValue: true,
		})
		for _, eaName := range eaNames {
			hs.Add(ck, eaName, 1)
		}
	}
	return hs
}

func (set *ConnectionSet) ExtendedAttributesWriteXLSX(filename string) error {
	hs := set.ExtendedAttributes()
	return hs.WriteXLSXPivot(filename, "Attributes", "Connector", true, false, false, true)
}

func (set *ConnectionSet) init() {
	if set.Map == nil {
		set.Map = ConnectionFile{}
	}
}
