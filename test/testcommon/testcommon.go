package testcommon

import (
	"io"
	"testing"
	"bytes"
)

func IsMapsIntStringEqual(a map[int]string,
	                      b map[int]string,
	                      ) (adiff map[int]string,
	                         bdiff map[int]string,
	                         missMatch map[int][2]string,
	                         isEqual bool) {

	adiff = map[int]string{}
	bdiff = map[int]string{}
	missMatch = map[int][2]string{}

	for k := range a {
		adiff[k] = a[k]
	}

	for k := range b {
		bdiff[k] = b[k]
	}

	for k, aval := range adiff {
		bval, found := bdiff[k]
		if !found {
			continue
		}
		delete(adiff, k)
		delete(bdiff, k)

		if aval != bval {
			missMatch[k] = [2]string{aval, bval}
		}
	}

	return adiff,
	       bdiff,
	       missMatch,
	       (len(adiff) == 0 && len(bdiff) == 0 && len(missMatch) == 0)
}

func MapsIntStringMustBeEqual(t *testing.T, ref map[int]string, res map[int]string) {

	refDiff, resDiff, missMatch, isEqual := IsMapsIntStringEqual(ref, res)
	if isEqual {
		return
	}

	t.Log("Expected:", ref)
	t.Log("Returned:", res)
	if len(refDiff) != 0 {
		t.Log("Return lacks of:", refDiff)
	}
	if len(resDiff) != 0 {
		t.Log("Return should not contain:", resDiff)
	}
	if len(missMatch) != 0 {
		t.Log("Return key vallues differentiates:", missMatch)
	}

	t.Fail()
}

func ReaderContentMustBeEqual(t *testing.T, ref io.Reader, res io.Reader) {
	var refBuff []byte = make([]byte, 4096)
	var resBuff []byte = make([]byte, 4096)
	var refReadErr, resReadErr error

	for {
		_, refReadErr = ref.Read(refBuff)
		_, resReadErr = res.Read(resBuff)

		if refReadErr != nil || resReadErr != nil {
			break
		}

		if bytes.Compare(refBuff, resBuff) != 0 {
			t.Errorf("FAIL: expected: %s", string(refBuff))
			t.Errorf("FAIL: returned: %s", string(resBuff))
			t.Fail()
			return
		}
	}

	if refReadErr != io.EOF || resReadErr != io.EOF {
		if refReadErr != nil {
			t.Errorf("FAIL: expected ref EOF got: %s", refReadErr.Error())
		}
		if resReadErr != nil {
			t.Errorf("FAIL: expected res EOF got: %s", resReadErr.Error())
		}

		t.Fail()
		return
	}
}
