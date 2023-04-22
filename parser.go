package damonuntil

import (
	"fmt"
	"io"
	"os"
)

type Result struct {
	Version uint32 // 1 or 2
	Records []Record
}

type Record struct {
	TargetId  uint64
	Snapshots []Snapshot
}

type Snapshot struct {
	TargetId  uint64
	StartTime uint64
	EndTime   uint64
	Regions   []Region
}

type Region struct {
	StartAddr        uint64
	EndAddr          uint64
	NumberOfAccesses uint32
	Age              int
	AgeUnit          string
}

func ParseDamonFile(filepath string) (*Result, error) {
	fp, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	// 16 bytes for the header "damon_recfmt_ver"
	s := make([]byte, 16)
	n, err := io.ReadFull(fp, s)
	if err != nil {
		return nil, err
	}
	if n != 16 {
		return nil, fmt.Errorf("read %d bytes for header, expected 16", n)
	}
	if string(s) != "damon_recfmt_ver" {
		return nil, fmt.Errorf("invalid header: %s", string(s))
	}

	result := &Result{}

	// 4 bytes for the version
	v, err := readUint32(fp)
	if err != nil {
		return nil, fmt.Errorf("version error: %s", err.Error())
	}
	if v != 2 {
		return nil, fmt.Errorf("invalid version: %d (only version 2 is supported)", v)
	}
	result.Version = v

	startTime := uint64(0)

	// parse records
	for {
		// 8bytes for the end time [sec]
		// 8bytes for the end time [nsec]
		sec, err := readUint64(fp)
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, fmt.Errorf("end time sec error: %s", err.Error())
		}
		nsec, err := readUint64(fp)
		if err != nil {
			return result, fmt.Errorf("end time nsec error: %s", err.Error())
		}
		endTime := sec*1000000000 + nsec

		// 4 bytes for the number of results
		nr, err := readUint32(fp)
		if err != nil {
			return result, fmt.Errorf("number of results error: %s", err.Error())
		}

		r := Record{}
		for i := 0; i < int(nr); i++ {
			snapshot, err := parseSnapshot(fp)
			if err != nil {
				return result, err
			}
			snapshot.StartTime = startTime
			snapshot.EndTime = endTime
			r.Snapshots = append(r.Snapshots, snapshot)
		}

		result.Records = append(result.Records, r)

		startTime = endTime // for the next record
	}

	return result, nil
}

func parseSnapshot(buf io.Reader) (Snapshot, error) {
	s := Snapshot{}

	// 8 bytes for the target id
	// note: the target id of version 1 is 4 bytes
	targetId, err := readUint64(buf)
	if err != nil {
		return s, fmt.Errorf("target id error: %s", err.Error())
	}
	if targetId != 0 {
		return s, fmt.Errorf("invalid target id: %d (only 0 is supported)", targetId)
	}
	s.TargetId = targetId

	// 4 bytes for the number of regions
	nr, err := readUint32(buf)
	if err != nil {
		return s, fmt.Errorf("number of regions error: %s", err.Error())
	}

	for i := 0; i < int(nr); i++ {
		region, err := parseRegion(buf)
		if err != nil {
			return s, err
		}
		s.Regions = append(s.Regions, region)
	}

	return s, nil
}

func parseRegion(buf io.Reader) (Region, error) {
	r := Region{Age: -1}

	startAddr, err := readUint64(buf)
	if err != nil {
		return r, fmt.Errorf("start address error: %s", err.Error())
	}
	r.StartAddr = startAddr

	endAddr, err := readUint64(buf)
	if err != nil {
		return r, fmt.Errorf("end address error: %s", err.Error())
	}
	r.EndAddr = endAddr

	accesses, err := readUint32(buf)
	if err != nil {
		return r, fmt.Errorf("number of accesses error: %s", err.Error())
	}
	r.NumberOfAccesses = accesses

	return r, nil
}

func readUint32(buf io.Reader) (uint32, error) {
	s := make([]byte, 4)
	n, err := io.ReadFull(buf, s)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, fmt.Errorf("read %d bytes for uint32, expected 4", n)
	}
	return uint32(s[0]) | uint32(s[1])<<8 | uint32(s[2])<<16 | uint32(s[3])<<24, nil
}

func readUint64(buf io.Reader) (uint64, error) {
	s := make([]byte, 8)
	n, err := io.ReadFull(buf, s)
	if err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, fmt.Errorf("read %d bytes for uint64, expected 8", n)
	}
	return uint64(s[0]) | uint64(s[1])<<8 | uint64(s[2])<<16 | uint64(s[3])<<24 |
		uint64(s[4])<<32 | uint64(s[5])<<40 | uint64(s[6])<<48 | uint64(s[7])<<56, nil
}
