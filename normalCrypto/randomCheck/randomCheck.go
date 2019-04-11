package randomCheck

/*
#cgo LDFLAGS: -ldl
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#define NUMOFTESTS 19
typedef int* (*check_t)(char* b);
void *handle = NULL;
int loadSo(){
    handle = dlopen("/usr/local/bin/libtest.so", RTLD_LAZY);
    if(!handle){
        printf("dlopen - %s\n", dlerror());
        return -1;
    }
    return 0;
}

int* callRandomCheck(char* data)
{
    check_t checkRandom = (check_t) dlsym(handle, "checkRandomNumbers");
    if(!checkRandom)
    {
        return NULL;
    }
    int* result = checkRandom(data);
    return result;
}

int readCheckResult(int* result, int index)
{
    if(index <= NUMOFTESTS)
    {
        return result[index];
    }
    return -1;
}
*/
import "C"
import (
    "fmt"
    "github.com/pkg/errors"
)

const (
    MOD_NIST = iota
    MOD_GM
)

const NUMOFTESTS = 19
const MINRAMDOMDATA  = 1000000

var	testNames = [NUMOFTESTS+1]string{ " ", "Frequency", "BlockFrequency", "CumulativeSums", "Runs", "LongestRun", "Rank",
"FFT", "NonOverlappingTemplate", "OverlappingTemplate", "Universal", "ApproximateEntropy", "RandomExcursions",
"RandomExcursionsVariant", "Serial", "LinearComplexity", "RunsDistribution", "Poker", "BinaryDerivative", "AutoCorrelation"}

type RandomCheckInfo struct {
    RandomCheckType int
    RandomData     []byte
}

func DealRandomCheck(randomType int, randoms []byte) (string, error) {
    laodResult := 0
    if C.handle == nil {
        laodResult = int(C.loadSo())
    }
    fmt.Println(laodResult)
    if laodResult == 0 {
        StrBytes := StringFilter(randoms)
        if len(StrBytes) < MINRAMDOMDATA {
            return "", errors.New(fmt.Sprintf("Random data must larger than %d", MINRAMDOMDATA))
        }
        randomsString := string(StrBytes)
        point := C.callRandomCheck(C.CString(randomsString))
        resultStr := ""
        resultCount := 0
        if randomType == MOD_NIST {
            resultCount = NUMOFTESTS - 4
        }else {
            resultCount = NUMOFTESTS
        }
        for i := 1; i <= resultCount; i++ {
            result := C.readCheckResult(point, C.int(i))
            resultStr += testNames[i] + "=" + fmt.Sprintf("%d", int(result))
            if i != resultCount {
                resultStr += "&"
            }
        }
        return resultStr, nil
    }

    return "", errors.New("Load  randomCheck lib failed")
}

func StringFilter(orgString []byte) []byte {
    return orgString
    var resultBytes []byte
    for i := 0; i < len(orgString); i++ {
        if orgString[i] == byte(48) || orgString[i] == byte(49) {
            newValue := byte(int(orgString[i]) -48)
            resultBytes = append(resultBytes, newValue)
        }
    }
    return resultBytes
}
